package directives

import (
  "github.com/computeportal/wtsuite/pkg/files"
)

type FileScope struct {
	permissive bool
  source files.Source
  cache *FileCache
  ScopeData
}

func NewFileScope(permissive bool, source files.Source, cache *FileCache) *FileScope {
  return &FileScope{permissive, source, cache, newScopeData(nil)}
}

func (s *FileScope) Permissive() bool {
  return s.permissive
}

func (s *FileScope) GetSource() files.Source {
  return s.source
}

func (s *FileScope) GetCache() *FileCache {
  return s.cache
}
