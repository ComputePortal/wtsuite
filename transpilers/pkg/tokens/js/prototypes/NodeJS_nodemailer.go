package prototypes

import (
	//"../values"

	//"../../context"
)

var NodeJS_nodemailer *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_nodemailerprototype() bool {
	*NodeJS_nodemailer = BuiltinPrototype{
		"nodemailer", nil,
		map[string]BuiltinFunction{
      "createTransport": NewStatic(Object, NodeJS_nodemailer_SMTPTransport),
      "SMTPTransport": NewStaticClassGetter(NodeJS_nodemailer_SMTPTransport),
		},
		nil,
	}

	return true
}

var _NodeJS_nodemailerOk = generateNodeJS_nodemailerprototype()
