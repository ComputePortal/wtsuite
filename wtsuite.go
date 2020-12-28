// exported wtsuite module package for 'online' transpilation by web-servers
package wtsuite

import (
  "errors"
  "sync"
  
  "github.com/computeportal/wtsuite/pkg/cache"
  "github.com/computeportal/wtsuite/pkg/directives"
  "github.com/computeportal/wtsuite/pkg/files"
  "github.com/computeportal/wtsuite/pkg/tokens/context"
  tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
  "github.com/computeportal/wtsuite/pkg/tree"
  "github.com/computeportal/wtsuite/pkg/tree/styles"
)

type Transpiler struct {
  source Source
  fileCache *directives.FileCache
  mutex *sync.RWMutex
  resultsCache map[string][]byte
  compact bool
  mathFontURL string
}

// abstraction of file system
type Source interface {
  files.Source
  // Search(callerPath string, srcPath string) (string, error) // returns abspath
  // Read(path string) ([]byte, error)
}

type FileSource struct {
  files.FileSource
}

func NewTranspiler(source Source, compact bool, mathFontURL string) *Transpiler {
  // XXX: should this be done via SetEnv-like function(s) instead?
  styles.MATH_FONT_URL = mathFontURL

  return &Transpiler{
    source, 
    directives.NewFileCache(), 
    &sync.RWMutex{}, 
    make(map[string][]byte),
    compact,
    mathFontURL,
  }
}

// pwd is automatically included
func NewFileSource(include []string) *FileSource {
  inner := files.NewFileSource(include, files.UIPACKAGE_SUFFIX)

  return &FileSource{
    *inner,
  }
}

// template doesnt need to be exported though
func (t *Transpiler) TranspileTemplate(path string, name string, args_ map[string]interface{}, cacheResult bool) ([]byte, error) {
  ctx := context.NewDummyContext()

  // convert args to tokens representation
  // should sort keys internally
  args, err := tokens.GolangStringMapToRawDict(args_, ctx)
  if err != nil {
    return nil, err
  }

  var key string
  if cacheResult {
    key = args.Dump("")

    t.mutex.RLock()

    b, ok := t.resultsCache[key]

    t.mutex.RUnlock()

    if ok {
      return b, nil
    }
  }

  fileScope, _, err := directives.BuildFile(t.source, t.fileCache, path, "", false)
  if err != nil {
    return nil, err
  }

  if !fileScope.HasTemplate(name) {
    err := errors.New("Error: template " + name + " not found in " + path)
    return nil, err
  }

  root := tree.NewRoot(ctx)
  node := directives.NewRootNode(root, directives.HTML)

  if err := directives.BuildTemplate(fileScope, node, 
    tokens.NewTag(name, args, []*tokens.Tag{}, ctx)); err != nil {
    return nil, err
  }

  // no control, no cssUrl, no jsUrl
  _, cssBundleRules, err := directives.FinalizeRoot(node, "", "", "")
  if err != nil {
    return nil, err
  }

  htmlCache := cache.HTMLCache{}
	// update the cache with the cssBundleRules
	for _, rules := range cssBundleRules { // added to file later
		htmlCache.AddCssEntry(rules, path) // path argument is irrelevant
	}

  cssContent := htmlCache.WriteCSSBundle(t.mathFontURL)

  // add to the Root
  if err := root.IncludeStyle(cssContent); err != nil {
    return nil, err
  }

  var output string
  if t.compact {
    output = root.Write("", "", "")
  } else {
    output = root.Write("", tree.NL, tree.TAB)
  }

  b := []byte(output)

  if cacheResult {
    t.mutex.Lock()

    t.resultsCache[key] = b

    t.mutex.Unlock()
  }

  return b, nil
}

func (t *Transpiler) ClearCache() {
  t.fileCache.Clear()

  t.mutex.Lock()

  t.resultsCache = make(map[string][]byte)

  t.mutex.Unlock()
}
