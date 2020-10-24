package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_nodemailer_SMTPTransport *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_nodemailer_SMTPTransportPrototype() bool {
	*NodeJS_nodemailer_SMTPTransport = BuiltinPrototype{
		"nodemailer.SMTPTransport", nil,
		map[string]BuiltinFunction{
      "sendMail": NewNormalFunction(Object,
        func(stack values.Stack, this *values.Instance,
        args []values.Value, ctx context.Context) (values.Value, error) {
          return NewVoidPromise(ctx)
        }),
		},
		nil,
	}

	return true
}

var _NodeJS_nodemailer_SMTPTransportOk = generateNodeJS_nodemailer_SMTPTransportPrototype()
