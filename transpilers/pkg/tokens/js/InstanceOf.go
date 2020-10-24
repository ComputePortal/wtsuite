package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// in a file by itself because it is more complex than the typical operator
// XXX: there is not yet an equivalent for checking for interfaces

type InstanceOf struct {
	BinaryOp
	maybeString  bool
	maybeInt     bool
	maybeNumber  bool
	maybeBoolean bool
	interf       values.Interface // starts as nil
}

func NewInstanceOf(a Expression, b Expression, ctx context.Context) *InstanceOf {
	return &InstanceOf{
		BinaryOp{"instanceof", a, b, TokenData{ctx}},
		false,
		false,
		false,
		false,
		nil,
	}
}

func (t *InstanceOf) WriteExpression() string {
	var b strings.Builder

	b.WriteString("(")

	a := t.a.WriteExpression()

	firstDone := false
	if t.maybeString {
		b.WriteString("typeof(")
		b.WriteString(a)
		b.WriteString(")==='string'")
		firstDone = true
	}

	if t.maybeInt && !t.maybeNumber {
		if firstDone {
			b.WriteString("||")
		}

		b.WriteString("Number.isInteger(")
		b.WriteString(a)
		b.WriteString(")")
		firstDone = true
	} else if t.maybeNumber {
		if firstDone {
			b.WriteString("||")
		}
		b.WriteString("typeof(")
		b.WriteString(a)
		b.WriteString(")==='number'")
		firstDone = true
	}

	if t.maybeBoolean {
		if firstDone {
			b.WriteString("||")
		}
		b.WriteString("typeof(")
		b.WriteString(a)
		b.WriteString(")==='boolean'")
		firstDone = true
	}

	if t.interf == nil {
		if firstDone {
			b.WriteString("||")
		}
		b.WriteString(a)
		b.WriteString(" instanceof ")
		b.WriteString(t.b.WriteExpression())
	} else {
		interf, ok := t.interf.(*ClassInterface)
		if !ok || interf.explicitImplements == nil {
			panic("unexpected")
		}

		if len(interf.explicitImplements) == 0 {
			b.WriteString("false")
		} else {
			for i, proto := range interf.explicitImplements {
				if i != 0 || firstDone {
					b.WriteString("||")
				}
				b.WriteString(a)
				b.WriteString(" instanceof ")
				b.WriteString(proto.Name())
			}
		}
	}

	b.WriteString(")")

	return b.String()
}

func (t *InstanceOf) ResolveExpressionNames(scope Scope) error {
	if err := t.a.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if err := t.b.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if b, ok := t.b.(*VarExpression); ok {
		ci_ := b.GetVariable().GetObject()

		if ci_ != nil {
			if ci, ok := ci_.(*ClassInterface); ok {
				t.interf = ci
			}
		}
	}

	return nil
}

func (t *InstanceOf) evalExpression(stack values.Stack) (values.Value, values.Interface, error) {
	a, err := t.a.EvalExpression(stack)
	if err != nil {
		return nil, nil, err
	}

	// lhs must be an instance
	if !a.IsInstance() {
		errCtx := a.Context()
		return nil, nil, errCtx.NewError("Error: not an instance")
	}

	// rhs must be a concrete prototype (not an interface!)
	b, err := t.b.EvalExpression(stack)
	if err != nil {
		return nil, nil, err
	}

	var interf values.Interface = nil
	if b.IsInterface() {
		// the interface should be the same for each call to EvalExpression
		var ok bool
		interf, ok = b.GetClassInterface()
		if !ok {
			panic("unexpected")
		}

		if t.interf != nil && t.interf != interf {
			panic("interface differs, interface cannot be used as value!")
		}

		classInterf, ok := interf.(*ClassInterface)
		if !ok || classInterf.explicitImplements == nil {
			errCtx := t.b.Context()
			return nil, nil, errCtx.NewError("Error: not an explicitly implemented interface")
		}
	} else if b.IsClass() {
		bClasses := values.UnpackMulti([]values.Value{b})

		if t.interf != nil {
			panic("previously interface, now class?")
		}

		for _, bClass := range bClasses {
			bProto, ok := bClass.GetClassPrototype()
			if !ok {
				panic("unexpected")
			}

			if len(bClasses) == 1 {
				interf = bProto
			}

			if bProto.HasAncestor(prototypes.String) {
				t.maybeString = true
			}

			if bProto.HasAncestor(prototypes.Int) {
				// beware that whole Numbers still evaluate to Int!
				t.maybeInt = true
			} else if bProto.HasAncestor(prototypes.Number) {
				t.maybeNumber = true
			}

			if bProto.HasAncestor(prototypes.Boolean) {
				t.maybeBoolean = true
			}
		}
	} else {
		errCtx := b.Context()
		return nil, nil, errCtx.NewError("Error: not a class or an interface")
	}

	return prototypes.NewInstance(prototypes.Boolean, t.Context()), interf, nil
}

func (t *InstanceOf) EvalExpression(stack values.Stack) (values.Value, error) {
	v, _, err := t.evalExpression(stack)
	return v, err
}

func (t *InstanceOf) CollectTypeGuards(stack values.Stack, c map[interface{}]values.Interface) (bool, error) {
	// only if lhs is VarExpression
	if lhs, ok := t.a.(*VarExpression); ok {
		ref := lhs.GetVariable()

		_, interf, err := t.evalExpression(stack)
		if err != nil {
			return false, err
		}

		// only if rhs is a single interface/class
		if interf != nil { // in case of multiple or no classes
			if _, ok := c[ref]; !ok {
				c[ref] = interf
				return true, nil
			} // else: c already contains another type guard for the same variable -> void all
		}
	}

	// evalExpression wil be called a second time elsewhere
	return false, nil
}

func (t *InstanceOf) Walk(fn WalkFunc) error {
  if err := t.BinaryOp.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
