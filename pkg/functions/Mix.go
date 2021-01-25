package functions

import (
	"math"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Mix(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, NewInterface([]string{"a", "b", "f"}, ctx))
  if err != nil {
    return nil, err
  }

	color1, err := tokens.AssertColor(args[0])
	if err != nil {
		return nil, err
	}

	color2, err := tokens.AssertColor(args[1])
	if err != nil {
		return nil, err
	}

	factor, err := tokens.AssertFractionFloat(args[2])
	if err != nil {
		return nil, err
	}

	r1, g1, b1, a1 := color1.Values()
	r2, g2, b2, a2 := color2.Values()
	f := factor.Value()

	mix := func(a int, b int) int {
		m := float64(a)*(1.0-f) + float64(b)*f

		mInt := int(math.Round(m))

		if mInt < 0 {
			return 0
		} else if mInt > 255 {
			return 255
		} else {
			return mInt
		}
	}

	rm := mix(r1, r2)
	gm := mix(g1, g2)
	bm := mix(b1, b2)
	am := mix(a1, a2)

	return tokens.NewValueColor(rm, gm, bm, am, ctx), nil
}
