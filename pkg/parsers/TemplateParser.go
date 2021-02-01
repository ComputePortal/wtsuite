package parsers

import (
  "errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeTemplateFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	return nil, ctx.NewError("Error: can't have backtick formula in ui markup")
}

func tokenizeTemplateWordsAndLiterals(s string, ctx context.Context) (raw.Token, error) {
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
	case patterns.IsTemplateWord(s):
		return raw.NewWord(s, ctx)
	default:
		return nil, ctx.NewError("Syntax Error: unparseable")
	}
}

var uiParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.TEMPLATE_STRING_OR_COMMENT_REGEXP,
		groups: []quotedGroupSettings{
			quotedGroupSettings{
				maskType:        STRING,
				groupPattern:    patterns.SQ_STRING_GROUP,
				assertStopMatch: false,
				info:            "single quoted",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        STRING,
				groupPattern:    patterns.DQ_STRING_GROUP,
				assertStopMatch: false,
				info:            "double quoted",
				trackStarts:     true,
			},
			quotedGroupSettings{
				maskType:        STRING,
				groupPattern:    patterns.BT_FORMULA_GROUP,
				assertStopMatch: false,
				info:            "backtick quoted",
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
		},
	},
	formulas: formulasSettings{
		tokenizer: tokenizeTemplateFormulas,
	},
	// same as html
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.TEMPLATE_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeTemplateWordsAndLiterals,
	},
	// handled by FormulaParser
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  patterns.TEMPLATE_SYMBOLS_REGEXP,
	},
	operators:                newOperatorsSettings([]operatorSettings{}), // handled by FormulaParser
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
  tokenizeWhitespace:       true,
}

type TemplateParser struct {
	helper *HTMLParser
	Parser
}

func NewTemplateParser(path string) (*TemplateParser, error) {
	rawBytes, err := ioutil.ReadFile(path)
	if err != nil {
    return nil, errors.New("Error: problem reading \"" + path + "\" (" + err.Error() + ")")
	}

	raw := string(rawBytes)
	src := context.NewSource(raw)

	ctx := context.NewContext(src, path)
	p := &TemplateParser{
		NewEmptyHTMLParser(context.NewDummyContext()),
		newParser(raw, uiParserSettings, ctx),
	}

	if err := p.maskQuoted(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *TemplateParser) BuildTags() ([]*html.Tag, error) {
  ts, err := p.tokenizeFlat()
  if err != nil {
    return nil, err
  }

	result := make([]*html.Tag, 0)

  var tag *html.Tag
  tag, ts, err = p.buildPermissiveDirective(ts)
  if err != nil {
    return nil, err
  }

  if tag != nil {
    result = append(result, tag)
  }

	indents := make([]int, 0)
	stack := make([]*html.Tag, 0)

	pushStack := func(t *html.Tag, indent int) error {
		if len(result) == 0 {
			errCtx := t.Context()
			return errCtx.NewError("Internal Error: cannot increase indentation without previous tags")
		}

		stack = append(stack, t)
		indents = append(indents, indent)

		return nil
	}

	appendTagInner := func(t *html.Tag) error {
		n := len(stack)
		if n == 0 {
			result = append(result, t)
		} else if n != 0 {
			if err := stack[n-1].AppendChild(t); err != nil {
				return err
			}
		}

		return nil
	}

	// automatically create the ifelse tag at the same indent as an if
	//  dont pop that ifelse for any subsequent else or elseif tags at the same indent
	popStack := func(t *html.Tag, indent int) error {
		inBranch := t.Name() == "else" || t.Name() == "elseif"
		for i, _ := range stack {
			if indents[i] >= indent {
				// dont pop ifelse on same indent
				if !(inBranch && indents[i] == indent && stack[i].Name() == ".ifelse") {
					stack = stack[0:i]
					indents = indents[0:i]
					break
				}
			}
		}

		if t.Name() == "if" && (len(stack) == 0 || stack[len(stack)-1].Name() != ".ifelse") {
			ifElseCtx := t.Context()
			ifElseTag := html.NewDirectiveTag(".ifelse", html.NewEmptyRawDict(ifElseCtx),
				[]*html.Tag{}, ifElseCtx)

			if err := appendTagInner(ifElseTag); err != nil {
				return err
			}

			if err := pushStack(ifElseTag, indent); err != nil {
				return err
			}
		}

		return nil
	}

	appendTag := func(t *html.Tag, indent int) error {
		if t == nil {
			panic("tag is nil")
		}

		if err := popStack(t, indent); err != nil {
			return err
		}

		if err := appendTagInner(t); err != nil {
			return err
		}

		if err := pushStack(t, indent); err != nil {
			return err
		}

		return nil
	}

	// start at col 0 on an empty line
	for len(ts) > 0 {
    var indent int
		ts, indent = p.eatWhitespace(ts)
    if len(ts) > 0 {
      tag, ts, err = p.buildTag(indent, ts)
      if err != nil {
        return nil, err
      }

      if tag.Name() == "parameters" {
        if len(result) > 0 {
          if result[0].Name() == "permissive" {
            if len(result) > 1 {
              errCtx := tag.Context()
              return nil, errCtx.NewError("Error: parameters directive must come first in file (after permissive though)")
            }
          } else {
            errCtx := tag.Context()
            return nil, errCtx.NewError("Error: parameters directive must come first in file (" + result[0].Name() + " is first)")
          }
        }
      }

      if err := appendTag(tag, indent); err != nil {
        return nil, err
      }
    }
	}

	return result, nil
}

func (p *TemplateParser) buildTextTag(inline bool, ts []raw.Token) (*html.Tag, []raw.Token, error) {
  ctx := ts[0].Context()

  var expr_ html.Token
  var rem []raw.Token
  var err error
  if !inline {
    expr_, rem, err = p.buildExpression(ts)
    if err != nil {
      return nil, nil, err
    }
  } else {
    expr_, rem, err = p.buildTextTagExpression(ts)
    if err != nil {
      return nil, nil, err
    }
  }

  if html.IsString(expr_) {
    expr, err := html.AssertString(expr_)
    if err != nil {
      panic(err)
    }

    str := expr.Value()

    //if len(rem) != len(ts) - 1 {

      //panic("algo error")
    //}

    return html.NewTextTag(str, ctx), rem, nil
  } else {
    // elaborate print tag
    attr := html.NewEmptyRawDict(ctx)
    attr.Set(html.NewValueInt(0, ctx), expr_)

    return html.NewDirectiveTag(".print", attr, []*html.Tag{}, ctx), rem, nil
  }
}

func (p *TemplateParser) buildImportExportNames(ts []raw.Token) (*html.RawDict, []raw.Token, error) {
  ctx := ts[0].Context()
  if r, ok := raw.FindGroupStop(ts, 1, ts[0]); ok {
    groups, err := p.nestGroups(p.removeWhitespace(ts[0:r[1]+1]))
    if err != nil {
      return nil, nil, err
    }

    if len(groups) != 1 {
      errCtx := raw.MergeContexts(ts...)
      return nil, nil, errCtx.NewError("Error: unexpected")
    }

    group, err := raw.AssertBracesGroup(groups[0])
    if err != nil {
      return nil, nil, err
    }

    if len(group.Fields) == 0 {
      errCtx := group.Context()
      return nil, nil, errCtx.NewError("Error: empty brace")
    }

    names := html.NewEmptyRawDict(ctx)

    for _, field := range group.Fields {
      newName_ := field[len(field)-1]

      newName, err := raw.AssertWord(newName_)
      if err != nil {
        return nil, nil, err
      }

      newNameCtx := newName.Context()
      if !patterns.IsValidVar(newName.Value()) {
        errCtx := newNameCtx
        return nil, nil, errCtx.NewError("Error: not a valid var")
      }

      if len(field) == 1 {
        names.Set(
          html.NewValueString(newName.Value(), newNameCtx), 
          html.NewValueString(newName.Value(), newNameCtx))
        continue
      }

      if len(field) == 2 {
        errCtx := field[0].Context()
        if raw.IsWord(field[0], "as") {
          return nil, nil, errCtx.NewError("Error: expected more before as") 
        } else {
          return nil, nil, errCtx.NewError("Error: unexpected")
        }
      }

      if !raw.IsWord(field[len(field)-2], "as") {
        errCtx := field[len(field)-2].Context()
        return nil, nil, errCtx.NewError("Error: expected \"as\"")
      }

      oldNameExpr, rem, err := p.buildExpression(field[0:len(field)-2])
      if err != nil {
        return nil, nil, err
      }

      if len(rem) != 0 {
        errCtx := raw.MergeContexts(rem...)
        return nil, nil, errCtx.NewError("Error: unexpected tokens")
      }

      names.Set(oldNameExpr, html.NewValueString(newName.Value(), newNameCtx))
    }

    return names, ts[r[1]+1:], nil
  } else {
    errCtx := ctx
    return nil, nil, errCtx.NewError("Error: closing brace not found")
  }
}

// valid forms:
// import * as namespace from expression...
// export * from expression...
// import *(parameters) as namespace from expression...
// export *(parameters) from expression...
// import {name1, name2, name3 as alias3} from expression...
// export {name1, name2, name3 as alias3} from expression...
// export {name1, name2, name3 as alias3}
// import {name1, name2, name3 as alias3}(parameters) from expression...
// export {name1, name2, name3 as alias3}(parameters) from expression...
func (p *TemplateParser) buildImportExportDirective(dynamic bool, ts []raw.Token) (*html.Tag, []raw.Token, error) {
  nameToken, err := raw.AssertWord(ts[0])
  if err != nil {
    return nil, nil, err
  }
  name := nameToken.Value()
  ctx := ts[0].Context()

  attr := html.NewEmptyRawDict(ctx)
	attr.Set(html.NewValueString(".dynamic", ctx), html.NewValueBool(dynamic, ctx))

  if len(ts) < 4 {
    errCtx := ts[len(ts)-1].Context()
    return nil, nil, errCtx.NewError("Error: expected more tokens after")
  }

  ts, _ = p.eatWhitespace(ts[1:])

  names := html.NewEmptyRawDict(ts[0].Context())
  if raw.IsSymbol(ts[0], "*") {
    starCtx := ts[0].Context()
    ts, _ = p.eatWhitespace(ts[1:])
    if raw.IsSymbol(ts[0], patterns.PARENS_START) {
      // parameters
      var parameters *html.Parens
      parameters, ts, err = p.buildParens(ts)
      if err != nil {
        return nil, nil, err
      }

      attr.Set(html.NewValueString("parameters", ts[0].Context()), parameters)
    }

    if raw.IsWord(ts[0], "as") {
      if name == "export" {
        errCtx := ts[0].Context()
        return nil, nil, errCtx.NewError("Error: can't export nested namespace")
      }

      ts, _ = p.eatWhitespace(ts[1:])

      nameToken, err := raw.AssertWord(ts[0])
      if err != nil {
        return nil, nil, err
      }

      if !patterns.IsValidVar(nameToken.Value()) {
        errCtx := nameToken.Context()
        return nil, nil, errCtx.NewError("Error: bad namespace name")
      }

      names.Set(
        html.NewValueString("*", starCtx),
        html.NewValueString(nameToken.Value(), nameToken.Context()))

      ts, _ = p.eatWhitespace(ts[1:])
    } else {
      names.Set(
        html.NewValueString("*", starCtx),
        html.NewValueString("*", starCtx))
    }
  } else if raw.IsSymbol(ts[0], patterns.BRACES_START) {
    names, ts, err = p.buildImportExportNames(ts)
    if err != nil {
      return nil, nil, err
    }

    if raw.IsSymbol(ts[0], patterns.PARENS_START) {
      // parameters
      var parameters *html.Parens
      parameters, ts, err = p.buildParens(ts)
      if err != nil {
        return nil, nil, err
      }

      attr.Set(html.NewValueString("parameters", ts[0].Context()), parameters)
    }
  } else {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: bad " + name + " directive")
  }

  if !raw.IsWord(ts[0], "from") {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: bad " + name + " directive (expected \"from\")")
  }

  fromCtx := ts[0].Context()
  ts, _ = p.eatWhitespace(ts[1:])

  srcExpr, rem, err := p.buildExpression(ts)
  if err != nil {
    return nil, nil, err
  }

  attr.Set(html.NewValueString("names", names.Context()), names)
  attr.Set(html.NewValueString("from", fromCtx), srcExpr)
  
	return html.NewDirectiveTag(name, attr, []*html.Tag{}, ctx), rem, nil
}

// everything following the function keyword
// eat until first closing brace (containerCount == 0
func (p *TemplateParser) buildFunctionDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  ctx := ts[0].Context()

  if len(ts) < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expeced more tokens")
  }

  nameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, nil, err
  }

  // act as an expresion
  keywordToken := raw.NewValueWord("function", ctx)

  ts, _ = p.eatWhitespace(ts[2:])

  ts = append([]raw.Token{keywordToken}, ts...)

	fnValue, rem, err := p.buildExpression(ts)
	if err != nil {
		return nil, nil, err
	}

	varAttr := html.NewEmptyRawDict(ctx)

	nameHtmlToken := html.NewValueString(nameToken.Value(), nameToken.Context())
	varAttr.Set(nameHtmlToken, fnValue)

	return html.NewDirectiveTag("var", varAttr, []*html.Tag{}, ctx), rem, nil
}

func (p *TemplateParser) buildVarDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  ctx := ts[0].Context()

  if len(ts) < 4 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected more tokens for var directive")
  }

  nameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, nil, err
  }

  ts, _ = p.eatWhitespace(ts[2:])

  if !raw.IsSymbol(ts[0], "=") {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected =")
  }

  ts, _ = p.eatWhitespace(ts[1:])

  rhsExpr, rem, err := p.buildExpression(ts)
  if err != nil {
    return nil, nil, err
  }

	attr := html.NewEmptyRawDict(nameToken.Context())

	nameHtmlToken := html.NewValueString(nameToken.Value(), nameToken.Context())
	attr.Set(nameHtmlToken, rhsExpr)

	return html.NewDirectiveTag("var", attr, []*html.Tag{}, ctx), rem, nil
}

func (p *TemplateParser) buildPermissiveDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  for ;len(ts) > 0 && raw.IsWhitespace(ts[0]); {
    ts = ts[1:]
  }

  if len(ts) == 0 {
    return nil, ts, nil
  }

  if raw.IsWord(ts[0], "permissive") {
    iNext := 1
    if len(ts) > 1 {
      if !raw.IsNL(ts[1]) {
        errCtx := ts[1].Context()
        return nil, nil, errCtx.NewError("Error: unexpected token")
      } else {
        iNext += 1
      }
    }

    ctx := ts[0].Context()
    attr := html.NewEmptyRawDict(ctx)

    tag := html.NewDirectiveTag("permissive", attr, []*html.Tag{}, ctx)

    return tag, ts[iNext:], nil
  } else {
    return nil, ts, nil
  }
}

func (p *TemplateParser) buildParens(ts []raw.Token) (*html.Parens, []raw.Token, error) {
  if r, ok := raw.FindGroupStop(ts, 1, ts[0]); ok {
    groups, err := p.nestGroups(p.removeWhitespace(ts[0:r[1]+1]))
    if err != nil {
      return nil, nil, err
    }

    if len(groups) != 1 {
      errCtx := raw.MergeContexts(ts...)
      return nil, nil, errCtx.NewError("Error: unexpected")
    }

    group, err := raw.AssertParensGroup(groups[0])
    if err != nil {
      return nil, nil, err
    }

    ctx := group.Context()

    if group.IsSemiColon() {
      errCtx := group.Context()
      return nil, nil, errCtx.NewError("Error: expected comma separator, got semicolon")
    }

    n := len(group.Fields)
    values := make([]html.Token, n)
    alts := make([]html.Token, n)

    for i, field := range group.Fields {
      if len(field) == 1 {
        valExpr, rem, err := p.buildExpression(field)
        if err != nil {
          return nil, nil, err
        }

        if len(rem) != 0 {
          errCtx := raw.MergeContexts(rem...)
          return nil, nil, errCtx.NewError("Error: unexpected tokens")
        }

        values[i] = valExpr
        alts[i] = nil
      } else if len(field) > 2 && raw.IsAnyWord(field[0]) && raw.IsSymbol(field[1], "=") {
        keyToken, err := raw.AssertWord(field[0])
        if err != nil {
          return nil, nil, err
        }

        values[i] = html.NewValueString(keyToken.Value(), keyToken.Context())

        altExpr, rem, err := p.buildExpression(field[2:])
        if err != nil {
          return nil, nil, err
        }

        if len(rem) != 0 {
          errCtx := raw.MergeContexts(rem...)
          return nil, nil, errCtx.NewError("Error: unexpected tokens")
        }

        alts[i] = altExpr
      } else if len(field) > 2 && raw.IsAnyWord(field[0]) && raw.IsSymbol(field[1], "!=") {
        keyToken, err := raw.AssertWord(field[0])
        if err != nil {
          return nil, nil, err
        }

        values[i] = html.NewValueString(keyToken.Value() + "!", keyToken.Context())

        altExpr, rem, err := p.buildExpression(field[2:])
        if err != nil {
          return nil, nil, err
        }

        if len(rem) != 0 {
          errCtx := raw.MergeContexts(rem...)
          return nil, nil, errCtx.NewError("Error: unexpected tokens")
        }

        alts[i] = altExpr
      } else {
        valExpr, rem, err := p.buildExpression(field)
        if err != nil {
          return nil, nil, err
        }

        if len(rem) != 0 {
          errCtx := raw.MergeContexts(rem...)
          return nil, nil, errCtx.NewError("Error: unexpected tokens")
        }

        values[i] = valExpr
        alts[i] = nil
      }
    }

    if r[1] + 1 < len(ts) {
      return html.NewParens(values, alts, ctx), ts[r[1]+1:], nil
    } else {
      return html.NewParens(values, alts, ctx), []raw.Token{}, nil
    }
  } else {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: closing parens not found")
  }
}

func (p *TemplateParser) buildParametersDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  parens, rem, err := p.buildParens(ts[1:])
  if err != nil {
    return nil, nil, err
  }

  ctx := parens.Context()
  attr := html.NewEmptyRawDict(ctx)
  attr.Set(html.NewValueString("parameters", ctx), parens)

  tag := html.NewTag("parameters", attr, []*html.Tag{}, ctx)

  return tag, rem, nil
}

func (p *TemplateParser) buildTemplateDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  if !raw.IsWord(ts[0], "template") {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected template keyword")
  }

  ctx := ts[0].Context()

  nameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, nil, err
  }

  nameCtx := nameToken.Context()
	nameKey := html.NewValueString("name", nameCtx)
	nameVal := html.NewValueString(nameToken.Value(), nameCtx)

  attr := html.NewEmptyRawDict(ctx)

  // name was eaten before
  attr.Set(nameKey, nameVal);

  ts, _ = p.eatWhitespace(ts[2:])

  iSuper := -1
  for i, t := range ts {
    if raw.IsWord(t, "super") {
      iSuper = i
      break
    }
  }

  if iSuper == -1 {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: super keyword not found")
  }

  superCtx := ts[iSuper].Context()
  superParens, rem, err := p.buildParens(ts[iSuper+1:])
  if err != nil {
    return nil, nil, err
  }

  superAttr := superParens.ToRawDict()

  ts = ts[0:iSuper]
  if raw.IsSymbol(ts[0], patterns.PARENS_START) {
    // args
    argParens, argRem, err := p.buildParens(ts)
    if err != nil {
      return nil, nil, err
    }

    attr.Set(html.NewValueString("args", argParens.Context()), argParens)

    ts, _ = p.eatWhitespace(argRem)
  }

  if !raw.IsWord(ts[0], "extends") {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected extends keyword")
  }

  extendsCtx := ts[0].Context()

  extendsExpr, extendsRem, err := p.buildExpression(ts[1:])
  if err != nil {
    return nil, nil, err
  }

  if len(extendsRem) != 0 {
    errCtx := raw.MergeContexts(extendsRem...)
    return nil, nil, errCtx.NewError("Error: unexpected tokens")
  }


  attr.Set(html.NewValueString("extends", extendsCtx), extendsExpr)
  attr.Set(html.NewValueString("super", superCtx), superAttr)

	return html.NewDirectiveTag("template", attr, []*html.Tag{}, ctx), rem, nil
}

func (p *TemplateParser) buildForDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  ctx := ts[0].Context()
  attr := html.NewEmptyRawDict(ctx)

  vNameToken, err := raw.AssertWord(ts[1])
  if err != nil {
    return nil, nil, err
  }

  if raw.IsSymbol(ts[2], patterns.COMMA) {
    iNameToken := vNameToken

    attr.Set(html.NewValueString("iname", ts[1].Context()), html.NewValueString(iNameToken.Value(), iNameToken.Context()))

    vNameToken, err = raw.AssertWord(ts[3])
    if err != nil {
      return nil, nil, err
    }

    ts, _ = p.eatWhitespace(ts[4:])
  } else {
    ts, _ = p.eatWhitespace(ts[2:])
  }

  attr.Set(html.NewValueString("vname", vNameToken.Context()), html.NewValueString(vNameToken.Value(), vNameToken.Context()))

  if !raw.IsWord(ts[0], "in") {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected \"in\"")
  }

  rhsExpr, rem, err := p.buildExpression(ts[1:])
  if err != nil {
    return nil, nil, err
  }

  attr.Set(html.NewValueString("in", ts[0].Context()), rhsExpr)

	return html.NewDirectiveTag("for", attr, []*html.Tag{}, ctx), rem, nil
}

func (p *TemplateParser) buildSingleOrNoValueDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  firstToken, err := raw.AssertWord(ts[0])
  if err != nil {
    return nil, nil, err
  }

  ctx := firstToken.Context()
  attr := html.NewEmptyRawDict(ctx)

  if len(ts) > 1 {
    foundNL := false
    for i, t := range ts[1:] {
      if raw.IsNL(t) {
        if i > 0 {
          rhsExpr, rem, err := p.buildExpression(ts[1:])
          if err != nil {
            return nil, nil, err
          }

          attr.Set(html.NewValueInt(0, ctx), rhsExpr)
          return html.NewDirectiveTag(firstToken.Value(), attr, []*html.Tag{}, ctx), rem, nil
        }

        foundNL = true
        break
      }
    }

    if !foundNL {
      errCtx := ctx
      return nil, nil, errCtx.NewError("Error: no terminating nl found")
    }

    return html.NewDirectiveTag(firstToken.Value(), attr, []*html.Tag{}, ctx), ts[1:], nil
  } else {
    return html.NewDirectiveTag(firstToken.Value(), attr, []*html.Tag{}, ctx), []raw.Token{}, nil
  }
}

func (p *TemplateParser) buildExportedDirective(ts []raw.Token) (*html.Tag, []raw.Token, error) {
  if len(ts) < 2 {
    errCtx := ts[0].Context()
    return nil, nil, errCtx.NewError("Error: expected more tokens")
  }

	addExportAttr := func(attr *html.RawDict, exportCtx context.Context) {
		exportToken := html.NewValueString("export", exportCtx)
		flagToken := html.NewValueString("", exportCtx)
		attr.Set(exportToken, flagToken)
	}

  exportCtx := ts[0].Context()

  switch {
  case raw.IsWord(ts[1], "var"):
    tag, rem, err := p.buildVarDirective(ts[1:])
    if err != nil {
      return nil, nil, err
    }

    addExportAttr(tag.RawAttributes(), exportCtx)

    return tag, rem, nil
  case raw.IsWord(ts[1], "template"):
    tag, rem, err := p.buildTemplateDirective(ts[1:])
    if err != nil {
      return nil, nil, err
    }

    addExportAttr(tag.RawAttributes(), exportCtx)

    // add export flag to tag
    //return html.NewTag(name, attr, []*html.Tag{}, tagCtx), nil
    return tag, rem, nil
  case raw.IsWord(ts[1], "function"):
    tag, rem, err := p.buildFunctionDirective(ts[1:])
    if err != nil {
      return nil, nil, err
    }

    addExportAttr(tag.RawAttributes(), exportCtx)

    return tag, rem, nil
  default:
    if !raw.IsSymbol(ts[1], patterns.BRACES_START) && !raw.IsSymbol(ts[1], "*") {
      errCtx := ts[1].Context()
      return nil, nil, errCtx.NewError("Error: invalid export statement")
    }

    return p.buildImportExportDirective(false, ts)
  }
}

func (p *TemplateParser) buildGenericTag(inline bool, ts []raw.Token) (*html.Tag, []raw.Token, error) {
  nameToken, err := raw.AssertWord(ts[0])
  if err != nil {
    return nil, nil, err
  }

  ctx := nameToken.Context()
  attr := html.NewEmptyRawDict(ctx)

  ts = ts[1:]
  if raw.IsSymbol(ts[0], patterns.PARENS_START) {
    parens, rem, err := p.buildParens(ts[0:])
    if err != nil {
      return nil, nil, err
    }

    attr = parens.ToRawDict()

    ts = rem
  } 

  tag := html.NewTag(nameToken.Value(), attr, []*html.Tag{}, ctx)

	// inline management is done by first none inline parent
	if inline {
		return tag, ts, nil
	}

	// while line is not empty find the children
	lineIsEmpty := func() bool {
    for _, t := range ts {
			if raw.IsIndent(t) {
				continue
			} else if raw.IsNL(t) {
				return true
			} else {
				return false
			}
		}

		return true
	}

	stack := make([]*html.Tag, 1)
	stack[0] = tag

	for !lineIsEmpty() {
		if raw.IsSymbol(ts[0], "<") {
			// pop the stack
			if len(stack) == 1 {
				errCtx := ts[0].Context()
				return nil, nil, errCtx.NewError("Error: cannot decrease inline stack before first child")
			} else {
				stack = stack[0 : len(stack)-1]
        ts = ts[1:]
			}
		} else {
			inlineTag, rem, err := p.buildTag(-1, ts)
			if err != nil {
				return nil, nil, err
			}

			if err := stack[len(stack)-1].AppendChild(inlineTag); err != nil {
				return nil, nil, err
			}

			// dont append text tags to the stack, this is a nuisance
			if !inlineTag.IsText() {
				stack = append(stack, inlineTag)
			}

      ts = rem
		}
	}

	return tag, ts, nil
}

// indent -1 means that we are inlining
// also returns the remaining tokens
func (p *TemplateParser) buildTag(indent int, ts []raw.Token) (*html.Tag, []raw.Token, error) {
  // no more eating needed in here?
	tagCtx := ts[0].Context()

  followedByParens := len(ts) > 1 && raw.IsSymbol(ts[1], patterns.PARENS_START)

	switch {
    case raw.IsLiteralString(ts[0]) || raw.IsSymbol(ts[0], patterns.DOLLAR) || raw.IsSymbol(ts[0], patterns.PARENS_START):
      return p.buildTextTag(indent == -1, ts)
    case indent != -1 && raw.IsAnyWord(ts[0]):
      keyToken, err := raw.AssertWord(ts[0])
      if err != nil {
        panic(err)
      }

      key := keyToken.Value()
      switch {
        case key == "export":
          // exports can be indented so they can benefit from branching
          return p.buildExportedDirective(ts)
        case key == "permissive":
          return nil, nil, tagCtx.NewError("Error: 'permissive' must be first word, and can't be indented")
        case key == "parameters":
          if indent != 0 {
            return nil, nil, tagCtx.NewError("Error: 'parameters' cannot be indented")
          }
          return p.buildParametersDirective(ts)
        case key == "import":
          return p.buildImportExportDirective(indent != 0, ts)
        case !followedByParens && key == "template":
          return p.buildTemplateDirective(ts)
        case !followedByParens && key == "var":
          return p.buildVarDirective(ts)
        case key == "function":
          return p.buildFunctionDirective(ts)
        case key == "for":
          return p.buildForDirective(ts)
        case key == "if" || key == "elseif" || key == "else" || key == "append" || key == "replace" || key == "prepend" || key == "block" || key == "switch" || key == "case" || key == "default":
          return p.buildSingleOrNoValueDirective(ts)
        default:
          return p.buildGenericTag(false, ts)
      }
    case indent == -1 && raw.IsAnyWord(ts[0]):
      return p.buildGenericTag(true, ts)
    default:
      errCtx := ts[0].Context()
      return nil, nil, errCtx.NewError("Error: not a tag")
	}
}

func (p *TemplateParser) removeWhitespace(ts []raw.Token) []raw.Token {
  tsInner := []raw.Token{}

  // filter out the whitespace
  for _, t := range ts {
    if !raw.IsWhitespace(t) {
      tsInner = append(tsInner, t)
    }
  }

  return tsInner
}

func (p *TemplateParser) buildExpressionInternal(ts []raw.Token) (html.Token, error) {
  ctx := ts[0].Context()
  tsInner := p.removeWhitespace(ts)

  // dummy FormulaParser for operators
  fp, err := NewFormulaParser("", ctx)
	if err != nil {
		return nil, err
	}

  tsInner, err = fp.nestGroups(tsInner)
	if err != nil {
		return nil, err
	}


	tsInner, err = fp.nestOperators(tsInner)
	if err != nil {
		return nil, err
	}

	tsInner = p.expandTmpGroups(tsInner)

  resExpr, err := p.helper.buildToken(tsInner)

  if err != nil {
    for _, t := range tsInner {
      fmt.Println(t.Dump(""))
    }
  }
  return resExpr, err
}

func (p *TemplateParser) buildExpression(ts []raw.Token) (html.Token, []raw.Token, error) {
  if len(ts) == 0 {
    panic("no expression tokens")
  }

  isExpectsMoreOp := func(t raw.Token) bool {
    if raw.IsAnySymbol(t) {
      s, err := raw.AssertAnySymbol(t)
      if err != nil {
        panic(err)
      }

      v := s.Value()

      // only operators, so can't use symbols pattern

      if v == ":" || v == "$" || v == "?" || v == "<" || v == ">" || v == "+" || v == "-" || v == "/" || v == "*" || v == "!=" || v == "==" || v == ":=" || v == "<=" || v == ">=" || v == "!" || v == "===" || v == "&&" || v == "||" || v == "!!" || v == "??" {
        return true
      } else {
        return false
      }
    } else {
      return false
    }
  }

  // find the end
  iStop := 0
  expectsMore := true
  groupCount := 0
  for i, t := range ts {
    if raw.IsSymbol(t, patterns.BRACES_START) || raw.IsSymbol(t, patterns.PARENS_START) || raw.IsSymbol(t, patterns.BRACKETS_START) {
      groupCount += 1
      expectsMore = true
    } else if raw.IsSymbol(t, patterns.BRACES_STOP) || raw.IsSymbol(t, patterns.PARENS_STOP) || raw.IsSymbol(t, patterns.BRACKETS_STOP) {
      if groupCount == 0 {
        errCtx := t.Context()
        return nil, nil, errCtx.NewError("Error: unmatched closing tag")
      }

      groupCount -= 1

      if groupCount == 0 {
        expectsMore = false
      }
    } else if isExpectsMoreOp(t) {
      expectsMore = true
    } else if raw.IsNL(t) && !expectsMore {
      iStop = i
      break
    } else if groupCount == 0 && !raw.IsWhitespace(t) {
      expectsMore = false
    }

    iStop = i+1
  }

  resExpr, err := p.buildExpressionInternal(ts[0:iStop])

  return resExpr, ts[iStop:], err
}

func (p *TemplateParser) buildTextTagExpression(ts []raw.Token) (html.Token, []raw.Token, error) {
  if len(ts) == 0 {
    panic("no expression tokens")
  }

  isExpectsMoreOp := func(t raw.Token) bool {
    if raw.IsAnySymbol(t) {
      s, err := raw.AssertAnySymbol(t)
      if err != nil {
        panic(err)
      }

      v := s.Value()

      // only operators, so can't use symbols pattern
      if v == ":" || v == "$" || v == "?" || v == "+" || v == "!=" || v == "==" || v == "!" || v == "===" || v == "&&" || v == "||" || v == "!!" || v == "??" {
        return true
      } else {
        return false
      }
    } else {
      return false
    }
  }

  // find the end
  iStop := 0
  expectsMore := true
  groupCount := 0
  for i, t := range ts {
    if raw.IsSymbol(t, patterns.BRACES_START) || raw.IsSymbol(t, patterns.PARENS_START) || raw.IsSymbol(t, patterns.BRACKETS_START) {
      groupCount += 1
      expectsMore = true
    } else if raw.IsSymbol(t, patterns.BRACES_STOP) || raw.IsSymbol(t, patterns.PARENS_STOP) || raw.IsSymbol(t, patterns.BRACKETS_STOP) {
      if groupCount == 0 {
        errCtx := t.Context()
        return nil, nil, errCtx.NewError("Error: unmatched closing tag")
      }

      groupCount -= 1

      if groupCount == 0 {
        expectsMore = false
      }
    } else if isExpectsMoreOp(t) {
      expectsMore = true
    } else if (raw.IsNL(t) || raw.IsAnyWord(t) || raw.IsLiteralString(t) || raw.IsSymbol(t, "<")) && !expectsMore  {
      iStop = i
      break
    } else if (raw.IsAnyWord(t) || raw.IsLiteralString(t) || raw.IsSymbol(t, "<")) && groupCount == 0 {
      expectsMore = false
    }

    iStop = i+1
  }

  resExpr, err := p.buildExpressionInternal(ts[0:iStop])

  return resExpr, ts[iStop:], err
}

// return indent from first non-empty line
func (p *TemplateParser) eatWhitespace(ts []raw.Token) ([]raw.Token, int) {
  for i, t := range ts {
    if !raw.IsWhitespace(t) {
      if i == 0 {
        return ts, 0
      } else {
        prev := ts[i-1]
        if raw.IsIndent(prev) {
          indent, err := raw.AssertIndent(prev)
          if err != nil {
            panic(err)
          }

          return ts[i:], indent.N()
        } else {
          return ts[i:], 0
        }
      }
    }
  }

  return []raw.Token{}, 0
}

func (p *TemplateParser) DumpTokens() {
	fmt.Println("\nRaw tokens:")
	fmt.Println("===========")

  ts, err := p.tokenizeFlat()
  if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
  }

	for _, t := range ts {
		fmt.Println(t.Dump(""))
	}

	fmt.Println("\nTemplate tokens:")
	fmt.Println("===========")

	p.Reset()
	tags, err := p.BuildTags()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, tag := range tags {
		fmt.Println(tag.Dump(""))
	}
}
