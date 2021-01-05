package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Scalar struct {
  typeName string
  ValueData
}

func newScalar(typeName string, ctx context.Context) Scalar {
  return Scalar{typeName, newValueData(ctx)}
}

func NewScalar(typeName string, ctx context.Context) Value {
  sc := newScalar(typeName, ctx)

  return &sc
}

func NewFloat(ctx context.Context) Value {
  return NewScalar("float", ctx)
}

func NewInt(ctx context.Context) Value {
  return NewScalar("int", ctx)
}

func NewBool(ctx context.Context) Value {
  return NewScalar("bool", ctx)
}

func (v *Scalar) TypeName() string {
  return v.typeName
}

func (v *Scalar) Check(other_ Value, ctx context.Context) error {
  other_ = UnpackContextValue(other_)

  if other_.TypeName() == v.TypeName() {
    return nil
  }

  return ctx.NewError("Error: expected " + v.TypeName() + ", got " + other_.TypeName())
}

func (v *Scalar) EvalFunction(args []Value, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: not a function")
}

func (v *Scalar) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: " + v.typeName + "." + key + " not found")
}

func (v *Scalar) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: " + v.typeName + "." + key + " not found")
}

func (v *Scalar) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of " + v.typeName)
}

func (v *Scalar) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of " + v.typeName)
}

func (v *Scalar) Length() int {
  return 1
}

func isScalar(v_ Value, typeName string) bool {
  v_ = UnpackContextValue(v_)

  if v, ok := v_.(*Scalar); ok {
    return v.TypeName() == typeName
  }

  return false
}

func assertScalar(v_ Value, typeName string) (*Scalar, error) {
  errCtx := v_.Context()
  v_ = UnpackContextValue(v_)

  if v, ok := v_.(*Scalar); ok {
    if v.TypeName() == typeName {
      return v, nil
    } 
  }

  return nil, errCtx.NewError("Error: expected " + typeName + ", got " + v_.TypeName())
}

func IsScalar(v_ Value) bool {
  v_ = UnpackContextValue(v_)

  if _, ok := v_.(*Scalar); ok {
    return true
  } else if _, ok := v_.LiteralIntValue(); ok {
    return true
  } else {
    return false
  }
}

func IsBool(v_ Value) bool {
  return isScalar(v_, "bool")
}

func AssertBool(v_ Value) (*Scalar, error) {
  return assertScalar(v_, "bool")
}

func IsInt(v_ Value) bool {
  _, isLit := v_.LiteralIntValue()

  return isScalar(v_, "int") || isLit
}

func AssertInt(v_ Value) (*Scalar, error) {
  _, isLit := v_.LiteralIntValue()
  if isLit {
    v, err := AssertLiteralInt(v_)
    if err != nil {
      return nil, err
    }

    return &v.Scalar, nil
  } else {
    return assertScalar(v_, "int")
  }
}

func IsFloat(v_ Value) bool {
  return isScalar(v_, "float")
}

func AssertFloat(v_ Value) (*Scalar, error) {
  return assertScalar(v_, "float")
}

func IsSimple(v_ Value) bool {
  return IsScalar(v_) || IsVec(v_)
}
