package prototypes

import (
	//"../values"

	//"../../context"
)

var NodeJS_mysql_Query *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_mysql_QueryPrototype() bool {
	*NodeJS_mysql_Query = BuiltinPrototype{
		"mysql.Query", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
		},
		nil,
	}

	return true
}

var _NodeJS_mysql_QueryOk = generateNodeJS_mysql_QueryPrototype()
