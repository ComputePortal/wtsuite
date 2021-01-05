package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Pointer struct {
  typeExpr *TypeExpression
  nameExpr *VarExpression
  length int
  TokenData
}

func newPointer(typeExpr *TypeExpression, nameExpr *VarExpression, length int, ctx context.Context) Pointer {
  return Pointer{typeExpr, nameExpr, length, newTokenData(ctx)}
}

func (t *Pointer) Name() string {
  return t.nameExpr.Name()
}

func (t *Pointer) GetVariable() Variable {
  return t.nameExpr.GetVariable()
}

func (t *Pointer) ResolveStatementNames(scope Scope) error {
  if err := t.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  variable := t.GetVariable()

  if err := scope.SetVariable(t.nameExpr.Name(), variable); err != nil {
    return err
  }

  return nil
}

func (t *Pointer) EvalStatement() error {
  variable := t.GetVariable()

  val, err := t.typeExpr.Instantiate(t.Context())
  if err != nil {
    return err
  }

  if t.length > 0 {
    val = values.NewArray(val, t.length, t.Context())
  }

  variable.SetValue(val)

  return nil
}

func (t *Pointer) ResolveStatementActivity(usage Usage) error {
  return usage.Rereference(t.nameExpr.GetVariable(), t.Context())
}

// likely exported, so keep original name
func (t *Pointer) UniqueStatementNames(ns Namespace) error {
  return ns.OrigName(t.GetVariable())
}
