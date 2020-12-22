package scripts

import (
	"strings"

	"../../tokens/js"
)

type ControlFileScript struct {
	FileScriptData
}

func NewControlFileScript(relPath string, caller string) (*ControlFileScript, error) {
	fileScriptData, err := newFileScriptData(relPath, caller)
	if err != nil {
		return nil, err
	}

	return &ControlFileScript{fileScriptData}, nil
}

func (s *ControlFileScript) Write() (string, error) {
	var b strings.Builder

	// wrap module in a function
	hash := s.Hash()

	b.WriteString(js.NL)
	b.WriteString("function ")
	b.WriteString(hash)
	//b.WriteString("(__documentVariables__, __elementConstructors__){")
	b.WriteString("(){")
	b.WriteString(js.NL)

	// tabs are useless, because module isnt tabbed
	/*b.WriteString("document.getVariable = (s => __documentVariables__[s]);")
	b.WriteString(js.NL)
	b.WriteString("document.newElement = (s => __elementConstructors__[s]());")

	b.WriteString(js.NL)*/

	str, err := s.module.Write()
	if err != nil {
		return "", err
	}

	b.WriteString(str)
	b.WriteString("}")

	return b.String(), nil
}
