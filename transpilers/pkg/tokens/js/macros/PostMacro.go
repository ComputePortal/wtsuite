package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type PostMacro struct {
	ToInstance
	Macro
}

func newPostMacro(args []js.Expression, ctx context.Context) PostMacro {
	return PostMacro{newToInstance(), newMacro(args, ctx)}
}

func (m *PostMacro) writeExpression(fnName string) string {
	var b strings.Builder

	b.WriteString(fnName)
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(",")
	b.WriteString(m.args[1].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *PostMacro) evalExpression(stack values.Stack, msg values.Value,
	classValue values.Value) (values.Value, error) {
	if !isAnObject(msg) {
		errCtx := m.Context()
		return nil,
			errCtx.NewError("Error: expected Object or instance of class that extends Object for argument 2, got " +
				msg.TypeName())
	}

	classProto, ok := classValue.GetClassPrototype()
	if !ok {
		errCtx := m.Context()
		return nil, errCtx.NewError("Error: argument 3 is not a class, got instance of " + classValue.TypeName())
	}

	resolveValue, err := classProto.GenerateInstance(stack, nil, nil, classValue.Context())
	if err != nil {
		context.AppendContextString(err, "Info: needed here", m.Context())
		return nil, err
	}

	rejectValue := prototypes.NewInstance(prototypes.Error, m.Context())

	var univErr error = nil
	resolveValue.LoopNestedPrototypes(func(proto values.Prototype) {
		if univErr == nil && !proto.IsUniversal() {
			errCtx := m.Context()
			univErr = errCtx.NewError("Error: class " + proto.Name() + " is not universal (hint: use 'universe'")
		}
	})

	if univErr != nil {
		return nil, univErr
	}

	promiseProps := values.NewPromiseProperties(m.Context())
	if err := promiseProps.SetResolveArgs([]values.Value{resolveValue}, m.Context()); err != nil {
		return nil, err
	}
	promiseProps.SetRejectArgs([]values.Value{rejectValue})

	return values.NewInstance(prototypes.Promise, promiseProps, m.Context()), nil
}
