package macros

import (
  "strings"

  "../prototypes"
  "../values"

  "../../context"
  "../../js"
)

type RPCServer struct {
  RPCMacro
}

func NewRPCServer(args []js.Expression, ctx context.Context) (js.Expression, error) {
  if len(args) != 2 {

    return nil, ctx.NewError("Error: expected 2 arguments")
  }

  rpcMacro, err := newRPCMacro(args, ctx)
  if err != nil {
    return nil, err
  }

  return &RPCServer{rpcMacro}, nil
}

func (m *RPCServer) Dump(indent string) string {
  return indent + "RPCServer(...)\n"
}

func (m *RPCServer) WriteExpression() string {
  var b strings.Builder

  b.WriteString("new ")
  b.WriteString(rpcServerHeader.Name())
  b.WriteString("(")
  b.WriteString(m.interf.Name())
  b.WriteString(",")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(")")

  return b.String()
}

func (m *RPCServer) EvalExpression() (values.Value, error) {
  ctx := m.Context()

  args, err := m.evalArgs()
  if err != nil {
    return nil, err
  }

  if !m.interf.IsRPC() {
    errCtx := m.interfExpr.Context()
    return nil, errCtx.NewError("Error: " + m.interf.Name() + " is not an rpc interface")
  }

  checkVal := values.NewInstance(m.interf, m.Context())
  if err := checkVal.Check(args[0], ctx); err != nil {
    return nil, err
  }

  return prototypes.NewRPCServer(m.Context()), nil
}

func (m *RPCServer) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(rpcServerHeader, m.Context())

  return m.RPCMacro.ResolveExpressionActivity(usage)
}

func (m *RPCServer) UniqueExpressionNames(ns js.Namespace) error {
  if err := UniqueHeaderNames(rpcServerHeader, ns); err != nil {
    return err
  }

  return m.RPCMacro.UniqueExpressionNames(ns)
}
