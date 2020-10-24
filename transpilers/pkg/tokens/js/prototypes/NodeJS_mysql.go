package prototypes

var NodeJS_mysql *BuiltinPrototype = allocBuiltinPrototype()

// is actually a builtin nodejs module
func generateNodeJS_mysqlPrototype() bool {
	*NodeJS_mysql = BuiltinPrototype{
		"mysql", nil,
		map[string]BuiltinFunction{
			"createConnection": NewStatic(Object, NodeJS_mysql_Connection),
			"createPool": NewStatic(Object, NodeJS_mysql_Pool),
			"Connection":       NewStaticClassGetter(NodeJS_mysql_Connection),
      "Error":            NewStaticClassGetter(NodeJS_mysql_Error),
			"FieldPacket":      NewStaticClassGetter(NodeJS_mysql_FieldPacket),
			"Query":            NewStaticClassGetter(NodeJS_mysql_Query),
			"Pool":            NewStaticClassGetter(NodeJS_mysql_Pool),
		},
		nil,
	}

	return true
}

var _NodeJS_mysqlOk = generateNodeJS_mysqlPrototype()
