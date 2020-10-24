package js

import (
	"./values"
)

type Statement interface {
	Token

	AddStatement(st Statement) // panics if Statement is not Block-like

	WriteStatement(indent string) string

	HoistNames(scope Scope) error

	ResolveStatementNames(scope Scope) error

	// only needed for function statements
	HoistValues(stack values.Stack) error

	EvalStatement(stack values.Stack) error

	// usage is resolved in reverse order, so that unused 'mutations' (i.e. variable assignments) can be detected
	ResolveStatementActivity(usage Usage) error

	// universal names need to be registered before other unique names are generated
	UniversalStatementNames(ns Namespace) error

	UniqueStatementNames(ns Namespace) error

  // used be refactoring tools
  Walk(fn WalkFunc) error
}
