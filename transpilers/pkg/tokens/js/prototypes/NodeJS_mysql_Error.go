package prototypes

import (
	//"../values"

	//"../../context"
)

var NodeJS_mysql_Error *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_mysql_ErrorPrototype() bool {
	*NodeJS_mysql_Error = BuiltinPrototype{
		"mysql.Error", Error,
		map[string]BuiltinFunction{
      "code": NewGetter(String),
      "errno": NewGetter(Int),
      "sqlMessage": NewGetter(String),
      "sqlState": NewGetter(String),
      "index": NewGetter(Int),
      "sql": NewGetter(String),
		},
		nil,
	}

	return true
}

var _NodeJS_mysql_ErrorOk = generateNodeJS_mysql_ErrorPrototype()
