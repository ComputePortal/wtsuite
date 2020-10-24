package scripts

import (
	"../../files"
	"../../tokens/js"
	"../../tokens/js/values"
)

type FileBundleScope struct {
	globals *js.GlobalScopeData
	b       *FileBundle
}

// caller can be taken from a context
func (bs *FileBundleScope) GetModule(caller string, path string) (js.Module, error) {
	absPath, err := files.Search(caller, path)
	if err != nil {
		return nil, err
	}

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

func (bs *FileBundleScope) IsBreakable() bool {
	return bs.globals.IsBreakable() // false of course
}

func (bs *FileBundleScope) IsContinueable() bool {
	return bs.globals.IsContinueable() // false of course
}

func (bs *FileBundleScope) IsAsync() bool {
	return bs.globals.IsAsync()
}
