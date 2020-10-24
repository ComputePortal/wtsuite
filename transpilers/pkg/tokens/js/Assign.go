package js

import (
	"strings"

	"./values"

	"../context"
)

type Assign struct {
	lhs Expression
  typeExpr *TypeExpression // can be nil
	rhs Expression
	op  string // eg. "+" for "+=" or "-" for "-=", defaults to empty string
	TokenData
}

func NewAssign(lhs Expression, rhs Expression, op string, ctx context.Context) *Assign {
	if op == ":" || op == "!" || op == "=" || op == "==" || op == "!=" || op == ">" || op == "<" {
		err := ctx.NewError("not a valid assign op '" + op + "'")
		panic(err)
	}

	return &Assign{lhs, nil, rhs, op, TokenData{ctx}}
}

func NewTypedAssign(lhs Expression, typeExpr *TypeExpression, rhs Expression, ctx context.Context) *Assign {
	return &Assign{lhs, typeExpr, rhs, "", TokenData{ctx}}
}

func (t *Assign) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Assign (")
	b.WriteString(t.op)
	b.WriteString("=\n")

  b.WriteString(t.lhs.Dump(indent + "  lhs:"))
  if t.typeExpr != nil {
    b.WriteString(t.typeExpr.Dump(indent + "  type:"))
  }
  b.WriteString(t.rhs.Dump(indent + "  rhs:"))

	return b.String()
}

func (t *Assign) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.lhs.WriteExpression())
	b.WriteString(t.op)
	b.WriteString("=")
	b.WriteString(t.rhs.WriteExpression())

	return b.String()
}

func (t *Assign) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.WriteExpression())

	return b.String()
}

func (t *Assign) Args() []Token {
	return []Token{t.lhs, t.rhs}
}

func (t *Assign) Precedence() int {
	return _binaryPrecedenceMap["="] // same precedence if op != ""
}

func (t *Assign) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Assign) ResolveExpressionNames(scope Scope) error {
	if err := t.lhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

  if t.typeExpr != nil {
    if err := t.typeExpr.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

	if err := t.rhs.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *Assign) HoistNames(scope Scope) error {
	return nil
}

func (t *Assign) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *Assign) HoistValues(stack values.Stack) error {
	return nil
}

// TODO: use the typeExpr for cast
func (t *Assign) EvalExpression(stack values.Stack) (values.Value, error) {
	var rhsValue values.Value
	var err error
	if t.op != "" {
		op, err := NewBinaryOp(t.op, t.lhs, t.rhs, t.Context())
		if err != nil {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: bad assign operator")
		}

		rhsValue, err = op.EvalExpression(stack)
		if err != nil {
			return nil, err
		}
	} else {
		rhsValue, err = t.rhs.EvalExpression(stack)
		if err != nil {
			return nil, err
		}
	}

	if rhsValue == nil {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: rhs is void")
	}

	// eval rhs literal function if its args list is completely typed (this is useful for certain callbacks)
	if rhsFn, ok := t.rhs.(*Function); ok {
		// check if arglists is complete
		if _, argsErr := rhsFn.GenerateArgInstances(stack, t.Context()); argsErr == nil {
			if err := rhsValue.EvalAsEntryPoint(stack, t.Context()); err != nil {
				return nil, err
			}
		}
	}

  // contrain using a type
  if t.typeExpr != nil {
    rhsValue, err = t.typeExpr.Constrain(stack, rhsValue)
    if err != nil {
      return nil, err
    }
  }

	switch lhs := t.lhs.(type) {
	case *VarExpression:
		if err := stack.SetValue(lhs.ref, rhsValue, true,
			t.Context()); err != nil {
			return nil, err
		}
	case *Member:
		if err := lhs.EvalSet(stack, rhsValue, t.Context()); err != nil {
			return nil, err
		}
	case *Index:
		if err := lhs.EvalSet(stack, rhsValue, t.Context()); err != nil {
			return nil, err
		}
	default:
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: unexpected assign lhs")
		//panic("unexpected assign lhs" + reflect.TypeOf(t.lhs).String())
	}

	return rhsValue, nil
}

func (t *Assign) EvalStatement(stack values.Stack) error {
	_, err := t.EvalExpression(stack)

	return err
}

func (t *Assign) IsRegular() bool {
	return t.op == ""
}

func (t *Assign) HasLhsVarExpression() bool {
	_, ok := t.lhs.(*VarExpression)
	return ok
}

func (t *Assign) GetLhsVarExpression() (*VarExpression, error) {
	if lhs, ok := t.lhs.(*VarExpression); ok {
		return lhs, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: not a simple assignment")
	}
}

func (t *Assign) resolveExpressionActivity(usage Usage, isNew bool) error {
	if err := t.rhs.ResolveExpressionActivity(usage); err != nil {
		return err
	}

  if t.typeExpr != nil {
    if err := t.typeExpr.ResolveExpressionActivity(usage); err != nil {
      return err
    }
  }

	switch lhs := t.lhs.(type) {
	case *VarExpression:
		if isNew {
			if err := usage.Rereference(lhs.GetVariable(), t.Context()); err != nil {
				return err
			}
		}
	default:
		if err := t.lhs.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (t *Assign) ResolveExpressionActivity(usage Usage) error {
	return t.resolveExpressionActivity(usage, false)
}

func (t *Assign) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *Assign) UniversalExpressionNames(ns Namespace) error {
	if err := t.lhs.UniversalExpressionNames(ns); err != nil {
		return err
	}

  if t.typeExpr != nil {
    if err := t.typeExpr.UniversalExpressionNames(ns); err != nil {
      return err
    }
  }

	return t.rhs.UniversalExpressionNames(ns)
}

func (t *Assign) UniqueExpressionNames(ns Namespace) error {
	if err := t.lhs.UniqueExpressionNames(ns); err != nil {
		return err
	}

  if t.typeExpr != nil {
    if err := t.typeExpr.UniqueExpressionNames(ns); err != nil {
      return err
    }
  }

	return t.rhs.UniqueExpressionNames(ns)
}

func (t *Assign) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *Assign) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *Assign) Walk(fn WalkFunc) error {
  if err := t.lhs.Walk(fn); err != nil {
    return err
  }

  if t.typeExpr != nil {
    if err := t.typeExpr.Walk(fn); err != nil {
      return err
    }
  }

  if err := t.rhs.Walk(fn); err != nil {
    return err
  }

  if err := fn(t); err != nil {
    return err
  }

  return nil
}

func IsAssign(t Expression) bool {
	_, ok := t.(*Assign)
	return ok
}

func IsSimpleAssign(t Expression) bool {
	assign, ok := t.(*Assign)

	if ok {
		return IsVarExpression(assign.lhs)
	} else {
		return false
	}
}
