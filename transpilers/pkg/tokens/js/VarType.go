package js

import (
	"../context"
)

type VarType int

const (
	CONST VarType = iota
	LET
	VAR
)

func StringToVarType(s string, ctx context.Context) (VarType, error) {
	switch s {
	case "const":
		return CONST, nil
	case "let":
		return LET, nil
	case "var":
		return VAR, nil
	default:
		return CONST, ctx.NewError("Error: expected 'var', 'let' or 'const', got " + s)
	}
}

func VarTypeToString(varType VarType) string {
	switch varType {
	case CONST:
		return "const"
	case LET:
		return "let"
	case VAR:
		return "var"
	default:
		panic("unhandled")
	}
}
