package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildEnum(ts []raw.Token) (*js.Enum, error) {
	enCtx := raw.MergeContexts(ts...)

	if len(ts) < 3 {
		errCtx := enCtx
		return nil, errCtx.NewError("Error: bad enum definition")
	}

	clType, ts, err := p.buildClassOrExtendsTypeExpression(ts[1:])
	if err != nil {
		return nil, err
	}

	// extends is inferred (should be Int or String)
	// XXX: if no good use is found for extends differing from Int or String, then the optional extends code should be completely removed
	//var extends *js.TypeExpression = nil
	/*if raw.IsWord(ts[0], "extends") {
		extends, ts, err = p.buildClassOrExtendsTypeExpression(ts[1:])
		if err != nil {
			return nil, err
		}
	}*/

	keys := make([]*js.Word, 0)
	uniqueKeys := make(map[string]*js.Word)
	values := make([]js.Expression, 0) // fill with to apply autocomplete later

	// asserts that all keys are unique
	appendKey := func(s string, c context.Context) error {
		if old, ok := uniqueKeys[s]; ok {
			errCtx := c
			err := errCtx.NewError("Error: key not unique")
			err.AppendContextString("Info: also defined here", old.Context())
			return err
		}

		w := js.NewWord(s, c)
		keys = append(keys, w)
		uniqueKeys[s] = w
		return nil
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[len(ts)-1])
	if err != nil {
		return nil, err
	}

	if bracesGroup.IsSemiColon() {
		errCtx := bracesGroup.Context()
		return nil, errCtx.NewError("Error: enums use comma separator")
	}

	// operators aren't yet nested, so we can detect '=' as separator
	for _, field := range bracesGroup.Fields {
		if len(field) == 0 {
			continue
		}

		switch {
		case len(field) == 1:
			key, err := raw.AssertWord(field[0])
			if err != nil {
				return nil, err
			}

			if err := appendKey(key.Value(), key.Context()); err != nil {
				return nil, err
			}
			values = append(values, nil)
		case len(field) > 2:
			key, err := raw.AssertWord(field[0])
			if err != nil {
				return nil, err
			}

			if !raw.IsSymbol(field[1], patterns.EQUAL) {
				errCtx := field[1].Context()
				return nil, errCtx.NewError("Error: expected equals or comma")
			}

			rhs, err := p.buildExpression(field[2:])
			if err != nil {
				return nil, err
			}

			if err := appendKey(key.Value(), key.Context()); err != nil {
				return nil, err
			}
			values = append(values, rhs)
		default:
			errCtx := raw.MergeContexts(field...)
			return nil, errCtx.NewError("Error: bad enum member")
		}
	}

	extendsName := ""
	// infer the type first, then complete the values
	for _, val := range values {
		if val == nil {
			continue
		}

		switch val.(type) {
		case *js.LiteralInt:
			if extendsName == "" {
				extendsName = "Int"
			} else if extendsName != "Int" {
				errCtx := val.Context()
				return nil, errCtx.NewError("Error: can't mix enum types")
			}
		case *js.LiteralString:
			if extendsName == "" {
				extendsName = "String"
			} else if extendsName != "String" {
				errCtx := val.Context()
				return nil, errCtx.NewError("Error: can't mix enum types")
			}
		default:
			errCtx := val.Context()
			return nil, errCtx.NewError("Error: expected literal string or literal int")
		}
	}

	// if no values are set we default to Int
	if extendsName == "" {
		extendsName = "Int"
	}

	// now complete the values
	for i, val := range values {
		if val != nil {
			continue
		}

		switch extendsName {
		case "Int":
			values[i] = js.NewLiteralInt(i, keys[i].Context())
		case "String":
			values[i] = js.NewLiteralString(keys[i].Value(), keys[i].Context())
		default:
			panic("not yet implemented")
		}
	}

	extends, err := js.NewTypeExpression(extendsName, nil, nil, enCtx)
  if err != nil {
    return nil, err
  }

	return js.NewEnum(clType, extends, keys, values, enCtx)
}

func (p *JSParser) buildEnumStatement(ts []raw.Token) (*js.Enum, []raw.Token, error) {
	for i, t := range ts {
		if raw.IsBracesGroup(t) {
			statement, err := p.buildEnum(ts[0 : i+1])
			if err != nil {
				return nil, nil, err
			}

			remaining := stripSeparators(i+1, ts, patterns.SEMICOLON)

			return statement, remaining, nil
		}
	}

	errCtx := raw.MergeContexts(ts...)
	return nil, nil, errCtx.NewError("Error: no enum body found")
}
