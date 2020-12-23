package main

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type InstanceNode struct {
  instance *values.Instance
}

func NewInstanceNode(instance *values.Instance) Node {
  return &InstanceNode{
    instance,
  }
}

func (n *InstanceNode) Name() string {
  return n.instance.TypeName()
}

func (n *InstanceNode) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(n.instance.TypeName()) // TODO: lower case first letters that are capitalized (except last?)
  b.WriteString("\" [shape=\"box\"];")

  return b.String()
}
