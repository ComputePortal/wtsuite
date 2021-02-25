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
      quotedGroupSettings{
        maskType:        FORMULA,
        groupPattern:    patterns.HTML_SCRIPT_GROUP,
        assertStopMatch: true,
        info:            "script",
        trackStarts:     true,
      },
      quotedGroupSettings{
        maskType:        FORMULA,
        groupPattern:    patterns.HTML_STYLE_GROUP,
        assertStopMatch: true,
        info:            "style",
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

  if err := p.maskFormulas(); err != nil {
    return nil, err
  }

  /*for i, m := range p.mask {
    if m == STRING {
      p.mask[i] = NONE // keep only the formulas
    }
  }*/

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
  sub := &XMLParser{p.refine(start, stop)}

  return sub
}

// used only for attributes
/*func (p *XMLParser) tokenize() ([]raw.Token, error) {
	ts, err := p.Parser.tokenize()
	if err != nil {
		return nil, err
	}

	return p.nestOperators(ts)
}*/

func (p *XMLParser) parseAttributes(ctx context.Context) (*html.RawDict, error) {
	result := html.NewEmptyRawDict(ctx)

	ts, err := p.tokenize()
	if err != nil {
		return result, err
	}

	ts = p.expandTmpGroups(ts)


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
			return result, err
		}

		if (i + 1) < len(ts) {
			switch t := ts[i+1].(type) {
			case *raw.Symbol:
				if _, err := raw.AssertSymbol(t, patterns.EQUAL); err != nil {
					return result, err
				}

				if (i + 2) < len(ts) {
					val := ts[i+2]
					if err := raw.AssertNotSymbol(val); err != nil {
						return result, err
					}

					vs := []raw.Token{val}

					if (i + 3) < len(ts) {
						if raw.IsGroup(ts[i+3]) {
							vs = append(vs, ts[i+3])
							i += 1
						}
					}

					if err := convertAppendKeyVal(key, vs); err != nil {
						return result, err
					}

					i += 3
				} else {
					errCtx := t.Context()
					return result, errCtx.NewError("Syntax Error: expected more")
				}
			case *raw.Word:
				if err := appendFlag(key); err != nil {
					return result, err
				}
				// leave ts[i+1] to next iteration
				i++
			default:
				errCtx := t.Context()
				return result, errCtx.NewError("Syntax Error: bad attribute")
			}
		} else {
			// append a flag
			if err := appendFlag(key); err != nil {
				return result, err
			}

			i++
		}
	}

	return result, nil
}

// script or 
func (p *XMLParser) maskFormula(tagName string) error {
  inSingleQuotes := false
  inTag := false
  inDoubleQuotes := false
  inComment := false

  //re := patterns.TAG_NAME_REGEXP

  formulaStarting := false
  formulaStart := -1

  tmpPos := p.pos
  
  p.pos = 0
  for p.pos < p.Len() {
    c := p.raw[p.pos]

    if inComment {
      if p.pos > 2 && c == '>' && p.raw[p.pos-1] == '-' && p.raw[p.pos-2] == '-' {
        inComment = false
      } 

      p.pos += 1
    } else if !inTag {
      if c == '<' {
        if p.pos < p.Len() - 3 && p.raw[p.pos+1] == '!' && p.raw[p.pos+2] == '-' && p.raw[p.pos+3] == '-' {
          inComment = true
          p.pos += 4
          continue
        }

        inTag = true

        closing := false
        if p.pos < p.Len() - 1 {
          if p.raw[p.pos+1] == '/' {
            closing = true
            p.pos += 1
          } 
        }

        pos := p.pos
        // TODO: replace by simply searching the next whitespace
        rname0 := pos+1
        rname1 := pos+1
        foundName := false
        for true {
          c := p.raw[rname1]
          if (c == ' ') || (c == '\n') || (c == '/') || (c == '>') || (c == '"') || (c == '\'') {
            foundName = true
            p.pos = rname1
            break
          } else {
            rname1 += 1
          }
        }

        /*rname, _, ok := p.nextMatch(re, false)

        if ok {
          if !foundName {
            panic("algo error 1")
          } else if rname[0] != rname0 {
            fmt.Println(rname, rname0)
            panic("algo error 2")
          } else if rname[1] != rname1 {
            nameCtx := p.NewContext(rname[0], rname[1])
            fmt.Println(rname, rname1)
            fmt.Println(nameCtx.NewError("name match").Error())
            panic("algo error 3")
          }
        }*/
        rname, ok := [2]int{rname0, rname1}, foundName

        name := p.Write(pos+1, rname[1])
        if ok && name == tagName {
          if !closing {
            if formulaStart == -1 {
              formulaStarting = true
            } else {
              inTag = false
            }
          } else {
            if formulaStart != -1 {
              p.SetMask(formulaStart, rname[0] - 2, FORMULA)
              formulaStart = -1
            }

            p.nextMatch(patterns.TAG_STOP_REGEXP, false)

            inTag = false
          }
        } else {
          p.pos = pos + 1

          if formulaStart != -1 {
            inTag = false
          }
        }
      } else {
        p.pos += 1
      }
    } else {
      if inSingleQuotes {
        if c == '\'' {
          inSingleQuotes = false
        }
      } else if inDoubleQuotes {
        if c == '"' {
          inDoubleQuotes = false
        }
      } else if c == '>' {
        inTag = false
        if formulaStarting {
          formulaStart = p.pos + 1
          formulaStarting = false
        }
      } else if c == '\'' {
        inSingleQuotes = true
      } else if c == '"' {
        inDoubleQuotes = true
      }

      p.pos += 1
    }
  }

  p.pos = tmpPos

  return nil
}

func (p *XMLParser) maskFormulas() error {
  if err := p.maskFormula("script"); err != nil {
    return err
  }

  if err := p.maskFormula("style"); err != nil {
    return err
  }

  return nil
}

// returns string of end
func (p *XMLParser) findTagEnd(stopSymbol string) ([2]int, string, bool) {
  inSingleQuotes := false
  inDoubleQuotes := false

  nStop := len(stopSymbol)

  isComment := stopSymbol == "-->"

  pos := p.pos
  for ;pos < p.Len(); pos++ {
    c := p.raw[pos]

    if isComment {
      // quotes dont matter inside xml comments
      if c == '>' {
        if pos > nStop && string(p.raw[pos-nStop+1:pos+1]) == stopSymbol {
          p.pos = pos + 1
          return [2]int{pos-nStop+1, pos+1}, stopSymbol, true
        } 
      }
    } else {
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
        if stopSymbol == ">" && p.raw[pos-1] == '/' {
          p.pos = pos+1
          return [2]int{pos-1, pos+1}, "/>", true
        } else if pos > nStop && string(p.raw[pos-nStop+1:pos+1]) == stopSymbol {
          p.pos = pos + 1
          return [2]int{pos-nStop+1, pos+1}, stopSymbol, true
        } else {
          continue
        }
      }
    }
  }

  return [2]int{0, 0,}, "", false
}

// p.pos is advanced if stop is found, otherwise it is unchanged
func (p *XMLParser) findStopTag(tagName string, inScript bool) ([2]int, bool) {
  count := 0
  
  inSingleQuotes := false
  inTag := false
  inDoubleQuotes := false
  inComment := false

  re := patterns.TAG_NAME_REGEXP

  start := p.pos
  for ;p.pos < p.Len(); {
    if p.mask[p.pos] == FORMULA {
      p.pos += 1
      continue
    }

    c := p.raw[p.pos]

    if inComment {
      if p.pos > 2 && c == '>' && p.raw[p.pos-1] == '-' && p.raw[p.pos-2] == '-' {
        inComment = false
      } 

      p.pos += 1
    } else if !inTag {
      if c == '<' {
        if p.pos < p.Len() - 3 && p.raw[p.pos+1] == '!' && p.raw[p.pos+2] == '-' && p.raw[p.pos+3] == '-' {
          inComment = true
          p.pos += 4
          continue
        }

        if !inScript {
          inTag = true
        }

        closing := false
        if p.pos < p.Len() - 1 {
          if p.raw[p.pos+1] == '/' {
            closing = true
            p.pos += 1
          } 
        }

        rname, name, ok := p.nextMatch(re, true)
        if ok && name == tagName {
          if !closing {
            if !inScript {
              count += 1
            }
          } else {

            if count == 0 {
              rend, _, ok := p.nextMatch(patterns.TAG_STOP_REGEXP, false)

              if ok {
                //p.pos = start
                return [2]int{rname[0] - 2, rend[1]}, true
              } else {
                return [2]int{rname[0] - 2, rname[1]}, true
              }
            }

            count -= 1
          }
        }
      } else {
        p.pos += 1
      }
    } else {
      if inSingleQuotes {
        if c == '\'' {
          inSingleQuotes = false
        }
      } else if inDoubleQuotes {
        if c == '"' {
          inDoubleQuotes = false
        }
      } else if c == '>' {
        inTag = false
      } else if c == '\'' {
        inSingleQuotes = true
      } else if c == '"' {
        inDoubleQuotes = true
      }

      p.pos += 1
    }
  }

  p.pos = start

  return [2]int{0, 0}, false
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
				stopSymbol := ">"
				if name == "?xml" {
					stopSymbol = "?>"
				} else if name == "!--" {
          stopSymbol = "-->"
        }

				if rr, s, ok := p.findTagEnd(stopSymbol); ok {
          if name == "!--" || (rname[0] > 0 && p.raw[rname[0]-1] == '/') {
            // skip if tag is comment, or if tag is redundant stop tag
            rprev = rr
            continue
          } 

					ctx := context.MergeContexts(p.NewContext(r[0], rname[1]), p.NewContext(rr[0], rr[1]))

					attrParser := p.Refine(rname[1], rr[0])
          attrParser.Reset()
          if err := attrParser.maskQuoted(); err != nil {
            return nil, err
          }

					attr, err := attrParser.parseAttributes(ctx) // this is where the magic happens
					if err != nil && attr == nil { // some attributes might be correctly parsed
						return nil, err
					}

					rprev = rr

					if name == "script" || name == "style" {
            // single and double quotes need to be matched, during search for stops
            // it is unlikely that the comments comment out the tags
            // the ScriptTagGroup keeps track of the quotes
            //if rrr, ok := p.nextGroupStopMatch(patterns.NewScriptTagGroup(name), true); ok {
            if rrr, ok := p.findStopTag(name, true); ok {

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

              if rrr, ok := p.findStopTag(name, false); ok {
              //if rrr, ok := p.nextGroupStopMatch(patterns.NewTagGroup(name), true); ok {
                ctx = context.MergeContexts(ctx, p.NewContext(rrr[0], rrr[1]))
                subParser = p.Refine(rr[1], rrr[0])
                rprev = rrr
              } else {
                // don't actually throw an error, just ignore
                //return nil, ctx.NewError("Syntax Error: stop tag not found (" + name + ")")
                p.pos = rr[1]
                continue
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
