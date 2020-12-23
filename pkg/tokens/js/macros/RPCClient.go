package macros

import (
  "strings"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
  "github.com/computeportal/wtsuite/pkg/tokens/js"
  "github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
  "github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type RPCClient struct {
  RPCMacro
}

func NewRPCClient(args []js.Expression, ctx context.Context) (js.Expression, error) {
  if len(args) != 2 {
    return nil, ctx.NewError("Error: expected 2 arguments")
  }

  rpcMacro, err := newRPCMacro(args, ctx)
  if err != nil {
    return nil, err
  }

  return &RPCClient{rpcMacro}, nil
}

func (m *RPCClient) Dump(indent string) string {
  return indent + "RPCClient(...)\n"
}

func (m *RPCClient) WriteExpression() string {
  var b strings.Builder

  b.WriteString("new ")
  b.WriteString(rpcClientHeader.Name())
  b.WriteString("(")
  b.WriteString(m.interf.Name())
  b.WriteString(",")
  b.WriteString(m.args[0].WriteExpression())
  b.WriteString(")")

  return b.String()
}

func (m *RPCClient) EvalExpression() (values.Value, error) {
  ctx := m.Context()

  args, err := m.evalArgs()
  if err != nil {
    return nil, err
  }

  if !m.interf.IsRPC() {
    errCtx := m.interfExpr.Context()
    return nil, errCtx.NewError("Error: " + m.interf.Name() + " is not an rpc interface")
  }

  s := prototypes.NewString(ctx)
  checkVal := values.NewFunction([]values.Value{s, prototypes.NewPromise(s, ctx)}, ctx)
  if err := checkVal.Check(args[0], ctx); err != nil {
    return nil, err
  }

  return values.NewInstance(m.interf, m.Context()), nil
}

func (m *RPCClient) ResolveExpressionActivity(usage js.Usage) error {
  ResolveHeaderActivity(rpcClientHeader, m.Context())

  return m.RPCMacro.ResolveExpressionActivity(usage)
}

func (m *RPCClient) UniqueExpressionNames(ns js.Namespace) error {
  if err := UniqueHeaderNames(rpcClientHeader, ns); err != nil {
    return err
  }

  return m.RPCMacro.UniqueExpressionNames(ns)
}
