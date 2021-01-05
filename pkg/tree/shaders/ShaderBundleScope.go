package shaders

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
)

type ShaderBundleScope struct {
	globals *glsl.GlobalScopeData
	b       *ShaderBundle
}

// caller can be taken from a context
func (bs *ShaderBundleScope) GetModule(absPath string) (glsl.Module, error) {
	// TODO: if this is slow -> use a map
	for _, s := range bs.b.shaders {
		if s.Path() == absPath {
			return s.Module(), nil
		}
	}

	// TODO: if module isnt yet included in scripts (e.g. 'dynamic' loading), build it on the fly

	panic("dependency not found")
}

func (bs *ShaderBundleScope) Parent() glsl.Scope {
	return bs.globals.Parent()
}

func (bs *ShaderBundleScope) GetVariable(name string) (glsl.Variable, error) {
	return bs.globals.GetVariable(name)
}

func (bs *ShaderBundleScope) HasVariable(name string) bool {
	return bs.globals.HasVariable(name)
}

func (bs *ShaderBundleScope) SetVariable(name string, v glsl.Variable) error {
	return bs.globals.SetVariable(name, v)
}

func (bs *ShaderBundleScope) GetFunction() *glsl.Function {
  return nil
}
