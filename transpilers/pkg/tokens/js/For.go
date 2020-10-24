package js

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type For struct {
	inits []Expression
	cond  Expression
	incrs []Expression
	ForBlock
}

// inits can be empty, cond can be nil, incrs can be empty
// in extreme case 'for(;;);' is written
func NewFor(varType VarType, inits []Expression, cond Expression, incrs []Expression,
	ctx context.Context) (*For, error) {

	return &For{inits, cond, incrs, newForBlock(varType, ctx)}, nil
}

func (t *For) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("For ")
	b.WriteString(VarTypeToString(t.varType))
	b.WriteString("\n")

	for _, expr := range t.inits {
		b.WriteString(expr.Dump(indent + "[init] "))
	}

	b.WriteString(t.cond.Dump(indent + "[cond] "))

	for _, expr := range t.incrs {
		b.WriteString(expr.Dump(indent + "[incr]  "))
	}

	for _, s := range t.statements {
		b.WriteString(s.Dump(indent + "{ "))
	}

	return b.String()
}

func (t *For) WriteStatement(indent string) string {
	// varType printed before inits, unless len(inits) == 0

	var b strings.Builder

	hasInits := len(t.inits) > 0
	b.WriteString(t.writeStatementHeader(indent, "", hasInits))

	if hasInits {
		for i, expr := range t.inits {
			b.WriteString(expr.WriteExpression())

			if i < len(t.inits)-1 {
				b.WriteString(",")
			}
		}
	}

	b.WriteString(";")

	if t.cond != nil {
		b.WriteString(t.cond.WriteExpression())
	}

	b.WriteString(";")

	if len(t.incrs) > 0 {
		for i, expr := range t.incrs {
			b.WriteString(expr.WriteExpression())

			if i < len(t.incrs)-1 {
				b.WriteString(",")
			}
		}
	}

	b.WriteString(t.writeStatementFooter(indent))

	return b.String()
}

func (t *For) getInitAssignments() []*Assign {
	inits := make([]*Assign, 0)

	for _, init_ := range t.inits {
		switch init := init_.(type) {
		case *Assign:
			inits = append(inits, init)
		default:
			panic("expected assign")
		}
	}

	return inits
}

func (t *For) HoistNames(scope Scope) error {
	if t.varType == VAR {
		inits := t.getInitAssignments()

		for _, init := range inits {
			lhs, err := init.GetLhsVarExpression()
			if err != nil {
				panic(err)
			}

			if err := scope.SetVariable(lhs.Name(), lhs.GetVariable()); err != nil {
				return err
			}
		}
	}

	return t.Block.HoistNames(scope)
}

func (t *For) ResolveStatementNames(scope Scope) error {
	subScope := NewLoopScope(scope)

	inits := t.getInitAssignments()

	// add inits to scope
	for _, init := range inits {
		lhs, err := init.GetLhsVarExpression()
		if err != nil {
			panic(err)
		}

		name := lhs.Name()

		switch t.varType {
		case CONST, LET:
			if err := subScope.SetVariable(name, lhs.GetVariable()); err != nil {
				return err
			}
		case VAR:
			if !scope.HasVariable(name) {
				panic("should've been added during construction")
			}
		default:
			panic("unhandled")
		}

		if err := init.rhs.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	if err := t.cond.ResolveExpressionNames(subScope); err != nil {
		return err
	}

	for _, incr := range t.incrs {
		if err := incr.ResolveExpressionNames(subScope); err != nil {
			return err
		}
	}

	return t.Block.ResolveStatementNames(subScope)
}

// return true if literal array loop was detected
func (t *For) evalLiteralArrayLoop(stack values.Stack) (bool, error) {
	return false, nil

	if len(t.inits) != 1 ||
		len(t.incrs) != 1 ||
		!IsSimpleAssign(t.inits[0]) ||
		!IsSimpleLT(t.cond) ||
		!IsSimplePostIncr(t.incrs[0]) ||
		t.varType != LET {
		return false, nil
	}

	init, _ := t.inits[0].(*Assign)
	cond, _ := t.cond.(*LTOp)
	incr, _ := t.incrs[0].(*PostIncrOp)

	vExpr1, _ := init.lhs.(*VarExpression)
	vExpr2, _ := cond.a.(*VarExpression)
	vExpr3, _ := incr.a.(*VarExpression)

	if vExpr1.Name() != vExpr2.Name() ||
		vExpr1.Name() != vExpr3.Name() {
		return false, nil
	}

	if !(vExpr1.ref == vExpr2.ref && vExpr1.ref == vExpr3.ref) {
		panic("refs don't correspond") // XXX: panic or return?
	}

	dummyStack := NewBranchStack(stack)

	instance1, err := t.inits[0].EvalExpression(dummyStack)
	if err != nil {
		return false, err
	}

	// length
	instance2, err := cond.b.EvalExpression(dummyStack)
	if err != nil {
		return false, err
	}

	instance3, err := t.incrs[0].EvalExpression(dummyStack)
	if err != nil {
		return false, err
	}

	if !instance1.IsInstanceOf(prototypes.Int) ||
		!instance2.IsInstanceOf(prototypes.Int) ||
		!instance3.IsInstanceOf(prototypes.Int) {
		return false, nil
	}

	i0, ok0 := instance1.LiteralIntValue()
	if !ok0 {
		return false, nil
	}

	i1, ok1 := instance2.LiteralIntValue()
	if !ok1 {
		return false, nil
	}

	indexCtx := vExpr1.Context()

	if VERBOSITY >= 3 || (VERBOSITY >= 1 && (i1-i0) > 3) {
		errCtx := t.Context()
		err := errCtx.NewError("Warning: evaluating literal loop " + strconv.Itoa(i1-i0) + " times (hint: wrap length with Int() to avoid this)")
		fmt.Fprintf(os.Stderr, err.Error())
	}

	if i0 < i1 {
		// there is a limit!
		for i := i0; i < i1 && i < 100; i++ {
			subStack := NewBranchStack(stack)

			indexVal := prototypes.NewLiteralInt(i, indexCtx)

			// dont allow branching!
			if t.varType == VAR {
				if err := stack.SetValue(vExpr1.ref, indexVal, false, indexCtx); err != nil {
					return false, err
				}
			} else {
				if err := subStack.SetValue(vExpr1.ref, indexVal, false, indexCtx); err != nil {
					return false, err
				}
			}

			if err := t.Block.EvalStatement(subStack); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (t *For) EvalStatement(stack values.Stack) error {
	done, err := t.evalLiteralArrayLoop(stack)
	if err != nil {
		return err
	}

	if done {
		return nil
	}

	subStack := NewBranchStack(stack)

	for _, init := range t.getInitAssignments() {
		lhs, err := init.GetLhsVarExpression()
		if err != nil {
			panic(err)
		}

		rhsVal, err := init.rhs.EvalExpression(subStack)
		if err != nil {
			return err
		}

		if t.varType == VAR {
			if err := stack.SetValue(lhs.GetVariable(), rhsVal, true,
				init.Context()); err != nil {
				return err
			}
		} else {
			if err := subStack.SetValue(lhs.GetVariable(), rhsVal, true,
				init.Context()); err != nil {
				return err
			}
		}
	}

	for _, incr := range t.incrs {
		if _, err := incr.EvalExpression(subStack); err != nil {
			return err
		}
	}

	condVal, err := t.cond.EvalExpression(subStack)
	if err != nil {
		return err
	}

	if !condVal.IsInstanceOf(prototypes.Boolean) {
		errCtx := condVal.Context()
		return errCtx.NewError("Error: expected boolean")
	}

	return t.Block.EvalStatement(subStack)
}

func (t *For) ResolveStatementActivity(usage Usage) error {
	// usage is resolved in reverse order (see Statement.go for more details)
	if err := t.Block.ResolveStatementActivity(usage); err != nil {
		return err
	}

	for i := len(t.incrs) - 1; i >= 0; i-- {
		incr := t.incrs[i]

		if err := incr.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	if err := t.cond.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	for i := len(t.inits) - 1; i >= 0; i-- {
		init := t.inits[i]

		if err := init.ResolveExpressionActivity(usage); err != nil {
			return err
		}

		/*
			XXX: why was this in forward order?
			switch init := init_.(type) {
			case *Assign:
				lhs, err := init.GetLhsVarExpression()
				if err != nil {
					panic(err)
				}

				if err := usage.Rereference(lhs.ref, lhs.Context()); err != nil {
					return err
				}

				if err := init.rhs.ResolveExpressionActivity(usage); err != nil {
					return err
				}
			default:
				panic("expected assign")
			}
		*/
	}

	return nil
}

func (t *For) UniversalStatementNames(ns Namespace) error {
	for _, init := range t.inits {
		if err := init.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	if err := t.cond.UniversalExpressionNames(ns); err != nil {
		return err
	}

	for _, incr := range t.incrs {
		if err := incr.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return t.Block.UniversalStatementNames(ns)
}

func (t *For) UniqueStatementNames(ns Namespace) error {
	subNs := ns.NewBlockNamespace()

	for _, init_ := range t.inits {
		switch init := init_.(type) {
		case *Assign:
			lhs, err := init.GetLhsVarExpression()
			if err != nil {
				panic(err)
			}

			switch t.varType {
			case LET, CONST:
				subNs.LetName(lhs.ref)
			case VAR:
				ns.VarName(lhs.ref)
			default:
				panic("unexpected")
			}

			if err := init.rhs.UniqueExpressionNames(ns); err != nil {
				return err
			}
		default:
			panic("unexpected")
		}
	}

	if err := t.cond.UniqueExpressionNames(subNs); err != nil {
		return err
	}

	for _, incr := range t.incrs {
		if err := incr.UniqueExpressionNames(subNs); err != nil {
			return err
		}
	}

	return t.Block.UniqueStatementNames(subNs)
}

func (t *For) Walk(fn WalkFunc) error {
  for _, init := range t.inits {
    if err := init.Walk(fn); err != nil {
      return err
    }
  }

  if err := t.cond.Walk(fn); err != nil {
    return err
  }

  for _, incr := range t.incrs {
    if err := incr.Walk(fn); err != nil {
      return err
    }
  }

  if err := t.ForBlock.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
