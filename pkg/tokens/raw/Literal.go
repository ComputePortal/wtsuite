package raw

func IsLiteral(t Token) bool {
	switch t.(type) {
	case *LiteralBool, *LiteralColor, *LiteralFloat, *LiteralInt, *LiteralNull, *LiteralString:
		return true
	}

	return false
}
