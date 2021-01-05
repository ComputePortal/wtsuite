package scripts

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
)

type InlineBundle struct {
	// in the order they appear in the htmlpp file
	scripts []Script
}

func NewInlineBundle() *InlineBundle {
	return &InlineBundle{make([]Script, 0)}
}

func (b *InlineBundle) Append(s Script) {
	b.scripts = append(b.scripts, s)
}

func (b *InlineBundle) IsEmpty() bool {
	return len(b.scripts) == 0
}

func (b *InlineBundle) Write() (string, error) {
	var sb strings.Builder

	for _, s := range b.scripts {
		str, err := s.Write()
		if err != nil {
			return sb.String(), err
		}

		sb.WriteString(str)
	}

	return sb.String(), nil
}

func (b *InlineBundle) Dependencies() []files.PathLang {
	// src's
	uniqueDeps := make(map[string]files.PathLang) // to make them unique

	for _, s := range b.scripts {
		deps := s.Dependencies()

		for _, pl := range deps {
      dep := pl.Path

			if _, ok := uniqueDeps[dep]; !ok {
				uniqueDeps[dep] = pl
			}
		}
	}

	result := make([]files.PathLang, 0)

	for _, pl := range uniqueDeps {
		result = append(result, pl)
	}

	return result
}
