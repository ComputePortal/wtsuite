package js

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Import struct {
  name string
  TokenData
}

func NewImport(name string, ctx context.Context) *Import {
  return &Import{name, newTokenData(ctx)}
}

func (t *Import) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Import(")
  b.WriteString(t.name)
  b.WriteString(")")

  return b.String()
}

func (t *Import) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *Import) AddStatement(st Statement) {
  panic("not available")
}

func (t *Import) HoistNames(scope Scope) error {
  return nil
}

func (t *Import) ResolveStatementNames(scope Scope) error {
  module := GetModule(scope)
  globals := GetGlobalScope(scope)

  if module == nil {
    panic("not inside a module")
  }

	if importedName, ok := module.importedNames[t.name]; ok {
    var v Variable 
		if importedName.v != nil {
      v = importedName.v
		} else {
			// get the module from which it is imported
			importedModule, err := globals.GetModule(importedName.dep.Value())
			if err != nil {
				errCtx := importedName.dep.Context()
				return errCtx.NewError("Error: module not found")
			}

			// cache elsewhere
			v, err = importedModule.GetExportedVariable(globals, importedName.old,
				importedName.dep.Context())
			if err != nil {
				return err
			}

			// package can be given the correct name now
			if pkg, ok := v.(*Package); ok {
				pkg.Rename(t.name)
			}

      // cache here too?
      importedName.v = v
      module.importedNames[t.name] = importedName
		}

    if err := scope.SetVariable(t.name, v); err != nil {
      return err
    }
	} else {
    errCtx := t.Context()
    err := errCtx.NewError("Internal Error: import not found in module")
    panic(err)
  }

  return nil
}

func (t *Import) EvalStatement() error {
  return nil
}

func (t *Import) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Import) UniversalStatementNames(ns Namespace) error {
  return nil
}

func (t *Import) UniqueStatementNames(ns Namespace) error {
  return nil
}

func (t *Import) Walk(fn WalkFunc) error {
  return fn(t)
}
