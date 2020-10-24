package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

var AUTO_DOC_TYPE = true

type DocType struct {
	tagData
}

func NewDocType(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("!DOCTYPE", true, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &DocType{td}, nil
}

func NewAutoDocType() *DocType {
	ctx := context.NewDummyContext()

	attr := tokens.NewEmptyStringDict(ctx)
	attr.Set("html", tokens.NewValueString("", ctx))

	td, err := newTag("!DOCTYPE", true, attr, ctx)
	if err != nil {
		panic(err)
	}

	return &DocType{td}
}

func (t *DocType) Validate() error {
	if t.NumChildren() != 0 {
		panic("should've been caught during construction")
	}

	t.attributes.Delete("__elementCount__")
	t.attributes.Delete("__elementCountFolded__")

	if t.attributes.Len() != 1 {
		errCtx := t.attributes.Context()
		return errCtx.NewError("Error: expected 1 attribute")
	}

	ok, err := tokens.DictHasFlag(t.attributes, "html")
	if err != nil {
		return err
	}

	if !ok {
		errCtx := t.attributes.Context()
		return errCtx.NewError("Error: html flag not found")
	}

	return nil
}
