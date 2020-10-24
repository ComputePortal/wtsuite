package js

import (
	"strings"

	"./values"

	"../context"
)

type Break struct {
	TokenData
}

func NewBreak(ctx context.Context) (*Break, error) {
	return &Break{TokenData{ctx}}, nil
}

func (t *Break) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Break\n")

	return b.String()
}

func (t *Break) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("break;")
	return b.String()
}

func (t *Break) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Break) HoistNames(scope Scope) error {
	return nil
}

func (t *Break) ResolveStatementNames(scope Scope) error {
	if !scope.IsBreakable() {
		errCtx := t.Context()
		return errCtx.NewError("Error: break not in breakable scope (i.e. switch, while or for)")
	}
	return nil
}

func (t *Break) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Break) EvalStatement(stack values.Stack) error {
	return nil
}

func (t *Break) ResolveStatementActivity(usage Usage) error {
	return nil
}

func (t *Break) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (t *Break) UniqueStatementNames(ns Namespace) error {
	return nil
}

func (t *Break) Walk(fn WalkFunc) error {
  return fn(t)
}
