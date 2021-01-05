package files

import (
  "errors"
  "io/ioutil"
)

type SemVerRange struct {
  min *SemVer // can be nil for -infty, inclusive
  max *SemVer // can be nil for +infty, exclusive
}

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
    return "", errors.New("Error: no valid package versions found for " + dir)
  } else {
    return files[iBest].Name(), nil
  }
}
