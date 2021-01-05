package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SetDynamicIndexFunction struct {
  vecType string // empty if just regular array
  typeName string
  DynamicIndexFunction
}

func NewSetDynamicIndexFunction(containerTypeName string, typeName string, length int, ctx context.Context) *SetDynamicIndexFunction {
  name := "setIndex_" + typeName + "_" + strconv.Itoa(length)

  vecType := ""

  if strings.HasSuffix(containerTypeName, "vec" + strconv.Itoa(length)) {
    vecType = containerTypeName
    name = "setIndex_" + containerTypeName + "_" + typeName
  }

  return &SetDynamicIndexFunction{vecType, typeName, newDynamicIndexFunction(name, length, ctx)}
}

func (t *SetDynamicIndexFunction) Dump(indent string) string {
  return indent + t.Name()
}

func (t *SetDynamicIndexFunction) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("void ")
  b.WriteString(t.Name())

  b.WriteString("(inout ")
  if t.vecType != "" {
    b.WriteString(t.vecType)
    b.WriteString(" ")
    b.WriteString("x")
  } else {
    b.WriteString(t.typeName)
    b.WriteString(" ")
    b.WriteString("x[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  b.WriteString(",in int i,in ")
  b.WriteString(t.typeName)
  b.WriteString(" a){")
  b.WriteString(nl)
  b.WriteString(t.DynamicIndexFunction.writeTree("i", indent + tab, nl, tab, func(i int) string {
    return "x[" + strconv.Itoa(i) + "]=a;"
  }))
  b.WriteString(indent)
  b.WriteString("}")

  return b.String()
}
