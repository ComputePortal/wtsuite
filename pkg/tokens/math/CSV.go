package math

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"
)

const (
	commaExtraPreSpacing  = 0.00
	commaExtraPostSpacing = 0.10
)

type CSV struct {
	args   []Token
	commas []Token
	TokenData
}

func NewCSV(args []Token, ctx context.Context) (Token, error) {
	return &CSV{args, make([]Token, 0), newTokenData(ctx)}, nil
}

func (t *CSV) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("CSV:\n")

	for _, arg := range t.args {
		b.WriteString(arg.Dump(indent + "  "))
	}

	return b.String()
}

func (t *CSV) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	var bbTotal boundingbox.BB = nil

	for i, arg := range t.args {
		bbArg, err := arg.GenerateTags(scope, x, y)
		if err != nil {
			return nil, err
		}

		if bbTotal == nil {
			bbTotal = bbArg
		} else {
			bbTotal = boundingbox.Merge(bbTotal, bbArg)
		}

		x = bbArg.Right()

		if i < len(t.args)-1 {
			x += commaExtraPreSpacing

			bbComma, err := t.commas[i].GenerateTags(scope, x, y)
			if err != nil {
				return nil, err
			}

			bbTotal = boundingbox.Merge(bbComma, bbTotal)

			x = bbComma.Right() + commaExtraPostSpacing
		}
	}

	return bbTotal, nil
}
