package parsers

import (
	"fmt"
	"strings"

	"../tokens/js"
	"../tokens/patterns"
	"../tokens/raw"
)

func (p *JSParser) buildLiteralStringExpression(t raw.Token) (js.Expression, error) {
	s, err := raw.AssertLiteralString(t)
	if err != nil {
		panic(err)
	}

	return js.NewLiteralString(s.Value(), s.Context()), nil
}

func (p *JSParser) buildTemplateStringExpression(t raw.Token) (js.Expression, error) {
	s, err := raw.AssertGroup(t)
	if err != nil {
		panic(err)
	}

	exprs := make([]js.Expression, len(s.Fields))

	for i, ts := range s.Fields {
		expr, err := p.buildExpression(ts)
		if err != nil {
			return nil, err
		}

		exprs[i] = expr
	}
	return js.NewTemplateString(exprs, s.Context())
}

func (p *JSParser) buildLiteralIntExpression(t raw.Token) (js.Expression, error) {
	i, err := raw.AssertLiteralInt(t)
	if err != nil {
		panic(err)
	}

	return js.NewLiteralInt(i.Value(), i.Context()), nil
}

func (p *JSParser) buildLiteralFloatExpression(t raw.Token) (js.Expression, error) {
	r, err := raw.AssertLiteralFloat(t, "")
	if err != nil {
		panic(err)
	}

	return js.NewLiteralFloat(r.Value(), r.Context()), nil
}

// eg. NaN and Infinity
func (p *JSParser) buildSpecialNumberExpression(t raw.Token) (js.Expression, error) {
	r, err := raw.AssertSpecialNumber(t)
	if err != nil {
		panic(err)
	}

	return js.NewSpecialNumber(r.Value(), r.Context()), nil
}

func (p *JSParser) buildLiteralBoolExpression(t raw.Token) (js.Expression, error) {
	b, err := raw.AssertLiteralBool(t)
	if err != nil {
		panic(err)
	}
	return js.NewLiteralBoolean(b.Value(), b.Context()), nil
}

func (p *JSParser) buildLiteralArrayExpression(t raw.Token) (js.Expression, error) {
	group, err := raw.AssertBracketsGroup(t)
	if err != nil {
		panic(err)
	}

	if group.IsSemiColon() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: literal Array can't use semicolons as separators")
	}

	items := make([]js.Expression, 0)
	for _, field := range group.Fields {
		item, err := p.buildExpression(field)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return js.NewLiteralArray(items, group.Context()), nil
}

func (p *JSParser) buildLiteralObjectExpression(t raw.Token) (js.Expression, error) {
	group, err := raw.AssertBracesGroup(t)
	if err != nil {
		panic(err)
	}

	if group.IsSemiColon() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: literal Object can't use semicolons as separators")
	}

	// keys are just temporary
	keys := make([]*js.Word, 0)
	values := make([]js.Expression, 0)

	for _, field := range group.Fields {
		components := p.splitBySeparator(field, patterns.COLON)
		if len(components) != 2 {
			errCtx := raw.MergeContexts(field...)
			return nil, errCtx.NewError("Error: bad dict key-value entry")
		}

		key_, err := p.buildExpression(components[0])
		if err != nil {
			return nil, err
		}

		switch key := key_.(type) {
		case *js.VarExpression:
			keys = append(keys, js.NewWord(key.Name(), key.Context()))
		case *js.LiteralString:
			keys = append(keys, js.NewWord(key.Value(), key.Context()))
		default:
			panic("unexpected")
		}

		value, err := p.buildExpression(components[1])
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}

	obj, err := js.NewLiteralObject(keys, values, group.Context())
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *JSParser) buildParensExpression(t raw.Token) (js.Expression, error) {
	group, err := raw.AssertParensGroup(t)
	if err != nil {
		panic(err)
	}

	if !group.IsSingle() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: simple parentheses must have exactly one entry")
	}

	expr, err := p.buildExpression(group.Fields[0])
	if err != nil {
		return nil, err
	}

	return js.NewParens(expr, group.Context()), nil
}

func (p *JSParser) buildTernaryOpExpression(t raw.Token) (js.Expression, error) {
	op, err := raw.AssertAnyTernaryOperator(t)
	if err != nil {
		panic(err)
	}

	jsName, ok := p.translateOpName(op.Name())
	if !ok {
		errCtx := op.Context()
		return nil, errCtx.NewError("Error: this ternary operator is not yet handled")
	}

	a, err := p.buildExpression(op.Args()[0:1])
	if err != nil {
		return nil, err
	}

	b, err := p.buildExpression(op.Args()[1:2])
	if err != nil {
		return nil, err
	}

	c, err := p.buildExpression(op.Args()[2:3])
	if err != nil {
		return nil, err
	}

	return js.NewTernaryOp(jsName, a, b, c, op.Context())
}

func (p *JSParser) buildBinaryOpExpression(t raw.Token) (js.Expression, error) {
	op, err := raw.AssertAnyBinaryOperator(t)
	if err != nil {
		panic(err)
	}

	jsName, ok := p.translateOpName(op.Name())
	if !ok {
		errCtx := op.Context()
		return nil, errCtx.NewError("Error: operator not yet handled")
	}

	a, err := p.buildExpression(op.Args()[0:1])
	if err != nil {
		return nil, err
	}

	b, err := p.buildExpression(op.Args()[1:2])
	if err != nil {
		return nil, err
	}

	return js.NewBinaryOp(jsName, a, b, op.Context())
}

func (p *JSParser) buildUnaryOpExpression(t raw.Token) (js.Expression, error) {
	op, err := raw.AssertAnyUnaryOperator(t)
	if err != nil {
		panic(err)
	}

	jsName, ok := p.translateOpName(op.Name())
	if !ok {
		errCtx := op.Context()
		return nil, errCtx.NewError("Error: operator not yet handled")
	}

	a, err := p.buildExpression(op.Args())
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(op.Name(), "post") {
		return js.NewPostUnaryOp(jsName, a, op.Context())
	} else {
		return js.NewPreUnaryOp(jsName, a, op.Context())
	}
}

func (p *JSParser) buildVarExpression(t raw.Token) (*js.VarExpression, error) {
	name, err := raw.AssertWord(t)
	if err != nil {
		return nil, err
	}

	return js.NewVarExpression(name.Value(), name.Context()), nil
}

func (p *JSParser) buildLiteralNullExpression(t raw.Token) (js.Expression, error) {
	return js.NewLiteralNull(t.Context()), nil
}

func (p *JSParser) buildIndexExpression(ts []raw.Token) (js.Expression, error) {
	n := len(ts)

	lhs, err := p.buildExpression(ts[0 : n-1])
	if err != nil {
		return nil, err
	}

	if err := js.AssertIndexable(lhs); err != nil {
		return nil, err
	}

	group, err := raw.AssertBracketsGroup(ts[n-1])
	if err != nil {
		return nil, err
	}

	if group.IsEmpty() {
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: index can't be empty")
	}

	if !group.IsSingle() { // comma's should've been combined in operators
		errCtx := group.Context()
		return nil, errCtx.NewError("Error: multi indexing not allowed")
	}

	field := group.Fields[0]

	index, err := p.buildExpression(field)
	if err != nil {
		return nil, err
	}

	return js.NewIndex(lhs, index, group.Context()), nil
}

func (p *JSParser) buildMemberExpression(ts []raw.Token) (js.Expression, error) {
	n := len(ts)

	if n < 3 {
		errCtx := ts[n-2].Context()
		return nil, errCtx.NewError("Error: member of nothing")
	}

	lhs, err := p.buildExpression(ts[0 : n-2])
	if err != nil {
		return nil, err
	}

	w, err := raw.AssertWord(ts[n-1])
	if err != nil {
		panic(err)
	}

	return js.NewMember(lhs, js.NewWord(w.Value(), w.Context()), ts[n-2].Context()), nil
}

func (p *JSParser) condensePackagePeriods(ts []raw.Token) (*raw.Word, []raw.Token, error) {
	nameToken, err := raw.AssertWord(ts[0])
	if err != nil {
		return nil, nil, err
	}

	i := 1
	for i < len(ts) {
		if raw.IsSymbol(ts[i], patterns.PERIOD) {
			if i == len(ts)-1 {
				errCtx := ts[i].Context()
				return nil, nil, errCtx.NewError("Error: expected tokens after .")
			}

			nextWord, err := raw.AssertWord(ts[i+1])
			if err != nil {
				return nil, nil, err
			}

			nameToken = raw.NewValueWord(nameToken.Value()+"."+nextWord.Value(), raw.MergeContexts(nameToken, ts[i], ts[i+1]))
			i += 2
		} else {
			break
		}
	}

	return nameToken, ts[i:], nil
}

func (p *JSParser) buildTypeExpression(ts []raw.Token) (*js.TypeExpression, error) {
	nameToken, ts, err := p.condensePackagePeriods(ts)
	if err != nil {
		return nil, err
	}

	var contentKeys []*js.Word = nil
	var contentTypes []*js.TypeExpression = nil
	if len(ts) == 1 {
		angled, err := raw.AssertAngledGroup(ts[0])
		if err != nil {
			return nil, err
		}

		contentTypes = make([]*js.TypeExpression, 0)

		if len(angled.Fields) == 0 || len(angled.Fields[0]) == 0 {
			// empty content types, usefull for empty objects
			contentKeys = make([]*js.Word, 0)
		} else {

			somePositional := false
			someKeyed := false
			for _, field := range angled.Fields {
				if !raw.ContainsSymbol(field, patterns.COLON) {
					if someKeyed {
						errCtx := field[0].Context()
						return nil, errCtx.NewError("Error: unexpected positional content type (hint: all must be positional or all must be keyed)")
					}

					somePositional = true

					contentType, err := p.buildTypeExpression(field)
					if err != nil {
						return nil, err
					}

					contentTypes = append(contentTypes, contentType)
				} else {
					if somePositional {
						errCtx := field[0].Context()
						return nil, errCtx.NewError("Error: unexpected keyed content type (hint: all must be positional or all must be keyed")
					}

					if !someKeyed {
						contentKeys = make([]*js.Word, 0)
						someKeyed = true
					}

					components := p.splitBySeparator(field, patterns.COLON)

					if len(components) != 2 || len(components[0]) != 1 {
						errCtx := raw.MergeContexts(field...)
						return nil, errCtx.NewError("Error: bad keyed content type key-value entry")
					}

					key, err := raw.AssertWord(components[0][0])
					if err != nil {
						return nil, err
					}

					contentType, err := p.buildTypeExpression(components[1])
					if err != nil {
						return nil, err
					}

					contentKeys = append(contentKeys, js.NewWord(key.Value(), key.Context()))
					contentTypes = append(contentTypes, contentType)
				}
			}
		}
	} else if len(ts) > 1 {
		errCtx := raw.MergeContexts(ts...)
		return nil, errCtx.NewError("Error: bad type expression")
	}

	return js.NewTypeExpression(nameToken.Value(), contentKeys, contentTypes, nameToken.Context())
}

// extracts angled brackets
func (p *JSParser) buildClassOrExtendsTypeExpression(ts []raw.Token) (*js.TypeExpression, []raw.Token, error) {
	if len(ts) == 1 {
		te, err := p.buildTypeExpression(ts)
		return te, []raw.Token{}, err
	} else {
		if raw.IsAnyWord(ts[0]) && raw.IsAngledGroup(ts[1]) {
			te, err := p.buildTypeExpression(ts[:2])
			return te, ts[2:], err
		} else if raw.IsAnyWord(ts[0]) {
			w, rem, err := p.condensePackagePeriods(ts)
			if err != nil {
				return nil, nil, err
			}

			te, err := p.buildTypeExpression([]raw.Token{w})
			return te, rem, err
		} else {
			errCtx := raw.MergeContexts(ts...)
			return nil, nil, errCtx.NewError("Error: bad type expression")
		}
	}
}

func (p *JSParser) buildExpression(ts []raw.Token) (js.Expression, error) {
	ts = p.expandAngledGroups(ts)

	ts, err := p.nestOperators(ts)
	if err != nil {
		return nil, err
	}

	ts = p.expandTmpGroups(ts)

	n := len(ts)
	switch {
	case n == 1:
		switch {
		case raw.IsLiteralString(ts[0]):
			return p.buildLiteralStringExpression(ts[0])
		case raw.IsTemplateGroup(ts[0]):
			return p.buildTemplateStringExpression(ts[0])
		case raw.IsLiteralInt(ts[0]):
			return p.buildLiteralIntExpression(ts[0])
		case raw.IsLiteralFloat(ts[0]):
			return p.buildLiteralFloatExpression(ts[0])
		case raw.IsSpecialNumber(ts[0]):
			return p.buildSpecialNumberExpression(ts[0])
		case raw.IsLiteralBool(ts[0]):
			return p.buildLiteralBoolExpression(ts[0])
		case raw.IsBracketsGroup(ts[0]):
			return p.buildLiteralArrayExpression(ts[0])
		case raw.IsBracesGroup(ts[0]):
			return p.buildLiteralObjectExpression(ts[0])
		case raw.IsParensGroup(ts[0]):
			return p.buildParensExpression(ts[0])
		case raw.IsAnyTernaryOperator(ts[0]):
			return p.buildTernaryOpExpression(ts[0])
		case raw.IsAnyBinaryOperator(ts[0]):
			return p.buildBinaryOpExpression(ts[0])
		case raw.IsAnyUnaryOperator(ts[0]):
			return p.buildUnaryOpExpression(ts[0])
		case raw.IsAnyWord(ts[0]): // variable
			return p.buildVarExpression(ts[0])
		case raw.IsLiteralNull(ts[0]): // null can be used as placeholder for any value type
			return p.buildLiteralNullExpression(ts[0])
		default:
			// nested group
			if raw.IsTmpGroup(ts[0]) {
				gr, err := raw.AssertGroup(ts[0])
				if err != nil {
					panic(err)
				}

				return p.buildExpression(gr.Fields[0])
			} else {
				errCtx := ts[0].Context()
				err := errCtx.NewError("Error: expression not yet supported")
				panic(err)
				return nil, err
			}
		}
	case raw.IsParensGroup(ts[n-1]):
		return p.buildCallExpression(ts)
	case raw.IsBracketsGroup(ts[n-1]):
		return p.buildIndexExpression(ts)
	case raw.IsSymbol(ts[n-2], patterns.PERIOD) &&
		raw.IsAnyWord(ts[n-1]):
		return p.buildMemberExpression(ts)
	case n > 2 &&
		raw.IsWord(ts[0], "function") &&
		raw.IsParensGroup(ts[n-2]) &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildFunctionExpression(ts, false)
	case n > 3 &&
		raw.IsWord(ts[0], "function") &&
		(raw.IsParensGroup(ts[1]) || raw.IsParensGroup(ts[2])) &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildFunctionExpression(ts, false)
	case n > 3 &&
		raw.IsWord(ts[0], "async") &&
		raw.IsWord(ts[1], "function") &&
		raw.IsParensGroup(ts[2]) &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildFunctionExpression(ts, false)
	case n > 2 &&
		raw.IsSymbol(ts[n-2], patterns.ARROW) &&
		raw.IsParensGroup(ts[0]) &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildFunctionExpression(raw.Concat(ts[0:n-2], ts[n-1:]), true)
	case n > 3 &&
		raw.IsSymbol(ts[n-2], patterns.ARROW) &&
		raw.IsWord(ts[0], "async") &&
		raw.IsParensGroup(ts[1]) &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildFunctionExpression(raw.Concat(ts[0:n-2], ts[n-1:]), true)
	case n > 1 &&
		raw.IsWord(ts[0], "class") &&
		raw.IsBracesGroup(ts[n-1]):
		return p.buildClassExpression(ts)
	default:
		for _, t := range ts {
			fmt.Println(t.Dump(""))
		}
		errCtx := raw.MergeContexts(ts...)
		err := errCtx.NewError("Error: unhandled expression (hint: missing a semicolon?)")
    panic(err)
		return nil, err
	}
}
