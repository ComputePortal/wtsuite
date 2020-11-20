package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type LiteralArray struct {
	items []Expression
	LiteralData
}

func NewLiteralArray(items []Expression, ctx context.Context) *LiteralArray {
	return &LiteralArray{items, newLiteralData(ctx)}
}

func (t *LiteralArray) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("LiteralArray\n")

	for _, item := range t.items {
		b.WriteString(item.Dump(indent + "| "))
	}

	return b.String()
}

func (t *LiteralArray) WriteExpression() string {
	var b strings.Builder

	b.WriteString("[")

	for i, item := range t.items {
		b.WriteString(item.WriteExpression())

		if i < len(t.items)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String()
}

func (t *LiteralArray) ResolveExpressionNames(scope Scope) error {
	for _, item := range t.items {
		if err := item.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralArray) EvalExpression() (values.Value, error) {
	items := make([]values.Value, len(t.items))

	for i, itemExpr := range t.items {
		item, err := itemExpr.EvalExpression()
		if err != nil {
			return nil, err
		}

		items[i] = item
	}

  common := values.CommonValue(items, t.Context())

	return prototypes.NewArray(common, t.Context()), nil
}

func (t *LiteralArray) ResolveExpressionActivity(usage Usage) error {
	for _, item := range t.items {
		if err := item.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralArray) UniversalExpressionNames(ns Namespace) error {
	for _, item := range t.items {
		if err := item.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralArray) UniqueExpressionNames(ns Namespace) error {
	for _, item := range t.items {
		if err := item.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralArray) Walk(fn WalkFunc) error {
  for _, item := range t.items {
    if err := item.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}
