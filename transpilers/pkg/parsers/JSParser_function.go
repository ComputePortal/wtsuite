package parsers

import (
	"../tokens/context"
	"../tokens/js"
	"../tokens/js/prototypes"
	"../tokens/patterns"
	"../tokens/raw"
)

func (p *JSParser) buildFunctionArgumentInner(nameToken raw.Token,
	typeTokens []raw.Token, defTokens []raw.Token) (*js.FunctionArgument, error) {
	name, err := raw.AssertWord(nameToken)
	if err != nil {
		panic(err)
	}

	var typeExpr *js.TypeExpression = nil
	if len(typeTokens) > 0 {
		typeExpr, err = p.buildTypeExpression(typeTokens)
		if err != nil {
			return nil, err
		}
	} else {
		typeExpr, err = js.NewTypeExpression("any", nil, nil, nameToken.Context())
    if err != nil {
      return nil, err
    }
	}

	var defArg js.Expression = nil
	if len(defTokens) > 0 {
		defArg, err = p.buildExpression(defTokens)
		if err != nil {
			return nil, err
		}
	}

	ctx := name.Context()

	return js.NewFunctionArgument(name.Value(), typeExpr, defArg, ctx)
}

func (p *JSParser) buildFunctionArgument(ts []raw.Token, 
	last bool) (*js.FunctionArgument, error) {
	switch {
	case len(ts) == 1 &&
		raw.IsAnyWord(ts[0]):
		return p.buildFunctionArgumentInner(ts[0], []raw.Token{}, []raw.Token{})
	case len(ts) == 2 &&
		raw.IsAnyWord(ts[0]) &&
		raw.IsAnyWord(ts[1]):
		return p.buildFunctionArgumentInner(ts[0], ts[1:], []raw.Token{})
	case len(ts) == 3 &&
		raw.IsAnyWord(ts[0]) &&
		raw.IsAnyWord(ts[1]) &&
		raw.IsAngledGroup(ts[2]):
		return p.buildFunctionArgumentInner(ts[0], ts[1:], []raw.Token{})
	case len(ts) > 2 &&
		raw.IsAnyWord(ts[0]) &&
		raw.IsSymbol(ts[1], patterns.EQUAL):
		return p.buildFunctionArgumentInner(ts[0], []raw.Token{}, ts[2:])
	case len(ts) > 3 &&
		raw.IsAnyWord(ts[0]) &&
		raw.IsAnyWord(ts[1]) &&
		raw.IsSymbol(ts[2], patterns.EQUAL):
		return p.buildFunctionArgumentInner(ts[0], ts[1:2], ts[3:])
	case len(ts) > 4 &&
		raw.IsAnyWord(ts[0]) &&
		raw.IsAnyWord(ts[1]) &&
		raw.IsAngledGroup(ts[2]) &&
		raw.IsSymbol(ts[3], patterns.EQUAL):
		return p.buildFunctionArgumentInner(ts[0], ts[1:3], ts[4:])
	case len(ts) == 2 && last &&
		raw.IsSymbol(ts[0], patterns.SPLAT) &&
		raw.IsAnyWord(ts[1]):
		// not typed, because it is always an Array
    errCtx := raw.MergeContexts(ts[0:]...)
    return nil, errCtx.NewError("Error: rest not yet supported")
	default:
		// TODO: replace all other non-rest cases with this one
		if len(ts) == 0 {
			panic("no tokens")
		}

		nameToken := ts[0]

		typeTokens := []raw.Token{}
		defTokens := []raw.Token{}

		if len(ts) > 1 {
			if raw.IsSymbol(ts[len(ts)-1], patterns.EQUAL) {
				errCtx := ts[len(ts)-1].Context()
				return nil, errCtx.NewError("Error: expected tokens after =")
			}

			for i := 1; i < len(ts)-1; i++ {
				if raw.IsSymbol(ts[i], patterns.EQUAL) {
					typeTokens = ts[1:i]
					defTokens = ts[i+1:]
					break
				} else if !(raw.IsSymbol(ts[i], patterns.PERIOD) ||
					raw.IsAngledGroup(ts[i]) ||
					raw.IsAnyWord(ts[i])) {
					errCtx := ts[i].Context()
					return nil, errCtx.NewError("Error: unexpected")
				}
			}

			if len(defTokens) == 0 {
				typeTokens = ts[1:]
			}
		}

		return p.buildFunctionArgumentInner(nameToken, typeTokens, defTokens)
	}
}

func (p *JSParser) isFunctionRoleKeyword(t raw.Token) bool {
  if !raw.IsAnyWord(t) {
    return false
  }

  if _, err := p.buildFunctionRole([]raw.Token{t}); err == nil {
    return true
  } else {
    return false
  }
}

func (p *JSParser) buildFunctionRole(ts []raw.Token) (prototypes.FunctionRole, error) {
	role := prototypes.NORMAL
	rolesDone := make(map[string]*raw.Word)

	for i := 0; i < len(ts); i++ {
		keyword, err := raw.AssertWord(ts[i])
		if err != nil {
			return role, err
		}

		switch keyword.Value() {
		case "const":
			role = role | prototypes.CONST
		case "static":
			role = role | prototypes.STATIC
		case "get":
			role = role | prototypes.GETTER
		case "set":
			role = role | prototypes.SETTER
		case "public":
			role = role | prototypes.PUBLIC
		case "private":
			role = role | prototypes.PRIVATE
		case "abstract":
			role = role | prototypes.ABSTRACT
		case "override":
			role = role | prototypes.OVERRIDE
		case "async":
			role = role | prototypes.ASYNC
		case "function":
			if i != len(ts)-1 {
				errCtx := keyword.Context()
				return role, errCtx.NewError("Error: function keyword must come after roles")
			}
		default:
			errCtx := keyword.Context()
			return role, errCtx.NewError("Error: unexpected role")
		}

		if prev, ok := rolesDone[keyword.Value()]; ok {
			errCtx := context.MergeContexts(prev.Context(), keyword.Context())
			return role, errCtx.NewError("Error: duplicate member role specification")
		} else {
			rolesDone[keyword.Value()] = keyword
		}
	}

	return role, nil
}

func (p *JSParser) buildFunctionInterface(ts []raw.Token,
	named bool, ctx context.Context) (*js.FunctionInterface, []raw.Token, error) {
	if len(ts) == 0 {
		return nil, nil, ctx.NewError("Error: bad function interface")
	}

	// find first parens and cut there or one before
	parensPos := -1
	for i, t := range ts {
		if raw.IsParensGroup(t) {
			parensPos = i
			break
		}
	}
	if parensPos == -1 {
		return nil, nil, ctx.NewError("Error: bad function interface")
	}

	fnName := ""
	rolePos := parensPos
	if named {
		if parensPos > 0 && raw.IsAnyWord(ts[parensPos-1]) {
			nameToken, err := raw.AssertWord(ts[parensPos-1])
			if err != nil {
				panic(err)
			}

			fnName = nameToken.Value()
			rolePos -= 1
		} else {
			errCtx := raw.MergeContexts(ts...)
			return nil, nil, errCtx.NewError("Error: no function name found")
		}
	}

	// the last keyword can be 'function', which is optional and thus ignored
	role, err := p.buildFunctionRole(ts[0:rolePos])
	if err != nil {
		return nil, nil, err
	}

	ts = ts[parensPos:]

	if len(ts) == 0 {
		panic("shouldn't happen")
		//return nil, nil, ctx.NewError("Error: bad function interface")
	}

	argGroup, err := raw.AssertParensGroup(ts[0])
	if err != nil {
		panic("shouldn't happen")
		//return nil, err
	}

	if argGroup.IsSemiColon() {
		errCtx := argGroup.Context()
		return nil, nil, errCtx.NewError("Error: can't use semicolons in function arg list")
	}

	ts = ts[1:]

	err = nil

	fnInterf := js.NewFunctionInterface(fnName, role, ctx)
	for i, field := range argGroup.Fields {
		if len(field) == 0 {
			errCtx := ctx
			return nil, nil, errCtx.NewError("Error: empty function argument")
		}

		arg, err := p.buildFunctionArgument(field, i == len(argGroup.Fields)-1)
		if err != nil {
			return nil, nil, err
		}

		fnInterf.AppendArg(arg)
	}

	// return value might be specified
	if len(ts) > 0 && raw.IsAnyWord(ts[0]) {
		iLastRetTypeToken := 1
		for i, t := range ts {
			if raw.IsSymbol(t, patterns.ARROW) || raw.IsBracesGroup(t) {
				break
			}
			iLastRetTypeToken = i
		}

		retType, err := p.buildTypeExpression(ts[0:iLastRetTypeToken+1])
		if err != nil {
			return nil, nil, err
		}

		fnInterf.SetReturnType(retType)
		ts = ts[iLastRetTypeToken+1:]
	}

	return fnInterf, ts, nil
}

// js.Function can be used both as expression and statement
func (p *JSParser) buildFunction(ts []raw.Token, named bool,
	isArrow bool) (*js.Function, []raw.Token, error) {

	// get the correct function context
	iFirstBrace := 0
	for i, t := range ts {
		if raw.IsBracesGroup(t) {
			iFirstBrace = i
			break
		}
	}

	fnCtx := raw.MergeContexts(ts[0 : iFirstBrace+1]...)

	fnInterf, ts, err := p.buildFunctionInterface(ts, named, fnCtx)
	if err != nil {
		return nil, nil, err
	}

	function, err := js.NewFunction(fnInterf, isArrow, fnCtx)
	if err != nil {
		return nil, nil, err
	}

	if len(ts) < 1 {
		return nil, nil, fnCtx.NewError("Error: bad function")
	}

	remaining := ts[1:]

	statementGroup, err := raw.AssertBracesGroup(ts[0])
	if err != nil {
		return nil, nil, err
	}

	if !statementGroup.IsEmpty() {
		statements, err := p.buildBlockStatements(statementGroup)
		if err != nil {
			return nil, nil, err
		}

		for _, statement := range statements {
			function.AddStatement(statement)
		}
	}

	return function, remaining, nil
}

// dont ever name function expression, too obscure functionality anyway
func (p *JSParser) buildFunctionExpression(ts []raw.Token, isArrow bool) (*js.Function,
	error) {
	expr, remaining, err := p.buildFunction(ts, false, isArrow)
	if err != nil {
		return nil, err
	}

	if len(remaining) != 0 {
		errCtx := raw.MergeContexts(remaining...)
		return nil, errCtx.NewError("Error: unexpected tokens after function expression")
	}

	if expr.Role() != prototypes.NORMAL && expr.Role() != prototypes.ASYNC {
		errCtx := expr.Context()
		return nil, errCtx.NewError("Error: illegal function expression role(s)")
	}
	return expr, nil
}

// should be able to handle all the keywords
func (p *JSParser) buildFunctionStatement(ts []raw.Token) (*js.Function,
	[]raw.Token, error) {
	if len(ts) < 4 {
		errCtx := raw.MergeContexts(ts...)
		return nil, nil, errCtx.NewError("Error: bad function statement")
	}

	//if !(raw.IsAnyWord(ts[1]) && raw.IsParensGroup(ts[2]) && raw.IsBracesGroup(ts[3])) {
	//errCtx := raw.MergeContexts(ts...)
	//return nil, nil, errCtx.NewError("Error: bad function statement")
	//}

	//remaining := p.stripSeparators(4, ts, patterns.SEMICOLON)

	fn, remaining, err := p.buildFunction(ts, true, false)
	if err != nil {
		return nil, nil, err
	}

	if fn.Role() != prototypes.NORMAL && fn.Role() != prototypes.ASYNC {
		errCtx := fn.Context()
		return nil, nil, errCtx.NewError("Error: illegal function statement role(s)")
	}

	return fn, remaining, nil
}
