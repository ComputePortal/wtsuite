package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Member struct {
  object Expression
  key *Word
  old string
  TokenData
}

func NewMember(object Expression, key *Word, ctx context.Context) *Member {
	mergedCtxs := context.MergeContexts(object.Context(), key.Context(), ctx)

	return &Member{
		object,
		key,
		key.value,
		newTokenData(mergedCtxs),
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
    return pkgMember.Name()
	} else {
		var b strings.Builder
		b.WriteString(t.object.WriteExpression())
		b.WriteString(".")
		b.WriteString(t.key.Value())

		return b.String()
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

func (t *Member) GetPackageMember() (Variable, error) {
  pkg, err := t.getPackage() 
  if err != nil {
    return nil, err
  }

  if pkg != nil {
    return pkg.getMember(t.key.Value(), t.key.Context())
  }

  return nil, nil
}

func (t *Member) ResolveExpressionNames(scope Scope) error {
	return t.object.ResolveExpressionNames(scope) // key is done later (requires type info)
}

func (t *Member) EvalExpression() (values.Value, error) {
  pkgMember, err := t.GetPackageMember()
  if err != nil {
    return nil, err
  } else if pkgMember != nil {
		// use a dummy VarExpression to retrieve the pkgMember value
		tmpVe := NewVarExpression("", t.Context()) // doesn't need a name
		tmpVe.variable = pkgMember

		return tmpVe.EvalExpression()
  } else {
    objectValue, err := t.object.EvalExpression()
    if err != nil {
      return nil, err
    }

    res, err := objectValue.GetMember(t.key.Value(), t.key.Context())
    if err != nil {
      return nil, err
    }

    return res, nil
  }
}

func (t *Member) EvalSet(rhsValue values.Value, ctx context.Context) error {
  pkgMember, err := t.GetPackageMember()
  if err != nil {
    return err
  } else if pkgMember != nil {
    errCtx := t.Context()
    return errCtx.NewError("Error: can't set package member")
  }

	objectValue, err := t.object.EvalExpression()
	if err != nil {
		return err
	}

	return objectValue.SetMember(t.key.Value(), rhsValue, t.key.Context())
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
