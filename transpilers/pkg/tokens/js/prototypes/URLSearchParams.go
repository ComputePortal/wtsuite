package prototypes

var URLSearchParams *BuiltinPrototype = allocBuiltinPrototype()

func generateURLSearchParamsPrototype() bool {
	*URLSearchParams = BuiltinPrototype{
		"URLSearchParams", nil,
		map[string]BuiltinFunction{
			"has": NewNormal(String, Boolean),
			"get": NewNormal(String, String),
		},
		nil,
	}

	return true
}

var _URLSearchParamsOk = generateURLSearchParamsPrototype()
