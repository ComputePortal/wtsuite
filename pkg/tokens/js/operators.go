package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
	"../patterns"
)

type Op interface {
	Args() []Token
	Expression
}

type UnaryOp struct {
	op string
	a  Expression
	TokenData
}

type PreUnaryOp struct {
	UnaryOp
}

type PostUnaryOp struct {
	UnaryOp
}

type BinaryOp struct {
	op   string
	a, b Expression
	TokenData
}

type TernaryOp struct {
	op0, op1 string
	a, b, c  Expression
	TokenData
}

// no longer used, but keep the code anyway
type NewOp struct {
	PreUnaryOp
}

type DeleteOp struct {
	PreUnaryOp
}

type TypeOfOp struct {
	PreUnaryOp
}

// InstanceOf if in InstanceOf.go
type InOp struct {
	BinaryOp
}

type AddOp struct {
	BinaryOp
}

type SubOp struct {
	BinaryOp
}

type DivOp struct {
	BinaryOp
}

type MulOp struct {
	BinaryOp
}

type RemainderOp struct {
	BinaryOp
}

type PowOp struct {
	BinaryOp
}

type BinaryBitOp struct {
	BinaryOp
}

type BitAndOp struct {
	BinaryBitOp
}

type BitOrOp struct {
	BinaryBitOp
}

type BitXorOp struct {
	BinaryBitOp
}

type BitNotOp struct {
	PreUnaryOp
}

type ShiftOp struct {
	BinaryBitOp
}

type LeftShiftOp struct {
	ShiftOp
}

type KeepSignRightShiftOp struct {
	ShiftOp
}

type DontKeepSignRightShiftOp struct {
	ShiftOp
}

type OrderCompareOp struct {
	BinaryOp
}

type LTOp struct {
	OrderCompareOp
}

type GTOp struct {
	OrderCompareOp
}

type LEOp struct {
	OrderCompareOp
}

type GEOp struct {
	OrderCompareOp
}

type EqCompareOp struct {
	BinaryOp
}

type EqOp struct {
	EqCompareOp
}

type NEOp struct {
	EqCompareOp
}

type StrictEqOp struct {
	EqCompareOp
}

type StrictNEOp struct {
	EqCompareOp
}

type PostIncrOp struct {
	PostUnaryOp
}

type PostDecrOp struct {
	PostUnaryOp
}

// cannot be used as statement
type PreIncrOp struct {
	PreUnaryOp
}

// cannot be used as statement
type PreDecrOp struct {
	PreUnaryOp
}

type NegOp struct {
	PreUnaryOp
}

type PosOp struct {
	PreUnaryOp
}

type LogicalNotOp struct {
	PreUnaryOp
}

type LogicalBinaryOp struct {
	BinaryOp
}

type LogicalAndOp struct {
	LogicalBinaryOp
}

type LogicalOrOp struct {
	LogicalBinaryOp
}

type IfElseOp struct {
	TernaryOp
}

func NewPostIncrOp(a Expression, ctx context.Context) *PostIncrOp {
	return &PostIncrOp{PostUnaryOp{UnaryOp{"++", a, TokenData{ctx}}}}
}

func NewPostDecrOp(a Expression, ctx context.Context) *PostDecrOp {
	return &PostDecrOp{PostUnaryOp{UnaryOp{"--", a, TokenData{ctx}}}}
}

func NewDeleteOp(a Expression, ctx context.Context) *DeleteOp {
	return &DeleteOp{PreUnaryOp{UnaryOp{"delete", a, TokenData{ctx}}}}
}

func NewBinaryOp(op string, a Expression, b Expression, ctx context.Context) (Op, error) {
	switch {
	case op == ".":
		panic("not handled as an operator")
	case op == ":=":
		panic("not handled as an operator")
	case op == "+":
		return &AddOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "-":
		return &SubOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "/":
		return &DivOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "*":
		return &MulOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "%":
		return &RemainderOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "**":
		return &PowOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "&":
		return &BitAndOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "|":
		return &BitOrOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "||":
		return &LogicalOrOp{LogicalBinaryOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "&&":
		return &LogicalAndOp{LogicalBinaryOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "^":
		return &BitXorOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "<<":
		return &LeftShiftOp{ShiftOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}}, nil
	case op == ">>":
		return &KeepSignRightShiftOp{ShiftOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}}, nil
	case op == ">>>":
		return &DontKeepSignRightShiftOp{ShiftOp{BinaryBitOp{BinaryOp{op, a, b, TokenData{ctx}}}}}, nil
	case op == ">":
		return &GTOp{OrderCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "<":
		return &LTOp{OrderCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "<=":
		return &LEOp{OrderCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == ">=":
		return &GEOp{OrderCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "in":
		return &InOp{BinaryOp{op, a, b, TokenData{ctx}}}, nil
	case op == "instanceof":
		return NewInstanceOf(a, b, ctx)
	case op == "==":
		return &StrictEqOp{EqCompareOp{BinaryOp{"===", a, b, TokenData{ctx}}}}, nil
	case op == "!=":
		return &StrictNEOp{EqCompareOp{BinaryOp{"!==", a, b, TokenData{ctx}}}}, nil
	case op == "===":
		errCtx := ctx
		return nil, errCtx.NewError("Error: use '==' instead (which compiles to ===)")
		//return &StrictEqOp{EqCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case op == "!==":
		errCtx := ctx
		return nil, errCtx.NewError("Error: use '!=' instead (which compiles to !==)")
		//return &StrictNEOp{EqCompareOp{BinaryOp{op, a, b, TokenData{ctx}}}}, nil
	case strings.HasSuffix(op, "="): // must come after other operators that end with an '='
		subOp := strings.TrimSuffix(op, "=")
		return NewAssign(a, b, subOp, ctx), nil
	default:
		return nil, ctx.NewError("Error: binary operator '" + op + "' not supported")
	}
}

func NewPostUnaryOp(op string, a Expression, ctx context.Context) (Op, error) {
	switch op {
	case "++":
		return NewPostIncrOp(a, ctx), nil
	case "--":
		return NewPostDecrOp(a, ctx), nil
	default:
		return nil, ctx.NewError("Error: postfix operator '" + op + "' not supported")
	}
}

func NewPreUnaryOp(op string, a Expression, ctx context.Context) (Op, error) {
	switch op {
	case "++":
		return &PreIncrOp{PreUnaryOp{UnaryOp{"++", a, TokenData{ctx}}}}, nil
	case "--":
		return &PreDecrOp{PreUnaryOp{UnaryOp{"--", a, TokenData{ctx}}}}, nil
	case "-":
		return &NegOp{PreUnaryOp{UnaryOp{"-", a, TokenData{ctx}}}}, nil
	case "+":
		return &PosOp{PreUnaryOp{UnaryOp{"+", a, TokenData{ctx}}}}, nil
	case "~":
		return &BitNotOp{PreUnaryOp{UnaryOp{"~", a, TokenData{ctx}}}}, nil
	case "!":
		return &LogicalNotOp{PreUnaryOp{UnaryOp{"!", a, TokenData{ctx}}}}, nil
	case "new":
		newCtx := context.MergeContexts(ctx, a.Context())
		if _, ok := a.(*Call); !ok {
			errCtx := newCtx
			return nil, errCtx.NewError("Error: new argument is not a function call")
		}

		return &NewOp{PreUnaryOp{UnaryOp{"new", a, TokenData{newCtx}}}}, nil
	case "delete":
		return NewDeleteOp(a, ctx), nil
	case "typeof":
		return &TypeOfOp{PreUnaryOp{UnaryOp{op, a, TokenData{ctx}}}}, nil
	case "await":
		return NewAwait(a, ctx)
	default:
		return nil, ctx.NewError("Error: prefix operator '" + op + "' not supported")
	}
}

func NewTernaryOp(op string, a Expression, b Expression, c Expression, ctx context.Context) (Op, error) {
	switch op {
	case "? :":
		return &IfElseOp{TernaryOp{"?", ":", a, b, c, TokenData{ctx}}}, nil
	default:
		return nil, ctx.NewError("Error: ternary operator '" + op + "' not supported")
	}
}

func (t *PreUnaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("PreUnaryOp(")
	b.WriteString(t.op)
	b.WriteString(")\n")

	b.WriteString(t.a.Dump(indent + "  "))

	return b.String()
}

func (t *PostUnaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("PostUnaryOp(")
	b.WriteString(t.op)
	b.WriteString(")\n")

	b.WriteString(t.a.Dump(indent + "  "))

	return b.String()
}

func (t *PostUnaryOp) AddStatement(st Statement) {
	panic("not a block")
}

func (t *PostUnaryOp) HoistNames(scope Scope) error {
	return nil
}

func (t *PreUnaryOp) AddStatement(st Statement) {
	panic("not a block")
}

func (t *PreUnaryOp) HoistNames(scope Scope) error {
	return nil
}

func (t *BinaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("BinaryOp(")
	b.WriteString(t.op)
	b.WriteString(")\n")

	b.WriteString(t.a.Dump(indent + "  "))
	b.WriteString(t.b.Dump(indent + "  "))

	return b.String()
}

func (t *TernaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("TernaryOp(")
	b.WriteString(t.op0 + " " + t.op1)
	b.WriteString(")\n")

	b.WriteString(t.a.Dump(indent + "  "))
	b.WriteString(t.b.Dump(indent + "  "))
	b.WriteString(t.c.Dump(indent + "  "))

	return b.String()
}

func (t *PreUnaryOp) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.op)

	if patterns.ALPHABET_REGEXP.MatchString(t.op) {
		b.WriteString(" ")
	}

	if aOp, ok := t.a.(Op); ok {
    b.WriteString(aOp.WriteExpression())
	} else {
		b.WriteString(t.a.WriteExpression())
	}
	return b.String()
}

func (t *PostUnaryOp) WriteExpression() string {
	var b strings.Builder

	if aOp, ok := t.a.(Op); ok {
    b.WriteString(aOp.WriteExpression())
	} else {
		b.WriteString(t.a.WriteExpression())
	}

	if patterns.ALPHABET_REGEXP.MatchString(t.op) {
		b.WriteString(" ")
	}
	b.WriteString(t.op)
	return b.String()
}

func (t *TernaryOp) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.a.WriteExpression())
	b.WriteString(t.op0)
	b.WriteString(t.b.WriteExpression())
	b.WriteString(t.op1)
	b.WriteString(t.c.WriteExpression())

	return b.String()
}

func (t *BinaryOp) WriteExpression() string {
	var b strings.Builder
	if aOp, ok := t.a.(Op); ok {
    b.WriteString(aOp.WriteExpression())
	} else {
		b.WriteString(t.a.WriteExpression())
	}

	isWordOp := patterns.ALPHABET_REGEXP.MatchString(t.op)
	if isWordOp {
		b.WriteString(" ")
	}

	b.WriteString(t.op)

	if isWordOp {
		b.WriteString(" ")
	}

	if bOp, ok := t.b.(Op); ok {
    b.WriteString(bOp.WriteExpression())
	} else {
		b.WriteString(t.b.WriteExpression())
	}

	return b.String()
}

func (t *TernaryOp) Args() []Token {
	return []Token{t.a, t.b, t.c}
}

func (t *BinaryOp) Args() []Token {
	return []Token{t.a, t.b}
}

func (t *UnaryOp) Args() []Token {
	return []Token{t.a}
}

func (t *PostIncrOp) WriteStatement(indent string) string {
	return indent + t.a.WriteExpression() + t.op
}

func (t *PostDecrOp) WriteStatement(indent string) string {
	return indent + t.a.WriteExpression() + t.op
}

func (t *DeleteOp) WriteStatement(indent string) string {
	return indent + t.op + " " + t.a.WriteExpression()
}

func (t *TernaryOp) ResolveExpressionNames(scope Scope) error {
	if err := t.a.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.b.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.c.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *BinaryOp) ResolveExpressionNames(scope Scope) error {
	if err := t.a.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.b.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *UnaryOp) ResolveExpressionNames(scope Scope) error {
	if err := t.a.ResolveExpressionNames(scope); err != nil {
		return err
	}

	return nil
}

func (t *PostIncrOp) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *PostDecrOp) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *DeleteOp) ResolveStatementNames(scope Scope) error {
	return t.ResolveExpressionNames(scope)
}

func (t *TernaryOp) ResolveExpressionActivity(usage Usage) error {
	if err := t.a.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	if err := t.b.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	if err := t.c.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *BinaryOp) ResolveExpressionActivity(usage Usage) error {
	if err := t.a.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	if err := t.b.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *UnaryOp) ResolveExpressionActivity(usage Usage) error {
	return t.a.ResolveExpressionActivity(usage)
}

func (t *PostIncrOp) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *PostDecrOp) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *DeleteOp) ResolveStatementActivity(usage Usage) error {
	return t.ResolveExpressionActivity(usage)
}

func (t *TernaryOp) UniversalExpressionNames(ns Namespace) error {
	if err := t.a.UniversalExpressionNames(ns); err != nil {
		return err
	}

	if err := t.b.UniversalExpressionNames(ns); err != nil {
		return err
	}

	if err := t.c.UniversalExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *BinaryOp) UniversalExpressionNames(ns Namespace) error {
	if err := t.a.UniversalExpressionNames(ns); err != nil {
		return err
	}

	if err := t.b.UniversalExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *UnaryOp) UniversalExpressionNames(ns Namespace) error {
	return t.a.UniversalExpressionNames(ns)
}

func (t *PostIncrOp) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *PostDecrOp) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *DeleteOp) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *TernaryOp) UniqueExpressionNames(ns Namespace) error {
	if err := t.a.UniqueExpressionNames(ns); err != nil {
		return err
	}

	if err := t.b.UniqueExpressionNames(ns); err != nil {
		return err
	}

	if err := t.c.UniqueExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *BinaryOp) UniqueExpressionNames(ns Namespace) error {
	if err := t.a.UniqueExpressionNames(ns); err != nil {
		return err
	}

	if err := t.b.UniqueExpressionNames(ns); err != nil {
		return err
	}

	return nil
}

func (t *UnaryOp) UniqueExpressionNames(ns Namespace) error {
	return t.a.UniqueExpressionNames(ns)
}

func (t *PostIncrOp) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *PostDecrOp) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *DeleteOp) UniqueStatementNames(ns Namespace) error {
	return t.UniqueExpressionNames(ns)
}

func (t *PostIncrOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsInt(a) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	ctx := t.Context()
	result := prototypes.NewInt(ctx)

	return result, nil
}

func (t *PostDecrOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsInt(a) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	ctx := t.Context()
	result := prototypes.NewInt(ctx)

	return result, nil
}

func (t *PostIncrOp) EvalStatement() error {
	_, err := t.EvalExpression()
	return err
}

func (t *PostIncrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *PostDecrOp) EvalStatement() error {
	_, err := t.EvalExpression()
	return err
}

func (t *PostDecrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *DeleteOp) EvalStatement() error {
	_, err := t.EvalExpression()
	return err
}

func (t *TernaryOp) evalArgs() (values.Value, values.Value, values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, nil, nil, err
	}

	b, err := t.b.EvalExpression()
	if err != nil {
		return nil, nil, nil, err
	}

	c, err := t.c.EvalExpression()
	if err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, nil
}

func (t *TernaryOp) Walk(fn WalkFunc) error {
  if err := t.a.Walk(fn); err != nil {
    return err
  }
  
  if err := t.b.Walk(fn); err != nil {
    return err
  }

  return t.c.Walk(fn)
}

func (t *BinaryOp) evalArgs() (values.Value, values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, nil, err
	}

	b, err := t.b.EvalExpression()
	if err != nil {
		return nil, nil, err
	}

  if a == nil {
    hereCtx := t.a.Context()
    panic(hereCtx.NewError("a can't be nil").Error())
  } else if b == nil {
    hereCtx := t.b.Context()
    panic(hereCtx.NewError("b can't be nil").Error())
  }

	return a, b, nil
}

func (t *BinaryOp) Walk(fn WalkFunc) error {
  if err := t.a.Walk(fn); err != nil {
    return err
  }
  
  return t.b.Walk(fn)
}

func (t *UnaryOp) evalArg() (values.Value, error) {
	if a, err := t.a.EvalExpression(); err != nil {
		return nil, err
	} else {
		return a, nil
	}
}

func (t *UnaryOp) Walk(fn WalkFunc) error {
  return t.a.Walk(fn)
}

func (t *NewOp) EvalExpression() (values.Value, error) {
	call, ok := t.a.(*Call)
	if !ok {
		panic("expected call")
	}

	lhsCallValue, err := call.lhs.EvalExpression()
	if err != nil {
		return nil, err
	}

	args, err := call.evalArgs()
	if err != nil {
		return nil, err
	}

	return lhsCallValue.EvalConstructor(args, t.Context())
}

func (t *NewOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *DeleteOp) EvalExpression() (values.Value, error) {
	switch t.a.(type) {
	case *Member, *Index:
	default:
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Member or Index rhs to delete")
	}

	if _, err := t.a.EvalExpression(); err != nil {
		return nil, err
	}

	return prototypes.NewBoolean(t.Context()), nil
}

func (t *DeleteOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *TypeOfOp) EvalExpression() (values.Value, error) {
	if _, err := t.a.EvalExpression(); err != nil {
		return nil, err
	}

	// always a string, for any type
	return prototypes.NewString(t.Context()), nil
}

func (t *TypeOfOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *InOp) EvalExpression() (values.Value, error) {
	if _, _, err := t.BinaryOp.evalArgs(); err != nil {
		return nil, err
	}

	return prototypes.NewBoolean(t.Context()), nil
}

func (t *InOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *AddOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

  isString := false
  isInt := false
  isNumber := false

	switch {
	case prototypes.IsString(a):
		if !prototypes.IsStringable(b) {
			return nil, ctx.NewError("Error: expected String for second argument (hint: first argument is String)")
		}
		return prototypes.NewString(ctx), nil
	case prototypes.IsString(b):
		if !prototypes.IsStringable(a) {
			return nil, ctx.NewError("Error: expected String for first argument (hint: second argument is String)")
		}

    isString = true
	case prototypes.IsInt(a) && prototypes.IsInt(b):
    isInt = true
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
    isNumber = true
	case prototypes.IsStringable(a) && prototypes.IsStringable(b):
    isString = true
	default:
		return nil, ctx.NewError("Error: invalid operands for '+' operator" +
			" (expected two Numbers, or a String and String/Boolean/Number, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}

  if isString && (isNumber || isInt) {
    return values.NewAny(ctx), nil
  } else if isNumber {
    return prototypes.NewNumber(ctx), nil
  } else if isInt {
    return prototypes.NewInt(ctx), nil
  } else {
    if !isString {
      panic("unexpected")
    }
    return prototypes.NewString(ctx), nil
  } 
}

func (t *AddOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *SubOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a) && prototypes.IsInt(b):
		return prototypes.NewInt(ctx), nil
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '-' operator" +
			" (expected two Numbers, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *SubOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *DivOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '/' operator" +
			" (expected two Numbers, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *DivOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *MulOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a) && prototypes.IsInt(b):
		return prototypes.NewInt(ctx), nil
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '*' operator" +
			" (expected two Numbers, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *MulOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *RemainderOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a) && prototypes.IsInt(b):
		return prototypes.NewInt(ctx), nil
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '%' operator" +
			" (expected two Numbers, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *RemainderOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *PowOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a) && prototypes.IsInt(b):
		return prototypes.NewInt(ctx), nil
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '**' operator" +
			" (expected two Numbers, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *PowOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

// >=, <=, >, <
func (t *OrderCompareOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
	case prototypes.IsString(a) && prototypes.IsString(b):
	case prototypes.IsBoolean(a) && prototypes.IsBoolean(b):
	default:
		return nil, ctx.NewError("Error: expected a 2 Numbers, 2 Strings or 2 Booleans" +
			" (got " + a.TypeName() + " and " + b.TypeName() + ")")
	}

  return prototypes.NewBoolean(ctx), nil
}

// TODO: implement for specific types
func (t *OrderCompareOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *EqCompareOp) EvalExpression() (values.Value, error) {
	_, _, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()
	return prototypes.NewBoolean(ctx), nil
}

// TODO: implement for specific types
func (t *EqCompareOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *StrictEqOp) EvalExpression() (values.Value, error) {
	ctx := t.Context()

	_, _, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	return prototypes.NewBoolean(ctx), nil
}

func (t *StrictEqOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *StrictNEOp) EvalExpression() (values.Value, error) {
	ctx := t.Context()

	_, _, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	return prototypes.NewBoolean(ctx), nil
}

func (t *StrictNEOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *NewOp) WriteExpression() string {
	return t.PreUnaryOp.WriteExpression()
}

func (t *PreIncrOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsInt(a) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	return a, nil
}

func (t *PreIncrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *PreDecrOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsInt(a) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	return a, nil
}

func (t *PreDecrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *NegOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a):
		return prototypes.NewInt(ctx), nil
  case prototypes.IsNumber(a):
		return prototypes.NewNumber(ctx), nil
	default:
		return nil, ctx.NewError("Error: expected a Number, got " + a.TypeName())
	}
}

func (t *NegOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *PosOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case prototypes.IsInt(a):
		return prototypes.NewInt(ctx), nil
  case prototypes.IsNumber(a):
    return prototypes.NewNumber(ctx), nil
  //case prototypes.IsStringable(a):
		//return prototypes.NewString(ctx), nil
	default:
		return nil, ctx.NewError("Error: expected a Number, got " + a.TypeName())
	}
}

func (t *PosOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *BinaryBitOp) EvalExpression() (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !(prototypes.IsInt(a) && prototypes.IsInt(b)) {
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected two Int arguments," +
			" got " + a.TypeName() + " and " + b.TypeName())
	}

	return prototypes.NewInt(ctx), nil
}

// TODO: implement for each special function
func (t *BinaryBitOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func (t *BitNotOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !prototypes.IsInt(a) {
		return nil, ctx.NewError("Error: expected Int argument, got " + a.TypeName())
	}

	return prototypes.NewInt(ctx), nil
}

func (t *BitNotOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func (t *LogicalNotOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !prototypes.IsBoolean(a) {
		return nil, ctx.NewError("Error: expected Boolean argument, got " + a.TypeName())
	}

	if litVal, ok := a.LiteralBooleanValue(); ok {
		return prototypes.NewLiteralBoolean(!litVal, ctx), nil
	}

	return prototypes.NewBoolean(ctx), nil
}

func (t *LogicalNotOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}


func (t *LogicalBinaryOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	b, err := t.b.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	// also allow two numbers (to absorb nans, nulls etc)
	switch {
	case prototypes.IsBoolean(a) && prototypes.IsBoolean(b):
		return prototypes.NewBoolean(ctx), nil
	case prototypes.IsInt(a) && prototypes.IsInt(b):
		return prototypes.NewInt(ctx), nil
	case prototypes.IsNumber(a) && prototypes.IsNumber(b):
		return prototypes.NewNumber(ctx), nil
	default:
		err := ctx.NewError("Error: expected two Booleans, or two Numbers (got " + a.TypeName() + " and " + b.TypeName() + ")")
		return nil, err
	}
}

func (t *LogicalOrOp) Walk(fn WalkFunc) error {
  if err := t.LogicalBinaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func (t *LogicalAndOp) Walk(fn WalkFunc) error {
  if err := t.LogicalBinaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func (t *LogicalAndOp) CollectTypeGuards(c map[Variable]values.Interface) (bool, error) {
	if _, err := t.EvalExpression(); err != nil {
		return false, err
	}

	if a, ok := t.a.(TypeGuard); ok {
		if b, ok := t.b.(TypeGuard); ok {
			ok, err := a.CollectTypeGuards(c)
			if err != nil {
				return false, err
			}

			if ok {
				ok, err := b.CollectTypeGuards(c)
				if err != nil {
					return false, err
				}

				return ok, nil
			}
		}
	}

	return false, nil
}

func (t *IfElseOp) EvalExpression() (values.Value, error) {
	a, err := t.a.EvalExpression()
	if err != nil {
		return nil, err
	}

	b, err := t.b.EvalExpression()
	if err != nil {
		return nil, err
	}

	c, err := t.c.EvalExpression()
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !prototypes.IsBoolean(a) {
		return nil, ctx.NewError("Error: expected Boolean first argument, got " + a.TypeName())
	}

  common := values.CommonValue([]values.Value{b, c}, t.Context())
	return values.NewContextValue(common, ctx), nil
}

func (t *IfElseOp) Walk(fn WalkFunc) error {
  if err := t.TernaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func IsSimpleLT(t Expression) bool {
	lt, ok := t.(*LTOp)
	if ok {
		return IsVarExpression(lt.a)
	} else {
		return false
	}
}

func IsSimplePostIncr(t Expression) bool {
	pi, ok := t.(*PostIncrOp)

	if ok {
		return IsVarExpression(pi.a)
	} else {
		return false
	}
}
