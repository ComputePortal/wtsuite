package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildLiteralBoolExpression(t raw.Token) (*glsl.LiteralBool, error) {
  b, err := raw.AssertLiteralBool(t)
  if err != nil {
    panic(err)
  }

  return glsl.NewLiteralBool(b.Value(), b.Context()), nil
}

func (p *GLSLParser) buildLiteralIntExpression(t raw.Token) (*glsl.LiteralInt, error) {
  i, err := raw.AssertLiteralInt(t)
  if err != nil {
    panic(err)
  }

  return glsl.NewLiteralInt(i.Value(), i.Context()), nil
}

func (p *GLSLParser) buildLiteralFloatExpression(t raw.Token) (*glsl.LiteralFloat, error) {
  f, err := raw.AssertLiteralFloat(t, "")
  if err != nil {
    panic(err)
  }

  if f.Unit() != "" {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: united floats not supported")
  }

  return glsl.NewLiteralFloat(f.Value(), f.Context()), nil
}

func (p *GLSLParser) buildParensExpression(t raw.Token) (*glsl.Parens, error) {
  group, err := raw.AssertParensGroup(t)
  if err != nil {
    panic(err)
  }

  if !group.IsSingle() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: simple parentheses must have exactly one entry")
  }

  expr, err := p.buildExpression(group.Fields[0])
  if err != nil {
    return nil, err
  }

  return glsl.NewParens(expr, group.Context()), nil
}

func (p *GLSLParser) buildUnaryOpExpression(t raw.Token) (glsl.Expression, error) {
  op, err := raw.AssertAnyUnaryOperator(t)
  if err != nil {
    panic(err)
  }

  a, err := p.buildExpression(op.Args()[0:1])
  if err != nil {
    return nil, err
  }

  ctx := op.Context()

  switch op.Name() {
  case "pre!":
    return glsl.NewNotOp(a, ctx), nil
  case "pre-":
    return glsl.NewNegOp(a, ctx), nil
  case "pre+":
    return glsl.NewPosOp(a, ctx), nil
  default:
    return nil, ctx.NewError("Error: unhandled preunary op")
  }
}

func (p *GLSLParser) buildBinaryOpExpression(t raw.Token) (glsl.Expression, error) {
  op, err := raw.AssertAnyBinaryOperator(t)
  if err != nil {
    panic(err)
  }

  a, err := p.buildExpression(op.Args()[0:1])
  if err != nil {
    return nil, err
  }
  
  b, err := p.buildExpression(op.Args()[1:2])
  if err != nil {
    return nil, err
  }

  ctx := op.Context()

  switch op.Name() {
  case "bin+":
    return glsl.NewAddOp(a, b, ctx), nil
  case "bin-":
    return glsl.NewSubOp(a, b, ctx), nil
  case "bin/":
    return glsl.NewDivOp(a, b, ctx), nil
  case "bin*":
    return glsl.NewMulOp(a, b, ctx), nil
  case "bin%":
    mod := glsl.NewVarExpression("mod", ctx)
    return glsl.NewCall(mod, []glsl.Expression{a, b}, ctx), nil
  case "bin&&":
    return glsl.NewAndOp(a, b, ctx), nil
  case "bin||":
    return glsl.NewOrOp(a, b, ctx), nil
  case "bin^^":
    return glsl.NewXorOp(a, b, ctx), nil
  case "bin<":
    return glsl.NewLTOp(a, b, ctx), nil
  case "bin>":
    return glsl.NewGTOp(a, b, ctx), nil
  case "bin<=":
    return glsl.NewLEOp(a, b, ctx), nil
  case "bin>=":
    return glsl.NewGEOp(a, b, ctx), nil
  case "bin==":
    return glsl.NewEqOp(a, b, ctx), nil
  case "bin!=":
    return glsl.NewNEOp(a, b, ctx), nil
  default:
    return nil, ctx.NewError("Error: unhandled binary op")
  }
}

func (p *GLSLParser) buildVarExpression(t raw.Token) (*glsl.VarExpression, error) {
  w, err := raw.AssertWord(t)
  if err != nil {
    return nil, err
  }

  return glsl.NewVarExpression(w.Value(), w.Context()), nil
}

func (p *GLSLParser) buildIndexExpression(ts []raw.Token) (*glsl.Index, error) {
	n := len(ts)

	lhs, err := p.buildExpression(ts[0 : n-1])
	if err != nil {
		return nil, err
	}

	group, err := raw.AssertBracketsGroup(ts[n-1])
	if err != nil {
		return nil, err
	}

	if group.IsEmpty() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: index can't be empty")
	}

	if !group.IsSingle() { // comma's should've been combined in operators
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: multi indexing not allowed")
	}

	field := group.Fields[0]

	index, err := p.buildExpression(field)
	if err != nil {
		return nil, err
	}

	return glsl.NewIndex(lhs, index, group.Context()), nil
}

func (p *GLSLParser) buildCall(ts []raw.Token) (*glsl.Call, error) {
  parens, err := raw.AssertGroup(ts[len(ts)-1])
  if err != nil {
    return nil, err
  }

  if parens.IsSemiColon() {
    errCtx := parens.Context()
    return nil, errCtx.NewError("Error: expected comma separators")
  }

  lhs, err := p.buildExpression(ts[0:len(ts)-1])
  if err != nil {
    return nil, err
  }

  args := []glsl.Expression{}

  for _, field := range parens.Fields {
    arg, err := p.buildExpression(field)
    if err != nil {
      return nil, err
    }

    args = append(args, arg)
  }

  ctx := raw.MergeContexts(ts...)
  return glsl.NewCall(lhs, args, ctx), nil
}

// can also return a macro
func (p *GLSLParser) buildCallExpression(ts []raw.Token) (glsl.Expression, error) {
  call, err := p.buildCall(ts)
  if err != nil {
    return nil, err
  }

  if glsl.IsMacroExpression(call) {
    return glsl.DispatchMacroExpression(call)
  } else if glsl.IsMacroStatement(call) {
    errCtx := call.Context()
    return nil, errCtx.NewError("Error: is a statement, not an expression")
  }

  return call, nil
}

func (p *GLSLParser) buildMemberExpression(ts []raw.Token) (*glsl.Member, error) {
	n := len(ts)

	if n < 3 {
		errCtx := ts[n-2].Context()
		return nil, errCtx.NewError("Error: member of nothing")
	}

	lhs, err := p.buildExpression(ts[0 : n-2])
	if err != nil {
		return nil, err
	}

	w, err := raw.AssertWord(ts[n-1])
	if err != nil {
		panic(err)
	}

	return glsl.NewMember(lhs, glsl.NewWord(w.Value(), w.Context()), ts[n-2].Context()), nil
}

func (p *GLSLParser) buildExpression(ts []raw.Token) (glsl.Expression, error) {
  ts, err := p.nestOperators(ts)
  if err != nil {
    return nil, err
  }

  ts = p.expandTmpGroups(ts)

  n := len(ts)
  switch {
  case n == 1:
    switch {
    case raw.IsLiteralString(ts[0]):
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: literal strings not supported (though maybe for certain macros in future)")
    case raw.IsLiteralBool(ts[0]):
      return p.buildLiteralBoolExpression(ts[0])
    case raw.IsLiteralInt(ts[0]):
      return p.buildLiteralIntExpression(ts[0])
    case raw.IsLiteralFloat(ts[0]):
      return p.buildLiteralFloatExpression(ts[0])
    case raw.IsParensGroup(ts[0]):
      return p.buildParensExpression(ts[0])
    case raw.IsAnyUnaryOperator(ts[0]):
      return p.buildUnaryOpExpression(ts[0])
    case raw.IsAnyBinaryOperator(ts[0]):
      return p.buildBinaryOpExpression(ts[0])
    case raw.IsAnyWord(ts[0]):
      return p.buildVarExpression(ts[0])
    default:
      // nested group
      if raw.IsTmpGroup(ts[0]) {
        gr, err := raw.AssertGroup(ts[0])
        if err != nil {
          panic(err)
        }

        return p.buildExpression(gr.Fields[0])
      } else {
				errCtx := ts[0].Context()
				err := errCtx.NewError("Error: expression not yet supported")
				return nil, err
      }
    }
  case raw.IsParensGroup(ts[n-1]):
    return p.buildCallExpression(ts)
  case raw.IsBracketsGroup(ts[n-1]):
    return p.buildIndexExpression(ts)
  case n > 2 && raw.IsSymbol(ts[n-2], patterns.PERIOD) && raw.IsAnyWord(ts[n-1]):
    return p.buildMemberExpression(ts)
  default:
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: unrecognized expression (hint: missing semicolon?")
  }
}

