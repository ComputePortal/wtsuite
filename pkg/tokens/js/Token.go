package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

var (
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

func AssertCallable(t Token) error {
	if !IsCallable(t) {
		errCtx := t.Context()
		return errCtx.NewError("Error: not callable")
	}

	return nil
}
