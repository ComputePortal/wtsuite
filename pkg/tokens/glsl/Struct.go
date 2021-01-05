package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type StructEntry struct {
  typeExpr *TypeExpression
  nameExpr *VarExpression
  length int
  ctx context.Context
}

type Struct struct {
  nameExpr *VarExpression
  entries []*StructEntry
  TokenData
}

func NewStructEntry(typeExpr *TypeExpression, nameExpr *VarExpression, length int, ctx context.Context) *StructEntry {
  return &StructEntry{typeExpr, nameExpr, length, ctx}
}

func NewStruct(nameExpr *VarExpression, entries []*StructEntry, ctx context.Context) (*Struct, error) {
  // check uniqueness of entries
  for i, entry := range entries {
    for j, other := range entries {
      if i == j {
        continue
      }

      if entry.Name() == other.Name() {
        errCtx := other.Context()
        err := errCtx.NewError("Error: duplicate struct entry name")
        err.AppendContextString("Info: other declared here", entry.Context())
        return nil, err
      }
    }
  }

  return &Struct{nameExpr, entries, newTokenData(ctx)}, nil
}

func (se *StructEntry) Context() context.Context {
  return se.ctx
}

func (se *StructEntry) Name() string {
  return se.nameExpr.Name()
}

func (se *StructEntry) Instantiate(ctx context.Context) (values.Value, error) {
  val, err := se.typeExpr.Instantiate(ctx)
  if err != nil {
    return nil, err
  }

  if se.length > 0 {
    val = values.NewArray(val, se.length, ctx)
  }

  return val, nil
}

func (t *Struct) Name() string {
  return t.nameExpr.Name()
}

func (t *Struct) GetVariable() Variable {
  return t.nameExpr.GetVariable()
}

func (t *StructEntry) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(t.typeExpr.Dump(indent + "t:"))
  b.WriteString(t.nameExpr.Dump(indent + "n:"))
  if t.length > 0 {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  return b.String()
}

func (t *Struct) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Struct(")
  b.WriteString(t.Name())
  b.WriteString(")\n")

  for _, entry := range t.entries {
    b.WriteString(entry.Dump(indent + "  "))
  }

  return b.String()
}

func (t *StructEntry) writeEntry() string {
  var b strings.Builder
  
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())
  
  if t.length > 0 {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  b.WriteString(";")

  return b.String()
}

func (t *Struct) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("struct ")
  b.WriteString(t.Name())
  b.WriteString("{")
  b.WriteString(nl)
  
  for _, entry := range t.entries {
    b.WriteString(indent + tab)
    b.WriteString(entry.writeEntry())
    b.WriteString(nl)
  }

  b.WriteString(indent)
  b.WriteString("};")
  b.WriteString(nl)

  return b.String()
}

func (se *StructEntry) ResolveNames(scope Scope) error {
  if err := se.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  return nil
}

func (t *Struct) ResolveStatementNames(scope Scope) error {
  for _, entry := range t.entries {
    if err := entry.ResolveNames(scope); err != nil {
      return err
    }
  }

  variable := t.GetVariable()

  if err := scope.SetVariable(t.Name(), variable); err != nil {
    return err
  }

  return nil
}

func (t *Struct) EvalStatement() error {
  variable := t.GetVariable()

  variable.SetValue(values.NewStructType(t, t.Context()))

  return nil
}

func (t *Struct) CheckConstruction(args []values.Value, ctx context.Context) error {
  if len(args) != len(t.entries) {
    errCtx := ctx
    return errCtx.NewError("Error: expected " + strconv.Itoa(len(t.entries)) + ", got " + strconv.Itoa(len(args)))
  }

  for i, entry := range t.entries {
    checkVal, err := entry.Instantiate(ctx)
    if err != nil {
      return err
    }

    if err := checkVal.Check(args[i], ctx); err != nil {
      return err
    }
  }

  return nil
}

func (t *Struct) GetMember(key string, ctx context.Context) (values.Value, error) {
  for _, entry := range t.entries {
    if entry.Name() == key {
      return entry.Instantiate(ctx)
    }
  }

  return nil, ctx.NewError("Error: " + t.Name() + "." + key + " not found")
}

func (t *Struct) SetMember(key string, arg values.Value, ctx context.Context) error {
  for _, entry := range t.entries {
    if entry.Name() == key {
      checkVal, err := entry.Instantiate(ctx)
      if err != nil {
        return err
      }

      return checkVal.Check(arg, ctx)
    }
  }

  return ctx.NewError("Error: " + t.Name() + "." + key + " not found")
}

func (t *Struct) ResolveStatementActivity(usage Usage) error {
  return nil
}

func (t *Struct) UniqueStatementNames(ns Namespace) error {
  ns.VarName(t.GetVariable())

  return nil
}
