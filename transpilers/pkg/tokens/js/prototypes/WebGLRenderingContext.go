package prototypes

import (
	"../values"

	"../../context"
)

var WebGLRenderingContext *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLRenderingContextPrototype() bool {
	getEnum := func(name string) BuiltinFunction {
		return NewGetterFunction(func(stack values.Stack, this *values.Instance,
			args []values.Value, ctx context.Context) (values.Value, error) {
			return values.NewInstance(GLEnum, values.NewPropertiesWithContent(map[string]values.Value{
				".name": NewLiteralString(name, ctx),
			}, ctx), ctx), nil
		})
	}

	getGLParameter := func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {

		proto := Float32Array

		if arg, ok := values.UnpackContextValue(args[0]).(*values.Instance); ok {
			props := arg.Properties()
			if name_, ok := props.GetProperty(".name"); ok {

				if name, ok := name_.LiteralStringValue(); ok {
					switch name {
					case "MAX_FRAGMENT_UNIFORM_VECTORS":
						proto = Int
					case "MAX_VERTEX_UNIFORM_VECTORS":
						proto = Int
					default:
						proto = Float32Array
					}
				} else {
					panic("expected literal string name")
				}
			} else {
				panic("expected name (expected GLEnum)")
			}
		}

		return NewInstance(proto, ctx), nil
	}

	*WebGLRenderingContext = BuiltinPrototype{
		"WebGLRenderingContext", nil,
		map[string]BuiltinFunction{
			"ACTIVE_ATTRIBUTES":                getEnum("ACTIVE_ATTRIBUTES"),
			"ACTIVE_UNIFORMS":                  getEnum("ACTIVE_UNIFORMS"),
			"ALPHA":                            getEnum("ALPHA"),
			"ALWAYS":                           getEnum("ALWAYS"),
			"ARRAY_BUFFER":                     getEnum("ARRAY_BUFFER"),
			"ATTACHED_SHADERS":                 getEnum("ATTACHED_SHADERS"),
			"BLEND":                            getEnum("BLEND"),
			"BLEND_COLOR":                      getEnum("BLEND_COLOR"),
			"CLAMP_TO_BORDERS":                 NewGetter(Int),
			"CLAMP_TO_EDGE":                    NewGetter(Int),
			"COLOR_BUFFER_BIT":                 NewGetter(Int),
			"COMPILE_STATUS":                   getEnum("COMPILE_STATUS"),
			"CONTEXT_LOST_WEBGL":               getEnum("CONTEXT_LOST_WEBGL"),
			"CULL_FACE":                        getEnum("CULL_FACE"),
			"DELETE_STATUS":                    getEnum("DELETE_STATUS"),
			"DEPTH_BUFFER_BIT":                 NewGetter(Int),
			"DEPTH_TEST":                       getEnum("DEPTH_TEST"),
			"DITHER":                           getEnum("DITHER"),
			"DST_ALPHA":                        getEnum("DST_ALPHA"),
			"DST_COLOR":                        getEnum("DST_COLOR"),
			"DYNAMIC_DRAW":                     getEnum("DYNAMIC_DRAW"),
			"ELEMENT_ARRAY_BUFFER":             getEnum("ELEMENT_ARRAY_BUFFER"),
			"EQUAL":                            getEnum("EQUAL"),
			"FLOAT":                            getEnum("FLOAT"),
			"FRAGMENT_SHADER":                  getEnum("FRAGMENT_SHADER"),
			"GEQUAL":                           getEnum("GEQUAL"),
			"GREATER":                          getEnum("GREATER"),
			"INVALID_ENUM":                     getEnum("INVALID_ENUM"),
			"INVALID_VALUE":                    getEnum("INVALID_VALUE"),
			"INVALID_OPERATION":                getEnum("INVALID_OPERATION"),
			"INVALID_FRAMEBUFFER_OPERATION":    getEnum("INVALID_FRAMEBUFFER_OPERATION"),
			"LEQUAL":                           getEnum("LEQUAL"),
			"LESS":                             getEnum("LESS"),
			"LINEAR":                           NewGetter(Int),
			"LINES":                            getEnum("LINES"),
			"LINE_LOOP":                        getEnum("LINE_LOOP"),
			"LINE_STRIP":                       getEnum("LINE_STRIP"),
			"LINK_STATUS":                      getEnum("LINK_STATUS"),
			"LUMINANCE":                        getEnum("LUMINANCE"),
			"LUMINANCE_ALPHA":                  getEnum("LUMINANCE_ALPHA"),
			"MAX_TEXTURE_IMAGE_UNITS":          getEnum("MAX_TEXTURE_IMAGE_UNITS"),
			"MAX_COMBINED_TEXTURE_IMAGE_UNITS": getEnum("MAX_COMBINED_TEXTURE_IMAGE_UNITS"),
			"MAX_VERTEX_TEXTURE_IMAGE_UNITS":   getEnum("MAX_VERTEX_TEXTURE_IMAGE_UNITS"),
			"MAX_VERTEX_UNIFORM_VECTORS":       getEnum("MAX_VERTEX_UNIFORM_VECTORS"),
			"MAX_FRAGMENT_UNIFORM_VECTORS":     getEnum("MAX_FRAGMENT_UNIFORM_VECTORS"),
			"MIRRORED_REPEAT":                  NewGetter(Int),
			"NEAREST":                          NewGetter(Int),
			"NEVER":                            getEnum("NEVER"),
			"NO_ERROR":                         getEnum("NO_ERROR"),
			"NOTEQUAL":                         getEnum("NOTEQUAL"),
			"ONE":                              getEnum("ONE"),
			"ONE_MINUS_DST_ALPHA":              getEnum("ONE_MINUS_DST_ALPHA"),
			"ONE_MINUS_DST_COLOR":              getEnum("ONE_MINUS_DST_COLOR"),
			"ONE_MINUS_SRC_ALPHA":              getEnum("ONE_MINUS_SRC_ALPHA"),
			"ONE_MINUS_SRC_COLOR":              getEnum("ONE_MINUS_SRC_COLOR"),
			"OUT_OF_MEMORY":                    getEnum("OUT_OF_MEMORY"),
			"POINTS":                           getEnum("POINTS"),
			"POLYGON_OFFSET_FILL":              getEnum("POLYGON_OFFSET_FILL"),
			"REPEAT":                           NewGetter(Int),
			"RGB":                              getEnum("RGB"),
			"RGBA":                             getEnum("RGBA"),
			"SAMPLE_ALPHA_TO_COVERAGE":         getEnum("SAMPLE_ALPHA_TO_COVERAGE"),
			"SAMPLE_COVERAGE":                  getEnum("SAMPLE_COVERAGE"),
			"SCISSOR_TEST":                     getEnum("SCISSOR_TEST"),
			"SHADER_TYPE":                      getEnum("SHADER_TYPE"),
			"SRC_ALPHA":                        getEnum("SRC_ALPHA"),
			"SRC_COLOR":                        getEnum("SRC_COLOR"),
			"STATIC_DRAW":                      getEnum("STATIC_DRAW"),
			"STENCIL_TEST":                     getEnum("STENCIL_TEST"),
			"STREAM_DRAW":                      getEnum("STREAM_DRAW"),
			"TEXTURE0":                         NewGetter(Int),
			"TEXTURE1":                         NewGetter(Int),
			"TEXTURE2":                         NewGetter(Int),
			"TEXTURE3":                         NewGetter(Int),
			"TEXTURE4":                         NewGetter(Int),
			"TEXTURE5":                         NewGetter(Int),
			"TEXTURE6":                         NewGetter(Int),
			"TEXTURE7":                         NewGetter(Int),
			"TEXTURE8":                         NewGetter(Int),
			"TEXTURE9":                         NewGetter(Int),
			"TEXTURE10":                        NewGetter(Int),
			"TEXTURE11":                        NewGetter(Int),
			"TEXTURE12":                        NewGetter(Int),
			"TEXTURE13":                        NewGetter(Int),
			"TEXTURE14":                        NewGetter(Int),
			"TEXTURE15":                        NewGetter(Int),
			"TEXTURE_2D":                       getEnum("TEXTURE_2D"),
			"TEXTURE_MAG_FILTER":               getEnum("TEXTURE_MAG_FILTER"),
			"TEXTURE_MIN_FILTER":               getEnum("TEXTURE_MIN_FILTER"),
			"TEXTURE_WRAP_S":                   getEnum("TEXTURE_WRAP_S"),
			"TEXTURE_WRAP_T":                   getEnum("TEXTURE_WRAP_T"),
			"TRIANGLES":                        getEnum("TRIANGLES"),
			"TRIANGLE_FAN":                     getEnum("TRIANGLE_FAN"),
			"TRIANGLE_STRIP":                   getEnum("TRIANGLE_STRIP"),
			"UNSIGNED_BYTE":                    getEnum("UNSIGNED_BYTE"),
			"UNSIGNED_INT":                     getEnum("UNSIGNED_INT"),
			"UNSIGNED_SHORT":                   getEnum("UNSIGNED_SHORT"),
			"VALIDATE_STATUS":                  getEnum("VALIDATE_STATUS"),
			"VERTEX_SHADER":                    getEnum("VERTEX_SHADER"),
			"ZERO":                             getEnum("ZERO"),

			"activeTexture":           NewNormal(Int, nil),
			"attachShader":            NewNormal(&And{WebGLProgram, WebGLShader}, nil),
			"bindBuffer":              NewNormal(&And{GLEnum, WebGLBuffer}, nil),
			"bindTexture":             NewNormal(&And{GLEnum, WebGLTexture}, nil),
			"blendFunc":               NewNormal(&And{GLEnum, GLEnum}, nil),
			"blendFuncSeparate":       NewNormal(&Many{4, GLEnum}, nil),
			"bufferData":              NewNormal(&And{GLEnum, &And{TypedArray, GLEnum}}, nil),
			"clear":                   NewNormal(Int, nil),
			"clearColor":              NewNormal(&Many{4, Number}, nil),
			"compileShader":           NewNormal(WebGLShader, nil),
			"createBuffer":            NewNormal(&None{}, WebGLBuffer),
			"createProgram":           NewNormal(&None{}, WebGLProgram),
			"createShader":            NewNormal(GLEnum, WebGLShader),
			"createTexture":           NewNormal(&None{}, WebGLTexture),
			"depthFunc":               NewNormal(GLEnum, nil),
			"disable":                 NewNormal(GLEnum, nil),
			"drawArrays":              NewNormal(&And{GLEnum, &And{Int, Int}}, nil),
			"drawElements":            NewNormal(&And{GLEnum, &And{Int, &And{GLEnum, Int}}}, nil),
			"enable":                  NewNormal(GLEnum, nil),
			"enableVertexAttribArray": NewNormal(Int, nil),
			"getAttribLocation":       NewNormal(&And{WebGLProgram, String}, Int),
			"getError":                NewNormal(&None{}, GLEnum),
			"getExtension":            NewNormal(String, WebGLExtension),
			"getParameter":            NewNormalFunction(GLEnum, getGLParameter),
			"getProgramInfoLog":       NewNormal(WebGLProgram, String),
			"getProgramParameter":     NewNormal(&And{WebGLProgram, GLEnum}, Boolean),
			"getShaderInfoLog":        NewNormal(WebGLShader, String),
			"getShaderParameter":      NewNormal(&And{WebGLShader, GLEnum}, Boolean),
			"getUniformLocation":      NewNormal(&And{WebGLProgram, String}, Int),
			"linkProgram":             NewNormal(WebGLProgram, nil),
			"scissor":                 NewNormal(&Many{4, Number}, nil),
			"shaderSource":            NewNormal(&And{WebGLShader, String}, nil),
			"texParameterf":           NewNormal(&And{GLEnum, &And{GLEnum, Number}}, nil),
			"texParameteri":           NewNormal(&And{GLEnum, &And{GLEnum, Int}}, nil),

			"texImage2D": NewNormal(&And{&And{GLEnum, &And{Int, GLEnum}}, &Or{
				&And{GLEnum, &And{GLEnum, &Or{ImageData, &Or{ArrayBuffer, HTMLCanvasElement}}}},
				&And{&Many{3, Int}, &And{GLEnum, &And{GLEnum, TypedArray}}},
			}}, nil),

			"uniform1f":       NewNormal(&And{Int, Number}, nil),
			"uniform1fv":      NewNormal(&And{Int, Array}, nil),
			"uniform1i":       NewNormal(&And{Int, Int}, nil),
			"uniform1iv":      NewNormal(&And{Int, Array}, nil),
			"uniform2f":       NewNormal(&And{Int, &And{Number, Number}}, nil),
			"uniform2fv":      NewNormal(&And{Int, Array}, nil),
			"uniform2i":       NewNormal(&And{Int, &And{Int, Int}}, nil),
			"uniform2iv":      NewNormal(&And{Int, Array}, nil),
			"uniform3f":       NewNormal(&And{Int, &Many{3, Number}}, nil),
			"uniform3fv":      NewNormal(&And{Int, Array}, nil),
			"uniform3i":       NewNormal(&And{Int, &Many{3, Int}}, nil),
			"uniform3iv":      NewNormal(&And{Int, Array}, nil),
			"uniform4f":       NewNormal(&And{Int, &Many{4, Number}}, nil),
			"uniform4fv":      NewNormal(&And{Int, Array}, nil),
			"uniform4i":       NewNormal(&And{Int, &Many{4, Int}}, nil),
			"uniform4iv":      NewNormal(&And{Int, Array}, nil),
			"useProgram":      NewNormal(WebGLProgram, nil),
			"validateProgram": NewNormal(WebGLProgram, nil),
			"vertexAttribPointer": NewMethodLikeNormal(&And{Int,
				&And{Int, &And{GLEnum, &And{Boolean, &And{Int, Int}}}}}, Boolean),
			"viewport": NewNormal(&Many{4, Number}, nil),
		},
		NewNoContentGenerator(WebGLRenderingContext),
	}

	return true
}

var _WebGLRenderingContextOk = generateWebGLRenderingContextPrototype()
