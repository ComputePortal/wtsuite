package functions

import (
	"fmt"
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

type AnonFun struct {
	scope       tokens.Scope
	args        []string
	argDefaults []tokens.Token
	value       tokens.Token
	ctx         context.Context
}

func NewAnonFun(scope tokens.Scope, args []string, argDefaults []tokens.Token, value tokens.Token, ctx context.Context) *AnonFun {
	return &AnonFun{scope, args, argDefaults, value, ctx}
}

func (f *AnonFun) Dump(indent string) string {
	return indent + "AnonFun(" + strings.Join(f.args, ",") + ")\n"
}

func (f *AnonFun) Eval(scope tokens.Scope) (tokens.Token, error) {
	return f, nil
}

func (f *AnonFun) Context() context.Context {
	return f.ctx
}

func (a *AnonFun) IsSame(other tokens.Token) bool {
	if b, ok := other.(*AnonFun); ok {
		if len(a.args) == len(b.args) {
			if !a.value.IsSame(b.value) {
				return false
			}
		}
	}

	return false
}

func (f *AnonFun) EvalFun(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) > len(f.args) {
		err := ctx.NewError(fmt.Sprintf("Error: too many arguments (expected %d, got %d)", len(f.args), len(args)))
		err.AppendContextString("Info: function defined here", f.Context())
		return nil, err
	}

	lambdaScope := NewLambdaScope(f.scope, scope)

	for i, arg := range args {
		// eval the incoming args too
		evaluatedArg, err := arg.Eval(scope)
		if err != nil {
			return nil, err
		}

		v := Var{evaluatedArg, false, true, false, false, arg.Context()}
		lambdaScope.SetVar(f.args[i], v)
	}

	// apply defaults for remainder
	for i := len(args); i < len(f.argDefaults); i++ {
		argDefault := f.argDefaults[i]
		if argDefault != nil {
			evaluatedArg, err := argDefault.Eval(lambdaScope)
			if err != nil {
				return nil, err
			}

			v := Var{evaluatedArg, false, true, false, false, argDefault.Context()}
			lambdaScope.SetVar(f.args[i], v)
		}

		// TODO: should we give an error if default is not available?
	}

	result, err := f.value.Eval(lambdaScope)
	if err != nil {
		context.AppendContextString(err, "Info: called here", ctx)
		return nil, err
	}

	return result, nil
}

func (f *AnonFun) Len() int {
	return len(f.args)
}

func IsAnonFun(t tokens.Token) bool {
	_, ok := t.(*AnonFun)
	return ok
}
