package main

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type InheritsEdge struct {
  child *js.Class
  parent values.Prototype // parent can potentially be builtin
}

func NewInheritsEdge(child *js.Class, parent values.Prototype) Edge {
  return &InheritsEdge{
    child,
    parent,
  }
}

func (e *InheritsEdge) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(e.child.Name())
  b.WriteString("\" -> \"")
  b.WriteString(e.parent.Name())
  b.WriteString("\" [shape=\"onormal\"];")

  return b.String()
}
