package styles

import (
	"encoding/base64"
	"io/ioutil"
  "strings"

	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/math/serif"
)

func writeMathFontFace(mathFontUrl string) string {
	var b strings.Builder

	if mathFontUrl != "" {
		b.WriteString("@font-face{font-family:")
		b.WriteString(directives.MATH_FONT)
		b.WriteString(";src:url(")
		b.WriteString(mathFontUrl)
		b.WriteString(")}")
		b.WriteString(patterns.NL)
  }

  return b.String()
}

func SaveMathFont(dst string) error {
	data, err := base64.StdEncoding.DecodeString(serif.Woff2Blob)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, 0644)
}
