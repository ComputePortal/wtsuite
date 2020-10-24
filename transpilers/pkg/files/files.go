package files

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"../tokens/context"
)

var htmlppPath = os.Getenv("HTMLPPPATH")
var includeDirs = filepath.SplitList(htmlppPath)

var (
	XML_SYNTAX = false
	JS_MODE    = false
	VERBOSITY  = 0
)

const (
  UIFILE_EXT = ".wtt"
  JSFILE_EXT = ".wts"

	UIPACKAGE_SUFFIX = "__init__" + UIFILE_EXT
	JSPACKAGE_SUFFIX = "__init__" + JSFILE_EXT
)

var StartCacheUpdate func(fname string) = nil
var AddCacheDependency func(fname string, dep string) = nil
var HasUpstreamCacheDependency func(thisPath string, upstreamPath string) bool = nil

func PrependIncludeDirs(dirs []string) {
	includeDirs = append(dirs, includeDirs...)
}

func AppendIncludeDirs(dirs []string) {
	includeDirs = append(includeDirs, dirs...)
}

func IsFile(fname string) bool {
	if info, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else if info.IsDir() {
		return false
	} else {
		return true
	}
}

func AssertFile(fname string) error {
	if info, err := os.Stat(fname); os.IsNotExist(err) {
		return errors.New("doesnt't exist")
	} else if err != nil {
		return err
	} else if info.IsDir() {
		return errors.New("is a directory")
	} else {
		return nil
	}
}

func IsDir(dname string) bool {
	if info, err := os.Stat(dname); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else if !info.IsDir() {
		return false
	} else {
		return true
	}
}

func AssertDir(dname string) error {
	if info, err := os.Stat(dname); os.IsNotExist(err) {
		return errors.New("doesn't exist")
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return errors.New("is not a directory")
	} else {
		return nil
	}
}

func search(currentFname string, refFname string, pkgSuffix string) (absFname string, isPackage bool, found error) {
	if filepath.IsAbs(refFname) {
		if err := AssertFile(refFname); err != nil {
			return "", false, err
		} else {
			return refFname, false, nil
		}
	}

	if !filepath.IsAbs(currentFname) {
    if currentFname == "" {
      panic("currentFname empty even though refFname isnt Abs: " + refFname)
    } else {
      panic("currentFname should be absolute, got: " + currentFname)
    }
	}

	currentDir := filepath.Dir(currentFname)

	searchDirs := append([]string{currentDir}, includeDirs...)

	foundPackage := false

	for _, dir := range searchDirs {
		if dir == "" {
			continue
		}

		fname := filepath.Join(dir, refFname)

		if pkgSuffix != "" && IsDir(fname) {
			fname = filepath.Join(fname, pkgSuffix)
			foundPackage = true
		}

		if err := AssertFile(fname); err == nil {
			if absFname, err := filepath.Abs(fname); err != nil {
				return "", false, err
			} else {
				return absFname, foundPackage, nil
			}
		}
	}

	err := errors.New("not found")
	return "", false, err
}

// currentFname is the caller, refFname is the file we are trying to find
func Search(currentFname string, refFname string) (string, error) {
	absPath, _, err := search(currentFname, refFname, "")
	return absPath, err
}

func SearchPackage(currentFname string, refFname string, pkgSuffix string) (string, bool, error) {
	return search(currentFname, refFname, pkgSuffix)
}

func Abbreviate(path string) string {
	return context.Abbreviate(path)
}

// path is just used for info
func WriteFile(path string, target string, content []byte) error {
	if VERBOSITY >= 2 {
		fmt.Println(Abbreviate(path) + " -> " + Abbreviate(target))
	}

	if !filepath.IsAbs(path) {
		panic("should be abs")
	}

	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(target, content, 0644); err != nil {
		return err
	}

	return nil
}

// guarantee that each file is only visited once
// ext includes the period (eg. '.wts' for script files)
func WalkFiles(dir string, ext string, fn func(string) error) error {
  done := make(map[string]string)

  if !filepath.IsAbs(dir) {
    var err error 
    dir, err = filepath.Abs(dir)
    if err != nil {
      return err
    }
  }

  if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return errors.New("Error: unable to walk file tree at \"" + dir + "\"")
    }

    if filepath.Ext(path) == ext && !info.IsDir() {

      if _, ok := done[path]; !ok {
        if err := fn(path); err != nil {
          return err
        }

        done[path] = path
      }
    }

    return nil
  }); err != nil {
    return err
  }

  return nil
}
