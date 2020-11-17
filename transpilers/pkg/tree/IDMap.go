package tree

import (
	"fmt"

	"./styles"

	"../tokens/context"
	"../tokens/js"
)

var IsAutoUID func(id string) bool = nil

// only requires a subset of that tag methods
type IDMapTag interface {
	Name() string
	ToJSType() string
	InnerHTML() string // empty if no direct descendant is text
	CollectStates() map[string][]string
	Context() context.Context
}

type IDMap interface {
	Has(id string) bool
	Get(id string) IDMapTag
	Set(id string, t IDMapTag)

	Dump() // for debugging
}

type IDMapData struct {
	tags map[string]IDMapTag
}

func NewIDMap() IDMap {
	return &IDMapData{make(map[string]IDMapTag)}
}

func (m *IDMapData) Has(id string) bool {
	_, ok := m.tags[id]
	return ok
}

func (m *IDMapData) Get(id string) IDMapTag {
	t, ok := m.tags[id]
	if !ok {
		panic("should've been caught before")
	}

	return t
}

func (m *IDMapData) Set(id string, t IDMapTag) {
	m.tags[id] = t
}

func (m *IDMapData) Dump() {
	for k, v := range m.tags {
		fmt.Println(k + ": " + v.Name())
	}
}
