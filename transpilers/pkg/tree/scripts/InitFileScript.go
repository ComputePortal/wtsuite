package scripts

import (
	"../../tokens/js"
	"../../tokens/js/values"
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

func (s *InitFileScript) EvalTypes(globals values.Stack) error {
	return s.module.EvalAsEntryPoint(globals)
}

func (s *InitFileScript) UniqueEntryPointNames(ns js.Namespace) error {
	return s.module.UniqueEntryPointNames(ns)
}
