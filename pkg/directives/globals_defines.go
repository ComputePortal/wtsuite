package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

var _defines map[string]*tokens.String = nil

func RegisterDefine(k, v string) {
	if _defines == nil {
		_defines = make(map[string]*tokens.String)
	}

	_defines[k] = tokens.NewValueString(v, context.NewDummyContext())
}

func HasDefine(k string) bool {
	_, ok := _defines[k]
	return ok
}

func GetDefine(k string) *tokens.String {
	return _defines[k]
}
