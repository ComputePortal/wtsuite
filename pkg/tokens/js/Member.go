package js

import (
	"strings"

	"./values"

	"../context"
)

type Member struct {
	object             Expression
	key                *Word
	lhsValue           values.Value
	old                string
	friendlyPrototypes []values.Prototype
	TokenData
}

func NewMember(object Expression, key *Word, ctx context.Context) *Member {
	mergedCtxs := context.MergeContexts(object.Context(), key.Context(), ctx)
	return &Member{
		object,
		key,
		nil,
		key.value,
		[]values.Prototype{},
		TokenData{mergedCtxs},
	}
}

func (t *Member) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Member(")
	b.WriteString(t.key.Value())
	b.WriteString(")\n")

	b.WriteString(t.object.Dump(indent + "  "))

	return b.String()
}

func (t *Member) WriteExpression() string {
	pkgMember, err := t.GetPackageMember()
	if err != nil {
		panic("should've been caught before")
	} else if pkgMember != nil {
    if t.IsBuiltinPackage() {
      return t.PackageName() + "." + pkgMember.Name()
    } else {
      return pkgMember.Name()
    }
	} else {
		var b strings.Builder
		// literals must be surrounded by brackets
		isLit := IsLiteral(t.object)

		if isLit {
			b.WriteString("(")
		}

		b.WriteString(t.object.WriteExpression())

		if isLit {
			b.WriteString(")")
		}

		b.WriteString(".")
		b.WriteString(t.key.Value())

		return b.String()
	}
}

func (t *Member) Args() []Token {
	return []Token{t.object}
}

// return empty name if object is not VarExpression
func (t *Member) ObjectNameAndKey() (string, string) {
	name := ""
	if ve, ok := t.object.(*VarExpression); ok {
		name = ve.Name()
	}

	return name, t.key.Value()
}

func (t *Member) ToTypeExpression() (*TypeExpression, error) {
  base, key := t.ObjectNameAndKey()
  if base == "" {
    return nil, nil
  }

  return NewTypeExpression(base + "." + key, nil, nil, t.Context())
}

func (t *Member) ResolveExpressionNames(scope Scope) error {
	t.friendlyPrototypes = scope.FriendlyPrototypes()

	return t.object.ResolveExpressionNames(scope) // key is done later (requires type info)
}

func (t *Member) havePrivateAccess(objectValue values.Value) bool {
  // is this?:
  if varExpr, ok := t.object.(*VarExpression); ok {
    if varExpr.Name() == "this" {
      return true
    }
  }

  proto := values.GetPrototype(objectValue)
  if proto != nil {
    for _, friendlyProto := range t.friendlyPrototypes {
      if values.PrototypeIsAncestorOf(proto, friendlyProto) {
        return true
      }
    }

    return false
  } else {
    return false
  }
}

func (t *Member) getPackage() (*Package, error) {
	switch obj := t.object.(type) {
	case *VarExpression:
		pkg_ := obj.GetVariable()
		if pkg, ok := pkg_.(*Package); ok {
			// now Member acts as a nested VarExpression
      return pkg, nil
		}
	case *Member:
		pkg_, err := obj.GetPackageMember()
		if err != nil {
			return nil, err
		}

		if pkg, ok := pkg_.(*Package); ok {
			return pkg, nil
		}
	}

  return nil, nil
}

// return nil if this is not a package member
func (t *Member) GetPackageMember() (Variable, error) {
  pkg, err := t.getPackage()
  if err != nil {
    return nil, err
  }

  if pkg != nil {
    return pkg.getMember(t.key.value, t.key.Context())
  }

	return nil, nil
}

// used for refactoring
func (t *Member) PackagePath() string {
  pkg, err := t.getPackage()
  if err != nil || pkg == nil {
    return ""
  } else {
    return pkg.Path()
  }
}

func (t *Member) IsBuiltinPackage() bool {
  pkg, err := t.getPackage()
  if err != nil || pkg == nil {
    return false
  } else {
    return pkg.IsBuiltin()
  }
}

func (t *Member) PackageName() string {
  pkg, err := t.getPackage()
  if err != nil || pkg == nil {
    return ""
  } else {
    return pkg.Name()
  }
}

// used for refactoring
func (t *Member) ObjectContext() context.Context {
  return t.object.Context()
}

// used for refactoring
func (t *Member) KeyContext() context.Context {
  return t.key.Context()
}

func (t *Member) EvalExpression() (values.Value, error) {
	pkgMember, err := t.GetPackageMember()
	if err != nil {
		return nil, err
	} else if pkgMember != nil {
		// use a dummy VarExpression to retrieve the pkgMember value
		tmpVe := NewConstantVarExpression("", t.Context()) // doesn't need a name
		tmpVe.variable = pkgMember

		return tmpVe.EvalExpression()
	}

	objectValue, err := t.object.EvalExpression()
	if err != nil {
		return nil, err
	}

	includePrivate := t.havePrivateAccess(objectValue)

	res, err := objectValue.GetMember(t.key.value, includePrivate, t.key.Context())
	if err != nil {
		return nil, err
	}

	return res, nil
	//return values.NewContextValue(res, t.Context()), nil
}

func (t *Member) EvalSet(rhsValue values.Value, ctx context.Context) error {
	objectValue, err := t.object.EvalExpression()
	if err != nil {
		return err
	}

	includePrivate := t.havePrivateAccess(objectValue)

	err = objectValue.SetMember(t.key.value, includePrivate, rhsValue,
		t.key.Context())

	if err != nil {
		return err
	}

	return nil
}

func (t *Member) ResolveExpressionActivity(usage Usage) error {
	pkgMember, err := t.GetPackageMember()
	if err != nil {
		return err
	} else if pkgMember != nil {
		return usage.Use(pkgMember, t.Context())
	}

	if err := t.object.ResolveExpressionActivity(usage); err != nil {
		return err
	}

	return nil
}

func (t *Member) UniversalExpressionNames(ns Namespace) error {
	return t.object.UniversalExpressionNames(ns)
}

func (t *Member) UniqueExpressionNames(ns Namespace) error {
	return t.object.UniqueExpressionNames(ns)
}

func (t *Member) Walk(fn WalkFunc) error {
  if err := t.key.Walk(fn); err != nil {
    return err
  }

  if err := t.object.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}

func IsAnyMember(t Expression) bool {
	_, ok := t.(*Member)
	return ok
}

func IsMember(t Expression, name string) bool {
	m, ok := t.(*Member)
	if ok {
		return name == m.key.Value()
	} else {
		return false
	}
}
