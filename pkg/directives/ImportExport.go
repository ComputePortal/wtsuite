package directives

import (
	"path/filepath"

	"../files"
	"../parsers"
	"../tokens/context"
	tokens "../tokens/html"
	"../tokens/patterns"
	"../tree"
)

type CachedScope struct {
	scope *ScopeData
	node  *RootNode
}

var _importCache = make(map[string]CachedScope)

func parseFile(path, caller string) ([]*tokens.Tag, context.Context, error) {
	if files.XML_SYNTAX {
		p, err := parsers.NewHTMLParser(path)
		if err != nil {
			return nil, context.Context{}, err
		}

		if caller != "" {
			p = p.ChangeCaller(caller)
		}

		tags, err := p.BuildTags()
		return tags, p.NewContext(0, 1), err
	} else {
		p, err := parsers.NewUIParser(path)
		if err != nil {
			return nil, context.Context{}, err
		}

		if caller != "" {
			p = p.ChangeCaller(caller)
		}

		tags, err := p.BuildTags()
		return tags, p.NewContext(0, 1), err
	}
}

// also used by NewRoot
// abs path, so we can use this to cache the import results
func BuildFile(path string, caller string, isRoot bool) (*ScopeData, *RootNode, error) {
	var fileScope *ScopeData = nil
	var node *RootNode = nil

	if cachedScope, ok := _importCache[path]; ok && !isRoot {
		fileScope = cachedScope.scope
		node = cachedScope.node
	} else {
		if caller != "" {
			files.StartCacheUpdate(path)
		}

		tags, fileCtx, err := parseFile(path, caller)
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
		fileScope = NewRootScope(permissive)

		autoCtx := fileCtx.NewContext(0, 1)

		// TODO: should we refactor these into the Node structure?
		SetFile(fileScope, path, autoCtx)
		//SetURL(fileScope, path, autoCtx) // this is file local url, only valid in the root scope if the path is effectively also used as a html document

		// this is where the magic happens
		for _, tag := range tags {
			if IsDirective(tag.Name()) || isRoot { // if not root we can't build regular tags, because __url__ would be wrong
				if err := BuildTag(fileScope, node, tag); err != nil {
					return nil, nil, err
				}
			}
		}

		//UnsetURL(fileScope)
		_importCache[path] = CachedScope{fileScope, node}
	}

	return fileScope, node, nil
}

func addCacheDependency(dynamic bool, thisPath string, importPath string) {
	// only add cache dependency if the other direction doesn't already exist
	// the other direction can span multiple files though, so must do a nested search
	// we can do this search this the dependency structure in the cache
	if !dynamic || !files.HasUpstreamCacheDependency(importPath, thisPath) {
		files.AddCacheDependency(thisPath, importPath)
	}
}

func importExport(scope Scope, node Node, export bool, tag *tokens.Tag) error {
	if err := tag.AssertEmpty(); err != nil {
		return err
	}

	attrScope := NewSubScope(scope)

	dynamic := false
	dynamicToken_, hasDynamic := tag.RawAttributes().Get(".dynamic")
	if hasDynamic {
		if dynamicToken, err := tokens.AssertBool(dynamicToken_); err != nil {
			return err
		} else {
			dynamic = dynamicToken.Value()
		}
	} else {
		panic("expected .dynamic attribute to be set for import/export statement")
	}

	asToken_, hasAs := tag.RawAttributes().Get("as")

	nAttr := tag.RawAttributes().Len()
	if nAttr != 2 && nAttr != 3 {
		errCtx := tag.RawAttributes().Context()
		return errCtx.NewError("Error: unexpected import attributes")
	} else if nAttr == 2 || hasAs {
		if export {
			errCtx := tag.Context()
			return errCtx.NewError("Error: aggregate export not allowed for packages (pointless)")
		}

		attr, err := tag.Attributes([]string{"src"})
		if err != nil {
			return err
		}

		attr, err = attr.EvalStringDict(attrScope)
		if err != nil {
			return err
		}
		attr.Delete(".dynamic")

		srcToken, err := tokens.DictString(attr, "src")
		if err != nil {
			return err
		}

		ctx := tag.Context()
		path, _, err := files.SearchPackage(ctx.Path(), srcToken.Value(), files.UIPACKAGE_SUFFIX)
		if err != nil {
			errCtx := srcToken.Context()
			return errCtx.NewError("Error: file " + err.Error())
		}

		namespace := filepath.Base(filepath.Dir(path)) + patterns.NAMESPACE_SEPARATOR
		if hasAs {
			namespaceToken, err := tokens.AssertString(asToken_)
			if err != nil {
				return err
			}

			namespace = namespaceToken.Value() + patterns.NAMESPACE_SEPARATOR
		}

		subScope, _, err := BuildFile(path, ctx.Caller(), false)
		if err != nil {
			return err
		}

		subScope.SyncPackage(scope, false, false, !export, namespace)
		addCacheDependency(dynamic, ctx.Path(), path)
	} else if nAttr == 3 {
		attr, err := tag.Attributes([]string{"names"})
		if err != nil {
			return err
		}

		attr, err = attr.EvalStringDict(attrScope)
		if err != nil {
			return err
		}
		attr.Delete(".dynamic")

		srcToken, err := tokens.DictString(attr, "from")
		if err != nil {
			return err
		}

		namesToken, err := tokens.DictList(attr, "names")
		if err != nil {
			return err
		}

		ctx := tag.Context()
		path, _, err := files.SearchPackage(ctx.Path(), srcToken.Value(), files.UIPACKAGE_SUFFIX)
		if err != nil {
			errCtx := srcToken.Context()
			return errCtx.NewError("Error: file " + err.Error())
		}

		subScope, _, err := BuildFile(path, ctx.Caller(), false)
		if err != nil {
			return err
		}

		if err := subScope.SyncFiltered(scope, false, false, !export, "", namesToken); err != nil {
			return err
		}
		addCacheDependency(dynamic, ctx.Path(), path)
	}

	return nil
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
