package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type MacroStatement func(args []Expression, ctx context.Context) (Statement, error) 

type MacroExpression func(args []Expression, ctx context.Context) (Expression, error) 

var _macroStatements = map[string]MacroStatement{
  "setIndex": NewSetDynamicIndex,
}

var _macroExpressions = map[string]MacroExpression{
  "getIndex": NewGetDynamicIndex,
}

func IsMacroStatement(call *Call) bool {
  key := call.Name()

  if key != "" {
    _, ok := _macroStatements[key]
    return ok
  } else {
    return false
  }
}

func IsMacroExpression(call *Call) bool {
  key := call.Name()

  if key != "" {
    _, ok := _macroExpressions[key]
    return ok
  } else {
    return false
  }
}

func DispatchMacroStatement(call *Call) (Statement, error) {
  key := call.Name()

  macro, ok := _macroStatements[key]
  ctx := call.Context()
  if ok {
    return macro(call.Args(), ctx)
  } else {
    return nil, ctx.NewError("Error: not a macro")
  }
}

func DispatchMacroExpression(call *Call) (Expression, error) {
  key := call.Name()

  macro, ok := _macroExpressions[key]
  ctx := call.Context()
  if ok {
    return macro(call.Args(), ctx)
  } else {
    return nil, ctx.NewError("Error: not a macro")
  }
}
