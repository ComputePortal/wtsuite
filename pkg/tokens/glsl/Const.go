package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Const struct {
  altRHS string
  rhsExpr Expression // can be nil
  isExport bool // useful for OrigName or not
  Pointer
}

func NewConst(typeExpr *TypeExpression, name string, n int, rhsExpr Expression, isExport bool, ctx context.Context) *Const {
  return &Const{"", rhsExpr, isExport, newPointer(typeExpr, NewVarExpression(name, ctx), n, ctx)}
}

func (t *Const) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Const(")

  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))

  if (t.length > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  if t.altRHS != "" {
    b.WriteString("=")
    b.WriteString(t.altRHS)
  } else if t.rhsExpr != nil {
    b.WriteString("=")
    b.WriteString(t.rhsExpr.Dump(""))
  }

  b.WriteString(")\n")

  return b.String()
}

func (t *Const) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("const ")
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  if (t.length > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  if t.altRHS != "" {
    b.WriteString("=")
    b.WriteString(t.altRHS)
  } else if t.rhsExpr != nil {
    b.WriteString("=")
    b.WriteString(t.rhsExpr.WriteExpression())
  }
  b.WriteString(";")

  return b.String()
}

func (t *Const) TypeName() string {
  return t.typeExpr.WriteExpression()
}

func (t *Const) Length() int {
  return t.length
}

func (t *Const) SetAltRHS(str string) {
  t.altRHS = str
}

func (t *Const) ResolveStatementNames(scope Scope) error {
  if err := t.Pointer.ResolveStatementNames(scope); err != nil {
    return err
  }

  variable := t.GetVariable()
  variable.SetObject(t)

  return nil
}

func (t *Const) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Const) UniqueStatementNames(ns Namespace) error {
  if t.isExport {
    return ns.OrigName(t.GetVariable())
  } else {
    ns.VarName(t.GetVariable())
    return nil
  }
}
