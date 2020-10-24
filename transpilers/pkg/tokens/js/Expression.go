package js

import (
	"./values"
)

type Expression interface {
	WriteExpression() string

	ResolveExpressionNames(scope Scope) error

	EvalExpression(stack values.Stack) (values.Value, error)

	ResolveExpressionActivity(usage Usage) error

	// universal names need to be registered before other unique names are generated
	UniversalExpressionNames(ns Namespace) error

	UniqueExpressionNames(ns Namespace) error

  Walk(fn WalkFunc) error

	Token
}

// it is easier to do this via assertion than to add a method to each expression
func IsLiteral(expr Expression) bool {
	switch expr.(type) {
	case *LiteralArray, *LiteralBoolean, *LiteralObject, *LiteralInt, *LiteralFloat, *LiteralString:
		return true
	default:
		return false
	}
}
