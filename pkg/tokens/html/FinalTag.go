package html

type FinalTag interface {
  NumChildren() int
  FinalParent() FinalTag
}
