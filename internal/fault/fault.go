// Package fault contains the [fault/Fault] interface and all others Faults struct.
// A [fault/Fault] represent an error in the context of this architecture.
// It contains much more information really useful for debugging and limiting complexity through this project
package fault

type Layer int

var layersString = []string{
	"Other",
	"Frameworks",
	"Adapters",
	"UseCases",
	"Commands",
}

func (l Layer) String() string {
	return layersString[l]
}

const (
	Other Layer = iota
	Frameworks
	Adapters
	UseCases
	Commands
)

type Fault interface {
	Code() string
	Layer() Layer
	Middleware() string
	Message() string
	Metadata() map[string]any
	Cause() error
	Error() string
}
