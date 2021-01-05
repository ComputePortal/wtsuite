package shaders

import (
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
)

var (
	VERBOSITY = 0
)

type ShaderFile interface {
	Write(usage glsl.Usage, nl string, tab string) (string, error)
	Dependencies() []files.PathLang // src fields in script or call
	ResolveNames(scope glsl.GlobalScope) error
  EvalTypes() error
  ResolveActivity(usage glsl.Usage) error
  FinalizeInjected(usage glsl.Usage) error
  UniqueEntryPointNames(ns glsl.Namespace) error
  UniqueNames(ns glsl.Namespace) error
  CollectVersion(version *glsl.Word) (*glsl.Word, error)
  CollectVaryings(varyings map[string]string) error
  FindExportedConst(name string) *glsl.Const

	Module() glsl.Module
	Path() string
}

type ShaderFileData struct {
	path   string
	module *glsl.ModuleData
}

// if relPath is already absolute, then caller can be left empty
func newShaderFileData(path string) (ShaderFileData, error) {
	// for caching
  if files.StartCacheUpdate != nil {
    files.StartCacheUpdate(path)
  }

	p, err := parsers.NewGLSLParser(path)
	if err != nil {
		return ShaderFileData{}, err
	}

	m, err := p.BuildModule()
	if err != nil {
		return ShaderFileData{}, err
	}

	return ShaderFileData{path, m}, nil
}

func NewShaderFile(absPath string) (ShaderFile, error) {
  s, err := newShaderFileData(absPath)
  if err != nil {
    return nil, err
  }

  return &s, nil
}

func (s *ShaderFileData) Dependencies() []files.PathLang {
  return s.module.Dependencies()
}

func (s *ShaderFileData) Module() glsl.Module {
	return s.module
}

func (s *ShaderFileData) Path() string {
	return s.path
}

func (s *ShaderFileData) Write(usage glsl.Usage, nl string, tab string) (string, error) {
	return s.module.Write(usage, nl, tab)
}

func (s *ShaderFileData) ResolveNames(scope glsl.GlobalScope) error {
	return s.module.ResolveNames(scope)
}

func (s *ShaderFileData) EvalTypes() error {
  return s.module.EvalTypes()
}

func (s *ShaderFileData) ResolveActivity(usage glsl.Usage) error {
  return s.module.ResolveActivity(usage)
}

func (s *ShaderFileData) FinalizeInjected(usage glsl.Usage) error {
  return s.module.FinalizeInjected(usage)
}

func (s *ShaderFileData) UniqueEntryPointNames(ns glsl.Namespace) error {
  return nil
}

func (s *ShaderFileData) UniqueNames(ns glsl.Namespace) error {
  return s.module.UniqueNames(ns)
}

func (s *ShaderFileData) CollectVersion(version *glsl.Word) (*glsl.Word, error) {
  return s.module.CollectVersion(version)
}

func (s *ShaderFileData) CollectVaryings(varyings map[string]string) error {
  return s.module.CollectVaryings(varyings)
}

func (s *ShaderFileData) FindExportedConst(name string) *glsl.Const {
  return nil
}
