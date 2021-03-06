package functions

import (
	"math"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/svg"
)

func SVGPathPos(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	var pathStr string = ""
	switch a := args[0].(type) {
	case *tokens.List:
		pathToken, err := Str(scope, tokens.NewParens([]tokens.Token{a, tokens.NewValueString(" ", ctx)}, nil, ctx), ctx)
		if err != nil {
			return nil, err
		}

		pathStr_, err := tokens.AssertString(pathToken)
		if err != nil {
			panic(err)
		}

		pathStr = pathStr_.Value()
	case *tokens.String:
		pathStr = a.Value()
	default:
		return nil, ctx.NewError("Error: expected List or String for first argument")
	}

	b, err := tokens.AssertIntOrFloat(args[1])
	if err != nil {
		return nil, err
	}

	// now parse into a path
	pcs, err := svg.ParsePathString(pathStr, ctx)
	if err != nil {
		return nil, err
	}

	segments, err := svg.GenerateSegments(pcs, ctx)
	if err != nil {
		return nil, err
	}

	l := svg.SegmentsLength(segments)
	tVal := math.Min(l, math.Max(0.0, b.Value()*l))

	x, y := svg.SegmentsPosition(segments, tVal)

	return tokens.NewValuesList([]*tokens.Float{tokens.NewValueFloat(x, ctx), tokens.NewValueFloat(y, ctx)}, ctx), nil
}
