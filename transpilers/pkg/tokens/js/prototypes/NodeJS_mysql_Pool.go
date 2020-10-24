package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_mysql_Pool *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_mysql_PoolPrototype() bool {
	*NodeJS_mysql_Pool = BuiltinPrototype{
		"mysql.Pool", NodeJS_mysql_Connection,
		map[string]BuiltinFunction{
      "getConnection": NewNormalFunction(&Function{}, 
        func(stack values.Stack, this *values.Instance, args []values.Value,
          ctx context.Context) (values.Value, error) {
          if err := args[0].EvalMethod(stack, []values.Value{
            NewInstance(NodeJS_mysql_Error, ctx),
            NewInstance(NodeJS_mysql_Connection, ctx),
          }, ctx); err != nil {
            return nil, err
          }

          return nil, nil
        }),
		},
		nil,
	}

	return true
}

var _NodeJS_mysql_PoolOk = generateNodeJS_mysql_PoolPrototype()
