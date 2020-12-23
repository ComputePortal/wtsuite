package parsers

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/raw"
)

var VERBOSITY = 0

type RuneMask int

const (
	NONE RuneMask = iota
	COMMENT
	STRING
	FORMULA
	WORD_OR_LITERAL
	SYMBOL
)

type Parser struct {
	pos  int // equivalent of Seek()/Peek() position
	raw  []rune
	mask []RuneMask

	settings ParserSettings
	ctx      context.Context
}

func newParser(raw string, settings ParserSettings, ctx context.Context) Parser {
	n := len(raw)

	return Parser{0, []rune(raw), make([]RuneMask, n), settings, ctx}
}

func (p *Parser) NewContext(start, stop int) context.Context {
	if stop == -1 {
		stop = p.Len()
	}

	return p.ctx.NewContext(start, stop)
}

func (p *Parser) refine(start, stop int) Parser {
	if stop == -1 {
		stop = p.Len()
	}

	return Parser{0, p.raw[start:stop], p.mask[start:stop], p.settings, p.NewContext(start, stop)}
}

func (p *Parser) changeCaller(caller string) Parser {
	return Parser{p.pos, p.raw[:], p.mask[:], p.settings, p.ctx.ChangeCaller(caller)}
}

func (p *Parser) NewError(start, stop int, msg string) error {
	errCtx := p.NewContext(start, stop)
	return errCtx.NewError(msg)
}

func (p *Parser) Reset() {
	p.SeekToStart()
	for i, m := range p.mask {
		if m == WORD_OR_LITERAL || m == SYMBOL {
			p.mask[i] = NONE
		}
	}
}

func (p *Parser) Len() int {
	return len(p.raw)
}

func (p *Parser) isWhiteSpace(i int) bool {
	r := p.raw[i]

	return r == 10 || r == 13 || r == 9 || r == 32
}

func (p *Parser) isEmpty(start, end int) bool {
	if p.Len() == 0 {
		return true
	} else {
		for i := start; i < end; i++ {
			if !(p.mask[i] == COMMENT || p.isWhiteSpace(i)) {
				return false
			}
		}

		return true
	}
}

func (p *Parser) IsEmpty() bool {
	return p.isEmpty(0, p.Len())
}

func (p *Parser) AssertEmpty() error {
	if !p.IsEmpty() {
		errCtx := p.NewContext(0, -1)
		return errCtx.NewError("Error: unexpected content")
	} else {
		return nil
	}
}

func (p *Parser) assertEmpty(start, end int) error {
	if !p.isEmpty(start, end) {
		errCtx := p.NewContext(start, end)
		return errCtx.NewError("Error: unexpected content")
	} else {
		return nil
	}
}

// regexp match interface function
func (p *Parser) ReadRune() (rune, int, error) {
	if p.pos >= p.Len() {
		return 0, 1, errors.New("EOF")
	}

	c := p.raw[p.pos]

	if p.mask[p.pos] != NONE {
		c = 32 // space
	}
	p.pos += 1

	return c, 1, nil
}

func (p *Parser) SeekToStart() {
	p.pos = 0
}

func (p *Parser) SetMask(start, stop int, mask RuneMask) {
	n := stop
	if n == -1 {
		n = p.Len()
	}

	for i := start; i < stop; i++ {
		p.mask[i] = mask
	}
}

// used by BuildText, quotes/ticks are left unchanged
func (p *Parser) Write(start, stop int) string {
	if stop == -1 {
		stop = p.Len()
	}

	// remove the comment parts
	var b strings.Builder
	if start == -1 {
		return ""
	}

	for i := start; i < stop; i++ {
		if p.mask[i] != COMMENT {
			b.WriteRune(p.raw[i])
		}
	}

	return b.String()
}

func (p *Parser) writeWithoutFormulasAndStrings(start, stop int) string {
	if stop == -1 {
		stop = p.Len()
	}

	// remove the comment parts
	var b strings.Builder
	for i := start; i < stop; i++ {
		if p.mask[i] != COMMENT && p.mask[i] != FORMULA && p.mask[i] != STRING {
			b.WriteRune(p.raw[i])
		}
	}

	return b.String()
}

// the returned string is the match
func (p *Parser) nextMatch(re *regexp.Regexp, removeFormulasAndStrings bool) ([2]int, string, bool) {
	// matched indices are relative to initial pos
	pos := p.pos

	// parser implements io.Reader interface
	relIndices := re.FindReaderIndex(p)

	if relIndices == nil {
		// only change the position if we have a match
		p.pos = pos
		return [2]int{0, 0}, "", false
	} else {
		start := relIndices[0] + pos
		stop := relIndices[1] + pos
		p.pos = stop
		// also return the matched string
		if removeFormulasAndStrings {
			return [2]int{start, stop}, p.writeWithoutFormulasAndStrings(start, stop), true
		} else {
			return [2]int{start, stop}, p.Write(start, stop), true
		}
	}
}

func (p *Parser) nextGroupStopMatch(group patterns.Group, trackStarts bool) ([2]int, bool) {
	// assume we just detected an open, and want the next close
	// the regexp detects both
	count := 0
	for true {
		// if group.IsTagGroup() == true => strings and formulas are returned as spaces in 's'
		//  because they might contain symbols that confuse of the match
		if r, s, ok := p.nextMatch(group.StartStopRegexp(), group.IsTagGroup()); ok {
			switch {
			case group.MatchStop(s):
				if count == 0 {
					return r, true
				} else {
					count--
				}
			case group.MatchStart(s):
				if trackStarts {
					count++
				}
			default:
				panic("bad regexp")
			}
		} else {
			break
		}
	}

	return [2]int{0, 0}, false
}

func (p *Parser) maskQuoted() error {
	quotedGroupMap := make(map[string]quotedGroupSettings)
	assertStopMatchMap := make(map[string]quotedGroupSettings)

	for _, v := range p.settings.quotedGroups.groups {
		startKey := v.groupPattern.Start()

		quotedGroupMap[startKey] = v

		if v.assertStopMatch {
			stopKey := v.groupPattern.Stop()
			assertStopMatchMap[stopKey] = v
		}
	}

	for true {
		if r, s, ok := p.nextMatch(p.settings.quotedGroups.pattern, false); ok {
			if group, ok := quotedGroupMap[s]; ok {
				if rr, ok := p.nextGroupStopMatch(group.groupPattern, group.trackStarts); ok {
					p.SetMask(r[0], rr[1], group.maskType)
				} else {
					return p.NewError(r[0], r[1], "Error: unmatched "+group.info)
				}
			} else if group, ok := assertStopMatchMap[s]; ok {
				return p.NewError(r[0], r[1], "Error: unmatched stop of "+group.info)
			} else {
				panic("quotedGroups.pattern/quotedGroups.groups inconsistency")
			}
		} else {
			break
		}
	}

	// reset
	p.SeekToStart()

	return nil
}

func (p *Parser) tokenizeQuoted(ts []tokens.Token, maskType RuneMask, fn func(start, stop int) ([]tokens.Token, error)) ([]tokens.Token, error) {
	appendToken := func(start, stop int) error {
		t, err := fn(start, stop)

		if err != nil {
			return err
		}

		ts = append(ts, t...)
		return nil
	}

	b := false
	start := -1

	for i, x := range p.mask {
		if x == maskType {
			if !b {
				start = i
				b = true
			}
		} else {
			if b {
				b = false

				if err := appendToken(start, i); err != nil {
					return ts, err
				}
			}
		}
	}

	if b {
		if err := appendToken(start, p.Len()); err != nil {
			return ts, err
		}
	}

	p.SeekToStart()

	return ts, nil
}

func (p *Parser) tokenizeStrings(ts []tokens.Token) ([]tokens.Token, error) {
	return p.tokenizeQuoted(ts, STRING, func(start, stop int) ([]tokens.Token, error) {
		s := p.Write(start+1, stop-1)

		strToken, err := tokens.NewLiteralString(s, p.NewContext(start, stop))
		return []tokens.Token{strToken}, err
	})
}

func (p *Parser) tokenizeFormulas(ts []tokens.Token) ([]tokens.Token, error) {
	return p.tokenizeQuoted(ts, FORMULA, func(start, stop int) ([]tokens.Token, error) {
		s := p.Write(start+1, stop-1)

		return p.settings.formulas.tokenizer(s, p.NewContext(start+1, stop-1))
	})
}

func (p *Parser) tokenizeWordsAndLiterals(ts []tokens.Token) ([]tokens.Token, error) {
	for true {
		if r, s, ok := p.nextMatch(p.settings.wordsAndLiterals.pattern, false); ok {
			t, err := p.settings.wordsAndLiterals.tokenizer(s, p.NewContext(r[0], r[1]))
			if err != nil {
				return ts, err
			}
			p.SetMask(r[0], r[1], p.settings.wordsAndLiterals.maskType)
			ts = append(ts, t)
		} else {
			break
		}
	}

	p.SeekToStart()

	return ts, nil
}

func (p *Parser) tokenizeSymbols(ts []tokens.Token) []tokens.Token {
	for true {
		if r, s, ok := p.nextMatch(p.settings.symbols.pattern, false); ok {
			p.SetMask(r[0], r[1], p.settings.symbols.maskType)
			ts = append(ts, tokens.NewSymbol(s, false, p.NewContext(r[0], r[1])))
		} else {
			break
		}
	}

	p.SeekToStart()

	return ts
}

func (p *Parser) sortTokens(ts []tokens.Token) {
	sort.Slice(ts, func(i, j int) bool {
		ci := ts[i].Context()
		cj := ts[j].Context()
		return ci.Less(&cj)
	})
}

func (p *Parser) assertNoStrayCharacters() error {
	// check for stray characters
	ctxs := []context.Context{}
	for i, m := range p.mask {
		if m == NONE && !p.isWhiteSpace(i) {
			ctxs = append(ctxs, p.NewContext(i, i+1))
		}
	}

	if len(ctxs) > 0 {
		ctx := context.MergeContexts(ctxs...)
		return ctx.NewError("Syntax Error: stray characters")
	}

	return nil
}

// proc can be used to handle unary/binary operators
func (p *Parser) nestGroups(ts []tokens.Token) ([]tokens.Token, error) {
	result := make([]tokens.Token, 0)

	isrc := 0 // position in ts
	for isrc < len(ts) {
		t := ts[isrc]
		if tokens.IsAnyNonAngledGroupStart(t) {
			r, ok := tokens.FindGroupStop(ts, isrc+1, t)
			if !ok {
				errCtx := t.Context()
				return result, errCtx.NewError("Syntax Error: unmatched container start")
			}

			sub, err := p.nestGroups(ts[r[0]:r[1]])
			if err != nil {
				return result, err
			}

			res, err := tokens.NewGroupFromTokens(tokens.Concat(t, sub, ts[r[1]]))
			if err != nil {
				return result, err
			}

			result = append(result, res)

			isrc = r[1] + 1
		} else if tokens.IsSymbol(t, patterns.ANGLED_START) && p.settings.tmpGroupAngled {
			tsinner, istop, ok := tokens.SuggestAngledGroup(ts, isrc+1, t)
			if !ok {
				result = append(result, t)
				// not actually an error, just for debugging
				//errCtx := t.Context()
				//panic(errCtx.NewError("Syntax Error: unmatched container start")) // for debugging
				isrc++
			} else {
				sub, err := p.nestGroups(tsinner)
				if err != nil {
					return result, err
				}

				res, err := tokens.NewGroupFromTokens(tokens.Concat(t, sub, ts[istop]))
				if err != nil {
					return result, err
				}

				result = append(result, res)

				isrc = istop + 1
			}
		} else if tokens.IsAnyNonAngledGroupStop(t) {
			// ANGLED_STOP (i.e. >) can also be for math
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: unmatched group")
		} else {
			result = append(result, t)

			isrc++
		}
	}

	return result, nil
}

func (p *Parser) nestOperatorRightToLeft(ts []tokens.Token, osm *operatorSettingsMap) ([]tokens.Token, error) {
	result := make([]tokens.Token, 0)

	// scan right to left
	n := len(ts)
	isrc := n - 1

	prependSingular := func(name string, ctx context.Context) {
		result = append([]tokens.Token{tokens.NewSingularOperator(name, ctx)}, result...)
	}

	prependUnary := func(name string, arg tokens.Token, ctx context.Context) {
		result = append([]tokens.Token{tokens.NewUnaryOperator(name, arg, ctx)}, result...)
		isrc--
	}

	replaceLastWithUnary := func(name string, arg tokens.Token, ctx context.Context) {
		result[0] = tokens.NewUnaryOperator(name, arg, ctx)
	}

	replaceLastWithBinary := func(name string, a tokens.Token, b tokens.Token, ctx context.Context) {
		result[0] = tokens.NewBinaryOperator(name, a, b, ctx)
		isrc -= 1
	}

	replaceLastWithTernary := func(name string, a tokens.Token, b tokens.Token,
		c tokens.Token, ctx context.Context) {
		result[0] = tokens.NewTernaryOperator(name, a, b, c, ctx)
		isrc -= 3
	}

	for isrc >= 0 {
		t := ts[isrc]
		ctx := t.Context()

		if osm.Has(t) {
			if isrc == 0 {
				if isrc == n-1 {
					if name, ok := osm.GetSingular(t); ok {
						prependSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else {
					right := result[0]

					if !tokens.IsOperable(right) {
						if name, ok := osm.GetSingular(t); ok {
							prependSingular(name, ctx)
						} else {
							return nil, ctx.NewError("Error: not a singular operator")
						}
					} else if name, ok := osm.GetPreUnary(t); ok {
						replaceLastWithUnary(name, right, ctx)
					} else {
						return nil, ctx.NewError("Error: not a pre-unary operator")
					}
				}
			} else if isrc == n-1 {
				left := ts[isrc-1]

				if !tokens.IsOperable(left) {
					if name, ok := osm.GetSingular(t); ok {
						prependSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else if name, ok := osm.GetPostUnary(t); ok {
					prependUnary(name, left, ctx)
				} else if tokens.IsAnyWord(t) {
					// for js "new", TODO: also for js "typeof" etc.
					result = append([]tokens.Token{t}, result...)
					isrc--
					continue
				} else {
					return nil, ctx.NewError("Error: not a post-unary operator")
				}
			} else {
				left := ts[isrc-1]
				right := result[0]

				if !tokens.IsOperable(left) && !tokens.IsOperable(right) {
					if name, ok := osm.GetSingular(t); ok {
						prependSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else if !tokens.IsOperable(left) {
					if name, ok := osm.GetPreUnary(t); ok {
						replaceLastWithUnary(name, right, ctx)
					}
				} else if !tokens.IsOperable(right) {
					if name, ok := osm.GetPostUnary(t); ok {
						prependUnary(name, left, ctx)
					}
				} else {
					isTernary := false
					if isrc > 2 {
						if name, ok := osm.GetTernaryR2L(ts[isrc-2], t); ok {
							replaceLastWithTernary(name, ts[isrc-3], left, right, ctx)
							isTernary = true
						}
					}

					if !isTernary {
						if name, ok := osm.GetBinary(t); ok {
							replaceLastWithBinary(name, left, right, ctx)
						} else {
							// to be handled by another case (eg. neg vs sub)
							result = append([]tokens.Token{t}, result...)
						}
					}
				}
			}
		} else {
			result = append([]tokens.Token{t}, result...)
		}

		isrc--
	}

	return result, nil
}

func (p *Parser) nestOperatorLeftToRight(ts []tokens.Token, osm *operatorSettingsMap) ([]tokens.Token, error) {
	result := make([]tokens.Token, 0)

	// scan right to left
	n := len(ts)
	isrc := 0

	appendSingular := func(name string, ctx context.Context) {
		result = append(result, tokens.NewSingularOperator(name, ctx))
	}

	appendUnary := func(name string, arg tokens.Token, ctx context.Context) {
		result = append(result, tokens.NewUnaryOperator(name, arg, ctx))
		isrc++
	}

	replaceLastWithUnary := func(name string, arg tokens.Token, ctx context.Context) {
		result[len(result)-1] = tokens.NewUnaryOperator(name, arg, ctx)
	}

	replaceLastWithBinary := func(name string, a tokens.Token, b tokens.Token, ctx context.Context) {
		result[len(result)-1] = tokens.NewBinaryOperator(name, a, b, ctx)
		isrc += 1
	}

	replaceLastWithTernary := func(name string, a tokens.Token, b tokens.Token, c tokens.Token, ctx context.Context) {
		result[len(result)-1] = tokens.NewTernaryOperator(name, a, b, c, ctx)
		isrc += 3
	}

	for isrc < n {
		t := ts[isrc]
		ctx := t.Context()

		if osm.Has(t) {
			if isrc == n-1 {
				if isrc == 0 {
					if name, ok := osm.GetSingular(t); ok {
						appendSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else {
					left := result[len(result)-1]

					if !tokens.IsOperable(left) {
						if name, ok := osm.GetSingular(t); ok {
							appendSingular(name, ctx)
						} else {
							return nil, ctx.NewError("Error: not a singular operator")
						}
					} else if name, ok := osm.GetPreUnary(t); ok {
						replaceLastWithUnary(name, left, ctx)
					} else {
						return nil, ctx.NewError("Error: not a pre-unary operator")
					}
				}
			} else if isrc == 0 {
				right := ts[isrc+1]

				if !tokens.IsOperable(right) {
					if name, ok := osm.GetSingular(t); ok {
						appendSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else if name, ok := osm.GetPostUnary(t); ok {
					appendUnary(name, right, ctx)
				} else {
					return nil, ctx.NewError("Error: not a post-unary operator")
				}
			} else {
				right := ts[isrc+1]
				left := result[len(result)-1]

				if !tokens.IsOperable(left) && !tokens.IsOperable(right) {
					if name, ok := osm.GetSingular(t); ok {
						appendSingular(name, ctx)
					} else {
						return nil, ctx.NewError("Error: not a singular operator")
					}
				} else if !tokens.IsOperable(right) {
					if name, ok := osm.GetPostUnary(t); ok {
						replaceLastWithUnary(name, left, ctx)
					}
				} else if !tokens.IsOperable(left) {
					if name, ok := osm.GetPreUnary(t); ok {
						appendUnary(name, right, ctx)
					}
				} else {
					isTernary := false
					if isrc < n-3 {
						if name, ok := osm.GetTernaryL2R(t, ts[isrc+2]); ok {
							replaceLastWithTernary(name, left, right, ts[isrc+3], ctx)
							isTernary = true
						}
					}

					if !isTernary {
						if name, ok := osm.GetBinary(t); ok {
							replaceLastWithBinary(name, left, right, ctx)
						} else {
							// to be handled by another case (eg. neg vs sub)
							result = append(result, t)
						}
					}
				}
			}
		} else {
			result = append(result, t)
		}

		isrc++
	}

	return result, nil
}

func (p *Parser) isNestableToken(i int, t tokens.Token) bool {
	if tokens.IsAnySymbol(t) { // some symbols can also be words, so check before IsAnyWord
		if p.settings.tmpGroupPeriods && tokens.IsSymbol(t, patterns.PERIOD) {
			return false
		}

		if p.settings.tmpGroupArrows && tokens.IsSymbol(t, patterns.ARROW) {
			return false
		}

		if p.settings.tmpGroupDColons && tokens.IsSymbol(t, patterns.DCOLON) {
			return false
		}

		return true
	} else if tokens.IsAnyWord(t) {
		return (i > 0) && !p.settings.tmpGroupWords
	} else {
		return false
	}
}

func (p *Parser) nestOperators(ts []tokens.Token) ([]tokens.Token, error) {
	if p.settings.recursivelyNestOperators {
		for _, t := range ts {
			if group, ok := t.(*tokens.Group); ok {
				for j, field := range group.Fields {
					newField, err := p.nestOperators(field)
					if err != nil {
						return nil, err
					}
					group.Fields[j] = newField
				}
			}
		}
	}

	// compact consecutive []() ...
	compact := make([]tokens.Token, 0)
	tmp := make([]tokens.Token, 0)

	for i, t := range ts {
		if p.isNestableToken(i, t) {
			if len(tmp) != 0 {
				if len(tmp) > 1 {
					compact = append(compact, tokens.NewTmpGroup(tmp))
				} else {
					compact = append(compact, tmp[0])
				}
				tmp = make([]tokens.Token, 0)
			}

			compact = append(compact, t)
		} else {
			tmp = append(tmp, t)
		}
	}
	if len(tmp) != 0 {
		if len(tmp) > 1 {
			compact = append(compact, tokens.NewTmpGroup(tmp))
		} else {
			compact = append(compact, tmp[0])
		}
	}
	ts = compact

	var err error
	for _, op := range p.settings.operators.sortedOperators {
		if op.LeftToRight() {
			ts, err = p.nestOperatorLeftToRight(ts, &op)
			if err != nil {
				return nil, err
			}
		} else {
			ts, err = p.nestOperatorRightToLeft(ts, &op)
			if err != nil {
				return nil, err
			}
		}
	}

	return ts, nil
}

func (p *Parser) tokenizeFlat() ([]tokens.Token, error) {
	ts := make([]tokens.Token, 0)
	var err error

	ts, err = p.tokenizeStrings(ts)
	if err != nil {
		return nil, err
	}

	ts, err = p.tokenizeFormulas(ts)
	if err != nil {
		return nil, err
	}

	ts, err = p.tokenizeWordsAndLiterals(ts)
	if err != nil {
		return nil, err
	}

	ts = p.tokenizeSymbols(ts)

	p.sortTokens(ts)

	if err := p.assertNoStrayCharacters(); err != nil {
		return nil, err
	}

	return ts, nil
}

// tokenize everything
func (p *Parser) tokenize() ([]tokens.Token, error) {
	ts, err := p.tokenizeFlat()
	if err != nil {
		return nil, err
	}

	return p.nestGroups(ts)
}

// only expand once
func (p *Parser) expandTmpGroups(ts []tokens.Token) []tokens.Token {
	return tokens.ExpandTmpGroups(ts)
}

func (p *Parser) expandAngledGroups(ts []tokens.Token) []tokens.Token {
	return tokens.ExpandAngledGroups(ts)
}
