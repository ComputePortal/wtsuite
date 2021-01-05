package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Package struct {
  name string 

  path string
  members map[string]Variable

  TokenData
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

func (t *Package) Name() string {
	return t.name
}

func (t *Package) Rename(n string) {
	t.name = n
}

func (t *Package) Constant() bool {
	return true
}

func (t *Package) SetConstant() {
}

func (t *Package) SetValue(v values.Value) {
	panic("not applicable")
}

func (t *Package) GetValue() values.Value {
	panic("not applicable")
}

func (t *Package) SetObject(obj interface{}) {
	panic("not applicable")
}

func (t *Package) GetObject() interface{} {
	panic("not applicable")
}
