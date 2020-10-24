package tree

import (
	tokens "../tokens/html"
)

type ClassMap interface {
	HasTag(class string) bool
	GetTag(class string) []VisibleTag
	AppendTag(class string, t VisibleTag)

	HasScriptClass(class string) bool
}

type ClassMapData struct {
	tags          map[string][]VisibleTag
	scriptClasses map[string]*tokens.String
}

func NewClassMap() ClassMap {
	return &ClassMapData{make(map[string][]VisibleTag), make(map[string]*tokens.String)}
}

func (m *ClassMapData) HasTag(class string) bool {
	_, ok := m.tags[class]
	return ok
}

func (m *ClassMapData) HasScriptClass(class string) bool {
	_, ok := m.scriptClasses[class]
	return ok
}

func (m *ClassMapData) GetTag(class string) []VisibleTag {
	ts, ok := m.tags[class]
	if !ok {
		panic("should've been caught before")
	}

	return ts
}

func (m *ClassMapData) AppendTag(class string, t VisibleTag) {
	if !m.HasTag(class) {
		m.tags[class] = make([]VisibleTag, 0)
	}

	m.tags[class] = append(m.tags[class], t)
}
