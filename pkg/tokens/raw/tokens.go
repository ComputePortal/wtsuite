package raw

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Token interface {
	Dump(indent string) string // to inspect the syntax-tree
	Context() context.Context
}

type TokenData struct {
	ctx context.Context
}

func (t *TokenData) Context() context.Context {
	return t.ctx
}

func IsOperable(t Token) bool {
	return !IsAnyNonWordSymbol(t)
}

func AssertOperable(t Token) error {
	if !IsOperable(t) {
		errCtx := t.Context()
		return errCtx.NewError("Error: not operable")
	}

	return nil
}

func Concat(args ...interface{}) []Token {
	result := make([]Token, 0)

	for _, arg := range args {
		switch t := arg.(type) {
		case []Token:
			result = append(result, t...)
		case Token:
			result = append(result, t)
		default:
			panic("not Token or []Token")
		}
	}

	return result
}

func MergeContexts(ts ...Token) context.Context {
	ctxs := make([]context.Context, len(ts))

	for i, t := range ts {
		ctxs[i] = t.Context()
	}

	return context.MergeContexts(ctxs...)
}

func ShortHash(content string) string {
	hasher := sha1.New()
	hasher.Write([]byte(content))
	fullHash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	hash := ""

	for _, c := range fullHash {
		if len(hash) > 5 {
			break
		}
		// not a digit
		if (c > 0x40 && c < 0x5b) || (c > 0x60 && c < 0x7b) {
			hash += string(c)
		}
	}

	if hash == "" {
		panic("hash failed")
	}

	return hash
}
