package js

import (
  "strings"

  "./values"

	"../context"
)

// a js.Package acts like a collection of variables
// technically packages can be nested, but this is not recommended

type Package struct {

	name string

  path string
	members map[string]Variable

	TokenData
}

// for builtin packages
func NewBuiltinPackage(name string) *Package {
  ctx := context.NewDummyContext()

  return &Package{name, "", make(map[string]Variable), TokenData{ctx}}
}

// user packages start nameless
func NewPackage(path string, ctx context.Context) *Package {
	return &Package{"", path, make(map[string]Variable), TokenData{ctx}}
}

func (t *Package) addMember(key string, v Variable) error {
	if other, ok := t.members[key]; ok {
		errCtx := v.Context()
		err := errCtx.NewError("Error: package already contains " + key)
		err.AppendContextString("Info: previously defined here", other.Context())
		return err
	}

	t.members[key] = v

	return nil
}

func (t *Package) getMember(key string, ctx context.Context) (Variable, error) {
	if v, ok := t.members[key]; ok {
		return v, nil
	} else {
		return nil, ctx.NewError("Error: " + t.Name() + "." + key + " undefined")
	}
}

func (t *Package) AddPrototype(proto values.Prototype) {
  memberName := proto.Name()

  if strings.ContainsAny(memberName, "<>") {
    panic("package prototype can't have type parameters")
  }

  ctx := proto.Context()
  variable := NewVariable(memberName, true, ctx)
  variable.SetObject(proto)
  classValue, err := proto.GetClassValue()
  if err != nil {
    panic("should've been caught before")
  }
  variable.SetValue(classValue)

  if err := t.addMember(memberName, variable); err != nil {
    panic(err)
  }
}

func (t *Package) AddValue(memberName string, v values.Value) {
  ctx := v.Context()

  variable := NewVariable(memberName, true, ctx)
  variable.SetValue(v)

  if err := t.addMember(memberName, variable); err != nil {
    panic(err)
  }
}

func (t *Package) Dump(indent string) string {
	return indent + "Package " + t.name
}

func (t *Package) Name() string {
	return t.name
}

func (t *Package) Constant() bool {
	return true
}

func (t *Package) SetConstant() {
}

func (t *Package) Rename(newName string) {
	t.name = newName
}

func (t *Package) SetObject(ptr interface{}) {
	panic("not applicable")
}

func (t *Package) GetObject() interface{} {
	return nil
}

func (t *Package) SetValue(v values.Value) {
	panic("not applicable")
}

func (t *Package) GetValue() values.Value {
	panic("not applicable")
}

func (t *Package) Path() string {
  // used for renaming packages
  return t.path
}
