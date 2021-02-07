package styles

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
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
