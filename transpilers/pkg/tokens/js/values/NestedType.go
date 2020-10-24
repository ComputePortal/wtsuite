package values

import (
	"strings"

	"../../context"
)

// used for nested type cast
type NestedType struct {
	key      string
	interf   Interface
	children []*NestedType // can be nil

	ctx context.Context
}

func NewNestedType(key string, interf Interface, children []*NestedType, ctx context.Context) *NestedType {
	return &NestedType{key, interf, children, ctx}
}

func (t *NestedType) Key() string {
	return t.key
}

func (t *NestedType) SetKey(k string) {
	t.key = k
}

func (t *NestedType) Name() string {
	return t.interf.Name()
}

func (t *NestedType) Interface() Interface {
	return t.interf
}

func (t *NestedType) Children() []*NestedType {
	return t.children
}

func (t *NestedType) Context() context.Context {
	return t.ctx
}

// for debugging
func (t *NestedType) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	if t.key != "" {
		b.WriteString(t.key)
		b.WriteString(":")
	}
	b.WriteString(t.interf.Name())
	b.WriteString("\n")

	if t.children != nil {
		for _, child := range t.children {
			b.WriteString(child.Dump(indent + "  "))
		}
	}

	return b.String()
}
