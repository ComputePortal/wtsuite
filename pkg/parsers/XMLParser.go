package parsers

// TODO: change this parser so that style and script tags can contain any crap

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeXMLWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
	switch {
	case patterns.IsXMLWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

func tokenizeXMLFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	return nil, ctx.NewError("Error: can't have backtick formula in xml markup")
}

// this is a bad approach, we better just base ourselves on <>
var xmlParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.XML_STRING_REGEXP,
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
		tokenizer: tokenizeXMLFormulas,
	},
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.XML_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeXMLWordsAndLiterals,
	},
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.XML_SYMBOLS_REGEXP,
	},
	operators: newOperatorsSettings([]operatorSettings{}),
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
  tokenizeWhitespace:       false,
}

type XMLParser struct {
	Parser
}

func NewXMLParserFromBytes(rawBytes []byte, path string) (*XMLParser, error) {
  raw := string(rawBytes)

	src := context.NewSource(raw)

	ctx := context.NewContext(src, path)
	p := &XMLParser{newParser(raw, xmlParserSettings, ctx)}

  // dont mask anything yet, because style/script/text tags can contain anything

	return p, nil
}

// can be a url, in which case it is fetched
func NewXMLParser(path string) (*XMLParser, error) {
  if !filepath.IsAbs(path) {
    panic("path should be absolute")
  }

  rawBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  return NewXMLParserFromBytes(rawBytes, path)
}

func NewEmptyXMLParser(ctx context.Context) *XMLParser {
	return &XMLParser{newParser("", xmlParserSettings, ctx)}
}

func (p *XMLParser) Refine(start, stop int) *XMLParser {
	return &XMLParser{p.refine(start, stop)}
}

// used only for attributes
func (p *XMLParser) tokenize() ([]raw.Token, error) {
	ts, err := p.Parser.tokenize()
	if err != nil {
		return nil, err
	}

	return p.nestOperators(ts)
}

func (p *XMLParser) parseAttributes(ctx context.Context) (*html.RawDict, error) {
	ts, err := p.tokenize()
	if err != nil {
		return nil, err
	}

	ts = p.expandTmpGroups(ts)

	result := html.NewEmptyRawDict(ctx)

	appendKeyVal := func(k *raw.Word, v html.Token) error {
		if other, otherValue_, ok := result.GetKeyValue(k.Value()); ok {
      // duplicate is not a problem, just extend

      if otherValue, ok := otherValue_.(*html.String); ok {
        if vStr, okV := v.(*html.String); okV  {
          s, _ := html.NewString(k.Value(), k.Context())
          result.Set(s, html.NewValueString(otherValue.Value() + " " + vStr.Value(), otherValue.Context()))
          return nil
        } 
      }

      errCtx := context.MergeContexts(k.Context(), other.Context())
      return errCtx.NewError("Error: duplicate (" + k.Value() + ")")
		} else {
			s, _ := html.NewString(k.Value(), k.Context())
			result.Set(s, v)
			return nil
		}
	}

	appendFlag := func(k *raw.Word) error {
		return appendKeyVal(k, html.NewFlag(k.Context()))
	}

  // TODO: only accept strings
	convertAppendKeyVal := func(k *raw.Word, vs []raw.Token) error {
    if len(vs) == 0{
      errCtx := k.Context()
      return errCtx.NewError("Error: expected value after attribute key")
    } else if len(vs) > 1 {
      errCtx := raw.MergeContexts(vs[1:]...)
      return errCtx.NewError("Error: unexpected value tokens")
    }

    v, err := raw.AssertLiteralString(vs[0])
    if err != nil {
      return err
    }

		return appendKeyVal(k, html.NewValueString(v.Value(), v.Context()))
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

// returns string of end
func (p *XMLParser) findTagEnd() ([2]int, string, bool) {
  inSingleQuotes := false
  inDoubleQuotes := false

  pos := p.pos
  for ;pos < p.Len(); pos++ {
    c := p.raw[pos]
    if inDoubleQuotes {
      if c == '"' {
        inDoubleQuotes = false
      }
    } else if inSingleQuotes {
      if c == '\'' {
        inSingleQuotes = false
      }
    } else if c == '"' {
      inDoubleQuotes = true
    } else if c == '\'' {
      inSingleQuotes = true
    } else if c == '>' {

      if p.raw[pos-1] == '/' {
        return [2]int{pos-1, pos+1}, "/>", true
      } else {
        return [2]int{pos, pos+1}, ">", true
      }
    }
  }

  return [2]int{0, 0,}, "", false
}

func (p *XMLParser) BuildTags() ([]*html.Tag, error) {
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

			if rname, name, ok := p.nextMatch(patterns.TAG_NAME_REGEXP, false); ok {
				/*stopRegexp := patterns.TAG_STOP_REGEXP
				if name == "?xml" {
					stopRegexp = patterns.XML_HEADER_STOP_REGEXP
				} else if name == "!--" {
          stopRegexp = patterns.XML_COMMENT_STOP_REGEXP
        }*/

				if rr, s, ok := p.findTagEnd(); ok {
          if name == "!--" {
            rprev = rr
            continue
          } 

					ctx := context.MergeContexts(p.NewContext(r[0], rname[1]), p.NewContext(rr[0], rr[1]))

					attrParser := p.Refine(rname[1], rr[0])
          if err := attrParser.maskQuoted(); err != nil {
            return nil, err
          }

					attr, err := attrParser.parseAttributes(ctx) // this is where the magic happens
					if err != nil {
						return nil, err
					}

					rprev = rr

					if name == "script" || name == "style" {
            // single and double quotes need to be matched, during search for stops
            // it is unlikely that the comments comment out the tags
            // the ScriptTagGroup keeps track of the quotes
            if rrr, ok := p.nextGroupStopMatch(patterns.NewScriptTagGroup(name), true); ok {

              ctx = context.MergeContexts(ctx, p.NewContext(rrr[0], rrr[1]))
              subParser := p.Refine(rr[1], rrr[0])
              rprev = rrr
              subTag := html.NewScriptTag(strings.ToLower(name), attr, subParser.Write(0, -1),
                subParser.NewContext(0, -1), ctx)
              appendTag(subTag)
            } else {
              return nil, ctx.NewError("Syntax Error: unmatched script/style tag (" + name + ")")
            }
          } else {
            var subParser *XMLParser = nil
            if patterns.IsSelfClosing(name, s) {
              subParser = p.Refine(rr[1], rr[1])
            } else {
              if name == "!--" {
                panic("shouldn't get here")
              }

              if rrr, ok := p.nextGroupStopMatch(patterns.NewTagGroup(name), true); ok {
                ctx = context.MergeContexts(ctx, p.NewContext(rrr[0], rrr[1]))
                subParser = p.Refine(rr[1], rrr[0])
                rprev = rrr
              } else {
                return nil, ctx.NewError("Syntax Error: unmatched tag (" + name + ")")
              }
            }

            subTags, err := subParser.BuildTags()
            if err != nil {
              return nil, err
            }

            subTag := html.NewTag(strings.ToLower(name), attr, subTags, ctx)
            appendTag(subTag)
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

func (p *XMLParser) DumpTokens() {
	fmt.Println("\nXML tokens:")
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
