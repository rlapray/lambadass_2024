package sql //nolint:revive,max-public-structs // Many public struct, we know

import (
	"context"
	"reflect"

	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ExecMapKey struct {
	Q string
	D any
}

type SelectMapKey struct {
	Q  string
	DA any
}

type SelectMapValue struct {
	D any
	F fault.Fault
}

type ExecMapValue struct {
	RA int64
	F  fault.Fault
}

type ExecOneRowAffectedMapValue struct {
	F fault.Fault
}

type MockClient[T any, U any] struct {
	logger                       *zerolog.Logger
	execMap                      map[ExecMapKey]ExecMapValue
	execMapCounter               map[ExecMapKey]int
	ExecOneRowAffectedMap        map[ExecMapKey]ExecOneRowAffectedMapValue
	ExecOneRowAffectedMapCounter map[ExecMapKey]int
	selectMap                    map[SelectMapKey]SelectMapValue
	selectMapCounter             map[SelectMapKey]int
}

/******************************************************************************
***** Exec
******************************************************************************/
func (m *MockClient[T, U]) Exec(query string, data any) (int64, fault.Fault) {
	key := ExecMapKey{Q: query, D: data}
	if res, exists := m.execMap[key]; exists {
		if counter, exists := m.execMapCounter[key]; exists {
			m.execMapCounter[key] = counter + 1
		} else {
			m.execMapCounter[key] = 1
		}
		return res.RA, res.F
	}
	return 0, fault.NewSQL(m.logger, "MOCK_DATA_NOT_FOUND", "Mock data not found", map[string]any{"key": key}, nil)
}

func (m *MockClient[T, U]) MockExecMap(key ExecMapKey, value ExecMapValue) {
	if m.execMap == nil {
		m.execMap = make(map[ExecMapKey]ExecMapValue)
	}
	m.execMap[key] = value
}

/******************************************************************************
***** ExecOneRowAffected
******************************************************************************/

func (m *MockClient[T, U]) ExecOneRowAffected(query string, data any) fault.Fault {
	key := ExecMapKey{Q: query, D: data}
	if res, exists := m.ExecOneRowAffectedMap[key]; exists {
		if counter, exists := m.ExecOneRowAffectedMapCounter[key]; exists {
			m.ExecOneRowAffectedMapCounter[key] = counter + 1
		} else {
			m.ExecOneRowAffectedMapCounter[key] = 1
		}
		return res.F
	}
	return fault.NewSQL(m.logger, "MOCK_DATA_NOT_FOUND", "Mock data not found", map[string]any{"key": key}, nil)
}

func (m *MockClient[T, U]) MockExecOneRowAffectedMap(key ExecMapKey, value ExecOneRowAffectedMapValue) {
	if m.execMap == nil {
		m.ExecOneRowAffectedMap = make(map[ExecMapKey]ExecOneRowAffectedMapValue)
		m.ExecOneRowAffectedMapCounter = make(map[ExecMapKey]int)
	}
	m.ExecOneRowAffectedMap[key] = value
}

/******************************************************************************
***** Select
******************************************************************************/

func (m *MockClient[T, U]) Select(query string, data, destination any) fault.Fault {
	key := SelectMapKey{Q: query, DA: data}
	if res, exists := m.selectMap[key]; exists {
		if counter, exists := m.selectMapCounter[key]; exists {
			m.selectMapCounter[key] = counter + 1
		} else {
			m.selectMapCounter[key] = 1
		}

		val := reflect.ValueOf(destination)
		if val.Kind() == reflect.Ptr {
			newVal := reflect.ValueOf(res.D)
			if newVal.Elem().Kind() == reflect.Interface {
				val.Elem().Set(newVal.Elem().Elem())
			} else {
				val.Elem().Set(newVal.Elem()) // Not yet tested
			}
		}
		return res.F
	}
	return fault.NewSQL(m.logger, "MOCK_DATA_NOT_FOUND", "Mock data not found", map[string]any{"key": key}, nil)
}

func (m *MockClient[T, U]) MockSelectMap(key SelectMapKey, flt fault.Fault, destination any) {
	if m.selectMap == nil {
		m.selectMap = make(map[SelectMapKey]SelectMapValue)
		m.selectMapCounter = make(map[SelectMapKey]int)
	}
	m.selectMap[key] = SelectMapValue{F: flt, D: &destination}
}

/******************************************************************************
***** Middleware
******************************************************************************/
func (m *MockClient[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	ll := log.Logger.With().Str("framework", "SQL").Logger()
	m.logger = &ll
	return nil
}

func (*MockClient[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	return nil
}

func (*MockClient[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	return err
}

func (*MockClient[T, U]) OnShutdown() {
}
