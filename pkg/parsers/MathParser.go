package parsers

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/math"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeMathWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
	switch {
	case patterns.IsInt(s):
		return raw.NewLiteralInt(s, ctx)
	case patterns.IsPlainFloat(s):
		return raw.NewLiteralFloat(s, ctx)
	// TODO: keyword operators
	case patterns.IsWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

func tokenizeMathFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	return nil, ctx.NewError("Error: can't have backtick formula within math")
}

var mathParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.MATH_STRING_OR_COMMENT_REGEXP,
		groups: []quotedGroupSettings{
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
		tokenizer: tokenizeMathFormulas,
	},
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.MATH_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeMathWordsAndLiterals,
	},
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.MATH_SYMBOLS_REGEXP,
	},
	operators: newOperatorsSettings([]operatorSettings{
		operatorSettings{17, "-", PRE},
		operatorSettings{16, "^", BIN | L2R},
		operatorSettings{16, "_", BIN | L2R},
		operatorSettings{14, "/", BIN | L2R},
		operatorSettings{14, "*", BIN | L2R},
		operatorSettings{14, ".", BIN | L2R},
		operatorSettings{13, "-", BIN | L2R},
		operatorSettings{13, "+", BIN | L2R},
		operatorSettings{11, "<<", BIN | L2R},
		operatorSettings{11, "<", BIN | L2R},
		operatorSettings{11, "<=", BIN | L2R},
		operatorSettings{11, ">>", BIN | L2R},
		operatorSettings{11, ">", BIN | L2R},
		operatorSettings{11, ">=", BIN | L2R},
		operatorSettings{10, "!=", BIN},
		operatorSettings{10, "==", BIN},
		operatorSettings{10, "~=", BIN},
		operatorSettings{5, "=", BIN},
		operatorSettings{4, "->", BIN},
		operatorSettings{4, "=>", BIN},
	}),
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
}

var mathOperatorMap = map[string]string{
	"bin^":  "^",
	"bin_":  "_",
	"bin*":  "*",
	"bin+":  "+",
	"bin-":  "-",
	"bin/":  "/", // used for frac, not subscript!
	"bin=":  "=",
	"bin!=": "!=",
	"bin~=": "~=",
	"bin<=": "<=",
	"bin<<": "<<",
	"bin<":  "<",
	"bin>>": ">>",
	"bin>":  ">",
	"bin>=": ">=",
	"bin->": "->",
	"bin=>": "=>",
	"pre-":  "-",
}

type MathParser struct {
	Parser
}

func NewMathParser(s string, ctx context.Context) (*MathParser, error) {
	p := &MathParser{newParser(s, mathParserSettings, ctx)}

	return p, nil
}

func (p *MathParser) tokenize() ([]raw.Token, error) {
	ts, err := p.Parser.tokenize()
	if err != nil {
		return nil, err
	}

	return p.nestOperators(ts)
}

func (p *MathParser) build(ts []raw.Token) (math.Token, error) {
	ts = p.expandTmpGroups(ts)

	switch len(ts) {
	case 1:
		switch {
		case raw.IsAnyWord(ts[0]):
			w, err := raw.AssertWord(ts[0])
			if err != nil {
				panic(err)
			}

			return math.NewWord(w.Value(), w.Context())
		case raw.IsLiteralFloat(ts[0]):
			r, err := raw.AssertLiteralFloat(ts[0], "")
			if err != nil {
				panic(err)
			}

			return math.NewFloat(r.Value(), r.Context())
		case raw.IsLiteralInt(ts[0]):
			r, err := raw.AssertLiteralInt(ts[0])
			if err != nil {
				panic(err)
			}

			return math.NewFloat(float64(r.Value()), r.Context())
		case raw.IsAnyUnaryOperator(ts[0]):
			op, err := raw.AssertAnyUnaryOperator(ts[0])
			if err != nil {
				panic(err)
			}

			name, ok := mathOperatorMap[op.Name()]
			if !ok {
				errCtx := op.Context()
				return nil, errCtx.NewError("Error: operator " + op.Name() + " not yet handled")
			}

			a, err := p.build(op.Args()[0:1])
			if err != nil {
				return nil, err
			}

			if strings.HasPrefix(op.Name(), "pre") {
				return math.NewPreUnaryOp(name, a, op.Context())
			} else {
				return math.NewPostUnaryOp(name, a, op.Context())
			}
		case raw.IsAnyBinaryOperator(ts[0]):
			op, err := raw.AssertAnyBinaryOperator(ts[0])
			if err != nil {
				panic(err)
			}

			name, ok := mathOperatorMap[op.Name()]
			if !ok {
				errCtx := op.Context()
				return nil, errCtx.NewError("Error: operator " + op.Name() + " not yet handled")
			}

			ats := op.Args()[0:1]
			bts := op.Args()[1:2]
			isFrac := false

			// expand parentheses in special cases
			switch name {
			case "^":
				bts = raw.ExpandParensGroup(op.Args()[1])
			case "_":
				bts = raw.ExpandParensGroup(op.Args()[1])
			case "/":
				if raw.IsParensGroup(op.Args()[0]) && raw.IsParensGroup(op.Args()[1]) {
					isFrac = true
				}
				ats = raw.ExpandParensGroup(op.Args()[0])
				bts = raw.ExpandParensGroup(op.Args()[1])
			}

			var a math.Token = nil
			if len(ats) > 1 {
				aParts := make([]math.Token, 0)
				for _, at := range ats {
					aPart, err := p.build([]raw.Token{at})
					if err != nil {
						return nil, err
					}
					aParts = append(aParts, aPart)
				}

				a, err = math.NewGroup(aParts, ",", op.Context())
				if err != nil {
					return nil, err
				}
			} else {
				a, err = p.build(ats)
				if err != nil {
					return nil, err
				}
			}

			var b math.Token = nil
			if len(bts) > 1 {
				bParts := make([]math.Token, 0)
				for _, bt := range bts {
					bPart, err := p.build([]raw.Token{bt})
					if err != nil {
						return nil, err
					}
					bParts = append(bParts, bPart)
				}

				b, err = math.NewGroup(bParts, ",", op.Context())
				if err != nil {
					return nil, err
				}
			} else {
				b, err = p.build(bts)
				if err != nil {
					return nil, err
				}
			}

			if isFrac {
				return math.NewFracOp(a, b, op.Context())
			} else {
				return math.NewBinaryOp(name, a, b, op.Context())
			}
		case raw.IsParensGroup(ts[0]):
			gr, err := raw.AssertParensGroup(ts[0])
			if err != nil {
				panic(err)
			}

			if !gr.IsSingle() {
				errCtx := gr.Context()
				return nil, errCtx.NewError("Error: bad parens (expected 1 token as content)")
			}

			content, err := p.build(gr.Fields[0])
			if err != nil {
				return nil, err
			}

			// if content is frac, return that directly
			return math.NewParens(content, gr.Context())
		default:
			panic("not yet handled")
		}
	case 2:
		switch {
		case raw.IsAnyWord(ts[0]) && raw.IsParensGroup(ts[1]):
			w, err := raw.AssertWord(ts[0])
			if err != nil {
				panic(err)
			}

			args := make([]math.Token, 0)

			gr, err := raw.AssertParensGroup(ts[1])
			if err != nil {
				panic(err)
			}

			if w.Value() == "text" {
				// dont build child tokens!
				contentCtx := ts[1].Context()
				contentStr := contentCtx.Content()
				nContentChars := len(contentStr)
				contentStr = contentStr[1 : nContentChars-1]

				contentToken, err := math.NewWord(contentStr, contentCtx)
				if err != nil {
					panic(err)
				}
				return math.NewCall(w.Value(), []math.Token{contentToken}, context.MergeContexts(w.Context(), contentCtx))
			}

			if gr.IsSemiColon() {
				errCtx := gr.Context()
				return nil, errCtx.NewError("Error: expected only comma separator")
			}

			for _, f := range gr.Fields {
				arg, err := p.build(f)
				if err != nil {
					return nil, err
				}

				args = append(args, arg)
			}

			return math.NewCall(w.Value(), args, context.MergeContexts(w.Context(), gr.Context()))
		default:
			for _, t := range ts {
				fmt.Println(t.Dump(""))
			}
			panic("not yet handled")
		}

	default:
		for _, t := range ts {
			fmt.Println(t.Dump(""))
		}
		if len(ts) == 0 {
			panic("no input tokens")
		}
		errCtx := raw.MergeContexts(ts...)
		return nil, errCtx.NewError("Error: bad math expression")
	}
}

func (p *MathParser) Build() (math.Token, error) {
	ts, err := p.tokenize()
	if err != nil {
		return nil, err
	}

	if len(ts) < 1 {
		return nil, errors.New("Error: empty math")
	}

	return p.build(ts)
}

func (p *MathParser) DumpTokens() {
	fmt.Println("Raw tokens:")
	fmt.Println("===========")

	ts, err := p.tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, t := range ts {
		fmt.Println(t.Dump(""))
	}

	fmt.Println("\nMath tokens:")
	fmt.Println("============")

	mt, err := p.build(ts)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(mt.Dump(""))
}
