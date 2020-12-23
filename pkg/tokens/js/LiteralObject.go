package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralObjectMember struct {
	key   *Word
	value Expression
}

type LiteralObject struct {
	items []*LiteralObjectMember
	TokenData
}

func NewLiteralObject(keys []*Word, values []Expression, ctx context.Context) (*LiteralObject, error) {
	// check validity of items
	n := len(keys)
	if n != len(values) {
		panic("keys and values not same length")
	}

	items := make([]*LiteralObjectMember, n)
	for i, key := range keys {
		value := values[i]
		if key == nil {
			panic("nil key")
		}
		if value == nil {
			panic("nil value")
		}

		items[i] = &LiteralObjectMember{key, value}
	}

	return &LiteralObject{items, TokenData{ctx}}, nil
}

func (t *LiteralObject) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("LiteralObject\n")

	for _, item := range t.items {
		b.WriteString(item.key.Dump(indent) + "  ")
		b.WriteString(item.value.Dump(indent + ": "))
	}

	return b.String()
}

func (t *LiteralObject) WriteExpression() string {
	var b strings.Builder

	b.WriteString("{")

	for i, item := range t.items {
		b.WriteString("'")
		b.WriteString(item.key.value)
		b.WriteString("'")

		b.WriteString(":")
		b.WriteString(item.value.WriteExpression())

		if i < len(t.items)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("}")
	return b.String()
}

func (t *LiteralObject) ResolveExpressionNames(scope Scope) error {
	for _, item := range t.items {
		if err := item.value.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralObject) EvalExpression() (values.Value, error) {
	props := make(map[string]values.Value)

	for _, item := range t.items {
		itemValue, err := item.value.EvalExpression()
		if err != nil {
			return nil, err
		}

		if prev, ok := props[item.key.Value()]; ok {
			errCtx := item.key.Context()
			err := errCtx.NewError("Error: key already set")
			err.AppendContextString("Info: set here", prev.Context())
			return nil, err
		}

		props[item.key.Value()] = itemValue
	}

  return prototypes.NewObject(props, t.Context()), nil
}

func (t *LiteralObject) ResolveExpressionActivity(usage Usage) error {
	// in reverse order
	for i := len(t.items) - 1; i >= 0; i-- {
		item := t.items[i]
		if err := item.value.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (t *LiteralObject) UniversalExpressionNames(ns Namespace) error {
	for _, item := range t.items {
		if err := item.value.UniversalExpressionNames(ns); err != nil {
			return err
		}
		// keys are never renamed
	}

	return nil
}

func (t *LiteralObject) UniqueExpressionNames(ns Namespace) error {
	for _, item := range t.items {
		if err := item.value.UniqueExpressionNames(ns); err != nil {
			return err
		}
		// keys are never renamed
	}

	return nil
}

func (t *LiteralObject) Walk(fn WalkFunc) error {
  for _, item := range t.items {
    if err := item.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}

func (m *LiteralObjectMember) Walk(fn WalkFunc) error {
  if err := m.key.Walk(fn); err != nil {
    return err
  }

  if err := m.value.Walk(fn); err != nil {
    return err
  }

  return fn(m) 
}
