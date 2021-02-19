package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Module interface {
	GetExportedVariable(gs GlobalScope, name string,
		nameCtx context.Context) (Variable, error)

  SymbolDependencies(allModules map[string]Module, name string) []string // non-unique, not sorted
  MinimalDependencies(allModules map[string]Module) []string // non-unique, not-sorted

	Context() context.Context
}

type ImportedVariable struct {
	old string
  new string
	dep *LiteralString // abs path to a file (not a directory, nor pkg/module path)
  lang files.Lang // we need to remember this so we can tell tree/scripts/NewFileScript exactly what language to transpile
	v   Variable       // cache it, so we don't need to keep searching for it
  // can also be used during refactoring

	ctx context.Context
}

type ExportedVariable struct {
	inner string
	v     Variable
	ctx   context.Context
}

type ModuleData struct {
	importedNames    map[string]*ImportedVariable
	exportedNames    map[string]*ExportedVariable
	aggregateExports map[string]*ImportedVariable // "*<path>" is used as key for unamed aggregate exports
	Block
}

type NotFoundMarker struct {
}

func NewNotFoundMarker() *NotFoundMarker {
  return &NotFoundMarker{}
}

func IsNotFoundError(err_ error) bool {
  if err, ok := err_.(*context.ContextError); ok {
    obj := err.GetObject()
    if obj != nil {
      if _, okInner := obj.(*NotFoundMarker); okInner {
        return true
      }
    }
  }

  return false
}

func NewModule(ctx context.Context) *ModuleData {
	// statements are added later
	return &ModuleData{
		make(map[string]*ImportedVariable),
		make(map[string]*ExportedVariable),
		make(map[string]*ImportedVariable),
		newBlock(ctx),
	}
}

func (m *ModuleData) newScope(globals GlobalScope) Scope {
	return &ModuleScope{m, globals, newScopeData(globals)}
}

// called from within other module
func (m *ModuleData) GetExportedVariable(gs GlobalScope, name string,
	nameCtx context.Context) (Variable, error) {
	if name == "*" {
    // prepare self as a package
    ctx := m.Context()
		pkg := NewPackage(ctx.Path(), nameCtx)
		for name, _ := range m.exportedNames {
			v, err := m.GetExportedVariable(gs, name, nameCtx)
			if err != nil {
				return nil, err
			}

			if pkgErr := pkg.addMember(name, v); pkgErr != nil {
				return nil, err
			}
		}

		// also add the aggregate exports
		for name, _ := range m.aggregateExports {
			v, err := m.GetExportedVariable(gs, name, nameCtx)
			if err != nil {
				return nil, err
			}

			if pkgErr := pkg.addMember(name, v); pkgErr != nil {
				return nil, err
			}
		}

		return pkg, nil
	} else if exportedName, ok := m.exportedNames[name]; ok {
		if exportedName.v == nil {
			panic("export should've been set by the ResolveNames stage")
		}

		return exportedName.v, nil
	} else if aggregateExport, ok := m.aggregateExports[name]; ok {
		if aggregateExport.v != nil {
			return aggregateExport.v, nil
		} else {
			importedModule, err := gs.GetModule(aggregateExport.dep.Value())
			if err != nil {
				errCtx := aggregateExport.dep.Context()
				return nil, errCtx.NewError("Error: module not found")
			}

			v, err := importedModule.GetExportedVariable(gs, aggregateExport.old,
				aggregateExport.dep.Context())
			if err != nil {
				return nil, err
			}

			aggregateExport.v = v
			m.aggregateExports[name] = aggregateExport
			return v, nil
		}
	} else {
    // look in any unamed aggregate exports
    var found Variable = nil
    var foundDep *LiteralString = nil
    for key, aggregateExport := range m.aggregateExports {
      if strings.HasPrefix(key, "*") {
        depModule, err := gs.GetModule(aggregateExport.dep.Value())
        if err != nil {
          errCtx := aggregateExport.dep.Context()
          return nil, errCtx.NewError("Error: module not found")
        }

        v, err := depModule.GetExportedVariable(gs, name, aggregateExport.dep.Context())
        if err == nil {
          if found != nil {
            errCtx := aggregateExport.dep.Context()
            err := errCtx.NewError("Error: " + name + " exported twice")
            err.AppendContextString("Info: also exported by this module", foundDep.Context())
          }

          found = v
          foundDep = aggregateExport.dep
        } else if IsNotFoundError(err) {
          continue
        } else {
          return nil, err
        }
      }
    }

    if found == nil {
      notFoundError := nameCtx.NewError("Error: '" + name + "' not exported by this module")
      notFoundError.SetObject(NewNotFoundMarker())
      return nil, notFoundError
    } else {
      // save a copy of the export as an aggregate export in this module too
      aggregateExport := newImportedVariable(name, name, foundDep, files.SCRIPT, nameCtx)
      aggregateExport.v = found

      m.aggregateExports[name] = aggregateExport
      return found, nil
    }
	}
}

func (m *ModuleData) Dump() string {
	var b strings.Builder

	if len(m.importedNames) > 0 {
		b.WriteString("#Module imported names:\n")
		for k, v := range m.importedNames {
			b.WriteString("#  ")
			b.WriteString(v.old)
			b.WriteString(" as \u001b[1m")
			b.WriteString(k)
			b.WriteString("\u001b[0m from '")
			b.WriteString(v.dep.Value())
			b.WriteString("'\n")
		}
	}

	if len(m.exportedNames) > 0 {
		b.WriteString("#Module exported names:\n")
		for k, v := range m.exportedNames {
			b.WriteString("#  ")
			b.WriteString(v.inner)
			b.WriteString(" as \u001b[1m")
			b.WriteString(k)
			b.WriteString("\u001b[0m\n")
		}
	}

	for _, s := range m.statements {
		b.WriteString(s.Dump(""))
	}

	return b.String()
}

func (m *ModuleData) Parent() Scope {
	return nil
}

// import statements must be toplevel, so we could instead loop the statements
func (m *ModuleData) Dependencies() []files.PathLang {
	result := make([]files.PathLang, 0)
	done := make(map[string]bool)

  fn := func(iv *ImportedVariable) {
    pathVal := iv.dep.Value()
		if _, ok := done[pathVal]; !ok {
			result = append(result, files.PathLang{pathVal, iv.lang, iv.dep.Context()})
			done[pathVal] = true
		}
  }

  for _, iv := range m.importedNames {
    fn(iv)
  }

	for _, iv := range m.aggregateExports {
    fn(iv)
	}

	return result
}

func (m *ModuleData) MinimalDependencies(allModules map[string]Module) []string {
  res := make([]string, 0)

  // in this case ignore the aggregate exports, and just look at the imports
  for _, iv := range m.importedNames {
    oldName := iv.old
    // oldName can also be "*" or ""

    depModule, ok := allModules[iv.dep.Value()]
    if !ok {
      panic("all modules should be available at this point")
    }

    res = append(res, depModule.SymbolDependencies(allModules, oldName)...)
  }

  return res
}

func (m *ModuleData) SymbolDependencies(allModules map[string]Module, name string) []string {
  thisCtx := m.Context()
  thisPath := thisCtx.Path()

  if name == "" || name == "*" { // include all Dependencies, and self too
    res_ := m.Dependencies()

    res := []string{thisPath}

    for _, r := range res_ {
      res = append(res, r.Path)
    }

    return res
  } else {
    if ae, ok := m.aggregateExports[name]; ok {
      newName := ae.new

      depModule, ok := allModules[ae.dep.Value()]
      if !ok {
        panic("all modules should be available")
      }

      return depModule.SymbolDependencies(allModules, newName)
    } else if _, ok := m.exportedNames[name]; ok {
      // this file is definitely needed, and thus also all imports
      res := []string{thisPath}

      for _, iv := range m.importedNames {
        depModule, ok := allModules[iv.dep.Value()]
        if !ok {
          panic("all modules should be available")
        }

        res = append(res, depModule.SymbolDependencies(allModules, iv.old)...)
      }

      return res
    } else {
      // look in any of the unamed aggregate exports
      // no errors are thrown here, if not found simply no dependencies are added
      res := make([]string, 0)

      for key, ae := range m.aggregateExports {
        if strings.HasPrefix(key, "*") {
          depModule, ok := allModules[ae.dep.Value()]
          if !ok {
            panic("all modules should be available")
          }
          
          res = append(res, depModule.SymbolDependencies(allModules, name)...)
        }
      }

      return res
    }
  }
}

func (m *ModuleData) Write(usage Usage, nl string, tab string) (string, error) {
	var b strings.Builder

	// TODO: write standard library imports

	b.WriteString(m.writeBlockStatements(usage, "", nl, tab))

	if b.Len() != 0 {
		b.WriteString(";")
		b.WriteString(nl)
	}

	return b.String(), nil
}

func newImportedVariable(oldName, newName string, pathLiteral *LiteralString, lang files.Lang, ctx context.Context) *ImportedVariable {
  return &ImportedVariable{oldName, newName, pathLiteral, lang, nil, ctx}
}

func NewImportedVariable(oldName, newName string, pathLiteral *LiteralString, lang files.Lang, ctx context.Context) (*ImportedVariable, error) {
  path := pathLiteral.Value()
  // make relative paths absolute
  absPath := path
  var err error
  if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
    // language doesnt matter, these should be files
    absPath, err = files.Search(ctx.Path(), path)
  } else {
    switch lang {
    case files.SCRIPT:
      absPath, err = files.SearchScript(ctx.Path(), path)
    case files.TEMPLATE:
      absPath, err = files.SearchTemplate(ctx.Path(), path) 
    default:
      err = ctx.NewError("Error: unimportable language")
    }
  }

  if err != nil {
    errCtx := pathLiteral.Context()
    return nil, errCtx.NewError(err.Error())
  }

  pathLiteral = NewLiteralString(absPath, pathLiteral.Context())

  return newImportedVariable(oldName, newName, pathLiteral, lang, ctx), nil
}

func (m *ModuleData) AddImportedName(newName, oldName string, pathLiteral *LiteralString, lang files.Lang, ctx context.Context) error {
	if newName != "" {
		if oldName == "" {
			panic("invalid oldname")
		}

		if other, ok := m.importedNames[newName]; ok {
			err := ctx.NewError("Error: imported variable already imported")
			err.AppendContextString("Info: imported here", other.ctx)
			return err
		}

    iv, err := NewImportedVariable(oldName, newName, pathLiteral, lang, ctx)
    if err != nil {
      return err
    }

		m.importedNames[newName] = iv
	}

	return nil
}

func (m *ModuleData) AddExportedName(outerName, innerName string, v Variable, ctx context.Context) error {
  // exportName might've been added during parsing (to be able to do initial tree shaking), and again during resolve names to have the exact variable
	if other, ok := m.exportedNames[outerName]; ok && !ctx.Same(&other.ctx) {
		err := ctx.NewError("Error: exported variable name already used")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

	if other, ok := m.aggregateExports[outerName]; ok {
		err := ctx.NewError("Error: name already exported as aggregate")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

	m.exportedNames[outerName] = &ExportedVariable{innerName, v, ctx}

	return nil
}

func (m *ModuleData) AddAggregateExport(newName, oldName string, pathLiteral *LiteralString, lang files.Lang, ctx context.Context) error {
	if newName == "" || oldName == "" {
		panic("bad names")
	}

	if other, ok := m.exportedNames[newName]; ok {
		err := ctx.NewError("Error: name already exported")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

	if other, ok := m.aggregateExports[newName]; ok {
		err := ctx.NewError("Error: name already exported as aggregate")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

  iv, err := NewImportedVariable(oldName, newName, pathLiteral, lang, ctx)
  if err != nil {
    return err
  }

	m.aggregateExports[newName] = iv

	return nil
}

func (m *ModuleData) ResolveNames(gs GlobalScope) error {
	// wrap GlobalScope in a ModuleScope, so that we can add variables
	ms := m.newScope(gs)

	// cache all imports
	/*for name, imported := range m.importedNames {
		v, err := ms.GetVariable(name)
		if err != nil {
			return err
		}

		imported.v = v
		m.importedNames[name] = imported
	}*/

	for name, ae := range m.aggregateExports {
    if !strings.HasPrefix(name, "*") {
      if _, err := m.GetExportedVariable(gs, name, ae.ctx); err != nil {
        return err
      }
    } // unamed aggregate export are handled on name by name basis, as needed by entry points
	}

	return m.Block.HoistAndResolveStatementNames(ms)
}

func (m *ModuleData) EvalTypes() error {
	if err := m.Block.EvalStatement(); err != nil {
		return err
	}

	return nil
}

func (m *ModuleData) ResolveActivity(usage Usage) error {
	return m.Block.ResolveStatementActivity(usage)
}

func (m *ModuleData) UniversalNames(ns Namespace) error {
	return m.Block.UniversalStatementNames(ns)
}

func (m *ModuleData) UniqueNames(ns Namespace) error {
	return m.Block.UniqueStatementNames(ns)
}

func (m *ModuleData) Walk(fn WalkFunc) error {
  if err := m.Block.Walk(fn); err != nil {
    return err
  }

  for _, iv := range m.importedNames {
    if err := iv.Walk(fn); err != nil {
      return err
    }
  }

  for _, iv := range m.aggregateExports {
    if err := iv.Walk(fn); err != nil {
      return err
    }
  }

  return fn(m)
}

func (iv *ImportedVariable) Walk(fn WalkFunc) error {
  return fn(iv)
}

func (iv *ImportedVariable) AbsPath() string {
  return iv.dep.Value()
}

func (iv *ImportedVariable) GetVariable() Variable {
  return iv.v
}

func (iv *ImportedVariable) PathLiteral() *LiteralString {
  return iv.dep
}

func (iv *ImportedVariable) PathContext() context.Context {
  ctx := iv.dep.Context()

  // remove the quotes
  return ctx.NewContext(1, len(ctx.Content())-1)
}

func (iv *ImportedVariable) Context() context.Context {
  return iv.ctx
}

func (m *ModuleData) UniqueEntryPointNames(ns Namespace) error {
	for newName, ae := range m.aggregateExports {
		if err := ns.LibName(ae.v, newName); err != nil {
			return err
		}
	}

	for newName, ex := range m.exportedNames {
		if err := ns.LibName(ex.v, newName); err != nil {
			return err
		}
	}

	return nil
}
