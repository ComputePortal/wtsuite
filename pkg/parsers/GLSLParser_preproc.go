package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

// assumed to be called from top-level
func (p *GLSLParser) buildDirective(ts []raw.Token) ([]raw.Token, error) {
  if len(ts) < 2 {
    panic("unexpected")
  }

  pragmaWord, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, err
  }

  hashContext := ts[0].Context()
  pragmaContext := ts[1].Context()

  // the hash symbol must come at the beginning of a line
  if !hashContext.IsAtLineStart() {
    errCtx := context.MergeContexts(hashContext, pragmaContext)
    return nil, errCtx.NewError("Error: not at start of line")
  }

  // there can't be a space between the hash symbol and the directive name
  if hashContext.Distance(&pragmaContext) != 0 {
    errCtx := context.MergeContexts(hashContext, pragmaContext)
    return nil, errCtx.NewError("Error: can't have space between hash symbol and directive name")
  }

  switch pragmaWord.Value() {
  case "version":
    if !hashContext.IsAtSourceStart() {
      errCtx := context.MergeContexts(hashContext, pragmaContext)
      return nil, errCtx.NewError("Error: not at start of file")
    }

    return p.buildVersion(ts)
  case "extension":
    return p.buildExtension(ts)
  default:
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: unrecognized preprocessor directive")
  }
}

// nDir is number of tokens related to directive
func (p *GLSLParser) getDirectiveContext(ts []raw.Token, nDir int) (context.Context, error) {
  dirCtx := raw.MergeContexts(ts[0:nDir]...)
  
  if !dirCtx.IsSingleLine() {
    return dirCtx, dirCtx.NewError("Error: directive not on a single line")
  }

  if nDir < len(ts) {
    nextCtx := raw.MergeContexts(ts[0:nDir+1]...)
    if nextCtx.IsSingleLine() {
      return dirCtx, nextCtx.NewError("Error: directive has tokens after on same line")
    }
  }

  return dirCtx, nil
}

func (p *GLSLParser) buildVersion(ts []raw.Token) ([]raw.Token, error) {
  if len(ts) < 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: bad version directive")
  }

  numberToken, err := raw.AssertLiteralInt(ts[2])
  if err != nil {
    return nil, err
  }

  esToken, err := raw.AssertWord(ts[3])
  if err != nil {
    return nil, err
  }

  if esToken.Value() != "es" {
    errCtx := esToken.Context()
    return nil, errCtx.NewError("Error: must be es for webgl")
  }

  ctx, err := p.getDirectiveContext(ts, 4)
  if err != nil {
    return nil, err
  }

  st := glsl.NewVersion(numberToken.Value(), esToken.Value(), ctx)

  p.module.AddStatement(st)

  return ts[4:], nil
}

func (p *GLSLParser) buildExtension(ts []raw.Token) ([]raw.Token, error) {
  if len(ts) < 5 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: bad extension directive")
  }

  extensionNameToken, err := raw.AssertWord(ts[2])
  if err != nil {
    return nil, err
  }

  if !raw.IsSymbol(ts[3], patterns.COLON) {
    errCtx := ts[3].Context()
    return nil, errCtx.NewError("Error: expected :")
  }

  behaviorToken, err := raw.AssertWord(ts[4])
  if err != nil {
    return nil, err
  }

  ctx, err := p.getDirectiveContext(ts, 5)
  if err != nil {
    return nil, err
  }

  st := glsl.NewExtension(extensionNameToken.Value(), behaviorToken.Value(), ctx)

  p.module.AddStatement(st)

  return ts[5:], nil
}
