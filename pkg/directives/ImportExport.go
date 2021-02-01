package directives

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type CachedScope struct {
	scope *FileScope
	node  *RootNode
}

func parseFile(path string) ([]*tokens.Tag, context.Context, error) {
  p, err := parsers.NewTemplateParser(path)
  if err != nil {
    return nil, context.Context{}, err
  }

  tags, err := p.BuildTags()
  return tags, p.NewContext(0, 1), err
}

// returns abs path
func searchFile(relPath string, ctx context.Context) (string, error) {
  var absPath string
  var err error = nil
  // TODO: make this safer for windows
  if strings.HasPrefix(relPath, "./") || strings.HasPrefix(relPath, "../") {
    absPath, err = files.Search(ctx.Path(), relPath)
  } else if strings.HasPrefix(relPath, "/") { 
    absPath = relPath
  } else {
    absPath, err = files.SearchTemplate(ctx.Path(), relPath)
  } 

  return absPath, err
}

func evalParameters(fileScope *FileScope, parTag *tokens.Tag, parameters *tokens.Parens) error {
  parParens_, ok := parTag.RawAttributes().Get("parameters")
  if !ok {
    panic("unexpected")
  }

  parParens, err := tokens.AssertParens(parParens_)
  if err != nil {
    panic("unexpected")
  }

  if err := functions.CompleteArgsAndFillScope(fileScope, parameters, parParens); err != nil {
    return err
  }

  return nil
}

// also used by NewRoot
// abs path, so we can use this to cache the import results
// incoming parameters should be evaluated
func BuildFile(cache *FileCache, path string, isRoot bool, parameters *tokens.Parens) (*FileScope, *RootNode, error) {
	var fileScope *FileScope = nil
	var node *RootNode = nil

	if cache.IsCached(path, parameters) && !isRoot {
		fileScope, node = cache.Get(path, parameters)
	} else {
		if files.StartCacheUpdate != nil {
			files.StartCacheUpdate(path)
		}

		tags, fileCtx, err := parseFile(path)
		if err != nil {
			return nil, nil, err
		}

    permissive := false
    if len(tags) > 0 && tags[0].Name() == "permissive" {
      permissive = true
      tags = tags[1:]
    }

		root := tree.NewRoot(fileCtx)
		node = NewRootNode(root, HTML)
		fileScope = NewFileScope(permissive, cache)

		autoCtx := fileCtx.NewContext(0, 1)

		// TODO: should we refactor these into the Node structure?
		SetFile(fileScope, path, autoCtx)
		//SetURL(fileScope, path, autoCtx) // this is file local url, only valid in the root scope if the path is effectively also used as a html document

    if len(tags) > 0 && tags[0].Name() == "parameters" {
      if err := evalParameters(fileScope, tags[0], parameters); err != nil {
        return nil, nil, err
      }
      tags = tags[1:]
    } else if parameters != nil {
      errCtx := parameters.Context()
      return nil, nil, errCtx.NewError("Error: module \"" + files.Abbreviate(path) + "\" doesn't accept parameters")
    }

    // allow circular imports, by already setting the result here
    cache.Set(path, parameters, fileScope, node)

		// this is where the magic happens
		for _, tag := range tags {
			if IsDirective(tag.Name()) || isRoot { // if not root we can't build regular tags, because __url__ would be wrong
				if err := BuildTag(fileScope, node, tag); err != nil {
					return nil, nil, err
				}
			}
		}

		//UnsetURL(fileScope)
    //cache.Set(path, parameters, fileScope, node)
	}

	return fileScope, node, nil
}

func addCacheDependency(dynamic bool, thisPath string, importPath string) {
	// only add cache dependency if the other direction doesn't already exist
	// the other direction can span multiple files though, so must do a nested search
	// we can do this search in the dependency structure of the cache
  if files.StartCacheUpdate != nil {
    if !dynamic || !files.HasUpstreamCacheDependency(importPath, thisPath) {
      if thisPath != importPath {
          files.AddCacheDependency(thisPath, importPath)
      }
    }
  }
}

func importExport(dstScope Scope, node Node, export bool, tag *tokens.Tag) error {
  ctx := tag.Context()

	if err := tag.AssertEmpty(); err != nil {
		return err
	}

	attrScope := NewSubScope(dstScope)

  rawAttr := tag.RawAttributes()

  // extract "names" and "parameters", so that they can be evaluated correctly
  namesToken__, ok := rawAttr.Get("names")
  if !ok {
    panic("expected names")
  }
  namesToken_, err := tokens.AssertRawDict(namesToken__)
  if err != nil {
    panic(err)
  }
  namesToken, err := namesToken_.EvalStringDict(attrScope)
  if err != nil {
    return err
  }

  var parameters *tokens.Parens = nil
  if parameters_, ok := rawAttr.Get("parameters"); ok {
    parameters, err = tokens.AssertParens(parameters_)
    if err != nil {
      return err
    }

    parameters, err = parameters.EvalAsArgs(attrScope)
    if err != nil {
      return err
    }
  }

  dynamicToken_, ok := rawAttr.Get(".dynamic")
  if err != nil {
    panic(err)
  }
  dynamicToken, err := tokens.AssertBool(dynamicToken_)
  if err != nil {
    panic(err)
  }
	dynamic := dynamicToken.Value()

  fromToken__, ok := rawAttr.Get("from")
  if !ok {
    panic("don't know what to do")
  }
  fromToken_, err := fromToken__.Eval(attrScope)
  if err != nil {
    return err
  }
  fromToken, err := tokens.AssertString(fromToken_)
  if err != nil {
    return err
  }

  path := fromToken.Value()
  absPath, err := searchFile(path, ctx)
  if err != nil {
    return err
  }

  srcScope, _, err := BuildFile(dstScope.GetCache(), absPath, false, parameters)
  if err != nil {
    return err
  }

  addCacheDependency(dynamic, ctx.Path(), absPath)

  if namespaceToken, ok := namesToken.Get("*"); ok {
    namespaceToken, err := tokens.AssertString(namespaceToken)
    if err != nil {
      return err
    }

    if namespaceToken.Value() == "*" {
      srcScope.SyncPackage(dstScope, false, false, !export, "")
    } else {
      srcScope.SyncPackage(dstScope, false, false, !export, namespaceToken.Value() + patterns.NAMESPACE_SEPARATOR)
    }
  } else {
		if err := srcScope.SyncFiltered(dstScope, false, false, !export, "", namesToken); err != nil {
			return err
		}
  }

	return nil
}

// expects two strings
func evalDynamicImport(scope Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  var err error
  args_, err = args_.EvalAsArgs(scope)
  if err != nil {
    return nil, err
  }

  args, err := functions.CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

  if len(args) != 2 {
    return nil, ctx.NewError("Error: expected 2 arguments")
  }

  fileToken, err := tokens.AssertString(args[0])
  if err != nil {
    return nil, err
  }

  relPath := fileToken.Value()
  if relPath == "." {
    return nil, ctx.NewError("Error: can't import dynamically from self")
  }

  absPath, err := searchFile(relPath, ctx)
  if err != nil {
    return nil, err
  }

  nameToken, err := tokens.AssertString(args[1])
  if err != nil {
    return nil, err
  }

  // parameters not (yet) permitted
  importedScope, _, err := BuildFile(scope.GetCache(), absPath, false, nil)
  if err != nil {
    return nil, err
  }

  if importedScope.HasTemplate(nameToken.Value()) {
    return nil, ctx.NewError("Error: can't dynamically import \"" + nameToken.Value() + "\" from \"" + files.Abbreviate(absPath) + "\" because it is a template")
  }

  if !importedScope.HasVar(nameToken.Value()) {
    return nil, ctx.NewError("Error: \"" + nameToken.Value() + "\" not found in \"" + files.Abbreviate(absPath) + "\"")
  }

  valVar := importedScope.GetVar(nameToken.Value())
  return valVar.Value, nil
}

// doesnt change the node, but node can be used for elementCount
func Import(scope Scope, node Node, tag *tokens.Tag) error {
	return importExport(scope, node, false, tag)
}

// doesnt change the node, but node can be used for elementCount
func Export(scope Scope, node Node, tag *tokens.Tag) error {
	return importExport(scope, node, true, tag)
}

var _importOk = registerDirective("import", Import)
var _exportOk = registerDirective("export", Export)
