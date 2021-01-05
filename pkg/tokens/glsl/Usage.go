package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Usage interface {
  Use(v Variable, ctx context.Context) error
	Rereference(v Variable, ctx context.Context) error
  IsUsed(v Variable) bool

  GetInjectedStatements(name string) []InjectedStatement
  InjectStatement(name string, variable Variable, deps []Variable, st Statement)
  PopInjectedStatements(dep Variable) []InjectedStatement

  DetectUnused() error
}

type InjectedStatement struct {
  name string
  variable Variable
  deps []Variable
  statement Statement
}

type UsageState struct {
  used bool
  ctx context.Context
}

type UsageData struct {
  vars map[Variable]UsageState
  injected []InjectedStatement
}

func NewUsage() Usage {
  return &UsageData{make(map[Variable]UsageState), make([]InjectedStatement, 0)}
}

func (u *UsageData) Use(v Variable, ctx context.Context) error {
  u.vars[v] = UsageState{true, ctx}

  return nil
}

func (u *UsageData) Rereference(v Variable, ctx context.Context) error {
  if _, ok := u.vars[v]; !ok {
    u.vars[v] = UsageState{false, ctx}
  }

  return nil
}

func (u *UsageData) IsUsed(v Variable) bool {
  if us, ok := u.vars[v]; ok {
    return us.used
  } else {
    return false
  }
}

func (u *UsageData) GetInjectedStatements(name string) []InjectedStatement {
  res := make([]InjectedStatement, 0)

  for _, injSt := range u.injected {
    if injSt.name == name {
      res = append(res, injSt)
    }
  }

  return res
}

func (u *UsageData) InjectStatement(name string, variable Variable, deps []Variable, st Statement) {
  u.injected = append(u.injected, InjectedStatement{
    name,
    variable,
    deps,
    st,
  })
}

func (u *UsageData) PopInjectedStatements(dep Variable) []InjectedStatement {
  rem := make([]InjectedStatement, 0)
  res := make([]InjectedStatement, 0)

  for _, st := range u.injected {
    for i, injDep := range st.deps {
      if dep == injDep {
        if i == 0 {
          st.deps = st.deps[i+1:]
        } else if i < len(st.deps) - 1 {
          st.deps = append(st.deps[0:i], st.deps[i+1:]...)
        } else {
          st.deps = st.deps[0:i]
        }
        break
      }
    }

    if len(st.deps) == 0 {
      res = append(res, st)
    } else {
      rem = append(rem, st)
    }
  }

  u.injected = rem

  return res
}

func (u *UsageData) DetectUnused() error {
	for _, state := range u.vars {
		if !state.used {
			err := state.ctx.NewError("Error: declared but not used")
			return err
		}
	}

  if len(u.injected) > 0 {
    errCtx := u.injected[0].statement.Context()
    return errCtx.NewError("Error: unused injected statements")
  }

	return nil
}
