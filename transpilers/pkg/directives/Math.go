package directives

import (
	"fmt"
	"os"

	"../functions"
	"../parsers"
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
	"../tree/styles"
	"../tree/svg"
)

func Math(scope Scope, node Node, tag *tokens.Tag) error {
	attrScope := NewSubScope(scope, node)

	// eval the incoming attr
	attr, err := tag.Attributes([]string{"value", "inline"}) // inline defaults to true
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(attrScope)
	if err != nil {
		return err
	}

	valueToken, err := tokens.DictString(attr, "value")
	if err != nil {
		return err
	}

	attr.Delete("value")

	isInline := true
	if inlineToken_, ok := attr.Get("inline"); ok {
		inlineToken, err := tokens.AssertBool(inlineToken_)
		if err != nil {
			return err
		}

		isInline = inlineToken.Value()
	}

	ctx := tag.Context()

	isInSVG := node.Type() == SVG
	svgAttr := tokens.NewEmptyStringDict(ctx) // filled later, depends on BB
	svgTag, err := tree.BuildTag("svg", svgAttr, ctx)
	if err != nil {
		return err
	}
	if err := node.AppendChild(svgTag); err != nil {
		return nil
	}

	mathParser, err := parsers.NewMathParser(valueToken.Value(), valueToken.InnerContext())
	if err != nil {
		return err
	}

	mt, err := mathParser.Build()
	if err != nil {
		return err
	}

	var mStyle *tokens.StringDict = nil
	if styleToken_, ok := attr.Get("style"); ok {
		styleToken, err := tokens.AssertStringDict(styleToken_)
		if err != nil {
			return err
		}

		mStyle = styleToken
	}
	mNode := NewMathNode(svgTag, node, mStyle)

	totalBB, err := mt.GenerateTags(mNode, 0.0, 0.0)
	if err != nil {
		return err
	}

	// fill the svg attributes
	svgAttr.Set("overflow", tokens.NewValueString("visible", ctx))

	styleValue := tokens.NewEmptyStringDict(ctx)
	styleValue.Set("font-family", tokens.NewValueString(styles.MATH_FONT_FAMILY, ctx))

	if !isInSVG {
		if isInline {
			paddingLeft := 0.15 // based on typical advance width
			paddingRight := 0.15

			inlineHeight := 1.0
			viewBoxValue := tokens.NewValueString(fmt.Sprintf("%g %g %g %g",
				totalBB.Left()-paddingLeft,
				-inlineHeight,
				totalBB.Width()+paddingLeft+paddingRight,
				inlineHeight), ctx)
			svgAttr.Set("viewBox", viewBoxValue)

			heightVal := tokens.NewValueUnitFloat(inlineHeight*1.0, "em", ctx)

			widthVal := tokens.NewValueUnitFloat((totalBB.Width()+paddingLeft+paddingRight)*1.0, "em", ctx)

			styleValue.Set("height", heightVal)
			styleValue.Set("width", widthVal)
		} else {
			viewBoxValue := tokens.NewValueString(fmt.Sprintf("%g %g %g %g",
				totalBB.Left(), totalBB.Top(), totalBB.Width(), totalBB.Height()),
				ctx)
			svgAttr.Set("viewBox", viewBoxValue)

			heightVal := tokens.NewValueUnitFloat(totalBB.Height(), "em", ctx)

			widthVal := tokens.NewValueUnitFloat(totalBB.Width(), "em", ctx)

			styleValue.Set("height", heightVal)
			styleValue.Set("width", widthVal)
		}

		// search for the color in incoming style
		// TODO: should math inside svg use the fills directly?
		colorValue, err := SearchStyle(scope, attr, "color", ctx)
		if err != nil {
			return err
		}

		if !tokens.IsNull(colorValue) {
			styleValue.Set("fill", colorValue)
		}
	} else {
		viewBoxValue := tokens.NewValueString(fmt.Sprintf("%g %g %g %g",
			totalBB.Left(), totalBB.Top(), totalBB.Width(), totalBB.Height()),
			ctx)
		svgAttr.Set("viewBox", viewBoxValue)

		inputHeight := -1.0
		inputWidth := -1.0

		heightToken_, hasHeight := attr.Get("height")
		widthToken_, hasWidth := attr.Get("width")
		fontSizeToken_, hasFontSize := attr.Get("font-size")
		if !hasHeight && !hasWidth {
			if !hasFontSize {
				errCtx := attr.Context()
				return errCtx.NewError("Error: must specifiy either height or width or font-size when including math in an svg")
			}
		} else if hasFontSize {
			warningCtx := attr.Context()
			fmt.Fprintf(os.Stderr, "%s\n", warningCtx.NewError("Warning: font-size ignored in favour of height/width").Error())
		}

		if hasHeight {
			h, err := tokens.AssertIntOrFloat(heightToken_)
			if err != nil {
				return err
			}

			inputHeight = h.Value()
			if inputHeight <= 0.0 {
				errCtx := h.Context()
				return errCtx.NewError("Error: non-positive input height")
			}
		}

		if hasWidth {
			w, err := tokens.AssertIntOrFloat(widthToken_)
			if err != nil {
				return err
			}

			inputWidth = w.Value()

			if inputWidth <= 0.0 {
				errCtx := w.Context()
				return errCtx.NewError("Error: non-positive input width")
			}
		}

		if hasFontSize {
			fs, err := tokens.AssertIntOrFloat(fontSizeToken_)
			if err != nil {
				return err
			}

			inputHeight = fs.Value() * totalBB.Height()
			if inputHeight <= 0 {
				errCtx := fs.Context()
				return errCtx.NewError("Error: non-positive input font-size")
			}
		}

		resultHeight := -1.0
		resultWidth := -1.0
		if inputHeight > 0.0 {
			resultHeight = inputHeight
			resultWidth = totalBB.Width() / totalBB.Height() * inputHeight
		}

		if inputWidth > 0.0 {
			if (resultWidth > 0.0 && inputWidth < resultWidth) || resultWidth < 0.0 {
				resultWidth = inputWidth
				resultHeight = totalBB.Height() / totalBB.Width() * inputWidth
			}
		}

		heightVal := tokens.NewValueFloat(resultHeight, ctx)
		widthVal := tokens.NewValueFloat(resultWidth, ctx)
		svgAttr.Set("height", heightVal)
		svgAttr.Set("width", widthVal)

		// anchors are only relevant in an svg
		horAnchor, verAnchor, anchorOffset, err := parseMathAnchors(attr)
		if err != nil {
			return err
		}

		if x_, ok := attr.Get("x"); ok {
			x, err := tokens.AssertIntOrFloat(x_)
			if err != nil {
				return err
			}

			x = tokens.NewValueFloat(
				x.Value()-
					0.5*resultWidth*float64(1-horAnchor)+
					float64(horAnchor)*anchorOffset,
				x.Context(),
			)
			attr.Set("x", x)
		}

		if y_, ok := attr.Get("y"); ok {
			y, err := tokens.AssertIntOrFloat(y_)
			if err != nil {
				return err
			}

			y = tokens.NewValueFloat(
				y.Value()-
					0.5*resultHeight*float64(1-verAnchor)+
					float64(verAnchor)*anchorOffset,
				y.Context(),
			)
			attr.Set("y", y)
		}
	}

	svgAttr.Set("style", styleValue)

	// merge using input attributes

	if err := functions.MergeStringDictsInplace(svgAttr, attr, ctx); err != nil {
		return err
	}

	return nil
}

// returned anchorOffset is one leg (hor or ver) of manhatten distance, not euclidean distance
func parseMathAnchors(attr *tokens.StringDict) (int, int, float64, error) {
	horAnchor := 1 // -1/0/1
	verAnchor := -1
	if anchorToken_, ok := attr.Get("anchor"); ok {
		anchorToken, err := tokens.AssertString(anchorToken_)
		if err != nil {
			return 0, 0, 0.0, err
		}

		str := anchorToken.Value()
		if len(str) != 2 {
			errCtx := anchorToken.Context()
			return 0, 0, 0.0, errCtx.NewError("Error: expected two characters (eg. cc)")
		}

		horChar := str[0:1]
		verChar := str[1:2]

		switch horChar {
		case "c":
			horAnchor = 0
		case "l":
			horAnchor = -1
		case "r":
			horAnchor = 1
		default:
			errCtx := anchorToken.Context()
			return 0, 0, 0.0, errCtx.NewError("Error: expected c/l/r for first char, got " + horChar)
		}

		switch verChar {
		case "c":
			verAnchor = 0
		case "t":
			verAnchor = -1
		case "b":
			verAnchor = 1
		default:
			errCtx := anchorToken.Context()
			return 0, 0, 0.0, errCtx.NewError("Error: expected c/t/b for second char, got " + verChar)
		}
	}

	anchorOffset := 0.0
	if anchorOffsetToken_, ok := attr.Get("anchor-offset"); ok {
		anchorOffsetToken, err := tokens.AssertIntOrFloat(anchorOffsetToken_)
		if err != nil {
			return 0, 0, 0.0, err
		}

		anchorOffset = anchorOffsetToken.Value()
	}

	return horAnchor, verAnchor, anchorOffset, nil
}

// assume it is used for inline, wrap around
func evalMathURI(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	mathAttr := tokens.NewEmptyStringDict(ctx)
	mathAttr.Set("value", args[0])

	uriNode := NewURINode(scope.GetNode())

	mathTag := tokens.NewTag("math", mathAttr.ToRaw(), []*tokens.Tag{}, ctx)
	if err := Math(scope, uriNode, mathTag); err != nil {
		return nil, err
	}

	tag := uriNode.tag

	// XXX: data-uri svg's with @font-face styles are not actually supported
	if styles.MATH_FONT_URL != "" {
		// add style child for math font import
		defs, err := svg.BuildTag("defs", tokens.NewEmptyStringDict(ctx), ctx)
		if err != nil {
			panic(err)
		}
		importFontStyle, err := tree.NewStyle(tokens.NewEmptyStringDict(ctx),
			"@font-face{font-family:"+styles.MATH_FONT+";src:url("+styles.MATH_FONT_URL+");}",
			ctx)
		defs.AppendChild(importFontStyle)

		textTag := tag.Children()[0]
		textAttr := textTag.Attributes()
		textAttr.Set("font-family", tokens.NewValueString(styles.MATH_FONT_FAMILY, ctx))

		tag.InsertChild(0, defs)
	}

	return svgToURI(tag, ctx)
}

var _mathOk = registerDirective("math", Math)
