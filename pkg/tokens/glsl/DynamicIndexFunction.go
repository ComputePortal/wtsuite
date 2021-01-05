package glsl

import (
  "math"
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type DynamicIndexFunction struct {
  length int
  MacroFunction
}

func newDynamicIndexFunction(name string, length int, ctx context.Context) DynamicIndexFunction {
  return DynamicIndexFunction{length, newMacroFunction(name, ctx)}
}

func (t *DynamicIndexFunction) writeTreeRecursively(i string, indent string, a int, b int, nl string, tab string, fnLeaf func(int) string) string {
  if (a == b - 1) {
    return indent + fnLeaf(a)  + nl
  } else {
    h := int(math.Ceil(0.5*(float64(a) + float64(b))))

    var s strings.Builder

    s.WriteString(indent)
    s.WriteString("if(")
    s.WriteString(i)
    s.WriteString("<")
    s.WriteString(strconv.Itoa(h))
    s.WriteString("){")
    s.WriteString(nl)
    s.WriteString(t.writeTreeRecursively(i, indent + tab, a, h, nl, tab, fnLeaf))
    s.WriteString(indent)
    s.WriteString("}else{")
    s.WriteString(nl)
    s.WriteString(t.writeTreeRecursively(i, indent + tab, h, b, nl, tab, fnLeaf))
    s.WriteString(indent)
    s.WriteString("}")
    s.WriteString(nl)

    return s.String()
  }
}

func (t *DynamicIndexFunction) writeTree(i string, indent string, nl string, tab string, fnLeaf func(int) string) string {

  return t.writeTreeRecursively(i, indent, 0, t.length, nl, tab, fnLeaf)
}
