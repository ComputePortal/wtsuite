package prototypes

type FunctionRole int

type FunctionWithRole interface {
	Role() FunctionRole
}

const (
	NORMAL   FunctionRole = 0
	CONST                 = 1 << 0
	STATIC                = 1 << 1
	GETTER                = 1 << 2
	SETTER                = 1 << 3
	PUBLIC                = 1 << 4
	PRIVATE               = 1 << 5
	ABSTRACT              = 1 << 6
	OVERRIDE              = 1 << 7
	ASYNC                 = 1 << 8
  PROPERTY              = 1 << 9
)

func IsNormal(m FunctionWithRole) bool {
	return !(IsStatic(m) || IsGetter(m) || IsSetter(m))
}

func IsStatic(m FunctionWithRole) bool {
	return (m.Role() & STATIC) > 0
}

func IsStaticGetter(m FunctionWithRole) bool {
	r := m.Role()

	return ((r & STATIC) > 0) && ((r & GETTER) > 0)
}

func IsGetter(m FunctionWithRole) bool {
	return !IsStatic(m) && ((m.Role() & GETTER) > 0)
}

func IsSetter(m FunctionWithRole) bool {
	return !IsStatic(m) && ((m.Role() & SETTER) > 0)
}

func IsConst(m FunctionWithRole) bool {
	return m.Role()&CONST > 0
}

func IsPrivate(m FunctionWithRole) bool {
	return m.Role()&PRIVATE > 0
}

func IsPublic(m FunctionWithRole) bool {
	return (m.Role()&PUBLIC > 0) || !IsPrivate(m)
}

func IsOverride(m FunctionWithRole) bool {
	return m.Role()&OVERRIDE > 0
}

func IsAbstract(m FunctionWithRole) bool {
	return m.Role()&ABSTRACT > 0
}

func IsAsync(m FunctionWithRole) bool {
	return m.Role()&ASYNC > 0
}

func IsProperty(m FunctionWithRole) bool {
  return m.Role()&PROPERTY > 0
}
