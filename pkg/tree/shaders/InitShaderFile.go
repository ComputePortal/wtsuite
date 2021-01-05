package shaders

import (
	"github.com/computeportal/wtsuite/pkg/tokens/glsl"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type InitShaderFile struct {
  main glsl.Variable
  ShaderFileData
}

func NewInitShaderFile(path string) (*InitShaderFile, error) {
  shaderFileData, err := newShaderFileData(path)
  if err != nil {
    return nil, err
  }

  return &InitShaderFile{nil, shaderFileData}, nil
}

func (s *InitShaderFile) ResolveNames(scope glsl.GlobalScope) error {
  mainVar, err := s.module.ResolveEntryNames(scope)
  if err != nil {
    return err
  }

  // check the value
  if err := values.AssertMainFunction(mainVar.GetValue()); err != nil {
    return err
  }

  s.main = mainVar

  return nil
}

func (s *InitShaderFile) ResolveActivity(usage glsl.Usage) error {
  if err := usage.Use(s.main, s.main.Context()); err != nil {
    return err
  }

  return s.module.ResolveActivity(usage)
}

func (s *InitShaderFile) UniqueEntryPointNames(ns glsl.Namespace) error {
  return ns.OrigName(s.main)
}

func (s *InitShaderFile) FindExportedConst(name string) *glsl.Const {
  return s.module.FindExportedConst(name)
}
