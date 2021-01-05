package prototypes

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type WebGLRenderingContext struct {
  BuiltinPrototype
}

func NewWebGLRenderingContextPrototype() values.Prototype {
  return &WebGLRenderingContext{newBuiltinPrototype("WebGLRenderingContext")}
}

func NewWebGLRenderingContext(ctx context.Context) values.Value {
  return values.NewInstance(NewWebGLRenderingContextPrototype(), ctx)
}

func (p *WebGLRenderingContext) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*WebGLRenderingContext); ok {
    return nil
  } else {
    return checkParent(p, other_, ctx)
  }
}

func IsWebGLRenderingContext(v values.Value) bool {
  ctx := context.NewDummyContext()

  checkVal := NewWebGLRenderingContext(ctx)

  return checkVal.Check(v, ctx) == nil
}

func (p *WebGLRenderingContext) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  b := NewBoolean(ctx)
  f := NewNumber(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)
  enum := NewGLEnum(ctx)

  switch key {
  case "ACTIVE_ATTRIBUTES", "ACTIVE_UNIFORMS", "ALPHA", "ALWAYS", "ARRAY_BUFFER", "ATTACHED_SHADERS", "BLEND", "BLEND_COLOR", "COMPILE_STATUS", "CONTEXT_LOST_WEBGL", "CULL_FACE", "DELETE_STATUS", "DEPTH_TEST", "DITHER", "DST_ALPHA", "DST_COLOR", "DYNAMIC_DRAW", "ELEMENT_ARRAY_BUFFER", "EQUAL", "FLOAT", "FRAGMENT_SHADER", "GEQUAL", "GREATER", "INVALID_ENUM", "INVALID_VALUE", "INVALID_OPERATION", "INVALID_FRAMEBUFFER_OPERATION", "LEQUAL", "LESS", "LINES", "LINE_LOOP", "LINE_STRIP", "LINK_STATUS", "LUMINANCE", "LUMINANCE_ALPHA", "MAX_TEXTURE_IMAGE_UNITS", "MAX_COMBINED_TEXTURE_IMAGE_UNITS", "MAX_VERTEX_TEXTURE_IMAGE_UNITS", "MAX_VERTEX_UNIFORM_VECTORS", "MAX_FRAGMENT_UNIFORM_VECTORS", "NEVER", "NO_ERROR", "NOTEQUAL", "ONE", "ONE_MINUS_DST_ALPHA", "ONE_MINUS_DST_COLOR", "ONE_MINUS_SRC_ALPHA", "ONE_MINUS_SRC_COLOR", "OUT_OF_MEMORY", "POINTS", "POLYGON_OFFSET_FILL", "RGB", "RGBA", "SAMPLE_ALPHA_TO_COVERAGE", "SAMPLE_COVERAGE", "SCISSOR_TEST", "SHADER_TYPE", "SRC_ALPHA", "SRC_COLOR", "STATIC_DRAW", "STENCIL_TEST", "STREAM_DRAW", "TEXTURE_2D", "TEXTURE_MAG_FILTER", "TEXTURE_MIN_FILTER", "TEXTURE_WRAP_S", "TEXTURE_WRAP_T", "TRIANGLES", "TRIANGLE_FAN", "TRIANGLE_STRIP", "UNSIGNED_BYTE", "UNSIGNED_INT", "UNSIGNED_SHORT", "VALIDATE_STATUS", "VERTEX_SHADER", "ZERO":
    return NewNamedGLEnum(key, ctx), nil
  case "CLAMP_TO_BORDERS", "CLAMP_TO_EDGE", "COLOR_BUFFER_BIT", "DEPTH_BUFFER_BIT", "LINEAR", "MIRRORED_REPEAT", "NEAREST", "REPEAT", "TEXTURE0", "TEXTURE1", "TEXTURE2", "TEXTURE3", "TEXTURE4", "TEXTURE5", "TEXTURE6", "TEXTURE7", "TEXTURE8", "TEXTURE9", "TEXTURE10", "TEXTURE11", "TEXTURE12", "TEXTURE13", "TEXTURE14", "TEXTURE15":
    return i, nil
  case "activeTexture", "clear", "enableVertexAttribArray":
    return values.NewFunction([]values.Value{i, nil}, ctx), nil
  case "attachShader":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), NewWebGLShader(ctx), nil}, ctx), nil
  case "bindBuffer":
    return values.NewFunction([]values.Value{enum, NewWebGLBuffer(ctx), nil}, ctx), nil
  case "bindTexture":
    return values.NewFunction([]values.Value{enum, NewWebGLTexture(ctx), nil}, ctx), nil
  case "blendFunc":
    return values.NewFunction([]values.Value{enum, enum, nil}, ctx), nil
  case "blendFuncSeparate":
    return values.NewFunction([]values.Value{enum, enum, enum, enum, nil}, ctx), nil
  case "bufferData":
    return values.NewFunction([]values.Value{enum, NewTypedArray(ctx), enum, nil}, ctx), nil
  case "compileShader":
    return values.NewFunction([]values.Value{NewWebGLShader(ctx), nil}, ctx), nil
  case "createBuffer":
    return values.NewFunction([]values.Value{NewWebGLBuffer(ctx)}, ctx), nil
  case "createProgram":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx)}, ctx), nil
  case "createShader":
    return values.NewFunction([]values.Value{enum, NewWebGLShader(ctx)}, ctx), nil
  case "createTexture":
    return values.NewFunction([]values.Value{NewWebGLTexture(ctx)}, ctx), nil
  case "depthFunc", "disable", "enable":
    return values.NewFunction([]values.Value{enum, nil}, ctx), nil
  case "drawArrays":
    return values.NewFunction([]values.Value{enum, i, i, nil}, ctx), nil
  case "drawElements":
    return values.NewFunction([]values.Value{enum, i, enum, i, nil}, ctx), nil
  case "getAttribLocation":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), s, i}, ctx), nil
  case "getError":
    return values.NewFunction([]values.Value{enum}, ctx), nil
  case "getExtension":
    return values.NewFunction([]values.Value{s, NewWebGLExtension(ctx)}, ctx), nil
  case "getParameter":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{NewNamedGLEnum("MAX_FRAGMENT_UNIFORM_VECTORS", ctx), i},
      []values.Value{NewNamedGLEnum("MAX_VERTEX_UNIFORM_VECTORS", ctx), i},
      []values.Value{enum, NewFloat32Array(ctx)},
    }, ctx), nil
  case "getProgramInfoLog":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), s}, ctx), nil
  case "getProgramParameter":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), enum, b}, ctx), nil
  case "getShaderInfoLog":
    return values.NewFunction([]values.Value{NewWebGLShader(ctx), s}, ctx), nil
  case "getShaderParameter":
    return values.NewFunction([]values.Value{NewWebGLShader(ctx), enum, b}, ctx), nil
  case "getUniformLocation":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), s, i}, ctx), nil
  case "linkProgram":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), nil}, ctx), nil
  case "scissor", "viewport", "clearColor":
    return values.NewFunction([]values.Value{f, f, f, f, nil}, ctx), nil
  case "shaderSource":
    return values.NewFunction([]values.Value{NewWebGLShader(ctx), s, nil}, ctx), nil
  case "texParameterf":
    return values.NewFunction([]values.Value{enum, enum, f, nil}, ctx), nil
  case "texParameteri":
    return values.NewFunction([]values.Value{enum, enum, i, nil}, ctx), nil
  case "texImage2D":
    return values.NewOverloadedFunction([][]values.Value{
      []values.Value{enum, i, enum, i, i, i, enum, enum, NewTypedArray(ctx), nil},
      []values.Value{enum, i, enum, enum, enum, NewImageData(ctx), nil},
      []values.Value{enum, i, enum, enum, enum, NewArrayBuffer(ctx), nil},
      []values.Value{enum, i, enum, enum, enum, NewHTMLCanvasElement(ctx), nil},
    }, ctx), nil
  case "uniform1f":
    return values.NewFunction([]values.Value{i, f, nil}, ctx), nil
  case "uniform2f":
    return values.NewFunction([]values.Value{i, f, f, nil}, ctx), nil
  case "uniform3f":
    return values.NewFunction([]values.Value{i, f, f, f, nil}, ctx), nil
  case "uniform4f":
    return values.NewFunction([]values.Value{i, f, f, f, f, nil}, ctx), nil
  case "uniform1i":
    return values.NewFunction([]values.Value{i, i, nil}, ctx), nil
  case "uniform2i":
    return values.NewFunction([]values.Value{i, i, i, nil}, ctx), nil
  case "uniform3i":
    return values.NewFunction([]values.Value{i, i, i, i, nil}, ctx), nil
  case "uniform4i":
    return values.NewFunction([]values.Value{i, i, i, i, i, nil}, ctx), nil
  case "uniform1fv", "uniform2fv", "uniform3fv", "uniform4fv":
    return values.NewFunction([]values.Value{i, NewArray(f, ctx), nil}, ctx), nil
  case "uniform1iv", "uniform2iv", "uniform3iv", "uniform4iv":
    return values.NewFunction([]values.Value{i, NewArray(i, ctx), nil}, ctx), nil
  case "uniformMatrix2fv", "uniformMatrix3fv", "uniformMatrix4fv":
    return values.NewFunction([]values.Value{i, b, NewArray(f, ctx), nil}, ctx), nil
  case "useProgram", "validateProgram":
    return values.NewFunction([]values.Value{NewWebGLProgram(ctx), nil}, ctx), nil
  case "vertexAttribPointer":
    return values.NewOverloadedMethodLikeFunction([][]values.Value{
      []values.Value{i, i, enum, b, i, i, b},
    }, ctx), nil
  default:
    return nil, nil
  }
}

func (p *WebGLRenderingContext) GetClassValue() (*values.Class, error) {
  ctx := p.Context()
  return values.NewUnconstructableClass(NewWebGLRenderingContextPrototype(), ctx), nil
}
