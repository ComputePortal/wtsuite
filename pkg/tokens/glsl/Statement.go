package glsl

type Statement interface {
  Token

  WriteStatement(usage Usage, indent string, nl string, tab string) string

  ResolveStatementNames(scope Scope) error

  EvalStatement() error

  ResolveStatementActivity(usage Usage) error

  UniqueStatementNames(ns Namespace) error
}
