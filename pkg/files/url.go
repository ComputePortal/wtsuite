package files

import (
  "regexp"
)

func IsURL(s string) bool {
  re := regexp.MustCompile("^[a-zA-Z]*[:][/][/]")
  return re.MatchString(s)
}
