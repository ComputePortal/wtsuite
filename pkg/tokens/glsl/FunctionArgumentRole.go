package glsl

import (
  "strings"
)

type FunctionArgumentRole int 

const (
  NO_ROLE FunctionArgumentRole = 0
  IN_ROLE                      = 1 << 0
  OUT_ROLE                     = 1 << 1
)

func RoleToString(role FunctionArgumentRole) string {
  var b strings.Builder

  if role & IN_ROLE > 0 {
    b.WriteString("in ")
  }

  if role & OUT_ROLE > 0 {
    b.WriteString("out ")
  }

  return b.String()
}
