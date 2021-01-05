package macros

import (
  "strconv"
  "strings"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
  "github.com/computeportal/wtsuite/pkg/tokens/js"
  "github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type WebGLProgram struct {
  vertexSource string
  fragmentSource string
  Macro
}

type TranspileWebGLShadersFunc func(callerPath string, vertexPath *js.Word, vertexConsts map[string]values.Value,
  fragmentPath *js.Word, fragmentConsts map[string]values.Value) (string, string, error)

var transpileWebGLShaders TranspileWebGLShadersFunc = nil
  
func RegisterTranspileWebGLShaders(fn TranspileWebGLShadersFunc) bool {
  transpileWebGLShaders = fn

  return true
}

func NewWebGLProgram(args []js.Expression, ctx context.Context) (js.Expression, error) {
  if len(args) < 3 || len(args) > 5 {
    errCtx := ctx
    return nil, errCtx.NewError("Error: expected 3, 4 or 5 arguments, got " + strconv.Itoa(len(args)))
  }

  return &WebGLProgram{"", "", newMacro(args, ctx)}, nil
}

func (m *WebGLProgram) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("new WebGLProgram(...)")
  for _, arg := range m.args {
    b.WriteString(arg.Dump(indent + "  "))
  }

  return b.String()
}

func (m *WebGLProgram) WriteExpression() string {
  var b strings.Builder

  b.WriteString("((function(gl,v,f){")
  b.WriteString("return ")
  b.WriteString(webGLProgramHeader.Name())
  b.WriteString("(gl,")
  b.WriteString(m.vertexSource)
  b.WriteString(",")
  b.WriteString(m.fragmentSource)
  b.WriteString(")})(")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(",")
  if (len(m.args) > 3) {
    b.WriteString(m.args[2].WriteExpression())
  } else {
    b.WriteString("{}")
  }
  b.WriteString(",")
  if (len(m.args) > 4) {
    b.WriteString(m.args[4].WriteExpression())
  } else {
    b.WriteString("{}")
  }
  b.WriteString("))")

  return b.String()
}

func (m *WebGLProgram) ResolveExpressionNames(scope js.Scope) error {
  return m.Macro.ResolveExpressionNames(scope)
}

func (m *WebGLProgram) EvalExpression() (values.Value, error) {
  args, err := m.evalArgs()
  if err != nil {
    return nil, err
  }

  if !prototypes.IsWebGLRenderingContext(args[0]) {
    errCtx := m.args[0].Context()
    return nil, errCtx.NewError("Error: expected WebGLRenderingContext, got " + args[0].TypeName())
  }

  vertexPath_, ok := args[1].LiteralStringValue()
  if !ok {
    errCtx := m.args[1].Context()
    return nil, errCtx.NewError("Error: expected literal string, got " + args[1].TypeName())
  }
  vertexPath := js.NewWord(vertexPath_, args[1].Context())

  vertexConsts := make(map[string]values.Value)

  ifrag := 2
  if len(args) > 3 {
    ifrag = 3
    if vertexConsts, err = prototypes.GetLiteralObjectMembers(args[2]); err != nil {
      return nil, err
    }
  }

  fragmentPath_, ok := args[ifrag].LiteralStringValue()
  if !ok {
    errCtx := m.args[ifrag].Context()
    return nil, errCtx.NewError("Error: expected literal string, got " + args[ifrag].TypeName())
  }
  fragmentPath := js.NewWord(fragmentPath_, args[ifrag].Context())

  fragmentConsts := make(map[string]values.Value)

  if len(args) > 4 {
    if fragmentConsts, err = prototypes.GetLiteralObjectMembers(args[4]); err != nil {
      return nil, err
    }
  }

  if transpileWebGLShaders == nil {
    panic("transpileWebGLShaders not registered")
  }

  ctx := m.Context()
  callerPath := ctx.Path()
  m.vertexSource, m.fragmentSource, err = transpileWebGLShaders(
    callerPath, vertexPath, vertexConsts, fragmentPath, fragmentConsts)
  if err != nil {
    return nil, err
  }

  return prototypes.NewWebGLProgram(m.Context()), nil
}

func (m *WebGLProgram) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(webGLProgramHeader, m.Context())

  return m.Macro.ResolveExpressionActivity(usage)
}

func (m *WebGLProgram) UniqueExpressionNames(ns js.Namespace) error {
  if err := UniqueHeaderNames(webGLProgramHeader, ns); err != nil {
    return err
  }

  return m.Macro.UniqueExpressionNames(ns)
}
