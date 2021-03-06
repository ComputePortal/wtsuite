package directives

import (
	"fmt"
	"os"

	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/svg"
)

var MATH_FONT = "FreeSerif"
var MATH_FONT_FAMILY = "FreeSerif"
var MATH_FONT_URL = ""

func Math(scope Scope, node Node, tag *tokens.Tag) error {
	attrScope := NewSubScope(scope)

  if err := tag.AssertEmpty(); err != nil {
    return err
  }

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

  var x *tokens.Float = nil
  var y *tokens.Float = nil
  if x_, ok := attr.Get("x"); ok {
    x, err = tokens.AssertIntOrFloat(x_)
    if err != nil {
      return err
    }

    attr.Delete("x")
  }

  if y_, ok := attr.Get("y"); ok {
    y, err = tokens.AssertIntOrFloat(y_)
    if err != nil {
      return err
    }

    attr.Delete("y")
  }

  if isInSVG {
    if duplicate_, ok := attr.Get("duplicate"); ok && !tokens.IsNull(duplicate_) {
      if err := tokens.AssertFlag(duplicate_); err != nil {
        return err
      }

      attr.Delete("duplicate")

      svgTag, err := buildMathSVGTag(node, valueToken, attrScope, attr, x, y, isInSVG, isInline, ctx)
      if err != nil {
        return err
      }

      if err := node.AppendChild(svgTag); err != nil {
        return nil
      }

      svgTag.Attributes().Set("class", tokens.NewValueString("math duplicate", ctx))
    }
  }

  svgTag, err := buildMathSVGTag(node, valueToken, attrScope, attr, x, y, isInSVG, isInline, ctx)
  if err != nil {
    return err
  }

	if err := node.AppendChild(svgTag); err != nil {
		return err
	}

  return nil
}

func buildMathSVGTag(parentNode Node, valueToken *tokens.String, attrScope Scope, attr *tokens.StringDict, x, y *tokens.Float, isInSVG bool, isInline bool, ctx context.Context) (tree.Tag, error) {
	svgAttr := tokens.NewEmptyStringDict(ctx) // filled later, depends on BB
	svgTag, err := tree.BuildTag("svg", svgAttr, ctx)
	if err != nil {
		return nil, err
	}

	mathParser, err := parsers.NewMathParser(valueToken.Value(), valueToken.InnerContext())
	if err != nil {
		return nil, err
	}

	mt, err := mathParser.Build()
	if err != nil {
		return nil, err
	}

	mNode := NewMathNode(svgTag, parentNode)

	totalBB, err := mt.GenerateTags(mNode, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	// fill the svg attributes
	svgAttr.Set("overflow", tokens.NewValueString("visible", ctx))
	svgAttr.Set("class", tokens.NewValueString("math", ctx))

	//styleValue := tokens.NewEmptyStringDict(ctx)
	//styleValue.Set("font-family", tokens.NewValueString(MATH_FONT_FAMILY, ctx))

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

			svgAttr.Set("height", heightVal)
			svgAttr.Set("width", widthVal)
		} else {
			viewBoxValue := tokens.NewValueString(fmt.Sprintf("%g %g %g %g",
				totalBB.Left(), totalBB.Top(), totalBB.Width(), totalBB.Height()),
				ctx)
			svgAttr.Set("viewBox", viewBoxValue)

			heightVal := tokens.NewValueUnitFloat(totalBB.Height(), "em", ctx)

			widthVal := tokens.NewValueUnitFloat(totalBB.Width(), "em", ctx)

			svgAttr.Set("height", heightVal)
			svgAttr.Set("width", widthVal)
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
				return nil, errCtx.NewError("Error: must specifiy either height or width or font-size when including math in an svg")
			}
		} else if hasFontSize {
			warningCtx := attr.Context()
			fmt.Fprintf(os.Stderr, "%s\n", warningCtx.NewError("Warning: font-size ignored in favour of height/width").Error())
		}

		if hasHeight {
			h, err := tokens.AssertIntOrFloat(heightToken_)
			if err != nil {
				return nil, err
			}

			inputHeight = h.Value()
			if inputHeight <= 0.0 {
				errCtx := h.Context()
				return nil, errCtx.NewError("Error: non-positive input height")
			}
		}

		if hasWidth {
			w, err := tokens.AssertIntOrFloat(widthToken_)
			if err != nil {
				return nil, err
			}

			inputWidth = w.Value()

			if inputWidth <= 0.0 {
				errCtx := w.Context()
				return nil, errCtx.NewError("Error: non-positive input width")
			}
		}

		if hasFontSize {
			fs, err := tokens.AssertIntOrFloat(fontSizeToken_)
			if err != nil {
				return nil, err
			}

			inputHeight = fs.Value() * totalBB.Height()
			if inputHeight <= 0 {
				errCtx := fs.Context()
				return nil, errCtx.NewError("Error: non-positive input font-size")
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
			return nil, err
		}

		if x != nil {
      x = tokens.NewValueFloat(
        x.Value()-
          0.5*resultWidth*float64(1-horAnchor)+
          float64(horAnchor)*anchorOffset,
        x.Context(),
      )

			svgAttr.Set("x", x)
		}

		if y != nil {
      y = tokens.NewValueFloat(
        y.Value()-
          0.5*resultHeight*float64(1-verAnchor)+
          float64(verAnchor)*anchorOffset,
        y.Context(),
      )

			svgAttr.Set("y", y)
		}
	}

	// merge using input attributes

	if err := functions.MergeStringDictsInplace(attrScope, svgAttr, attr, ctx); err != nil {
		return nil, err
	}

	return svgTag, nil
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
func evalMathURI(scope Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := functions.CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	mathAttr := tokens.NewEmptyStringDict(ctx)
	mathAttr.Set("value", args[0])

	uriNode := NewURINode()

	mathTag := tokens.NewTag("math", mathAttr.ToRaw(), []*tokens.Tag{}, ctx)
	if err := Math(scope, uriNode, mathTag); err != nil {
		return nil, err
	}

	tag := uriNode.tag

	// XXX: data-uri svg's with @font-face styles are not actually supported
  //  so the following doesnt work in any browser (maybe in future ...)
	if MATH_FONT_URL != "" {
		// add style child for math font import
		defs, err := svg.BuildTag("defs", tokens.NewEmptyStringDict(ctx), ctx)
		if err != nil {
			panic(err)
		}
		importFontStyle, err := tree.NewStyle(tokens.NewEmptyStringDict(ctx),
			"@font-face{font-family:"+MATH_FONT+";src:url("+MATH_FONT_URL+");}",
			ctx)
		defs.AppendChild(importFontStyle)

		textTag := tag.Children()[0]
		textAttr := textTag.Attributes()
		textAttr.Set("font-family", tokens.NewValueString(MATH_FONT_FAMILY, ctx))

		tag.InsertChild(0, defs)
	}

	return svgToURI(tag, ctx)
}

var _mathOk = registerDirective("math", Math)
