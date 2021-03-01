package files

import (
  "errors"
  "io/ioutil"
  "path/filepath"
)

var LATEST bool = false

type SemVerRange struct {
  min *SemVer // can be nil for -infty, inclusive
  max *SemVer // can be nil for +infty, exclusive
}

func NewSemVerRange(min *SemVer, max *SemVer) *SemVerRange {
  return &SemVerRange{min, max}
}

func (sr *SemVerRange) Min() *SemVer {
  return sr.min
}

func (sr *SemVerRange) Max() *SemVer {
  return sr.max
}

// returns empty string if no relevant version found
func (sr *SemVerRange) FindBestVersion(dir string) (string, error) {
  files, err := ioutil.ReadDir(dir)
  if err != nil {
    return "", err
  }

  iBest := -1

  for i, file := range files {
    if !file.IsDir() {
      continue
    }

    semVer, err := ParseSemVer(file.Name())
    if err != nil {
      return "", errors.New("Error: package " + dir + " version is not a semver")
    }

    if sr.min != nil {
      if sr.min.After(semVer) {
        continue
      }
    }

    if sr.max == nil || sr.max.After(semVer) { 
      iBest = i
    }
  }

  if iBest == -1 {
    return "", nil
  } else {
    return filepath.Join(dir, files[iBest].Name()), nil
  }
}

func (sr *SemVerRange) Contains(semVer *SemVer) bool {
  if sr.min != nil {
    if sr.min.After(semVer) {
      return false
    }
  }

  if sr.max != nil {
    if sr.max.After(semVer) {
      return true
    } else {
      return false
    }
  } else {
    return true
  }
}
