package js

import (
  "strings"

  "./prototypes"
  "./values"

  "../context"
)

type ClassProperty struct {
  name *Word
  typeExpr *TypeExpression // nil if not specified
}

func NewClassProperty(name *Word, typeExpr *TypeExpression) *ClassProperty {
  return &ClassProperty{name, typeExpr}
}

func (p *ClassProperty) Name() string {
  return p.name.Value()
}

func (p *ClassProperty) Role() prototypes.FunctionRole {
  if strings.HasPrefix("_", p.Name()) {
    return prototypes.PRIVATE | prototypes.PROPERTY
  } else {
    return prototypes.PROPERTY
  }
}

func (p *ClassProperty) IsUniversal() bool {
  if p.typeExpr != nil {
    val, err := p.typeExpr.EvalExpression()
    if err != nil {
      panic("should've been checked before")
    }

    if proto := values.GetPrototype(val); proto != nil {
      return proto.IsUniversal()
    } else {
      return false
    }
  } else {
    return false
  }
}

func (p *ClassProperty) Context() context.Context {
  return p.name.Context()
}

func (p *ClassProperty) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(p.Name())
  if p.typeExpr != nil {
    b.WriteString(" ")
    b.WriteString(p.typeExpr.Dump(""))
  }
  b.WriteString("\n")

  return b.String()
}

func (p *ClassProperty) WriteStatement(indent string) string {
  return ""
}

func (p *ClassProperty) ResolveNames(scope Scope) error {
  if p.typeExpr != nil {
    if err := p.typeExpr.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

  return nil
}

func (p *ClassProperty) GetValue(ctx context.Context) (values.Value, error) {
  if p.typeExpr != nil {
    val, err := p.typeExpr.EvalExpression()
    if err != nil {
      return nil, err
    }

    return values.NewContextValue(val, ctx), nil
  } else {
    return values.NewAny(ctx), nil
  }
}

func (p *ClassProperty) SetValue(v values.Value, ctx context.Context) error {
  if p.typeExpr != nil {
    checkVal, err := p.typeExpr.EvalExpression()
    if err != nil {
      return err
    }

    return checkVal.Check(v, ctx)
  } else {
    return nil
  }
}

func (p *ClassProperty) Eval() error {
  if _, err := p.GetValue(p.Context()); err != nil {
    return err
  }

  return nil
} 

func (p *ClassProperty) ResolveActivity(usage Usage) error {
  return nil
}

func (p *ClassProperty) UniversalNames(ns Namespace) error {
  return nil
}

func (p *ClassProperty) UniqueNames(ns Namespace) error {
  // TODO: give members other names when compacting
  return nil
}

func (p *ClassProperty) Walk(fn WalkFunc) error {
  if err := p.name.Walk(fn); err != nil {
    return err
  }

  if p.typeExpr != nil {
    if err := p.typeExpr.Walk(fn); err != nil {
      return err
    }
  }

  return fn(p)
}

