package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildBlockStatements(bracesGroup *raw.Group) ([]js.Statement, error) {
	return p.buildBlockStatementsInternal(bracesGroup.Fields)
}

func (p *JSParser) buildBlockStatementsInternal(fields [][]raw.Token) ([]js.Statement, error) {
	statements := make([]js.Statement, 0)

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
