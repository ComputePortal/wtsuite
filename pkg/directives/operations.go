package directives

import (
  "strconv"
  "strings"

	"../tokens/context"
	tokens "../tokens/html"
)

type Operation interface {
	Target() string
  SetTarget(t string)
	Merge(other Operation) (Operation, error)
	Context() context.Context
	Apply(origScope Scope, newNode Node, childTokens []*tokens.Tag) error
  Popped() bool
  SetPopped() 
}

type OperationData struct {
	target string
  popped bool // operations can be used multiple times, but must be popped once
  // note that operations that have not been renamed using a 'white-space' prefix, can only be popped once
}

type ReplaceChildrenOp struct {
	tags  [][]*tokens.Tag
	scopes []Scope
	OperationData
}

type AppendOp struct {
	tags   [][]*tokens.Tag
	scopes []Scope
	OperationData
}

/*type PrependOp struct {
	OperationData
}*/

func newOperationData(target string) OperationData {
  return OperationData{target, false}
}

func (op *OperationData) Popped() bool {
  return op.popped
}

func (op *OperationData) SetPopped() {
  op.popped = true
}

func (op *OperationData) Target() string {
	return op.target
}

func (op *OperationData) SetTarget(t string) {
	op.target = t
}

func (op *ReplaceChildrenOp) Context() context.Context {
	return op.tags[0][0].Context()
}

func (op *AppendOp) Context() context.Context {
	return op.tags[0][0].Context()
}

func NewAppendToDefaultOp(scope Scope, tags []*tokens.Tag) (*AppendOp, error) {
  // can't use subscope!
	return &AppendOp{[][]*tokens.Tag{tags}, []Scope{scope}, newOperationData("default")}, nil
}

func (op *AppendOp) Merge(other_ Operation) (Operation, error) {
	if other_.Target() != op.Target() {
		panic("targets dont correspond")
	}
	switch other := other_.(type) {
	case *ReplaceChildrenOp:
		//errCtx := other.Context()
		//return nil, errCtx.NewError("Error: append is being overridden by replace children")
		return other, nil
	case *AppendOp:
		op.tags = append(op.tags, other.tags...)
		op.scopes = append(op.scopes, other.scopes...)
		return op, nil
	default:
		panic("unrecognize")
	}
}

func (op *ReplaceChildrenOp) Merge(other_ Operation) (Operation, error) {
	if other_.Target() != op.Target() {
		panic("targets dont correspond")
	}
	switch other := other_.(type) {
	case *ReplaceChildrenOp:
    return other, nil
	case *AppendOp:
    return &ReplaceChildrenOp{
      append(op.tags, other.tags...),
      append(op.scopes, other.scopes...),
      newOperationData(op.target),
    }, nil
	default:
		panic("unrecognized")
	}
}

func (op *ReplaceChildrenOp) Apply(origScope Scope, node Node, childTokens []*tokens.Tag) error {
  for i, tags := range op.tags {
    // original scope.node is probably another tag, and needs to be changed
    scope := op.scopes[i]
    for _, child := range tags {
      if err := BuildTag(scope, node, child); err != nil {
        return err
      }
    }
  }

	return nil
}

func (op *AppendOp) Apply(origScope Scope, node Node, childTokens []*tokens.Tag) error {
	for _, child := range childTokens {
		if err := BuildTag(origScope, node, child); err != nil {
			return err
		}
	}

	for i, tags := range op.tags {
    // original scope.node is probably another tag, and needs to be changed
		scope := op.scopes[i]
		for _, child := range tags {
			if err := BuildTag(scope, node, child); err != nil {
				return err
			}
		}
	}

	return nil
}

var _uniqueOpCount = 0

func NewUniqueOpTargetName() string {
  // initial whitespace makes sure there can never be a naming conflict
  res := " " + strconv.Itoa(_uniqueOpCount)

  _uniqueOpCount += 1

  return res
}

func IsUniqueOpTargetName(t string) bool {
  return strings.HasPrefix(t, " ")
}

func getOpNameTarget(key string, tag *tokens.Tag) (string, error) {
  attr, err := tag.Attributes([]string{key})
  if err != nil {
    return "", err
  }

  if attr.Len() != 1 {
    errCtx := tag.Context()
    return "", errCtx.NewError("Error: expected only " + key + " attribute")
  }

  resToken, ok := attr.Get(key)
  if !ok {
    errCtx := tag.Context()
    return "", errCtx.NewError("Error: " + key + " attribute not found")
  }

  resString, err := tokens.AssertString(resToken)
  if err != nil {
    return "", err
  }

  res := resString.Value()

  return res, nil
}

func AppendToBlock(scope Scope, node *TemplateNode, tag *tokens.Tag) error {
  name, err := getOpNameTarget("target", tag)
  if err != nil {
    return err
  }

  subScope := NewSubScope(scope)

  op := &AppendOp{[][]*tokens.Tag{tag.Children()}, []Scope{subScope}, newOperationData(name)}

  return node.PushOp(op)
}

func ReplaceBlockChildren(scope Scope, node *TemplateNode, tag *tokens.Tag) error {
  name, err := getOpNameTarget("target", tag)
  if err != nil {
    return err
  }

  subScope := NewSubScope(scope)

  op := &ReplaceChildrenOp{[][]*tokens.Tag{tag.Children()}, []Scope{subScope}, 
    newOperationData(name)}

  return node.PushOp(op)
}
