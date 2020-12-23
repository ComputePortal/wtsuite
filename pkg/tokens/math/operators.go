package math

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type UnaryOp struct {
	name string
	a    Token
	TokenData
}

type BinaryOp struct {
	name string
	a    Token
	b    Token
	TokenData
}

func NewBinaryOp(name string, a Token, b Token, ctx context.Context) (Token, error) {
	switch name {
	case "-":
		return NewSubOp(a, b, ctx)
	case "+":
		return NewAddOp(a, b, ctx)
	case "*":
		return NewMulOp(a, b, ctx)
	case "^":
		return NewPowOp(a, b, ctx)
	case "_":
		return NewIndexOp(a, b, ctx)
	case "=":
		return NewEqualsOp(a, b, ctx)
	case "!=":
		return NewNEOp(a, b, ctx)
	case "~=":
		return NewApproxOp(a, b, ctx)
	case "/":
		return NewDivOp(a, b, ctx)
	case ">=":
		return NewGEOp(a, b, ctx)
	case ">":
		return NewGTOp(a, b, ctx)
	case ">>":
		return NewBinSymbolOp(genericBinSymbolSpacing, genericBinSymbolSpacing, newUnicodeSymbol(">>", 0x226b, ctx), a, b, ctx)
	case "<":
		return NewLTOp(a, b, ctx)
	case "<<":
		return NewBinSymbolOp(genericBinSymbolSpacing, genericBinSymbolSpacing, newUnicodeSymbol("<<", 0x226a, ctx), a, b, ctx)
	case "<=":
		return NewLEOp(a, b, ctx)
	case "->":
		return NewRightArrow1(a, b, ctx)
	case "=>":
		return NewRightArrow2(a, b, ctx)
	default:
		panic("binary op " + name + " not yet handled")
	}
}

func (t *UnaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.name)
	b.WriteString("()\n")

	b.WriteString(t.a.Dump(indent + "  "))

	return b.String()
}

func (t *BinaryOp) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(t.name)
	b.WriteString("()\n")

	b.WriteString(t.a.Dump(indent + "  "))
	b.WriteString(t.b.Dump(indent + "  "))

	return b.String()
}

func NewPreUnaryOp(name string, a Token, ctx context.Context) (Token, error) {
	switch name {
	case "-":
		return NewNegOp(a, ctx)
	default:
		panic("preunary op " + name + " not yet handled")
	}
}

func NewPostUnaryOp(name string, a Token, ctx context.Context) (Token, error) {
	panic("postunary op " + name + " not yet handled")
}
