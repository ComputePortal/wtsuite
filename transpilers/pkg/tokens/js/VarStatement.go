package js

import (
	"strings"

	"./values"

	"../context"
)

type VarStatement struct {
	varType VarType
	exprs   []Expression // contains VarExpressions (or Assigns to VarExpressions)
	TokenData
}

func NewVarStatement(varType VarType, exprs []Expression,
	ctx context.Context) (*VarStatement, error) {
	// check that all expressions are VarExpressions or Assign to VarExpressions
	for _, expr_ := range exprs {
		switch expr := expr_.(type) {
		case *VarExpression:
			if varType == CONST {
				expr.ref.SetConstant()
			}
		case *Assign:
			lhs, err := expr.GetLhsVarExpression()
			if err != nil {
				return nil, err
			}

			if varType == CONST {
				lhs.ref.SetConstant()
			}

		default:
			errCtx := expr.Context()
			return nil, errCtx.NewError("Error: not a VarExpression or Assign to VarExpression")
		}
	}

	return &VarStatement{varType, exprs, TokenData{ctx}}, nil
}

func (t *VarStatement) GetVariables() map[string]Variable {
	// collect the variables from the expressions

	variables := make(map[string]Variable)

	for _, expr_ := range t.exprs {
		switch expr := expr_.(type) {
		case *VarExpression:
			variables[expr.Name()] = expr.GetVariable()
		case *Assign:
			lhs, err := expr.GetLhsVarExpression()
			if err != nil {
				panic("should've been caught during construction")
			}
			variables[lhs.Name()] = lhs.GetVariable()
		default:
			panic("should've been caught during construction")
		}
	}

	return variables
}

func (t *VarStatement) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("VarStatement(")
	b.WriteString(VarTypeToString(t.varType))
	b.WriteString(")\n")

	for _, expr := range t.exprs {
		b.WriteString(expr.Dump(indent + "  "))
	}

	return b.String()
}

func (t *VarStatement) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(VarTypeToString(t.varType))
	b.WriteString(" ")

	for i, expr := range t.exprs {
		b.WriteString(expr.WriteExpression())
		if i < len(t.exprs)-1 {
			b.WriteString(",")
		}
	}

	return b.String()
}

func (t *VarStatement) AddStatement(st Statement) {
	panic("not a block")
}

func (t *VarStatement) assertUnique(scope Scope, name string) error {
	if scope.HasVariable(name) {
		prev, _ := scope.GetVariable(name)

		errCtx := t.Context()
		err := errCtx.NewError("Error: already defined")

		err.AppendContextString("Info: defined here", prev.Context())

		return err
	}

	return nil
}

func (t *VarStatement) HoistNames(scope Scope) error {
	if t.varType == VAR {
		for _, expr_ := range t.exprs {
			switch expr := expr_.(type) {
			case *Assign:
				lhs, err := expr.GetLhsVarExpression()
				if err != nil {
					return err
				}

				if err := t.assertUnique(scope, lhs.Name()); err != nil {
					return err
				}

				if err := scope.SetVariable(lhs.Name(), lhs.GetVariable()); err != nil {
					return err
				}
			case *VarExpression:
				if err := t.assertUnique(scope, expr.Name()); err != nil {
					return err
				}

				if err := scope.SetVariable(expr.Name(), expr.GetVariable()); err != nil {
					return err
				}
			default:
				panic("invalid VarStatement expr (should be Assign or VarExpression)")
			}
		}
	}

	return nil
}

func (t *VarStatement) ResolveStatementNames(scope Scope) error {
	setVar := func(name string, ref Variable) error {
		switch t.varType {
		case LET, CONST:
			if err := t.assertUnique(scope, name); err != nil {
				return err
			}

			if err := scope.SetVariable(name, ref); err != nil {
				return err
			}
		case VAR:
			if !scope.HasVariable(name) {
				panic("should've been added during construction")
			}
		default:
			panic("unhandled")
		}

		return nil
	}

	for _, expr_ := range t.exprs {
		switch expr := expr_.(type) {
		case *Assign:
			lhs, err := expr.GetLhsVarExpression()
			if err != nil {
				return err
			}

			if err := setVar(lhs.Name(), lhs.GetVariable()); err != nil {
				return err
			}

			if err := expr.rhs.ResolveExpressionNames(scope); err != nil {
				return err
			}
		case *VarExpression:
			if err := setVar(expr.Name(), expr.GetVariable()); err != nil {
				return err
			}
		default:
			panic("invalid VarStatement expr (should be Assign or VarExpression)")
		}
	}

	return nil
}

func (t *VarStatement) EvalStatement(stack values.Stack) error {
	for _, expr_ := range t.exprs {
		switch expr := expr_.(type) {
		case *Assign:
			rhsValue, err := expr.rhs.EvalExpression(stack)
			if err != nil {
				return err
			}

			nameExpr, err := expr.GetLhsVarExpression()
			if err != nil {
				panic(err)
			}

			v := nameExpr.GetVariable()

			if err := stack.SetValue(v, rhsValue, false, expr.Context()); err != nil {
				return err
			}
		case *VarExpression:
			ref := expr.GetVariable()

			if err := stack.SetValue(ref, nil, false, expr.Context()); err != nil {
				return err
			}
		default:
			panic("unhandled")
		}
	}

	return nil
}

func (t *VarStatement) HoistValues(stack values.Stack) error {
	return nil
}

func (t *VarStatement) ResolveStatementActivity(usage Usage) error {
	for i := len(t.exprs) - 1; i >= 0; i-- {
		expr_ := t.exprs[i]

		switch expr := expr_.(type) {
		case *Assign:
			if err := expr.resolveExpressionActivity(usage, true); err != nil {
				return err
			}

			// XXX: why was this in regular forward order?
			/*
				lhs, err := expr.GetLhsVarExpression()
				if !ok {
					panic("unexpected")
				}

				if err := usage.Rereference(lhs.ref, lhs.Context()); err != nil {
					return err
				}

				if err := expr.rhs.ResolveExpressionActivity(usage); err != nil {
					return err
				}
			*/
		case *VarExpression:
			ref := expr.GetVariable()

			if ref == nil {
				panic("ref expected to be set")
			}

			if err := usage.Rereference(ref, expr.Context()); err != nil {
				return err
			}
		default:
			panic("unexpected")
		}
	}

	return nil
}

func (t *VarStatement) UniversalStatementNames(ns Namespace) error {
	for _, expr := range t.exprs {
		if err := expr.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *VarStatement) UniqueStatementNames(ns Namespace) error {
	for _, expr_ := range t.exprs {
		switch expr := expr_.(type) {
		case *Assign:
			lhs, err := expr.GetLhsVarExpression()
			if err != nil {
				panic(err)
			}

			if err := lhs.uniqueDeclarationName(ns, t.varType); err != nil {
				return err
			}

			// like in assign, first left then right
			if err := expr.rhs.UniqueExpressionNames(ns); err != nil {
				return err
			}
		case *VarExpression:
			if err := expr.uniqueDeclarationName(ns, t.varType); err != nil {
				return err
			}
		default:
			panic("unexpected")
		}
	}

	return nil
}

func (t *VarStatement) Walk(fn WalkFunc) error {
  for _, expr := range t.exprs {
    if err := expr.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}
