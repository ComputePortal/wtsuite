package macros

import (
	"../prototypes"
	"../values"

	"../../context"
)

type ToInstance struct {
	protos []values.Prototype // TODO: collect these prototypes during ResolveNames stage, not during EvalTypes stage
}

func newToInstance() ToInstance {
	return ToInstance{make([]values.Prototype, 0)}
}

func (m *ToInstance) evalInstancePrototype(stack values.Stack, protoArg values.Value, ctx context.Context) (values.Value, error) {
	classValues, ok := protoArg.LiteralArrayValues()
	if !ok {
		if !protoArg.IsInstanceOf(prototypes.Array) {
			return nil, ctx.NewError("Error: argument 2 not an array")
		}
		// try the regular array items
		mv, err := protoArg.GetIndex(stack, prototypes.NewLiteralInt(1, ctx), ctx)
		if err != nil {
			return nil, err
		}

		vs := values.UnpackMulti([]values.Value{mv})
		classValues = vs

		//err := ctx.NewError("Error: argument 2 is not a literal array")
		//return nil, err
	}

	if len(classValues) == 0 {
		err := ctx.NewError("Error: expected at least one entry in class list")
		return nil, err
	}

	res := make([]values.Value, len(classValues))

	for i, classValue := range classValues {
		proto, ok := classValue.GetClassPrototype()
		if !ok {
			return nil, ctx.NewError("Error: expected only classes in array, got " + classValue.TypeName())
		}

		if !proto.IsUniversal() {
			errCtx := ctx
			return nil, errCtx.NewError("Error: class " + proto.Name() + " is not universal (hint: use 'universe')")
		}

		inst, err := proto.GenerateInstance(stack, nil, nil, classValue.Context())
		if err != nil {
			context.AppendContextString(err, "Info: needed in this macro", ctx)
			return nil, err
		}

		res[i] = inst

		var univErr error = nil
		inst.LoopNestedPrototypes(func(nestedProto values.Prototype) {
			if univErr == nil && !proto.IsUniversal() {
				errCtx := ctx
				univErr = errCtx.NewError("Error: class " + proto.Name() + " is not universal (hint: use 'universe')")
			}
		})

		if univErr != nil {
			return nil, univErr
		}
	}

	return values.NewMulti(res, ctx), nil
}
