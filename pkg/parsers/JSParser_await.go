package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildAwait(ts []raw.Token) (*js.Await, []raw.Token, error) {
	exprTokens, remainingTokens := splitByNextSeparator(ts[1:], patterns.SEMICOLON)

	if len(exprTokens) > 0 {
		expr, err := p.buildExpression(exprTokens)
		if err != nil {
			return nil, nil, err
		}

		await, err := js.NewAwait(expr, ts[0].Context())
		if err != nil {
			return nil, nil, err
		}

		return await, remainingTokens, nil
	} else {
		errCtx := ts[0].Context()
		return nil, nil, errCtx.NewError("Error: expected 1 argument")
	}
}

func (p *JSParser) buildAwaitExpression(ts []raw.Token) (js.Expression, error) {
	await, rem, err := p.buildAwait(ts)
	if err != nil {
		return nil, err
	}

	if len(rem) != 0 {
		errCtx := raw.MergeContexts(rem...)
		return nil, errCtx.NewError("Error: unexpected tokens after await expression (hint: did you forget semicolon?")
	}

	return await, nil
}

func (p *JSParser) buildAwaitStatement(ts []raw.Token) (js.Statement, []raw.Token, error) {
	return p.buildAwait(ts)
}
