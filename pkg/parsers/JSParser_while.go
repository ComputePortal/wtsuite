package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildWhileStatement(ts []raw.Token) (*js.While, []raw.Token, error) {
	first := ts[0]

	if len(ts) < 3 {
		errCtx := first.Context()
		return nil, nil, errCtx.NewError("Error: expected while(...){...}")
	}

	parensGroup, err := raw.AssertParensGroup(ts[1])
	if err != nil {
		return nil, nil, err
	}

	if !parensGroup.IsSingle() {
		errCtx := parensGroup.Context()
		return nil, nil, errCtx.NewError("Error: expected single condition")
	}

	cond, err := p.buildExpression(parensGroup.Fields[0])
	if err != nil {
		return nil, nil, err
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[2])
	if err != nil {
		return nil, nil, err
	}

	whileCtx := context.MergeContexts(first.Context(),
		parensGroup.Context(), bracesGroup.Context())

	whileStatement, err := js.NewWhile(cond, whileCtx)
	if err != nil {
		return nil, nil, err
	}

	statements, err := p.buildBlockStatements(bracesGroup)
	if err != nil {
		return nil, nil, err
	}

	for _, st := range statements {
		whileStatement.AddStatement(st)
	}

	ts = p.stripSeparators(3, ts, patterns.SEMICOLON)

	return whileStatement, ts, nil
}
