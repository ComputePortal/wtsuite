package js

import (
	"./prototypes"
	"./values"

	"../context"
)

type SpecialNumber struct {
	value string
	TokenData
}

func NewSpecialNumber(value string, ctx context.Context) *SpecialNumber {
	return &SpecialNumber{value, TokenData{ctx}}
}

func (t *SpecialNumber) Value() string {
	return t.value
}

func (t *SpecialNumber) Dump(indent string) string {
	return indent + "SpecialNumber(" + t.WriteExpression() + ")\n"
}

func (t *SpecialNumber) WriteExpression() string {
	return t.value
}

func (t *SpecialNumber) ResolveExpressionNames(scope Scope) error {
	return nil
}

func (t *SpecialNumber) EvalExpression() (values.Value, error) {
	return prototypes.NewNumber(t.Context()), nil
}

func (t *SpecialNumber) ResolveExpressionActivity(usage Usage) error {
	return nil
}

func (t *SpecialNumber) UniversalExpressionNames(ns Namespace) error {
	return nil
}

func (t *SpecialNumber) UniqueExpressionNames(ns Namespace) error {
	return nil
}

func (t *SpecialNumber) Walk(fn WalkFunc) error {
  return fn(t)
}
