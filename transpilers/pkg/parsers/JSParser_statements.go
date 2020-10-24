package parsers

import (
	"strings"

	"../tokens/js"
	"../tokens/js/macros"
	"../tokens/patterns"
	"../tokens/raw"
)

func (p *JSParser) buildVarStatement(ts []raw.Token,
	varType js.VarType) (*js.VarStatement, []raw.Token, error) {
	n := len(ts)
	if n < 2 {
		errCtx := ts[0].Context()
		return nil, nil, errCtx.NewError("Error: expected at least a name")
	}

	ts, remainingTokens := p.splitByNextSeparator(ts, patterns.SEMICOLON)
	fields := p.splitBySeparator(ts[1:], patterns.COMMA)

	expressions := make([]js.Expression, 0)

	for _, field := range fields {
		if len(field) < 1 {
			errCtx := ts[0].Context()
			return nil, nil, errCtx.NewError("Error: bad " +
				js.VarTypeToString(varType) + " declaration")
		}

		switch {
		case len(field) == 1 &&
			raw.IsAnyWord(field[0]):
			name, err := raw.AssertWord(field[0])
			if err != nil {
				panic(err)
			}

			if varType == js.CONST {
				errCtx := name.Context()
				return nil, nil, errCtx.NewError("Error: pointless const declaration " +
					"(hint: include rhs)")
			}

			if macros.IsOnlyMacro(name.Value()) {
				errCtx := name.Context()
				return nil, nil, errCtx.NewError("Error: macro " + name.Value() + " is not variable")
			}

			expressions = append(expressions, js.NewVarExpression(name.Value(), name.Context()))
		case len(field) > 2 &&
			raw.IsAnyWord(field[0]) &&
			raw.IsSymbol(field[1], patterns.EQUAL):
			name, err := raw.AssertWord(field[0])
			if err != nil {
				panic(err)
			}

			// turn rhs into an expression
			rhs, err := p.buildExpression(field[2:])
			if err != nil {
				return nil, nil, err
			}

			if macros.IsOnlyMacro(name.Value()) {
				errCtx := name.Context()
				return nil, nil, errCtx.NewError("Error: macro " + name.Value() + " is not variable")
			}

			expressions = append(expressions, js.NewAssign(js.NewVarExpression(name.Value(),
				name.Context()), rhs, "", field[1].Context()))
		default:
      iEqual := 0
			for i, t := range field {
				if raw.IsSymbol(t, patterns.EQUAL) {
          if i == len(field)-1 {
            errCtx := t.Context()
            return nil, nil, errCtx.NewError("Error: expected expression after")
          }

          iEqual = i
					break
				}
			}

      if raw.IsAnyWord(field[0]) {
        name, err := raw.AssertWord(field[0])
        if err != nil {
          panic(err)
        }

        if macros.IsOnlyMacro(name.Value()) {
          errCtx := name.Context()
          return nil, nil, errCtx.NewError("Error: macro " + name.Value() + " is not variable")
        }

        typeExpr, err := p.buildTypeExpression(field[1:iEqual], false)
        if err != nil {
          return nil, nil, err
        }

        rhs, err := p.buildExpression(field[iEqual+1:])
        if err != nil {
          return nil, nil, err
        }

        expressions = append(expressions, 
          js.NewTypedAssign(js.NewVarExpression(name.Value(),
            name.Context()), typeExpr, rhs, field[iEqual].Context(),
          ),
        )
      } else {
        errCtx := raw.MergeContexts(field...)
        err := errCtx.NewError("Error: not yet supported")
        return nil, nil, err
      }
		}
	}

	statement, err := js.NewVarStatement(varType, expressions, ts[0].Context())
	if err != nil {
		return nil, nil, err
	}

	return statement, remainingTokens, nil
}

func (p *JSParser) buildReturnStatement(ts []raw.Token) (*js.Return, []raw.Token, error) {
	exprTokens, remainingTokens := p.splitByNextSeparator(ts[1:], patterns.SEMICOLON)
	var expr js.Expression = nil
	if len(exprTokens) > 0 {
		var err error
		expr, err = p.buildExpression(exprTokens)
		if err != nil {
			return nil, nil, err
		}
	}

	retStatement, err := js.NewReturn(expr, ts[0].Context())
	if err != nil {
		return nil, nil, err
	}

	return retStatement, remainingTokens, nil
}

func (p *JSParser) buildThrowStatement(ts []raw.Token) (*js.Throw, []raw.Token, error) {
	exprTokens, remainingTokens := p.splitByNextSeparator(ts[1:], patterns.SEMICOLON)
	if len(exprTokens) > 0 {
		expr, err := p.buildExpression(exprTokens)
		if err != nil {
			return nil, nil, err
		}

		throwStatement, err := js.NewThrow(expr, ts[0].Context())
		if err != nil {
			return nil, nil, err
		}

		return throwStatement, remainingTokens, nil
	} else {
		errCtx := ts[0].Context()
		return nil, nil, errCtx.NewError("Error: expected 1 argument")
	}
}

func (p *JSParser) buildBreakStatement(ts []raw.Token) (*js.Break, []raw.Token, error) {
	exprTokens, remainingTokens := p.splitByNextSeparator(ts, patterns.SEMICOLON)
	if len(exprTokens) != 1 {
		errCtx := raw.MergeContexts(ts...)
		return nil, nil, errCtx.NewError("Error: bad break statement")
	}

	breakStatement, err := js.NewBreak(exprTokens[0].Context())
	if err != nil {
		return nil, nil, err
	}

	return breakStatement, remainingTokens, nil
}

func (p *JSParser) buildContinueStatement(ts []raw.Token) (*js.Continue, []raw.Token, error) {
	exprTokens, remainingTokens := p.splitByNextSeparator(ts, patterns.SEMICOLON)
	if len(exprTokens) != 1 {
		errCtx := raw.MergeContexts(ts...)
		return nil, nil, errCtx.NewError("Error: bad continue statement")
	}

	continueStatement, err := js.NewContinue(exprTokens[0].Context())
	if err != nil {
		return nil, nil, err
	}

	return continueStatement, remainingTokens, nil
}

func (p *JSParser) buildAssignStatementLHS(ts []raw.Token) (js.Expression, *js.TypeExpression, []raw.Token, error) {
  iEqual := -1
  for i, t := range ts {
    if raw.IsSymbolThatEndsWith(t, patterns.EQUAL) {
      iEqual = i
      break
    }
  }

  if iEqual <= 0 || iEqual == len(ts) - 1 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, nil, errCtx.NewError("Error: invalid assign")
  }

  lhsTokens := ts[0:iEqual]
  nonLHSTokens := ts[iEqual:] // includes *Equals, so op can be extracted

  // if it only contains words, periods, and angled groups, then try parsing TypeExpression after second consecutive word
  prevWasWord := false
  for i := 0; i < len(lhsTokens); i++ {
    t := lhsTokens[i]
    if !(raw.IsAnyWord(t) || raw.IsSymbol(t, patterns.PERIOD)) {
      break
    }

    isWord := raw.IsAnyWord(t)
    if isWord {
      if prevWasWord {
        lhs, err := p.buildExpression(lhsTokens[0:i])
        if err != nil {
          return nil, nil, nil, err
        }

        typeExpr, err := p.buildTypeExpression(lhsTokens[i:], false)
        if err != nil {
          return nil, nil, nil, err
        }

        return lhs, typeExpr, nonLHSTokens, nil
      } 
    }

    prevWasWord = isWord
  }

  // no typeExpression
  lhs, err := p.buildExpression(lhsTokens)
  if err != nil {
    return nil, nil, nil, err
  }

  return lhs, nil, nonLHSTokens, nil
}

func (p *JSParser) buildAssignStatement(ts_ []raw.Token) (*js.Assign, []raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts_, patterns.SEMICOLON)

	if len(ts) < 3 {
		errCtx := ts[1].Context()
		return nil, nil, errCtx.NewError("Error: assign statement expects rhs")
	}

  lhs, typeExpr, nonLHSTokens, err := p.buildAssignStatementLHS(ts)
  if err != nil {
    return nil, nil, err
  }

	symbol, err := raw.AssertAnySymbol(nonLHSTokens[0])
	if err != nil {
		panic(err)
	}

	op := strings.TrimSuffix(symbol.Value(), patterns.EQUAL)

	rhs, err := p.buildExpression(nonLHSTokens[1:])
	if err != nil {
		return nil, nil, err
	}

  if op != "" && typeExpr != nil {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: mutating assign can't be typed")
  } else if op == "!" || op == "<" || op == ">" || op == "!=" || op == "=" || op == "==" {
		errCtx := raw.MergeContexts(ts...)
		return nil, nil, errCtx.NewError("Error: invalid statement " +
			"(hint: did you forget return)")
	}
  var assign *js.Assign = nil
  if typeExpr == nil {
    assign = js.NewAssign(lhs, rhs, op, nonLHSTokens[0].Context())
  } else {
    assign = js.NewTypedAssign(lhs, typeExpr, rhs, nonLHSTokens[0].Context())
  }

	return assign, remainingTokens, nil
}

func (p *JSParser) buildImplicitLetStatement(ts_ []raw.Token) (*js.VarStatement,
	[]raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts_, patterns.SEMICOLON)

	if len(ts) < 3 {
		errCtx := ts[1].Context()
		return nil, nil, errCtx.NewError("Error: implicit let statement expects rhs")
	}

	ctx := ts[1].Context()
	letTokens := raw.Concat(raw.NewValueWord("let", ctx),
		ts[0], raw.NewSymbol(patterns.EQUAL, false, ctx), ts[2:])
	varStatement, unexpectedRemaining, err := p.buildVarStatement(letTokens, js.LET)
	if err != nil {
		return nil, nil, err
	}

	if len(unexpectedRemaining) != 0 {
		errCtx := raw.MergeContexts(unexpectedRemaining...)
		return nil, nil, errCtx.NewError("Error: unexpected tokens " +
			"(hint: missing a semicolon?)")
	}

	return varStatement, remainingTokens, nil
}

func (p *JSParser) buildPostIncrOpStatement(ts_ []raw.Token) (*js.PostIncrOp,
	[]raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts_, patterns.SEMICOLON)

	if len(ts) < 2 {
		panic("should've been caught before")
	} else if len(ts) > 2 {
		errCtx := raw.MergeContexts(ts[1], ts[2])
		return nil, nil, errCtx.NewError("Error: unexpected token after postfix ++")
	}

	lhs, err := p.buildExpression(ts[0:1])
	if err != nil {
		return nil, nil, err
	}

	return js.NewPostIncrOp(lhs, ts[1].Context()), remainingTokens, nil
}

func (p *JSParser) buildPostDecrOpStatement(ts_ []raw.Token) (*js.PostDecrOp,
	[]raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts_, patterns.SEMICOLON)

	if len(ts) < 2 {
		panic("should've been caught before")
	} else if len(ts) > 2 {
		errCtx := raw.MergeContexts(ts[1], ts[2])
		return nil, nil, errCtx.NewError("Error: unexpected token after postfix --")
	}

	lhs, err := p.buildExpression(ts[0:1])
	if err != nil {
		return nil, nil, err
	}

	return js.NewPostDecrOp(lhs, ts[1].Context()), remainingTokens, nil
}

func (p *JSParser) buildVoidStatement(ts []raw.Token) (*js.Void,
	[]raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts, patterns.SEMICOLON)

	expr, err := p.buildExpression(ts[1:])
	if err != nil {
		return nil, nil, err
	}

	return js.NewVoidStatement(expr, ts[0].Context()), remainingTokens, nil
}

func (p *JSParser) buildDeleteStatement(ts []raw.Token) (*js.DeleteOp,
	[]raw.Token, error) {
	ts, remainingTokens := p.splitByNextSeparator(ts, patterns.SEMICOLON)

	expr, err := p.buildExpression(ts[1:])
	if err != nil {
		return nil, nil, err
	}

	return js.NewDeleteOp(expr, ts[0].Context()), remainingTokens, nil
}

func (p *JSParser) buildStatement(ts []raw.Token) (js.Statement, []raw.Token, error) {
	ts = p.expandTmpGroups(ts)

	if raw.IsAnyWord(ts[0]) {
		firstWord, err := raw.AssertWord(ts[0])
		if err != nil {
			panic(err)
		}

		switch firstWord.Value() {
		case "const", "let", "var":
			varType, err := js.StringToVarType(firstWord.Value(), firstWord.Context())
			if err != nil {
				panic(err)
			}
			return p.buildVarStatement(ts, varType)
		case "class":
			return p.buildClassStatement(ts)
		case "enum":
			return p.buildEnumStatement(ts)
		case "interface":
			return p.buildInterfaceStatement(ts)
		case "async":
			return p.buildFunctionStatement(ts)
		case "function":
			return p.buildFunctionStatement(ts)
		case "return":
			return p.buildReturnStatement(ts)
		case "throw":
			return p.buildThrowStatement(ts)
		case "break":
			return p.buildBreakStatement(ts)
		case "continue":
			return p.buildContinueStatement(ts)
		case "await":
			return p.buildAwaitStatement(ts)
		case "else":
			errCtx := firstWord.Context()
			return nil, nil, errCtx.NewError("Error: stray else")
		case "if":
			return p.buildIfStatement(ts)
		case "try":
			return p.buildTryCatchStatement(ts)
		case "catch":
			errCtx := firstWord.Context()
			return nil, nil, errCtx.NewError("Error: stray catch")
		case "finally":
			errCtx := firstWord.Context()
			return nil, nil, errCtx.NewError("Error: stray finally")
		case "switch":
			return p.buildSwitchStatement(ts)
		case "while":
			return p.buildWhileStatement(ts)
		case "for":
			return p.buildForStatement(ts)
		case "void":
			return p.buildVoidStatement(ts)
		case "delete":
			return p.buildDeleteStatement(ts) // only a statement, not an expression
		default:
			if len(ts) < 2 {
				errCtx := raw.MergeContexts(ts...)
				return nil, nil, errCtx.NewError("Error: bad statement")
			}

			ilast := p.nextSeparatorPosition(ts, patterns.SEMICOLON)

			switch {
			case raw.IsSymbolThatEndsWith(ts[1], patterns.EQUAL) &&
				!raw.IsSymbol(ts[1], patterns.COLON_EQUAL):
				return p.buildAssignStatement(ts)
			case raw.IsSymbol(ts[1], patterns.COLON_EQUAL):
				return p.buildImplicitLetStatement(ts)
			case raw.IsSymbol(ts[1], patterns.PLUS_PLUS):
				return p.buildPostIncrOpStatement(ts)
			case raw.IsSymbol(ts[1], patterns.MINUS_MINUS):
				return p.buildPostDecrOpStatement(ts)
			case !raw.ContainsSymbol(ts[0:ilast], patterns.EQUAL) &&
				(raw.IsAnyGroup(ts[ilast-2]) || raw.IsAnyWord(ts[ilast-2])) &&
				raw.IsParensGroup(ts[ilast-1]):
				return p.buildCallStatement(ts)
			default:
        // why not use the buildassignstatement function directly?
        return p.buildAssignStatement(ts)

				/*statementTokens, remaining := p.splitByNextSeparator(ts, patterns.SEMICOLON)
				// try to find an equals
				assignStatement_, err := p.buildExpression(statementTokens)
				if err != nil {
					return nil, nil, err
				}

				if assignStatement, ok := assignStatement_.(*js.Assign); ok {
					return assignStatement, remaining, nil
				}

				errCtx := raw.MergeContexts(statementTokens...)
				return nil, nil, errCtx.NewError("Error: invalid statement")*/
			}
		}
	} else {
		ilast := p.nextSeparatorPosition(ts, patterns.SEMICOLON)

		switch {
		case ilast == 0:
			return nil, ts[ilast+1:], nil
		case raw.IsParensGroup(ts[ilast-1]):
			return p.buildCallStatement(ts)
		case ilast == 1 && raw.IsTmpGroup(ts[0]):
			// nested group
			group, err := raw.AssertGroup(ts[0])
			if err != nil {
				return nil, nil, err
			}

			statement, innerRemaining, err := p.buildStatement(group.Fields[0])
			if err != nil {
				return nil, nil, err
			}

			if len(innerRemaining) != 0 {
				errCtx := raw.MergeContexts(innerRemaining...)
				return nil, nil, errCtx.NewError("Error: unexpected statement tokens")
			}

			return statement, ts[ilast:], nil
		//case raw.ContainsSymbol(ts[0:ilast], patterns.EQUAL):
		//return p.buildAssignStatement(ts)
		default:
			// allowed as statement if last part is call
			errCtx := raw.MergeContexts(ts...) //[0].Context()
			return nil, nil, errCtx.NewError("Error: unhandled statement")
		}
	}
}
