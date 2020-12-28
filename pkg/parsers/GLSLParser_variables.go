package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildTypeExpression(ts []raw.Token) (*glsl.TypeExpression, error) {
  if len(ts) != 1 {
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: package type expression not yet supported")
  }

  w, err := raw.AssertWord(ts[0])
  if err != nil {
    return nil, err
  }

  return glsl.NewTypeExpression(w.Value(), w.Context()), nil
}

func (p *GLSLParser) buildAttribute(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  if len(ts) != 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 tokens")
  }

  typeExpr, err := p.buildTypeExpression(ts[1:2])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[2])
  if err != nil {
    return nil, err
  }

  st := glsl.NewAttribute(typeExpr, name.Value(), name.Context())

  p.module.AddStatement(st)

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

func (p *GLSLParser) buildVarying(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  
  if len(ts) != 3 && len(ts) != 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 or 4 tokens")
  }
  
  var precType glsl.PrecisionType = glsl.DEFAULTP
  if len(ts) == 4 {
    var err error
    precType, err = p.buildPrecisionType(ts[1])
    if err != nil {
      return nil, err
    }

    ts = ts[1:]
  }

  typeExpr, err := p.buildTypeExpression(ts[1:2])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[2])
  if err != nil {
    return nil, err
  }

  st := glsl.NewVarying(precType, typeExpr, name.Value(), name.Context())

  p.module.AddStatement(st)

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

func (p *GLSLParser) buildUniform(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  if len(ts) != 3 && len(ts) != 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 or 4 tokens")
  }

  typeExpr, err := p.buildTypeExpression(ts[1:2])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[2])
  if err != nil {
    return nil, err
  }

  n := -1

  if len(ts) == 4 {
    n, err = p.buildArraySize(ts[3])
    if err != nil {
      return nil, err
    }
  }

  st := glsl.NewUniform(typeExpr, name.Value(), n, name.Context())
  
  p.module.AddStatement(st)

  return remainingTokens, nil
}

func (p *GLSLParser) buildConst(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  if len(ts) != 3 && len(ts) != 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 or 4 tokens")
  }

  typeExpr, err := p.buildTypeExpression(ts[1:2])
  if err != nil {
    return nil, err
  }

  name, err := raw.AssertWord(ts[2])
  if err != nil {
    return nil, err
  }

  n := -1

  if len(ts) == 4 {
    n, err = p.buildArraySize(ts[3])
    if err != nil {
      return nil, err
    }
  }

  st := glsl.NewConst(typeExpr, name.Value(), n, name.Context())
  
  p.module.AddStatement(st)

  return remainingTokens, nil
}
