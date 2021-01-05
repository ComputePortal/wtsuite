package scripts

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
)

type InitFileScript struct {
	FileScriptData
}

func NewInitFileScript(absPath string) (*InitFileScript, error) {
	fileScriptData, err := newFileScriptData(absPath)
	if err != nil {
		return nil, err
	}

	return &InitFileScript{fileScriptData}, nil
}

func (s *InitFileScript) EvalTypes() error {
	return s.module.EvalTypes()
}

func (s *InitFileScript) UniqueEntryPointNames(ns js.Namespace) error {
	return s.module.UniqueEntryPointNames(ns)
}
