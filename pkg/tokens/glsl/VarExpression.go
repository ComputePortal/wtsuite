package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type VarExpression struct {
  variable Variable
  origName string // used for refactoring
  pkgRef *Package // used for refactoring

  TokenData
}

func newVarExpression(name string, ctx context.Context) VarExpression {
  return VarExpression{NewVariable(name, ctx), name, nil, newTokenData(ctx)}
}

func NewVarExpression(name string, ctx context.Context) *VarExpression {
  ve := newVarExpression(name, ctx)

  return &ve
}

func (t *VarExpression) Dump(indent string) string {
  s := indent + "Var(" + t.Name() + ")\n"
  return s
}

func (t *VarExpression) GetVariable() Variable {
  return t.variable
}

func (t *VarExpression) Name() string {
  return t.variable.Name()
}

func (t *VarExpression) WriteExpression() string {
  return t.Name()
}

func (t *VarExpression) resolvePackageMember(scope Scope, parts []string) error {
	if !scope.HasVariable(parts[0]) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: package '" + parts[0] + "' undefined")
		return err
	}

	base := parts[0]
	pkg_, err := scope.GetVariable(base)
	if err != nil {
		panic(err)
	}

	pkg, ok := pkg_.(*Package)
	if !ok {
		errCtx := t.Context()
		return errCtx.NewError("Error: '" + base + "' is not a package")
	}

  t.pkgRef = pkg

	var member Variable = nil
	parts = parts[1:]
	for i, part := range parts {
		member, err = pkg.getMember(part, t.Context())
		if err != nil {
			return err
		}

		if i < len(parts)-1 {
			pkg, ok = member.(*Package)
			if !ok {
				errCtx := t.Context()
				return errCtx.NewError("Error '" + strings.Join(append([]string{base}, parts[:i+1]...), ".") + "' is not a package")
			}
		}
	}

	if _, ok := member.(*Package); ok {
		errCtx := t.Context()
		return errCtx.NewError("Error: can't use package like a variable")
	}

	t.variable = member
	return nil
}

func (t *VarExpression) ResolveExpressionNames(scope Scope) error {
	name := t.Name()

	// variables that begin with a period might be interal hidden vars
	// (eg. variables created by the import() macro
	parts := strings.Split(name, ".")
	if len(parts) > 1 && !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, ".") {
		return t.resolvePackageMember(scope, parts)
	}

	if !scope.HasVariable(name) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: '" + name + "' undefined")
		return err
	}

	var err error
	t.variable, err = scope.GetVariable(name)
	if err != nil {
		return err
	}

	if t.variable == nil {
		panic("nil variable")
	}

	return nil
}

func (t *VarExpression) EvalExpression() (values.Value, error) {
	if t.variable == nil {
		panic("var is nil")
	}

  res := t.variable.GetValue()

  return values.NewContextValue(res, t.Context()), nil
}

func (t *VarExpression) EvalSet(v values.Value, ctx context.Context) error {
  if t.variable.Constant() {
    return ctx.NewError("Error: can't assign to const")
  }

  thisVal := t.variable.GetValue()

  return thisVal.Check(v, ctx)
}

func (t *VarExpression) ResolveExpressionActivity(usage Usage) error {
  return usage.Use(t.GetVariable(), t.Context())
}
