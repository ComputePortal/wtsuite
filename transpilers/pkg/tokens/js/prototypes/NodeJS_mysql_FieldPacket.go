package prototypes

import (
	//"../values"

	//"../../context"
)

var NodeJS_mysql_FieldPacket *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_mysql_FieldPacketPrototype() bool {
	*NodeJS_mysql_FieldPacket = BuiltinPrototype{
		"mysql.FieldPacket", nil,
		map[string]BuiltinFunction{
      "catalog": NewGetter(String),
      "db": NewGetter(String),
      "table": NewGetter(String),
      "name": NewGetter(String),
      "length": NewGetter(Int),
      "type": NewGetter(Int),
      "flags": NewGetter(Int),
		},
		nil,
	}

	return true
}

var _NodeJS_mysql_FieldPacketOk = generateNodeJS_mysql_FieldPacketPrototype()
