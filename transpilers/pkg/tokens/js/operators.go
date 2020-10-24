package js

import (
	"strconv"
	"strings"

	"./prototypes"
	"./values"

	"../context"
	"../patterns"
)

var _binaryPrecedenceMap = map[string]int{
	".":          19, // not actually used, treated separately
	"**":         16,
	"*":          15,
	"/":          15,
	"%":          15,
	"+":          14,
	"-":          14,
	"<<":         13,
	">>":         13,
	">>>":        13,
	"<":          12,
	"<=":         12,
	">":          12,
	">=":         12,
	"in":         12,
	"instanceof": 12,
	"==":         11,
	"!=":         11,
	"===":        11,
	"!==":        11,
	"&":          10,
	"^":          9,
	"|":          8,
	"&&":         6,
	"||":         5,
	"=":          3,
}

var _preUnaryPrecedenceMap = map[string]int{
	"new":    19,
	"!":      17,
	"-":      17,
	"~":      17,
	"+":      17,
	"++":     17,
	"--":     17,
	"typeof": 17,
	"delete": 17,
	"await":  17,
}

var _postUnaryPrecedenceMap = map[string]int{
	"++": 18,
	"--": 18,
}

var _ternaryPrecedenceMap = map[string]int{
	"? :": 4,
}

type Op interface {
	Args() []Token
	Precedence() int
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
		return NewInstanceOf(a, b, ctx), nil
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

func (t *BinaryOp) Precedence() int {
	return _binaryPrecedenceMap[t.op]
}

func (t *TernaryOp) Precedence() int {
	return _ternaryPrecedenceMap[t.op0+" "+t.op1]
}

func (t *PreUnaryOp) Precedence() int {
	return _preUnaryPrecedenceMap[t.op]
}

func (t *PostUnaryOp) Precedence() int {
	return _postUnaryPrecedenceMap[t.op]
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

func (t *PostUnaryOp) HoistValues(stack values.Stack) error {
	return nil
}

func (t *PreUnaryOp) AddStatement(st Statement) {
	panic("not a block")
}

func (t *PreUnaryOp) HoistNames(scope Scope) error {
	return nil
}

func (t *PreUnaryOp) HoistValues(stack values.Stack) error {
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
		if aOp.Precedence() < t.Precedence() {
			panic("unexpected")
			b.WriteString("(")
			b.WriteString(aOp.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(aOp.WriteExpression())
		}
	} else {
		b.WriteString(t.a.WriteExpression())
	}
	return b.String()
}

func (t *PostUnaryOp) WriteExpression() string {
	var b strings.Builder

	if aOp, ok := t.a.(Op); ok {
		if aOp.Precedence() < t.Precedence() {
			panic("unexpected")
			b.WriteString("(")
			b.WriteString(aOp.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(aOp.WriteExpression())
		}
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
		if aOp.Precedence() < t.Precedence() {
			panic("unexpected")
			b.WriteString("(")
			b.WriteString(aOp.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(aOp.WriteExpression())
		}
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
		if bOp.Precedence() < t.Precedence() {
			errCtx := t.Context()
			panic(errCtx.NewError("unexpected precedence for b (" + strconv.Itoa(bOp.Precedence()) + " is less than " + strconv.Itoa(t.Precedence()) + ")").Error())
			b.WriteString("(")
			b.WriteString(bOp.WriteExpression())
			b.WriteString(")")
		} else {
			b.WriteString(bOp.WriteExpression())
		}
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

func (t *PostIncrOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	if !a.IsInstanceOf(prototypes.Int) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	ctx := t.Context()
	result := prototypes.NewInt(ctx)

	switch lhs := t.a.(type) {
	case *VarExpression:
		if err := stack.SetValue(lhs.ref, result, true, ctx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (t *PostDecrOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	if !a.IsInstanceOf(prototypes.Int) {
		errCtx := t.a.Context()
		return nil, errCtx.NewError("Error: expected Int, got " + a.TypeName())
	}

	ctx := t.Context()

	result := prototypes.NewInt(ctx)

	switch lhs := t.a.(type) {
	case *VarExpression:
		if err := stack.SetValue(lhs.ref, result, true, ctx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (t *PostIncrOp) EvalStatement(stack values.Stack) error {
	_, err := t.EvalExpression(stack)
	return err
}

func (t *PostIncrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *PostDecrOp) EvalStatement(stack values.Stack) error {
	_, err := t.EvalExpression(stack)
	return err
}

func (t *PostDecrOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *DeleteOp) EvalStatement(stack values.Stack) error {
	_, err := t.EvalExpression(stack)
	return err
}

func (t *TernaryOp) evalArgs(stack values.Stack) (values.Value, values.Value, values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, nil, nil, err
	}

	b, err := t.b.EvalExpression(stack)
	if err != nil {
		return nil, nil, nil, err
	}

	c, err := t.c.EvalExpression(stack)
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

func (t *BinaryOp) evalArgs(stack values.Stack) (values.Value, values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, nil, err
	}

	b, err := t.b.EvalExpression(stack)
	if err != nil {
		return nil, nil, err
	}

	return a, b, nil
}

func (t *BinaryOp) Walk(fn WalkFunc) error {
  if err := t.a.Walk(fn); err != nil {
    return err
  }
  
  return t.b.Walk(fn)
}

func (t *UnaryOp) evalArg(stack values.Stack) (values.Value, error) {
	if a, err := t.a.EvalExpression(stack); err != nil {
		return nil, err
	} else {
		return a, nil
	}
}

func (t *UnaryOp) Walk(fn WalkFunc) error {
  return t.a.Walk(fn)
}

func (t *NewOp) EvalExpression(stack values.Stack) (values.Value, error) {
	call, ok := t.a.(*Call)
	if !ok {
		panic("expected call")
	}

	lhsCallValue, err := call.lhs.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	args, err := call.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	return lhsCallValue.EvalConstructor(stack, args, t.Context())
}

func (t *NewOp) Walk(fn WalkFunc) error {
  if err := t.UnaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *DeleteOp) EvalExpression(stack values.Stack) (values.Value, error) {
	switch t.a.(type) {
	case *Member, *Index:
	default:
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Member of Index rhs to delete")
	}

	if _, err := t.a.EvalExpression(stack); err != nil {
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

func (t *TypeOfOp) EvalExpression(stack values.Stack) (values.Value, error) {
	if _, err := t.a.EvalExpression(stack); err != nil {
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

func (t *InOp) EvalExpression(stack values.Stack) (values.Value, error) {
	if _, _, err := t.BinaryOp.evalArgs(stack); err != nil {
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

func (t *AddOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	aStr, aOk := a.LiteralStringValue()
	bStr, bOk := b.LiteralStringValue()
	if aOk && bOk {
		// literal concatenation
		return prototypes.NewLiteralString(aStr+bStr, t.Context()), nil
	}

	switch {
	case values.IsAllNull(a) && values.IsAllNull(b):
		return values.NewAllNull(ctx), nil
	case values.IsAllNull(a):
		if !b.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) {
			return nil, ctx.NewError("Error: expected String, Number, of Boolean for second argument")
		}
		return values.NewContextValue(b.RemoveLiteralness(true), ctx), nil
	case values.IsAllNull(b):
		if !a.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) {
			return nil, ctx.NewError("Error: expected String, Number, of Boolean for first argument")
		}
		return values.NewContextValue(a.RemoveLiteralness(true), ctx), nil
	case a.IsInstanceOf(prototypes.String):
		if !b.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) {
			return nil, ctx.NewError("Error: expected String for second argument (hint: first argument is String)")
		}
		return prototypes.NewString(ctx), nil
	case b.IsInstanceOf(prototypes.String):
		if !a.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) {
			return nil, ctx.NewError("Error: expected String for first argument (hint: second argument is String)")
		}
		return prototypes.NewString(ctx), nil
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
		return prototypes.NewNumber(ctx), nil
	case a.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) &&
		b.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean):
		return prototypes.NewString(ctx), nil
	default:
		return nil, ctx.NewError("Error: invalid operands for '+' operator" +
			" (expected two Numbers, or a String and String/Boolean/Number, got " +
			a.TypeName() + " and " + b.TypeName() + ")")
	}
}

func (t *AddOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *SubOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case values.IsAllNull(a) && values.IsAllNull(b):
		return values.NewAllNull(ctx), nil
	case values.IsAllNull(a):
		if !b.IsInstanceOf(prototypes.Number) {
			return nil, ctx.NewError("Error: expected Number for second argument")
		}
		return values.NewContextValue(b, ctx), nil
	case values.IsAllNull(b):
		if !a.IsInstanceOf(prototypes.Number) {
			return nil, ctx.NewError("Error: expected Number for first argument")
		}
		return values.NewContextValue(a, ctx), nil
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
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

func (t *DivOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
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

func (t *MulOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
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

func (t *RemainderOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
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

func (t *PowOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
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
func (t *OrderCompareOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
		return prototypes.NewBoolean(ctx), nil
	case a.IsInstanceOf(prototypes.String) && b.IsInstanceOf(prototypes.String):
		return prototypes.NewBoolean(ctx), nil
	case a.IsInstanceOf(prototypes.Boolean) && b.IsInstanceOf(prototypes.Boolean):
		return prototypes.NewBoolean(ctx), nil
	default:
		return nil, ctx.NewError("Error: expected a 2 Numbers, 2 Strings or 2 Booleans" +
			" (got " + a.TypeName() + " and " + b.TypeName() + ")")
	}
}

// TODO: implement for specific types
func (t *OrderCompareOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *EqCompareOp) EvalExpression(stack values.Stack) (values.Value, error) {
	_, _, err := t.BinaryOp.evalArgs(stack)
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

func (t *StrictEqOp) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := t.Context()

	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	// return literal bool if both a and b are null
	if _, ok := t.a.(*LiteralNull); ok {
		if b.IsNull() && !values.IsAllNull(b) {
			return prototypes.NewLiteralBoolean(true, ctx), nil
		}
	} else if _, ok := t.b.(*LiteralNull); ok {
		if a.IsNull() && !values.IsAllNull(a) {
			return prototypes.NewLiteralBoolean(true, ctx), nil
		}
	}

	return prototypes.NewBoolean(ctx), nil
}

func (t *StrictEqOp) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }
  
  return fn(t)
}

func (t *StrictNEOp) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := t.Context()

	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	// return literal bool if both a and b are null
	if _, ok := t.a.(*LiteralNull); ok {
		if b.IsNull() { //&& !values.IsAllNull(b) {
			return prototypes.NewLiteralBoolean(false, ctx), nil
		}
	} else if _, ok := t.b.(*LiteralNull); ok {
		if a.IsNull() { //&& !values.IsAllNull(a) {
			return prototypes.NewLiteralBoolean(false, ctx), nil
		}
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

func (t *PreIncrOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	if !a.IsInstanceOf(prototypes.Int) {
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

func (t *PreDecrOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	if !a.IsInstanceOf(prototypes.Int) {
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

func (t *NegOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean, prototypes.Int):
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

func (t *PosOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	switch {
	case a.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean, prototypes.Int):
		return prototypes.NewNumber(ctx), nil
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

func (t *BinaryBitOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, b, err := t.BinaryOp.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !(a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int)) {
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

func (t *BitNotOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !a.IsInstanceOf(prototypes.Int) {
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

func (t *LogicalNotOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !a.IsInstanceOf(prototypes.Boolean) {
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


func (t *LogicalBinaryOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	b, err := t.b.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	// also allow two numbers (to absorb nans, nulls etc)
	switch {
	case a.IsInstanceOf(prototypes.Boolean) && b.IsInstanceOf(prototypes.Boolean):
		return prototypes.NewBoolean(ctx), nil
	case a.IsInstanceOf(prototypes.Int) && b.IsInstanceOf(prototypes.Int):
		return prototypes.NewInt(ctx), nil
	case a.IsInstanceOf(prototypes.Number) && b.IsInstanceOf(prototypes.Number):
		return prototypes.NewNumber(ctx), nil
	default:
		err := ctx.NewError("Error: expected two Booleans, or two Numbers")
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

func (t *LogicalAndOp) CollectTypeGuards(stack values.Stack, c map[interface{}]values.Interface) (bool, error) {
	// eval expression isn't strictly necessary, but the shortcircuit Boolean/Int/Number evals might give an error
	if _, err := t.EvalExpression(stack); err != nil {
		return false, err
	}

	if a, ok := t.a.(TypeGuard); ok {
		if b, ok := t.b.(TypeGuard); ok {
			ok, err := a.CollectTypeGuards(stack, c)
			if err != nil {
				return false, err
			}

			if ok {
				ok, err := b.CollectTypeGuards(stack, c)
				if err != nil {
					return false, err
				}

				return ok, nil
			}
		}
	}

	return false, nil
}

func (t *IfElseOp) EvalExpression(stack values.Stack) (values.Value, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	b, err := t.b.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	c, err := t.c.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	ctx := t.Context()

	if !a.IsInstanceOf(prototypes.Boolean) {
		return nil, ctx.NewError("Error: expected Boolean first argument, got " + a.TypeName())
	}

	return values.NewMulti([]values.Value{b, c}, ctx), nil
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
