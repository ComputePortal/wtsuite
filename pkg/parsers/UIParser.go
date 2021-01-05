package parsers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func tokenizeUIFormulas(s string, ctx context.Context) ([]raw.Token, error) {
	return nil, ctx.NewError("Error: can't have backtick formula in ui markup")
}

var uiParserSettings = ParserSettings{
	quotedGroups: quotedGroupsSettings{
		pattern: patterns.UI_STRING_OR_COMMENT_REGEXP,
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
		tokenizer: tokenizeUIFormulas,
	},
	// same as html
	wordsAndLiterals: wordsAndLiteralsSettings{
		maskType:  WORD_OR_LITERAL,
		pattern:   patterns.HTML_WORD_OR_LITERAL_REGEXP,
		tokenizer: tokenizeHTMLWordsAndLiterals,
	},
	// handled by FormulaParser
	symbols: symbolsSettings{
		maskType: SYMBOL,
		pattern:  nil,
	},
	operators:                newOperatorsSettings([]operatorSettings{}), // handled by FormulaParser
	tmpGroupWords:            true,
	tmpGroupPeriods:          false,
	tmpGroupArrows:           false,
	tmpGroupDColons:          false,
	tmpGroupAngled:           false,
	recursivelyNestOperators: true,
}

type UIParser struct {
	helper *HTMLParser
	Parser
}

func NewUIParser(path string) (*UIParser, error) {
	rawBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := string(rawBytes)
	src := context.NewSource(raw)

	ctx := context.NewContext(src, path)
	p := &UIParser{
		NewEmptyHTMLParser(context.NewDummyContext()),
		newParser(raw, uiParserSettings, ctx),
	}

	if err := p.maskQuoted(); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *UIParser) BuildTags() ([]*html.Tag, error) {
	result := make([]*html.Tag, 0)

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
				if !(inBranch && indents[i] == indent && stack[i].Name() == "ifelse") {
					stack = stack[0:i]
					indents = indents[0:i]
					break
				}
			}
		}

		if t.Name() == "if" && (len(stack) == 0 || stack[len(stack)-1].Name() != "ifelse") {
			ifElseCtx := t.Context()
			ifElseTag := html.NewDirectiveTag("ifelse", html.NewEmptyRawDict(ifElseCtx),
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
	for true {
		indent, eof := p.eatWhitespace()
		if eof {
			break
		}

		tag, err := p.buildTag(indent)
		if err != nil {
			return nil, err
		}

		if err := appendTag(tag, indent); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (p *UIParser) buildTextTag(str string, inline bool, ctx context.Context) (*html.Tag, error) {
	// literal strings are text tags
	// they cannot be followed by anything else
	if !inline && p.eatRestOfLine() {
		return nil, ctx.NewError("Error: unexpected content after")
	}

	// cut off the quotes
	str = str[1 : len(str)-1]
	return html.NewTextTag(str, ctx), nil
}

func (p *UIParser) buildGenericDirective(name string, ctx context.Context) (*html.Tag, error) {
	fnPre := func(_ html.Token, ts []raw.Token) []raw.Token {
		return p.convertWords(ts)
	}

	attr, err := p.buildDirectiveAttributes(-1, -1, ctx, fnPre)
	if err != nil {
		return nil, err
	}

	// children are appended later
	return html.NewDirectiveTag(name, attr, []*html.Tag{}, ctx), nil
}

func (p *UIParser) buildImportExportDirective(name string, dynamic bool, ctx context.Context) (*html.Tag, error) {
	fnPre := func(_ html.Token, ts []raw.Token) []raw.Token {
		return p.convertWords(ts)
	}

	attr, err := p.buildDirectiveAttributes(-1, -1, ctx, fnPre)
	if err != nil {
		return nil, err
	}

	attr.Set(html.NewValueString(".dynamic", ctx), html.NewValueBool(dynamic, ctx))

	// children are appended later
	return html.NewDirectiveTag(name, attr, []*html.Tag{}, ctx), nil
}

func (p *UIParser) matchFunctionDirective(ctx context.Context) ([2]int, error) {
	containerCount := 0

	firstBrace := -1
	start := p.pos
	for ; p.pos < p.Len(); p.pos++ {
		if p.mask[p.pos] == COMMENT || p.mask[p.pos] == STRING {
			continue
		}

		c := p.raw[p.pos]

		if c == '{' {
			if containerCount == 0 {
				firstBrace = p.pos
			}
			containerCount += 1
		} else if c == '}' {
			containerCount -= 1
			if containerCount < 0 {
				errCtx := p.NewContext(p.pos, p.pos+1)
				return [2]int{0, 0}, errCtx.NewError("Error: unmatched container end")
			} else if containerCount == 0 {
				p.pos += 1
				return [2]int{start, p.pos}, nil
			}
		}
	}

	if firstBrace == -1 {
		return [2]int{0, 0}, ctx.NewError("Error: not function brace group found")
	} else if containerCount != 0 {
		errCtx := p.NewContext(firstBrace, firstBrace+1)
		return [2]int{0, 0}, errCtx.NewError("Error: unmatched start brace")
	} else {
		panic("shouldn't be possible")
	}
}

// everything following the function keyword
// eat until first closing brace (containerCount == 0
func (p *UIParser) buildFunctionDirective(ctx context.Context) (*html.Tag, error) {
	r, err := p.matchFunctionDirective(ctx)
	if err != nil {
		return nil, err
	}

	// s := p.writeWithoutLineContinuation(r[0], r[1])
	//fmt.Printf("function value from %d to %d: \"%s\"\n", r[0], r[1], s)

	innerCtx := p.NewContext(r[0], r[1])

	ts, err := p.tokenizePartial(r[0], r[1], p.convertWords)
	if err != nil {
		return nil, err
	}

	if len(ts) < 3 {
		return nil, ctx.NewError("Error: insufficient tokens (expected at least 3)")
	}

	// first should be a strings
	nameToken, err := raw.AssertWord(ts[0])
	if err != nil {
		return nil, err
	}

	ts[0] = raw.NewValueWord("function", ctx)

	fnValue, err := p.buildToken(ts, false)
	if err != nil {
		return nil, err
	}

	varAttr := html.NewEmptyRawDict(innerCtx)

	nameHtmlToken := html.NewValueString(nameToken.Value(), nameToken.Context())
	varAttr.Set(nameHtmlToken, fnValue)

	return html.NewDirectiveTag("var", varAttr, []*html.Tag{}, ctx), nil
}

func (p *UIParser) buildVarDirective(ctx context.Context) (*html.Tag, error) {
	return p.buildGenericDirective("var", ctx)
}

func (p *UIParser) buildPermissiveDirective(ctx context.Context) (*html.Tag, error) {
  tag, err := p.buildGenericDirective("permissive", ctx)
  if err != nil {
    return nil, err
  }

  if tag.RawAttributes().Len() != 0 {
    errCtx := tag.RawAttributes().Context()
    return nil, errCtx.NewError("Error: unexpected attributes")
  }

  return tag, nil
}

func (p *UIParser) buildTagAttributes(ts_ []raw.Token, ctx context.Context) (*html.RawDict, error) {

	if len(ts_) != 1 || !raw.IsParensGroup(ts_[0]) {
		errCtx := ctx
		return nil, errCtx.NewError("Error: bad tag attributes")
	}

	group, err := raw.AssertParensGroup(ts_[0])
	if err != nil {
		panic(err)
	}

  if group.IsSemiColon() {
    errCtx := group.Context()
    return nil, errCtx.NewError("Error: expected commas as separators")
  }

	tss := group.Fields

	attr := html.NewEmptyRawDict(ctx)

	if len(tss) == 0 {
		return attr, nil
	}

	hadOpts := false
	posCount := 0

	// positional comes first
	for _, ts := range tss {
		switch {
		case len(ts) == 1 && raw.IsOperator(ts[0], "bin="):
			hadOpts = true

			op, err := raw.AssertAnyBinaryOperator(ts[0])
			if err != nil {
				return nil, err
			}

			w, err := raw.AssertWord(op.Args()[0])
			if err != nil {
				panic(err)
			}

			posKey := html.NewValueString(w.Value(), w.Context())

			posVal, err := p.buildToken(op.Args()[1:], true)
			if err != nil {
				return nil, err
			}

			attr.Set(posKey, posVal)
    case len(ts) == 1 && raw.IsOperator(ts[0], "bin!="):
			hadOpts = true

			op, err := raw.AssertAnyBinaryOperator(ts[0])
			if err != nil {
				return nil, err
			}

			w, err := raw.AssertWord(op.Args()[0])
			if err != nil {
				panic(err)
			}

			posKey := html.NewValueString(w.Value() + "!", w.Context())

			posVal, err := p.buildToken(op.Args()[1:], true)
			if err != nil {
				return nil, err
			}

      // must then be treated specially later
			attr.Set(posKey, posVal)
		case len(ts) == 1:
			if hadOpts {
				errCtx := ts[0].Context()
				return nil, errCtx.NewError("Error: positional must come before optional")
			}

			posVal, err := p.buildToken(ts, true)
			if err != nil {
				return nil, err
			}

			posKey := html.NewValueInt(posCount, ts[0].Context())
			attr.Set(posKey, posVal)
			posCount += 1
		case len(ts) == 0:
			errCtx := ctx
			return nil, errCtx.NewError("Error: unexpected empty attr field")
		default:
			errCtx := ts[0].Context()
			return nil, errCtx.NewError("Error: bad tag call field")
		}
	}

	return attr, nil
}

func (p *UIParser) buildTemplateDirective(ctx context.Context) (*html.Tag, error) {
  p.eatWhitespace() 
	rName, isString := p.eatNonWhitespace()
  if rName[0] == rName[1] {
    errCtx := ctx
    return nil, errCtx.NewError("Error: couldn't find the template name")
  } else if isString {
    errCtx := p.NewContext(rName[0], rName[1])
    return nil, errCtx.NewError("Error: expected name of template, got literal string instead")
  }
  nameCtx := p.NewContext(rName[0], rName[1])
	nameKey := html.NewValueString("name", nameCtx)
	nameVal := html.NewValueString(p.Write(rName[0], rName[1]), nameCtx)

	start := p.pos

	// look for the closing super, and get the attributes from that range
	rSuper, _, superOk := p.nextMatch(patterns.UI_TEMPLATE_SUPER_REGEXP, true)
	if !superOk {
		return nil, ctx.NewError("Error: no super keyword found")
	}
	p.pos -= 1 // include opening parens for next match

	_, _, superOpenOk := p.nextMatch(patterns.PARENS_OPEN_REGEXP, true)
	if !superOpenOk {
		errCtx := p.NewContext(rSuper[0], rSuper[1])
		return nil, errCtx.NewError("Error: bad format for super")
	}

	rSuperParens, superCloseOk := p.nextGroupStopMatch(patterns.PARENS_GROUP, true)
	if !superCloseOk {
		errCtx := p.NewContext(rSuper[0], rSuper[1])
		return nil, errCtx.NewError("Error: bad format for super")
	}

	fnPre := func(key html.Token, ts []raw.Token) []raw.Token {
		if html.IsInt(key) {
			iKey, err := html.AssertInt(key)
			if err != nil {
				panic(err)
			}

			if iKey.Value() == 0 {
				res := make([]raw.Token, 0)

				for i := 0; i < len(ts); i++ {
          var tNext raw.Token = nil
          if i < len(ts)-1 {
            tNext = ts[i+1]
          }
          t := p.convertWord(ts[i], tNext)
          res = append(res, t)
				}

				return res
			}
		}

		return p.convertWords(ts)
	}

	attr, err := p.buildDirectiveAttributes(start, rSuper[0], ctx, fnPre)
	if err != nil {
		return nil, err
	}

  // name was eaten before
  attr.Set(nameKey, nameVal);

	// parse the super
	// positional must come before optional args
	superCtx := p.NewContext(rSuper[0], rSuper[1])
	// include the parens themselves
	ts, err := p.tokenizePartial(rSuper[1]-1, rSuperParens[1], nil)
	if err != nil {
		panic(err)
		return nil, err
	}

	superKey := html.NewValueString("super", superCtx)
	superVal, err := p.buildTagAttributes(ts, superCtx)
	if err != nil {
		return nil, err
	}

	attr.Set(superKey, superVal)

	p.pos = rSuperParens[1]

	// search for the optional this(...)
	p.pos += 1
	rThis, isString := p.eatNonWhitespace()
	if !isString && rThis[0] != rThis[1] && p.Write(rThis[0], rThis[1]) == "this" {
		_, _, thisOpenOk := p.nextMatch(patterns.PARENS_OPEN_REGEXP, true)
		if !thisOpenOk {
			errCtx := p.NewContext(rThis[0], rThis[1])
			return nil, errCtx.NewError("Error: bad format for this")
		}

		rThisParens, thisCloseOk := p.nextGroupStopMatch(patterns.PARENS_GROUP, true)
		if !thisCloseOk {
			errCtx := p.NewContext(rSuper[0], rSuper[1])
			return nil, errCtx.NewError("Error: bad format for this")
		}

		thisCtx := p.NewContext(rThis[0], rThis[1])
		// include the parens themselves
		ts, err := p.tokenizePartial(rThis[1], rThisParens[1], nil)
		if err != nil {
			panic(err)
			return nil, err
		}

		thisKey := html.NewValueString("this", thisCtx)
		thisVal, err := p.buildTagAttributes(ts, thisCtx)
		if err != nil {
			panic(err)
			return nil, err
		}

		attr.Set(thisKey, thisVal)

		p.pos = rThisParens[1]

	} else {
		p.pos = rSuperParens[1]
	}

	return html.NewDirectiveTag("template", attr, []*html.Tag{}, ctx), nil
}

func (p *UIParser) buildExportedDirective(ctx context.Context) (*html.Tag, error) {
	addExportAttr := func(attr *html.RawDict, exportCtx context.Context) {
		exportToken := html.NewValueString("export", exportCtx)
		flagToken := html.NewValueString("", exportCtx)
		attr.Set(exportToken, flagToken)
	}

	prevPos := p.pos
	rName, isString := p.eatNonWhitespace()
	if isString || rName[0] == rName[1] {
		// is regular export
		p.pos = prevPos
		return p.buildGenericDirective("export", ctx)
	} else {
		exportCtx := ctx
		tagCtx := p.NewContext(rName[0], rName[1])

		name := p.Write(rName[0], rName[1])
		switch name {
		case "var":
			tag, err := p.buildVarDirective(tagCtx)
			if err != nil {
				return nil, err
			}

			addExportAttr(tag.RawAttributes(), exportCtx)

			return tag, nil
		case "template":
			tag, err := p.buildTemplateDirective(tagCtx)
			if err != nil {
				return nil, err
			}

			addExportAttr(tag.RawAttributes(), exportCtx)

			// add export flag to tag
			//return html.NewTag(name, attr, []*html.Tag{}, tagCtx), nil
			return tag, nil
		case "function":
			tag, err := p.buildFunctionDirective(ctx)
			if err != nil {
				return nil, err
			}

			addExportAttr(tag.RawAttributes(), exportCtx)

			return tag, nil
		default:
      if !strings.HasPrefix(name, "[") && !strings.HasPrefix(name, "{") {
        errCtx := exportCtx
        return nil, errCtx.NewError("Error: invalid export statement")
      }

			p.pos = prevPos
			return p.buildImportExportDirective("export", false, ctx)
		}
	}
}

func (p *UIParser) buildGenericTag(name string, inline bool, ctx context.Context) (*html.Tag, error) {
	// find parens on this line!
	// or none at all
	parensFound := false
	parensStart := -1
	for ; p.pos < p.Len(); p.pos++ {
		c := p.raw[p.pos]
		if c == ' ' || c == '\t' {
			continue
		}

		if c == '(' {
			parensFound = true
			parensStart = p.pos
			p.pos += 1
		}
		break
	}

	var tag *html.Tag = nil
	if !parensFound {
		tag = html.NewTag(name, html.NewEmptyRawDict(ctx), []*html.Tag{}, ctx)
	} else {
		startCtx := p.NewContext(p.pos, p.pos+1)

		if r, ok := p.nextGroupStopMatch(patterns.PARENS_GROUP, true); ok {
			ts, err := p.tokenizePartial(parensStart, r[1], nil)
			if err != nil {
				return nil, err
			}

			attr, err := p.buildTagAttributes(ts, startCtx)
			if err != nil {
				return nil, err
			}

			tag = html.NewTag(name, attr, []*html.Tag{}, ctx)
		} else {
			errCtx := startCtx
			return nil, errCtx.NewError("Error: closing parens not found")
		}
	}

	// inline management is done by first none inline parent
	if inline {
		return tag, nil
	}

	// while line is not empty find the children

	lineIsEmpty := func() bool {
		for ; p.pos < p.Len(); p.pos++ {
			c := p.raw[p.pos]
			if c == ' ' || c == '\t' || (p.mask[p.pos] == COMMENT && c != '\n') {
				continue
			}

			if c == '\n' {
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
		if p.raw[p.pos] == '<' {
			p.pos += 1
			// pop the stack
			if len(stack) == 1 {
				errCtx := p.NewContext(p.pos-1, p.pos)
				return nil, errCtx.NewError("Error: cannot decrease inline stack before first child")
			} else {
				stack = stack[0 : len(stack)-1]
			}
		} else {
			inlineTag, err := p.buildTag(-1)
			if err != nil {
				return nil, err
			}

			if err := stack[len(stack)-1].AppendChild(inlineTag); err != nil {
				return nil, err
			}

			// dont append text tags to the stack, this is a nuisance
			if !inlineTag.IsText() {
				stack = append(stack, inlineTag)
			}
		}
	}

	return tag, nil
}

// indent -1 means that we are inlining
func (p *UIParser) buildTag(indent int) (*html.Tag, error) {
	r, isString := p.eatNonWhitespace()
	tagCtx := p.NewContext(r[0], r[1])

	str := p.Write(r[0], r[1]) // name or literal string
  followedByParens := (r[1] < len(p.raw)) && (p.raw[r[1]] == '(') // non directive versions of template/var must be followed by parens immediately

	switch {
	case isString:
		return p.buildTextTag(str, indent == -1, tagCtx)
	case indent != -1 && str == "export":
		if indent != 0 {
			return nil, tagCtx.NewError("Error: export cannot be indented")
		}
		return p.buildExportedDirective(tagCtx)
  case indent != -1 && str == "permissive":
    if indent != 0 {
      return nil, tagCtx.NewError("Error: 'permissive' cannot be indented")
    } else if r[0] != 0 {
      return nil, tagCtx.NewError("Error: 'permissive' must be first word, and can't be indented")
    }
    return p.buildPermissiveDirective(tagCtx)
	case indent != -1 && str == "import":
		return p.buildImportExportDirective(str, indent != 0, tagCtx)
	case indent != -1 && !followedByParens && str == "template":
		return p.buildTemplateDirective(tagCtx)
	case indent != -1 && !followedByParens && str == "var":
		return p.buildVarDirective(tagCtx)
	case indent != -1 && str == "function":
		return p.buildFunctionDirective(tagCtx)
	case indent != -1 && (str == "for" || str == "if" || str == "else" || str == "elseif" || str == "switch" || str == "case" || str == "default" || str == "print" || str == "replace" || str == "append" || str == "prepend" || str == "block"):
		// dummy is like a regular tag!
		return p.buildGenericDirective(str, tagCtx)
	default:
		return p.buildGenericTag(str, indent == -1, tagCtx)
	}
}

// eat each next token (consequetive non-whitespace)
func (p *UIParser) buildDirectiveAttributes(start, end int, tagCtx context.Context,
	fnPre func(html.Token, []raw.Token) []raw.Token) (*html.RawDict, error) {

	keys := make([]html.Token, 0)
	values := make([]html.Token, 0)

	argPos := 0

	isPos := true

	if start >= 0 {
		p.pos = start
	}

	for true {
		rVal, rNextKey, err := p.matchNextDirectiveAttribute(end)
		if err != nil {
			return nil, err
		}

		if !isPos && rVal[0] == rVal[1] {
			errCtx := keys[len(keys)-1].Context()
			return nil, errCtx.NewError("Error: no value found for attr")
		}

		if rVal[0] == rVal[1] && rNextKey[0] == rNextKey[1] {
			return html.NewValuesRawDict(keys, values, tagCtx), nil
		}

		if rVal[0] != rVal[1] {
			var key html.Token = nil
			if isPos { // generate the positional key
				key = html.NewValueInt(argPos, p.NewContext(rVal[0], rVal[1]))
				keys = append(keys, key)
				argPos += 1
			} else {
				key = keys[len(keys)-1]
			}

			fnPre_ := func(ts []raw.Token) []raw.Token {
				return fnPre(key, ts)
			}

			val, err := p.buildDirectiveAttributeValue(rVal, fnPre_)
			if err != nil {
				return nil, err
			}

			values = append(values, val)

			// next value might be positional
			isPos = true
		}

		if rNextKey[0] != rNextKey[1] {
			keyStr := p.Write(rNextKey[0], rNextKey[1])

			keyCtx := p.NewContext(rNextKey[0], rNextKey[1])
			if !patterns.IsValidVar(keyStr) {
				return nil, keyCtx.NewError("Error: invalid attr name")
			}

			key := html.NewValueString(keyStr, keyCtx)

			keys = append(keys, key)

			// next value is not positional
			isPos = false
		}
	}

	panic("shouldn't get here")
}

func (p *UIParser) writeWithoutLineContinuation(start, stop int) string {
	if stop == -1 {
		stop = p.Len()
	}

	// remove the comment parts
	var b strings.Builder
	for i := start; i < stop; i++ {
		if p.raw[i] == '\\' || p.mask[i] == COMMENT {
			b.WriteRune(' ')
		} else {
			b.WriteRune(p.raw[i])
		}
	}

	return b.String()
}

func (p *UIParser) tokenizePartial(start, end int,
	fnPre func([]raw.Token) []raw.Token) ([]raw.Token, error) {
	s := p.writeWithoutLineContinuation(start, end)
	ctx := p.NewContext(start, end)

	fp, err := NewFormulaParser(s, ctx)
	if err != nil {
		return nil, err
	}

	ts, err := fp.Parser.tokenizeFlat()
	if err != nil {
		return nil, err
	}

	if fnPre != nil {
		ts = fnPre(ts)
	}

	// now nest the groups and the operators
	ts, err = fp.nestGroups(ts)
	if err != nil {
		return nil, err
	}

	ts, err = fp.nestOperators(ts)
	if err != nil {
		return nil, err
	}

	return p.expandTmpGroups(ts), nil
}

func (p *UIParser) convertWord(t raw.Token, tNext raw.Token) raw.Token {
	if raw.IsAnyWord(t) {
		w, err := raw.AssertWord(t)
		if err != nil {
			panic(err)
		}

		if strings.HasPrefix(w.Value(), "$") {
			return raw.NewValueWord(w.Value()[1:], w.Context())
		} else {
			followedByParens := false
			if tNext != nil && (raw.IsSymbol(tNext, "(") || raw.IsParensGroup(tNext)) {
				// use contexts to check if there is no space between
				thisCtx := t.Context()
				parensCtx := tNext.Context()
				followedByParens = thisCtx.IsConsecutive(parensCtx)
			}

			followedByBrackets := false
			if tNext != nil && (raw.IsSymbol(tNext, "[") || raw.IsBracketsGroup(tNext)) {
				thisCtx := t.Context()
				bracketsCtx := tNext.Context()
				followedByBrackets = thisCtx.IsConsecutive(bracketsCtx)
			}

			followedByNew := tNext != nil && raw.IsSymbol(tNext, ":=")

			if !(followedByParens || followedByBrackets || followedByNew) {
				return raw.NewWordLiteralString(w.Value(), w.Context())
			} else {
				return t
			}
		}
	} else {
		return t
	}
}

// tNext is needed to check if word is followed by parens (
func (p *UIParser) convertWords(ts []raw.Token) []raw.Token {

	res := make([]raw.Token, 0)

	for i := 0; i < len(ts); i++ {
		var tNext raw.Token = nil
		if i < len(ts)-1 {
			tNext = ts[i+1]
		}
		t := p.convertWord(ts[i], tNext)
		res = append(res, t)
	}

	return res
}

func (p *UIParser) convertWordsNested(t_ raw.Token, tNext raw.Token) raw.Token {
	switch t := t_.(type) {
	case *raw.Operator:
		args := t.Args()
		for i, arg := range args {
			var argNext raw.Token = nil
			if i < len(args)-1 {
				argNext = args[i+1]
			}

			args[i] = p.convertWordsNested(arg, argNext)
		}
		return t
	case *raw.Group:
		for _, field := range t.Fields {
			for i, sub := range field {
				var subNext raw.Token = nil
				if i < len(field)-1 {
					subNext = field[i+1]
				}

				field[i] = p.convertWordsNested(sub, subNext)
			}
		}
		return t
	default:
		return p.convertWord(t_, tNext)
	}
}

func (p *UIParser) buildToken(ts []raw.Token, convert bool) (html.Token, error) {
	// do conversion now, or do it before
	if convert {
		for i, t := range ts {
			var tNext raw.Token = nil
			if i < len(ts)-1 {
				tNext = ts[i+1]
			}
			ts[i] = p.convertWordsNested(t, tNext)
		}
	}

	return p.helper.buildToken(ts)
}

func (p *UIParser) parseRawValue(s string, ctx context.Context) ([]raw.Token, error) {
	fp, err := NewFormulaParser(s, ctx)
	if err != nil {
		return nil, err
	}

	ts, err := fp.Parser.tokenizeFlat()
	if err != nil {
		return nil, err
	}

	// turn words not beginning with a dollar sign into strings except expects when directly followed by brackets
	// remove dollar signs from remaining words
	for i := 0; i < len(ts); i++ {
		t := ts[i]

		if raw.IsAnyWord(t) {
			w, err := raw.AssertWord(t)
			if err != nil {
				panic(err)
			}

			str := w.Value()
			if strings.HasPrefix(str, "$") {
				ts[i] = raw.NewValueWord(str[1:], w.Context())
			} else {
				followedByParens := false
				if i < len(ts)-1 && raw.IsSymbol(ts[i+1], "(") {
					// use contexts to check if there is no space between
					thisCtx := t.Context()
					parensCtx := ts[i+1].Context()
					followedByParens = thisCtx.IsConsecutive(parensCtx)
				}

				if !followedByParens {
					ts[i] = raw.NewWordLiteralString(str, w.Context())
				}
			}
		}
	}

	// now nest the groups and the operators
	ts, err = fp.nestGroups(ts)
	if err != nil {
		return nil, err
	}

	return fp.nestOperators(ts)
}

func (p *UIParser) buildDirectiveAttributeValue(r [2]int,
	fnPre func(ts []raw.Token) []raw.Token) (html.Token, error) {
	//s := p.writeWithoutLineContinuation(r[0], r[1])
	//fmt.Printf("value from %d to %d: \"%s\"\n", r[0], r[1], s)

	//ctx := p.NewContext(r[0], r[1])

	ts, err := p.tokenizePartial(r[0], r[1], fnPre)
	if err != nil {
		return nil, err
	}

	return p.buildToken(ts, false)
}

// eat up to next equals sign (preceded by whitespace or word char)
// return prev value range, and range for next attr key
// either can be zero, if both are zero then there are no more attributes
func (p *UIParser) matchNextDirectiveAttribute(end int) ([2]int, [2]int, error) {
	// all symbols expect more
	// containers need to be matched
	// preunary symbols are a difficulty and require an extra check
	containerCount := 0

	start := -1 // first non-white initializes the start
	expectsMore := false

	prevC := ' '
	prevIsNonOp := false

	ignoreNL := true
	if end < 0 {
		ignoreNL = false
		end = p.Len()
	}

	for ; p.pos < end; p.pos++ {
		c := p.raw[p.pos]

		// this is ok for comments, but bad for strings
		if p.mask[p.pos] == COMMENT && !(c == '\n' || c == '\r') {
			prevC = c
			continue
		}

		if ignoreNL && (c == '\n' || c == '\r') {
			c = ' '
		}

		if c != ' ' && c != '\t' && c != '\n' && c != '\r' && start == -1 {
			start = p.pos
		}

		if p.mask[p.pos] == STRING {
			prevC = c
			continue
		} else if c == '(' || c == '{' || c == '[' {
			if containerCount == 0 && !expectsMore && !prevIsNonOp && p.pos != start {
				p.pos -= 1
				return [2]int{start, p.pos}, [2]int{p.pos, p.pos}, nil
			}
			prevIsNonOp = true
			containerCount += 1
		} else if c == ')' || c == '}' || c == ']' {
			prevIsNonOp = true
			containerCount -= 1
			if containerCount < 0 {
				errCtx := p.NewContext(p.pos, p.pos+1)
				return [2]int{0, 0}, [2]int{0, 0}, errCtx.NewError("Error: unmatched container end")
			} else {
				prevC = c
				continue
			}
		} else if containerCount != 0 {
			prevIsNonOp = false
			prevC = c
			continue
		} else if c == ' ' || c == '\t' || ((c == '\n' || c == '\r') && expectsMore) {
			if prevIsNonOp {
				expectsMore = false
			}
			prevIsNonOp = false
			prevC = c
			continue
		} else if c == '>' || c == '<' || c == '+' || c == '*' || c == '/' || (c == '=' && (prevC == ':' || prevC == '!' || prevC == '<' || prevC == '>' || prevC == '=')) {
			prevIsNonOp = false
			expectsMore = true
		} else if c == '\\' {
			prevIsNonOp = false
			expectsMore = true
		} else if c == '-' && p.pos < p.Len()-1 {
			prevIsNonOp = false
			if p.pos < p.Len()-1 {
				nextC := p.raw[p.pos+1]
				if nextC == ' ' || nextC == '\t' || nextC != '\n' {
					expectsMore = true
				} else if start != -1 && !expectsMore && (prevC == ' ' || prevC == '\t') {
					p.pos -= 1

					return [2]int{start, p.pos}, [2]int{p.pos, p.pos}, nil
				}
			} else {
				errCtx := p.NewContext(p.pos, p.pos+1)
				return [2]int{0, 0}, [2]int{0, 0}, errCtx.NewError("Error: stray negation")
			}
		} else if c == '=' {
			prevIsNonOp = false
			if p.pos < p.Len()-1 {
				nextC := p.raw[p.pos+1]
				if nextC == '=' {
					prevC = c
					continue
				} else {
					rKey, err := p.eatPrevWord()
					if err != nil {
						return [2]int{0, 0}, [2]int{0, 0}, err
					}

					p.pos += 1

					if start != rKey[0] {
						errCtx := p.NewContext(start, rKey[1])
						return [2]int{0, 0}, [2]int{0, 0}, errCtx.NewError("Error: incomplete value")
					}

					return [2]int{rKey[0], rKey[0]}, rKey, nil
				}
			} else {
				errCtx := p.NewContext(p.pos, p.pos+1)
				return [2]int{0, 0}, [2]int{0, 0}, errCtx.NewError("Error: stray attribute assignment")
			}
		} else if !expectsMore && (c == '\n' || c == '\r') {
			if start == -1 {
				start = p.pos
			}
			end := p.pos
			if start == end {
				p.pos += 1
				if p.pos < p.Len()-1 && p.raw[p.pos+1] == '\r' {
					p.pos += 1
				}
			}
			return [2]int{start, end}, [2]int{end, end}, nil
		} else if c == ',' || c == ';' {
			prevIsNonOp = false
		} else {
			// all other chars are part of variables
			if !expectsMore && !prevIsNonOp && p.pos != start {
				p.pos -= 1
				return [2]int{start, p.pos}, [2]int{p.pos, p.pos}, nil
			}

			expectsMore = false
			prevIsNonOp = true
		}

		prevC = c
	}

	if start == -1 {
		start = p.pos
	}

	return [2]int{start, p.pos}, [2]int{p.pos, p.pos}, nil
}

func (p *UIParser) eatPrevWord() ([2]int, error) {
	init := p.pos

	start := -1

	p.pos--
	for ; p.pos >= 0; p.pos-- {
		c := p.raw[p.pos]

		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\\' {
			if start != -1 {
				r := [2]int{p.pos + 1, start}
				p.pos = init
				return r, nil
			} else {
				continue
			}
		} else if c == '+' || c == '*' || c == '/' || c == '!' || c == '=' {
			errCtx := p.NewContext(p.pos, p.pos+1)
			return [2]int{0, 0}, errCtx.NewError("Error: unexpected symbol")
		} else {
			// alphanum char
			if start == -1 {
				start = p.pos + 1
			}
		}
	}

	errCtx := p.NewContext(p.pos, init)
	return [2]int{0, 0}, errCtx.NewError("Error: unable to find attribute key")
}

// return indent from first non-empty line
func (p *UIParser) eatWhitespace() (int, bool) {
	nlPos := p.pos

	for ; p.pos < p.Len(); p.pos++ {
		c := p.raw[p.pos]

		if c == '\n' || c == '\r' {
			nlPos = p.pos + 1
		} else if c != ' ' && c != '\t' && p.mask[p.pos] != COMMENT {
			indent := p.pos - nlPos
			return indent, false
		} else {
			continue
		}
	}

	return 0, true
}

func (p *UIParser) eatNonWhitespace() ([2]int, bool) {
	start := -1
	end := -1
	isString := true

	for ; p.pos < p.Len(); p.pos++ {
		c := p.raw[p.pos]

		if p.mask[p.pos] != STRING && (c == '\n' || c == '\r' || c == ' ' || c == '\t' || p.mask[p.pos] == COMMENT || c == '(') {
			if start != -1 {
				end = p.pos
				break
			}

			continue
		} else {
			if p.mask[p.pos] != STRING {
				isString = false
			}

			if start == -1 {
				start = p.pos
			}
		}
	}

	if start == -1 {
		//panic("cant be -1")
		start = p.pos
		end = p.pos
		isString = false
	}

	return [2]int{start, end}, isString
}

func (p *UIParser) eatRestOfLine() bool {
	isNotEmpty := false

	for ; p.pos < p.Len(); p.pos++ {
		c := p.raw[p.pos]

		if c == '\n' {
			return isNotEmpty
		} else if c == '\r' || c == ' ' || c == '\t' || p.mask[p.pos] == COMMENT {
			continue
		} else {
			isNotEmpty = true
		}
	}

	return isNotEmpty
}

func (p *UIParser) DumpTokens() {
	fmt.Println("\nUI tokens:")
	fmt.Println("===========")

	tags, err := p.BuildTags()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, tag := range tags {
		fmt.Println(tag.Dump(""))
	}
}
