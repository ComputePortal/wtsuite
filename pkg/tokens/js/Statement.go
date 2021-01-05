package js

import ()

type Statement interface {
	Token

	AddStatement(st Statement) // panics if Statement is not Block-like

	WriteStatement(usage Usage, indent string, nl string, tab string) string

	HoistNames(scope Scope) error

	ResolveStatementNames(scope Scope) error

	EvalStatement() error

	// usage is resolved in reverse order, so that unused 'mutations' (i.e. variable assignments) can be detected
	ResolveStatementActivity(usage Usage) error

	// universal names need to be registered before other unique names are generated
	UniversalStatementNames(ns Namespace) error

	UniqueStatementNames(ns Namespace) error

  // used be refactoring tools
  Walk(fn WalkFunc) error
}
