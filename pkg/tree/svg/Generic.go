package svg

import (
	"strconv"

	"../../tree"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type Generic struct {
	tree.SVGTagData
}

func NewGeneric(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData(key, attr, ctx)
	return &Generic{svgTag}, err
}

func AssertFloatAttribute(attr *tokens.StringDict, key string) error {
	if val_, ok := attr.Get(key); ok {
		var val *tokens.Float = nil
		if tokens.IsString(val_) {
			// convert to Float first
			valStr, err := tokens.AssertString(val_)
			if err != nil {
				panic(err)
			}

			fl, err := strconv.ParseFloat(valStr.Value(), 64)
			if err != nil {
				errCtx := val_.Context()
				return errCtx.NewError("Error: unable to parse as float")
			}

			val = tokens.NewValueFloat(fl, val_.Context())
		} else {
			var err error
			val, err = tokens.AssertIntOrFloat(val_)
			if err != nil {
				return err
			}
		}

		attr.Set(key, val)
	}

	return nil
}

func (t *Generic) Validate() error {
	attr := t.Attributes()

	if err := AssertFloatAttribute(attr, "x"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "y"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "cx"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "cy"); err != nil {
		return err
	}

	return t.SVGTagData.Validate()
}

func CompressFloatAttribute(attr *tokens.StringDict, key string,
	fn func(float64) float64) {
	if val_, ok := attr.Get(key); ok {
		var val *tokens.Float = nil
		if tokens.IsString(val_) {
			// convert to Float first
			valStr, err := tokens.AssertString(val_)
			if err != nil {
				panic(err)
			}

			fl, err := strconv.ParseFloat(valStr.Value(), 64)
			if err != nil {
				panic("unable to parse as float")
			}

			val = tokens.NewValueFloat(fl, val_.Context())
		} else {
			var err error
			val, err = tokens.AssertIntOrFloat(val_)
			if err != nil {
				panic("should've been caught before")
			}
		}

		val = tokens.NewValueFloat(fn(val.Value()), val.Context())
		attr.Set(key, val)
	}
}

func CompressStyles(attr *tokens.StringDict, vb tree.SVGViewBox) {
	if style_, ok := attr.Get("style"); ok && tokens.IsStringDict(style_) {
		style, err := tokens.AssertStringDict(style_)
		if err != nil {
			panic(err)
		}

		CompressFloatAttribute(style, "stroke-width", vb.CompressScalar)
	}
}

func (t *Generic) Compress(vb tree.SVGViewBox) {
	attr := t.Attributes()

	CompressFloatAttribute(attr, "x", vb.CompressX)
	CompressFloatAttribute(attr, "y", vb.CompressY)
	CompressFloatAttribute(attr, "cx", vb.CompressX)
	CompressFloatAttribute(attr, "cy", vb.CompressY)

	CompressStyles(attr, vb)

	t.SVGTagData.Compress(vb)
}

func (t *Generic) Minify() bool {
	attr := t.Attributes()

	attr.Delete("aria-label")
	attr.Delete("inkscape:label")
	attr.Delete("inkscape:groupmode")

	t.SetID("") // probably not used

	return t.SVGTagData.Minify()
}

func MinifyStyles(attr *tokens.StringDict, containsText bool) {
	if style_, ok := attr.Get("style"); ok && tokens.IsStringDict(style_) {
		style, err := tokens.AssertStringDict(style_)
		if err != nil {
			panic(err)
		}

		if !containsText {
			style.Delete("font-style")
			style.Delete("font-variant")
			style.Delete("font-weight")
			style.Delete("font-stretch")
			style.Delete("font-size")
			style.Delete("font-family")
			style.Delete("-inkscape-font-specification")
			style.Delete("word-spacing")
			style.Delete("letter-spacing")
			style.Delete("line-height")
		}

		// XXX: can we always do this?
		style.Delete("display")

		if stroke_, ok := style.Get("stroke"); ok && tokens.IsString(stroke_) {
			stroke, err := tokens.AssertString(stroke_)
			if err != nil {
				panic(err)
			}

			if stroke.Value() == "none" {
				style.Delete("stroke-linecap")
				style.Delete("stroke-miterlimit")
				style.Delete("stroke-dasharray")
				style.Delete("stroke-opacity")
				style.Delete("stroke-width")
			}
		}

		if fill_, ok := style.Get("fill"); ok && tokens.IsString(fill_) {
			fill, err := tokens.AssertString(fill_)
			if err != nil {
				panic(err)
			}

			if fill.Value() == "none" {
				style.Delete("fill-opacity")
			}
		}

		opacityThreshhold := 0.95
		// delete opacities very close to 1
		filterOpacity := func(key string) {
			if val_, ok := style.Get(key); ok && tokens.IsIntOrFloat(val_) {
				val, err := tokens.AssertAnyIntOrFloat(val_)
				if err != nil {
					panic(err)
				}

				if val.Value() > opacityThreshhold {
					style.Delete(key)
				}
			}
		}

		filterOpacity("fill-opacity")
		filterOpacity("opacity")
	}
}
