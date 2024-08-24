package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // For the database driver
	"github.com/jmoiron/sqlx"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/* ****************************************************************************
***** Structs
******************************************************************************/

type GenericClient[T any, U any] struct {
	database        *sqlx.DB
	mainTransaction *sqlx.Tx
	logger          *zerolog.Logger
}

type Client[T any, U any] interface {
	Exec(query string, data any) (int64, fault.Fault)
	ExecOneRowAffected(query string, data any) fault.Fault
	Select(query string, data, destination any) fault.Fault

	OnSetup(ctx context.Context, firstRequest *T) fault.Fault
	OnBefore(ctx context.Context, request *T) fault.Fault
	OnAfter(response *U, flt fault.Fault) fault.Fault
	OnShutdown()
}

/* ****************************************************************************
***** Functions
******************************************************************************/

// Execute an SQL query and return how many rows were affected. Useful for INSERT, UPDATE or DELETE queries.
func (m *GenericClient[T, U]) Exec(query string, data any) (int64, fault.Fault) {
	var sqlRes sql.Result
	var stmt *sqlx.NamedStmt
	var err error
	metadata := make(map[string]any)
	var rowAffected int64
	ll := m.logger.With().Str("query", query).Logger()
	logger := &ll

	logger.Debug().Msg("Executing SQL...")

	dur, _ := m.duration(func() {
		stmt, err = m.mainTransaction.PrepareNamed(query)
	})
	if err != nil {
		metadata["duration"] = dur
		return 0, fault.NewSQL(logger, "PREPARED_STATEMENT_FAILED", "Prepared statement cannot be created", metadata, err)
	}

	dur, durStr := m.duration(func() {
		sqlRes, err = stmt.Exec(data)
	})
	if err != nil {
		metadata["duration"] = dur
		return 0, fault.NewSQL(logger, GetPGError(err), "Error while executing SQL", metadata, err)
	}

	logger.Debug().Int64("duration", dur).Msgf("Executed SQL in %v", durStr)

	dur, _ = m.duration(func() {
		rowAffected, err = sqlRes.RowsAffected()
	})
	if err != nil {
		metadata["duration"] = dur
		return 0, fault.NewSQL(logger, "ROW_AFFECTED_UNKNOWN", "Cannot get the number of row affected", metadata, err)
	}
	return rowAffected, nil
}

// Like Exec, but will fail if not exactly 1 row is affected. Useful for INSERTs.
func (m *GenericClient[T, U]) ExecOneRowAffected(query string, data any) fault.Fault {
	ll := m.logger.With().Str("query", query).Logger()
	logger := &ll
	rowAffected, err := m.Exec(query, data)
	if err != nil {
		return err
	}
	if rowAffected != 1 {
		return fault.NewSQL(logger, "ROW_AFFECTED_NOT_ONE", "The number of row affected is not 1", nil, err)
	}
	return nil
}

func (m *GenericClient[T, U]) Select(query string, data, destination any) fault.Fault {
	var stmt *sqlx.NamedStmt
	var err error
	metadata := make(map[string]any)
	ll := m.logger.With().Str("query", query).Logger()
	logger := &ll

	logger.Debug().Msg("Executing SQL...")

	dur, _ := m.duration(func() {
		stmt, err = m.mainTransaction.PrepareNamed(query)
	})
	if err != nil {
		metadata["duration"] = dur
		return fault.NewSQL(logger, "PREPARED_STATEMENT_FAILED", "Prepared statement cannot be created", metadata, err)
	}

	dur, durStr := m.duration(func() {
		err = stmt.Select(destination, data)
	})
	if err != nil {
		metadata["duration"] = dur
		return fault.NewSQL(logger, GetPGError(err), "Error while executing SQL", metadata, err)
	}

	logger.Debug().Int64("duration", dur).Msgf("Executed SQL in %v", durStr)

	if err != nil {
		metadata["duration"] = dur
		return fault.NewSQL(logger, "ROW_AFFECTED_UNKNOWN", "Cannot get the number of row affected", metadata, err)
	}
	return nil
}

/******************************************************************************
***** Middleware
******************************************************************************/

func (m *GenericClient[T, U]) actualSetup(ctx context.Context) fault.Fault {
	user := os.Getenv("SQL_USER")
	pwd := os.Getenv("SQL_PASSWORD")
	host := os.Getenv("SQL_HOST")
	port := os.Getenv("SQL_PORT")
	database := os.Getenv("SQL_DATABASE")
	connectionMaxIdleTime, err := time.ParseDuration(os.Getenv("SQL_CONNECTION_MAX_IDLE_TIME"))
	if err != nil {
		return fault.NewSQL(m.logger, "NO_SQL_CONNECTION_MAX_IDLE_TIME", "Cannot parse duration from SQL_CONNECTION_MAX_IDLE_TIME", nil, err)
	}
	connectionMaxLifeTime, err := time.ParseDuration(os.Getenv("SQL_CONNECTION_MAX_LIFE_TIME"))
	if err != nil {
		return fault.NewSQL(m.logger, "NO_SQL_CONNECTION_MAX_LIFE_TIME", "cannot parse duration from SQL_CONNECTION_MAX_LIFE_TIME", nil, err)
	}

	db, err := sqlx.ConnectContext(ctx, "pgx", fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, pwd, host, port, database))
	if err != nil {
		return fault.NewSQL(m.logger, "SQL_CONNECTION_ERROR", "Cannot connect to the database", nil, err)
	}
	m.database = db
	m.database.SetConnMaxIdleTime(connectionMaxIdleTime)
	m.database.SetConnMaxLifetime(connectionMaxLifeTime)
	m.database.SetMaxIdleConns(1)
	m.database.SetMaxOpenConns(1)
	return nil
}

func (m *GenericClient[T, U]) OnSetup(ctx context.Context, _ *T) fault.Fault {
	ll := log.Logger.With().Str("framework", "SQL").Logger()
	m.logger = &ll
	m.logger.Trace().Msg("OnSetup")

	var closureErr fault.Fault

	m.logger.Debug().Msg("Connecting to the database...")
	dur, durStr := m.duration(func() {
		closureErr = m.actualSetup(ctx)
	})

	if closureErr != nil {
		m.logger.Warn().
			Int64("duration", dur).
			Msgf("Cannot connect to the database at the moment (%v)", durStr)
	} else {
		m.logger.Debug().Msgf("Connected to the database in %v", durStr)
	}
	return closureErr
}

func (m *GenericClient[T, U]) OnBefore(ctx context.Context, _ *T) fault.Fault {
	m.logger.Trace().Msg("OnBefore")
	txx, err := m.database.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return fault.NewSQL(m.logger, "NEW_TRANSACTION_ERROR", "Cannot create new transaction", nil, err)
	}
	m.mainTransaction = txx
	return nil
}

func (m *GenericClient[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	m.logger.Trace().Err(err).Msg("OnAfter")
	if err != nil {
		m.logger.Warn().Err(err).Msg("Rollback main transaction")
		if m.mainTransaction != nil {
			err2 := m.mainTransaction.Rollback()
			if err2 != nil {
				err = fault.NewSQL(m.logger, "SQL_ROLLBACK_ERROR", "Rollbacking main transaction raised an error", nil, err2)
			}
		} else {
			err = fault.NewSQL(m.logger, "SQL_ROLLBACK_NIL_TRANSACTION", "Rollbacking main transaction is impossible because it's nil", nil, err)
		}
	} else {
		m.logger.Info().Msg("Commit main transaction")
		err2 := m.mainTransaction.Commit()
		if err2 != nil {
			err = fault.NewSQL(m.logger, "SQL_COMMIT_ERROR", "Commit raised an error", nil, err2)
		}
	}
	return err
}

func (m *GenericClient[T, U]) OnShutdown() {
	m.logger.Trace().Msg("OnShutdown")
	m.logger.Debug().Msg("Closing database connections...")
	dur, durStr := m.duration(func() {
		m.database.Close()
	})
	m.logger.Debug().Int64("duration", dur).Msgf("Closed database connections in %v", durStr)
}

func (*GenericClient[T, U]) duration(f func()) (durationInt64 int64, durationObj string) {
	start := time.Now().UnixMilli()
	f()
	durationInt64 = time.Now().UnixMilli() - start
	durationObj = (time.Duration(durationInt64) * time.Millisecond).String()
	return durationInt64, durationObj
}
