package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Import struct {
  newName *Word
  varExpr *VarExpression
  path    *LiteralString
  TokenData
}

func resolveImportPath(pathLit *LiteralString) (*LiteralString, error) {
  ctx := pathLit.Context()
  path := pathLit.Value()

  absPath := path
  var err error
  if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
    absPath, err = files.Search(ctx.Path(), path)
  } else {
    absPath, err = files.SearchShader(ctx.Path(), path)
  }
  if err != nil {
    return nil, err
  }

  pathLit = NewLiteralString(absPath, ctx)

  return pathLit, nil
}

func newImport(newName *Word, varExpr *VarExpression, path *LiteralString, ctx context.Context) Import {
  return Import{newName, varExpr, path, newTokenData(ctx)}
}

func NewImport(newName *Word, varExpr *VarExpression, path_ *LiteralString, ctx context.Context) (*Import, error) {
  path, err := resolveImportPath(path_)
  if err != nil {
    return nil, err
  }

  imp := newImport(newName, varExpr, path, ctx)

  return &imp, nil
}

func (t *Import) Path() *LiteralString {
  return t.path
}

func (t *Import) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Import")
  b.WriteString("\n")

  if t.newName.Value() == t.varExpr.Name() {
    b.WriteString(t.varExpr.Dump(indent + "  "))
  } else {
    b.WriteString(t.newName.Dump(indent + "  "))
    b.WriteString("\n")
    b.WriteString(t.varExpr.Dump(indent + "as"))
    b.WriteString("\n")
  }

  b.WriteString(t.path.Value())
  b.WriteString("\n")

  return b.String()
}

func (t *Import) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *Import) getImportedVariable(scope Scope) (Variable, error) {
  module := GetModule(scope)
  globals := GetGlobalScope(scope)

  if module == nil {
    panic("not inside a module")
  }

  // get the module from which it is imported
  importedModule, err := globals.GetModule(t.path.Value())
  if err != nil {
    errCtx := t.path.Context()
    return nil, errCtx.NewError("Error: module not found")
  }

  v, err := importedModule.GetExportedVariable(globals, t.varExpr.Name(),
    t.path.Context())
  if err != nil {
    return nil, err
  }

  // package can be given the correct name now
  if pkg, ok := v.(*Package); ok {
    pkg.Rename(t.newName.Value())
  }

  return v, nil
}

func (t *Import) ResolveStatementNames(scope Scope) error {
  v, err := t.getImportedVariable(scope)
  if err != nil {
    return err
  }

  if err := scope.SetVariable(t.newName.Value(), v); err != nil {
    return err
  }

  return nil
}

func (t *Import) EvalStatement() error {
  return nil
}

func (t *Import) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Import) UniqueStatementNames(ns Namespace) error {
  return nil
}

