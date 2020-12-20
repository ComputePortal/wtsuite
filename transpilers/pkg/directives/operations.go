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
}

type OperationData struct {
	target string
}

type ReplaceChildrenOp struct {
	tags  []*tokens.Tag
	scope Scope
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

func (op *OperationData) Target() string {
	return op.target
}

func (op *OperationData) SetTarget(t string) {
	op.target = t
}

func (op *ReplaceChildrenOp) Context() context.Context {
	return op.tags[0].Context()
}

func (op *AppendOp) Context() context.Context {
	return op.tags[0][0].Context()
}

func NewAppendToDefaultOp(scope Scope, tags []*tokens.Tag) (*AppendOp, error) {
	subScope := NewSubScope(scope, scope.GetNode())

	return &AppendOp{[][]*tokens.Tag{tags}, []Scope{subScope}, OperationData{"default"}}, nil
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
    return &AppendOp{
      append([][]*tokens.Tag{op.tags}, other.tags...),
      append([]Scope{op.scope}, other.scopes...),
      OperationData{op.target},
    }, nil
	default:
		panic("unrecognized")
	}
}

func (op *ReplaceChildrenOp) Apply(origScope Scope, node Node, childTokens []*tokens.Tag) error {
	for _, child := range op.tags {
		if err := BuildTag(op.scope, node, child); err != nil {
			return err
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

  op := &AppendOp{[][]*tokens.Tag{tag.Children()}, []Scope{scope}, OperationData{name}}

  return node.PushOp(op)
}

func ReplaceBlockChildren(scope Scope, node *TemplateNode, tag *tokens.Tag) error {
  name, err := getOpNameTarget("target", tag)
  if err != nil {
    return err
  }

  op := &ReplaceChildrenOp{tag.Children(), scope, OperationData{name}}

  return node.PushOp(op)
}
