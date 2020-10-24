package js

import (
	"reflect"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// very similar to VarExpression, but with extra
type TypeExpression struct {
  content []*TypeExpressionMember // can be nil, can't be empty
	interf        values.Interface  // starts as nil, evaluated later
	VarExpression                   // the base type will be a class variable
}

type TypeExpressionMember struct {
  key *Word // can be nil
  texpr *TypeExpression // can't be nil
}

func NewTypeExpression(name string, contentKeys []*Word, contentTypes []*TypeExpression,
	ctx context.Context) *TypeExpression {

	if contentKeys != nil && contentTypes != nil && len(contentKeys) != len(contentTypes) {
		panic("contentKeys and contentTypes should have same length")
	}

  var content []*TypeExpressionMember = nil
  if contentTypes != nil {
    if len(contentTypes) == 0 {
      panic("contentTypes can't be empty (but can be nil)")
    }

    content = make([]*TypeExpressionMember, len(contentTypes))

    for i, contentType := range contentTypes {
      if contentKeys == nil {
        content[i] = &TypeExpressionMember{
          nil,
          contentType,
        }
      } else {
        content[i] = &TypeExpressionMember{
          contentKeys[i],
          contentType,
        }
      }
    }
  }

	return &TypeExpression{content, nil, newVarExpression(name, true, ctx)}
}

func (t *TypeExpression) hasKeys() bool {
  if t.content != nil {
    return t.content[0].key != nil
  }

  return false
}

func (t *TypeExpression) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	//b.WriteString("Type(")
	b.WriteString(t.VarExpression.Name())

	if t.hasKeys() {
		b.WriteString("<")
		for _, cont := range t.content {
			contentType := cont.texpr

			b.WriteString(cont.key.Value())
			b.WriteString(":")
			b.WriteString(contentType.Dump(""))
			b.WriteString(",")
		}

		b.WriteString(">")
	} else if t.content != nil {
		b.WriteString("<")

		for _, cont := range t.content {
			b.WriteString(cont.texpr.Dump(""))
		}

		b.WriteString(">")
	}

	b.WriteString("\n")

	return b.String()
}

func (t *TypeExpression) ResolveExpressionNames(scope Scope) error {
	if t.Name() == "any" || t.Name() == "function" || t.Name() == "class" || t.Name() == "void" {
		return nil
	}

	if err := t.VarExpression.ResolveExpressionNames(scope); err != nil {
		return err
	}

	if t.content != nil {
		for _, cont := range t.content {
			if err := cont.texpr.ResolveExpressionNames(scope); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TypeExpression) EvalExpression(stack values.Stack) (values.Value, error) {
	if t.Name() == "any" || t.Name() == "function" || t.Name() == "class" || t.Name() == "void" {
		if t.content != nil {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: " + t.Name() + " can't have content types")
		}
		return values.NewClass(values.NewDummyPrototype(t.Name()), t.Context()), nil
	} else {
		return t.VarExpression.EvalExpression(stack)
	}
}

func (t *TypeExpression) Walk(fn WalkFunc) error {
  if t.content != nil {
    for _, cont := range t.content {
      if err := cont.Walk(fn); err != nil {
        return err
      }
    }
  }

  if err := t.VarExpression.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func (t *TypeExpressionMember) Walk(fn WalkFunc) error {
  if t.key != nil {
    if err := t.key.Walk(fn); err != nil {
      return err
    }
  }

  if err := t.texpr.Walk(fn); err != nil {
    return err
  }

  return nil
}

func (t *TypeExpression) GenerateInstance(stack values.Stack, ctx context.Context) (values.Value, error) {
	if t.Name() == "class" {
		return nil, ctx.NewError("Error: can't generate class")
	}

	if t.Name() == "any" {
		return nil, ctx.NewError("Error: can't generate generic instance")
	}

	if t.Name() == "function" {
		return nil, ctx.NewError("Error: can't generate function")
	}

	cl, err := stack.GetValue(t.ref, t.Context())
	if err != nil {
		return nil, err
	}

	if cl == nil {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: class variable declared, but doesn't have a value")
	}

	var keys []string = nil
	var args []values.Value = nil

	if t.hasKeys() {
		keys = make([]string, 0)
		for _, cont := range t.content {
			keys = append(keys, cont.key.Value())
		}
	}

	if t.content != nil {
		args = make([]values.Value, 0)
		for _, cont := range t.content {
			contentArg, err := cont.texpr.GenerateInstance(stack, t.Context())
			if err != nil {
				return nil, err
			}

			args = append(args, contentArg)
		}
	}

	if cl.IsClass() {
		proto_, ok := cl.GetClassPrototype()
		if !ok {
			panic("unexpected")
		}

		switch proto := proto_.(type) {
		case *prototypes.BuiltinPrototype:
		case *Class:
			if proto.IsAbstract() {
				errCtx := t.Context()
				return nil, errCtx.NewError("Error: can't instantiate an abstract class")
			}
		}

		res, err := proto_.GenerateInstance(stack, keys, args, t.Context())
		if err != nil {
			return nil, err
		}

		return values.NewContextValue(res, t.Context()), nil
	} else if cl.IsInterface() {
		ci, ok := cl.GetClassInterface()
		if !ok {
			panic("unexpected")
		}

		return ci.GenerateInstance(stack, keys, args, t.Context())
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: not a class nor an explicit interface")
	}
}

func (t *TypeExpression) CollectNestedTypes(stack values.Stack) (*values.NestedType, error) {

	var interf values.Interface = nil
	if t.Name() == "any" || t.Name() == "class" || t.Name() == "function" || t.Name() == "void" {
		interf = values.NewDummyPrototype(t.Name())
	} else {
		classVal, err := t.EvalExpression(stack)
		if err != nil {
			return nil, err
		}

		var ok bool
		interf, ok = classVal.GetClassInterface()
		if !ok {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: contraint is not an interface (" + reflect.TypeOf(values.UnpackContextValue(classVal)).String() + ")")
		}
	}

	t.interf = interf
	var children []*values.NestedType = nil

	if t.content != nil {
		children = make([]*values.NestedType, len(t.content))
		for i, cont := range t.content {
			nested, err := cont.texpr.CollectNestedTypes(stack)
			if err != nil {
				return nil, err
			}

			if cont.key != nil {
				nested.SetKey(cont.key.Value())
			}

			children[i] = nested
		}
	}

	return values.NewNestedType("", interf, children, t.Context()), nil
}

func (t *TypeExpression) Constrain(stack values.Stack, v values.Value) (values.Value, error) {
	nestedType, err := t.CollectNestedTypes(stack)
	if err != nil {
		return nil, err
	}

  res, err := v.Cast(nestedType, t.Context())

  // attempt to generate instance for AllNull, don't bother for typed nulls (those should'be been created by casting an AllNull anyway)
  if values.IsAllNull(v) {
    genRes, err := t.GenerateInstance(stack, v.Context())
    if err == nil {
      return genRes, nil
    }
  }

  return res, nil
}

// the dual of Constrain() for Function<...> types (void part is handled in FunctionInterface)
/*func (t *TypeExpression) CheckInterface(wanted *values.NestedType, ctx context.Context) error {
	if t.Name() == "any" && wanted.Name() != "any" {
		errCtx := wanted.Context()
		return errCtx.NewError("Error: have any, want " + wanted.Name())
	}

	if t.Name() != "any" && wanted.Name() == "any" {
		errCtx := wanted.Context()
		return errCtx.NewError("Error: have " + t.Name() + ", want " + wanted.Name())
	}

	// types should be identical, but we can't check based on name
	if t.interf != wanted.Interface() {
		errCtx := wanted.Context()
		return errCtx.NewError("Error: have " + t.Name() + ", want " + wanted.Name())
	}

	wantedChildren := wanted.Children()
	if wantedChildren != nil {
		if t.contentKeys == nil {
			errCtx := wanted.Context()
			return errCtx.NewError("Error: don't have content, want content")
		}

		if len(wantedChildren) != len(t.contentTypes) {
			errCtx := wanted.Context()
			return errCtx.NewError(fmt.Sprintf("Error: have %d content type, want %d content types", len(t.contentTypes), len(wantedChildren)))
		}

		if t.contentKeys != nil {
			// check that all the keys are the same and exist
			for i, key := range t.contentKeys {
				found := false
				for _, wantedType := range wantedChildren {

					if key.Value() == wantedType.Key() {
						found = true

						if err := t.contentTypes[i].CheckInterface(wantedType, ctx); err != nil {
							return err
						}
						break
					}
				}

				if !found {
					errCtx := wanted.Context()
					return errCtx.NewError("Error: subtype " + key.Value() + " not wanted")
				}
			}
		} else {
			for i, contentType := range t.contentTypes {
				if err := contentType.CheckInterface(wantedChildren[i], ctx); err != nil {
					return err
				}
			}
		}
	}

	return nil
}*/
