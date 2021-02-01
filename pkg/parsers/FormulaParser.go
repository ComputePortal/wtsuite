package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeFormulaWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
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

func tokenizeFormulaFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	return nil, ctx.NewError("Error: can't have formula within formula")
}

var formulaParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.FORMULA_STRING_OR_COMMENT_REGEXP,
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
				maskType:        SL_COMMENT,
				groupPattern:    patterns.SL_COMMENT_GROUP,
				assertStopMatch: false,
				info:            "single-line comment",
				trackStarts:     false,
			},
			quotedGroupSettings{
				maskType:        ML_COMMENT,
				groupPattern:    patterns.ML_COMMENT_GROUP,
				assertStopMatch: true,
				info:            "js-style multiline comment",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        ML_COMMENT,
				groupPattern:    patterns.XML_COMMENT_GROUP,
				assertStopMatch: true,
				info:            "xml-style multiline comment",
				trackStarts:     true,
			},
		},
	},
	formulas: formulasSettings{
		tokenizer: tokenizeFormulaFormulas,
	},
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.FORMULA_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeFormulaWordsAndLiterals,
	},
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.FORMULA_SYMBOLS_REGEXP,
	},
	operators: newOperatorsSettings([]operatorSettings{
    operatorSettings{17, "$", PRE},
		operatorSettings{16, "-", PRE},
		operatorSettings{16, "!", PRE},
		operatorSettings{14, "/", BIN | L2R},
		operatorSettings{14, "*", BIN | L2R},
		operatorSettings{13, "-", BIN | L2R},
		operatorSettings{13, "+", BIN | L2R},
		operatorSettings{11, "<", BIN | L2R},
		operatorSettings{11, "<=", BIN | L2R},
		operatorSettings{11, ">", BIN | L2R},
		operatorSettings{11, ">=", BIN | L2R},
		operatorSettings{10, "!=", BIN | L2R},
		operatorSettings{10, "==", BIN | L2R},
		operatorSettings{10, "===", BIN | L2R},
		operatorSettings{8, "&&", BIN | L2R},
		operatorSettings{7, "||", BIN | L2R},
		operatorSettings{6, "!!", BIN | L2R},
    operatorSettings{5, "??", BIN | L2R},
		operatorSettings{4, ":=", BIN}, // so we can use new in ternary operators
		operatorSettings{3, "?", BIN},  // so we can use ternary operator inside dicts
		operatorSettings{2, ":", SING | PRE | POST | BIN},
		operatorSettings{1, "=", BIN},
	}),
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
  tokenizeWhitespace:       false,
}

type FormulaParser struct {
	Parser
}

func NewFormulaParser(s string, ctx context.Context) (*FormulaParser, error) {
	p := &FormulaParser{newParser(s, formulaParserSettings, ctx)}

	if err := p.maskQuoted(); err != nil {
		return nil, err
	}

	for i, m := range p.mask {
		if m == FORMULA {
			errCtx := context.MergeContexts(p.NewContext(0, 1), p.NewContext(i, i+1))
			return nil, errCtx.NewError("Error: formula can't contain formula")
		}
	}
	return p, nil
}

func (p *FormulaParser) tokenize() ([]raw.Token, error) {
	ts, err := p.Parser.tokenize()
	if err != nil {
		return nil, err
	}

	return p.nestOperators(ts)
}
