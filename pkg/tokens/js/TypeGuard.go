package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type TypeGuard interface {
	// collect all variables/interfaces
	// return false if all typeguards should be voided
	// (should also do everything EvalExpression does)
	CollectTypeGuards(c map[Variable]values.Interface) (bool, error)
}
