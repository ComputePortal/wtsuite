package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Export struct {
  newName *Word
  varExpr *VarExpression
  TokenData
}

func NewExport(newName *Word, varExpr *VarExpression, ctx context.Context) *Export {
  return &Export{newName, varExpr, newTokenData(ctx)}
}

func (t *Export) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Export")
  b.WriteString("\n")

  if t.newName.Value() == t.varExpr.Name() {
    b.WriteString(t.varExpr.Dump(indent + "  "))
  } else {
    b.WriteString(t.newName.Dump(indent + "  "))
    b.WriteString("\n")
    b.WriteString(t.varExpr.Dump(indent + "as"))
    b.WriteString("\n")
  }

  return b.String()
}

func (t *Export) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *Export) ResolveStatementNames(scope Scope) error {
  module := GetModule(scope)

  if module == nil {
    panic("not inside a module")
  }

  if err := t.varExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if err := module.SetExportedVariable(t.newName.Value(), t.varExpr.GetVariable(), t.newName.Context()); err != nil {
    return err
  }

  return nil
}

func (t *Export) EvalStatement() error {
  return nil
}

func (t *Export) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Export) UniqueStatementNames(ns Namespace) error {
  return nil
}
