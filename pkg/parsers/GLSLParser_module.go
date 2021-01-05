package parsers

import (
  "errors"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildImportExportList(t raw.Token, fn func(*glsl.Word, *glsl.Word, context.Context) error, errMsg string) error {
	group, err := raw.AssertBracesGroup(t)
	if err != nil {
		panic(err)
	}

	if !(group.IsSingle() || group.IsComma()) {
		errCtx := group.Context()
		return errCtx.NewError("Error: expected single or comma braces")
	}

	for _, bracesField := range group.Fields {
    asPos := -1
    for i, t := range bracesField {
      if raw.IsWord(t, "as") {
        asPos = i
        break
      }
    }

		if asPos == -1 {
			name, rem, err := condensePackagePeriods(bracesField)
			if err != nil {
				return err
			}

      if len(rem) != 0 {
        errCtx := raw.MergeContexts(rem...)
        return errCtx.NewError("Error: unexpected tokens")
      }

			if err := fn(glsl.NewWord(name.Value(), name.Context()), 
        glsl.NewWord(name.Value(), name.Context()),
				context.MergeFill(t.Context(), name.Context())); err != nil {
				return err
			}
		} else if asPos < len(bracesField) - 1 {
      oldName, rem, err := condensePackagePeriods(bracesField[0:asPos])
      if err != nil {
        return err
      }

      if len(rem) != 0 {
        errCtx := raw.MergeContexts(rem...)
        return errCtx.NewError("Error: unexpected tokens")
      }

			newName, err := raw.AssertWord(bracesField[asPos+1])
			if err != nil {
				return err
			}

			if err := fn(glsl.NewWord(newName.Value(), newName.Context()), 
        glsl.NewWord(oldName.Value(), oldName.Context()),
				context.MergeFill(t.Context(), newName.Context())); err != nil {
				return err
			}
		} else {
			errCtx := raw.MergeContexts(bracesField...)
			return errCtx.NewError(errMsg)
		}
	}

	return nil
}

func (p *GLSLParser) buildExportList(ts []raw.Token) error {
  fnAdd := func(newName *glsl.Word, oldName *glsl.Word, ctx context.Context) error {
    st := glsl.NewExport(newName, glsl.NewVarExpression(oldName.Value(), oldName.Context()), ctx)

    p.module.AddStatement(st)

    return nil
  }

  if err := p.buildImportExportList(ts[1], fnAdd, "Error: bad export list"); err != nil {
    return err
  }

  return nil
}

func (p *GLSLParser) buildAggregateExport(ts []raw.Token) error {
  path_, err := raw.AssertLiteralString(ts[3])
  if err != nil {
    return err

  }

  path := glsl.NewLiteralString(path_.Value(), path_.Context())

  fnAdd := func(newName *glsl.Word, oldName *glsl.Word, ctx context.Context) error {
    st, err := glsl.NewImportExport(newName, glsl.NewVarExpression(oldName.Value(), oldName.Context()), path, ctx)
    if err != nil {
      return err
    }

    p.module.AddStatement(st)

    return nil
  }

  if err := p.buildImportExportList(ts[1], fnAdd, "Error: bad aggregate export"); err != nil {
    return err
  }

  return nil
}

func (p *GLSLParser) buildExportStatement(ts []raw.Token) ([]raw.Token, error) {
  ilast := len(ts)
  for i, t := range ts {
    if raw.IsSymbol(t, patterns.SEMICOLON) {
      ilast = i
      break
    } else if raw.IsBracesGroup(t) && raw.IsParensGroup(ts[i-1]) {
      ilast = i + 1
      break
    }
  }

  remainingTokens := ts[ilast:]
  ts = ts[0:ilast]

  switch {
  case len(ts) == 2 && raw.IsBracesGroup(ts[1]):
    if err := p.buildExportList(ts); err != nil {
      return nil, err
    }
    return remainingTokens, nil
  case len(ts) == 4 && raw.IsBracesGroup(ts[1]) && raw.IsWord(ts[2], "from") && raw.IsLiteralString(ts[3]):
    if err := p.buildAggregateExport(ts); err != nil {
      return nil, err
    }
    return remainingTokens, nil
  case raw.IsWord(ts[1], "export"):
    errCtx := ts[1].Context()
    return nil, errCtx.NewError("Error: unexpected")
  default:
    // cut off the export part
    rem, err := p.buildModuleStatement(ts[1:], true)
    if err != nil {
      return nil, err
    }

    if len(rem) != 0 {
      errCtx := raw.MergeContexts(rem...)
      return nil, errCtx.NewError("Error: unexpected tokens")
    }

    return remainingTokens, nil
  }
}

func (p *GLSLParser) buildSimpleExport(name string, ctx context.Context) error {
  st := glsl.NewExport(glsl.NewWord(name, ctx), glsl.NewVarExpression(name, ctx), ctx)

  p.module.AddStatement(st)

  return nil
}

func (p *GLSLParser) buildNamedImports(ts []raw.Token) error {
  path_, err := raw.AssertLiteralString(ts[3])
  if err != nil {
    return err
  }

  path := glsl.NewLiteralString(path_.Value(), path_.Context())

  fnAdd := func(newName *glsl.Word, oldName *glsl.Word, ctx context.Context) error {
    st, err := glsl.NewImport(newName, glsl.NewVarExpression(oldName.Value(), oldName.Context()), path, ctx)
    if err != nil {
      return err
    }

    p.module.AddStatement(st)
    return nil
  }

  if err := p.buildImportExportList(ts[1], fnAdd, "Error: bad named import"); err != nil {
    return err
  }

  return nil
}

func (p *GLSLParser) buildNamespaceImport(ts []raw.Token) error {
  path_, err := raw.AssertLiteralString(ts[5])
  if err != nil {
    return err
  }

  path := glsl.NewLiteralString(path_.Value(), path_.Context())

  nameToken, err := raw.AssertWord(ts[3])
  if err != nil {
    return err
  }

  nameWord := glsl.NewWord(nameToken.Value(), nameToken.Context())

  st, err := glsl.NewImport(nameWord, glsl.NewVarExpression("*", ts[1].Context()), path, raw.MergeContexts(ts...))
  if err != nil {
    return err
  }

  p.module.AddStatement(st)

  return nil
}

func (p *GLSLParser) buildImportStatement(ts []raw.Token) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  switch {
  case len(ts) == 4 && raw.IsBracesGroup(ts[1]) && raw.IsWord(ts[2], "from") && raw.IsLiteralString(ts[3]):
    if err := p.buildNamedImports(ts); err != nil {
      return nil, err
    }
    return remainingTokens, nil
  case len(ts) == 6 && raw.IsSymbol(ts[1], "*") && raw.IsWord(ts[2], "as") && raw.IsAnyWord(ts[3]) && raw.IsWord(ts[4], "from") && raw.IsLiteralString(ts[5]):
    if err := p.buildNamespaceImport(ts); err != nil {
      return nil, err
    }
    return remainingTokens, nil
  default:
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: bad import statement")
  }
}

func (p *GLSLParser) buildModuleStatement(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts = p.expandTmpGroups(ts)

  for raw.IsSymbol(ts[0], patterns.SEMICOLON) {
    ts = ts[1:]
  }

  if raw.IsAnyWord(ts[0]) {
    firstWord, err := raw.AssertWord(ts[0])
    if err != nil {
      panic(err)
    }

    switch firstWord.Value() {
    case "export":
      if isExport {
        errCtx := ts[0].Context()
        err := errCtx.NewError("Error: unexpected token")
        return nil, err
      }
      return p.buildExportStatement(ts)
    case "import":
      if isExport {
        errCtx := ts[0].Context()
        return nil, errCtx.NewError("Error: unexpected token")
      }

      return p.buildImportStatement(ts)
    case "attribute":
      return p.buildAttribute(ts, isExport)
    case "varying":
      return p.buildVarying(ts, isExport)
    case "uniform":
      return p.buildUniform(ts, isExport)
    case "const":
      return p.buildConst(ts, isExport)
    case "struct":
      return p.buildStruct(ts, isExport)
    case "precision":
      if isExport {
        errCtx := ts[0].Context()
        return nil, errCtx.NewError("Error: can't be exported")
      }

      return p.buildPrecision(ts)
    default:
      if iparens := raw.FindFirstParensGroup(ts, 0); iparens != -1 && iparens < len(ts) - 1 && raw.IsBracesGroup(ts[iparens + 1]) {
        allWordsOrPeriods := true
        for _, t := range ts[0:iparens] {
          if !raw.IsAnyWord(t) && !raw.IsSymbol(t, patterns.PERIOD) {
            allWordsOrPeriods = false
            break
          }
        }

        if allWordsOrPeriods {
          return p.buildFunction(ts, isExport)
        } 
      }

      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: unrecognized top level statement")
    }
  } else if len(ts) >= 2 && raw.IsSymbol(ts[0], patterns.HASH) && raw.IsAnyWord(ts[1]) {
    if isExport {
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: can't be exported")
    }
    return p.buildDirective(ts)
  } else {
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: all top level statements start with a word")
  }
}

func (p *GLSLParser) BuildModule() (*glsl.ModuleData, error) {
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
    ts, err = p.buildModuleStatement(ts, false)
    if err != nil {
      return nil, err
    }
  }

  return p.module, nil
}
