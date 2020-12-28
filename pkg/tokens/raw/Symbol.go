package raw

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

// this is a temporary token, and shouldnt be available outside the parser package
type Symbol struct {
	value     string
	isKeyWord bool
	isDummy   bool // for >> and > after a >>> or >>
	TokenData
}

func NewSymbol(value string, isKeyWord bool, ctx context.Context) *Symbol {
	return &Symbol{value, isKeyWord, false, TokenData{ctx}}
}

func NewDummySymbol(value string, ctx context.Context) *Symbol {
	return &Symbol{value, false, true, TokenData{ctx}}
}

func (t *Symbol) Value() string {
	return t.value
}

func IsAnySymbol(t Token) bool {
	_, ok := t.(*Symbol)
	return ok
}

func IsAnyNonWordSymbol(t Token) bool {
	if sym, ok := t.(*Symbol); ok {
		return !patterns.IsSimpleWord(sym.Value())
	}

	return false
}

func IsSymbol(t Token, value string) bool {
	symbol, ok := t.(*Symbol)
	if !ok || symbol.value != value {
		return false
	} else {
		return true
	}
}

func IsDummySymbol(t Token) bool {
	if symbol, ok := t.(*Symbol); ok {
		return symbol.isDummy
	} else {
		return false
	}
}

func IsSymbolThatEndsWith(t Token, value string) bool {
	symbol, ok := t.(*Symbol)
	if ok {
		return strings.HasSuffix(symbol.value, value)
	} else {
		return false
	}
}

func AssertSymbol(t Token, value string) (*Symbol, error) {
	symbol, ok := t.(*Symbol)
	if !ok || symbol.value != value {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected \"" + value + "\"")
		return nil, err
	}

	return symbol, nil
}

func AssertAnySymbol(t Token) (*Symbol, error) {
	symbol, ok := t.(*Symbol)
	if !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected a symbol")
	}

	return symbol, nil
}

func AssertNotSymbol(t Token) error {
	if s, ok := t.(*Symbol); ok {
		errCtx := t.Context()
		return errCtx.NewError("Error: unexpected symbol '" + s.value + "'")
	}

	return nil
}

func ChangeSymbol(t Token, newValue string, newIsKeyWord bool) {
	if s, ok := t.(*Symbol); ok {
		s.value = newValue
		s.isKeyWord = newIsKeyWord
	} else {
		panic("expected *Symbol")
	}
}

func (t *Symbol) Dump(indent string) string {
	return indent + "Symbol " + t.value + "\n"
}

func (a *Symbol) IsSame(other Token) bool {
	if b, ok := other.(*Symbol); ok {
		return a.value == b.value
	} else {
		return false
	}
}

func ContainsSymbol(ts []Token, s string) bool {
	for _, t := range ts {
		if IsSymbol(t, s) {
			return true
		}
	}

	return false
}

func ContainsSymbolThatEndsWith(ts []Token, s string) bool {
	for _, t := range ts {
		if IsSymbolThatEndsWith(t, s) {
			return true
		}
	}

	return false
}

// first part includes the symbol
func SplitByFirstSymbol(ts []Token, s string) ([]Token, []Token, error) {
	for i, t := range ts {
		if IsSymbol(t, s) {
			if i == len(ts)-1 {
				return ts[0 : i+1], []Token{}, nil
			} else {
				return ts[0 : i+1], ts[i+1:], nil
			}
		}
	}

	errCtx := MergeContexts(ts...)
	return nil, nil, errCtx.NewError("Error: symbol '" + s + "' not found")
}
