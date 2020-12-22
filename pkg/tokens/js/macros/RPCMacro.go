package macros

import (
  "../../context"
  "../../js"
)

type RPCMacro struct {
  ToInstance
}

func newRPCMacro(args []js.Expression, ctx context.Context) (RPCMacro, error) {
  interfExpr, err := getTypeExpression(args[0])
  if err != nil {
    return RPCMacro{}, err
  }

  return RPCMacro{newToInstance(args[1:], interfExpr, ctx)}, nil
}

func (m *RPCMacro) ResolveExpressionActivity(usage js.Usage) error {
  if err := m.Macro.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  if !m.interf.IsRPC() {
    panic("should've been checked during eval phase")
  }

  return nil
}
