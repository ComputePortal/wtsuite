package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type UnaryOp struct {
  op string
  a Expression
  TokenData
}

type PreUnaryOp struct {
  UnaryOp
}

type PostUnaryOp struct {
  UnaryOp
}

type BinaryOp struct {
  op string
  a, b Expression
  TokenData
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

type NegOp struct {
  PreUnaryOp
}

type PosOp struct {
  PreUnaryOp
}

type NotOp struct {
  PreUnaryOp
}

type LogicalBinaryOp struct {
  BinaryOp
}

type AndOp struct {
  LogicalBinaryOp
}

type OrOp struct {
  LogicalBinaryOp
}

type XorOp struct {
  LogicalBinaryOp
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

type PostIncrOp struct {
  PostUnaryOp
}

type PostDecrOp struct {
  PostUnaryOp
}

func newUnaryOp(op string, a Expression, ctx context.Context) UnaryOp {
  return UnaryOp{op, a, newTokenData(ctx)}
}

func newPreUnaryOp(op string, a Expression, ctx context.Context) PreUnaryOp {
  return PreUnaryOp{newUnaryOp(op, a, ctx)}
}

func newPostUnaryOp(op string, a Expression, ctx context.Context) PostUnaryOp {
  return PostUnaryOp{newUnaryOp(op, a, ctx)}
}

func newBinaryOp(op string, a, b Expression, ctx context.Context) BinaryOp {
  return BinaryOp{op, a, b, newTokenData(ctx)}
}

func NewAddOp(a, b Expression, ctx context.Context) *AddOp {
  return &AddOp{newBinaryOp("+", a, b, ctx)}
}

func NewSubOp(a, b Expression, ctx context.Context) *SubOp {
  return &SubOp{newBinaryOp("-", a, b, ctx)}
}

func NewDivOp(a, b Expression, ctx context.Context) *DivOp {
  return &DivOp{newBinaryOp("/", a, b, ctx)}
}

func NewMulOp(a, b Expression, ctx context.Context) *MulOp {
  return &MulOp{newBinaryOp("*", a, b, ctx)}
}

func NewNegOp(a Expression, ctx context.Context) *NegOp {
  return &NegOp{newPreUnaryOp("-", a, ctx)}
}

func NewPosOp(a Expression, ctx context.Context) *PosOp {
  return &PosOp{newPreUnaryOp("+", a, ctx)}
}

func NewNotOp(a Expression, ctx context.Context) *NotOp {
  return &NotOp{newPreUnaryOp("!", a, ctx)}
}

func newLogicalBinaryOp(op string, a, b Expression, ctx context.Context) LogicalBinaryOp {
  return LogicalBinaryOp{newBinaryOp(op, a, b, ctx)}
}

func NewAndOp(a, b Expression, ctx context.Context) *AndOp {
  return &AndOp{newLogicalBinaryOp("&&", a, b, ctx)}
}

func NewOrOp(a, b Expression, ctx context.Context) *OrOp {
  return &OrOp{newLogicalBinaryOp("||", a, b, ctx)}
}

func NewXorOp(a, b Expression, ctx context.Context) *XorOp {
  return &XorOp{newLogicalBinaryOp("^^", a, b, ctx)}
}

func newOrderCompareOp(op string, a, b Expression, ctx context.Context) OrderCompareOp {
  return OrderCompareOp{newBinaryOp(op, a, b, ctx)}
}

func NewLTOp(a, b Expression, ctx context.Context) *LTOp {
  return &LTOp{newOrderCompareOp("<", a, b, ctx)}
}

func NewGTOp(a, b Expression, ctx context.Context) *GTOp {
  return &GTOp{newOrderCompareOp(">", a, b, ctx)}
}

func NewLEOp(a, b Expression, ctx context.Context) *LEOp {
  return &LEOp{newOrderCompareOp("<=", a, b, ctx)}
}

func NewGEOp(a, b Expression, ctx context.Context) *GEOp {
  return &GEOp{newOrderCompareOp(">=", a, b, ctx)}
}

func newEqCompareOp(op string, a, b Expression, ctx context.Context) EqCompareOp {
  return EqCompareOp{newBinaryOp(op, a, b, ctx)}
}

func NewEqOp(a, b Expression, ctx context.Context) *EqOp {
  return &EqOp{newEqCompareOp("==", a, b, ctx)}
}

func NewNEOp(a, b Expression, ctx context.Context) *NEOp {
  return &NEOp{newEqCompareOp("!=", a, b, ctx)}
}

func NewPostIncrOp(a Expression, ctx context.Context) *PostIncrOp {
  return &PostIncrOp{newPostUnaryOp("++", a, ctx)}
}

func NewPostDecrOp(a Expression, ctx context.Context) *PostDecrOp {
  return &PostDecrOp{newPostUnaryOp("--", a, ctx)}
}

// dump functions
func (t *UnaryOp) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(t.op)
  b.WriteString("\n")
  b.WriteString(t.a.Dump(indent + "  "))

  return b.String()
}

// dump functions
func (t *BinaryOp) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(t.op)
  b.WriteString("\n")
  b.WriteString(t.a.Dump(indent + "  "))
  b.WriteString(t.b.Dump(indent + "  "))

  return b.String()
}

func (t *PreUnaryOp) WriteExpression() string {
  var b strings.Builder

  b.WriteString(t.op)
  b.WriteString(t.a.WriteExpression())

  return b.String()
}

func (t *PostUnaryOp) WriteExpression() string {
  var b strings.Builder

  b.WriteString(t.a.WriteExpression())
  b.WriteString(t.op)

  return b.String()
}

func (t *BinaryOp) WriteExpression() string {
  var b strings.Builder

  b.WriteString(t.a.WriteExpression())
  b.WriteString(t.op)
  b.WriteString(t.b.WriteExpression())

  return b.String()
}

func (t *PostUnaryOp) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder
  b.WriteString(indent)
  b.WriteString(t.a.WriteExpression())
  b.WriteString(t.op)

  return b.String()
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

func (t *PostUnaryOp) ResolveStatementNames(scope Scope) error {
  return t.ResolveExpressionNames(scope)
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

  return a, b, nil
}

func (t *AddOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a):
    if _, err := values.AssertInt(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  case values.IsFloat(a):
    if _, err := values.AssertFloat(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't add " + a.TypeName() + " and " + b.TypeName())
  }
}

func (t *SubOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a):
    if _, err := values.AssertInt(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  case values.IsFloat(a):
    if _, err := values.AssertFloat(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't subtract " + a.TypeName() + " and " + b.TypeName())
  }
}

func (t *DivOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsFloat(a):
    if _, err := values.AssertFloat(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't divide " + a.TypeName() + " and " + b.TypeName())
  }
}

func (t *MulOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a):
    if _, err := values.AssertInt(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  case values.IsFloat(a):
    if _, err := values.AssertFloat(b); err != nil {
      return nil, err
    }
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't multiply " + a.TypeName() + " and " + b.TypeName())
  }
}

func (t *NegOp) EvalExpression() (values.Value, error) {
  a, err := t.a.EvalExpression()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a):
    return values.NewContextValue(a, t.Context()), nil
  case values.IsFloat(a):
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't negate " + a.TypeName())
  }
}

func (t *PosOp) EvalExpression() (values.Value, error) {
  a, err := t.a.EvalExpression()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a):
    return values.NewContextValue(a, t.Context()), nil
  case values.IsFloat(a):
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't + " + a.TypeName())
  }
}

func (t *NotOp) EvalExpression() (values.Value, error) {
  a, err := t.a.EvalExpression()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsBool(a):
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't + " + a.TypeName())
  }
}

func (t *LogicalBinaryOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsBool(a) && values.IsBool(b):
    return values.NewContextValue(a, t.Context()), nil
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: " + a.TypeName() + t.op + b.TypeName() + " is illegal")
  }
}

func (t *OrderCompareOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  switch {
  case values.IsInt(a) && values.IsInt(b):
  case values.IsFloat(a) && values.IsFloat(b):
  default:
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: " + a.TypeName() + t.op + b.TypeName() + " is illegal")
  }

  return values.NewScalar("bool", t.Context()), nil
}

func (t *EqCompareOp) EvalExpression() (values.Value, error) {
  a, b, err := t.evalArgs()
  if err != nil {
    return nil, err
  }

  if err := a.Check(b, t.Context()); err != nil {
    return nil, err
  }

  return values.NewScalar("bool", t.Context()), nil
}

func (t *PostIncrOp) EvalStatement() error {
  a, err := t.a.EvalExpression()
  if err != nil {
    return err
  }

  if _, err := values.AssertInt(a); err != nil {
    return err
  }

  return nil
}

func (t *PostDecrOp) EvalStatement() error {
  a, err := t.a.EvalExpression()
  if err != nil {
    return err
  }

  if _, err := values.AssertInt(a); err != nil {
    return err
  }

  return nil
}

func (t *UnaryOp) ResolveExpressionActivity(usage Usage) error {
  return t.a.ResolveExpressionActivity(usage)
}

func (t *BinaryOp) ResolveExpressionActivity(usage Usage) error {
  if err := t.a.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  return t.b.ResolveExpressionActivity(usage)
}

func (t *PostUnaryOp) ResolveStatementActivity(usage Usage) error {
  return t.ResolveExpressionActivity(usage)
}

func (t *PostUnaryOp) UniqueStatementNames(ns Namespace) error {
  return nil
}
