package js

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

func (t *Export) AddStatement(st Statement) {
  panic("not available")
}

func (t *Export) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  return ""
}

func (t *Export) HoistNames(scope Scope) error {
  return nil
}

func (t *Export) ResolveStatementNames(scope Scope) error {
  module := GetModule(scope)

  if module == nil {
    panic("not inside a module")
  }

  if err := t.varExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  if err := module.AddExportedName(t.newName.Value(), t.varExpr.Name(), t.varExpr.GetVariable(), t.newName.Context()); err != nil {
    return err
  }

  return nil
}

func (t *Export) EvalStatement() error {
  return nil
}

func (t *Export) ResolveStatementActivity(usage Usage) error {
  if err := t.varExpr.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  return nil
}

func (t *Export) UniversalStatementNames(ns Namespace) error {
  if err := t.varExpr.UniversalExpressionNames(ns); err != nil {
    return err
  }
  
  return nil
}

func (t *Export) UniqueStatementNames(ns Namespace) error {
  if err := t.varExpr.UniqueExpressionNames(ns); err != nil {
    return err
  }
  
  return nil
}

func (t *Export) Walk(fn WalkFunc) error {
  if err := t.varExpr.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
