package functions

import (
	"errors"
	"fmt"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Var struct {
	Value    tokens.Token
	Constant bool
	Auto     bool
	Imported bool
	Exported bool
	Ctx      context.Context
}

func (v Var) ToJSTypeValue() (string, string, error) {
	switch t := v.Value.(type) {
	case *tokens.Int:
		return "Int", fmt.Sprintf("%d", t.Value()), nil
	case *tokens.Float:
		if t.Unit() != "" {
			return "", "", errors.New("united float not supported")
		}

		return "Number", fmt.Sprintf("%f", t.Value()), nil
	case *tokens.String:
		return "String", "\"" + t.Value() + "\"", nil
	default:
		return "", "", errors.New("not (yet) supported")
	}
}
