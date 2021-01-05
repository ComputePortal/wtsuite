package parsers

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *GLSLParser) buildReturnStatement(ts []raw.Token) (*glsl.Return, []raw.Token, error) {
	exprTokens, remainingTokens := splitByNextSeparator(ts[1:], patterns.SEMICOLON)
	var expr glsl.Expression = nil
	if len(exprTokens) > 0 {
		var err error
		expr, err = p.buildExpression(exprTokens)
		if err != nil {
			return nil, nil, err
		}
	}

	return glsl.NewReturn(expr, ts[0].Context()), remainingTokens, nil
}

func (p *GLSLParser) buildIfStatement(ts []raw.Token) (*glsl.If, []raw.Token, error) {
  ifStatement := glsl.NewIf(ts[0].Context())

  for raw.IsWord(ts[0], "if") {
		if len(ts) < 3 {
			errCtx := raw.MergeContexts(ts...)
			return nil, nil, errCtx.NewError("Error: expected 'if(...){...}'")
		}

		condGroup, err := raw.AssertParensGroup(ts[1])
		if err != nil {
			return nil, nil, err
		}

		condField, err := condGroup.FlattenCommas()
		if err != nil {
			return nil, nil, err
		}

		cond, err := p.buildExpression(condField)
		if err != nil {
			return nil, nil, err
		}

		ifStatement.AddCondition(cond)

		bracesGroup, err := raw.AssertBracesGroup(ts[2])
		if err != nil {
			return nil, nil, err
		}

		statements, err := p.buildBlockStatements(bracesGroup)
		if err != nil {
			return nil, nil, err
		}

		for _, st := range statements {
			ifStatement.AddStatement(st)
		}

		if len(ts) >= 4 && raw.IsWord(ts[3], "else") {
			if len(ts) < 5 {
				errCtx := ts[3].Context()
				return nil, nil, errCtx.NewError("Error: bad else statement")
			}

			if raw.IsWord(ts[4], "if") {
				ts = ts[4:]
				continue
			}

			finalGroup, err := raw.AssertBracesGroup(ts[4])
			if err != nil {
				return nil, nil, err
			}

			ifStatement.AddElse()

			statements, err := p.buildBlockStatements(finalGroup)
			if err != nil {
				return nil, nil, err
			}

			for _, st := range statements {
				ifStatement.AddStatement(st)
			}

			ts = ts[5:]
			break
		} else {
			ts = ts[3:]
			break
		}
	}

	ts = stripSeparators(0, ts, patterns.SEMICOLON)

	return ifStatement, ts, nil
}

func (p *GLSLParser) buildVarStatement(ts []raw.Token) (*glsl.VarStatement, []raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  if len(ts) < 2 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected at least 2 tokens")
  }

  iequal := nextSeparatorPosition(ts, patterns.EQUAL)

  if iequal < 2 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: bad var statement")
  }

  iname := iequal - 1
  arraySize := -1
  if raw.IsBracketsGroup(ts[iname]) {
    var err error
    arraySize, err = p.buildArraySize(ts[iname])
    if err != nil {
      return nil, nil, err
    }
    iname = iname - 1
  }

  nameToken, err := raw.AssertWord(ts[iname])
  if err != nil {
    return nil, nil, err
  }

  typeExpr, err := p.buildTypeExpression(ts[0:iname])
  if err != nil {
    return nil, nil, err
  }

  var rhsExpr glsl.Expression = nil
  if iequal < len(ts) - 1 {
    rhsExpr, err = p.buildExpression(ts[iequal+1:])
    if err != nil {
      return nil, nil, err
    }
  }

  return glsl.NewVarStatement(typeExpr, nameToken.Value(), arraySize, rhsExpr, nameToken.Context()), remainingTokens, nil
}

func (p *GLSLParser) buildAssignStatement(ts []raw.Token) (*glsl.Assign, []raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  if len(ts) < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected at least 3 tokens")
  } 

  iequal := nextSymbolPositionThatEndsWith(ts, patterns.EQUAL)

  if iequal == 0 || iequal == len(ts) - 1{
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: bad assign statement")
  }

  a, err := p.buildExpression(ts[0:iequal])
  if err != nil {
    return nil, nil, err
  }

  b, err := p.buildExpression(ts[iequal+1:])
  if err != nil {
    return nil, nil, err
  }

  symbol, err := raw.AssertAnySymbol(ts[iequal])
  if err != nil {
    panic(err)
  }

  op := strings.TrimSuffix(symbol.Value(), patterns.EQUAL)

  return glsl.NewAssign(a, b, op, raw.MergeContexts(ts...)), remainingTokens, nil
}

func (p *GLSLParser) buildIncrStatement(ts []raw.Token) (*glsl.PostIncrOp, []raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  if len(ts) != 2 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected 2 tokens")
  }

  if !raw.IsAnyWord(ts[0]) {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected single word (can't increment package exports)")
  }

  lhs, err := p.buildExpression(ts[0:1])
  if err != nil {
    return nil, nil, err
  }

  if !raw.IsSymbol(ts[1], patterns.PLUS_PLUS) {
    errCtx := ts[1].Context()
    return nil, nil, errCtx.NewError("Error: expected ++")
  }

  return glsl.NewPostIncrOp(lhs, raw.MergeContexts(ts...)), remainingTokens, nil
}

func (p *GLSLParser) buildDecrStatement(ts []raw.Token) (*glsl.PostDecrOp, []raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)
  if len(ts) != 2 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected 2 tokens")
  }

  if !raw.IsAnyWord(ts[0]) {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected single word (can't increment package exports)")
  }

  lhs, err := p.buildExpression(ts[0:1])
  if err != nil {
    return nil, nil, err
  }

  if !raw.IsSymbol(ts[1], patterns.MINUS_MINUS) {
    errCtx := ts[1].Context()
    return nil, nil, errCtx.NewError("Error: expected --")
  }

  return glsl.NewPostDecrOp(lhs, raw.MergeContexts(ts...)), remainingTokens, nil
}

func (p *GLSLParser) buildCallStatement(ts []raw.Token) (glsl.Statement, []raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  call, err := p.buildCall(ts)
  if err != nil {
    return nil, nil, err
  }

  if glsl.IsMacroStatement(call) {
    macroStatement, err := glsl.DispatchMacroStatement(call)
    if err != nil {
      return nil, nil, err
    }

    return macroStatement, remainingTokens, nil
  } else if glsl.IsMacroExpression(call) {
    errCtx := call.Context()
    return nil, nil, errCtx.NewError("Error: is an expression, not a statement")
  } 

  return call, remainingTokens, nil
}

func (p *GLSLParser) buildStatement(ts []raw.Token) (glsl.Statement, []raw.Token, error) {
  ts = p.expandTmpGroups(ts)

  if raw.IsAnyWord(ts[0]) {
    firstWord, err := raw.AssertWord(ts[0])
    if err != nil {
      return nil, nil, err
    }

    switch firstWord.Value() {
    case "return":
      return p.buildReturnStatement(ts)
		case "else":
			errCtx := firstWord.Context()
			return nil, nil, errCtx.NewError("Error: stray else")
    case "if":
      return p.buildIfStatement(ts)
    default:
      ilast := nextSeparatorPosition(ts, patterns.SEMICOLON)

      if raw.ContainsSymbol(ts[0:ilast], patterns.EQUAL) {
        iequal := nextSeparatorPosition(ts[0:ilast], patterns.EQUAL)

        if (iequal >= 2 && raw.IsAnyWord(ts[iequal-1]) && raw.IsAnyWord(ts[iequal-2])) || (iequal >= 3 && raw.IsAnyWord(ts[iequal-2]) && raw.IsAnyWord(ts[iequal-3]) && raw.IsBracketsGroup(ts[iequal-1])) {
          return p.buildVarStatement(ts)
        } else {
          return p.buildAssignStatement(ts)
        }
      } else if raw.ContainsSymbol(ts[0:ilast], patterns.PLUS_EQUAL) || raw.ContainsSymbol(ts[0:ilast], patterns.MUL_EQUAL) || raw.ContainsSymbol(ts[0:ilast], patterns.MINUS_EQUAL) {
        return p.buildAssignStatement(ts)
      } else if ilast >= 2 && raw.IsSymbol(ts[1], patterns.PLUS_PLUS) {
        return p.buildIncrStatement(ts)
      } else if ilast >= 2 && raw.IsSymbol(ts[1], patterns.MINUS_MINUS) {
        return p.buildDecrStatement(ts)
      } else if ilast >= 2 && raw.IsParensGroup(ts[ilast-1]) && raw.IsAnyWord(ts[ilast-2]) {
        return p.buildCallStatement(ts)
      } else if (ilast >= 2 && raw.IsAnyWord(ts[ilast-1]) && raw.IsAnyWord(ts[ilast-2])) || (ilast >= 3 && raw.IsAnyWord(ts[ilast-2]) && raw.IsAnyWord(ts[ilast-3]) && raw.IsBracketsGroup(ts[ilast-1])) {
        return p.buildVarStatement(ts)
      } else {
        errCtx := raw.MergeContexts(ts[0:ilast]...)
        return nil, nil, errCtx.NewError("Error: bad statement")
      }
    }
  } else {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: bad statement")
  }
}

func (p *GLSLParser) buildBlockStatements(bracesGroup *raw.Group) ([]glsl.Statement, error) {
	return p.buildBlockStatementsInternal(bracesGroup.Fields)
}

func (p *GLSLParser) buildBlockStatementsInternal(fields [][]raw.Token) ([]glsl.Statement, error) {
	statements := make([]glsl.Statement, 0)

	for _, field := range fields {
		if len(field) == 0 {
			continue
		}

		statement, remaining, err := p.buildStatement(field)
		if err != nil {
			return nil, err
		}

		statements = append(statements, statement)

		for len(remaining) > 0 {
			statement, remaining, err = p.buildStatement(remaining)
			if err != nil {
				return nil, err
			}

			statements = append(statements, statement)
		}

		if len(remaining) > 0 {
			errCtx := remaining[0].Context()
			return nil, errCtx.NewError("Error: unexpected remaining tokens")
		}
	}

	return statements, nil
}

func (p *GLSLParser) buildStructEntry(ts []raw.Token) (*glsl.StructEntry, error) {
  n := len(ts)

  ctx := raw.MergeContexts(ts...)

  if n < 2 {
    errCtx := ctx
    return nil, errCtx.NewError("Error: expected at least two tokens")
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

  nameExpr, err := p.buildVarExpression(ts[n-1])
  if err != nil {
    return nil, err
  }

  typeExpr, err := p.buildTypeExpression(ts[0:n-1])
  if err != nil {
    return nil, err
  }

  return glsl.NewStructEntry(typeExpr, nameExpr, arraySize, ctx), nil
}

func (p *GLSLParser) buildStruct(ts []raw.Token, isExport bool) ([]raw.Token, error) {
  ts, remainingTokens := splitByNextSeparator(ts, patterns.SEMICOLON)

  if len(ts) != 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, errCtx.NewError("Error: expected 3 tokens")
  }

  name, err := p.buildVarExpression(ts[1])
  if err != nil {
    return nil, err
  }

  entries := make([]*glsl.StructEntry, 0)

  brace, err := raw.AssertBracesGroup(ts[2])
  if err != nil {
    return nil, err
  }

  if brace.IsComma() {
    errCtx := brace.Context()
    return nil, errCtx.NewError("Error: expected semicolon separators")
  }

  for _, field := range brace.Fields {
    structEntry, err := p.buildStructEntry(field)
    if err != nil {
      return nil, err
    }

    entries = append(entries, structEntry)
  }

  ctx := raw.MergeContexts(ts...)

  st, err := glsl.NewStruct(name, entries, ctx)
  if err != nil {
    return nil, err
  }

  p.module.AddStatement(st)

  if isExport {
    if err := p.buildSimpleExport(name.Name(), name.Context()); err != nil {
      return nil, err
    }
  }

  return remainingTokens, nil
}
