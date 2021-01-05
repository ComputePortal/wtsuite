package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type Function struct {
	fi      *FunctionInterface
	this    Variable // unique per function (although value might not be!)
	isArrow bool     // arrow function use 'this' from parent
  ret     []*Return // registered via FunctionScope
	Block
}

func NewFunction(fi *FunctionInterface, isArrow bool,
	ctx context.Context) (*Function, error) {

	this := NewVariable("this", true, ctx)

	return &Function{fi, this, isArrow, make([]*Return, 0), newBlock(ctx)}, nil
}

func (t *Function) NewScope(parent Scope) *FunctionScope {
	return NewFunctionScope(t, parent)
}

func (t *Function) Name() string {
	return t.fi.Name()
}

func (t *Function) Length() int {
	return t.fi.Length()
}

func (t *Function) GetVariable() Variable {
	return t.fi.GetVariable()
}

func (t *Function) Role() prototypes.FunctionRole {
	return t.fi.Role()
}

func (t *Function) GetThisVariable() Variable {
	return t.this
}

func (t *Function) Interface() *FunctionInterface {
	return t.fi
}

func (t *Function) IsAsync() bool {
	return prototypes.IsAsync(t)
}

func (t *Function) IsVoid() bool {
  return t.fi.IsVoid()
}

func (t *Function) RegisterReturn(ret *Return) {
  t.ret = append(t.ret, ret)
}

func (t *Function) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Function")

	if t.Name() != "" {
		b.WriteString(" ")
	}

	b.WriteString(t.fi.Dump())

	for _, st := range t.statements {
		b.WriteString(st.Dump(indent + "{ "))
	}

	return b.String()
}

func (t *Function) writeBody(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(t.fi.Write())

	if t.isArrow {
		b.WriteString("=>")
	}

	b.WriteString("{")

	s := t.Block.writeBlockStatements(usage, indent+tab, nl, tab)

	if s != "" {
		b.WriteString(nl)
		b.WriteString(s)
		b.WriteString(nl)
		b.WriteString(indent)
	}

	b.WriteString("}")

	return b.String()
}

func (t *Function) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	if t.IsAsync() {
		b.WriteString("async ")
	}

	if !t.isArrow {
		b.WriteString("function ")
		b.WriteString(t.Name())
	}
  b.WriteString(t.writeBody(usage, indent, nl, tab))

	return b.String()
}

func (t *Function) WriteExpression() string {
	// named function expression are only really useful for runtime debugging, which are trying to entirely avoid with this new language

	var b strings.Builder

	if t.IsAsync() {
		b.WriteString("async ")
	}
	if !t.isArrow {
		b.WriteString("function")
	}
	b.WriteString(t.writeBody(nil, "", "", ""))

	return b.String()
}

func (t *Function) HoistNames(scope Scope) error {
	if scope.HasVariable(t.Name()) {
		errCtx := t.Context()
		return errCtx.NewError("Error: \"" + t.Name() + "\" already defined")
	}

	return scope.SetVariable(t.Name(), t.GetVariable())
}

func (t *Function) resolveExpressionNames(outer Scope, inner Scope) error {
	if err := t.fi.ResolveNames(inner); err != nil {
		return err
	}

	if t.isArrow && outer.HasVariable("this") {
		this, err := outer.GetVariable("this")
		if err != nil {
			return err
		}
		t.this = this
	}

	if err := inner.SetVariable("this", t.this); err != nil {
		return err
	}

	if err := t.Block.HoistAndResolveStatementNames(inner); err != nil {
		return err
	}

  // set the value of the variable right
	return nil
}

func (t *Function) ResolveExpressionNames(outer Scope) error {
	// wrap the scope
	inner := t.NewScope(outer)

  if err := t.resolveExpressionNames(outer, inner); err != nil {
    return err
  }

  // register value to variable so it is available in a scope that used the hoisted variable
  // must be done here so value is available in class
  fn, err := t.GetFunctionValue()
  if err != nil {
    return err
  }

  variable := t.GetVariable()
  variable.SetValue(fn)

  return nil
}

func (t *Function) ResolveStatementNames(scope Scope) error {
	if !scope.HasVariable(t.Name()) {
		panic("function should've been hoisted before")
	}

  if err := t.ResolveExpressionNames(scope); err != nil {
    return err
  }

	return nil
}

func (t *Function) GetFunctionValue() (*values.Function, error) {
  return t.fi.GetFunctionValue()
}

// pre async
func (t *Function) getReturnValue() (values.Value, error) {
  return t.fi.getReturnValue()
}

// post async
// return type can be nil in case of void
func (t *Function) GetReturnValue() (values.Value, error) {
  return t.fi.GetReturnValue()
}

func (t *Function) GetArgValues() ([]values.Value, error) {
	return t.fi.GetArgValues()
}

func (t *Function) assertLastStatementReturns(lastStatement Statement) error {
	switch st := lastStatement.(type) {
	case *Return:
		return nil
	case *Throw:
		return nil
	case *If:
		if st.conds[len(st.conds)-1] != nil {
			errCtx := st.Context()
			return errCtx.NewError("Error: not every branch returns a value")
		}

		for i, _ := range st.conds {
			groupStatements := st.grouped[i]

			if err := t.assertLastStatementReturns(groupStatements[len(groupStatements)-1]); err != nil {
				return err
			}
		}
	case *While:
		if !IsLiteralTrue(st.cond) {
			errCtx := st.Context()
			return errCtx.NewError("Error: while as final returning statement only makes sense for infinite loop")
		}

		if err := t.assertLastStatementReturns(st.statements[len(st.statements)-1]); err != nil {
			return err
		}
	case *For:
		if !IsLiteralTrue(st.cond) {
			errCtx := st.Context()
			return errCtx.NewError("Error: for as final returning statement only makes sense for infinite loop")
		}

		if err := t.assertLastStatementReturns(st.statements[len(st.statements)-1]); err != nil {
			return err
		}
	default:
		errCtx := lastStatement.Context()
		return errCtx.NewError("Error: missing return statement (hint: return statement must come last in every branch)")
	}

	return nil
}

func (t *Function) EvalExpression() (values.Value, error) {
  if err := t.fi.Eval(); err != nil {
    return nil, err
  }

	if err := t.Block.EvalStatement(); err != nil {
		return nil, err
	}

  // pre async (post async might give Promise<void>)
  retVal, err := t.fi.getReturnValue()
  if err != nil {
    return nil, err
  }

  if retVal != nil {
    n := len(t.statements)
    if n == 0 {
      errCtx := t.Context()
      return nil, errCtx.NewError("Error: expected return statement")
    } else {
      if err := t.assertLastStatementReturns(t.statements[n-1]); err != nil {
        return nil, err
      }
    }

    /*if len(t.ret) == 0 {
      errCtx := t.Context()
      return nil, errCtx.NewError("Error: expected non-void return value, but no return statement found")
    }*/
  }

  // function value was created before
  variable := t.GetVariable()
  return variable.GetValue(), nil
}

func (t *Function) EvalStatement() error {
	_, err := t.EvalExpression()
	if err != nil {
		return err
	}

	return nil
}

func (t *Function) ResolveExpressionActivity(usage Usage) error {
	tmp := usage.InFunction()
	usage.SetInFunction(true)

	err := t.Block.ResolveStatementActivity(usage)

  // TODO: fi.arg defaults?

	usage.SetInFunction(tmp)

	if err != nil {
		return err
	}

	if err := usage.DetectUnused(); err != nil {
		return err
	}

	return nil
}

func (t *Function) ResolveStatementActivity(usage Usage) error {
	if usage.InFunction() {
		ref := t.GetVariable()

		if err := usage.Rereference(ref, t.Context()); err != nil {
			return err
		}
	}

	return t.ResolveExpressionActivity(usage)
}

func (t *Function) UniversalExpressionNames(ns Namespace) error {
	if err := t.fi.UniversalNames(ns); err != nil {
		return err
	}

	return t.Block.UniversalStatementNames(ns)
}

func (t *Function) UniqueExpressionNames(ns Namespace) error {
	subNs := ns.NewFunctionNamespace()

	if err := t.fi.UniqueNames(subNs); err != nil {
		return err
	}

	return t.Block.UniqueStatementNames(subNs)
}

func (t *Function) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *Function) UniqueStatementNames(ns Namespace) error {
	ns.FunctionName(t.GetVariable())

	return t.UniqueExpressionNames(ns)
}

func (t *Function) Walk(fn WalkFunc) error {
  if err := t.fi.Walk(fn); err != nil {
    return err
  }

  if err := t.Block.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
