package scripts

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

type FileBundleScope struct {
	globals *js.GlobalScopeData
	b       *FileBundle
}

func (bs *FileBundleScope) GetModule(absPath string) (js.Module, error) {
	// TODO: if this is slow -> use a map
	for _, s := range bs.b.scripts {
		if s.Path() == absPath {
			return s.Module(), nil
		}
	}

	// TODO: if module isnt yet included in scripts (e.g. 'dynamic' loading), build it on the fly

	panic("dependency not found")
}

func (bs *FileBundleScope) Parent() js.Scope {
	return bs.globals.Parent()
}

func (bs *FileBundleScope) GetVariable(name string) (js.Variable, error) {
	return bs.globals.GetVariable(name)
}

func (bs *FileBundleScope) HasVariable(name string) bool {
	return bs.globals.HasVariable(name)
}

func (bs *FileBundleScope) SetVariable(name string, v js.Variable) error {
	return bs.globals.SetVariable(name, v)
}

func (bs *FileBundleScope) FriendlyPrototypes() []values.Prototype {
	return bs.globals.FriendlyPrototypes()
}

func (bs *FileBundleScope) GetFunction() *js.Function {
  return nil
}

func (bs *FileBundleScope) IsBreakable() bool {
	return bs.globals.IsBreakable() // false of course
}

func (bs *FileBundleScope) IsContinueable() bool {
	return bs.globals.IsContinueable() // false of course
}

func (bs *FileBundleScope) IsAsync() bool {
	return bs.globals.IsAsync()
}
