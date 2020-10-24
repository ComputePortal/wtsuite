package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// simply prints the variable name
type VarExpression struct {
	ref Variable // might be overwritten during ResolveNames
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
	if t.ref == nil {
		panic("ref shouldn't be nil")
	}

	return t.ref.Name()
}

func (t *VarExpression) GetVariable() Variable {
	return t.ref
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

	pkgObject := pkg_.GetObject()
	if pkgObject != nil {
		nodejsModule, maybeNodeJSModule := pkgObject.(*prototypes.BuiltinPrototype)
		if maybeNodeJSModule {
			origName := nodejsModule.Name()

			if _, isNodeJSModule := GetNodeJSModule(origName); isNodeJSModule {
				// these tend to only go one level
				if len(parts) != 2 {
					errCtx := t.Context()
					return errCtx.NewError("Error: nodejs module  members are of form module.member")
				}

				// member must be a class
				dummyStack := NewDummyStack()
				value_, err := nodejsModule.GetMember(dummyStack, nil, parts[1], false, t.Context())
				if err != nil {
					return err
				}

				// value_ must be values.Class
				valueClass, ok := value_.GetClassPrototype()
				if !ok {
					errCtx := t.Context()
					return errCtx.NewError("Error: nodejs module member " + parts[1] + " is not a class")
				}

				// use the builtin name as reference
				classVariable, ok := GetNodeJSClassVariable(valueClass.Name())
				if ok {
					t.ref = classVariable
					return nil
				}
			}
		}
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

	t.ref = member
	return nil
}

func (t *VarExpression) ResolveExpressionNames(scope Scope) error {
	name := t.Name()

	// variables that begin with a period might be interal hidden vars
	// (eg. variables created by the Parser.importDefault() macro
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
	t.ref, err = scope.GetVariable(name)
	if err != nil {
		return err
	}

	if t.ref == nil {
		panic("nil variable")
	}

	return nil
}

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

// as rhs
func (t *VarExpression) EvalExpression(stack values.Stack) (values.Value, error) {
	if t.ref == nil {
		panic("ref is still nil")
	}

	if _, isPkg := t.ref.(*Package); isPkg {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: package can't be used as a variable")
	}

	res, err := stack.GetValue(t.ref, t.Context())
	if err != nil {
		return nil, err
	}

	if res == nil {
		errCtx := t.Context()
		err := errCtx.NewError("Error: ref value is nil")
		return nil, err
	}

	return values.NewContextValue(res, t.Context()), nil
}

func (t *VarExpression) ResolveExpressionActivity(usage Usage) error {
	if _, isPkg := t.ref.(*Package); isPkg {
		errCtx := t.Context()
		return errCtx.NewError("Error: package can't be used as a variable")
	}

	return usage.Use(t.ref, t.Context())
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
  case CONST, LET:
    ns.LetName(t.ref)
  case VAR:
    ns.VarName(t.ref)
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
