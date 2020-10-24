package scripts

import (
	"strings"

	"../../tokens/js"
)

type ControlFileScript struct {
	views []string
	FileScriptData
}

func NewControlFileScript(relPath string, caller string, views []string) (*ControlFileScript, error) {
	fileScriptData, err := newFileScriptData(relPath, caller)
	if err != nil {
		return nil, err
	}

	return &ControlFileScript{views, fileScriptData}, nil
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

func (s *ControlFileScript) hasView(v string) bool {
	for _, view := range s.views {
		if view == v {
			return true
		}
	}

	return false
}
