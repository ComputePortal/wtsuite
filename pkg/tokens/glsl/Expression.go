package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Expression interface {
  Token

  WriteExpression() string

  ResolveExpressionNames(scope Scope) error

  EvalExpression() (values.Value, error)

  ResolveExpressionActivity(usage Usage) error
}
