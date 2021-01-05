package parsers

import (
	"errors"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) assertValidPath(t raw.Token) (*js.LiteralString, error) {
	path, err := raw.AssertLiteralString(t)
	if err != nil {
		return nil, err
	}


	pathLiteral := js.NewLiteralString(path.Value(), path.Context())

	return pathLiteral, nil
}

func (p *JSParser) buildImportOrAggregateExport(t raw.Token,
	pathLiteral *js.LiteralString, lang files.Lang, fnAdder func(string, string,
		*js.LiteralString, files.Lang, context.Context) error, errMsg string) error {

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

			if err := fnAdder(name.Value(), name.Value(),
				pathLiteral, lang, context.MergeFill(t.Context(), name.Context())); err != nil {
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

			if err := fnAdder(newName.Value(), oldName.Value(),
				pathLiteral, lang, context.MergeFill(t.Context(), newName.Context())); err != nil {
				return err
			}
		} else {
			errCtx := raw.MergeContexts(bracesField...)
			return errCtx.NewError(errMsg)
		}
	}

	return nil
}

func (p *JSParser) buildRegularImportStatement(ts []raw.Token, lang files.Lang) error {
	n := len(ts)

	pathLiteral, err := p.assertValidPath(ts[n-1])
	if err != nil {
		return err
	}

	switch {
	case n == 2: // simple import
    if err := p.module.AddImportedName("", "", pathLiteral, files.SCRIPT,
      context.MergeFill(ts[0].Context(), pathLiteral.Context())); err != nil {
      return err
    }
	case n >= 3 && raw.IsWord(ts[n-2], "from"):
		fields := splitBySeparator(ts[1:n-2], patterns.COMMA)

		if len(fields) != 1 {
			errCtx := raw.MergeContexts(ts...)
			return errCtx.NewError("Error: bad import statement")
		}

    field := fields[0]

    if len(field) < 1 {
      errCtx := raw.MergeContexts(ts...)
      return errCtx.NewError("Error: bad import statement")
    }

    if len(field) == 3 && 
      raw.IsSymbol(field[0], "*") &&
			raw.IsWord(field[1], "as") &&
			raw.IsAnyWord(field[2]) {
      name, err := raw.AssertWord(field[2])
      if err != nil {
        panic(err)
      }

      if !patterns.IsJSWord(name.Value()) {
        errCtx := name.Context()
        return errCtx.NewError("Error: not a valid name")
      }

      if err := p.module.AddImportedName(name.Value(), "*",
        pathLiteral, lang, context.MergeFill(ts[0].Context(), name.Context())); err != nil {
        return err
      }

      st := js.NewImport(name.Value(), name.Context())

      p.module.AddStatement(st)
    } else if len(field) == 1 && raw.IsBracesGroup(field[0]) {
      fnAdd := func(newName, oldName string, path *js.LiteralString, lang_ files.Lang, ctx context.Context) error {
        if err := p.module.AddImportedName(newName, oldName, path, lang_, ctx); err != nil {
          return err
        }

        st := js.NewImport(newName, ctx)

        p.module.AddStatement(st)

        return nil
      }

      if err := p.buildImportOrAggregateExport(field[0], pathLiteral, lang,
        fnAdd, "Error: bad import"); err != nil {
        return err
      }

    } else {
      errCtx := raw.MergeContexts(ts...)
      return errCtx.NewError("Error: bad import statement")
    }
	}

	return nil
}

func (p *JSParser) buildNodeJSImportStatement(ts []raw.Token) error {
	n := len(ts)

	switch {
  case n == 6 && raw.IsSymbol(ts[1], "*") && raw.IsWord(ts[2], "as") && raw.IsAnyWord(ts[3]) && raw.IsWord(ts[4], "from"):

    nameToken, err := raw.AssertWord(ts[3])
    if err != nil {
      panic(err)
    }

		path, err := raw.AssertLiteralString(ts[n-1])
		if err != nil {
			return err
		}

		expr := js.NewVarExpression(nameToken.Value(), path.Context())
		statement := js.NewNodeJSImport(path.Value(), expr, path.Context())

		p.module.AddStatement(statement)
	case n == 2:
    errCtx := raw.MergeContexts(ts...)
    return errCtx.NewError("Error: no longer supported, use import * as name from '...' instead")
	default:
    errCtx := raw.MergeContexts(ts...)
    return errCtx.NewError("Error: unsupported import for builtin NodeJS module")
	}

	return nil
}

func (p *JSParser) buildImportStatement(ts []raw.Token) ([]raw.Token, error) {
	ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

	n := len(ts)
	if n < 2 {
		errCtx := ts[0].Context()
		return nil, errCtx.NewError("Error: expected more than just import;")
	}

  lang := files.SCRIPT

  if raw.IsWord(ts[n-2], "lang") && raw.IsAnyWord(ts[n-1]) {
    lastWord, err := raw.AssertWord(ts[n-1])
    if err != nil {
      panic(err)
    }
      
    switch lastWord.Value() {
    case "script", "wts", "wtscript":
      lang = files.SCRIPT
    case "template", "wtt", "wttemplate":
      lang = files.TEMPLATE
    default:
      errCtx := lastWord.Context()
      return nil, errCtx.NewError("Error: don't know how to import \"" + lastWord.Value() + "\"")
    }

    ts = ts[0:n-2]
  }

	if !raw.IsLiteralString(ts[n-1]) {
		// probably forgot semicolon
		for _, t := range ts {
			if raw.IsLiteralString(t) {
				errCtx := t.Context()
				return nil, errCtx.NewError("Error: invalid import statement, did you forget semicolon?")
			}
		}

		errCtx := ts[0].Context()
		return nil, errCtx.NewError("Error: invalid import statement, no path literal found")
	} else {
		for i, t := range ts {
			if raw.IsLiteralString(t) && i != n-1 {
				errCtx := t.Context()
				return nil, errCtx.NewError("Error: invalid import statement, did you forget semicolon?")
			}
		}
	}

	pathLiteral, err := raw.AssertLiteralString(ts[n-1])
	if err != nil {
		return nil, err
	}

	if js.IsNodeJSPackage(pathLiteral.Value()) && lang == files.SCRIPT {
		return remainingTokens, p.buildNodeJSImportStatement(ts)
	} else {
    // add literal as invisible statement, so refactoring methods can change it using the context

		return remainingTokens, p.buildRegularImportStatement(ts, lang)
	}
}

func (p *JSParser) buildExportVarStatement(ts []raw.Token, 
  varType js.VarType) ([]raw.Token, error) {
	statement, remaining, err := p.buildVarStatement(ts[1:], varType)
	if err != nil {
		return nil, err
	}

	p.module.AddStatement(statement)

	variables := statement.GetVariables()

  // export, but not default
  for k, v := range variables {
    if err := p.module.AddExportedName(k, k, v, v.Context()); err != nil {
      return nil, err
    }
  }

	return remaining, nil
}

func (p *JSParser) buildExportFunctionStatement(ts []raw.Token) ([]raw.Token, error) {
	fn, remaining, err := p.buildFunctionStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	fnVar := fn.GetVariable()

  if err := p.module.AddExportedName(fn.Name(), fn.Name(),
    fnVar, fn.Context()); err != nil {
    return nil, err
  }

	p.module.AddStatement(fn)

	return remaining, nil
}

func (p *JSParser) buildExportClassStatement(ts []raw.Token) ([]raw.Token, error) {
	cl, remaining, err := p.buildClassStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	clVar := cl.GetVariable()

  if err := p.module.AddExportedName(cl.Name(), cl.Name(),
    clVar, cl.Context()); err != nil {
    return nil, err
  }

	p.module.AddStatement(cl)

	return remaining, nil
}

func (p *JSParser) buildExportEnumStatement(ts []raw.Token) ([]raw.Token, error) {
	en, remaining, err := p.buildEnumStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	enVar := en.GetVariable()

  if err := p.module.AddExportedName(en.Name(), en.Name(),
    enVar, en.Context()); err != nil {
    return nil, err
  }

	p.module.AddStatement(en)

	return remaining, nil
}

func (p *JSParser) buildExportInterfaceStatement(ts []raw.Token) ([]raw.Token, error) {
	interf, remaining, err := p.buildInterfaceStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	interfVar := interf.GetVariable()

  if err := p.module.AddExportedName(interf.Name(), interf.Name(),
    interfVar, interf.Context()); err != nil {
    return nil, err
  }

	p.module.AddStatement(interf)

	return remaining, nil
}

func (p *JSParser) buildExportList(ts []raw.Token) ([]raw.Token, error) {
  ts, remaining := splitByNextSeparator(ts, patterns.SEMICOLON)

  // check that old names dont appear twice
  // and check that new names dont appear twice
  prevNewNames := make([]string, 0)
  prevOldNames := make([]string, 0)
  fnAdd := func(newName string, oldName string, path *js.LiteralString, lang_ files.Lang, ctx context.Context) error {
    for _, prevOldName := range prevOldNames {
      if prevOldName == oldName {
        errCtx := ctx
        return errCtx.NewError("Error: duplicate variable in export list")
      }
    }

    for _, prevNewName := range prevNewNames {
      if prevNewName == newName {
        errCtx := ctx
        return errCtx.NewError("Error: duplicate variable in export list")
      }
    }

    prevNewNames = append(prevNewNames, newName)
    prevOldNames = append(prevOldNames, oldName)

    newNameToken := js.NewWord(newName, ctx)
    varExpr := js.NewVarExpression(oldName, ctx)

    st := js.NewExport(newNameToken, varExpr, raw.MergeContexts(ts...))

    p.module.AddStatement(st)

    return nil
  }

  if err := p.buildImportOrAggregateExport(ts[1], nil, files.SCRIPT, fnAdd, "Error: bad export list"); err != nil {
    return nil, err
  }


  return remaining, nil
}

func (p *JSParser) buildExportStatement(ts []raw.Token) ([]raw.Token, error) {
	if len(ts) < 2 {
		errCtx := ts[0].Context()
		return nil, errCtx.NewError("Error: empty export statement")
	}

	switch {
	case raw.IsAnyWord(ts[1]):
		w1, err := raw.AssertWord(ts[1])
		if err != nil {
			panic(err)
		}

		switch w1.Value() {
		case "const", "let", "var":
			varType, err := js.StringToVarType(w1.Value(), w1.Context())
			if err != nil {
				panic(err)
			}

			return p.buildExportVarStatement(ts, varType)
		case "function":
			return p.buildExportFunctionStatement(ts)
		case "async":
			return p.buildExportFunctionStatement(ts)
		case "class", "abstract", "final":
			return p.buildExportClassStatement(ts)
		case "enum":
			return p.buildExportEnumStatement(ts)
		case "interface":
			return p.buildExportInterfaceStatement(ts)
		default:
      if len(ts) > 3 && raw.IsWord(ts[1], "rpc") && raw.IsWord(ts[2], "interface") {
        return p.buildExportInterfaceStatement(ts)
      }

			errCtx := ts[1].Context()
			return nil, errCtx.NewError("Error: unrecognized export statement")
		}
	// aggregate exports
	case raw.IsWord(ts[2], "from"):
		ts, remaining := splitByNextSeparator(ts, patterns.SEMICOLON)

    pathLiteral, err := p.assertValidPath(ts[3])
    if err != nil {
      return nil, err
    }

		switch {
		case raw.IsBracesGroup(ts[1]):
			if err := p.buildImportOrAggregateExport(ts[1], pathLiteral, files.SCRIPT,
				p.module.AddAggregateExport, "Error: bad aggregate export"); err != nil {
				return nil, err
			}

			return remaining, nil
		default:
			errCtx := raw.MergeContexts(ts[1:]...)
			return nil, errCtx.NewError("Error: unhandled aggregate export statement")
		}
  case raw.IsBracesGroup(ts[1]):
    return p.buildExportList(ts)
	default:
		errCtx := ts[0].Context()
		return nil, errCtx.NewError("Error: not yet handled")
	}
}

func (p *JSParser) buildModuleStatement(ts []raw.Token) ([]raw.Token, error) {
	if raw.IsAnyWord(ts[0]) {
		firstWord, err := raw.AssertWord(ts[0])
		if err != nil {
			panic(err)
		}

		switch firstWord.Value() {
		case "import":
			return p.buildImportStatement(ts)
		case "export":
			if len(ts) < 2 {
				errCtx := ts[0].Context()
				return nil, errCtx.NewError("Error: bad export statement")
			}

			if raw.IsWord(ts[1], "default") {
        errCtx := ts[1].Context()
        return nil, errCtx.NewError("Error: default exports not supported")
			} 

      return p.buildExportStatement(ts)
		case "return":
			errCtx := ts[0].Context()
			return nil, errCtx.NewError("Error: unexpected toplevel statement")
		case "continue":
			errCtx := ts[0].Context()
			return nil, errCtx.NewError("Error: unexpected toplevel statement")
		case "break":
			errCtx := ts[0].Context()
			return nil, errCtx.NewError("Error: unexpected toplevel statement")
		}
	}

	// else
	st, remaining, err := p.buildStatement(ts) // statement can be nil in case of only semicolons for example
	if err != nil {
		return nil, err
	}

	if st != nil {
		p.module.AddStatement(st)
	}

	return remaining, nil
}

func (p *JSParser) BuildModule() (*js.ModuleData, error) {
	ts, err := p.tokenize()
	if err != nil {
		return nil, err
	}

	if len(ts) < 1 {
		return nil, errors.New("Error: empty module '" +
			files.Abbreviate(p.ctx.Path()) + "'\n")
	}

	p.module = js.NewModule(ts[0].Context())

	for len(ts) > 0 {
		ts, err = p.buildModuleStatement(ts)
		if err != nil {
			return nil, err
		}
	}

	return p.module, nil
}
