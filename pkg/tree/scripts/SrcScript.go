package scripts

import (
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SrcScript struct {
	src string
}

func NewSrcScript(src string) (*SrcScript, error) {
	return &SrcScript{src}, nil
}

func (s *SrcScript) Write() (string, error) {
	return "", nil
}

func (s *SrcScript) Dependencies() []files.PathLang {
	return []files.PathLang{files.PathLang{s.src, files.SCRIPT, context.NewDummyContext()}}
}
