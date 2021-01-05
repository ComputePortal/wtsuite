package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Usage interface {
	SetInFunction(bool)
	InFunction() bool
	Rereference(v Variable, ctx context.Context) error
	Use(v Variable, ctx context.Context) error
	DetectUnused() error
}

type UsageState struct {
	used bool
	ctx  context.Context
}

type UsageData struct {
	inFunction bool
	users      map[Variable]UsageState
}

func NewUsage() Usage {
	return &UsageData{false, make(map[Variable]UsageState)}
}

func (u *UsageData) SetInFunction(inFunction bool) {
	u.inFunction = inFunction
}

func (u *UsageData) InFunction() bool {
	return u.inFunction
}

func (u *UsageData) Rereference(v Variable, ctx context.Context) error {
	if u.inFunction {
		if _, ok := u.users[v]; !ok {
			u.users[v] = UsageState{false, ctx}
		}
	}

	return nil
}

func (u *UsageData) Use(v Variable, ctx context.Context) error {
	u.users[v] = UsageState{true, ctx}

	return nil
}

func (u *UsageData) DetectUnused() error {
	for _, state := range u.users {
		if !state.used {
			err := state.ctx.NewError("Error: declared but not used")
			return err
		}
	}

	return nil
}
