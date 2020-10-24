package values

type CopyCache map[Value]Value

func NewCopyCache() CopyCache {
  return make(map[Value]Value)
}
