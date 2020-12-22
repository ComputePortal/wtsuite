package scripts

import (
	"../../tokens/js"
)

type InitFileScript struct {
	FileScriptData
}

func NewInitFileScript(relPath string, caller string) (*InitFileScript, error) {
	fileScriptData, err := newFileScriptData(relPath, caller)
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
