package directives

import (
  "fmt"
  "reflect"
	"strings"

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

    if err := dst.SetVar(prefix+k, v); err != nil {
      return err
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
        c.final,
				c.ctx,
			}
		}

    if err := dst.SetTemplate(prefix+k, c); err != nil {
      return err
    }
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

    isExported := v.Exported
		if asImports {
			v.Imported = true
      v.Exported = false
		}

		if isExported {
      if err := dst.SetVar(prefix+k, v); err != nil {
        return err
      }
		}
	}

	for k, c := range scope.templates {
		if c.imported && !keepImports {
			continue
		}


		if c.exported {
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
          false,
          c.final,
          c.ctx,
        }
      }

      if err := dst.SetTemplate(prefix+k, c); err != nil {
        return err
      }
		}
	}

	return nil
}

func (scope *ScopeData) SyncFiltered(dst Scope, keepAutoVars, keepImports, asImports bool, prefix string, names *tokens.StringDict) error {
	found := make(map[string]bool, names.Len())

  // returned token is just for context
	filterImport := func(k string) (tokens.Token, bool, error) {
		b := false
		var entry tokens.Token = nil
		if err := names.Loop(func(oldName *tokens.String, newName tokens.Token, last bool) error {
			if oldName.Value() == k {
				if b || found[oldName.Value()] {
					errCtx := oldName.Context()
					return errCtx.NewError("Error: duplicate import")
				}

				b = true
				found[oldName.Value()] = true
				entry = oldName
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

    isExported := v.Exported
		if asImports {
			v.Imported = true
      v.Exported = false
		}

		ctxToken, ok, err := filterImport(k)
		if err != nil {
			return err
		}

		if ok {
			if !isExported {
				errCtx := ctxToken.Context()
        return errCtx.NewError("Error: var \"" + k + "\" not exported")
			}

      newName, err := tokens.DictString(names, k)
      if err != nil {
        return err
      }

      if err := dst.SetVar(prefix+newName.Value(), v); err != nil {
        return err
      }
		}
	}

	for k, c := range scope.templates {
		if c.imported && !keepImports {
			continue
		}

    isExported := c.exported

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
				false,
        c.final,
				c.ctx,
			}
		}

    // namesEntry only for context
		ctxToken, ok, err := filterImport(k)
		if err != nil {
			return err
		}

		if ok {
			if !isExported {
				errCtx := ctxToken.Context()
				return errCtx.NewError("Error: template \"" + k + "\" not exported")
			}

      newName, err := tokens.DictString(names, k)
      if err != nil {
        return err
      }

      if err := dst.SetTemplate(prefix+newName.Value(), c); err != nil {
        return err
      }
		}
	}

  if err := names.Loop(func(oldName *tokens.String, _ tokens.Token, last bool) error {
    if b, ok := found[oldName.Value()]; !ok || !b {
			errCtx := oldName.Context()
			return errCtx.NewError("Error: \"" + oldName.Value() + "\" not found")
    }
    return nil
  }); err != nil {
    return err
  }

	return nil
}

func (scope *ScopeData) SetVar(key string, v functions.Var) error {
  if v.Exported && scope.parent != nil { // if scope.parent == nil then this is the FileScope
    errCtx := v.Ctx
    fmt.Println(reflect.TypeOf(scope.parent).String(), scope.parent.Parent() == nil)
    err := errCtx.NewError("Error: can't be exported from this scope")
    panic(err)
    return err
  }

	if key != "_" { // never set dummy vars
		// always set at this level
		scope.vars[key] = v
	}

  return nil
}

func (scope *ScopeData) SetTemplate(key string, d Template) error {
  if d.exported && scope.parent != nil { // if scope.parent == nil then this is the FileScope
    errCtx := d.ctx
    return errCtx.NewError("Error: template can't be exported from this scope")
  }

	// always set at this level
	scope.templates[key] = d

  return nil
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

func (scope *ScopeData) Eval(key string, args *tokens.Parens,
	ctx context.Context) (tokens.Token, error) {
	return eval(scope, key, args, ctx)
}


func (scope *ScopeData) GetCache() *FileCache {
  return scope.parent.GetCache()
}
