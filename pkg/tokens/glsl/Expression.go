package glsl

type Expression interface {
  Token

  WriteExpression() string
}
