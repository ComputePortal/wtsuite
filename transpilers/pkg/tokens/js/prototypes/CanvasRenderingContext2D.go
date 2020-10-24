package prototypes

var CanvasRenderingContext2D *BuiltinPrototype = allocBuiltinPrototype()

func generateCanvasRenderingContext2D() bool {
	*CanvasRenderingContext2D = BuiltinPrototype{
		"CanvasRenderingContext2D", nil,
		map[string]BuiltinFunction{
			"arc":                     NewNormal(&And{&Many{5, Number}, &Opt{Boolean}}, nil),
			"arcTo":                   NewNormal(&Many{5, Number}, nil),
			"beginPath":               NewNormal(&None{}, nil),
			"bezierCurveTo":           NewNormal(&Many{6, Number}, nil),
			"clearRect":               NewNormal(&Many{4, Number}, nil),
			"clip":                    NewNormal(&Opt{String}, nil),
			"closePath":               NewNormal(&None{}, nil),
			"createLinearGradient":    NewNormal(&Many{4, Number}, CanvasGradient),
			"createRadialGradient":    NewNormal(&Many{6, Number}, CanvasGradient),
			"createPattern":           NewNormal(&And{&Any{}, String}, CanvasPattern),
			"direction":               NewSetter(String),
			"drawImage":               NewNormal(&And{&Or{HTMLImageElement, &Or{Image, HTMLCanvasElement}}, &And{&And{Number, Number}, &Opt{&And{&And{Number, Number}, &Opt{&Many{4, Number}}}}}}, nil),
			"ellipse":                 NewNormal(&And{&Many{7, Number}, &Opt{Boolean}}, nil),
			"fill":                    NewNormal(&Opt{String}, nil),
			"fillRect":                NewNormal(&Many{4, Number}, nil),
			"fillStyle":               NewSetter(&Or{String, &Or{CanvasGradient, CanvasPattern}}),
			"fillText":                NewNormal(&And{String, &And{&Many{2, Number}, &Opt{Number}}}, nil),
			"font":                    NewSetter(String),
			"getImageData":            NewNormal(&Many{4, Number}, ImageData),
			"getTransform":            NewNormal(&None{}, DOMMatrix),
			"globalAlpha":             NewSetter(Number),
			"globalCompositeOperator": NewSetter(String),
			"isPointInPath":           NewNormal(&And{Number, &And{Number, &Opt{String}}}, Boolean),
			"isPointInStroke":         NewNormal(&And{Number, Number}, Boolean),
			"lineCap":                 NewSetter(String),
			"lineJoin":                NewSetter(String),
			"lineTo":                  NewNormal(&And{Number, Number}, nil),
			"lineWidth":               NewSetter(Number),
			"measureText":             NewNormal(String, TextMetrics),
			"miterLimit":              NewSetter(Number),
			"moveTo":                  NewNormal(&And{Number, Number}, nil),
			"putImageData":            NewNormal(&And{ImageData, &Or{&Many{2, Number}, &Many{6, Number}}}, nil),
			"quadraticCurveTo":        NewNormal(&Many{4, Number}, nil),
			"rect":                    NewNormal(&Many{4, Number}, nil),
			"restore":                 NewNormal(&None{}, nil),
			"rotate":                  NewNormal(Number, nil),
			"save":                    NewNormal(&None{}, nil),
			"scale":                   NewNormal(&And{Number, Number}, nil),
			"setTransform":            NewNormal(&Or{&Many{6, Number}, DOMMatrix}, nil),
			"shadowBlir":              NewSetter(Number),
			"shadowColor":             NewSetter(String),
			"shadowOffsetX":           NewSetter(Number),
			"shadowOffsetY":           NewSetter(Number),
			"stroke":                  NewNormal(&None{}, nil),
			"strokeRect":              NewNormal(&Many{4, Number}, nil),
			"strokeStyle":             NewSetter(&Or{String, &Or{CanvasGradient, CanvasPattern}}),
			"strokeText":              NewNormal(&And{String, &And{&Many{2, Number}, &Opt{Number}}}, nil),
			"textAlign":               NewSetter(String),
			"textBaseline":            NewSetter(String),
			"translate":               NewNormal(&Many{6, Number}, nil),
		},
		nil,
	}

	return true
}

var _CanvasRenderingContext2DOk = generateCanvasRenderingContext2D()
