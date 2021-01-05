package scripts

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type ControlFileScript struct {
	FileScriptData
}

func NewControlFileScript(absPath string) (*ControlFileScript, error) {
	fileScriptData, err := newFileScriptData(absPath)
	if err != nil {
		return nil, err
	}

	return &ControlFileScript{fileScriptData}, nil
}

func (s *ControlFileScript) Write() (string, error) {
	var b strings.Builder

	// wrap module in a function
	hash := s.Hash()

	b.WriteString(patterns.NL)
	b.WriteString("function ")
	b.WriteString(hash)
	//b.WriteString("(__documentVariables__, __elementConstructors__){")
	b.WriteString("(){")
	b.WriteString(patterns.NL)

	// tabs are useless, because module isnt tabbed
	/*b.WriteString("document.getVariable = (s => __documentVariables__[s]);")
	b.WriteString(js.NL)
	b.WriteString("document.newElement = (s => __elementConstructors__[s]());")

	b.WriteString(js.NL)*/

	str, err := s.module.Write(nil, patterns.NL, patterns.TAB)
	if err != nil {
		return "", err
	}

	b.WriteString(str)
	b.WriteString("}")

	return b.String(), nil
}
