package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Continue struct {
	TokenData
}

func NewContinue(ctx context.Context) (*Continue, error) {
	return &Continue{TokenData{ctx}}, nil
}

func (t *Continue) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Continue\n")

	return b.String()
}

func (t *Continue) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("continue;")
	return b.String()
}

func (t *Continue) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Continue) HoistNames(scope Scope) error {
	return nil
}

func (t *Continue) ResolveStatementNames(scope Scope) error {
	if !scope.IsContinueable() {
		errCtx := t.Context()
		return errCtx.NewError("Error: not in continueable scope (i.e. for or while)")
	}

	return nil
}

func (t *Continue) EvalStatement() error {
	return nil
}

func (t *Continue) ResolveStatementActivity(usage Usage) error {
	return nil
}

func (t *Continue) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (t *Continue) UniqueStatementNames(ns Namespace) error {
	return nil
}

func (t *Continue) Walk(fn WalkFunc) error {
  return fn(t)
}
