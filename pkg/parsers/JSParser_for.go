package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildForInOfStatement(await bool, parensGroup *raw.Group,
	ts []raw.Token, forCtx context.Context) (js.Statement, error) {
	field := parensGroup.Fields[0]
	if len(field) < 4 {
		errCtx := raw.MergeContexts(field...)
		return nil, errCtx.NewError("Error: expected at least 4 tokens")
	}

	varTypeToken, err := raw.AssertWord(field[0])
	if err != nil {
		return nil, err
	}

	varType, err := js.StringToVarType(varTypeToken.Value(), varTypeToken.Context())
	if err != nil {
		return nil, err
	}

	nameToken, err := raw.AssertWord(field[1])
	if err != nil {
		return nil, err
	}

	nameExpr := js.NewVarExpression(nameToken.Value(), nameToken.Context())

	var inOfToken *raw.Word
	if raw.IsAnyWord(field[2]) {
		inOfToken, err = raw.AssertWord(field[2])
		if err != nil {
			panic(err)
		}
	} else if raw.IsSymbol(field[2], "in") {
		inOfToken = raw.NewValueWord("in", field[2].Context())
	} else {
		errCtx := field[2].Context()
		return nil, errCtx.NewError("Error: expected in or of")
	}

	if await {
		if inOfToken.Value() != "of" {
			errCtx := inOfToken.Context()
			return nil, errCtx.NewError("Error: expected of (await for loop)")
		}
	}

	rhs, err := p.buildExpression(field[3:])
	if err != nil {
		return nil, err
	}

	switch inOfToken.Value() {
	case "in":
		return js.NewForIn(varType, nameExpr, rhs, forCtx)
	case "of":
		return js.NewForOf(await, varType, nameExpr, rhs, forCtx)
	default:
		errCtx := inOfToken.Context()
		return nil, errCtx.NewError("Error: expected 'in' or 'of'")
	}
}

func (p *JSParser) buildForRegularStatement(parensGroup *raw.Group,
	ts []raw.Token, forCtx context.Context) (js.Statement, error) {
	initField := parensGroup.Fields[0]
	initExprs := make([]js.Expression, 0)

	implicitVarType := true
	varType := js.LET

	// add init expressions
	if len(initField) > 0 { // init can be empty
		if len(initField) < 3 {
			errCtx := context.MergeContexts(raw.MergeContexts(initField...))
			return nil, errCtx.NewError("Error: bad for loop init field")
		}

		if !raw.IsSymbol(initField[1], patterns.COLON_EQUAL) && raw.IsAnyWord(initField[0]) {
			varTypeToken, err := raw.AssertWord(initField[0])
			if err != nil {
				return nil, err
			}

			varType, err = js.StringToVarType(varTypeToken.Value(), varTypeToken.Context())
			if err != nil {
				return nil, err
			}

			implicitVarType = false

			initField = initField[1:]
		}

		components := p.splitBySeparator(initField, patterns.COMMA)
		for _, c := range components {
			if len(c) < 2 {
				return nil, forCtx.NewError("Error: unexpectedly short init expression")
			}

			isColonEqual := false
			if raw.IsSymbol(c[1], patterns.COLON_EQUAL) {
				if !implicitVarType {
					errCtx := c[1].Context()
					return nil, errCtx.NewError("Error: cannot combine := and var type keyword")
				}

				raw.ChangeSymbol(c[1], patterns.EQUAL, false)
				isColonEqual = true // now all assign expressions must be ':='
			}

			initExpr, err := p.buildExpression(c)
			if err != nil {
				return nil, err
			}

			initAssign, initAssignOk := initExpr.(*js.Assign)
			if initAssignOk && initAssign.HasLhsVarExpression() && implicitVarType && !isColonEqual {
				errCtx := initAssign.Context()
				return nil, errCtx.NewError("Error: use := for every var assign, or for none")
			}

			initExprs = append(initExprs, initExpr)
		}
	}

	// add condition
	condField := parensGroup.Fields[1]
	var cond js.Expression = nil
	if len(condField) > 0 {
		var err error
		cond, err = p.buildExpression(condField)
		if err != nil {
			return nil, err
		}
	}

	// add incr expressions
	incrField := parensGroup.Fields[2]
	incrExprs := make([]js.Expression, 0)
	if len(incrField) > 0 {
		if len(incrField) == 1 {
			errCtx := incrField[0].Context()
			return nil, errCtx.NewError("Error: bad for loop final statement, expected at least two tokens (eg. a++)")
		}

		components := p.splitBySeparator(incrField, patterns.COMMA)
		for _, c := range components {
			incrExpr, err := p.buildExpression(c)
			if err != nil {
				return nil, err
			}

			incrExprs = append(incrExprs, incrExpr)
		}
	}

	return js.NewFor(varType, initExprs, cond, incrExprs, forCtx)
}

func (p *JSParser) buildForStatement(ts []raw.Token) (js.Statement, []raw.Token, error) {
	await := false

	first := ts[0]

	if len(ts) < 3 {
		errCtx := first.Context()
		return nil, nil, errCtx.NewError("Error: invalid for statement")
	}

	if raw.IsWord(ts[1], "await") {
		await = true
		ts = ts[1:]
	}

	if len(ts) < 3 {
		errCtx := first.Context()
		return nil, nil, errCtx.NewError("Error: invalid for statement")
	}

	parensGroup, err := raw.AssertParensGroup(ts[1])
	if err != nil {
		return nil, nil, err
	}

	if parensGroup.IsEmpty() {
		errCtx := parensGroup.Context()
		return nil, nil, errCtx.NewError("Error: empty for loop parentheses")
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[2])
	if err != nil {
		return nil, nil, err
	}

	// handle the await
	if await && len(parensGroup.Fields) != 1 {
		errCtx := parensGroup.Context()
		return nil, nil, errCtx.NewError("Error: bad await for loop")
	}

	var forStatement js.Statement
	forCtx := context.MergeContexts(first.Context(),
		parensGroup.Context(), bracesGroup.Context())

	switch {
	case len(parensGroup.Fields) == 1:
		// (const v of|in ...)
		forStatement, err = p.buildForInOfStatement(await, parensGroup, ts, forCtx)
		if err != nil {
			return nil, nil, err
		}
	case len(parensGroup.Fields) == 3 && parensGroup.IsSemiColon():
		forStatement, err = p.buildForRegularStatement(parensGroup, ts, forCtx)
		if err != nil {
			return nil, nil, err
		}
	}

	statements, err := p.buildBlockStatements(bracesGroup)
	if err != nil {
		return nil, nil, err
	}

	for _, st := range statements {
		forStatement.AddStatement(st)
	}

	remaining := p.stripSeparators(3, ts, patterns.SEMICOLON)

	return forStatement, remaining, nil
}
