package raw

func IsWhitespace(t Token) bool {
  return IsIndent(t) || IsNL(t)
}

func RemoveWhitespace(ts []Token) []Token {
  tsFiltered := []Token{}

  // filter out the whitespace
  for _, t := range ts {
    if !IsWhitespace(t) {
      tsFiltered = append(tsFiltered, t)
    }
  }

  return tsFiltered
}
