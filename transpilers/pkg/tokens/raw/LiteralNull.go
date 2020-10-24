package raw

import (
	"../context"
)

type LiteralNull struct {
	TokenData
}

func NewLiteralNull(ctx context.Context) *LiteralNull {
	return &LiteralNull{TokenData{ctx}}
}

func (t *LiteralNull) Dump(indent string) string {
	return indent + "LiteralNull\n"
}

func IsLiteralNull(t Token) bool {
	switch t.(type) {
	case *LiteralNull:
		return true
	default:
		return false
	}
}
