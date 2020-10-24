package directives

import (
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
)

type Operation interface {
	MapBlocks(blocks *tokens.StringDict) error
	ID() string
	Merge(other Operation) (Operation, error)
	Context() context.Context
	Apply(origScope Scope, parentNode Node, newNode Node, tag tree.Tag, childTokens []*tokens.Tag) error
}

type OperationData struct {
	id string
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

type ReplaceOp struct {
	tags  []*tokens.Tag
	scope Scope
	OperationData
}

func (op *OperationData) ID() string {
	return op.id
}

func (op *ReplaceOp) Context() context.Context {
	return op.tags[0].Context()
}

func (op *ReplaceChildrenOp) Context() context.Context {
	return op.tags[0].Context()
}

func (op *AppendOp) Context() context.Context {
	return op.tags[0][0].Context()
}

func (op *OperationData) MapBlocks(blocks *tokens.StringDict) error {
	if op.id == "" {
		if val_, ok := blocks.Get("default"); ok {
			val, err := tokens.AssertString(val_)
			if err != nil {
				return err
			}

			op.id = val.Value()
		}
	} else {

		if val_, ok := blocks.Get(op.id); ok {
			val, err := tokens.AssertString(val_)
			if err != nil {
				return err
			}

			op.id = val.Value()
		}
	}

	return nil
}

func (op *ReplaceChildrenOp) MapBlocks(blocks *tokens.StringDict) error {
	if err := op.OperationData.MapBlocks(blocks); err != nil {
		return err
	}

	return nil
}

func NewReplaceChildrenOp(scope Scope, tag *tokens.Tag) (Operation, error) {
	subScope := NewSubScope(scope, scope.GetNode())
	if err := tag.AssertNoAttributes(); err != nil {
		return nil, err
	}

	id := tag.Name()[1:]
	if id == "" {
		errCtx := tag.Context()
		return nil, errCtx.NewError("Error: replace children id can't be empty")
	}
	return &ReplaceChildrenOp{tag.Children(), subScope, OperationData{id}}, nil
}

func buildAppendPrependAttributes(scope Scope, tag *tokens.Tag) (string, int, error) {
	attr, err := buildAttributes(scope, nil, tag, []string{"id", "pos"})
	if err != nil {
		return "", 0, err
	}

	idToken_, ok := attr.Get("id")
	if !ok {
		errCtx := attr.Context()
		return "", 0, errCtx.NewError("Error: id not found")
	}

	idToken, err := tokens.AssertString(idToken_)
	if err != nil {
		return "", 0, err
	}

	pos := -1
	//posToken_, ok := attr.Get("pos")
	if ok {
		errCtx := tag.Context()
		return "", 0, errCtx.NewError("Error: append pos no longer supported")
		/*posToken, err := tokens.AssertInt(posToken_)
		if err != nil {
			return "", 0, err
		}

		pos = posToken.Value()
		if pos < 0 {
			errCtx := tag.Context()
			return "", 0, errCtx.NewError("Error: negative position not allowed")
		}*/
	}

	return idToken.Value(), pos, nil
}

func NewAppendOp(scope Scope, tag *tokens.Tag) (Operation, error) {
	subScope := NewSubScope(scope, scope.GetNode())

	id, _, err := buildAppendPrependAttributes(subScope, tag)
	if err != nil {
		return nil, err
	}
	return &AppendOp{[][]*tokens.Tag{tag.Children()}, []Scope{subScope}, OperationData{id}}, nil
}

/*func NewPrependOp(scope Scope, tag *tokens.Tag) (Operation, error) {
	subScope := NewSubScope(scope, scope.GetNode())

	id, _, err := buildAppendPrependAttributes(subScope, tag)
	if err != nil {
		return nil, err
	}
	return &PrependOp{OperationData{id, tag.Children(), subScope}}, nil
}*/

func NewAppendToDefaultOp(scope Scope, tags []*tokens.Tag) (Operation, error) {
	subScope := NewSubScope(scope, scope.GetNode())

	return &AppendOp{[][]*tokens.Tag{tags}, []Scope{subScope}, OperationData{""}}, nil
}

func NewReplaceOp(scope Scope, tag *tokens.Tag) (Operation, error) {
	subScope := NewSubScope(scope, scope.GetNode())

	attr, err := buildAttributes(subScope, nil, tag, []string{"id"})
	if err != nil {
		return nil, err
	}

	idToken_, ok := attr.Get("id")
	if !ok {
		errCtx := attr.Context()
		return nil, errCtx.NewError("Error: id not found")
	}

	idToken, err := tokens.AssertString(idToken_)
	if err != nil {
		return nil, err
	}

	id := idToken.Value()
	if id == "" {
		errCtx := tag.Context()
		return nil, errCtx.NewError("Error: replace id cant be empty")
	}

	return &ReplaceOp{tag.Children(), subScope, OperationData{id}}, nil
}

func (op *ReplaceOp) Merge(other_ Operation) (Operation, error) {
	if other_.ID() != op.ID() {
		panic("ids dont correspond")
	}
	switch other := other_.(type) {
	case *ReplaceOp:
		errCtx := other.Context()
		return nil, errCtx.NewError("Error: replace conflict")
	case *ReplaceChildrenOp:
		return op, nil
	case *AppendOp:
		return op, nil
	default:
		panic("unrecognized")
	}
}

func (op *AppendOp) Merge(other_ Operation) (Operation, error) {
	if other_.ID() != op.ID() {
		panic("ids dont correspond")
	}
	switch other := other_.(type) {
	case *ReplaceOp:
		return other, nil
	case *ReplaceChildrenOp:
		errCtx := other.Context()
		return nil, errCtx.NewError("Error: append is being overridden by replace children")
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
	if other_.ID() != op.ID() {
		panic("ids dont correspond")
	}
	switch other := other_.(type) {
	case *ReplaceOp:
		return other, nil
	case *ReplaceChildrenOp:
		errCtx := other.Context()
		return nil, errCtx.NewError("Error: merging of replace children op not yet implemented")
	case *AppendOp:
		errCtx := other.Context()
		return nil, errCtx.NewError("Error: merging of replace children op not yet implemented")
		return op, nil
	default:
		panic("unrecognized")
	}
}

func (op *ReplaceOp) Apply(origScope Scope, parentNode Node, newNode Node, tag tree.Tag, childTokens []*tokens.Tag) error {
	// ignore all incoming args excepts parentNode
	for _, child := range op.tags {
		if err := BuildTag(op.scope, parentNode, child); err != nil {
			return err
		}
	}

	return nil
}

func (op *ReplaceChildrenOp) Apply(origScope Scope, parentNode Node, newNode Node, tag tree.Tag, childTokens []*tokens.Tag) error {
	if tag != nil {
		if err := parentNode.AppendChild(tag); err != nil {
			return err
		}
	}

	for _, child := range op.tags {
		if err := BuildTag(op.scope, newNode, child); err != nil {
			return err
		}
	}

	return nil
}

func (op *AppendOp) Apply(origScope Scope, parentNode Node, newNode Node, tag tree.Tag, childTokens []*tokens.Tag) error {
	if tag != nil {
		if err := parentNode.AppendChild(tag); err != nil {
			return err
		}
	}

	for _, child := range childTokens {
		if err := BuildTag(origScope, newNode, child); err != nil {
			return err
		}
	}

	for i, tags := range op.tags {
		scope := op.scopes[i]
		for _, child := range tags {
			if err := BuildTag(scope, newNode, child); err != nil {
				return err
			}
		}
	}

	return nil
}
