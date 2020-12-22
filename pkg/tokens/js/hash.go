package js

import (
  "../raw"
)

func HashControl(fname string) string {
  return raw.ShortHash(fname)
}
