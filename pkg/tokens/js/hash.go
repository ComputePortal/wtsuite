package js

import (
  "github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func HashControl(fname string) string {
  return raw.ShortHash(fname)
}
