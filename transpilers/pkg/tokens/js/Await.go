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

func (t *Await) evalInternal() (values.Value, error) {
	promise, err := t.expr.EvalExpression()
	if err != nil {
		return nil, err
	}

	// expecting a Promise
	if !prototypes.IsPromise(promise) {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected a Promise, got a " + promise.TypeName())
	}

	return promise.GetMember(".resolve", false, t.Context())
}

func (t *Await) EvalExpression() (values.Value, error) {
  res, err := t.evalInternal()
  if err != nil {
    return nil, err
  }

  if res == nil {
    errCtx := t.Context()
    return nil,  errCtx.NewError("Error: promise resolves to void")
  }

  return res, nil
}

func (t *Await) EvalStatement() error {
	res, err := t.evalInternal()
	if err != nil {
		return err
	}

  if res != nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: promise returns non-void (hint: use void)")
  }

  return nil
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
