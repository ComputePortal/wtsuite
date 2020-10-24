package prototypes

import (
	"../values"

	"../../context"
)

var WebAssemblyEnv *BuiltinPrototype = allocBuiltinPrototype()

func generateWebAssemblyEnvPrototype() bool {
	*WebAssemblyEnv = BuiltinPrototype{
		"WebAssemblyEnv", nil,
		map[string]BuiltinFunction{}, // TODO: io functions?
		NewConstructorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {

			if err := CheckInputs(WebAssemblyFS, args, ctx); err != nil {
				return nil, err
			}

			fs := args[0]

			// check all the necessary interface methods
			existsMember, err := fs.GetMember(stack, "exists", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			existsRes, err := existsMember.EvalFunction(stack, []values.Value{NewString(ctx)}, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if !existsRes.IsInstanceOf(Boolean) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Boolean from fs.exists(), got " + existsRes.TypeName())
				panic(err)
				return nil, err
			}

			openMember, err := fs.GetMember(stack, "open", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			openRes, err := openMember.EvalFunction(stack, []values.Value{NewString(ctx)}, ctx)
			if err != nil {
				return nil, err
			}

			if !openRes.IsInstanceOf(Int) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Int from fs.open(), got " + openRes.TypeName())
				panic(err)
				return nil, err
			}

			createMember, err := fs.GetMember(stack, "create", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			createRes, err := createMember.EvalFunction(stack, []values.Value{NewString(ctx)}, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if !createRes.IsInstanceOf(Int) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Int from fs.create(), got " + createRes.TypeName())
				panic(err)
				return nil, err
			}

			closeMember, err := fs.GetMember(stack, "close", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if err := closeMember.EvalMethod(stack, []values.Value{NewInt(ctx)}, ctx); err != nil {
				panic(err)
				return nil, err
			}

			readMember, err := fs.GetMember(stack, "read", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			readRes, err := readMember.EvalFunction(stack, []values.Value{NewInt(ctx), NewInt(ctx)}, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if !readRes.IsInstanceOf(Uint8Array) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Uint8Array from fs.read(), got " + readRes.TypeName())
				panic(err)
				return nil, err
			}

			writeMember, err := fs.GetMember(stack, "write", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if err := writeMember.EvalMethod(stack, []values.Value{NewInt(ctx), NewAltArray(Uint8Array, []values.Value{NewInt(ctx)}, ctx)}, ctx); err != nil {
				panic(err)
				return nil, err
			}

			seekMember, err := fs.GetMember(stack, "seek", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if err := seekMember.EvalMethod(stack, []values.Value{NewInt(ctx),
				NewInt(ctx)}, ctx); err != nil {
				panic(err)
				return nil, err
			}

			tellMember, err := fs.GetMember(stack, "tell", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			tellRes, err := tellMember.EvalFunction(stack, []values.Value{NewInt(ctx)}, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if !tellRes.IsInstanceOf(Int) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Int from fs.tell(), got " + tellRes.TypeName())
				panic(err)
				return nil, err
			}

			sizeMember, err := fs.GetMember(stack, "size", false, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			sizeRes, err := sizeMember.EvalFunction(stack, []values.Value{NewInt(ctx)}, ctx)
			if err != nil {
				panic(err)
				return nil, err
			}

			if !sizeRes.IsInstanceOf(Int) {
				errCtx := args[0].Context()
				err := errCtx.NewError("Error: expected Int from fs.size(), got " + sizeRes.TypeName())
				panic(err)
				return nil, err
			}

			return NewInstance(WebAssemblyEnv, ctx), nil
		}),
	}

	return true
}

var _WebAssemblyEnvOk = generateWebAssemblyEnvPrototype()
