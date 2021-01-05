package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// simply prints the variable name
type VarExpression struct {
	variable Variable // might be overwritten during ResolveNames
  origName string // used for refactoring
  pkgRef *Package // used for refactoring
	TokenData
}

func newVarExpression(name string, constant bool, ctx context.Context) VarExpression {
	return VarExpression{NewVariable(name, constant, ctx), name, nil, TokenData{ctx}}
}

func NewVarExpression(name string, ctx context.Context) *VarExpression {
  ve := newVarExpression(name, false, ctx)
	return &ve
}

// for Function and Class statements
func NewConstantVarExpression(name string, ctx context.Context) *VarExpression {
  ve := newVarExpression(name, true, ctx)
	return &ve
}

func (t *VarExpression) Name() string {
	if t.variable == nil {
		panic("ref shouldn't be nil")
	}

	return t.variable.Name()
}

func (t *VarExpression) GetVariable() Variable {
	return t.variable
}

func (t *VarExpression) ToTypeExpression() (*TypeExpression, error) {
  return NewTypeExpression(t.Name(), nil, nil, t.Context())
}

func (t *VarExpression) GetInterface() values.Interface {
  obj_ := t.GetVariable().GetObject()
  if obj_ == nil {
    // eg. for random expression
    return nil
  }

  obj, ok := obj_.(values.Interface)
  if ok {
    return obj
  } else {
    return nil
  }
}

func (t *VarExpression) GetPrototype() values.Prototype {
  obj_ := t.GetVariable().GetObject()
  obj, ok := obj_.(values.Prototype)
  if ok {
    return obj
  } else {
    return nil
  }
}

func (t *VarExpression) Dump(indent string) string {
	s := indent + "Var(" + t.Name() + ")\n"
	return s
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

// as rhs
func (t *VarExpression) EvalExpression() (values.Value, error) {
	if t.variable == nil {
		panic("ref is still nil")
	}

  if t.GetInterface() != nil && t.GetPrototype() == nil {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: can't use interface in an expression")
  }

  // note that both Object and Value must be set for builtin classes/interfaces
  res := t.variable.GetValue()

	return values.NewContextValue(res, t.Context()), nil
}

// as lhs
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

func (t *VarExpression) UniversalExpressionNames(ns Namespace) error {
	// nothing to be done
	return nil
}

func (t *VarExpression) UniqueExpressionNames(ns Namespace) error {
	return nil
}

func (t *VarExpression) Walk(fn WalkFunc) error {
  return fn(t)
}

// the following function is used where variables are declared (VarStatement, NodeJSModule require)
func (t *VarExpression) uniqueDeclarationName(ns Namespace, varType VarType) error {
  switch varType {
  case CONST, LET, AUTOLET:
    ns.LetName(t.variable)
  case VAR:
    ns.VarName(t.variable)
  default:
    panic("unexpected")
  }

	return nil
}

func (t *VarExpression) RefersToPackage(absPath string) bool {
  if t.pkgRef != nil {
    return t.pkgRef.Path() == absPath 
  }

  return false
}

func (t *VarExpression) PackagePath() string {
  if t.pkgRef != nil {
    return t.pkgRef.Path()
  } else {
    return ""
  }
}

// use by refactoring tools
func (t *VarExpression) PackageContext() context.Context {
  name := t.origName

  parts := strings.Split(name, ".")

  ctx := t.Context()

  if len(parts) > 1 && !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, ".") {
    subN := len(parts[0]) // so from start to this
    return ctx.NewContext(0, subN)
  } else {
    return ctx
  }
}

// use by refactoring tools
func (t *VarExpression) NonPackageContext() context.Context {
  name := t.origName

  parts := strings.Split(name, ".")

  ctx := t.Context()

  if len(parts) > 1 && !strings.HasPrefix(name, ".") && !strings.HasSuffix(name, ".") {
    startN := len(parts[0]) + 1 // without the dot

    return ctx.NewContext(startN, len(name))
  } else {
    return ctx
  }
}

func IsVarExpression(t Expression) bool {
	_, ok := t.(*VarExpression)
	return ok
}

func AssertVarExpression(t Token) (*VarExpression, error) {
	if ve, ok := t.(*VarExpression); ok {
		return ve, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected variable word")
	}
}

