package tree

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
	"../tokens/js"

	"./scripts"
	"./styles"
)

type HTML struct {
	classes []string
	tagData
}

func NewHTML(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("html", false, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &HTML{[]string{}, td}, nil
}

func (t *HTML) getHeadBody() (*Head, *Body, error) {
	var head *Head = nil
	var body *Body = nil

	for _, child := range t.children {
		switch tt := child.(type) {
		case *Head:
			if head != nil {
				errCtx := context.MergeContexts(child.Context(), head.Context())
				return nil, nil, errCtx.NewError("HTML Error: head defined twice")
			} else if body != nil {
				errCtx := context.MergeContexts(child.Context(), body.Context())
				return nil, nil, errCtx.NewError("HTML Error: body defined before head")
			}

			head = tt
		case *Body:
			if body != nil {
				errCtx := context.MergeContexts(child.Context(), body.Context())
				return nil, nil, errCtx.NewError("HTML Error: body defined twice")
			}

			body = tt
		default:
			errCtx := child.Context()
      err := errCtx.NewError("HTML Error: expected head or body (" + child.Name() + ")")
      context.AppendContextString(err, "Info: children of this tag", t.Context())
			return nil, nil, err
		}
	}

	if head == nil {
		return nil, nil, t.ctx.NewError("HTML Error: no head defined")
	}

	if body == nil {
		return nil, nil, t.ctx.NewError("HTML Error: no body defined")
	}

	return head, body, nil
}

func (t *HTML) Validate() error {
	head, body, err := t.getHeadBody()

	if err != nil {
		return err
	}

	if err := head.Validate(); err != nil {
		return err
	}

	if err := body.Validate(); err != nil {
		return err
	}

	return err
}

func (t *HTML) CollectIDs(idMap IDMap) error {
	_, body, err := t.getHeadBody()
	if err != nil {
		return err
	}

	idMap.Set("html", t)

	return body.CollectIDs(idMap)
}

func (t *HTML) CollectClasses(classMap ClassMap) error {
	classes, err := t.collectAttributeClasses()
	if err != nil {
		return err
	}

	if len(t.classes) != 0 && len(classes) != 0 {
		panic("unexpected")
	}

	t.classes = classes

	// XXX: does this tag need to be appended to the classMap?
	//for _, cl := range t.classes {
	//classMap.AppendTag(cl, t)
	//}

	_, body, err := t.getHeadBody()
	if err != nil {
		return err
	}

	return body.CollectClasses(classMap)
}

func (t *HTML) collectAllRule(ss styles.Sheet) error {
	// add the * {margin: 0, padding: 0} rule
	allStyle := tokens.NewEmptyStringDict(t.Context())

	// add the dict entries
	marginValue := tokens.NewValueInt(0, t.Context())
	paddingValue := tokens.NewValueInt(0, t.Context())
	boxSizingValue := tokens.NewValueString("inherit", t.Context())

	allStyle.Set("margin", marginValue)
	allStyle.Set("padding", paddingValue)
	allStyle.Set("box-sizing", boxSizingValue)

	rule, err := styles.NewAllRule(t, allStyle)
	if err != nil {
		return err
	}

	ss.Append(rule)

	return nil
}

func (t *HTML) collectAutoLinkRule(ss styles.Sheet) error {
	if !AUTO_LINK {
		return nil
	}

	ctx := t.Context()
	aStyle := tokens.NewEmptyStringDict(ctx)

	textDecoKey, err := tokens.NewString("text-decoration", ctx)
	if err != nil {
		panic(err)
	}

	textDecoVal, err := tokens.NewString("none", ctx)
	if err != nil {
		panic(err)
	}

	displayKey, err := tokens.NewString("display", ctx)
	if err != nil {
		panic(err)
	}

	// flex seems to wrap the tightest
	displayVal, err := tokens.NewString("inline", ctx)
	if err != nil {
		panic(err)
	}

	aStyle.Set(textDecoKey, textDecoVal)
	aStyle.Set(displayKey, displayVal)

	// dummy tag
	aTag, err := NewVisibleTag("a", false, tokens.NewEmptyStringDict(ctx), ctx)
	if err != nil {
		panic(err)
	}

	rule, err := styles.NewTagRule(&aTag, aStyle)
	if err != nil {
		return err
	}

	ss.Append(rule)

	return nil
}

// return the bundleable css rules
func (t *HTML) CollectStyles(idMap IDMap, classMap ClassMap, cssUrl string) ([][]string, error) {
	head, body, err := t.getHeadBody()
	if err != nil {
		return nil, err
	}

	ss := styles.NewDocSheet()

	// * {margin: 0, padding 0}
	if err := t.collectAllRule(ss); err != nil {
		return nil, err
	}

	// a {text-decoration: none; etc...}
	if err := t.collectAutoLinkRule(ss); err != nil {
		return nil, err
	}

	if styleToken_, ok := t.attributes.Get("style"); ok {
		if !tokens.IsNull(styleToken_) {
			styleToken, err := tokens.AssertStringDict(styleToken_)
			if err != nil {
				context.AppendContextString(err, "Info: needed here", t.Context())
				return nil, err
			}

			if !styleToken.IsEmpty() {
				rule, err := styles.NewTagRule(t, styleToken)
				if err != nil {
					return nil, err
				}
				ss.Append(rule)
			}
			t.attributes.Delete("style") // null isn't printed anyway, so doesnt need to be deletd
		}
	}

	if err := body.CollectStyles(ss); err != nil {
		return nil, err
	}

	if err := ss.Synchronize(); err != nil {
		return nil, err
	}

	if !ss.IsEmpty() {
		ctx := t.Context()

		// resolve any nested rules
		ssExpanded, bundleRules, err := ss.ExpandNested()
		if err != nil {
			return nil, err
		}

    // if cssUrl == "" then css must be added by caller using the t.IncludeStyle() method
		if len(bundleRules) != 0  && cssUrl != "" {
			bundleLinkTag, err := NewStyleSheetLink(cssUrl, ctx)
			if err != nil {
				return nil, err
			}

			head.AppendChild(bundleLinkTag)
		}

		if !ssExpanded.IsEmpty() {
			styleTag, err := NewStyle(tokens.NewEmptyStringDict(ctx), ssExpanded.Write(), ctx)
			if err != nil {
				return nil, err
			}

			head.AppendChild(styleTag)
		}

		return styles.WriteBundleRules(bundleRules), nil
	} else {
		return [][]string{}, nil
	}
}

func (t *HTML) CollectScripts(idMap IDMap, classMap ClassMap, bundle *scripts.InlineBundle) error {
	// bundle is only used here, but HTML must implement Tag interface (to be a child of Root), so that's why bundle is passed in as an argument

	head, body, err := t.getHeadBody()
	if err != nil {
		return err
	}

	if err := head.CollectScripts(idMap, classMap, bundle); err != nil {
		return err
	}

	if err := body.CollectScripts(idMap, classMap, bundle); err != nil {
		return err
	}

	if !bundle.IsEmpty() {
		ctx := head.Context()

		deps := bundle.Dependencies()
		for _, dep := range deps {
			srcScript, _ := NewSrcScript(dep, ctx)

			head.AppendChild(srcScript)
		}

		loaderContent, err := bundle.Write()
		if err != nil {
			return err
		}

		loaderScript, err := NewLoaderScript(loaderContent, ctx)
		if err != nil {
			return err
		}

		// place script at end of body to decrease loading times
		body.AppendChild(loaderScript)
	}

	return nil
}

func (t *HTML) ApplyControl(control string, jsUrl string) error {
	head, body, err := t.getHeadBody()
	if err != nil {
		return err
	}

	ctx := head.Context()
	srcScript, _ := NewSrcScript(jsUrl, ctx)
	head.AppendChild(srcScript)

	// write the loader content
	var b strings.Builder
	b.WriteString(js.HashControl(control))
	b.WriteString("();")
	loaderContent := b.String()

	loaderScript, err := NewLoaderScript(loaderContent, ctx)
	if err != nil {
		return err
	}

	body.AppendChild(loaderScript)

	return nil
}

func (t *HTML) IncludeStyle(style string) error {
	head, _, err := t.getHeadBody()
  if err != nil {
    return err
  }

	ctx := head.Context()
  styleTag, err := NewStyle(tokens.NewEmptyStringDict(ctx), style, ctx)
  if err != nil {
    return err
  }

  head.AppendChild(styleTag)
  return nil
}

func (t *HTML) IncludeControl(code string) error {
	head, _, err := t.getHeadBody()
  if err != nil {
    return err
  }

	ctx := head.Context()
	script, _ := NewScript(tokens.NewEmptyStringDict(ctx), code, ctx)
	head.AppendChild(script)

  return nil
}

func (t *HTML) GetClasses() []string {
	return t.classes
}

func (t *HTML) SetClasses(cs []string) {
	t.classes = cs
}

func (t *HTML) Write(indent string, nl, tab string) string {
	hasClasses := len(t.classes) > 0

	if hasClasses {
		keyToken := tokens.NewValueString("class", t.Context())
		valueToken := tokens.NewValueString(strings.Join(t.classes, " "), t.Context())
		t.attributes.Set(keyToken, valueToken)
	}

	result := t.tagData.Write(indent, nl, tab)

	if hasClasses {
		t.attributes.Delete("class")
	}

	return result
}

func (t *HTML) InnerHTML() string {
	return ""
}

func (t *HTML) ToJSType() string {
	return "HTMLElement"
}
