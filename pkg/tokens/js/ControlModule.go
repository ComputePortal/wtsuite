// TODO "ControlModule" should be probably be renamed to just "Module"
package js

import (
	"strings"

	"../context"

	"../../files"
)

type ImportedVariable struct {
	old string
  new string
	dep *LiteralString // path
	v   Variable       // cache it, so we don't need to keep searching for it
  // can also be used during refactoring

	ctx context.Context
}

type ExportedVariable struct {
	inner string
	v     Variable
	ctx   context.Context
}

type ControlModule struct {
	dependencies     []*LiteralString // relative paths!
	importedNames    map[string]*ImportedVariable
	exportedNames    map[string]*ExportedVariable
	aggregateExports map[string]*ImportedVariable
	Block
}

func NewControlModule(ctx context.Context) *ControlModule {
	// statements are added later
	return &ControlModule{
		make([]*LiteralString, 0),
		make(map[string]*ImportedVariable),
		make(map[string]*ExportedVariable),
		make(map[string]*ImportedVariable),
		newBlock(ctx),
	}
}

func (m *ControlModule) newScope(globals GlobalScope) Scope {
	return &ModuleScope{m, globals, newScopeData(globals)}
}

// called from within other module
func (m *ControlModule) GetExportedVariable(gs GlobalScope, name string,
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
			ctx := m.Context()
			importedModule, err := gs.GetModule(ctx.Path(), aggregateExport.dep.Value())
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
		return nil, nameCtx.NewError("Error: '" + name + "' not exported by this module")
	}
}

func (m *ControlModule) Dump() string {
	var b strings.Builder

	if len(m.dependencies) > 0 {
		b.WriteString("#Module dependencies:\n")
		for _, d := range m.dependencies {
			b.WriteString("#  '")
			b.WriteString(d.Value())
			b.WriteString("'\n")
		}
	}

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

func (m *ControlModule) Parent() Scope {
	return nil
}

func (m *ControlModule) Dependencies() []string {
	result := make([]string, 0)
	done := make(map[string]bool)

	for _, dep := range m.dependencies {
		if _, ok := done[dep.Value()]; !ok {
			result = append(result, dep.Value())
			done[dep.Value()] = true
		}
	}

	return result
}

func (m *ControlModule) Write() (string, error) {
	var b strings.Builder

	// TODO: write standard library imports

	b.WriteString(m.writeBlockStatements("", NL))

	if b.Len() != 0 {
		b.WriteString(";")
		b.WriteString(NL)
	}

	return b.String(), nil
}

func (m *ControlModule) addDependency(dep *LiteralString) error {
	ctx := m.Context()
	if _, err := files.Search(ctx.Path(), dep.Value()); err != nil {
		errCtx := dep.Context()
		return errCtx.NewError("Error: file not found")
	}

	for _, other := range m.dependencies {
		if other.Value() == dep.Value() {
			return nil
		}
	}

	m.dependencies = append(m.dependencies, dep)

	return nil
}

func (m *ControlModule) AddImportedName(newName, oldName string, pathLiteral *LiteralString, ctx context.Context) error {
	if err := m.addDependency(pathLiteral); err != nil {
		return err
	}

	if newName != "" {
		if oldName == "" {
			panic("invalid oldname")
		}

		if other, ok := m.importedNames[newName]; ok {
			err := ctx.NewError("Error: imported variable already imported")
			err.AppendContextString("Info: imported here", other.ctx)
			return err
		}

		m.importedNames[newName] = &ImportedVariable{oldName, newName, pathLiteral, nil, ctx}
	}

	return nil
}

func (m *ControlModule) AddExportedName(outerName, innerName string, v Variable, ctx context.Context) error {
	if other, ok := m.exportedNames[outerName]; ok {
		err := ctx.NewError("Error: exported variable name already used")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

	m.exportedNames[outerName] = &ExportedVariable{innerName, v, ctx}

	return nil
}

func (m *ControlModule) AddAggregateExport(newName, oldName string, pathLiteral *LiteralString, ctx context.Context) error {
	if err := m.addDependency(pathLiteral); err != nil {
		return nil
	}

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

	m.aggregateExports[newName] = &ImportedVariable{oldName, newName, pathLiteral, nil, ctx}

	return nil
}

func (m *ControlModule) ResolveNames(gs GlobalScope) error {
	// wrap GlobalScope in a ModuleScope, so that we can add variables
	ms := m.newScope(gs)

	// cache all imports
	for name, imported := range m.importedNames {
		v, err := ms.GetVariable(name)
		if err != nil {
			return err
		}

		imported.v = v
		m.importedNames[name] = imported
	}

	// get all aggregate exports, so they are ready in case of EvalAsEntryPoint
	for name, ae := range m.aggregateExports {
		if _, err := m.GetExportedVariable(gs, name, ae.ctx); err != nil {
			return err
		}
	}

	return m.Block.HoistAndResolveStatementNames(ms)
}

func (m *ControlModule) EvalTypes() error {
	if err := m.Block.EvalStatement(); err != nil {
		return err
	}

	return nil
}

func (m *ControlModule) ResolveActivity(usage Usage) error {
	return m.Block.ResolveStatementActivity(usage)
}

func (m *ControlModule) UniversalNames(ns Namespace) error {
	return m.Block.UniversalStatementNames(ns)
}

func (m *ControlModule) UniqueNames(ns Namespace) error {
	return m.Block.UniqueStatementNames(ns)
}

func (m *ControlModule) Walk(fn WalkFunc) error {
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

func (m *ControlModule) UniqueEntryPointNames(ns Namespace) error {
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
