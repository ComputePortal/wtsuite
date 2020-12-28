package glsl

type Statement interface {
  Token

  WriteStatement(usage Usage, indent string, nl string, tab string) string
}
