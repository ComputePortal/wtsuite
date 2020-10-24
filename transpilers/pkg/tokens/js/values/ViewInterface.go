package values

import (
	"../../context"
)

type ViewInterface interface {
	GetVarTypeInstance(stack Stack, name string, ctx context.Context) (Value, error)
	GetElemTypeInstance(stack Stack, name string, ctx context.Context) (Value, error)
	GetDefTypeInstance(stack Stack, name string, ctx context.Context) (Value, error)

	GetURL() string
	GetHTML(stack Stack, id string, ctx context.Context) (string, error)
	GetElemStates(id string) map[string][]string

	IsElem(key string) bool
}
