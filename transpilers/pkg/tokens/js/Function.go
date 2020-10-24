package js

import (
	"fmt"
	//"os"
	"reflect"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// implements several interfaces:
//  Statement
//  Expression
//  Function
type Function struct {
	fi      *FunctionInterface
	this    Variable // unique per function (although value might not be!)
	isArrow bool     // arrow function use 'this' from parent
	Block
}

func NewFunction(fi *FunctionInterface, isArrow bool,
	ctx context.Context) (*Function, error) {

	this := NewVariable("this", true, ctx)

	return &Function{fi, this, isArrow, newBlock(ctx)}, nil
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

func (t *Function) writeBody(indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(t.fi.Write())

	if t.isArrow {
		b.WriteString("=>")
	}

	b.WriteString("{")

	s := t.Block.writeBlockStatements(indent+tab, nl)

	if s != "" {
		b.WriteString(nl)
		b.WriteString(s)
		b.WriteString(nl)
		b.WriteString(indent)
	}

	b.WriteString("}")

	return b.String()
}

func (t *Function) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	if t.IsAsync() {
		b.WriteString("async ")
	}

	if !t.isArrow {
		b.WriteString("function ")
		b.WriteString(t.Name())
	}
	b.WriteString(t.writeBody(indent, NL, TAB))

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
	b.WriteString(t.writeBody("", "", ""))

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
	if err := t.fi.ResolveNames(outer, inner); err != nil {
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
		//context.AppendContextString(err, "Info: called here", t.Context())
		return err
	}

	return nil
}

func (t *Function) ResolveExpressionNames(outer Scope) error {
	// wrap the scope
	inner := t.NewScope(outer)

	return t.resolveExpressionNames(outer, inner)
}

func (t *Function) ResolveStatementNames(scope Scope) error {
	if !scope.HasVariable(t.Name()) {
		panic("function should've been hoisted before")
	}

	return t.ResolveExpressionNames(scope)
}

func (t *Function) GenerateArgInstances(stack values.Stack, ctx context.Context) ([]values.Value, error) {
	return t.fi.GenerateArgInstances(stack, ctx)
}

func (t *Function) EvalStatement(stack values.Stack) error {
	// dont do anything, add self during HoistValues
	return nil
}

func (t *Function) HoistValues(stack values.Stack) error {
	fnVal, err := t.EvalExpression(stack)
	if err != nil {
		return err
	}

	ref := t.GetVariable()
	if err := stack.SetValue(ref, fnVal, false, t.Context()); err != nil {
		return err
	}

	return nil
}

func (t *Function) EvalExpression(stack values.Stack) (values.Value, error) {
	// 'this' shouldn't be nil here, but perhaps get it from the stack?
	var thisVal *values.Instance = nil

	if t.isArrow {
		thisVal_, err := stack.GetValue(t.this, t.Context())
		context.AppendContextString(err, "Hint: arrow functions can only be used inside classes", t.Context())
		if err != nil {
			return nil, err
		}

		if thisVal_ == nil {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: this variable declared, but doesn't have a value")
		}

		thisVal = values.AssertInstance(thisVal_)
	}

	return values.NewFunction(t, stack, thisVal, t.Context()), nil
}

func (t *Function) assertLastStatementReturns(lastStatement Statement,
	retValue values.Value) error {
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

			if err := t.assertLastStatementReturns(groupStatements[len(groupStatements)-1], retValue); err != nil {
				return err
			}
		}
	case *While:
		if !IsLiteralTrue(st.cond) {
			errCtx := st.Context()
			return errCtx.NewError("Error: while as final returning statement only makes sense for infinite loop")
		}

		if err := t.assertLastStatementReturns(st.statements[len(st.statements)-1], retValue); err != nil {
			return err
		}
	case *For:
		if !IsLiteralTrue(st.cond) {
			errCtx := st.Context()
			return errCtx.NewError("Error: for as final returning statement only makes sense for infinite loop")
		}

		if err := t.assertLastStatementReturns(st.statements[len(st.statements)-1], retValue); err != nil {
			return err
		}
	default:
		errCtx := lastStatement.Context()
		return errCtx.NewError("Error: missing return statement (hint: return statement must come last in every branch)")
	}

	return nil
}

func (t *Function) evalInternalSync(parent values.Stack, this *values.Instance,
	args []values.Value, unlockThis bool, ctx context.Context) (values.Value, error) {
	// check recursiveness, must be done before this is set!
	if isRec, retVal := parent.IsRecursive(t); isRec {
		if retVal != nil {
			return retVal, nil
		} else {
			return values.NewAllNull(t.Context()), nil
		}
	}

	// function stack contains new values every time the function is invoked
	stack := NewFunctionStack(t, parent) // this is the stack we need for this mutation

	if this == nil {
		if err := stack.SetValue(t.GetThisVariable(), nil, false,
			ctx); err != nil {
			return nil, err
		}
	} else {
		if err := stack.SetValue(t.GetThisVariable(), this, false,
			ctx); err != nil {
			return nil, err
		}

		if unlockThis {
			this.UnlockForStack(stack)
		}
	}

	if err := t.fi.EvalArgs(stack, args, ctx); err != nil {
		return nil, err
	}

	if err := t.Block.HoistValues(stack); err != nil {
		return nil, err
	}

	if err := t.Block.EvalStatement(stack); err != nil {
		return nil, err
	}

	if stack.returnValue != nil && !values.IsVoid(stack.returnValue) {
		if err := t.assertLastStatementReturns(t.statements[len(t.statements)-1],
			stack.returnValue); err != nil {
			return nil, err
		}
	}

	return t.fi.ConstrainReturnValue(stack, stack.returnValue, ctx)
}

// TODO: clean this shit up
func (t *Function) evalInternal(stack values.Stack, this *values.Instance,
	args []values.Value, unlockThis bool, ctx context.Context) (values.Value, error) {
	ret, err := t.evalInternalSync(stack, this, args, unlockThis, ctx)
	if err != nil {
		if ar, ok := err.(*prototypes.AsyncRequest); ok {
			if !t.IsAsync() {
				errCtx := t.Context()
				err := errCtx.NewError("not an async function, should've been caught during resolve names stage")
				panic(err)
			}

			resolveFn := values.NewFunctionFunction(func(
				stack_ values.Stack, this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {
				// XXX: or should we use stack_?
				asyncStack, err := NewAsyncStack(stack, ar.GetAwait(), args_)
				if err != nil {
					return nil, err
				}

				ret_, err_ := t.evalInternal(asyncStack, this, args, unlockThis, ctx)
				if err_ != nil {
					return nil, err_
				}

				if ret_.IsInstanceOf(prototypes.Promise) {
					ret_ = values.UnpackContextValue(ret_)
					if ret, ok := ret_.(*values.Instance); ok {
						// unpack the promise
						retProps := values.AssertPromiseProperties(ret.Properties())
						return retProps.ResolveAwait()
					} else {
						panic("unexpected")
					}
				}

				return ret_, nil
			}, stack, this, ctx)

			ar.SetResolveFn(resolveFn)

			props := values.NewPromiseProperties(ctx)
			props.SetRejectArgs([]values.Value{prototypes.NewError(ctx)})
			promise := values.NewInstance(prototypes.Promise, props, ctx)
			return promise, nil
		} else {
			return nil, err
		}
	}

	if t.IsAsync() {
		props := values.NewPromiseProperties(ctx)

		if !values.IsVoid(ret) {
			if err := props.SetResolveArgs([]values.Value{ret}, ctx); err != nil {
				return nil, err
			}
		} else {
			if err := props.SetResolveArgs([]values.Value{}, ctx); err != nil {
				return nil, err
			}
		}

		// assume some error will be thrown at some point
		props.SetRejectArgs([]values.Value{prototypes.NewError(ctx)})

		promise := values.NewInstance(prototypes.Promise, props, ctx)
		return promise, nil
	} else {
		return ret, nil
	}
}

// implement values.Callable interface
func (t *Function) evalFunction(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	ret, err := t.evalInternal(stack, this, args, false, ctx)
	if err != nil {
		return nil, err
	}

	if ret == nil {
		return nil, ctx.NewError("Error: function doesn't return a value")
	} else if values.IsVoid(ret) {
		// XXX: recursive value might be caught higher up?
		//return ret, nil
		return nil, ctx.NewError(
			fmt.Sprintf("Error: function doesn't return a value (%p, %s, %s)", ret, reflect.TypeOf(ret).String(), this.TypeName()))
	}

	if values.IsAllNull(ret) {
		// AllNull might be returned if recursiveness is detected, try generating an instance in this case
		retAlt, err := t.fi.GenerateReturnInstance(stack, ctx)
		if err == nil {
			ret = retAlt
		} else {
			//if VERBOSITY >= 2 {
			//infoCtx := t.Context()
			//fmt.Fprintf(os.Stderr, "Warning: unable to generate return value instance in recursive call (%s)",
			//infoCtx.NewError("").Error())
			//}
		}
	}

	return values.NewContextValue(ret, t.Context()), nil
}

func (t *Function) EvalFunction(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	return t.evalFunction(stack, this, args, ctx)
}

func (t *Function) EvalFunctionNoReturn(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	ret, err := t.evalInternal(stack, this, args, false, ctx)
	if err != nil {
		return nil, err
	}

	if ret != nil && ret.IsVoid() {
		return nil, nil
	}

	return ret, nil
}

func (t *Function) evalMethod(stack values.Stack, this *values.Instance,
	args []values.Value, unlockThis bool, ctx context.Context) error {
	ret, err := t.evalInternal(stack, this, args, unlockThis, ctx)
	if err != nil {
		return err
	}

	if !(ret == nil || ret.IsVoid()) {
		if !prototypes.IsAsync(t) { // return value could be promise, but it could be a shoot and forget call
			errCtx := ret.Context()
			return errCtx.NewError("Error: unexpected return value")
		}
	}

	return nil
}

func (t *Function) evalFunctionAsEntryPoint(stack values.Stack, this *values.Instance,
	ctx context.Context) (values.Value, error) {
	args, err := t.GenerateArgInstances(stack, ctx)
	if err != nil {
		return nil, err
	}

	return t.evalFunction(stack, this, args, ctx)
}

func (t *Function) evalMethodAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) error {
	args, err := t.GenerateArgInstances(stack, ctx)
	if err != nil {
		return err
	}

	return t.evalMethod(stack, this, args, false, ctx)
}

func (t *Function) EvalMethod(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) error {
	return t.evalMethod(stack, this, args, false, ctx)
}

func (t *Function) EvalAsConstructor(stack values.Stack, this *values.Instance, args []values.Value,
	ctx context.Context) error {
	return t.evalMethod(stack, this, args, true, ctx)
}

func (t *Function) EvalAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) error {
	if _, err := t.evalAsEntryPoint(stack, this, ctx); err != nil {
		return err
	}

	return nil
}

func (t *Function) evalAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) (values.Value, error) {
	args, err := t.GenerateArgInstances(stack, ctx)
	if err != nil {
		return nil, err
	}

	return t.evalInternal(stack, this, args, false, ctx)
}

func (t *Function) ResolveExpressionActivity(usage Usage) error {
	tmp := usage.InFunction()
	usage.SetInFunction(true)

	err := t.Block.ResolveStatementActivity(usage)

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
