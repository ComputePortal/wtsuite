package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Scope interface {
	BuildMathText(x float64, y float64, fontSize float64, value string, ctx context.Context) error
	BuildMathPath(value string, ctx context.Context) error

	NewSubScope() SubScope
}

// new subscope starts at xOffset = 0, yOffset = 0, xScale = 1, yScale = 1
type SubScope interface {
	Scope

	// first offset, then scale
	Transform(xOffset float64, yOffset float64, xScale float64, yScale float64)
}
