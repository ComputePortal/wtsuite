package directives

import ()

type FileScope struct {
	permissive bool
  cache *FileCache
  ScopeData
}

func NewFileScope(permissive bool, cache *FileCache) *FileScope {
  return &FileScope{permissive, cache, newScopeData(nil)}
}

func (s *FileScope) Permissive() bool {
  return s.permissive
}

func (s *FileScope) GetCache() *FileCache {
  return s.cache
}
