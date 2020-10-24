package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type TryCatch struct {
	try     []Statement // cannot be nil
	arg     *FunctionArgument
	catch   []Statement // can be nil
	finally []Statement // can be nil
	Block
}

func NewTryCatch(ctx context.Context) (*TryCatch, error) {
	return &TryCatch{make([]Statement, 0), nil, nil, nil, newBlock(ctx)}, nil
}

// arg can be nil
func (t *TryCatch) AddCatch(arg *FunctionArgument) error {
	if t.arg != nil || t.catch != nil {
		errCtx := arg.Context()
		return errCtx.NewError("Error: catch block already defined")
	}

	if arg != nil {
		t.arg = arg

		if t.arg.def != nil {
			errCtx := arg.Context()
			return errCtx.NewError("Error: catch arg cant have default")
		}
	}

	t.catch = make([]Statement, 0)
	return nil
}

func (t *TryCatch) AddFinally() error {
	if t.finally != nil {
		errCtx := t.Context()
		// TODO: get context of duplicate for nicer error message
		return errCtx.NewError("Error: finally block already defined")
	}

	t.finally = make([]Statement, 0)
	return nil
}

func (t *TryCatch) AddStatement(statement Statement) {
	if t.finally != nil {
		t.finally = append(t.finally, statement)
	} else if t.catch != nil {
		t.catch = append(t.catch, statement)
	} else {
		t.try = append(t.try, statement)
	}
}

func (t *TryCatch) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Try\n")
	t.statements = t.try
	b.WriteString(t.Block.Dump(indent + "  "))

	if t.catch != nil {
		b.WriteString(indent)
		b.WriteString("Catch")
		if t.arg != nil {
			b.WriteString("(")
			b.WriteString(strings.Replace(t.arg.Write(false), "\n", "", -1))
			b.WriteString(")")
		}
		b.WriteString("\n")
		t.statements = t.catch
		b.WriteString(t.Block.Dump(indent + "  "))
	}

	if t.finally != nil {
		b.WriteString(indent)
		b.WriteString("Finally\n")
		t.statements = t.finally
		b.WriteString(t.Block.Dump(indent + "  "))
	}

	return b.String()
}

func (t *TryCatch) WriteStatement(indent string) string {
	var b strings.Builder

	writeBlock := func(statements []Statement) {
		b.WriteString("{")
		b.WriteString(NL)
		t.statements = statements
		b.WriteString(t.writeBlockStatements(indent+TAB, NL))
		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString("}")
	}

	b.WriteString(indent)
	b.WriteString("try")
	writeBlock(t.try)

	if t.catch != nil {
		b.WriteString("catch")
		if t.arg != nil {
			b.WriteString("(")
			b.WriteString(t.arg.Write(false))
			b.WriteString(")")
		}
		writeBlock(t.catch)
	}

	if t.finally != nil {
		b.WriteString("finally")
		writeBlock(t.finally)
	}

	return b.String()
}

func (t *TryCatch) HoistNames(scope Scope) error {
	// try part is exactly like if
	t.statements = t.try
	if err := t.Block.HoistNames(scope); err != nil {
		return err
	}

	if t.catch != nil {
		t.statements = t.catch
		if err := t.Block.HoistNames(scope); err != nil {
			return err
		}
	}

	if t.finally != nil {
		t.statements = t.finally
		if err := t.Block.HoistNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *TryCatch) ResolveStatementNames(scope Scope) error {
	if t.catch == nil && t.finally == nil {
		errCtx := t.Context()
		return errCtx.NewError("Error: neither catch or finally specified")
	}

	subScope := NewBranchScope(scope)
	t.statements = t.try
	if err := t.Block.ResolveStatementNames(subScope); err != nil {
		return err
	}

	if t.catch != nil {
		subScope := NewBranchScope(scope)
		if t.arg != nil {
			if err := t.arg.ResolveNames(scope, subScope); err != nil {
				return err
			}
		}

		t.statements = t.catch
		if err := t.Block.ResolveStatementNames(subScope); err != nil {
			return err
		}
	}

	if t.finally != nil {
		subScope := NewBranchScope(scope)
		t.statements = t.finally

		if err := t.Block.ResolveStatementNames(subScope); err != nil {
			return err
		}
	}

	return nil
}

func (t *TryCatch) HoistValues(stack values.Stack) error {
	return nil
}

func (t *TryCatch) EvalStatement(stack values.Stack) error {
	subStack := NewBranchStack(stack)
	if err := t.Block.evalStatements(t.try, subStack); err != nil {
		return err
	}

	if t.catch != nil {
		subStack := NewBranchStack(stack)

		if t.arg != nil && t.arg.Name() != "_" {
			ctx := t.arg.Context()
			arg := values.NewAllNull(ctx)
			if t.arg.constraint != nil {
				var err error
				arg, err = t.arg.GenerateArgInstance(stack, ctx)
				if err != nil {
					return err
				}
			}

			if !arg.IsInstanceOf(prototypes.Error) {
				return ctx.NewError("Error: expected Error, got " + arg.TypeName())
			}

			subStack.SetValue(t.arg.name.GetVariable(), arg, false, ctx)
		}

		if err := t.Block.evalStatements(t.catch, subStack); err != nil {
			return err
		}
	}

	if t.finally != nil {
		subStack := NewBranchStack(stack)

		if err := t.Block.evalStatements(t.finally, subStack); err != nil {
			return err
		}
	}

	return nil
}

func (t *TryCatch) ResolveStatementActivity(usage Usage) error {
	if t.finally != nil {
		t.statements = t.finally
		if err := t.Block.ResolveStatementActivity(usage); err != nil {
			return err
		}
	}

	if t.catch != nil {
		t.statements = t.catch
		if err := t.Block.ResolveStatementActivity(usage); err != nil {
			return err
		}
	}

	t.statements = t.try
	if err := t.Block.ResolveStatementActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *TryCatch) UniversalStatementNames(ns Namespace) error {
	t.statements = t.try
	if err := t.Block.UniversalStatementNames(ns); err != nil {
		return err
	}

	if t.catch != nil {
		if t.arg != nil {
			if err := t.arg.UniversalNames(ns); err != nil {
				return err
			}
		}

		t.statements = t.catch
		if err := t.Block.UniversalStatementNames(ns); err != nil {
			return err
		}
	}

	if t.finally != nil {
		t.statements = t.finally
		if err := t.Block.UniversalStatementNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *TryCatch) UniqueStatementNames(ns Namespace) error {
	t.statements = t.try
	if err := t.Block.UniqueStatementNames(ns); err != nil {
		return err
	}

	if t.catch != nil {
		if t.arg != nil {
			if err := t.arg.UniqueNames(ns); err != nil {
				return err
			}
		}

		t.statements = t.catch
		if err := t.Block.UniqueStatementNames(ns); err != nil {
			return err
		}
	}

	if t.finally != nil {
		t.statements = t.finally
		if err := t.Block.UniqueStatementNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *TryCatch) Walk(fn WalkFunc) error {
	t.statements = t.try
	if err := t.Block.Walk(fn); err != nil {
		return err
	}

	if t.catch != nil {
		if t.arg != nil {
			if err := t.arg.Walk(fn); err != nil {
				return err
			}
		}

		t.statements = t.catch
		if err := t.Block.Walk(fn); err != nil {
			return err
		}
	}

	if t.finally != nil {
		t.statements = t.finally
		if err := t.Block.Walk(fn); err != nil {
			return err
		}
	}

	return fn(t)
}
