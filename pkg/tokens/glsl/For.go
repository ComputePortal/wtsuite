package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)


type For struct {
  init Statement
  comp Expression
  incr Statement
  Block
}

func NewFor(init Statement, comp Expression, incr Statement, statements []Statement, ctx context.Context) *For {
  fr := &For{
    init, 
    comp,
    incr, 
    newBlock(ctx),
  }

  for _, st := range statements {
    fr.AddStatement(st)
  }

  return fr
}

func (t *For) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("for(")
  b.WriteString(t.init.Dump(indent + "  "))
  b.WriteString(t.comp.Dump(indent + "  "))
  b.WriteString(t.incr.Dump(indent + "  "))
  b.WriteString(t.Block.Dump(indent + "{ "))

  return b.String()
}

func (t *For) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("for(")
  b.WriteString(t.init.WriteStatement(usage, "", "", ""))
  b.WriteString(";")
  b.WriteString(t.comp.WriteExpression())
  b.WriteString(";")
  b.WriteString(t.incr.WriteStatement(usage, "", "", ""))
  b.WriteString("){")
  b.WriteString(t.Block.writeBlockStatements(usage, indent + tab, nl, tab));
  b.WriteString(indent)
  b.WriteString("}")

  return b.String()
}

func (t *For) ResolveStatementNames(scope Scope) error {
  subScope := NewScope(scope)

  if err := t.init.ResolveStatementNames(subScope); err != nil {
    return err
  }

  if err := t.comp.ResolveExpressionNames(subScope); err != nil {
    return err
  }

  if err := t.incr.ResolveStatementNames(subScope); err != nil {
    return err
  }

  if err := t.Block.ResolveStatementNames(subScope); err != nil {
    return err
  }

  return nil
}

func (t *For) EvalStatement() error {
  if err := t.init.EvalStatement(); err != nil {
    return err
  }

  compVal, err := t.comp.EvalExpression()
  if err != nil {
    return err
  }

  if !values.IsBool(compVal) {
    errCtx := t.comp.Context()
    return errCtx.NewError("Error: expected bool condition, got " + compVal.TypeName())
  }

  if err := t.incr.EvalStatement(); err != nil {
    return err
  }

  if err := t.Block.evalStatements(); err != nil {
    return err
  }

  return err
}

func (t *For) ResolveStatementActivity(usage Usage) error {

  if err := t.Block.ResolveStatementActivity(usage); err != nil {
    return err
  }

  if err := t.incr.ResolveStatementActivity(usage); err != nil {
    return err
  }

  if err := t.comp.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  if err := t.init.ResolveStatementActivity(usage); err != nil {
    return err
  }

  return nil
}

func (t *For) UniqueStatementNames(ns Namespace) error {
  subNs := ns.NewBlockNamespace()

  if err := t.init.UniqueStatementNames(subNs); err != nil {
    return err
  }

  if err := t.incr.UniqueStatementNames(subNs); err != nil {
    return err
  }

  if err := t.Block.UniqueStatementNames(subNs); err != nil {
    return err
  }

  return nil
}
