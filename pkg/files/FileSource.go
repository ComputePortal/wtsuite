package files

import (
  "errors"
  "io/ioutil"
	"path/filepath"
)

type FileSource struct {
  include []string
  pkgSuffix string // eg. __init__.wts
}

func NewFileSource(include []string, pkgSuffix string) *FileSource {
  return &FileSource{include, pkgSuffix}
}

func (s *FileSource) Search(callerPath string, srcPath string) (string, error) {
	if filepath.IsAbs(srcPath) {
		if err := AssertFile(srcPath); err != nil {
			return "", err
		} else {
			return srcPath, nil
		}
	}

	if !filepath.IsAbs(callerPath) {
    if callerPath == "" {
      panic("currentFname empty even though refFname isnt Abs: " + srcPath)
    } else {
      panic("currentFname should be absolute, got: " + callerPath)
    }
	}

	currentDir := filepath.Dir(callerPath)

	searchDirs := append([]string{currentDir}, s.include...)

	for _, dir := range searchDirs {
		if dir == "" {
			continue
		}

		fname := filepath.Join(dir, srcPath)

		if s.pkgSuffix != "" && IsDir(fname) {
			fname = filepath.Join(fname, s.pkgSuffix)
		}

		if err := AssertFile(fname); err == nil {
			if absFname, err := filepath.Abs(fname); err != nil {
				return "", err
			} else {
				return absFname, nil
			}
		}
	}

	err := errors.New(srcPath + " not found")

	return "", err
}

func (s *FileSource) Read(path string) ([]byte, error) {
  return ioutil.ReadFile(path)
}
