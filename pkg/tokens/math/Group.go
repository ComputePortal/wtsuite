package math

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	GROUP_SEP_LEFT_SPACING  = 0.05
	GROUP_SEP_RIGHT_SPACING = 0.1
)

type Group struct {
	content []Token
	sep     *Symbol
	TokenData
}

func NewGroup(content []Token, sepStr string, ctx context.Context) (Token, error) {
	sep := newSymbol(sepStr, ctx)

	return &Group{content, sep, newTokenData(ctx)}, nil
}

func (t *Group) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Group()\n")

	for _, part := range t.content {
		b.WriteString(part.Dump(indent + "  "))
	}

	return b.String()
}

func (t *Group) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	// dont need subscope, simple forward integration

	left := x
	var bbTot boundingbox.BB = nil
	for i, part := range t.content {
		bbPart, err := part.GenerateTags(scope, left, y)
		if err != nil {
			return nil, err
		}

		if bbTot == nil {
			bbTot = bbPart
		} else {
			bbTot = boundingbox.Merge(bbTot, bbPart)
		}

		left = bbPart.Right()

		if i < len(t.content)-1 {
			bbSep, err := t.sep.GenerateTags(scope, left+GROUP_SEP_LEFT_SPACING, y)
			if err != nil {
				return nil, err
			}

			bbTot = boundingbox.Merge(bbTot, bbSep)
			left = bbSep.Right() + GROUP_SEP_RIGHT_SPACING
		}
	}

	if bbTot == nil {
		panic("empty group not yet handled")
	}

	return bbTot, nil
}
