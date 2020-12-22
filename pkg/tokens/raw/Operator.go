package raw

import (
	"strings"

	"../context"
)

type Operator struct {
	name string
	args []Token
	TokenData
}

func NewSingularOperator(name string, ctx context.Context) *Operator {
	return &Operator{name, []Token{}, TokenData{ctx}}
}

func NewUnaryOperator(name string, a Token, ctx context.Context) *Operator {
	return &Operator{name, []Token{a}, TokenData{ctx}}
}

func NewBinaryOperator(name string, a Token, b Token, ctx context.Context) *Operator {
	return &Operator{name, []Token{a, b}, TokenData{ctx}}
}

func NewTernaryOperator(name string, a Token, b Token, c Token, ctx context.Context) *Operator {
	return &Operator{name, []Token{a, b, c}, TokenData{ctx}}
}

func (t *Operator) Name() string {
	return t.name
}

func (t *Operator) Args() []Token {
	return t.args
}

func (t *Operator) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Operator ")
	b.WriteString(t.name)
	b.WriteString("\n")
	for _, arg := range t.args {
		b.WriteString(arg.Dump(indent + "  "))
	}

	return b.String()
}

func IsOperator(t Token, name string) bool {
	if op, ok := t.(*Operator); ok {
		return op.name == name
	}
	return false
}

func ChangeOperator(t Token, newName string) {
	if op, ok := t.(*Operator); ok {
		op.name = newName
	} else {
		panic("expected *Operator")
	}
}

func AssertAnyOperator(t Token) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		return op, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected operator")
	}
}

func IsAnyOperator(t Token) bool {
	_, ok := t.(*Operator)
	return ok
}

func IsSingularOperator(t Token, name string) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 0 && op.name == name
	}
	return false
}

func IsAnySingularOperator(t Token) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 0
	}
	return false
}

func AssertSingularOperator(t Token, name string) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 0 && op.name == name {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected singular operator")
}

func AssertAnySingularOperator(t Token) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 0 {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected singular operator")
}

func IsUnaryOperator(t Token, name string) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 1 && op.name == name
	}
	return false
}

func IsAnyUnaryOperator(t Token) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 1
	}
	return false
}

func AssertUnaryOperator(t Token, name string) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 1 && op.name == name {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected unary operator")
}

func AssertAnyUnaryOperator(t Token) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 1 {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected unary operator")
}

func IsBinaryOperator(t Token, name string) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 2 && op.name == name
	}
	return false
}

func IsAnyBinaryOperator(t Token) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 2
	}
	return false
}

func AssertBinaryOperator(t Token, name string) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 2 && op.name == name {
			return op, nil
		}
	}
	errCtx := t.Context()
	err := errCtx.NewError("Error: expected binary operator (" + name + ")")
	return nil, err
}

func AssertAnyBinaryOperator(t Token) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 2 {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected binary operator")
}

func IsTernaryOperator(t Token, name string) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 3 && op.name == name
	}
	return false
}

func IsAnyTernaryOperator(t Token) bool {
	if op, ok := t.(*Operator); ok {
		return len(op.args) == 3
	}
	return false
}

func AssertTernaryOperator(t Token, name string) (*Operator, error) {
	if op, ok := t.(*Operator); ok {

		if len(op.args) == 3 && op.name == name {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected ternary operator")
}

func AssertAnyTernaryOperator(t Token) (*Operator, error) {
	if op, ok := t.(*Operator); ok {
		if len(op.args) == 3 {
			return op, nil
		}
	}
	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected ternary operator")
}
