package js

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ClassMember interface {
  Context() context.Context
  Name() string
  Dump(indent string) string
  WriteStatement(usage Usage, indent string, nl string, tab string) string
  Role() prototypes.FunctionRole

  IsUniversal() bool // functions are always universal, properties not necessarily
  ResolveNames(scope Scope) error

  GetValue(ctx context.Context) (values.Value, error)
  SetValue(v values.Value, ctx context.Context) error
  Eval() error

  ResolveActivity(usage Usage) error
  UniversalNames(ns Namespace) error
  UniqueNames(ns Namespace) error
  Walk(fn WalkFunc) error
}
