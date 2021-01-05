package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type MacroFunction struct {
  nameExpr *VarExpression
  TokenData
}

func newMacroFunction(name string, ctx context.Context) MacroFunction {
  return MacroFunction{NewVarExpression(name, ctx), newTokenData(ctx)}
}

func (t *MacroFunction) Name() string {
  return t.nameExpr.Name()
}

func (t *MacroFunction) GetVariable() Variable {
  return t.nameExpr.GetVariable()
}

func (t *MacroFunction) ResolveStatementNames(scope Scope) error {
  panic("not available")
}

func (t *MacroFunction) EvalStatement() error {
  panic("not available")
}

func (t *MacroFunction) ResolveStatementActivity(usage Usage) error {
  panic("not available")
}

func (t *MacroFunction) UniqueStatementNames(ns Namespace) error {
  ns.FunctionName(t.GetVariable())

  return nil
}
