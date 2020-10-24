package js

import (
	"../context"

	"./values"
)

// implements the Statement interface
type NodeJSModule struct {
	expr   *VarExpression
	module values.Prototype
	TokenData
}

func NewNodeJSModule(expr *VarExpression, module values.Prototype, ctx context.Context) *NodeJSModule {
	return &NodeJSModule{
		expr,
		module,
		newTokenData(ctx),
	}
}

func (m *NodeJSModule) Dump(indent string) string {
	if m.expr.Name() == m.module.Name() {
		return indent + "NodeJSModule(" + m.module.Name() + ")\n"
	} else {
		return indent + "NodeJSModule(" + m.module.Name() + " as " + m.expr.Name() + ")\n"
	}
}

func (m *NodeJSModule) AddStatement(st Statement) {
	panic("not a block statement")
}

func (m *NodeJSModule) WriteStatement(indent string) string {
	return "const " + m.expr.Name() + "=require('" + m.module.Name() + "');"
}

func (m *NodeJSModule) HoistNames(scope Scope) error {
	return nil
}

func (m *NodeJSModule) ResolveStatementNames(scope Scope) error {
	variable := m.expr.GetVariable()
	variable.SetObject(m.module)

	if err := scope.SetVariable(m.expr.Name(), variable); err != nil {
		return err
	}

	return nil
}

func (m *NodeJSModule) HoistValues(stack values.Stack) error {
	return nil
}

func (m *NodeJSModule) EvalStatement(stack values.Stack) error {
	// abuse the values.class/prototype structure
	if err := stack.SetValue(m.expr.GetVariable(), values.NewClass(m.module, m.Context()), false, m.Context()); err != nil {
		return err
	}

	return nil
}

func (m *NodeJSModule) ResolveStatementActivity(usage Usage) error {
	return m.expr.ResolveExpressionActivity(usage)
}

func (m *NodeJSModule) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (m *NodeJSModule) UniqueStatementNames(ns Namespace) error {
	return m.expr.uniqueDeclarationName(ns, CONST)
}

func (m *NodeJSModule) Walk(fn WalkFunc) error {
  if err := m.expr.Walk(fn); err != nil {
    return err
  }

  return fn(m)
}
