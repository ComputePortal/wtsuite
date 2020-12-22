package js

import (
	"../context"
)

// intended for LiteralInt, LiteralBool, LiteralFloat, LiteralString and LiteralNull
type LiteralData struct {
	TokenData
}

func newLiteralData(ctx context.Context) LiteralData {
	return LiteralData{TokenData{ctx}}
}

func (t *LiteralData) ResolveExpressionNames(scope Scope) error {
	return nil
}

func (t *LiteralData) ResolveExpressionActivity(usage Usage) error {
	return nil
}

func (t *LiteralData) UniversalExpressionNames(ns Namespace) error {
	return nil
}

func (t *LiteralData) UniqueExpressionNames(ns Namespace) error {
	return nil
}
