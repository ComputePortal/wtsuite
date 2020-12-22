package main

import (
  "strings"

	"../../pkg/tokens/js"
	"../../pkg/tokens/js/values"
)

type ImplementsEdge struct {
  class *js.Class
  interf values.Interface // parent can potentially be builtin
}

func NewImplementsEdge(class *js.Class, interf values.Interface) Edge {
  return &ImplementsEdge{
    class,
    interf,
  }
}

func (e *ImplementsEdge) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(e.class.Name())
  b.WriteString("\" -> \"")
  b.WriteString(e.interf.Name())
  b.WriteString("\" [shape=\"onormal\",style=\"dotted\"];")

  return b.String()
}
