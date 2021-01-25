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

	for _, st := range t.statements {
		s := st.WriteStatement(usage, indent, nl, tab)

    if strings.HasPrefix(s, "#") { // is preproc directive
      b.WriteString("\n")
      b.WriteString(s)
      b.WriteString("\n")
    } else if s != "" {
			b.WriteString(s)
      b.WriteString(";")
      b.WriteString(nl)
		}
	}

	return b.String()
}

func (t *Block) ResolveStatementNames(scope Scope) error {
	for _, st := range t.statements {
		if err := st.ResolveStatementNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Block) evalStatements() error {
	for _, st := range t.statements {
		err := st.EvalStatement()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Block) ResolveStatementActivity(usage Usage) error {
	for i := len(t.statements) - 1; i >= 0; i-- {
		st := t.statements[i]

		if err := st.ResolveStatementActivity(usage); err != nil {
			return err
		}

		if i < len(t.statements)-1 {
			errCtx := context.MergeContexts(st.Context(), t.statements[i+1].Context())

			switch st.(type) {
			case *Return:
				return errCtx.NewError("Error: unreachable statement after return statement")
			}
		}
	}

	return nil
}

func (t *Block) UniqueStatementNames(ns Namespace) error {
  for _, st := range t.statements {
    if err := st.UniqueStatementNames(ns); err != nil {
      return err
    }
  }

  return nil
}
