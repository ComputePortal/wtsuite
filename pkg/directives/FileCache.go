package directives

import (
  "sync"
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

func (c *FileCache) IsCached(path string) bool {
  c.mutex.RLock()

  _, ok := c.entries[path]

  c.mutex.RUnlock()

  return ok
}

func (c *FileCache) Get(path string) (*FileScope, *RootNode) {
  c.mutex.RLock()

  entry := c.entries[path]

  c.mutex.RUnlock()

  return entry.scope, entry.node
}

func (c *FileCache) Set(path string, scope *FileScope, node *RootNode) {
  c.mutex.Lock()

  c.entries[path] = fileCacheEntry{scope, node}

  c.mutex.Unlock()
}

func (c *FileCache) Clear() {
  c.mutex.Lock()

  c.entries = make(map[string]fileCacheEntry)

  c.mutex.Unlock()
}
