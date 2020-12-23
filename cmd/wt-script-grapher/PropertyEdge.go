package main

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type PropertyEdge struct {
  // TODO: should edge be given a label?
  owner *values.Instance
  label string
  property *values.Instance // parent can potentially be builtin
}

func NewPropertyEdge(owner *values.Instance, label string, property *values.Instance) Edge {
  return &PropertyEdge{
    owner,
    label,
    property,
  }
}

func (e *PropertyEdge) Write(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("\"")
  b.WriteString(e.owner.TypeName())
  b.WriteString("\" -> \"")
  b.WriteString(e.property.TypeName())
  b.WriteString("\" [color=\"blue\"")
  if e.label != "" {
    b.WriteString(",label=\"")
    b.WriteString(e.label)
    b.WriteString("\"")
  }
  b.WriteString("];")

  return b.String()
}
