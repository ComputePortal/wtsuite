package directives

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/functions"
  "github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type ScopeData struct {
	parent Scope // nil is used to detect toplevel

	vars    map[string]functions.Var
	templates map[string]Template
  permissive bool
}

func newScopeData(parent Scope) ScopeData {
	return ScopeData{
    parent, 
    make(map[string]functions.Var), 
    make(map[string]Template), 
    false,
  }
}

func (s *ScopeData) Permissive() bool {
  if s.parent != nil {
    return s.parent.Permissive()
  } else {
    return s.permissive
  }
}

// including builtin functions
func (scope *ScopeData) listValidVarNames() string {
	var b strings.Builder

	for k, v := range scope.vars {
		b.WriteString(" \u001b[0m")
		b.WriteString(k)

		// write some type information
		b.WriteString(" ")
		if v.Imported {
			b.WriteString("(imported)")
		}
		b.WriteString("\n")
	}
	if scope.parent != nil {
		b.WriteString(scope.parent.listValidVarNames())
	}

	return b.String()
}

//func (scope *ScopeData) setBlockTarget(block *tokens.Tag, target string) {
  //return scope.parent.setBlockTarget(block, target)
//}

//func (scope *ScopeData) getBlockTarget(block *tokens.Tag) string {
  //return scope.parent.getBlockTarget(block)
//}

func (scope *ScopeData) Parent() Scope {
	return scope.parent
}

func (scope *ScopeData) Sync(dst Scope, keepAutoVars, keepImports, asImports bool, prefix string) error {
	for k, v := range scope.vars {
		if v.Imported && !keepImports {
			continue
		}

		if v.Auto && !keepAutoVars {
			continue
		}

		if asImports {
			v.Imported = true
		}

		dst.SetVar(prefix+k, v)
	}

	for k, c := range scope.templates {
		if c.imported && !keepImports {
			continue
		}

		if asImports {
			c = Template{
				c.name,
				c.extends,
				c.scope,
				c.args,
				c.argDefaults,
				c.superAttr,
				c.children,
				true,
				c.exported,
				c.ctx,
			}
		}

		dst.SetTemplate(prefix+k, c)
	}

	return nil
}

func (scope *ScopeData) SyncPackage(dst Scope, keepAutoVars, keepImports, asImports bool, prefix string) error {
	for k, v := range scope.vars {
		if v.Imported && !keepImports {
			continue
		}

		if v.Auto && !keepAutoVars {
			continue
		}

		if asImports {
			v.Imported = true
		}

		if v.Exported {
			dst.SetVar(prefix+k, v)
		}
	}

	for k, c := range scope.templates {
		if c.imported && !keepImports {
			continue
		}

		if asImports {
			c = Template{
				c.name,
				c.extends,
				c.scope,
				c.args,
				c.argDefaults,
				c.superAttr,
				c.children,
				true,
				c.exported,
				c.ctx,
			}
		}

		if c.exported {
			dst.SetTemplate(prefix+k, c)
		}
	}

	return nil
}

func (scope *ScopeData) SyncFiltered(dst Scope, keepAutoVars, keepImports, asImports bool, prefix string, lst *tokens.List) error {
	found := make([]bool, lst.Len())

	filterImport := func(k string) (tokens.Token, bool, error) {
		b := false
		var entry tokens.Token = nil
		if err := lst.Loop(func(i int, v tokens.Token, last bool) error {
			s, err := tokens.AssertString(v)
			if err != nil {
				return err
			}

			if s.Value() == k {
				if b || found[i] {
					errCtx := v.Context()
					return errCtx.NewError("Error: duplicate import")
				}

				b = true
				found[i] = true
				entry = v
			}

			return nil
		}); err != nil {
			return nil, false, err
		}

		return entry, b, nil
	}

	for k, v := range scope.vars {
		if v.Imported && !keepImports {
			continue
		}

		if v.Auto && !keepAutoVars {
			continue
		}

		if asImports {
			v.Imported = true
		}

		lstEntry, ok, err := filterImport(k)
		if err != nil {
			return err
		}

		if ok {
			if !v.Exported {
				errCtx := lstEntry.Context()
				return errCtx.NewError("Error: var not exported")
			}

			dst.SetVar(prefix+k, v)
		}
	}

	for k, c := range scope.templates {
		if c.imported && !keepImports {
			continue
		}

		if asImports {
			c = Template{
				c.name,
				c.extends,
				c.scope,
				c.args,
				c.argDefaults,
				c.superAttr,
				c.children,
				true,
				c.exported,
				c.ctx,
			}
		}

		lstEntry, ok, err := filterImport(k)
		if err != nil {
			return err
		}

		if ok {
			if !c.exported {
				errCtx := lstEntry.Context()
				return errCtx.NewError("Error: var not exported")
			}

			dst.SetTemplate(prefix+k, c)
		}
	}

	for i, b := range found {
		if !b {
			lstEntry, err := lst.Get(i)
			if err != nil {
				panic(err)
			}
			errCtx := lstEntry.Context()
			return errCtx.NewError("Error: not found")
		}
	}

	return nil
}

func (scope *ScopeData) SetVar(key string, v functions.Var) {
	if key != "_" { // never set dummy vars
		// always set at this level
		scope.vars[key] = v
	}
}

func (scope *ScopeData) SetTemplate(key string, d Template) {
	// always set at this level
	scope.templates[key] = d
}

func (scope *ScopeData) HasVar(key string) bool {
	if _, ok := scope.vars[key]; ok {
		return true
	} else if scope.parent != nil {
		return scope.parent.HasVar(key)
	} else {
		return false
	}
}

func (scope *ScopeData) HasTemplate(key string) bool {
	if _, ok := scope.templates[key]; ok {
		return true
	} else if scope.parent != nil {
		return scope.parent.HasTemplate(key)
	} else {
		return false
	}
}

func (scope *ScopeData) GetVar(key string) functions.Var {
	if v, ok := scope.vars[key]; ok {
		return v
	} else if scope.parent != nil {
		return scope.parent.GetVar(key)
	} else {
		panic("not found")
	}
}

func (scope *ScopeData) GetTemplate(key string) Template {
	if d, ok := scope.templates[key]; ok {
		return d
	} else if scope.parent != nil {
		return scope.parent.GetTemplate(key)
	} else {
		panic("not found")
	}
}

func (scope *ScopeData) Eval(key string, args []tokens.Token,
	ctx context.Context) (tokens.Token, error) {
	return eval(scope, key, args, ctx)
}


func (scope *ScopeData) GetCache() *FileCache {
  return scope.parent.GetCache()
}

func (scope *ScopeData) GetSource() files.Source {
  return scope.parent.GetSource()
}
