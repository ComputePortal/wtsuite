package js

type ModuleScope struct {
	module  *ModuleData
	globals GlobalScope // both as parent and as globals
	ScopeData
}

func (ms *ModuleScope) HasVariable(name string) bool {
	if _, ok := ms.module.importedNames[name]; ok {
		return true
	} else {
		return ms.ScopeData.HasVariable(name)
	}
}

func (ms *ModuleScope) GetVariable(name string) (Variable, error) {
	if importedName, ok := ms.module.importedNames[name]; ok {
		if importedName.v != nil {
			return importedName.v, nil
		} else {
			// get the module from which it is imported
			ctx := ms.module.Context()
			importedModule, err := ms.globals.GetModule(ctx.Path(), importedName.dep.Value())
			if err != nil {
				errCtx := importedName.dep.Context()
				return nil, errCtx.NewError("Error: module not found")
			}

			// cache elsewhere
			v, err := importedModule.GetExportedVariable(ms.globals, importedName.old,
				importedName.dep.Context())
			if err != nil {
				return nil, err
			}

			// package can be given the correct name now
			if pkg, ok := v.(*Package); ok {
				pkg.Rename(name)
			}

			return v, nil
		}
	} else {
		return ms.ScopeData.GetVariable(name)
	}
}
