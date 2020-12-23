package main

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type InterfaceNode struct {
  interf values.Interface
}

func NewInterfaceNode(interf values.Interface) Node {
  return &InterfaceNode{
    interf,
  }
}

func (n *InterfaceNode) Name() string {
  return n.interf.Name()
}

func (n *InterfaceNode) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(n.interf.Name())
  b.WriteString("\" [shape=\"box\",color=\"red\"];")

  return b.String()
}
