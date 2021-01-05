package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ImportExport struct {
  Import
}

func NewImportExport(newName *Word, varExpr *VarExpression, path_ *LiteralString, ctx context.Context) (*ImportExport, error) {
  path, err := resolveImportPath(path_)
  if err != nil {
    return nil, err
  }

  return &ImportExport{newImport(newName, varExpr, path, ctx)}, nil
}

func (t *ImportExport) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("ImportExport")
  b.WriteString("\n")
  b.WriteString(t.Import.Dump(indent + "  "))

  return b.String()
}

func (t *ImportExport) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *ImportExport) ResolveStatementNames(scope Scope) error {
  v, err := t.getImportedVariable(scope)
  if err != nil {
    return err
  }

  module := GetModule(scope)
  if module == nil {
    panic("not inside module")
  }

  if err := module.SetExportedVariable(t.newName.Value(), v, t.newName.Context()); err != nil {
    return err
  }

  return nil
}
