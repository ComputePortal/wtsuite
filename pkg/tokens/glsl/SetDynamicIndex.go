package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type SetDynamicIndex struct {
  fnVar Variable // variable of auto-generated function
  arg Expression
  Index
}

func NewSetDynamicIndex(args []Expression, ctx context.Context) (Statement, error) {
  if len(args) != 3 {
    return nil, ctx.NewError("Error: expected 3 args")
  }

  return &SetDynamicIndex{nil, args[2], newIndex(args[0], args[1], ctx)}, nil
}

func (t *SetDynamicIndex) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("SetDynamicIndex\n")

	b.WriteString(t.container.Dump(indent + "  "))
	b.WriteString(t.index.Dump(indent + "[ "))
  b.WriteString(t.arg.Dump(indent + "= "))

	return b.String()
}

func (t *SetDynamicIndex) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

  b.WriteString(indent)
	b.WriteString(t.fnVar.Name())
	b.WriteString("(")
	b.WriteString(t.container.WriteExpression())
  b.WriteString(",")
	b.WriteString(t.index.WriteExpression())
  b.WriteString(",")
  b.WriteString(t.arg.WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (t *SetDynamicIndex) ResolveStatementNames(scope Scope) error {
  if err := t.Index.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if err := t.arg.ResolveExpressionNames(scope); err != nil {
    return err
  }

  return nil
}

func (t *SetDynamicIndex) EvalStatement() error {
	containerValue, err := t.container.EvalExpression()
	if err != nil {
		return err
	}

	indexValue, err := t.index.EvalExpression()
	if err != nil {
		return err
	}

  if _, err := values.AssertInt(indexValue); err != nil {
    return err
  }

  argVal, err := t.arg.EvalExpression()
  if err != nil {
    return err
  }

  dummyIndex := values.NewLiteralInt(0, t.index.Context())

  return containerValue.SetIndex(dummyIndex, argVal, t.Context())
}

func (t *SetDynamicIndex) ResolveStatementActivity(usage Usage) error {
  if err := t.Index.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  if err := t.arg.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  containerValue, err := t.container.EvalExpression()
  if err != nil {
    return err
  }

  dummyIndex := values.NewLiteralInt(0, t.index.Context())

  contentVal, err := containerValue.GetIndex(dummyIndex, t.Context())
  if err != nil {
    return err
  }

  contentTypeName := contentVal.TypeName()

  length := containerValue.Length()

  fn := NewSetDynamicIndexFunction(containerValue.TypeName(), contentTypeName, length, t.Context())

  t.fnVar, err = injectDynamicIndexStatement(usage, fn.GetVariable(), fn, contentVal)
  return err
}

func (t *SetDynamicIndex) UniqueStatementNames(ns Namespace) error {
  return nil
}
