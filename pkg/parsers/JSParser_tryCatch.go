package parsers

import (
	"../tokens/js"
	"../tokens/raw"
)

func (p *JSParser) buildTryCatchStatement(ts []raw.Token) (*js.TryCatch, []raw.Token, error) {
	tryCatch, err := js.NewTryCatch(ts[0].Context())
	if err != nil {
		return nil, nil, err
	}

	firstTry := true
	for len(ts) > 0 && (raw.IsWord(ts[0], "try") || raw.IsWord(ts[0], "catch") || raw.IsWord(ts[0], "finally")) {
		switch {
		case raw.IsWord(ts[0], "try"):
			if !firstTry {
				errCtx := ts[0].Context()
				return nil, nil, errCtx.NewError("Error: duplicate try statement")
			}

			if len(ts) < 4 {
				errCtx := raw.MergeContexts(ts...)
				return nil, nil, errCtx.NewError("Error: expected try {...} catch | finally {...}")
			}

			bracesGroup, err := raw.AssertBracesGroup(ts[1])
			if err != nil {
				return nil, nil, err
			}

			statements, err := p.buildBlockStatements(bracesGroup)
			if err != nil {
				return nil, nil, err
			}

			for _, st := range statements {
				tryCatch.AddStatement(st)
			}

			ts = ts[2:]

			firstTry = false
		case raw.IsWord(ts[0], "catch"):
			errCtx := ts[0].Context()
			if len(ts) < 2 {
				return nil, nil, errCtx.NewError("Error: expected catch [(e)] {...}")
			}

			var catchArg *js.FunctionArgument = nil
			if raw.IsParensGroup(ts[1]) {
				condGroup, err := raw.AssertParensGroup(ts[1])
				if err != nil {
					panic(err)
				}

				condField, err := condGroup.FlattenCommas()
				if err != nil {
					return nil, nil, err
				}

				if len(condField) == 0 {
					return nil, nil, errCtx.NewError("Error: no catch argument specified (hint: leave out parens)")
				}

				catchArg, err = p.buildFunctionArgument(condField, false)
				if err != nil {
					return nil, nil, err
				}

				ts = ts[2:]
			} else {
				ts = ts[1:]
			}

			if err := tryCatch.AddCatch(catchArg); err != nil {
				return nil, nil, err
			}

			if len(ts) < 1 {
				return nil, nil, errCtx.NewError("Error: expected catch [(e)] {...}")
			}

			bracesGroup, err := raw.AssertBracesGroup(ts[0])
			if err != nil {
				return nil, nil, err
			}

			statements, err := p.buildBlockStatements(bracesGroup)
			if err != nil {
				return nil, nil, err
			}

			for _, st := range statements {
				tryCatch.AddStatement(st)
			}

			ts = ts[1:]
		case raw.IsWord(ts[0], "finally"):
			ctx := ts[0].Context()
			if len(ts) < 2 {
				errCtx := ctx
				return nil, nil, errCtx.NewError("Error: expected finally {...}")
			}

			bracesGroup, err := raw.AssertBracesGroup(ts[1])
			if err != nil {
				return nil, nil, err
			}

			statements, err := p.buildBlockStatements(bracesGroup)
			if err != nil {
				return nil, nil, err
			}

			if err := tryCatch.AddFinally(); err != nil {
				return nil, nil, err
			}

			for _, st := range statements {
				tryCatch.AddStatement(st)
			}

			ts = ts[2:]
		default:
			errCtx := ts[0].Context()
			return nil, nil, errCtx.NewError("Error: unexpected word")
		}
	}

	return tryCatch, ts, nil
}
