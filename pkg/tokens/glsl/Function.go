package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
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

func (t *Function) GetVariable() Variable {
	return t.fi.GetVariable()
}

func (t *Function) NewScope(parent Scope) *FunctionScope {
	return NewFunctionScope(t, parent)
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
  if !usage.IsUsed(t.GetVariable()) {
    return ""
  }

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

func (t *Function) RegisterReturn(ret *Return) {
  t.ret = append(t.ret, ret)
}

func (t *Function) ResolveStatementNames(outer Scope) error {
  // default SetVariable method isnt good enough
	if outer.HasVariable(t.Name()) {
		errCtx := t.Context()
		return errCtx.NewError("Error: \"" + t.Name() + "\" already defined")
	}

  variable := t.GetVariable()
  variable.SetConstant()

  if err := outer.SetVariable(t.Name(), variable); err != nil {
    return err
  }

	// wrap the scope
	inner := t.NewScope(outer)

	if err := t.fi.ResolveNames(inner); err != nil {
		return err
	}

	if err := t.Block.ResolveStatementNames(inner); err != nil {
		return err
	}

  fnVal := values.NewFunction(t, t.Context())
  variable.SetValue(fnVal)

  return nil
}

func (t *Function) assertLastStatementReturns(lastStatement Statement) error {
	switch st := lastStatement.(type) {
	case *Return:
		return nil
	case *If:
		if st.conds[len(st.conds)-1] != nil {
			errCtx := st.Context()
			return errCtx.NewError("Error: not every branch returns a value")
		}

		for i, _ := range st.conds {
			groupStatements := st.groups[i].statements

			if err := t.assertLastStatementReturns(groupStatements[len(groupStatements)-1]); err != nil {
				return err
			}
		}
	default:
		errCtx := lastStatement.Context()
		return errCtx.NewError("Error: missing return statement (hint: return statement must come last in every branch)")
	}

	return nil
}

func (t *Function) EvalStatement() error {
  retVal, err := t.fi.EvalCall(nil, t.Context())
  if err != nil {
    return err
  }

  if err := t.Block.evalStatements(); err != nil {
    return err
  }

  if retVal != nil {
    n := len(t.statements)
    if n == 0 {
      errCtx := t.Context()
      return errCtx.NewError("Error: expected return statement")
    } else {
      if err := t.assertLastStatementReturns(t.statements[n-1]); err != nil {
        return err
      }
    }

    if len(t.ret) == 0 {
      errCtx := t.Context()
      return errCtx.NewError("Error: expected non-void return value, but no return statement found")
    }
  }

  variable := t.GetVariable()
  variable.SetValue(values.NewFunction(t, t.Context()))

  return nil
}

// args == nil allows the caller to retrieve only the return value
func (t *Function) EvalCall(args []values.Value, ctx context.Context) (values.Value, error) {
  return t.fi.EvalCall(args, ctx)
}

func (t *Function) ResolveStatementActivity(usage Usage) error {
  if !usage.IsUsed(t.GetVariable()) {
    return nil
  }

  if err := t.Block.ResolveStatementActivity(usage); err != nil {
    return err
  }

  // fi args dont make a difference

  return nil
}

func (t *Function) UniqueStatementNames(ns Namespace) error {
  // if main: then should've already been set
  ns.FunctionName(t.GetVariable())

  subNs := ns.NewFunctionNamespace()

  if err := t.fi.UniqueNames(subNs); err != nil {
    return err
  }

  return t.Block.UniqueStatementNames(subNs)
}
