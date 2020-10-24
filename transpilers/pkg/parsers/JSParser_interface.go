package parsers

import (
	"../tokens/js"
	"../tokens/js/prototypes"
	"../tokens/patterns"
	"../tokens/raw"
)

func (p *JSParser) buildInterface(ts []raw.Token) (*js.ClassInterface, error) {
	interfCtx := raw.MergeContexts(ts...)

	if len(ts) == 2 && raw.IsBracesGroup(ts[1]) {
		errCtx := interfCtx
		return nil, errCtx.NewError("Error: missing interface name")
	} else if len(ts) < 3 {
		errCtx := interfCtx
		return nil, errCtx.NewError("Error: bad interface definition")
	}

	clType, ts, err := p.buildClassOrExtendsTypeExpression(ts[1:])
	if err != nil {
		return nil, err
	}

	var extends *js.TypeExpression = nil
	if raw.IsWord(ts[0], "extends") {
		extends, ts, err = p.buildClassOrExtendsTypeExpression(ts[1:])
		if err != nil {
			return nil, err
		}
	}

	classInterface, err := js.NewClassInterface(clType, extends, interfCtx)
	if err != nil {
		return nil, err
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[len(ts)-1])
	if err != nil {
		return nil, err
	}

	if bracesGroup.IsComma() {
		errCtx := bracesGroup.Context()
		return nil, errCtx.NewError("Error: interface uses semicolon separator")
	}

	for _, field := range bracesGroup.Fields {
		if len(field) == 0 {
			continue
		}

		fi, remaining, err := p.buildFunctionInterface(field, true, false, interfCtx)
		if err != nil {
			return nil, err
		}

		if fi.Role() != prototypes.NORMAL &&
			fi.Role() != prototypes.GETTER &&
			fi.Role() != prototypes.SETTER &&
			fi.Role() != prototypes.ASYNC {
			errCtx := fi.Context()
			return nil, errCtx.NewError("Error: illegal interface function role(s)")
		}

		if len(remaining) != 0 {
			errCtx := raw.MergeContexts(remaining...)
			return nil, errCtx.NewError("Error: unexpected tokens (hint: did forget a semicolon?)")
		}

		if err := classInterface.AddMember(fi); err != nil {
			return nil, err
		}
	}

	return classInterface, nil
}

func (p *JSParser) buildInterfaceStatement(ts []raw.Token) (*js.ClassInterface, []raw.Token, error) {
	for i, t := range ts {
		if raw.IsBracesGroup(t) {
			statement, err := p.buildInterface(ts[0 : i+1])
			if err != nil {
				return nil, nil, err
			}

			remaining := p.stripSeparators(i+1, ts, patterns.SEMICOLON)

			return statement, remaining, nil
		}
	}

	errCtx := raw.MergeContexts(ts...)
	return nil, nil, errCtx.NewError("Error: no interface body found")
}
