package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_http_IncomingMessage *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_http_IncomingMessagePrototype() bool {
	*NodeJS_http_IncomingMessage = BuiltinPrototype{
		"http.IncomingMessage", NodeJS_stream_Readable,
		map[string]BuiltinFunction{
			"aborted":     NewGetter(Boolean),
			"complete":    NewGetter(Boolean),
			"headers":     NewGetter(Object),
			"httpVersion": NewGetter(String),
			"method":      NewGetter(String),
			"rawHeaders": NewGetterFunction(func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				content := NewInstance(String, ctx)
				return NewArray([]values.Value{content}, ctx), nil
			}),
			"statusCode":    NewGetter(Int),
			"statusMessage": NewGetter(String),
			"url":           NewGetter(String),
		},
		nil,
	}

	return true
}

var _NodeJS_http_IncomingMessageOk = generateNodeJS_http_IncomingMessagePrototype()
