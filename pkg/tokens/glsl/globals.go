package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

var TARGET = "vertex"
 
func handleRegistrationError(err error) {
  if err != nil {
    if TARGET != "all" {
      panic(err)
    }
  }
}
 
func registerValue(scope Scope, name string, constant bool, val values.Value) {
  ctx := context.NewDummyContext()

  variable := NewVariable(name, ctx)
  if constant {
    variable.SetConstant()
  }

  variable.SetValue(val)

  handleRegistrationError(scope.SetVariable(name, variable))
}

func FillCoreScope(scope Scope) {
  ctx := context.NewDummyContext()

  registerValue(scope, "float", true, values.NewScalarType("float", ctx))
  registerValue(scope, "int"  , true, values.NewScalarType("int", ctx))
  registerValue(scope, "bool" , true, values.NewScalarType("bool", ctx))

  registerValue(scope, "vec2", true, values.NewVecType("float", 2, ctx))
  registerValue(scope, "vec3", true, values.NewVecType("float", 3, ctx))
  registerValue(scope, "vec4", true, values.NewVecType("float", 4, ctx))

  registerValue(scope, "bvec2", true, values.NewVecType("bool", 2, ctx))
  registerValue(scope, "bvec3", true, values.NewVecType("bool", 3, ctx))
  registerValue(scope, "bvec4", true, values.NewVecType("bool", 4, ctx))

  registerValue(scope, "ivec2", true, values.NewVecType("int", 2, ctx))
  registerValue(scope, "ivec3", true, values.NewVecType("int", 3, ctx))
  registerValue(scope, "ivec4", true, values.NewVecType("int", 4, ctx))


  // builtin functions
  registerValue(scope, "abs"        , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "acos"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "all"        , true, values.NewAnyAllFunction(ctx))
  registerValue(scope, "any"        , true, values.NewAnyAllFunction(ctx))
  registerValue(scope, "asin"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "atan"       , true, values.NewOneOrTwoToOneFunction(ctx))
  registerValue(scope, "ceil"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "clamp"      , true, values.NewClampFunction(ctx))
  registerValue(scope, "cos"        , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "cross"      , true, values.NewCrossFunction(ctx))
  registerValue(scope, "degrees"    , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "distance"   , true, values.NewDotFunction(ctx))
  registerValue(scope, "dot"        , true, values.NewDotFunction(ctx))
  registerValue(scope, "equal"      , true, values.NewCompareFunction(ctx))
  registerValue(scope, "exp"        , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "exp2"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "faceforward", true, values.NewThreeToOneFunction(ctx))
  registerValue(scope, "floor"      , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "fract"      , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "greaterThan", true, values.NewCompareFunction(ctx))
  registerValue(scope, "greaterThanEqual", true, values.NewCompareFunction(ctx))
  registerValue(scope, "inversesqrt", true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "length"     , true, values.NewLengthFunction(ctx))
  registerValue(scope, "lessThan"   , true, values.NewCompareFunction(ctx))
  registerValue(scope, "lessThanEqual", true, values.NewCompareFunction(ctx))
  registerValue(scope, "log"        , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "log2"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "max"        , true, values.NewMinMaxFunction(ctx))
  registerValue(scope, "min"        , true, values.NewMinMaxFunction(ctx))
  registerValue(scope, "mix"        , true, values.NewMixFunction(ctx))
  registerValue(scope, "mod"        , true, values.NewTwoToOneFunction(ctx))
  registerValue(scope, "normalize"  , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "not"        , true, values.NewNotFunction(ctx))
  registerValue(scope, "notEqual"   , true, values.NewCompareFunction(ctx))
  registerValue(scope, "pow"        , true, values.NewTwoToOneFunction(ctx))
  registerValue(scope, "radians"    , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "reflect"    , true, values.NewTwoToOneFunction(ctx))
  registerValue(scope, "refract"    , true, values.NewTwoToOneFunction(ctx))
  registerValue(scope, "sign"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "sin"        , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "smoothstep" , true, values.NewSmoothStepFunction(ctx))
  registerValue(scope, "sqrt"       , true, values.NewOneToOneFunction(ctx))
  registerValue(scope, "step"       , true, values.NewStepFunction(ctx))
  registerValue(scope, "tan"        , true, values.NewOneToOneFunction(ctx))
}

func FillVertexShaderScope(scope Scope) {
  ctx := context.NewDummyContext()

  FillCoreScope(scope)

  registerValue(scope, "gl_Position", false, values.NewVec("float", 4, ctx))
}

func FillFragmentShaderScope(scope Scope) {
  ctx := context.NewDummyContext()

  FillCoreScope(scope)

  registerValue(scope, "sampler2D", true, values.NewSampler2DType(ctx))
  registerValue(scope, "samplerCube", true, values.NewSamplerCubeType(ctx))

  registerValue(scope, "gl_FragCoord", true, values.NewVec("float", 4, ctx))
  registerValue(scope, "gl_FrontFacing" , true, values.NewScalar("bool", ctx))
  registerValue(scope, "gl_PointCoord" , true, values.NewVec("float", 2, ctx))

  registerValue(scope, "gl_FragColor", false, values.NewVec("float", 4, ctx))

  registerValue(scope, "texture2D"  , true, values.NewTexture2DFunction(ctx))
  registerValue(scope, "textureCube"  , true, values.NewTextureCubeFunction(ctx))
}

func FillGlobalScope(scope Scope) {
  switch TARGET {
  case "vertex":
    FillVertexShaderScope(scope)
  case "fragment":
    FillFragmentShaderScope(scope)
  default:
    panic("unhandled")
  }
}

func NewFilledGlobalScope() *GlobalScopeData {
  scope := &GlobalScopeData{newScopeData(nil)}

  FillGlobalScope(scope)

  return scope
}
