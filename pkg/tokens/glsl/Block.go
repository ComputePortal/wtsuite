package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Block struct {
  statements []Statement
  TokenData
}

func newBlock(ctx context.Context) Block {
  return Block{make([]Statement, 0), TokenData{ctx}}
}

func NewBlock(ctx context.Context) *Block {
  bl := newBlock(ctx)

  return &bl
}

func (t *Block) AddStatement(statement Statement) {
  t.statements = append(t.statements, statement)
}

func (t *Block) Dump(indent string) string {
  var b strings.Builder

  for _, statement := range t.statements {
    b.WriteString(statement.Dump(indent))
  }

  return b.String()
}

func (t *Block) writeBlockStatements(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

	prevWroteSomething := false
	for _, st := range t.statements {
		s := st.WriteStatement(usage, indent, nl, tab)

		if s != "" {
			if prevWroteSomething {
				b.WriteString(";")
				b.WriteString(nl)
			}

			b.WriteString(s)

			prevWroteSomething = true
		}
	}

	return b.String()
}
