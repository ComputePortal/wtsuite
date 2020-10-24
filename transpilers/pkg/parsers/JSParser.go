// note:
//  JSParser functions are spread over multiple source files
//  entry point to JSParser is in JSParser_module.go

package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"../tokens/context"
	"../tokens/js"
	"../tokens/patterns"
	"../tokens/raw"
)

func tokenizeJSWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
	switch {
	case patterns.IsHex(s):
		return raw.NewHexLiteralInt(s, ctx)
	case patterns.IsInt(s):
		return raw.NewLiteralInt(s, ctx)
	case patterns.IsFloat(s):
		return raw.NewLiteralFloat(s, ctx)
	case s == "NaN" || s == "Infinity":
		return raw.NewSpecialNumber(s, ctx), nil
	case patterns.IsBool(s):
		return raw.NewLiteralBool(s, ctx)
	case patterns.IsNull(s):
		return raw.NewLiteralNull(ctx), nil
	case s == "in" || s == "instanceof" || s == "typeof" || s == "await" || s == "new":
		return raw.NewSymbol(s, true, ctx), nil
	case patterns.IsJSWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

// for interpolation strings
func tokenizeJSFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	template := raw.NewTemplateGroup(ctx)

	appendString := func(s_ string, ctx_ context.Context) {
		template.Fields = append(template.Fields, []raw.Token{raw.NewValueLiteralString(s_, ctx_)})
	}

	appendTemplate := func(s_ string, ctx_ context.Context) error {
		subP, err := NewRawJSParser(s_, ctx_)
		if err != nil {
			return err
		}

		ts, err := subP.tokenize()
		if err != nil {
			return err
		}

		template.Fields = append(template.Fields, ts)

		return nil
	}

	// parser settings arent really used here
	p_ := newParser(s, ParserSettings{}, ctx)
	p := &p_
	rprev := [2]int{0, 0}
	for true {
		if r, _, ok := p.nextMatch(patterns.JS_STRING_TEMPLATE_START_REGEXP, false); ok {
			appendString(p.Write(rprev[1], r[0]), p.NewContext(r[1], r[0]))

			if rr, _, ok := p.nextMatch(patterns.JS_STRING_TEMPLATE_STOP_REGEXP, false); ok {
				if err := appendTemplate(p.Write(r[1], rr[0]), p.NewContext(r[1], rr[0])); err != nil {
					return nil, err
				}
				rprev = rr
			} else {
				return nil, p.NewError(r[0], r[1], "Syntax Error: string template not closed")
			}
		} else {
			break
		}
	}

	if rprev[1] < p.Len() {
		appendString(p.Write(rprev[1], -1), p.NewContext(rprev[1], -1))
	}

	return []raw.Token{template}, nil
}

var jsParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.JS_STRING_OR_COMMENT_REGEXP,
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
		},
	},
	formulas: formulasSettings{
		//tokenizer: tokenizeJSFormulas, // to avoid circular initialization loop: set later
	},
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.JS_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeJSWordsAndLiterals,
	},
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.JS_SYMBOLS_REGEXP,
	},
	operators: newOperatorsSettings([]operatorSettings{
		operatorSettings{19, "new", PRE},
		operatorSettings{18, "++", POST},
		operatorSettings{18, "--", POST},
		operatorSettings{17, "!", PRE},
		operatorSettings{17, "-", PRE},
		operatorSettings{17, "~", PRE},
		operatorSettings{17, "+", PRE},
		operatorSettings{17, "++", PRE},
		operatorSettings{17, "--", PRE},
		operatorSettings{17, "typeof", PRE},
		//operatorSettings{17, "delete", PRE},
		operatorSettings{17, "await", PRE},
		operatorSettings{16, "**", BIN},
		operatorSettings{15, "/", BIN | L2R},
		operatorSettings{15, "*", BIN | L2R},
		operatorSettings{15, "%", BIN | L2R},
		operatorSettings{14, "+", BIN | L2R},
		operatorSettings{14, "-", BIN | L2R},
		operatorSettings{13, "<<", BIN | L2R},
		operatorSettings{13, ">>", BIN | L2R},
		operatorSettings{13, ">>>", BIN | L2R},
		operatorSettings{12, "<", BIN | L2R},
		operatorSettings{12, "<=", BIN | L2R},
		operatorSettings{12, ">", BIN | L2R},
		operatorSettings{12, ">=", BIN | L2R},
		operatorSettings{12, "in", BIN | L2R},
		operatorSettings{12, "instanceof", BIN | L2R},
		operatorSettings{11, "==", BIN | L2R},
		operatorSettings{11, "!=", BIN | L2R},
		operatorSettings{11, "===", BIN | L2R},
		operatorSettings{11, "!==", BIN | L2R},
		operatorSettings{10, "&", BIN | L2R},
		operatorSettings{9, "^", BIN | L2R},
		operatorSettings{8, "|", BIN | L2R},
		operatorSettings{6, "&&", BIN | L2R},
		operatorSettings{5, "||", BIN | L2R},
		operatorSettings{4, "? :", TER | L2R},
		operatorSettings{3, "=", BIN},
		operatorSettings{3, "+=", BIN},
		operatorSettings{3, "-=", BIN},
		operatorSettings{3, "*=", BIN},
		operatorSettings{3, "/=", BIN},
		operatorSettings{3, "%=", BIN},
		operatorSettings{3, "**=", BIN},
		operatorSettings{3, "<<=", BIN},
		operatorSettings{3, ">>=", BIN},
		operatorSettings{3, ">>>=", BIN},
		operatorSettings{3, "&=", BIN},
		operatorSettings{3, "|=", BIN},
		operatorSettings{3, "^=", BIN},
	}),
	tmpGroupWords:            true,
	tmpGroupPeriods:          true,
	tmpGroupArrows:           true,
	tmpGroupDColons:          true,
	tmpGroupAngled:           true,
	recursivelyNestOperators: false,
}

// to avoid circular initialization loop
func setJSFormulasTokenizer() bool {
	jsParserSettings.formulas.tokenizer = tokenizeJSFormulas

	return true
}

// to avoid circular initialization loop
var jsFormulasOk_ = setJSFormulasTokenizer()

var jsOperatorMap = map[string]string{
	"bin+":          "+",
	"bin-":          "-",
	"bin/":          "/",
	"bin*":          "*",
	"bin%":          "%",
	"bin**":         "**",
	"bin<":          "<",
	"bin>":          ">",
	"bin>=":         ">=",
	"bin<=":         "<=",
	"bin>>":         ">>",
	"bin<<":         "<<",
	"bin>>>":        ">>>",
	"bin|":          "|",
	"bin&":          "&",
	"bin^":          "^",
	"bin.":          ".",
	"binin":         "in",
	"bininstanceof": "instanceof",
	"bin||":         "||",
	"bin&&":         "&&",
	"bin==":         "==",
	"bin!=":         "!=",
	"bin!==":        "!==",
	"bin===":        "===",
	"bin=":          "=",
	"bin+=":         "+=",
	"bin-=":         "-=",
	"bin*=":         "*=",
	"bin/=":         "/=",
	"bin%=":         "%=",
	"bin**=":        "**=",
	"bin<<=":        "<<=",
	"bin>>=":        ">>=",
	"bin>>>=":       ">>>=",
	"bin&=":         "&=",
	"bin|=":         "|=",
	"bin^=":         "^=",
	"post++":        "++",
	"post--":        "--",
	"pre--":         "--",
	"pre++":         "++",
	"pre+":          "+",
	"pre-":          "-",
	"pre~":          "~",
	"pre!":          "!",
	"prenew":        "new",
	"pretypeof":     "typeof",
	"preawait":      "await",
	//"predelete":     "delete",
	"ter? :": "? :",
}

type JSParser struct {
	module *js.ControlModule
	Parser
}

// path is just for context reference
func NewRawJSParser(raw string, ctx context.Context) (*JSParser, error) {
	p := &JSParser{nil, newParser(raw, jsParserSettings, ctx)}

	if err := p.maskQuoted(); err != nil {
		return nil, err
	}

	return p, nil
}

func NewJSParser(path string) (*JSParser, error) {
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

	return NewRawJSParser(raw, ctx)
}

func (p *JSParser) translateOpName(name string) (string, bool) {
	// from raw names to js names
	translatedName, ok := jsOperatorMap[name]
	return translatedName, ok
}

func (p *JSParser) nextSeparatorPosition(ts []raw.Token, sep string) int {
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			return i
		}
	}

	return len(ts)
}

func (p *JSParser) splitByNextSeparator(ts []raw.Token, sep string) ([]raw.Token,
	[]raw.Token) {
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			if i < len(ts)-1 {
				return ts[0:i], ts[i+1:]
			} else {
				return ts[0:i], []raw.Token{}
			}
		}
	}

	return ts, []raw.Token{}
}

func (p *JSParser) splitBySeparator(ts []raw.Token, sep string) [][]raw.Token {
	result := make([][]raw.Token, 0)

	prev := 0
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			result = append(result, ts[prev:i])
			prev = i + 1
		}
	}

	if prev < len(ts) {
		result = append(result, ts[prev:])
	}

	return result
}

func (p *JSParser) stripSeparators(iStart int, ts []raw.Token, symbol string) []raw.Token {
	if iStart > len(ts)-1 {
		return []raw.Token{}
	}

	iRemaining := iStart
	for i := iStart; i < len(ts); i++ {
		if !raw.IsSymbol(ts[i], symbol) {
			iRemaining = i
			break
		}
	}

	return ts[iRemaining:]
}

func (p *JSParser) DumpTokens() {
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

	fmt.Println("\nJS tokens:")
	fmt.Println("============")

	p.Reset()
	m, err := p.BuildModule()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(m.Dump())
}

var _jsTypeExpressionBuilderOk = js.SetTypeExpressionBuilder(func(str string) *js.TypeExpression {
	// tokenize the string
	p, err := NewRawJSParser(str, context.NewDummyContext())
	if err != nil {
		panic(err)
	}

	ts, err := p.tokenize()
	if err != nil {
		panic(err)
	}

	if len(ts) < 1 || len(ts) > 2 {
		panic("bad type string " + str)
	}

	te, err := p.buildTypeExpression(ts, false)
	if err != nil {
		panic(err)
	}

	return te
})
