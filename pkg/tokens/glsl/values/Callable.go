package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Callable interface {
  EvalCall(args []Value, ctx context.Context) (Value, error)
}
