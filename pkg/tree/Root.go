package tree

import (
	"fmt"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	//"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

// doesn't need to implement Tag interface
// TODO: maybe it is more convenient if it DOES implement the Tag interface
type Root struct {
	tagData
}

func NewRoot(ctx context.Context) *Root {
	return &Root{tagData{"", "", false, nil, nil, make([]Tag, 0), ctx}}
}

func (t *Root) GetDocTypeAndHTML() (*DocType, *HTML, error) {
	var docType *DocType = nil
	var html *HTML = nil

	for _, child := range t.children {
		switch tt := child.(type) {
		case *DocType:
			if docType != nil {
				errCtx := context.MergeContexts(child.Context(), docType.Context())
				return nil, nil, errCtx.NewError("HTML Error: DOCTYPE defined twice")
			} else if html != nil {
				errCtx := context.MergeContexts(child.Context(), html.Context())
				return nil, nil, errCtx.NewError("HTML Error: html defined before DOCTYPE")
			}

			docType = tt
		case *HTML:
			if html != nil {
				errCtx := context.MergeContexts(child.Context(), html.Context())
				return nil, nil, errCtx.NewError("HTML Error: html defined twice")
			}

			html = tt
		default:
			errCtx := child.Context()
			return nil, nil, errCtx.NewError("HTML Error: expected DOCTYPE or html")
		}
	}

	if docType == nil {
		if AUTO_DOC_TYPE {
			docType = NewAutoDocType()
			t.children = []Tag{docType, html}
		} else {
			err := t.ctx.NewError(fmt.Sprintf("HTML Error: no !DOCTYPE defined (nChildren: %d)",
				len(t.children)))
			return nil, nil, err
		}
	}

	if html == nil {
		return nil, nil, t.ctx.NewError("HTML Error: no html defined")
	}

	return docType, html, nil
}

func (t *Root) Validate() error {
	docType, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	if err := docType.Validate(); err != nil {
		return err
	}

	if err := html.Validate(); err != nil {
		return err
	}

	return err
}

func (t *Root) VerifyElementCount(i int, ecKey string) error {
	for i, child := range t.children {
		if err := child.VerifyElementCount(i, ecKey); err != nil {
			return err
		}
	}

	return nil
}

// dummy args are just for interface
func (t *Root) Write(indent string, nl string, tab string) string {
	return t.writeChildren(indent, nl, tab)
}

func (t *Root) CollectIDs(idMap IDMap) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	return html.CollectIDs(idMap)
}

// returns the globally bundleable styles
func (t *Root) SetStyleURL(cssUrl string) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	// return map of css keys/entries for writing to bundle
	return html.SetStyleURL(cssUrl)
}

// dummy is just to respect the interface
/*func (t *Root) CollectScripts(bundle *scripts.InlineBundle) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	// bundle is only used in html, but HTML must implement Tag interface (to be a child of Root), so that's why bundle is passed in as an argument
	return html.CollectScripts(bundle)
}*/

func (t *Root) ApplyControl(control string, jsUrl string) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	return html.ApplyControl(control, jsUrl)
}

func (t *Root) IncludeStyle(styles string) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	return html.IncludeStyle(styles)
}

func (t *Root) IncludeControl(code string) error {
	_, html, err := t.GetDocTypeAndHTML()
	if err != nil {
		return err
	}

	return html.IncludeControl(code)
}
