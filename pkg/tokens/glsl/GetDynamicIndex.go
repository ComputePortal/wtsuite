package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type GetDynamicIndex struct {
  fnVar Variable // variable of the auto-generated function
  Index
}

func NewGetDynamicIndex(args []Expression, ctx context.Context) (Expression, error) {
  if len(args) != 2 {
    return nil, ctx.NewError("Error: expected 2 args")
  }

  return &GetDynamicIndex{nil, newIndex(args[0], args[1], ctx)}, nil
}

func (t *GetDynamicIndex) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("GetDynamicIndex\n")

	b.WriteString(t.container.Dump(indent + "  "))
	b.WriteString(t.index.Dump(indent + "[ "))

	return b.String()
}

func (t *GetDynamicIndex) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.fnVar.Name())
	b.WriteString("(")
	b.WriteString(t.container.WriteExpression())
  b.WriteString(",")
	b.WriteString(t.index.WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (t *GetDynamicIndex) EvalExpression() (values.Value, error) {
	containerValue, err := t.container.EvalExpression()
	if err != nil {
		return nil, err
	}

	indexValue, err := t.index.EvalExpression()
	if err != nil {
		return nil, err
	}

  if _, err := values.AssertInt(indexValue); err != nil {
    return nil, err
  }

  dummyIndex := values.NewLiteralInt(0, t.index.Context())

  return containerValue.GetIndex(dummyIndex, t.Context())
}

func injectDynamicIndexStatement(usage Usage, fnVar Variable, fn Statement, contentVal values.Value) (Variable, error) {
  fnName := fnVar.Name()

  if strings.ContainsAny(fnName, " ()") {
    panic("bad name: " + fnName)
  }

  if values.IsStruct(contentVal) {
    str, err := values.AssertStruct(contentVal)
    if err != nil {
      return nil, err
    }

    strSt, ok := str.GetStructable().(*Struct)
    if !ok {
      panic("expected only glsl.Struct")
    }

    depVar := strSt.GetVariable()

    // if injected before, then reuse
    injectedStatements := usage.GetInjectedStatements(fnName) 
    for _, injectedStatement := range injectedStatements {
      if len(injectedStatement.deps) == 1 && (injectedStatement.deps[0] == depVar) {
        return injectedStatement.variable, nil
      }
    } 

    usage.InjectStatement(fnName, fnVar, []Variable{depVar}, fn)
    return fnVar, nil
  } else {
    // if injected before, then reuse
    injectedStatements := usage.GetInjectedStatements(fnName)
    for _, injectedStatement := range injectedStatements {
      if len(injectedStatement.deps) == 0 {
        return injectedStatement.variable, nil
      }
    } 

    usage.InjectStatement(fnName, fnVar, []Variable{}, fn)
    return fnVar, nil
  }
}

func (t *GetDynamicIndex) ResolveExpressionActivity(usage Usage) error {
  if err := t.Index.ResolveExpressionActivity(usage); err != nil {
    return err
  }

  containerValue, err := t.container.EvalExpression()
  if err != nil {
    return err
  }

  dummyIndex := values.NewLiteralInt(0, t.index.Context())

  contentVal, err := containerValue.GetIndex(dummyIndex, t.Context())
  if err != nil {
    return err
  }

  contentTypeName := contentVal.TypeName()

  length := containerValue.Length()

  fn := NewGetDynamicIndexFunction(containerValue.TypeName(), contentTypeName, length, t.Context())

  t.fnVar, err = injectDynamicIndexStatement(usage, fn.GetVariable(), fn, contentVal)

  return err
}
