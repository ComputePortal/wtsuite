package scripts

import (
	"github.com/computeportal/wtsuite/pkg/files"
)

var (
	VERBOSITY = 0
)

type Script interface {
	Write() (string, error)
	Dependencies() []files.PathLang // src fields in script or call
}
