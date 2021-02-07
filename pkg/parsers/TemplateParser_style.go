package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *TemplateParser) collectSelector(ts []raw.Token) ([]raw.Token, *raw.Group, []raw.Token, error) {
  for i, t := range ts {
    if raw.IsBracesGroup(t) {
      if i == 0 {
        errCtx := t.Context()
        return nil, nil, nil, errCtx.NewError("Error: expected selector before")
      } 

      brace, err := raw.AssertBracesGroup(t)
      if err != nil {
        panic(err)
      }

      if brace.IsComma() {
        errCtx := brace.Context()
        return nil, nil, nil, errCtx.NewError("Error: expected semicolon separator")
      }

      return ts[0:i], brace, ts[i+1:], nil
    }
  }

  errCtx := ts[len(ts)-1].Context()
  return nil, nil, nil, errCtx.NewError("Error: expected brace after selector")
}

func (p *TemplateParser) trimWhitespace(ts []raw.Token) []raw.Token {
  start := 0
  for start < len(ts) && raw.IsWhitespace(ts[start]) {
    start += 1
  }

  stop := len(ts) -1
  for stop >= 0 && raw.IsWhitespace(ts[stop]) {
    stop -= 1
  }

  if start >= stop {
    return []raw.Token{}
  } else {
    return ts[start:stop+1]
  }
}

func (p *TemplateParser) cssTokensToString(ts []raw.Token, befWS bool, aftWS bool) (*html.String, error) {
  ctx := ts[0].Context()

  if len(ts) > 1 {
    ctx = context.SimpleFill(ts[0].Context(), ts[len(ts)-1].Context())
  }

  if befWS {
    ctx = ctx.IncludeLeftSpace()
  }

  if aftWS {
    ctx = ctx.IncludeRightSpace()
  }

  return html.NewValueString(ctx.Content(), ctx), nil
}

func (p *TemplateParser) cssTokensToStrFn(ts []raw.Token) (html.Token, error) {
  // collect the dollar parts

  parts := []html.Token{}

  prev := ts[0]

  fillSpace := func(at raw.Token, bt raw.Token) {
    a := at.Context()
    b := bt.Context()
    if !b.IsConsecutive(a) {
      b = b.IncludeLeftSpace()
      a = a.IncludeRightSpace()
      fill := context.SimpleFill(b, a)
      parts = append(parts, html.NewValueString(fill.Content(), fill))
    }
  }

  // inner function that can be called recursively if nested groups are unpacked
  var fn func([]raw.Token) error 
  fn = func(ts_ []raw.Token) error {
    dollarStart := -1
    for i := 0; i < len(ts_); i++ {
      t := ts_[i]
      if dollarStart != -1 {
        if !(raw.IsParensGroup(t) || raw.IsBracketsGroup(t)) {
          exprTs, err := p.nestOperators(ts_[dollarStart:i])
          if err != nil {
            return err
          }

          expr, err := p.buildExpression(exprTs)
          if err != nil {
            return err
          }

          parts = append(parts, expr)

          dollarStart = -1
          prev = t

          if raw.IsSymbol(t, "$") {
            fillSpace(ts_[i-1], t)
          }
        } else {
          continue
        }
      } 

      if raw.IsSymbol(t, "$") {
        if i == len(ts_) - 1 {
          errCtx := t.Context()
          return errCtx.NewError("Error: expected word after $")
        }

        if !raw.IsAnyWord(ts_[i+1]) {
          errCtx := t.Context()
          return errCtx.NewError("Error: expected word after $")
        }

        if !raw.IsSymbol(prev, "$") {
          part, err := p.cssTokensToString([]raw.Token{prev, ts_[i-1]}, len(parts) > 0, true)
          if err != nil {
            return err
          }

          if part != nil {
            parts = append(parts, part)
          }

          // also add some white space
          fillSpace(ts_[i-1], t)
        }

        prev = t
        dollarStart = i
        i += 1
      } else if raw.IsGroup(t) {
        // unpack the group, and call self recursively
        tGroup, err := raw.AssertGroup(t)
        if err != nil {
          panic(err)
        }

        tGroupTokens := tGroup.ExpandOnce()
        if err := fn(tGroupTokens); err != nil {
          return err
        }

        i += 1
      } else if raw.IsSymbol(prev, "$") {
        prev = t
      }
    }

    if dollarStart != -1 {
      exprTs, err := p.nestOperators(ts_[dollarStart:])
      if err != nil {
        return err
      }

      expr, err := p.buildExpression(exprTs)
      if err != nil {
        return err
      }

      parts = append(parts, expr)

      dollarStart = -1
    }

    return nil
  }

  if err := fn(ts); err != nil {
    return nil, err
  }


  if !raw.IsSymbol(prev, "$") {
    part, err := p.cssTokensToString([]raw.Token{prev, ts[len(ts)-1]}, prev != ts[0], false)
    if err != nil {
      return nil, err
    }

    if part != nil {
      parts = append(parts, part)
    }
  }

  ctx := raw.MergeContexts(ts...)
  return html.NewFunction("str", []html.Token{html.NewValuesList(parts, ctx)}, ctx), nil
}

func (p *TemplateParser) cssBraceToDict(brace *raw.Group) (html.Token, error) {
  ctx := brace.Context()

  parts := []html.Token{}

  // while colon is found keep adding to current
  current := html.NewEmptyRawDict(ctx)

  addProperty := func(lhs_ []raw.Token, rhs_ []raw.Token) error {
    lhs, err := p.cssTokensToStrFn(lhs_)
    if err != nil {
      return err
    }

    rhs, err := p.cssTokensToStrFn(rhs_)
    if err != nil {
      return err
    }

    current.Set(lhs, rhs)

    return nil
  }

  addMixin := func(ts []raw.Token) error {
    exprTs, err := p.nestOperators(ts)
    if err != nil {
      return err
    }

    expr, err := p.buildExpression(exprTs)
    if err != nil {
      return err
    }

    //if current.Len() != 0 {
    parts = append(parts, current) // adding an empty dict assert that merge works on dicts

    current = html.NewEmptyRawDict(ctx)
    //}

    parts = append(parts, expr)

    return nil
  }

  nestRule := func(selTs []raw.Token, nBrace *raw.Group) error {
    sel, err := p.cssTokensToStrFn(selTs)
    if err != nil {
      return err
    }

    rule, err := p.cssBraceToDict(nBrace)
    if err != nil {
      return err
    }

    current.Set(sel, rule)

    return nil
  }

  for _, field := range brace.Fields {
    inner := field

    for len(inner) > 0 {
      // look for nested brace
      nestedBraceI := raw.FindFirstBracesGroup(inner, 0)

      if nestedBraceI == -1 {
        // look for colon
        sub := raw.SplitBySymbol(inner, ":")
        if len(sub) == 1 {
          if err := addMixin(sub[0]); err != nil {
            return nil, err
          }
        } else {
          rhs := sub[1]
          for _, s := range sub[2:] {
            rhs = append(rhs, s...)
          }

          if err := addProperty(sub[0], rhs); err != nil {
            return nil, err
          }
        }

        inner = []raw.Token{}
      } else {
        nestedBrace, err := raw.AssertBracesGroup(inner[nestedBraceI])
        if err != nil {
          panic(err)
        }

        if err := nestRule(inner[0:nestedBraceI], nestedBrace); err != nil {
          return nil, err
        }
        inner = inner[nestedBraceI+1:]
      }
    }
  }

  parts = append(parts, current)

  return html.NewFunction("merge", parts, ctx), nil
}

func (p *TemplateParser) buildStyleContent(indent int, ts []raw.Token, ctx context.Context) (html.Token, []raw.Token, error) {
  parts := []html.Token{}

  if len(ts) > 0 {
    ctx = raw.MergeContexts(ts...)
  }

  dict := html.NewEmptyRawDict(ctx)

  tsBef := ts
  nBef := len(ts)

  for len(ts) > 0 {
    rem, ruleIndent := p.eatWhitespace(ts)

    if ruleIndent <= indent {
      break
    } else {
      ts = p.eatLine(rem)
    }
  }

  nAft := len(ts)
  innerTs := tsBef[0:nBef-nAft]

  var err error
  innerTs, err = p.nestGroups(raw.RemoveWhitespace(innerTs))
  if err != nil {
    return nil, nil, err
  }

  if len(innerTs) == 1 && raw.IsBracesGroup(innerTs[0]) {
    brace, err := raw.AssertBracesGroup(innerTs[0])
    if err != nil {
      panic(err)
    }

    part, err := p.cssBraceToDict(brace)
    if err != nil {
      return nil, nil, err
    }

    parts = append(parts, part)

    innerTs = []raw.Token{}
  }

  for len(innerTs) > 0 {
    if len(innerTs) < 2 {
      errCtx := raw.MergeContexts(innerTs...)
      return nil, nil, errCtx.NewError("Error: expected more tokens")
    }

    // top level mixin is also allowed
    if raw.IsSymbol(innerTs[0], "$") {
      continueOuter := false
      for i, innerT := range innerTs {
        atEnd := (i == len(innerTs) - 1 && !raw.IsBracesGroup(innerT))
        if raw.IsSymbol(innerT, ";") || atEnd {
          stop := i
          if atEnd && !raw.IsSymbol(innerT, ";") {
            stop = len(innerTs)
          }

          exprTs := innerTs[0:stop]
          exprTs, err = p.nestOperators(exprTs)
          if err != nil {
            return nil, nil, err
          }

          expr, err := p.buildExpression(exprTs)
          if err != nil {
            return nil, nil, err
          }

          //if current.Len() != 0 {
          parts = append(parts, dict) // adding an empty dict assert that merge works on dicts

          //}

          parts = append(parts, expr)

          dict = html.NewEmptyRawDict(ctx)

          if stop < len(innerTs) {
            innerTs = innerTs[stop+1:]
          } else {
            innerTs = []raw.Token{}
          }

          continueOuter = true
          break
        } else if (raw.IsBracesGroup(innerT)) {
          break
        }
      }

      if continueOuter {
        continue
      }
    }

    selTs, brace, rem, err := p.collectSelector(innerTs)
    if err != nil {
      return nil, nil, err
    }

    // selTs has to be turned into something that evaluates into a string
    // str([parts...], "")
    sel, err := p.cssTokensToStrFn(selTs)
    if err != nil {
      return nil, nil, err
    }

    rule, err := p.cssBraceToDict(brace)
    if err != nil {
      return nil, nil, err
    }

    dict.Set(sel, rule)

    innerTs = rem
  }

  parts = append(parts, dict)

  return html.NewFunction("merge", parts, ctx), ts, nil
}

func (p *TemplateParser) buildStyleDirective(indent int, ts []raw.Token) (*html.Tag, []raw.Token, error) {
  ctx := ts[0].Context()

  if len(ts) < 3 {
    errCtx := raw.MergeContexts(ts...)
    return nil, nil, errCtx.NewError("Error: expected more tokens")
  }

  ts = ts[1:]
  var attr *html.RawDict = nil
  if raw.IsSymbol(ts[0], patterns.PARENS_START) {
    parens, rem, err := p.buildParens(ts[0:])
    if err != nil {
      return nil, nil, err
    }

    attr = parens.ToRawDict()

    ts = rem
  }

  var nameToken *raw.Word = nil
  var nameArgs *html.Parens = nil
  if attr != nil {
    if !raw.IsNL(ts[0]) {
      errCtx := ts[0].Context()
      return nil, nil, errCtx.NewError("Error: unexpected tokens on same line")
    }


  } else {
    var err error
    nameToken, err = raw.AssertWord(ts[0])
    if err != nil {
      return nil, nil, err
    }

    ts = ts[1:]

    if raw.IsSymbol(ts[0], patterns.PARENS_START) {
      var rem []raw.Token
      nameArgs, rem, err = p.buildParens(ts[0:])
      if err != nil {
        return nil, nil, err
      }

      ts = rem
    }
  }

  dict, rem, err := p.buildStyleContent(indent, ts, ctx)
  if err != nil {
    return nil, nil, err
  }

  ts = rem

  if attr != nil {
    // generic tag with .content attr
    attr.Set(html.NewValueString(".content", dict.Context()), dict)

    return html.NewTag("style", attr, []*html.Tag{}, ctx), ts, nil
  } else {
    nameHtmlToken := html.NewValueString(nameToken.Value(), nameToken.Context())
    varAttr := html.NewEmptyRawDict(ctx)

    if nameArgs != nil {
      //statements := []html.Token{dict}

      //list := html.NewValuesList(statements, ctx)
      //index := html.NewValueInt(0, ctx)
      //wrapper := html.NewFunction("get", []html.Token{list, index}, ctx)

      varAttr.Set(nameHtmlToken, html.NewFunction("function", []html.Token{nameArgs, dict}, ctx))
    } else {
      varAttr.Set(nameHtmlToken, dict)
    }

    return html.NewDirectiveTag("var", varAttr, []*html.Tag{}, ctx), ts, nil
  } 
}
