package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Index struct {
  container Expression
  index Expression
  TokenData
}

func newIndex(container Expression, index Expression, ctx context.Context) Index {
  return Index{container, index, newTokenData(ctx)}
}
func NewIndex(container Expression, index Expression, ctx context.Context) *Index {
  idx := newIndex(container, index, ctx)
  
	return &idx
}

func (t *Index) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Index\n")

	b.WriteString(t.container.Dump(indent + "  "))
	b.WriteString(t.index.Dump(indent + "[ "))

	return b.String()
}

func (t *Index) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.container.WriteExpression())
	b.WriteString("[")
	b.WriteString(t.index.WriteExpression())
	b.WriteString("]")

	return b.String()
}

func (t *Index) ResolveExpressionNames(scope Scope) error {
	if err := t.container.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.index.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *Index) EvalExpression() (values.Value, error) {
	containerValue, err := t.container.EvalExpression()
	if err != nil {
		return nil, err
	}

	indexValue_, err := t.index.EvalExpression()
	if err != nil {
		return nil, err
	}

  indexValue, err := values.AssertLiteralInt(indexValue_)
  if err != nil {
    return nil, err
  }

  return containerValue.GetIndex(indexValue, t.Context())
}

func (t *Index) EvalSet(rhsValue values.Value, ctx context.Context) error {
	containerValue, err := t.container.EvalExpression()
	if err != nil {
		return err
	}

	indexValue_, err := t.index.EvalExpression()
	if err != nil {
		return err
	}

  indexValue, err := values.AssertLiteralInt(indexValue_)
  if err != nil {
    context.AppendString(err, "Hint: use setIndex(arr, i, x)\n")
    return err
  }

  return containerValue.SetIndex(indexValue, rhsValue, t.Context())
}

func (t *Index) ResolveExpressionActivity(usage Usage) error {
  if err := t.index.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  if err := t.container.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  return nil
}
