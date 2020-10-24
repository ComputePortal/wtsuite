package js

import (
	"reflect"

	"../context"
)

var (
	NL             = "\n"
	TAB            = "  "
	COMPACT_NAMING = false
	VERBOSITY      = 0
)

type Token interface {
	Dump(indent string) string
	Context() context.Context
}

type TokenData struct {
	ctx context.Context
}

func newTokenData(ctx context.Context) TokenData {
	return TokenData{ctx}
}

func (t *TokenData) Context() context.Context {
	return t.ctx
}

func MergeContexts(ts ...Token) context.Context {
	ctxs := make([]context.Context, len(ts))

	for i, t := range ts {
		ctxs[i] = t.Context()
	}

	return context.MergeContexts(ctxs...)
}

// used by the parser
func IsCallable(t Token) bool {
	switch t.(type) {
	case *Function, *VarExpression, *Call, *Index, *Member, *Parens:
		return true
	case *LiteralBoolean, *LiteralInt, *LiteralFloat, *LiteralString, Op, *Class:
		return false
	default:
		panic("unhandled")
	}
}

func IsIndexable(t Token) bool {
	switch t.(type) {
	case *VarExpression, *Call, *LiteralString, *Index, *Member, *Parens, *LiteralArray, *LiteralObject:
		return true
	case *LiteralBoolean, *LiteralInt, *LiteralFloat, Op, *Function, *Class:
		return false
	default:
		panic("unhandled" + reflect.TypeOf(t).String())
	}
}

func AssertCallable(t Token) error {
	if !IsCallable(t) {
		errCtx := t.Context()
		return errCtx.NewError("Error: not callable")
	}

	return nil
}

func AssertIndexable(t Token) error {
	if !IsIndexable(t) {
		errCtx := t.Context()
		return errCtx.NewError("Error: not indexable")
	}

	return nil
}
