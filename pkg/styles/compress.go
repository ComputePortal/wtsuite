package styles

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

func compress(raw string) string {
	return compressNested(raw, true)
}

// scan the input string, removing the queries between "@....{....}" and putting those in amap
func compressNested(raw string, topLevel bool) string {
	queries := make(map[string]string)
	plainClasses := make(map[string]string) // not done in top level, comes before other output (inside queries), map key is actually body, map value is comma separated list of classes
	// plainClasses compression can be turned off by always setting topLevel==true

	var output strings.Builder
	var key strings.Builder
	var body strings.Builder

	inQueryKey := false
	inQueryBody := false
	inString := false

	inClassKey := false
	inClassBody := false

	var prevNonWhite byte = '}'

	braceCount := 0
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		if inQueryKey || (!topLevel && inClassKey) {
			if inQueryKey {
				if c == '{' {
					inQueryKey = false
					inQueryBody = true

					keyStr := key.String()
					if keyStr == "@font-face" {
						// not allowed to combine in this case
						inQueryBody = false
						output.WriteString(keyStr)
						output.WriteByte(c)
						key.Reset()
					}
				} else {
					key.WriteByte(c)
				}
			} else { // inClassKey && !topLevel
				if !inClassKey {
					panic("algo error")
				}

				if c == '{' {
					keyStr := key.String()
					inClassKey = false
					if !patterns.CSS_PLAIN_CLASS_REGEXP.MatchString(keyStr) {
						output.WriteString(keyStr) // first dot should be included in keyStr
						output.WriteByte(c)
						key.Reset()
					} else {
						inClassBody = true
					}
				} else {
					key.WriteByte(c)
				}
			}
		} else if inQueryBody || (!topLevel && inClassBody) {
			if inString {
				if c == '"' || c == '\'' {
					inString = false
				}
			} else {
				if c == '"' || c == '\'' {
					inString = true
				}
			}

			if !inString {
				if c == '}' {
					if braceCount == 0 {
						// include the next newline(s) too
						for i < len(raw)-1 && raw[i+1] == '\n' {
							//body.WriteByte('\n')
							i++
						}

						// save the key -> body
						keyStr := key.String()
						bodyStr := body.String()

						if inQueryBody {
							if prevBody, ok := queries[keyStr]; !ok {
								queries[keyStr] = bodyStr
							} else {
								// XXX: is this addition slow?
								queries[keyStr] = prevBody + bodyStr
							}
							inQueryBody = false
						} else {
							if !inClassBody {
								panic("algo error")
							}

							if prevLst, ok := plainClasses[bodyStr]; ok {
								plainClasses[bodyStr] = prevLst + "," + keyStr
							} else {
								plainClasses[bodyStr] = keyStr
							}
							inClassBody = false
						}

						key.Reset()
						body.Reset()
						prevNonWhite = '}'
					} else if braceCount < 1 {
						panic("bad brace count")
					} else {
						braceCount--
						body.WriteByte(c)
					}
				} else if c == '{' {
					prevNonWhite = c
					braceCount++
					body.WriteByte(c)
				} else {
					body.WriteByte(c)
				}
			} else {
				body.WriteByte(c)
			}
		} else if c == '@' && prevNonWhite == '}' {
			inQueryKey = true
			key.WriteByte(c)
		} else if c == '.' && prevNonWhite == '}' && !topLevel {
			inClassKey = true
			key.WriteByte(c)
		} else {
			if !(c == ' ' || c == '\n' || c == '\t') {
				prevNonWhite = c
			}

			output.WriteByte(c)
		}
	}

	var finalOutput strings.Builder

	for body, key := range plainClasses {
		finalOutput.WriteString(key)
		finalOutput.WriteByte('{')
		finalOutput.WriteString(body)
		finalOutput.WriteByte('}')
		finalOutput.WriteString(patterns.NL)
	}

	finalOutput.WriteString(output.String())

	for k, q := range queries {
		finalOutput.WriteString(k)
		finalOutput.WriteByte('{')
		finalOutput.WriteString(compressNested(q, false)) // set to true to avoid plainClass compression
		finalOutput.WriteByte('}')
		finalOutput.WriteString(patterns.NL)
	}

	return finalOutput.String()
}
