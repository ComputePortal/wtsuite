package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ForIn struct {
	ForInOf
}

func NewForIn(varType VarType, lhs *VarExpression, rhs Expression, ctx context.Context) (*ForIn, error) {
	return &ForIn{newForInOf(varType, lhs, rhs, ctx)}, nil
}

func (t *ForIn) Dump(indent string) string {
	return t.ForInOf.dump(indent, "in")
}

func (t *ForIn) WriteStatement(indent string) string {
	return t.ForInOf.writeStatement(indent, "", "in")
}

func (t *ForIn) EvalStatement() error {
	rhsValue, err := t.rhs.EvalExpression()
	if err != nil {
		return err
	}

  inValue, err := rhsValue.GetMember(".getin", false, t.Context())
  if err != nil {
    return err
  }

  variable := t.lhs.GetVariable()

  variable.SetValue(inValue)
  variable.SetConstant()

  return t.Block.EvalStatement()
}

func (t *ForIn) Walk(fn WalkFunc) error { 
  if err := t.ForInOf.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
