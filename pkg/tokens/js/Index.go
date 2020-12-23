package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Index struct {
	container Expression
	index     Expression
	TokenData
}

func NewIndex(container Expression, index Expression, ctx context.Context) *Index {
	return &Index{container, index, TokenData{ctx}}
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

///////////////////////////
// 1. Name resolution stage
///////////////////////////
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

	indexValue, err := t.index.EvalExpression()
	if err != nil {
		return nil, err
	}

  fn, err := containerValue.GetMember(".getindex", false, t.Context())
  if err != nil {
    return nil, err
  }

  result, err := fn.EvalFunction([]values.Value{indexValue}, false, t.Context())
  if err != nil {
		context.AppendContextString(err, "Info: container", t.container.Context())
		return nil, err
  }

  if result == nil {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: index returns void")
  }

	return result, nil
}

func (t *Index) EvalSet(rhsValue values.Value, ctx context.Context) error {
	containerValue, err := t.container.EvalExpression()
	if err != nil {
		return err
	}

	indexValue, err := t.index.EvalExpression()
	if err != nil {
		return err
	}

  fn, err := containerValue.GetMember(".setindex", false, t.Context())
  if err != nil {
    return err
  }

  if _, err := fn.EvalFunction([]values.Value{indexValue, rhsValue}, false, t.Context()); err != nil {
		context.AppendContextString(err, "Info: container", t.container.Context())
    return err
  }

	return nil
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

func (t *Index) UniversalExpressionNames(ns Namespace) error {
	if err := t.container.UniversalExpressionNames(ns); err != nil {
		return err
	}

	if err := t.index.UniversalExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *Index) UniqueExpressionNames(ns Namespace) error {
	if err := t.container.UniqueExpressionNames(ns); err != nil {
		return err
	}

	if err := t.index.UniqueExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *Index) Walk(fn WalkFunc) error {
  if err := t.container.Walk(fn); err != nil {
    return err
  }

  if err := t.index.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
