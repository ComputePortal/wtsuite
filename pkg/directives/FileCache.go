package directives

import (
  "sync"

  tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type fileCacheEntry struct {
  scope *FileScope
  node *RootNode
}

type FileCache struct {
  // key is the file path
  entries map[string]fileCacheEntry
  mutex *sync.RWMutex
}

func NewFileCache() *FileCache {
  return &FileCache{make(map[string]fileCacheEntry), &sync.RWMutex{}}
}

func pathAndParametersToKey(path string, parameters *tokens.Parens) string {
  key := path

  if parameters != nil {
    key += parameters.Dump("")
  }

  return key
}

func (c *FileCache) IsCached(path string, parameters *tokens.Parens) bool {
  key := pathAndParametersToKey(path, parameters)

  c.mutex.RLock()

  _, ok := c.entries[key]

  c.mutex.RUnlock()

  return ok
}

func (c *FileCache) Get(path string, parameters *tokens.Parens) (*FileScope, *RootNode) {
  key := pathAndParametersToKey(path, parameters)

  c.mutex.RLock()

  entry := c.entries[key]

  c.mutex.RUnlock()

  return entry.scope, entry.node
}

func (c *FileCache) Set(path string, parameters *tokens.Parens, scope *FileScope, node *RootNode) {
  key := pathAndParametersToKey(path, parameters)

  c.mutex.Lock()

  c.entries[key] = fileCacheEntry{scope, node}

  c.mutex.Unlock()
}

func (c *FileCache) Clear() {
  c.mutex.Lock()

  c.entries = make(map[string]fileCacheEntry)

  c.mutex.Unlock()
}
