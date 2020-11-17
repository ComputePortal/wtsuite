package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// very similar to VarExpression, but with extra
type TypeExpression struct {
  parameters []*TypeExpressionMember // can be nil, can't be empty
	interf        values.Interface  // starts as nil, evaluated later
	VarExpression                   // the base type will be a class variable
}

type TypeExpressionMember struct {
  key *Word // can be nil
  typeExpr *TypeExpression // can't be nil
}

func NewTypeExpression(name string, parameterKeys []*Word, 
  parameterTypes []*TypeExpression, ctx context.Context) (*TypeExpression, error) {
	if parameterKeys != nil && parameterTypes != nil && len(parameterKeys) != len(parameterTypes) {
		panic("parameterKeys and parameterTypes should have same length")
	}

  var parameters []*TypeExpressionMember = nil
  if parameterTypes != nil {
    if len(parameterTypes) == 0 {
      errCtx := ctx
      return nil, errCtx.NewError("Error: parameter types can't be empty")
    }

    parameters = make([]*TypeExpressionMember, len(parameterTypes))

    for i, parameterType := range parameterTypes {
      if parameterKeys == nil {
        parameters[i] = &TypeExpressionMember{
          nil,
          parameterType,
        }
      } else {
        parameters[i] = &TypeExpressionMember{
          parameterKeys[i],
          parameterType,
        }
      }
    }
  }

	return &TypeExpression{parameters, nil, newVarExpression(name, true, ctx)}, nil
}

func (t *TypeExpression) hasKeys() bool {
  if t.parameters != nil {
    if t.parameters[0].key != nil {
      return true
    } else {
      return false
    }
  } else {
    return false
  }
}

func (t *TypeExpression) assertKeys() error {
  if t.parameters != nil {
    if t.parameters[0].key != nil {
      return nil
    } else {
      errCtx := t.parameters[0].typeExpr.Context()
      return errCtx.NewError("Error: expected keyed type parameter")
    }
  } else {
    return nil
  }
}

func (t *TypeExpression) assertNoKeys() error {
  if t.parameters != nil {
    if t.parameters[0].key == nil {
      return nil
    } else {
      errCtx := t.parameters[0].key.Context()
      return errCtx.NewError("Error: unexpected keyed type parameter")
    }
  } else {
    return nil
  }
}

func (t *TypeExpression) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	//b.WriteString("Type(")
	b.WriteString(t.VarExpression.Name())

	if t.hasKeys() {
		b.WriteString("<")
		for _, cont := range t.parameters {
			parameterType := cont.typeExpr

			b.WriteString(cont.key.Value())
			b.WriteString(":")
			b.WriteString(parameterType.Dump(""))
			b.WriteString(",")
		}

		b.WriteString(">")
	} else if t.parameters != nil {
		b.WriteString("<")

		for _, cont := range t.parameters {
			b.WriteString(cont.typeExpr.Dump(""))
		}

		b.WriteString(">")
	}

	b.WriteString("\n")

	return b.String()
}

func (t *TypeExpression) ResolveExpressionNames(scope Scope) error {
	if t.Name() == "any" || t.Name() == "void" {
    if t.parameters != nil {
      errCtx := t.Context()
      return errCtx.NewError("Error: doesn't accept type parameters")
    } else {
      return nil
    }
	}

  if t.Name() != "function" && t.Name() != "class" {
    if err := t.VarExpression.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

  if t.Name() != "Object" {
    if err := t.assertNoKeys(); err != nil {
      return err
    }
  }

	if t.parameters != nil {
		for _, parameter := range t.parameters {
			if err := parameter.typeExpr.ResolveExpressionNames(scope); err != nil {
				return err
			}
		}
	}

	return nil
}

// class<...>
func (t *TypeExpression) generateClass() (values.Value, error) {
  ctx := t.Context()

  // overloads of constructor args
  var cArgs [][]values.Value = nil
  var proto values.Prototype = nil

  if t.parameters != nil {
    if len(t.parameters) < 1 {
      errCtx := ctx
      return nil, errCtx.NewError("Error: 0 class type parameters")
    }

    n := len(t.parameters)

    cArgs = make([][]values.Value, 1)

    cArgs[0] = make([]values.Value, n - 1)

    for i := 0; i < n - 1; i++ {
      parameter := t.parameters[i]

      cArg, err := parameter.typeExpr.EvalExpression()
      if err != nil {
        return nil, err
      }

      if cArg == nil {
        errCtx := parameter.typeExpr.Context()
        return nil, errCtx.NewError("Error: unexpected void value")
      }

      cArgs[0][i] = cArg
    }

    newVal, err := t.parameters[n-1].typeExpr.EvalExpression()
    if err != nil {
      return nil, err
    }

    proto := values.GetPrototype(newVal)
    if proto == nil {
      errCtx := t.parameters[n-1].typeExpr.Context()
      return nil, errCtx.NewError("Error: expected instance of prototype")
    }
  }

  return values.NewClass(cArgs, proto, ctx), nil
}

// function<..., retType>
func (t *TypeExpression) generateFunction() (values.Value, error) {
  ctx := t.Context()

  var args [][]values.Value = nil
  if t.parameters != nil {
    if len(t.parameters) < 1 {
      return nil, ctx.NewError("Error: expected more than 0 function type parameters")
    }

    n := len(t.parameters)

    args = make([][]values.Value, 1)
    args[0] = make([]values.Value, n)

    for i := 0; i < n; i++ {
      p := t.parameters[i]

      arg, err := p.typeExpr.EvalExpression()
      if err != nil {
        return nil, err
      }

      if arg == nil && i < n -1 {
        errCtx := p.typeExpr.Context()
        return nil, errCtx.NewError("Error: unexpected void value")
      } 

      args[0][i] = arg
    }
  }

  return values.NewOverloadedFunction(args, t.Context()), nil
}

func (t *TypeExpression) generateSingleParameterValue() (values.Value, error) {
  var content values.Value
  if t.parameters != nil {
    if len(t.parameters) != 1 {
      ctx := t.Context()
      return nil, ctx.NewError("Error: expected 1 type parameter")
    }

    arg, err := t.parameters[0].typeExpr.EvalExpression()
    if err != nil {
      return nil, err
    }

    if arg == nil {
      errCtx := t.parameters[0].typeExpr.Context()
      return nil, errCtx.NewError("Error: unexpected void value")
    }

    content = arg
  }

  return content, nil
}

func (t *TypeExpression) generateDoubleParameterValue() (values.Value, values.Value, error) {
  content := []values.Value{nil, nil}

  if t.parameters != nil {
    if len(t.parameters) != 2 {
      errCtx := t.Context()
      return nil, nil, errCtx.NewError("Error: expected 1 type parameter")
    }

    for i, p := range t.parameters {
      arg, err := p.typeExpr.EvalExpression()
      if err != nil {
        return nil, nil, err
      }

      if arg == nil {
        errCtx := p.typeExpr.Context()
        return nil, nil, errCtx.NewError("Error: unexpected void value")
      }

      content[i] = arg
    }
  }

  return content[0], content[1], nil
}

func (t *TypeExpression) generateObject() (values.Value, error) {
  var props map[string]values.Value = nil

  if t.parameters != nil {
    if len(t.parameters) < 1 {
      errCtx := t.Context()
      return nil, errCtx.NewError("Error: expected at least 1 type parameter")
    }

    if !t.hasKeys() {
      common, err := t.generateSingleParameterValue()
      if err != nil {
        return nil, err
      }

      return prototypes.NewMapLikeObject(common, t.Context()), nil
    } else {
      for _, p := range t.parameters {
        key := p.key.Value()
        val, err := p.typeExpr.EvalExpression()
        if err != nil {
          return nil, err
        }

        if prev, ok := props[key]; ok {
          errCtx := p.key.Context()
          err := errCtx.NewError("Error: duplicate Object key")
          err.AppendContextString("Info: previously defined here", prev.Context())
          return nil, err
        }

        props[key] = val
      }
    }
  }

  return prototypes.NewObject(props, t.Context()), nil
}

// Acts as the new generate instance
func (t *TypeExpression) EvalExpression() (values.Value, error) {
  ctx := t.Context()

  switch t.Name() {
  case "any":
    if t.parameters != nil {
      panic("should've been checked during resolve stage")
    }
    return values.NewAny(ctx), nil
  case "class": 
    return t.generateClass();
  case "function":
    return t.generateFunction()
  case "void":
    if t.parameters != nil {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: void can't have parameter types")
    }
    return nil, nil
  case "Array":
    content, err := t.generateSingleParameterValue()
    if err != nil {
      return nil, err
    }

    return prototypes.NewArray(content, ctx), nil
  case "Set":
    content, err := t.generateSingleParameterValue()
    if err != nil {
      return nil, err
    }

    return prototypes.NewSet(content, ctx), nil
  case "Map":
    key, val, err := t.generateDoubleParameterValue()
    if err != nil {
      return nil, err
    }

    return prototypes.NewMap(key, val, ctx), nil
  case "Promise":
    content, err := t.generateSingleParameterValue()
    if err != nil {
      return nil, err
    }

    if content == nil {
      if t.parameters != nil {
        return prototypes.NewVoidPromise(ctx), nil
      } else {
        return prototypes.NewPromise(nil, ctx), nil
      }
    } else {
      return  prototypes.NewPromise(content, ctx), nil
    }
  case "Event":
    content, err := t.generateSingleParameterValue()
    if err != nil {
      return nil, err
    }

    return prototypes.NewEvent(content, ctx), nil
  case "IDBRequest":
    content, err := t.generateSingleParameterValue()
    if err != nil {
      return nil, err
    }

    return prototypes.NewIDBRequest(content, ctx), nil
  case "Object":
    return t.generateObject()
  default:
    if t.parameters != nil {
			errCtx := ctx
			return nil, errCtx.NewError("Error: unexpected type parameters")
    }

    interf := t.GetInterface()
    if interf == nil {
      return nil, ctx.NewError("Error: expeceted an interface")
    }

    return values.NewInstance(interf, ctx), nil
	}
}

func (t *TypeExpression) Walk(fn WalkFunc) error {
  if t.parameters != nil {
    for _, cont := range t.parameters {
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

  if err := t.typeExpr.Walk(fn); err != nil {
    return err
  }

  return nil
}
