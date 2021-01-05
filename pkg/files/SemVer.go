package files

import (
  "errors"
  "strconv"
  "strings"
)

type SemVer struct {
  major int
  minor int
  patch int
  meta string
}

func ParseSemVer(str string) (*SemVer, error) {
  if str == "" {
    return nil, nil
  }

  parts := strings.Split(str, "-")

  meta := ""
  if len(parts) == 2 {
    meta = parts[1]
  } else if len(parts) != 1 {
    return nil, errors.New("Error: bad semver")
  }

  majorMinorPatch := strings.Split(parts[0], ".")

  if len(majorMinorPatch) != 3 {
    return nil, errors.New("Error: bad semver")
  }

  major, err := strconv.ParseUint(majorMinorPatch[0], 10, 64)
  if err != nil {
    return nil, err
  }

  minor, err := strconv.ParseUint(majorMinorPatch[0], 10, 64)
  if err != nil {
    return nil, err
  }

  patch, err := strconv.ParseUint(majorMinorPatch[0], 10, 64)
  if err != nil {
    return nil, err
  }

  return &SemVer{int(major), int(minor), int(patch), meta}, nil
}

func (s *SemVer) After(other *SemVer) bool {
  if other == nil {
    panic("other can't be nil")
  }

  if s.major > other.major {
    return true
  } else if s.major < other.major {
    return false
  } else {
    if s.minor > other.minor {
      return true
    } else if s.minor < other.minor {
      return false
    } else {
      if s.patch > other.patch {
        return true
      } else if s.patch < other.patch {
        return false
      } else {
        if s.meta == "" && other.meta != "" {
          return true
        } else {
          ci := strings.Compare(s.meta, other.meta)
          if ci > 0 {
            return true
          } else {
            return false
          }
        }
      }
    }
  }
}
