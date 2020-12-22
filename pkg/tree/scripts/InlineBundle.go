package scripts

import (
	"strings"
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

func (b *InlineBundle) Dependencies() []string {
	// src's
	uniqueDeps := make(map[string]string) // to make them unique

	for _, s := range b.scripts {
		deps := s.Dependencies()

		for _, dep := range deps {
			if _, ok := uniqueDeps[dep]; !ok {
				uniqueDeps[dep] = dep
			}
		}
	}

	result := make([]string, 0)

	for k, _ := range uniqueDeps {
		result = append(result, k)
	}

	return result
}
