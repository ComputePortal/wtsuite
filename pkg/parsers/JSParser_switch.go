package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildSwitchStatement(ts []raw.Token) (js.Statement, []raw.Token, error) {
	if len(ts) < 2 {
		errCtx := raw.MergeContexts(ts...)
		return nil, nil, errCtx.NewError("Error: bad switch statement (expected switch(...){...} or switch{...})")
	}

	switchContext := ts[0].Context()
	ifTypeSwitch := false
	var switchStatement *js.Switch = nil

	if raw.IsParensGroup(ts[1]) {
		if len(ts) < 3 {
			errCtx := raw.MergeContexts(ts...)
			return nil, nil, errCtx.NewError("Error: bad switch statement (expected switch(...){...})")
		}

		exprGroup, err := raw.AssertParensGroup(ts[1])
		if err != nil {
			panic(err)
		}

		exprField, err := exprGroup.FlattenCommas()
		if err != nil {
			return nil, nil, err
		}

		expr, err := p.buildExpression(exprField)
		if err != nil {
			return nil, nil, err
		}

		switchStatement, err = js.NewSwitch(expr, switchContext)
		if err != nil {
			return nil, nil, err
		}

		ts = ts[2:]
	} else {
		ts = ts[1:]

		ifTypeSwitch = true
		var err error
		switchStatement, err = js.NewSwitch(nil, switchContext)
		if err != nil {
			return nil, nil, err
		}
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[0])
	if err != nil {
		return nil, nil, err
	}

	remaining := ts[1:]

	for i, field := range bracesGroup.Fields {
		if len(field) >= 1 && (raw.IsWord(field[0], "case") || raw.IsWord(field[0], "default")) {
			if raw.IsWord(field[0], "case") {
				if len(field) < 3 {
					errCtx := raw.MergeContexts(field...)
					return nil, nil, errCtx.NewError("Error: bad switch case statement")
				}
			} else {
				if len(field) < 2 {
					errCtx := field[0].Context()
					return nil, nil, errCtx.NewError("Error: bad switch default statement")
				} else if !raw.IsSymbol(field[1], patterns.COLON) {
					errCtx := field[0].Context()
					return nil, nil, errCtx.NewError("Error: bad switch default statement, missing colon")
				}
			}

			clauseTokens, otherTokens := raw.SplitByFirstSymbol(field, patterns.COLON)
			if len(otherTokens) == 0 && !raw.ContainsSymbol(field, patterns.COLON) {
        errCtx := raw.MergeContexts(field...)
        return nil, nil, errCtx.NewError("Error: \":\" not found")
			}

			clauseHeaders := [][]raw.Token{clauseTokens}

			// eat the empty clauses that might be in otherTokens
			for len(otherTokens) > 0 && (raw.IsWord(otherTokens[0], "case") || raw.IsWord(otherTokens[0], "default")) {
        clauseTokens_, otherTokens_ := raw.SplitByFirstSymbol(otherTokens, patterns.COLON)
        if len(otherTokens_) == 0 && !raw.ContainsSymbol(otherTokens, patterns.COLON) {
          errCtx := raw.MergeContexts(otherTokens...)
          return nil, nil, errCtx.NewError("Error: \":\" not found")
				} else {
					clauseHeaders = append(clauseHeaders, clauseTokens_)
					otherTokens = otherTokens_
				}
			}

			// find next cast/default, and append everything before that to otherTokens
			clauseBody := [][]raw.Token{otherTokens}

			for iNext := i + 1; iNext < len(bracesGroup.Fields); iNext++ {
				nextField := bracesGroup.Fields[iNext]

				if len(nextField) > 0 && (raw.IsWord(nextField[0], "case") || raw.IsWord(nextField[0], "default")) {
					break
				}

				clauseBody = append(clauseBody, nextField)
			}

			for i, clauseHeader := range clauseHeaders {
				if i < len(clauseHeaders)-1 {
					if err := p.buildSwitchClause(switchStatement, clauseHeader, [][]raw.Token{}); err != nil {
						return nil, nil, err
					}
				} else {
					if err := p.buildSwitchClause(switchStatement, clauseHeader, clauseBody); err != nil {
						return nil, nil, err
					}
				}
			}
		}
	}

	if ifTypeSwitch {
		// convert to actual if-else statement
		return switchStatement.ConvertToIf(), remaining, nil
	} else {
		return switchStatement, remaining, nil
	}
}

func (p *JSParser) buildSwitchClause(switchStatement *js.Switch, clauseTokens []raw.Token,
	otherTokens [][]raw.Token) error {
	// cut off the colon
	clauseTokens = clauseTokens[:len(clauseTokens)-1]

	if raw.IsWord(clauseTokens[0], "case") {
		clauseExpr, err := p.buildExpression(clauseTokens[1:])
		if err != nil {
			return err
		}

		if err := switchStatement.AddCase(clauseExpr); err != nil {
			return err
		}
	} else if raw.IsWord(clauseTokens[0], "default") {
		if len(clauseTokens) != 1 {
			errCtx := raw.MergeContexts(clauseTokens...)
			return errCtx.NewError("Error: bad switch default clause")
		}

		if err := switchStatement.AddDefault(clauseTokens[0].Context()); err != nil {
			return err
		}
	} else {
		errCtx := raw.MergeContexts(clauseTokens...)
		return errCtx.NewError("Error: bad switch clause")
	}

	statements, err := p.buildBlockStatementsInternal(otherTokens)
	if err != nil {
		return err
	}

	for _, statement := range statements {
		switchStatement.AddStatement(statement)
	}

	return nil
}
