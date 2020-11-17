package js

import (
	"strings"

	"../context"
)

type Block struct {
	statements []Statement
	TokenData
}

func newBlock(ctx context.Context) Block {
	return Block{make([]Statement, 0), TokenData{ctx}}
}

func (t *Block) AddStatement(statement Statement) {
	t.statements = append(t.statements, statement)
}

func (t *Block) NewScope(parent Scope) Scope {
	return NewBlockScope(parent)
}

func (t *Block) Dump(indent string) string {
	var b strings.Builder

	for _, statement := range t.statements {
		b.WriteString(statement.Dump(indent))
	}

	return b.String()
}

func (t *Block) writeBlockStatements(indent string, nl string) string {
	var b strings.Builder

	prevWroteSomething := false
	for _, st := range t.statements {
		s := st.WriteStatement(indent)

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

///////////////////////////
// 1. Name resolution stage
///////////////////////////
func (t *Block) HoistNames(scope Scope) error {
	for _, st := range t.statements {
		if err := st.HoistNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Block) ResolveStatementNames(scope Scope) error {
	for _, st := range t.statements {
		if err := st.ResolveStatementNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *Block) HoistAndResolveStatementNames(scope Scope) error {
	if err := t.HoistNames(scope); err != nil {
		return err
	}

	return t.ResolveStatementNames(scope)
}

func (t *Block) evalStatements(statements []Statement) error {
	for _, st := range statements {
		err := st.EvalStatement()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Block) EvalStatement() error {
	return t.evalStatements(t.statements)
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
				// TODO: break and continue
			}
		}
	}

	return nil
}

func (t *Block) UniversalStatementNames(ns Namespace) error {
	for _, st := range t.statements {
		if err := st.UniversalStatementNames(ns); err != nil {
			return err
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

func (t *Block) Walk(fn WalkFunc) error {
	for _, st := range t.statements {
		if err := st.Walk(fn); err != nil {
			return err
		}
	}

	return nil
}
