package parsers

import (
  "errors"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildModuleStatement(ts []raw.Token) ([]raw.Token, error) {
  ts = p.expandTmpGroups(ts)

  if raw.IsAnyWord(ts[0]) {
    firstWord, err := raw.AssertWord(ts[0])
    if err != nil {
      panic(err)
    }

    switch firstWord.Value() {
    case "export":
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: export not yet supported")
    case "import":
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: import not yet supported")
    case "void", "int":
      if len(ts) > 2 && raw.IsAnyWord(ts[1]) && raw.IsParensGroup(ts[2]) && raw.IsBracesGroup(ts[3]) {
        return p.buildFunction(ts)
      } else {
        errCtx := ts[0].Context()
        return nil, errCtx.NewError("Error: expected function after " + firstWord.Value())
      }
    case "attribute":
      return p.buildAttribute(ts)
    case "varying":
      return p.buildVarying(ts)
    case "uniform":
      return p.buildUniform(ts)
    case "const":
      return p.buildConst(ts)
    default:
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: unrecognized top level statement")
    }
  } else {
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: all top level statements start with a word")
  }
}

func (p *GLSLParser) BuildModule() (*glsl.Module, error) {
  ts, err := p.tokenize()
  if err != nil {
    return nil, err
  }

	if len(ts) < 1 {
		return nil, errors.New("Error: empty module '" +
			files.Abbreviate(p.ctx.Path()) + "'\n")
	}

  p.module = glsl.NewModule(ts[0].Context())

  for len(ts) > 0 {
    ts, err = p.buildModuleStatement(ts)
    if err != nil {
      return nil, err
    }
  }

  return p.module, nil
}
