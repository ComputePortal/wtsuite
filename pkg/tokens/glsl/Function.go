package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Function struct {
  fi *FunctionInterface
  ret []*Return // registered via scope
  Block
}

func NewFunction(fi *FunctionInterface, statements []Statement, ctx context.Context) *Function {
  fn := &Function{
    fi,
    make([]*Return, 0),
    newBlock(ctx),
  }

  for _, st := range statements {
    fn.AddStatement(st)
  }

  return fn
}

func (t *Function) Name() string {
  return t.fi.Name()
}

func (t *Function) Dump(indent string) string {
  var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Function")

	if t.Name() != "" {
    // name itself is dumped by function interface
		b.WriteString(" ")
	}

	b.WriteString(t.fi.Dump(indent))

  // override t.Block.Dump()
	for _, st := range t.statements {
		b.WriteString(st.Dump(indent + "{ "))
	}

	return b.String()
}

func (t *Function) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(t.fi.WriteInterface())

  b.WriteString("{")
  b.WriteString(nl)
  b.WriteString(t.writeBlockStatements(usage, indent+tab, nl, tab))
  b.WriteString(nl)
  b.WriteString(indent)
  b.WriteString("}")
  b.WriteString(nl)

  return b.String()
}
