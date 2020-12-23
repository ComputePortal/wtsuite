package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type ForOf struct {
	await bool
	ForInOf
}

func NewForOf(await bool, varType VarType, lhs *VarExpression, rhs Expression,
	ctx context.Context) (*ForOf, error) {
	return &ForOf{await, newForInOf(varType, lhs, rhs, ctx)}, nil
}

func (t *ForOf) Dump(indent string) string {
	op := "of"
	if t.await {
		op += "await"
	}
	return t.ForInOf.dump(indent, op)
}

func (t *ForOf) WriteStatement(indent string) string {
	extra := ""
	if t.await {
		extra = "await"
	}
	return t.ForInOf.writeStatement(indent, extra, "of")
}

func (t *ForOf) EvalStatement() error {
	rhsValue, err := t.rhs.EvalExpression()
	if err != nil {
		return err
	}

  ofValue, err := rhsValue.GetMember(".getof", false, t.Context())
  if err != nil {
    return err
  }

  variable := t.lhs.GetVariable()
  variable.SetValue(ofValue)
  variable.SetConstant()

  return t.Block.EvalStatement()
}

func (t *ForOf) Walk(fn WalkFunc) error {
  if err := t.ForInOf.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
