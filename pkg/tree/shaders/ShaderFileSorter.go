package shaders

import (
  "strings"
)

type ShaderFileSorter struct {
  shaders []ShaderFile
}

func NewShaderFileSorter(fs []ShaderFile) *ShaderFileSorter {
  return &ShaderFileSorter{fs}
}

func (fss *ShaderFileSorter) Len() int {
  return len(fss.shaders)
}

func (fss *ShaderFileSorter) Less(i, j int) bool {
  return strings.Compare(fss.shaders[i].Path(), fss.shaders[j].Path()) < 0
}

func (fss *ShaderFileSorter) Swap(i, j int) {
  tmp := fss.shaders[i]
  fss.shaders[i] = fss.shaders[j]
  fss.shaders[j] = tmp
}

func (fss *ShaderFileSorter) Result() []ShaderFile {
  return fss.shaders
}
