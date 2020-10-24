package prototypes

var TextMetrics *BuiltinPrototype = allocBuiltinPrototype()

func generateTextMetricsPrototype() bool {
	*TextMetrics = BuiltinPrototype{
		"TextMetrics", nil,
		map[string]BuiltinFunction{
			"width": NewGetter(Number),

			// XXX: not yet supported by all major browsers
			/*
				"actualBoundingBoxLeft": NewGetter(Number),
				"actualBoundingBoxRight": NewGetter(Number),
				"fontBoundingBoxAscent": NewGetter(Number),
				"fontBoundingBoxDescent": NewGetter(Number),
				"actualBoundingBoxAscent": NewGetter(Number),
				"actualBoundingBoxDescent": NewGetter(Number),
				"emHeightAscent": NewGetter(Number),
				"emHeightDescent": NewGetter(Number),
				"hangingBaseline": NewGetter(Number),
				"alphabeticBaseline": NewGetter(Number),
				"ideagraphicBaseline": NewGetter(Number),
			*/
		},
		nil,
	}

	return true
}

var _TextMetricsOk = generateTextMetricsPrototype()
