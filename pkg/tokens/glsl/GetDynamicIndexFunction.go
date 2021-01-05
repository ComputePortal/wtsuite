package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type GetDynamicIndexFunction struct {
  vecType string // empty if just regular array
  typeName string
  DynamicIndexFunction
}

func NewGetDynamicIndexFunction(containerTypeName string, typeName string, length int, ctx context.Context) *GetDynamicIndexFunction {
  name := "getIndex_" + typeName + "_" + strconv.Itoa(length)

  vecType := ""

  if strings.HasSuffix(containerTypeName, "vec" + strconv.Itoa(length)) {
    vecType = containerTypeName
    name = "getIndex_" + containerTypeName + "_" + typeName
  }

  return &GetDynamicIndexFunction{vecType, typeName, newDynamicIndexFunction(name, length, ctx)}
}

func (t *GetDynamicIndexFunction) Dump(indent string) string {
  return indent + t.Name()
}

func (t *GetDynamicIndexFunction) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(t.typeName)
  b.WriteString(" ")
  b.WriteString(t.Name())

  b.WriteString("(in ")
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

  b.WriteString(",in int i){")
  b.WriteString(nl)
  b.WriteString(t.DynamicIndexFunction.writeTree("i", indent + tab, nl, tab, func(i int) string {
    return "return x[" + strconv.Itoa(i) + "];"
  }))
  b.WriteString(indent)
  b.WriteString("}")

  return b.String()
}
