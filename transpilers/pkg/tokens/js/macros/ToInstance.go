package macros

import (
	"../values"
)

type ToInstance struct {
	protos []values.Prototype // TODO: collect these prototypes during ResolveNames stage, not during EvalTypes stage, for safer ToInstance in server code
}

func newToInstance() ToInstance {
	return ToInstance{make([]values.Prototype, 0)}
}
