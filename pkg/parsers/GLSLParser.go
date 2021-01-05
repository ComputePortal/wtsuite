package parsers

import (
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"

  "github.com/computeportal/wtsuite/pkg/tokens/context"
  "github.com/computeportal/wtsuite/pkg/tokens/glsl"
  "github.com/computeportal/wtsuite/pkg/tokens/patterns"
  "github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeGLSLWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
  switch {
  case patterns.IsHex(s):
    return raw.NewHexLiteralInt(s, ctx)
  case patterns.IsInt(s):
    return raw.NewLiteralInt(s, ctx)
  case patterns.IsFloat(s):
    return raw.NewLiteralFloat(s, ctx)
  case patterns.IsBool(s):
    return raw.NewLiteralBool(s, ctx)
  case patterns.IsGLSLWord(s):
    return raw.NewWord(s, ctx)
  default:
    return nil, ctx.NewError("Syntax Error: unparseable")
  }
}

func tokenizeGLSLFormulas(s string, ctx context.Context) ([]raw.Token, error) {
  return nil, ctx.NewError("Error: can't have backtick formula in glsl")
}

var glslParserSettings = ParserSettings{
  quotedGroups: quotedGroupsSettings{
		pattern: patterns.GLSL_STRING_OR_COMMENT_REGEXP,
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
    tokenizer: tokenizeGLSLFormulas,
  },
  wordsAndLiterals: wordsAndLiteralsSettings{
    maskType: WORD_OR_LITERAL,
    pattern: patterns.GLSL_WORD_OR_LITERAL_REGEXP,
    tokenizer: tokenizeGLSLWordsAndLiterals,
  },
  symbols: symbolsSettings{
    maskType: SYMBOL,
    pattern: patterns.GLSL_SYMBOLS_REGEXP,
  },
  operators: newOperatorsSettings([]operatorSettings{
		operatorSettings{17, "!", PRE},
		operatorSettings{17, "-", PRE},
		operatorSettings{17, "+", PRE},
    operatorSettings{15, "/", BIN | L2R},
    operatorSettings{15, "*", BIN | L2R},
    operatorSettings{15, "%", BIN | L2R},
    operatorSettings{14, "+", BIN | L2R},
    operatorSettings{14, "-", BIN | L2R},
		operatorSettings{12, "<", BIN | L2R},
		operatorSettings{12, "<=", BIN | L2R},
		operatorSettings{12, ">", BIN | L2R},
		operatorSettings{12, ">=", BIN | L2R},
		operatorSettings{11, "==", BIN | L2R},
		operatorSettings{11, "!=", BIN | L2R},
		operatorSettings{6, "&&", BIN | L2R},
		operatorSettings{5, "||", BIN | L2R},
  }),
  tmpGroupWords: true,
  tmpGroupPeriods: true,
  tmpGroupArrows: false,
  tmpGroupDColons: false,
  tmpGroupAngled: false,
  recursivelyNestOperators: false,
}

type GLSLParser struct {
  module *glsl.ModuleData
  Parser
}

func NewRawGLSLParser(raw string, ctx context.Context) (*GLSLParser, error) {
  p := &GLSLParser{nil, newParser(raw, glslParserSettings, ctx)}

  if err := p.maskQuoted(); err != nil {
    return nil, err
  }

  return p, nil
}

func NewGLSLParser(path string) (*GLSLParser, error) {
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

  return NewRawGLSLParser(raw, ctx)
}

func (p *GLSLParser) DumpTokens() {
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

	fmt.Println("\nGLSL tokens:")
	fmt.Println("==============")

	p.Reset()
	m, err := p.BuildModule()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Println(m.Dump(""))
}
