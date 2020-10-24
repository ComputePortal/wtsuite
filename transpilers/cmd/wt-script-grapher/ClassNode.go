package main

import (
  "strings"

	"../../pkg/tokens/js"
)

type ClassNode struct {
  class *js.Class
}

func NewClassNode(class *js.Class) Node {
  return &ClassNode{
    class,
  }
}

func (n *ClassNode) Name() string {
  return n.class.Name()
}

func (n *ClassNode) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(n.class.Name())
  b.WriteString("\" [shape=\"box\"];")

  return b.String()
}
