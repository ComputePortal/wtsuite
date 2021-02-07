package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeCSSWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
	switch {
	case patterns.IsCSSWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

// tokenizer for CSS selectors
var cssParserSettings = ParserSettings{
  quotedGroups: quotedGroupsSettings{
    pattern: patterns.CSS_STRING_REGEXP,
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
    },
  },
  formulas: formulasSettings{
    tokenizer: nil,
  },
  wordsAndLiterals: wordsAndLiteralsSettings{
    maskType: WORD_OR_LITERAL,
    pattern: patterns.CSS_WORD_OR_LITERAL_REGEXP,
    tokenizer: tokenizeCSSWordsAndLiterals,
  },
  symbols: symbolsSettings{
    maskType: SYMBOL,
    pattern: patterns.CSS_SYMBOLS_REGEXP,
  },
  operators: newOperatorsSettings([]operatorSettings{}), // irrelevant because we only tokenize and nestGroups
  tmpGroupWords:   false,
  tmpGroupPeriods:  false,
  tmpGroupArrows:  false,
  tmpGroupDColons: false,
  tmpGroupAngled:  false,
  recursivelyNestOperators: false,
  tokenizeWhitespace: false,
}

type CSSSelectorParser struct {
  Parser
}

func NewCSSSelectorParser(s string, ctx context.Context) (*CSSSelectorParser, error) {
  p := &CSSSelectorParser{newParser(s, cssParserSettings, ctx)}

  if err := p.maskQuoted(); err != nil {
    return nil, err
  }

  return p, nil
}

func (p *CSSSelectorParser) Tokenize() ([]raw.Token, error) {
  // also nests groupds
  return p.tokenize()
}
