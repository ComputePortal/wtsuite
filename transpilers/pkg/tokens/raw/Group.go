package raw

import (
	"strings"

	"../context"
	"../patterns"
)

type GroupType int

const (
	PARENS GroupType = iota
	BRACKETS
	BRACES
	ANGLED
	TMP
	TEMPLATE
)

type ContentType int

const (
	EMPTY ContentType = iota
	SINGLE
	COMMA
	SEMICOLON
	TMPCONTENT
)

type Group struct {
	gt       GroupType
	st       ContentType
	original []Token // only used for angled
	Fields   [][]Token
	TokenData
}

func NewTmpGroup(raw []Token) *Group {
	ctx := MergeContexts(raw...)
	return &Group{TMP, TMPCONTENT, raw, [][]Token{raw}, TokenData{ctx}}
}

func NewTemplateGroup(ctx context.Context) *Group {
	return &Group{TEMPLATE, TMPCONTENT, []Token{}, [][]Token{}, TokenData{ctx}}
}

func NewGroupFromTokens(raw []Token) (*Group, error) {
	n := len(raw)

	if n < 2 {
		panic("not enough tokens to form a group")
	}

	start, err := AssertAnySymbol(raw[0])
	if err != nil {
		return nil, err
	}
	// assume stops are consistent
	if !isGroupStop(start, raw[n-1]) {
		panic("bad stop")
	}

	ctx := MergeContexts(start, raw[n-1])

	var gt GroupType
	switch start.value {
	case patterns.PARENS_START:
		gt = PARENS
	case patterns.BRACES_START:
		gt = BRACES
	case patterns.BRACKETS_START:
		gt = BRACKETS
	case patterns.ANGLED_START:
		gt = ANGLED
	default:
		panic("bad group start")
	}

	content := raw[1 : n-1]

	fields := make([][]Token, 0)
	buffer := make([]Token, 0)

	st := EMPTY

	// prefer SEMICOLON
	for i, t := range content {
		switch {
		case IsSymbol(t, patterns.SEMICOLON):
			switch st {
			case EMPTY:
				st = SEMICOLON
				fields = append(fields, buffer)
			case COMMA:
				st = SEMICOLON
				// flatten current fields
				fields = [][]Token{content[0:i]}
			default:
				fields = append(fields, buffer)
			}
			buffer = make([]Token, 0)
		case IsSymbol(t, patterns.COMMA):
			switch st {
			case EMPTY:
				st = COMMA
				fallthrough
			case COMMA:
				fields = append(fields, buffer)
				buffer = make([]Token, 0)
			case SEMICOLON:
				buffer = append(buffer, t)
			}
		default:
			buffer = append(buffer, t)
		}
	}

	if len(buffer) != 0 {
		switch st {
		case EMPTY:
			st = SINGLE
		}
		fields = append(fields, buffer)
	}

	return &Group{gt, st, raw, fields, TokenData{ctx}}, nil
}

func (t *Group) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	switch t.gt {
	case PARENS:
		b.WriteString("ParensGroup ")
	case BRACES:
		b.WriteString("BracesGroup ")
	case BRACKETS:
		b.WriteString("BracketsGroup ")
	case ANGLED:
		b.WriteString("AngledGroup ")
	case TMP:
		b.WriteString("TmpGroup ")
	case TEMPLATE:
		b.WriteString("TemplateGroup ")
	default:
		panic("unhandled")
	}

	switch t.st {
	case COMMA:
		b.WriteString(" , ")
	case SEMICOLON:
		b.WriteString(" ; ")
	case SINGLE, EMPTY:
		b.WriteString(" no sep ")
	case TMPCONTENT:
	default:
		panic("unhandled")
	}

	b.WriteString("\n")

	for _, field := range t.Fields {
		for _, v := range field {
			b.WriteString(v.Dump(indent + "  "))
		}
	}

	return b.String()
}

func IsAnyNonAngledGroupStart(t Token) bool {
	symbol, ok := t.(*Symbol)

	if ok {
		return (symbol.value == patterns.PARENS_START ||
			symbol.value == patterns.BRACES_START ||
			symbol.value == patterns.BRACKETS_START)
	} else {
		return false
	}
}

func IsAnyNonAngledGroupStop(t Token) bool {
	symbol, ok := t.(*Symbol)

	if ok {
		return (symbol.value == patterns.PARENS_STOP ||
			symbol.value == patterns.BRACES_STOP ||
			symbol.value == patterns.BRACKETS_STOP)
	} else {
		return false
	}
}

func IsGroup(t Token) bool {
	_, ok := t.(*Group)
	return ok
}

func AssertGroup(t Token) (*Group, error) {
	if gr, ok := t.(*Group); ok {
		return gr, nil
	}

	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected group")
}

func AssertParensGroup(t Token) (*Group, error) {
	gr, err := AssertGroup(t)
	if err != nil || !gr.IsParens() {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected (...)")
	}

	return gr, nil
}

func AssertBracesGroup(t Token) (*Group, error) {
	gr, err := AssertGroup(t)
	if err != nil || !gr.IsBraces() {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected {...}")
	}

	return gr, nil
}

func AssertBracketsGroup(t Token) (*Group, error) {
	gr, err := AssertGroup(t)
	if err != nil || !gr.IsBrackets() {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected [...]")
	}

	return gr, nil
}

func AssertAngledGroup(t Token) (*Group, error) {
	gr, err := AssertGroup(t)
	if err != nil || !gr.IsAngled() {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected <...>")
	}

	return gr, nil
}

func isGroupStart(ref Token, tok Token) bool {
	r, ok := ref.(*Symbol)
	if !ok {
		return false
	}

	t, ok := tok.(*Symbol)
	if !ok {
		return false
	}

	return r.value == t.value
}

func isGroupStop(ref Token, tok Token) bool {
	r, ok := ref.(*Symbol)
	if !ok {
		return false
	}

	t, ok := tok.(*Symbol)
	if !ok {
		return false
	}

	switch r.value {
	case patterns.PARENS_START:
		return t.value == patterns.PARENS_STOP
	case patterns.BRACES_START:
		return t.value == patterns.BRACES_STOP
	case patterns.BRACKETS_START:
		return t.value == patterns.BRACKETS_STOP
	case patterns.ANGLED_START:
		return t.value == patterns.ANGLED_STOP
	default:
		return false
	}
}

func FindGroupStop(ts []Token, istart int, tstart Token) ([2]int, bool) {
	// doesn't work for angled groups because >> and >>> are actual symbols
	count := 0
	for istop := istart; istop < len(ts); istop++ {
		switch {
		case isGroupStop(tstart, ts[istop]):
			if count == 0 {
				return [2]int{istart, istop}, true
			} else {
				count--
			}
		case isGroupStart(tstart, ts[istop]):
			count++
		}
	}

	return [2]int{istart, istart}, false
}

// returns inner, and istop, and success
func SuggestAngledGroup(ts []Token, istart int, tstart Token) ([]Token, int, bool) {
	count := 0
	for i := istart; i < len(ts); i++ {
		switch {
		case IsSymbol(ts[i], patterns.SEMICOLON): // stop immediately, angled group can never contain semicolon
			return nil, istart, false
		case IsSymbol(ts[i], patterns.ANGLED_STOP): // can be dummy
			if count == 0 {
				return ts[istart:i], i, true
			} else {
				count--
			}
		case IsSymbol(ts[i], patterns.ANGLED_STOP2): // can be dummy
			if count == 0 {
				return nil, istart, false
			} else if count == 1 {
				return append(ts[istart:i], NewDummySymbol(patterns.ANGLED_STOP, ts[i].Context())),
					i, true
			} else {
				count -= 2
			}
		case IsSymbol(ts[i], patterns.ANGLED_STOP3): // cannot be dummy
			if count == 0 || count == 1 {
				return nil, istart, false
			} else if count == 2 {
				return append(ts[istart:i], NewDummySymbol(patterns.ANGLED_STOP2, ts[i].Context())),
					i, true
			}
		case IsSymbol(ts[i], patterns.ANGLED_START):
			count++
		}
	}

	return nil, istart, false
}

func IsAnyGroup(t Token) bool {
	_, ok := t.(*Group)
	return ok
}

func IsParensGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsParens()
	}

	return false
}

func IsBracesGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsBraces()
	}

	return false
}

func IsBracketsGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsBrackets()
	}

	return false
}

func IsAngledGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsAngled()
	}

	return false
}

func IsTmpGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsTmp()
	}

	return false
}

func IsTemplateGroup(t Token) bool {
	if g, ok := t.(*Group); ok {
		return g.IsTemplate()
	}

	return false
}

func (t *Group) IsParens() bool {
	return t.gt == PARENS
}

func (t *Group) IsBraces() bool {
	return t.gt == BRACES
}

func (t *Group) IsBrackets() bool {
	return t.gt == BRACKETS
}

func (t *Group) IsAngled() bool {
	return t.gt == ANGLED
}

func (t *Group) IsEmpty() bool {
	return t.st == EMPTY
}

func (t *Group) IsSingle() bool {
	return t.st == SINGLE
}

func (t *Group) IsComma() bool {
	return t.st == COMMA
}

func (t *Group) IsTmp() bool {
	return t.gt == TMP
}

func (t *Group) IsTemplate() bool {
	return t.gt == TEMPLATE
}

func (t *Group) IsSemiColon() bool {
	return t.st == SEMICOLON
}

func (t *Group) FlattenCommas() ([]Token, error) {
	if t.IsSemiColon() {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: unexpected semicolons")
	}

	result := make([]Token, 0)

	if t.IsSingle() {
		result = t.Fields[0]
	} else if t.IsComma() {
		for i, f := range t.Fields {
			result = append(result, f...)

			if i < len(t.Fields)-1 {
				result = append(result, NewSymbol(patterns.COMMA, false, t.Context()))
			}
		}
	}

	return result, nil
}

func ExpandTmpGroups(ts []Token) []Token {
	result := make([]Token, 0)

	for _, t := range ts {
		if IsTmpGroup(t) {
			tmp, err := AssertGroup(t)
			if err != nil {
				panic(err)
			}

			// Fields[0] is not the same as original due to nesting
			result = append(result, tmp.Fields[0]...)
		} else {
			result = append(result, t)
		}
	}

	return result
}

func ExpandAngledGroups(ts []Token) []Token {
	result := make([]Token, 0)

	for _, t := range ts {
		if IsAngledGroup(t) {
			angled, err := AssertGroup(t)
			if err != nil {
				panic(err)
			}

			result = append(result, angled.original[0]) // opening <

			n := len(angled.original)
			result = append(result, ExpandAngledGroups(angled.original[1:n-1])...)
			last := angled.original[n-1]

			if !IsDummySymbol(last) {
				result = append(result, last)
			}
		} else {
			result = append(result, t)
		}
	}

	return result
}

func ExpandParensGroup(t Token) []Token {
	if IsParensGroup(t) {
		parens, err := AssertParensGroup(t)
		if err != nil {
			panic(err)
		}
		res := make([]Token, 0)

		for _, fs := range parens.Fields {
			res = append(res, fs...)
		}

		return res
	} else {
		return []Token{t}
	}
}
