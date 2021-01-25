package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildTypeExpression(ts []raw.Token) (*glsl.TypeExpression, error) {
  w, rem, err := condensePackagePeriods(ts)
  if err != nil {
    return nil, err
  }

  if len(rem) != 0 {
    errCtx := raw.MergeContexts(rem...)
    return nil, errCtx.NewError("Error: bad type expression")
  }

  return glsl.NewTypeExpression(w.Value(), w.Context()), nil
}

func (p *GLSLParser) buildAttribute(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  n := len(ts)
  if n < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected at least 3 tokens")
  }

  typeExpr, err := p.buildTypeExpression(ts[1:n-1])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[n-1])
  if err != nil {
    return nil, err
  }

  st := glsl.NewAttribute(typeExpr, name.Value(), name.Context())

  p.module.AddStatement(st)

  if isExport {
    if err := p.buildSimpleExport(name.Value(), name.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}

func (p *GLSLParser) buildPrecisionType(t raw.Token) (glsl.PrecisionType, error) {
  w, err := raw.AssertWord(t)
  if err != nil {
    return glsl.DEFAULTP, err
  }

  switch w.Value() {
  case "lowp":
    return glsl.LOWP, nil
  case "mediump":
    return glsl.MEDIUMP, nil
  case "highp":
    return glsl.HIGHP, nil
  default:
    errCtx := t.Context()
    return glsl.DEFAULTP, errCtx.NewError("Error: unrecognized precision type")
  }
}

func (p *GLSLParser) buildVarying(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  
  n := len(ts)
  if n < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected at least 3")
  }
  
  var precType glsl.PrecisionType = glsl.DEFAULTP
  if raw.IsWord(ts[1], "highp") || raw.IsWord(ts[1], "mediump") || raw.IsWord(ts[1], "lowp") {
    var err error
    precType, err = p.buildPrecisionType(ts[1])
    if err != nil {
      return nil, err
    }

    n = n - 1
    ts = ts[1:]
  }

  typeExpr, err := p.buildTypeExpression(ts[1:n-1])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[n-1])
  if err != nil {
    return nil, err
  }

  st := glsl.NewVarying(precType, typeExpr, name.Value(), name.Context())

  p.module.AddStatement(st)

  if isExport {
    if err := p.buildSimpleExport(name.Value(), name.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}

func (p *GLSLParser) buildLiteralIndex(t raw.Token) (int, error) {
  if brackets, err := raw.AssertBracketsGroup(t); err != nil {
    return 0, err
  } else {
    if !brackets.IsSingle() {
      errCtx := t.Context()
      return 0, errCtx.NewError("Error: expected single argument")
    }

    content_ := brackets.Fields[0]
    if len(content_) != 1 {
      errCtx := t.Context()
      return 0, errCtx.NewError("Error: expected single argument")
    }

    content := content_[0]

    litContent, err := raw.AssertLiteralInt(content)
    if err != nil {
      return 0, err
    }

    return litContent.Value(), nil
  }
}

func (p *GLSLParser) buildArraySize(t raw.Token) (int, error) {
  n, err := p.buildLiteralIndex(t)
  if err != nil {
    return 0, err
  }

  if n <= 0 {
    errCtx := t.Context()
    return 0, errCtx.NewError("Error: invalid literal array size")
  }

  return n, nil
}

func (p *GLSLParser) buildUniform(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  n := len(ts)
  if n < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected at least 3 tokens")
  }

  arraySize := -1
  if raw.IsBracketsGroup(ts[n-1]) {
    var err error
    arraySize, err = p.buildArraySize(ts[n-1])
    if err != nil {
      return nil, err
    }
    n = n - 1
  }

  typeExpr, err := p.buildTypeExpression(ts[1:n-1])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[n-1])
  if err != nil {
    return nil, err
  }

  st := glsl.NewUniform(typeExpr, name.Value(), arraySize, name.Context())
  
  p.module.AddStatement(st)

  if isExport {
    if err := p.buildSimpleExport(name.Value(), name.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}

func (p *GLSLParser) buildConst(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  n := len(ts)
  if n < 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected at least 4 tokens")
  }

  iequal := nextSeparatorPosition(ts, patterns.EQUAL)

  if iequal == n {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: bad const statement")
  }

  iname := iequal - 1
  arraySize := -1
  if raw.IsBracketsGroup(ts[iname]) {
    var err error
    arraySize, err = p.buildArraySize(ts[iname])
    if err != nil {
      return nil, err
    }
    iname = iname - 1
  }

  typeExpr, err := p.buildTypeExpression(ts[1:iname])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[iname])
  if err != nil {
    return nil, err
  }

  rhsExpr, err := p.buildExpression(ts[iequal+1:])
  if err != nil {
    return nil, err
  }

  st := glsl.NewConst(typeExpr, name.Value(), arraySize, rhsExpr, isExport, name.Context())
  
  p.module.AddStatement(st)

  if isExport {
    if err := p.buildSimpleExport(name.Value(), name.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}

func (p *GLSLParser) buildPrecision(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  
  if len(ts) != 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 tokens")
  }

  precType, err := p.buildPrecisionType(ts[1])
  if err != nil {
    return nil, err
  }

  typeExpr, err := p.buildTypeExpression(ts[2:])
  if err != nil {
    return  nil, err
  }
  
  st := glsl.NewPrecision(precType, typeExpr, raw.MergeContexts(ts...))

  p.module.AddStatement(st)

  return remainingTokens, nil
}
