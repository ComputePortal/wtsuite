package scripts

import (
  "strings"
)

type FileScriptSorter struct {
  scripts []FileScript
}

func NewFileScriptSorter(fs []FileScript) *FileScriptSorter {
  return &FileScriptSorter{fs}
}

func (fss *FileScriptSorter) Len() int {
  return len(fss.scripts)
}

func (fss *FileScriptSorter) Less(i, j int) bool {
  return strings.Compare(fss.scripts[i].Path(), fss.scripts[j].Path()) < 0
}

func (fss *FileScriptSorter) Swap(i, j int) {
  tmp := fss.scripts[i]
  fss.scripts[i] = fss.scripts[j]
  fss.scripts[j] = tmp
}

func (fss *FileScriptSorter) Result() []FileScript {
  return fss.scripts
}
