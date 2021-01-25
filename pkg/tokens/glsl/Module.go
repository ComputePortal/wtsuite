package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Module interface {
	GetExportedVariable(gs GlobalScope, name string, nameCtx context.Context) (Variable, error)

	Context() context.Context
}

type ExportedVariable struct {
	v     Variable
	ctx   context.Context
}

type ModuleData struct {
  exported map[string]*ExportedVariable // key is exported name

  Block
}

func NewModule(ctx context.Context) *ModuleData {
  return &ModuleData{
    make(map[string]*ExportedVariable),
    newBlock(ctx),
  }
}

func (m *ModuleData) newScope(globals GlobalScope) Scope {
	return &ModuleScope{m, globals, newScopeData(globals)}
}

func (m *ModuleData) Dependencies() []files.PathLang {
  result := make([]files.PathLang, 0)
	done := make(map[string]bool)

  fn := func(path *LiteralString) {
    if _, ok := done[path.Value()]; !ok {
      result = append(result, files.PathLang{path.Value(), files.SHADER, path.Context()})
      done[path.Value()] = true
    }
  }
  
  for _, st_ := range m.statements {
    switch st := st_.(type) {
    case *Import:
      fn(st.Path())
    case *ImportExport:
      fn(st.Path())
    }
  }

	return result
}

func (m *ModuleData) SetExportedVariable(outerName string, v Variable, ctx context.Context) error {
	if other, ok := m.exported[outerName]; ok {
		err := ctx.NewError("Error: name already exported as aggregate")
		err.AppendContextString("Info: exported here", other.ctx)
		return err
	}

	m.exported[outerName] = &ExportedVariable{v, ctx}

	return nil
}

// called from within other module
func (m *ModuleData) GetExportedVariable(gs GlobalScope, name string,
	nameCtx context.Context) (Variable, error) {
	if name == "*" {
    // prepare self as a package
    ctx := m.Context()
		pkg := NewPackage(ctx.Path(), nameCtx)
		for name, _ := range m.exported {
      // call self recursively
			v, err := m.GetExportedVariable(gs, name, nameCtx)
			if err != nil {
				return nil, err
			}

			if pkgErr := pkg.addMember(name, v); pkgErr != nil {
				return nil, err
			}
		}

		return pkg, nil
	} else if v, ok := m.exported[name]; ok {
		return v.v, nil
	} else {
		return nil, nameCtx.NewError("Error: '" + name + "' not exported by this module")
	}
}

func (m *ModuleData) ResolveStatementNames(scope Scope) error {
  panic("call Block.ResolveStatementNames instead")
}

func (m *ModuleData) ResolveNames(gs GlobalScope) error {
	// wrap GlobalScope in a ModuleScope, so that we can add variables
	ms := m.newScope(gs)

  return m.Block.ResolveStatementNames(ms)
}

func (m *ModuleData) ResolveEntryNames(gs GlobalScope) (Variable, error) {
	// wrap GlobalScope in a ModuleScope, so that we can add variables
	ms := m.newScope(gs)

  if err := m.Block.ResolveStatementNames(ms); err != nil {
    return nil, err
  }

  if !ms.HasVariable("main") {
    errCtx := m.Context()
    return nil, errCtx.NewError("Error: no main function found")
  }

  return ms.GetVariable("main")
}

func (m *ModuleData) Write(usage Usage, nl string, tab string) (string, error) {
  var b strings.Builder

	for _, st := range m.statements {
		s := st.WriteStatement(usage, "", nl, tab)

    if strings.HasPrefix(s, "#") { // is preproc directive
      b.WriteString("\n")
      b.WriteString(s)
      b.WriteString("\n")
    } else if s != "" {
			b.WriteString(s)
      b.WriteString("\n")
		}
	}

	return b.String(), nil
}

func (m *ModuleData) EvalTypes() error {
  return m.Block.evalStatements()
}

func (m *ModuleData) ResolveActivity(usage Usage) error {
  return m.Block.ResolveStatementActivity(usage)
}

func (m *ModuleData) FinalizeInjected(usage Usage) error {
  injectedStatements := usage.PopInjectedStatements(nil)

  statements := make([]Statement, 0)

  for _, injectedStatement := range injectedStatements {
    statements = append(statements, injectedStatement.statement)
  }

  for _, st_ := range m.statements {
    statements = append(statements, st_)

    if st, ok := st_.(*Struct); ok {
      stVar := st.GetVariable()

      injectedStatements := usage.PopInjectedStatements(stVar)

      for _, injectedStatement := range injectedStatements {
        statements = append(statements, injectedStatement.statement)
      }
    }
  }

  m.statements = statements

  return nil
}

func (m *ModuleData) UniqueNames(ns Namespace) error {
  return m.Block.UniqueStatementNames(ns)
}

func (m *ModuleData) CollectVersion(version *Word) (*Word, error) {
  for i, st_ := range m.statements {
    if st, ok := st_.(*Version); ok {
      if i != 0 {
        panic("should always be first statement")
      }

      return st.CollectVersion(version)
    }
  }

  return version, nil
}

func (m *ModuleData) CollectVaryings(varyings map[string]string) error {
  for _, st_ := range m.statements {
    if st, ok := st_.(*Varying); ok {
      if err := st.Collect(varyings); err != nil {
        return err
      }
    }
  }

  return nil
}

func (m *ModuleData) FindExportedConst(name string) *Const {
  if exported, ok := m.exported[name]; ok {
    variable := exported.v

    if obj_ := variable.GetObject(); obj_ != nil {
      if obj, ok := obj_.(*Const); ok {
        return obj
      }
    }
  }

  return nil
}
