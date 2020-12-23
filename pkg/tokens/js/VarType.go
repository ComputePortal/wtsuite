package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type VarType int

const (
	CONST VarType = iota
	LET
	VAR
  AUTOLET
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
	case LET, AUTOLET:
		return "let"
	case VAR:
		return "var"
	default:
		panic("unhandled")
	}
}
