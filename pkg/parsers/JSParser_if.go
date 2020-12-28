package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildIfStatement(ts []raw.Token) (*js.If, []raw.Token, error) {
	ifStatement, err := js.NewIf(ts[0].Context())
	if err != nil {
		return nil, nil, err
	}

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
