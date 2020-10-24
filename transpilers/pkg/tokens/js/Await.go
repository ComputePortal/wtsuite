package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type Await struct {
	expr Expression
	TokenData
}

func NewAwait(expr Expression, ctx context.Context) (*Await, error) {
	return &Await{expr, TokenData{ctx}}, nil
}

func (t *Await) Args() []Token {
	return []Token{t.expr}
}

func (t *Await) Precedence() int {
	return _preUnaryPrecedenceMap["await"]
}

func (t *Await) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Await\n")

	b.WriteString(t.expr.Dump(indent + "  "))

	return b.String()
}

func (t *Await) WriteExpression() string {
	var b strings.Builder

	b.WriteString("await ")
	b.WriteString(t.expr.WriteExpression())

	return b.String()
}

func (t *Await) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.WriteExpression())

	return b.String()
}

func (t *Await) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Await) HoistNames(scope Scope) error {
	return nil
}

func (t *Await) ResolveExpressionNames(scope Scope) error {
	if !scope.IsAsync() {
		errCtx := t.Context()
		return errCtx.NewError("Error: await not in async scope")
	}

	return t.expr.ResolveExpressionNames(scope)
}

func (t *Await) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *Await) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Await) evalInternal(stack values.Stack) (values.Value, error) {
	exprValue, err := t.expr.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	// expecting a Promise
	if !exprValue.IsInstanceOf(prototypes.Promise) {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected a Promise, got a " + exprValue.TypeName())
	}

	return exprValue, nil
}

func (t *Await) EvalStatement(stack values.Stack) error {
	if res, ok, err := stack.ResolveAwait(t); ok {
		if err != nil {
			return err
		}

		if res != nil {
			errCtx := t.Context()
			return errCtx.NewError("Error: unexpected return value (hint: use void)")
		}
		return nil
	} else if err != nil {
		return err
	}

	exprValue, err := t.evalInternal(stack)
	if err != nil {
		return err
	}

	fn, err := exprValue.GetMember(stack, ".awaitMethod", false, t.Context())
	if err != nil {
		panic(err)
	}

	err = fn.EvalMethod(stack, []values.Value{}, t.Context())
	if err != nil {
		if ar, ok := err.(*prototypes.AsyncRequest); ok {
			ar.SetAwait(t)
			return ar
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (t *Await) EvalExpression(stack values.Stack) (values.Value, error) {
	if res, ok, err := stack.ResolveAwait(t); ok {
		if err != nil {
			return nil, err
		}

		if res == nil {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: expected a return value, got nothing")
		}
		return res, nil
	} else if err != nil {
		return nil, err
	}

	exprValue, err := t.evalInternal(stack)
	if err != nil {
		return nil, err
	}

	// if underlying object is AllNull (eg. uninitialized global object), then we cant be sure of anything
	if values.IsAllNull(exprValue) {
		return exprValue, nil
	}

	//if exprValue.IsNull() {
	//errCtx := t.Context()
	//panic(errCtx.NewError("Error: returned Promise can't be null"))
	//}

	fn, err := exprValue.GetMember(stack, ".awaitFunction", false, t.Context())
	if err != nil {
		panic(err)
	}

	if values.IsAllNull(fn) {
		errCtx := t.Context()
		panic(errCtx.NewError("Error: .awaitFunction can't be null"))
	}

	res, err := fn.EvalFunction(stack, []values.Value{}, t.Context())
	if err != nil {
		if ar, ok := err.(*prototypes.AsyncRequest); ok {
			ar.SetAwait(t)
			return nil, ar
		} else {
			return nil, err
		}
	} else {
		return res, nil
	}
}

func (t *Await) ResolveExpressionActivity(usage Usage) error {
	return t.expr.ResolveExpressionActivity(usage)
}
func (t *Await) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *Await) UniversalExpressionNames(ns Namespace) error {
	return t.expr.UniversalExpressionNames(ns)
}

func (t *Await) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *Await) UniqueExpressionNames(ns Namespace) error {
	return t.expr.UniqueExpressionNames(ns)
}

func (t *Await) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *Await) Walk(fn WalkFunc) error {
  if err := t.expr.Walk(fn); err != nil {
    return err
  }
  
  if err := fn(t); err != nil {
    return err
  }

  return nil
}
