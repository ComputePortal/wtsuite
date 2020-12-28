package parsers

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) assertValidPath(t raw.Token) (*js.LiteralString, bool, error) {
	path, err := raw.AssertLiteralString(t)
	if err != nil {
		return nil, false, err
	}

	pathCtx := path.Context()
	absPath, isPackage, err := files.SearchPackage(pathCtx.Caller(), path.Value(), files.JSPACKAGE_SUFFIX)
	if err != nil {
		if absPath, isPackage, err = files.SearchPackage(pathCtx.Caller(), path.Value(), files.UIPACKAGE_SUFFIX); err != nil {
			errCtx := path.InnerContext()
			return nil, false, errCtx.NewError("Error: file not found")
		}
	}

	pathLiteral := js.NewLiteralString(absPath, pathCtx)
	return pathLiteral, isPackage, nil
}

func (p *JSParser) buildImportOrAggregateExport(t raw.Token,
	pathLiteral *js.LiteralString, fnAdder func(string, string,
		*js.LiteralString, context.Context) error, errMsg string) error {

	group, err := raw.AssertBracesGroup(t)
	if err != nil {
		panic(err)
	}

	if !(group.IsSingle() || group.IsComma()) {
		errCtx := group.Context()
		return errCtx.NewError("Error: expected single or comma braces")
	}

	for _, bracesField := range group.Fields {
		if len(bracesField) == 1 && raw.IsAnyWord(bracesField[0]) {
			name, err := raw.AssertWord(bracesField[0])
			if err != nil {
				panic(err)
			}

			if err := fnAdder(name.Value(), name.Value(),
				pathLiteral, context.MergeFill(t.Context(), name.Context())); err != nil {
				return err
			}
		} else if len(bracesField) == 3 &&
			raw.IsAnyWord(bracesField[0]) &&
			raw.IsWord(bracesField[1], "as") &&
			raw.IsAnyWord(bracesField[2]) {
			oldName, err := raw.AssertWord(bracesField[0])
			if err != nil {
				panic(err)
			}

			newName, err := raw.AssertWord(bracesField[2])
			if err != nil {
				panic(err)
			}

			if err := fnAdder(newName.Value(), oldName.Value(),
				pathLiteral, context.MergeFill(t.Context(), newName.Context())); err != nil {
				return err
			}
		} else {
			errCtx := raw.MergeContexts(bracesField...)
			return errCtx.NewError(errMsg)
		}
	}

	return nil
}

func (p *JSParser) buildRegularImportStatement(ts []raw.Token) error {
	n := len(ts)

	pathLiteral, isPackage, err := p.assertValidPath(ts[n-1])
	if err != nil {
		return err
	}

	switch {
	case n == 2: // simple import
		// if the path is a directory then we use this as a syntactic sugar for a package import
		// otherwise this import is just for side-effects
		if isPackage {
			name := filepath.Base(filepath.Dir(pathLiteral.Value()))
			if err := p.module.AddImportedName(name, "*", pathLiteral,
				context.MergeFill(ts[0].Context(), pathLiteral.Context())); err != nil {
				return err
			}
		} else {
			if err := p.module.AddImportedName("", "", pathLiteral,
				context.MergeFill(ts[0].Context(), pathLiteral.Context())); err != nil {
				return err
			}
		}
	case n >= 3 && raw.IsWord(ts[n-2], "from"):
		fields := splitBySeparator(ts[1:n-2], patterns.COMMA)

		if len(fields) < 1 || len(fields) > 3 {
			errCtx := raw.MergeContexts(ts...)
			return errCtx.NewError("Error: bad import statement fields")
		}

		starDone := false
		bracesDone := false

		for i, field := range fields {
			if len(field) < 1 {
				errCtx := raw.MergeContexts(ts...)
				return errCtx.NewError("Error: bad import statement fields")
			}

			switch {
			case len(field) == 1 && raw.IsAnyWord(field[0]): // default field (todo: avoid keywords
				if i != 0 {
					errCtx := field[0].Context()
					return errCtx.NewError("Error: default import must come first")
				}

				name, err := raw.AssertWord(field[0])
				if err != nil {
					panic(err)
				}

				if err := p.module.AddImportedName(name.Value(), "default",
					pathLiteral, context.MergeFill(ts[0].Context(), name.Context())); err != nil {
					return err
				}
			case len(field) == 3 &&
				raw.IsSymbol(field[0], "*") &&
				raw.IsWord(field[1], "as") &&
				raw.IsAnyWord(field[2]):
				if starDone {
					errCtx := field[0].Context()
					return errCtx.NewError("Error: wildcard import already done")
				}

				name, err := raw.AssertWord(field[2])
				if err != nil {
					panic(err)
				}

				// TODO: a valid name check

				if err := p.module.AddImportedName(name.Value(), "*",
					pathLiteral, context.MergeFill(ts[0].Context(), name.Context())); err != nil {
					return err
				}
				starDone = true
			case len(field) == 1 && raw.IsBracesGroup(field[0]):
				if bracesDone {
					errCtx := field[0].Context()
					return errCtx.NewError("Error: please combine all brace parts into one")
				}

				if err := p.buildImportOrAggregateExport(field[0], pathLiteral,
					p.module.AddImportedName, "Error: bad import"); err != nil {
					return err
				}

				bracesDone = true
			}
		}
	}

	return nil
}

func (p *JSParser) buildNodeJSImportStatement(ts []raw.Token) error {
	n := len(ts)

	switch n {
	case 2:
		path, err := raw.AssertLiteralString(ts[n-1])
		if err != nil {
			return err
		}

    if VERBOSITY >= 3 {
      fmt.Println("importing nodejs module " + path.Value())
    }

		expr := js.NewVarExpression(path.Value(), path.Context())
		statement := js.NewNodeJSImport(expr, path.Context())

		p.module.AddStatement(statement)
	default:
		panic("not yet implemented")
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

	if js.IsNodeJSPackage(pathLiteral.Value()) {
		return remainingTokens, p.buildNodeJSImportStatement(ts)
	} else {
    // add literal as invisible statement, so refactoring methods can change it using the context

		return remainingTokens, p.buildRegularImportStatement(ts)
	}
}

func (p *JSParser) buildExportVarStatement(ts []raw.Token, varType js.VarType,
	isDefault bool) ([]raw.Token, error) {
	statement, remaining, err := p.buildVarStatement(ts[1:], varType)
	if err != nil {
		return nil, err
	}

	p.module.AddStatement(statement)

	variables := statement.GetVariables()

	if isDefault {
		if len(variables) > 1 {
			var first js.Variable = nil
			var second js.Variable = nil

			for _, v := range variables {
				if first == nil {
					first = v
				} else if second == nil {
					second = v
				} else {
					break
				}
			}

			errCtx := context.MergeContexts(ts[0].Context(), first.Context(), second.Context())
			return nil, errCtx.NewError("Error: there can only be one default")
		} else if len(variables) == 0 {
			errCtx := ts[0].Context()
			return nil, errCtx.NewError("Error: there must be at least one variable after 'default'")
		}

		for k, v := range variables {
			if err := p.module.AddExportedName("default", k, v, v.Context()); err != nil {
				return nil, err
			}
		}
	} else {
		// export, but not default
		for k, v := range variables {
			if err := p.module.AddExportedName(k, k, v, v.Context()); err != nil {
				return nil, err
			}
		}
	}

	return remaining, nil
}

func (p *JSParser) buildExportFunctionStatement(ts []raw.Token,
	isDefault bool) ([]raw.Token, error) {
	fn, remaining, err := p.buildFunctionStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	fnVar := fn.GetVariable()

	if isDefault {
		if err := p.module.AddExportedName("default", fn.Name(),
			fnVar, fn.Context()); err != nil {
			return nil, err
		}
	} else {
		if err := p.module.AddExportedName(fn.Name(), fn.Name(),
			fnVar, fn.Context()); err != nil {
			return nil, err
		}
	}

	p.module.AddStatement(fn)

	return remaining, nil
}

func (p *JSParser) buildExportClassStatement(ts []raw.Token,
	isDefault bool) ([]raw.Token, error) {
	cl, remaining, err := p.buildClassStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	clVar := cl.GetVariable()

	if isDefault {
		if err := p.module.AddExportedName("default", cl.Name(),
			clVar, cl.Context()); err != nil {
			return nil, err
		}
	} else {
		if err := p.module.AddExportedName(cl.Name(), cl.Name(),
			clVar, cl.Context()); err != nil {
			return nil, err
		}
	}

	p.module.AddStatement(cl)

	return remaining, nil
}

func (p *JSParser) buildExportEnumStatement(ts []raw.Token,
	isDefault bool) ([]raw.Token, error) {
	en, remaining, err := p.buildEnumStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	enVar := en.GetVariable()

	if isDefault {
		if err := p.module.AddExportedName("default", en.Name(),
			enVar, en.Context()); err != nil {
			return nil, err
		}
	} else {
		if err := p.module.AddExportedName(en.Name(), en.Name(),
			enVar, en.Context()); err != nil {
			return nil, err
		}
	}

	p.module.AddStatement(en)

	return remaining, nil
}

func (p *JSParser) buildExportInterfaceStatement(ts []raw.Token,
	isDefault bool) ([]raw.Token, error) {
	interf, remaining, err := p.buildInterfaceStatement(ts[1:])
	if err != nil {
		return nil, err
	}

	interfVar := interf.GetVariable()

	if isDefault {
		if err := p.module.AddExportedName("default", interf.Name(),
			interfVar, interf.Context()); err != nil {
			return nil, err
		}
	} else {
		if err := p.module.AddExportedName(interf.Name(), interf.Name(),
			interfVar, interf.Context()); err != nil {
			return nil, err
		}
	}

	p.module.AddStatement(interf)

	return remaining, nil
}

func (p *JSParser) buildExportStatement(ts []raw.Token,
	isDefault bool) ([]raw.Token, error) {
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

			return p.buildExportVarStatement(ts, varType, isDefault)
		case "function":
			return p.buildExportFunctionStatement(ts, isDefault)
		case "async":
			return p.buildExportFunctionStatement(ts, isDefault)
		case "class", "abstract", "final":
			return p.buildExportClassStatement(ts, isDefault)
		case "enum":
			return p.buildExportEnumStatement(ts, isDefault)
		case "interface":
			return p.buildExportInterfaceStatement(ts, isDefault)
		default:
      if len(ts) > 3 && raw.IsWord(ts[1], "rpc") && raw.IsWord(ts[2], "interface") {
        return p.buildExportInterfaceStatement(ts, isDefault)
      }

			errCtx := ts[1].Context()
			return nil, errCtx.NewError("Error: unrecognized export statement")
		}
	// aggregate exports
	case raw.IsWord(ts[2], "from"):
		ts, remaining := splitByNextSeparator(ts, patterns.SEMICOLON)

		pathLiteral, _, err := p.assertValidPath(ts[3])
		if err != nil {
			return nil, err
		}

		switch {
		case raw.IsBracesGroup(ts[1]):
			if err := p.buildImportOrAggregateExport(ts[1], pathLiteral,
				p.module.AddAggregateExport, "Error: bad aggregate export"); err != nil {
				return nil, err
			}

			return remaining, nil
		default:
			errCtx := raw.MergeContexts(ts[1:]...)
			return nil, errCtx.NewError("Error: unhandled aggregate export statement")
		}
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
				return p.buildExportStatement(ts[1:], true)
			} else {
				return p.buildExportStatement(ts, false)
			}
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
