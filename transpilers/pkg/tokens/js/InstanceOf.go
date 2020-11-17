package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// in a file by itself because it is more complex than the typical operator
type InstanceOf struct {
	BinaryOp

	interf       values.Interface // starts as nil
}

func NewInstanceOf(a Expression, b Expression, ctx context.Context) *InstanceOf {
	return &InstanceOf{
		BinaryOp{"instanceof", a, b, TokenData{ctx}},
		nil,
	}
}

func (t *InstanceOf) WriteExpression() string {
	var b strings.Builder

	b.WriteString("(")

	a := t.a.WriteExpression()

	firstDone := false
	if t.interf != nil && t.interf.Name() == "String" {
		b.WriteString("typeof(")
		b.WriteString(a)
		b.WriteString(")==='string'")
		firstDone = true
	} else if t.interf != nil && t.interf.Name() == "Int" {
		if firstDone {
			b.WriteString("||")
		}

		b.WriteString("Number.isInteger(")
		b.WriteString(a)
		b.WriteString(")")
		firstDone = true
	} else if t.interf != nil && t.interf.Name() == "Number" {
		if firstDone {
			b.WriteString("||")
		}
		b.WriteString("typeof(")
		b.WriteString(a)
		b.WriteString(")==='number'")
		firstDone = true
	} else if t.interf != nil && t.interf.Name() == "Boolean" {
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
    protos, err := t.interf.GetPrototypes()
    if err != nil {
      panic("should've been caught before")
    }
    // check if interf itself is included
    if interfProto, ok := t.interf.(values.Prototype); ok {
      selfIncluded := false
      for _, proto := range protos {
        if proto == interfProto {
          selfIncluded = true
          break
        }
      }

      if !selfIncluded {
        protos = append(protos, interfProto)
      }
    }

		if len(protos) == 0 {
			b.WriteString("false")
		} else {
			for i, proto := range protos {
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
    t.interf = b.GetInterface()
	}

	return nil
}

func (t *InstanceOf) evalInternal() error {
	a, err := t.a.EvalExpression()
	if err != nil {
		return err
	}

  if !values.IsInstance(a) {
    errCtx := t.a.Context()
    return errCtx.NewError("Error: not an instance")
  }

  if t.interf == nil {
    b, err := t.b.EvalExpression()
    if err != nil {
      return err
    }

    if !values.IsClass(b) {
      errCtx := t.b.Context()
      return errCtx.NewError("Error: not a class or interface")
    }
  }

	return nil
}

func (t *InstanceOf) EvalExpression() (values.Value, error) {
	if err := t.evalInternal(); err != nil {
    return nil, err
  }

	return prototypes.NewBoolean(t.Context()), nil
}

func (t *InstanceOf) CollectTypeGuards(c map[Variable]values.Interface) (bool, error) {
	// only if lhs is VarExpression
	if lhs, ok := t.a.(*VarExpression); ok {
		ref := lhs.GetVariable()

		if err := t.evalInternal(); err != nil {
			return false, err
		}

		// only if rhs is a single interface/class
		if t.interf != nil { // in case of multiple or no classes
			if _, ok := c[ref]; !ok {
				c[ref] = t.interf
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
