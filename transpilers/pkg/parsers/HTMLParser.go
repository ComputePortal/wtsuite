package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../tokens/context"
	"../tokens/html"
	"../tokens/patterns"
	"../tokens/raw"
)

func tokenizeHTMLWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
	switch {
	case patterns.IsColor(s):
		return raw.NewLiteralColor(s, ctx)
	case patterns.IsInt(s):
		return raw.NewLiteralInt(s, ctx)
	case patterns.IsFloat(s):
		return raw.NewLiteralFloat(s, ctx)
	case patterns.IsBool(s):
		return raw.NewLiteralBool(s, ctx)
	case patterns.IsNull(s):
		return raw.NewLiteralNull(ctx), nil
	case patterns.IsWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

func tokenizeHTMLFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	fp, err := NewFormulaParser(s, ctx)
	if err != nil {
		return nil, err
	}

	return fp.tokenize()
}

var htmlParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.HTML_STRING_OR_COMMENT_REGEXP,
		groups: []quotedGroupSettings{
			quotedGroupSettings{
				maskType:        STRING,
				groupPattern:    patterns.SQ_STRING_GROUP,
				assertStopMatch: false,
				info:            "single quotes",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        STRING,
				groupPattern:    patterns.DQ_STRING_GROUP,
				assertStopMatch: false,
				info:            "double quotes",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        FORMULA,
				groupPattern:    patterns.BT_FORMULA_GROUP,
				assertStopMatch: false,
				info:            "backtick formula",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        COMMENT,
				groupPattern:    patterns.SL_COMMENT_GROUP,
				assertStopMatch: false,
				info:            "single-line comment",
				trackStarts:     false,
			},
			quotedGroupSettings{
				maskType:        COMMENT,
				groupPattern:    patterns.ML_COMMENT_GROUP,
				assertStopMatch: true,
				info:            "js-style multiline comment",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        COMMENT,
				groupPattern:    patterns.XML_COMMENT_GROUP,
				assertStopMatch: true,
				info:            "xml-style multiline comment",
				trackStarts:     true,
			},
		},
	},
	formulas: formulasSettings{
		tokenizer: tokenizeHTMLFormulas,
	},
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.HTML_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeHTMLWordsAndLiterals,
	},
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.HTML_SYMBOLS_REGEXP,
	},
	operators: newOperatorsSettings([]operatorSettings{
		//operatorSettings{4, ":", BIN},
	}),
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
}

var htmlFunctionMap = map[string]string{
	"pre-":   "neg",
	"pre!":   "not",
	"bin/":   "div",
	"bin*":   "mul",
	"bin-":   "sub",
	"bin+":   "add",
	"bin<":   "lt",
	"bin<=":  "le",
	"bin>":   "gt",
	"bin>=":  "ge",
	"bin!=":  "ne",
	"bin==":  "eq",
	"bin===": "issame",
	"bin||":  "or",
	"bin&&":  "and",
	// ":", "?", ":=" and "=" are treated explicitely
}

type HTMLParser struct {
	Parser
}

func NewHTMLParser(path string) (*HTMLParser, error) {
	if !filepath.IsAbs(path) {
		panic("path should be absolute")
	}

	rawBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := string(rawBytes)
	src := context.NewSource(raw)

	ctx := context.NewContext(src, path)
	p := &HTMLParser{newParser(raw, htmlParserSettings, ctx)}

	if err := p.maskQuoted(); err != nil {
		return nil, err
	}

	return p, nil
}

func NewEmptyHTMLParser(ctx context.Context) *HTMLParser {
	return &HTMLParser{newParser("", htmlParserSettings, ctx)}
}

func (p *HTMLParser) Refine(start, stop int) *HTMLParser {
	return &HTMLParser{p.refine(start, stop)}
}

func (p *HTMLParser) ChangeCaller(caller string) *HTMLParser {
	return &HTMLParser{p.changeCaller(caller)}
}

func (p *HTMLParser) tokenize() ([]raw.Token, error) {
	ts, err := p.Parser.tokenize()
	if err != nil {
		return nil, err
	}

	return p.nestOperators(ts)
}

func (p *HTMLParser) buildOperatorToken(v *raw.Operator) (html.Token, error) {
	switch {
	case v.Name() == "bin:" && raw.IsBinaryOperator(v.Args()[0], "bin?"): // actually a ternary operator

		ab, err := raw.AssertBinaryOperator(v.Args()[0], "bin?")
		if err != nil {
			return nil, err
		}

		a, err := p.buildToken(ab.Args()[0:1])
		if err != nil {
			return nil, err
		}

		b, err := p.buildToken(ab.Args()[1:2])
		if err != nil {
			return nil, err
		}

		c, err := p.buildToken(v.Args()[1:2])
		if err != nil {
			return nil, err
		}

		return html.NewValueFunction("ifelse", []html.Token{a, b, c}, v.Context()), nil
	case v.Name() == "bin:=": // lhs must be word
		// accept word or string
		arg0 := v.Args()[0]

		if raw.IsAnyWord(arg0) {
			a, err := raw.AssertWord(arg0)
			if err != nil {
				panic("unexpected")
			}

			b, err := p.buildToken(v.Args()[1:2])
			if err != nil {
				return nil, err
			}

			return html.NewValueFunction("new", []html.Token{html.NewValueString(a.Value(), a.Context()), b}, v.Context()), nil
		} else {
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: lhs must be word (hint: missing semicolon?)")
		}
	case strings.HasPrefix(v.Name(), "bin"):
		a, err := p.buildToken(v.Args()[0:1])
		if err != nil {
			return nil, err
		}
		b, err := p.buildToken(v.Args()[1:2])
		if err != nil {
			return nil, err
		}
		if fnName, ok := htmlFunctionMap[v.Name()]; ok {
			return html.NewValueFunction(fnName, []html.Token{a, b}, v.Context()), nil
		} else {
			errCtx := v.Context()
			err := errCtx.NewError("Error: binary operator '" + strings.TrimLeft(v.Name(), "bin") + "' not recognized")
			panic(err)
			return nil, err
		}
	case strings.HasPrefix(v.Name(), "pre"):
		a, err := p.buildToken(v.Args())
		if err != nil {
			return nil, err
		}
		if fnName, ok := htmlFunctionMap[v.Name()]; ok {
			return html.NewValueFunction(fnName, []html.Token{a}, v.Context()), nil
		} else {
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: pre unary operator '" + strings.TrimLeft(v.Name(), "pre") + "' not recognized")
		}
	case strings.HasPrefix(v.Name(), "post"):
		a, err := p.buildToken(v.Args())
		if err != nil {
			return nil, err
		}
		if fnName, ok := htmlFunctionMap[v.Name()]; ok {
			return html.NewValueFunction(fnName, []html.Token{a}, v.Context()), nil
		} else {
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: post unary operator '" + strings.TrimLeft(v.Name(), "post") + "' not recognized")
		}
	case strings.HasPrefix(v.Name(), "sing"):
		if fnName, ok := htmlFunctionMap[v.Name()]; ok {
			return html.NewValueFunction(fnName, []html.Token{}, v.Context()), nil
		} else {
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: singular operator '" + strings.TrimLeft(v.Name(), "sing") + "' not recognized")
		}
	default:
		errCtx := v.Context()
		return nil, errCtx.NewError("Error: unrecognized operator '" + v.Name() + "'")
	}
}

func (p *HTMLParser) buildParensGroupToken(v *raw.Group) (*html.Parens, error) {
	if v.IsParens() && (v.IsEmpty() || v.IsSingle() || v.IsComma()) {
		values := make([]html.Token, 0)
		alts := make([]html.Token, 0) // if first token is string, and second is '=', then remainder is alt, otherwise nil
		for _, field := range v.Fields {
			if raw.IsBinaryOperator(field[0], "bin=") {
				eq, err := raw.AssertBinaryOperator(field[0], "bin=")
				if err != nil {
					panic(err)
				}

				a, err := p.buildToken(eq.Args()[0:1])
				if err != nil {
					return nil, err
				}

				b, err := p.buildToken(eq.Args()[1:2])
				if err != nil {
					return nil, err
				}

				values = append(values, a)
				alts = append(alts, b)
			} else {
				val, err := p.buildToken(field)
				if err != nil {
					return nil, err
				}

				values = append(values, val)
				alts = append(alts, nil)
			}
		}

		return html.NewParens(values, alts, v.Context()), nil
	} else {
		errCtx := v.Context()
		return nil, errCtx.NewError("Error: bad parens")
	}
}

func (p *HTMLParser) buildBracesGroupToken(v *raw.Group) (*html.RawDict, error) {
	if v.IsBraces() && (v.IsEmpty() || v.IsSingle() || v.IsComma()) {
		keys := make([]html.Token, 0)
		values := make([]html.Token, 0)

		for _, field := range v.Fields {
			if len(field) != 1 {
				errCtx := v.Context()
				return nil, errCtx.NewError("Error: bad dict content")
			}
			colon, err := raw.AssertBinaryOperator(field[0], "bin:")
			if err != nil {
				return nil, err
			}

			a, err := p.buildToken(colon.Args()[0:1])
			if err != nil {
				return nil, err
			}

			b, err := p.buildToken(colon.Args()[1:2])
			if err != nil {
				return nil, err
			}

			keys = append(keys, a)
			values = append(values, b)
		}

		return html.NewValuesRawDict(keys, values, v.Context()), nil
	} else {
		errCtx := v.Context()
		return nil, errCtx.NewError("Error: bad braces")
	}
}

// return value can be List or seq(...) function call
func (p *HTMLParser) buildBracketsGroupToken(v *raw.Group) (html.Token, error) {
	if v.IsBrackets() && (v.IsEmpty() || v.IsSingle() || v.IsComma()) {
		if v.IsSingle() && len(v.Fields[0]) == 1 && raw.IsOperator(v.Fields[0][0], "bin:") {
			op, err := raw.AssertAnyOperator(v.Fields[0][0])
			if err != nil {
				panic(err)
			}

			start, err := p.buildToken(op.Args()[0:1])
			if err != nil {
				return nil, err
			}

			if raw.IsOperator(op.Args()[1], "bin:") {
				op2, err := raw.AssertAnyOperator(op.Args()[1])
				if err != nil {
					panic(err)
				}

				incr, err := p.buildToken(op2.Args()[0:1])
				if err != nil {
					return nil, err
				}

				stop, err := p.buildToken(op2.Args()[1:])
				if err != nil {
					return nil, err
				}

				return html.NewValueFunction("seq", []html.Token{start, incr, stop},
					context.MergeContexts(v.Context(), op.Context(), op2.Context())), nil
			} else {
				errCtx := v.Context()
				return nil, errCtx.NewError("Error: forming sequence like this is not allowed, because it is too easily confused with ','")
				/*stop, err := p.buildToken(op.Args()[1:])
				if err != nil {
					return nil, err
				}

				return html.NewValueFunction("seq", []html.Token{start, stop},
					context.MergeContexts(v.Context(), op.Context())), nil*/
			}
		} else {
			values := make([]html.Token, 0)

			for _, field := range v.Fields {
				a, err := p.buildToken(field)
				if err != nil {
					return nil, err
				}

				values = append(values, a)
			}

			return html.NewValuesList(values, v.Context()), nil
		}
	} else {
		errCtx := v.Context()
		return nil, errCtx.NewError("Error: bad brackets")
	}
}

func (p *HTMLParser) buildGroupToken(v *raw.Group) (html.Token, error) {
	switch {
	case v.IsParens():
		return p.buildParensGroupToken(v)
	case v.IsBraces():
		return p.buildBracesGroupToken(v)
	case v.IsBrackets():
		return p.buildBracketsGroupToken(v)
	default:
		errCtx := v.Context()
		return nil, errCtx.NewError("Error: unhandled group type")
	}
}

func (p *HTMLParser) buildDefineFunctionToken(vs []raw.Token) (html.Token, []raw.Token, error) {
	// new function
	a, err := raw.AssertWord(vs[0])
	if err != nil {
		panic("unexpected")
	}

	if a.Value() != "function" {
		errCtx := a.Context()
		return nil, nil, errCtx.NewError("Error: expected function keyword")
	}

	argsGroup, err := raw.AssertParensGroup(vs[1])
	if err != nil {
		panic("unexpected")
	}

	argsWithDefaults, err := p.buildParensGroupToken(argsGroup)
	if err != nil {
		return nil, nil, err
	}

	statementsGroup, err := raw.AssertBracesGroup(vs[2])
	if err != nil {
		return nil, nil, err
	}

	if !(statementsGroup.IsSingle() || statementsGroup.IsSemiColon()) {
		errCtx := vs[2].Context()
		return nil, nil, errCtx.NewError("Error: bad statements for function def")
	}

	statements := make([]html.Token, 0)

	for _, field := range statementsGroup.Fields {
		st, err := p.buildToken(field)
		if err != nil {
			return nil, nil, err
		}

		statements = append(statements, st)
	}

	// wrap statements in a get
	ctx := vs[2].Context()
	list := html.NewValuesList(statements, ctx)
	index := html.NewValueInt(len(statements)-1, ctx)
	wrapper := html.NewValueFunction("get", []html.Token{list, index}, ctx)

	return html.NewValueFunction("function", []html.Token{argsWithDefaults, wrapper}, ctx), vs[3:], nil
}

// also return the remaining
func (p *HTMLParser) buildFunctionToken(vs []raw.Token) (html.Token, []raw.Token, error) {
	// new function
	a, err := raw.AssertWord(vs[0])
	if err != nil {
		panic("unexpected")
	}

	argsGroup, err := raw.AssertParensGroup(vs[1])
	if err != nil {
		panic("unexpected")
	}

	if !(argsGroup.IsEmpty() || argsGroup.IsSingle() || argsGroup.IsComma()) {
		errCtx := vs[1].Context()
		return nil, nil, errCtx.NewError("Error: bad function args")
	}

	args_, err := p.buildGroupToken(argsGroup)
	if err != nil {
		return nil, nil, err
	}

	args, err := html.AssertParens(args_)
	if err != nil {
		panic(err)
	}

	// check that all alts are nil
	for _, alt := range args.Alts() {
		if alt != nil {
			errCtx := alt.Context()
			return nil, nil, errCtx.NewError("Error: unexpected arg expression")
		}
	}

	return html.NewValueFunction(a.Value(), args.Values(),
		context.MergeContexts(a.Context(), vs[1].Context())), vs[2:], nil
}

func (p *HTMLParser) buildIndexedToken(vs []raw.Token) (html.Token, []raw.Token, error) {
	indices := make([]html.Token, 0)

	varName, err := raw.AssertWord(vs[0])
	if err != nil {
		return nil, nil, err
	}

	obj := html.NewValueString(varName.Value(), varName.Context())

	ctx := vs[0].Context()
	for _, v := range vs[1:] {
		if raw.IsBracketsGroup(v) {
			indexGroup, err := raw.AssertBracketsGroup(v)
			if err != nil {
				return nil, nil, err
			}

			if !indexGroup.IsSingle() {
				errCtx := indexGroup.Context()
				return nil, nil, errCtx.NewError("Error: bad index (hint: multi indexing not supported)")
			}

			field := indexGroup.Fields[0]
			if len(field) == 1 && (raw.IsOperator(field[0], "sing:") || raw.IsOperator(field[0], "pre:") || raw.IsOperator(field[0], "post:") || raw.IsOperator(field[0], "bin:")) {
				break // these require a slice instead
			}

			index, err := p.buildToken(field)
			if err != nil {
				return nil, nil, err
			}

			indices = append(indices, index)
		} else {
			break
		}
	}

	// nest these, so have get(get(get(dictname), "key"), index) etc.
	res := html.NewValueFunction("get", []html.Token{obj}, ctx)
	for _, index := range indices {
		res = html.NewValueFunction("get", []html.Token{res, index}, index.Context())
	}

	return res, vs[(len(indices) + 1):], nil
}

func (p *HTMLParser) buildEvalsAndIndexing(obj html.Token, vs []raw.Token) (html.Token, error) {
	for _, v := range vs {
		if !raw.IsGroup(v) {
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: unexpected")
		}

		gr, err := raw.AssertGroup(v)
		if err != nil {
			panic(err)
		}

		switch {
		case gr.IsBrackets() && gr.IsSingle():
			field := gr.Fields[0]

			if len(field) == 1 && (raw.IsOperator(field[0], "sing:") ||
				raw.IsOperator(field[0], "post:") ||
				raw.IsOperator(field[0], "pre:") ||
				raw.IsOperator(field[0], "bin:")) {
				op, err := raw.AssertAnyOperator(field[0])
				if err != nil {
					panic(err)
				}

				switch {
				case op.Name() == "sing:":
					ctx := op.Context()
					obj = html.NewValueFunction("slice", []html.Token{obj, html.NewNull(ctx),
						html.NewValueInt(1, ctx), html.NewNull(ctx)}, ctx)
				case op.Name() == "post:":
					ctx := op.Context()
					a, err := p.buildToken(op.Args())
					if err != nil {
						return nil, err
					}

					obj = html.NewValueFunction("slice", []html.Token{obj, a,
						html.NewValueInt(1, ctx), html.NewNull(ctx)}, ctx)
				case op.Name() == "pre:" || op.Name() == "bin:":
					var start html.Token
					if op.Name() == "pre:" {
						start = html.NewNull(op.Context())
					} else {
						start, err = p.buildToken(op.Args()[0:1])
						if err != nil {
							return nil, err
						}
					}

					op2_ := op.Args()[0]
					if raw.IsOperator(op2_, "post:") || raw.IsOperator(op2_, "bin:") {
						op2, err := raw.AssertAnyOperator(op2_)
						if err != nil {
							panic(err)
						}

						switch {
						case op2.Name() == "post:":
							incr, err := p.buildToken(op2.Args())
							if err != nil {
								return nil, err
							}

							stop := html.NewNull(op2.Context())

							obj = html.NewValueFunction("slice", []html.Token{obj, start,
								incr, stop}, context.MergeContexts(op.Context(), op2.Context()))
						case op2.Name() == "bin:":
							incr, err := p.buildToken(op2.Args()[0:1])
							if err != nil {
								return nil, err
							}

							stop, err := p.buildToken(op2.Args()[1:])
							if err != nil {
								return nil, err
							}

							obj = html.NewValueFunction("slice", []html.Token{obj, start,
								incr, stop}, context.MergeContexts(op.Context(), op2.Context()))
						}
					} else {
						var stop html.Token
						if op.Name() == "pre:" {
							stop, err = p.buildToken(op.Args()[0:1])
							if err != nil {
								return nil, err
							}
						} else {
							stop, err = p.buildToken(op.Args()[1:2])
							if err != nil {
								return nil, err
							}
						}

						obj = html.NewValueFunction("slice", []html.Token{obj, start,
							html.NewValueInt(1, op.Context()), stop}, op.Context())
					}
				default:
					panic("unhandled")
				}
			} else {
				index, err := p.buildToken(gr.Fields[0])
				if err != nil {
					return nil, err
				}

				obj = html.NewValueFunction("get", []html.Token{obj, index}, gr.Context())
			}
		case gr.IsParens() && (gr.IsEmpty() || gr.IsSingle() || gr.IsComma()):
			args := make([]html.Token, 0)

			for _, field := range gr.Fields {
				arg, err := p.buildToken(field)
				if err != nil {
					return nil, err
				}

				args = append(args, arg)
			}

			obj = html.NewValueFunction("eval", []html.Token{obj, html.NewValuesList(args, gr.Context())}, gr.Context())
		default:
			errCtx := gr.Context()
			err := errCtx.NewError("Error: bad indexing/evaluating")
			panic(err)
			return nil, err
		}
	}

	return obj, nil
}

func (p *HTMLParser) buildToken(vs []raw.Token) (html.Token, error) {
	if len(vs) == 1 {
		if tmp, ok := vs[0].(*raw.Group); ok {
			if tmp.IsTmp() {
				vs = tmp.Fields[0]
			}
		}
	}

	switch len(vs) {
	case 0:
		panic("expected at least one token")
	case 1:
		switch v := vs[0].(type) {
		case *raw.LiteralBool:
			return html.NewValueBool(v.Value(), v.Context()), nil
		case *raw.LiteralColor:
			r, g, b, a := v.Values()
			return html.NewValueColor(r, g, b, a, v.Context()), nil
		case *raw.LiteralFloat:
			return html.NewValueUnitFloat(v.Value(), v.Unit(), v.Context()), nil
		case *raw.LiteralInt:
			return html.NewValueInt(v.Value(), v.Context()), nil
		case *raw.LiteralNull:
			return html.NewNull(v.Context()), nil
		case *raw.LiteralString:
			if v.WasWord() {
				return html.NewWordString(v.Value(), v.Context()), nil
			} else {
				return html.NewValueString(v.Value(), v.Context()), nil
			}
		case *raw.Word:
			return html.NewValueFunction("get", []html.Token{html.NewValueString(v.Value(), v.Context())},
				v.Context()), nil
		case *raw.Operator:
			// NOTE: raw.Operator tokens are generated by FormulaParser
			// bin= is a special case, and the bin= function call is just used a placeholder
			return p.buildOperatorToken(v)
		case *raw.Group:
			if v.IsTmp() {
				return p.buildToken(v.Fields[0])
			} else {
				return p.buildGroupToken(v)
			}
		default:
			errCtx := v.Context()
			return nil, errCtx.NewError("Error: invalid syntax")
		}
	default:
		if len(vs) >= 3 && raw.IsAnyWord(vs[0]) && raw.IsParensGroup(vs[1]) && raw.IsBracesGroup(vs[2]) {
			fn, remaining, err := p.buildDefineFunctionToken(vs)
			if err != nil {
				return nil, err
			}

			if len(remaining) == 0 {
				return fn, nil
			} else {
				return p.buildEvalsAndIndexing(fn, remaining)
			}
		} else if len(vs) >= 2 && raw.IsAnyWord(vs[0]) && raw.IsParensGroup(vs[1]) {
			fn, remaining, err := p.buildFunctionToken(vs)
			if err != nil {
				return nil, err
			}

			if len(remaining) == 0 {
				return fn, nil
			} else {
				return p.buildEvalsAndIndexing(fn, remaining)
			}
		} else if len(vs) >= 2 && raw.IsAnyWord(vs[0]) && raw.IsBracketsGroup(vs[1]) {
			obj, remaining, err := p.buildIndexedToken(vs)
			if err != nil {
				return nil, err
			}

			if len(remaining) == 0 {
				return obj, nil
			} else {
				return p.buildEvalsAndIndexing(obj, remaining)
			}
		} else {
			obj, err := p.buildToken(vs[0:1])
			if err != nil {
				return nil, err
			}

			remaining := vs[1:]
			if len(remaining) == 0 {
				return obj, nil
			} else {
				return p.buildEvalsAndIndexing(obj, remaining)
			}
		}
	}
}

func (p *HTMLParser) parseAttributes(ctx context.Context) (*html.RawDict, error) {
	ts, err := p.tokenize()
	if err != nil {
		return nil, err
	}

	ts = p.expandTmpGroups(ts)

	result := html.NewEmptyRawDict(ctx)

	appendKeyVal := func(k *raw.Word, v html.Token) error {
		if other, _, ok := result.GetKeyValue(k.Value()); ok {
			errCtx := context.MergeContexts(k.Context(), other.Context())
			return errCtx.NewError("Error: duplicate")
		} else {
			s, _ := html.NewString(k.Value(), k.Context())
			result.Set(s, v)
			return nil
		}
	}

	appendFlag := func(k *raw.Word) error {
		return appendKeyVal(k, html.NewFlag(k.Context()))
	}

	convertAppendKeyVal := func(k *raw.Word, vs []raw.Token) error {
		v, err := p.buildToken(vs)
		if err != nil {
			return err
		}

		return appendKeyVal(k, v)
	}

	i := 0
	for i < len(ts) {
		key, err := raw.AssertWord(ts[i])
		if err != nil {
			return nil, err
		}

		if (i + 1) < len(ts) {
			switch t := ts[i+1].(type) {
			case *raw.Symbol:
				if _, err := raw.AssertSymbol(t, patterns.EQUAL); err != nil {
					return nil, err
				}

				if (i + 2) < len(ts) {
					val := ts[i+2]
					if err := raw.AssertNotSymbol(val); err != nil {
						return nil, err
					}

					vs := []raw.Token{val}

					if (i + 3) < len(ts) {
						if raw.IsGroup(ts[i+3]) {
							vs = append(vs, ts[i+3])
							i += 1
						}
					}

					if err := convertAppendKeyVal(key, vs); err != nil {
						return nil, err
					}

					i += 3
				} else {
					errCtx := t.Context()
					return nil, errCtx.NewError("Syntax Error: expected more")
				}
			case *raw.Word:
				if err := appendFlag(key); err != nil {
					return nil, err
				}
				// leave ts[i+1] to next iteration
				i++
			default:
				errCtx := t.Context()
				return nil, errCtx.NewError("Syntax Error: bad attribute")
			}
		} else {
			// append a flag
			if err := appendFlag(key); err != nil {
				return nil, err
			}

			i++
		}
	}

	return result, nil
}

func (p *HTMLParser) BuildTags() ([]*html.Tag, error) {
	rprev := [2]int{0, 0}

	result := make([]*html.Tag, 0)

	appendTag := func(t *html.Tag) {
		if t == nil {
			panic("tag is nil")
		}

		result = append(result, t)
	}

	for true {
		if r, _, ok := p.nextMatch(patterns.TAG_START_REGEXP, false); ok {
			// handle non-tag text that wasn't matched
			if r[0] > rprev[1] {
				subContent := p.Refine(rprev[1], r[0])
				if !subContent.IsEmpty() {
					appendTag(html.NewTextTag(p.Write(rprev[1], r[0]),
						p.NewContext(rprev[1], r[0])))
				}
			}

			rprev = r

			if rr, _, ok := p.nextMatch(patterns.DUMMY_TAG_NAME_REGEXP, false); ok {
				ctx := p.NewContext(rr[0], rr[1])
				rprev = rr

				if rrr, ok := p.nextGroupStopMatch(patterns.NewTagGroup(""), true); ok {
					ctx = context.MergeContexts(ctx, p.NewContext(rrr[0], rrr[1]))

					subParser := p.Refine(rr[1], rrr[0])
					subTags, err := subParser.BuildTags()
					if err != nil {
						return nil, err
					}

					rprev = rrr

					appendTag(html.NewTag(patterns.HTML_DUMMY_TAG_NAME, html.NewEmptyRawDict(ctx),
						subTags, ctx))
				}
			} else if rname, name, ok := p.nextMatch(patterns.TAG_NAME_REGEXP, false); ok {
				stopRegexp := patterns.TAG_STOP_REGEXP
				if name == "?xml" {
					stopRegexp = patterns.XML_HEADER_STOP_REGEXP
				}

				if rr, s, ok := p.nextMatch(stopRegexp, false); ok {
					ctx := context.MergeContexts(p.NewContext(r[0], rname[1]), p.NewContext(rr[0], rr[1]))

					attrParser := p.Refine(rname[1], rr[0])
					attr, err := attrParser.parseAttributes(ctx) // this is where the magic happens
					if err != nil {
						return nil, err
					}

					rprev = rr

					var subParser *HTMLParser = nil
					if patterns.IsSelfClosing(name, s) {
						subParser = p.Refine(rr[1], rr[1])
					} else {
						if rrr, ok := p.nextGroupStopMatch(patterns.NewTagGroup(name), true); ok {
							ctx = context.MergeContexts(ctx, p.NewContext(rrr[0], rrr[1]))
							subParser = p.Refine(rr[1], rrr[0])
							rprev = rrr
						} else {
							return nil, ctx.NewError("Syntax Error: unmatched tag")
						}
					}

					if name == "math" || name == "script" || name == "style" {
						appendTag(html.NewScriptTag(name, attr, subParser.Write(0, -1),
							subParser.NewContext(0, -1), ctx))
					} else {
						subTags, err := subParser.BuildTags()
						if err != nil {
							return nil, err
						}

						appendTag(html.NewTag(name, attr, subTags, ctx))
					}
				} else {
					return nil, p.NewError(r[0], rname[1], "Syntax Error: tag not closed")
				}
			} else {
				return nil, p.NewError(r[0], r[1], "Syntax Error: tag name not found")
			}

		} else {
			break
		}
	}

	if rprev[1] < p.Len() {
		subParser := p.Refine(rprev[1], -1)
		if !subParser.IsEmpty() {
			appendTag(html.NewTextTag(subParser.Write(0, -1), subParser.NewContext(0, -1)))
		}
	}

	return result, nil
}

func (p *HTMLParser) DumpTokens() {
	fmt.Println("\nHTML tokens:")
	fmt.Println("============")

	tags, err := p.BuildTags()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, tag := range tags {
		fmt.Println(tag.Dump(""))
	}
}

// only used for style dicts in conventional html "key1:value1;key2:value2"
// values are always strings
func ParseInlineDict(rawInput string, ctx context.Context) (*html.StringDict, error) {
	empty := html.NewEmptyStringDict(ctx)

	pairs := strings.Split(rawInput, ";")

	for _, pair_ := range pairs {

		if pair_ != "" {
			pair := strings.Split(pair_, ":")

			if len(pair) != 2 {
				return nil, ctx.NewError("Error: bad dict string")
			}

			var val html.Token = nil

			s := pair[1]
			switch {
			case patterns.IsColor(s):
				c, err := raw.NewLiteralColor(s, ctx)
				if err != nil {
					return nil, err
				}
				r, g, b, a := c.Values()
				val = html.NewValueColor(r, g, b, a, ctx)
			case patterns.IsInt(s):
				rawInt, err := raw.NewLiteralInt(s, ctx)
				if err != nil {
					return nil, err
				}
				val = html.NewValueInt(rawInt.Value(), ctx)
			case patterns.IsFloat(s):
				rawFloat, err := raw.NewLiteralFloat(s, ctx)
				if err != nil {
					return nil, err
				}
				val = html.NewValueUnitFloat(rawFloat.Value(), rawFloat.Unit(), ctx)
			default:
				val = html.NewValueString(s, ctx)
			}
			empty.Set(pair[0], val)
		}
	}

	return empty, nil
}
