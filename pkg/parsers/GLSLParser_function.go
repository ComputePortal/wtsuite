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
  n := len(ts)
  if n < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: bad function definition")
  }

  var retTypeExpr *glsl.TypeExpression = nil
  if !raw.IsWord(ts[0], "void") {
    var err error
    retTypeExpr, err = p.buildTypeExpression(ts[0:n-2])
    if err != nil {
      return nil, err
    }
  }

  nameToken, err := raw.AssertWord(ts[n-2])
  if err != nil {
    return nil, err
  }

  argParens, err := raw.AssertParensGroup(ts[n-1])
  if err != nil {
    return nil, err
  }

  if argParens.IsSemiColon() {
    errCtx := ts[n-1].Context()
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

func (p *GLSLParser) buildFunction(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  iparens := raw.FindFirstParensGroup(ts, 0)
  if iparens == -1 {
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: bad function definition")
  }

  ibrace := iparens + 1

  functionInterf, err := p.buildFunctionInterface(ts[0:ibrace])
  if err != nil {
    return nil, err
  }

  remainingTokens := stripSeparators(0, ts[ibrace+1:], patterns.SEMICOLON)

  // build the statements
  contentBrace, err := raw.AssertBracesGroup(ts[ibrace])
  if err != nil {
    return nil, err
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

  if isExport {
    if err := p.buildSimpleExport(functionInterf.Name(), functionInterf.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}
