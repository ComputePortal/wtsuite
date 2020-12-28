package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildFunctionArgumentRole(ts  []raw.Token) (glsl.FunctionArgumentRole, []raw.Token, error) {
  role := glsl.NO_ROLE

  iRemaining := 0
  for i := 0; i < len(ts); i++ {
    t := ts[i]
    if raw.IsAnyWord(t) {
      w, err := raw.AssertWord(t)
      if err != nil {
        panic(err)
      }

      switch w.Value() {
      case "in":
        role = role & glsl.IN_ROLE
        iRemaining = i+1
        continue
      case "out":
        role = role & glsl.OUT_ROLE
        iRemaining = i+1
        continue
      default:
        break
      }
    }

    break
  }

  if iRemaining < len(ts) {
    ts = ts[iRemaining:]
  } else {
    ts = []raw.Token{}
  }

  return role, ts, nil
}

func (p *GLSLParser) buildFunctionArgument(ts []raw.Token) (*glsl.FunctionArgument, error) {
  var err error
  var role glsl.FunctionArgumentRole
  role, ts, err = p.buildFunctionArgumentRole(ts)
  if err != nil {
    return nil, err
  }

  typeExpr, err := p.buildTypeExpression(ts[0:1])
  if err != nil {
    return nil, err
  }

  nameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, err
  }

  n := -1
  if len(ts) == 2 {
    n, err = p.buildArraySize(ts[2])
    if err != nil {
      return nil, err
    }
  } else if len(ts) > 2 {
    errCtx := ts[3].Context()
    return nil, errCtx.NewError("Error: unexpected tokens")
  }

  return glsl.NewFunctionArgument(role, typeExpr, nameToken.Value(), n, nameToken.Context()), nil
}

func (p *GLSLParser) buildFunctionInterface(ts []raw.Token) (*glsl.FunctionInterface, error) {
  if len(ts) != 3 {
    panic("internal error")
  }

  retTypeExpr, err := p.buildTypeExpression(ts[0:1])
  if err != nil {
    return nil, err
  }

  nameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    panic(err)
  }

  argParens, err := raw.AssertParensGroup(ts[2])
  if err != nil {
    panic(err)
  }

  if argParens.IsSemiColon() {
    errCtx := ts[2].Context()
    return nil, errCtx.NewError("Error: expected comma separators")
  }

  fArgs := []*glsl.FunctionArgument{}

  for _, field := range argParens.Fields {
    fArg, err := p.buildFunctionArgument(field)
    if err != nil {
      return nil, err
    }

    fArgs = append(fArgs, fArg)
  }

  return glsl.NewFunctionInterface(retTypeExpr, nameToken.Value(), fArgs, nameToken.Context()), nil
}

func (p *GLSLParser) buildFunction(ts []raw.Token) ([]raw.Token, error) {
  if len(ts) < 4 || !raw.IsAnyWord(ts[1]) || !raw.IsParensGroup(ts[2]) || !raw.IsBracesGroup(ts[3]) {
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: bad function")
  }

  functionInterf, err := p.buildFunctionInterface(ts[0:3])
  if err != nil {
    return nil, err
  }

  remainingTokens := stripSeparators(0, ts[4:], patterns.SEMICOLON)

  // build the statements
  contentBrace, err := raw.AssertBracesGroup(ts[3])
  if err != nil {
    panic(err)
  }

  if contentBrace.IsComma() {
    errCtx := contentBrace.Context()
    return nil, errCtx.NewError("Error: expected semicolon separators")
  }

  statements, err := p.buildBlockStatements(contentBrace)
  if err != nil {
    return nil, err
  }
  
  fn := glsl.NewFunction(functionInterf, statements, raw.MergeContexts(ts...))

  p.module.AddStatement(fn)

  return remainingTokens, nil
}
