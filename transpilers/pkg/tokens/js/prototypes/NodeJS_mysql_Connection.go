package prototypes

import (
	"../values"

	"../../context"

  "strings"
)

var NodeJS_mysql_Connection *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_mysql_ConnectionPrototype() bool {
	*NodeJS_mysql_Connection = BuiltinPrototype{
		"mysql.Connection", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
      "connect": NewNormal(&None{}, nil),
      "end": NewNormalFunction(&Function{}, 
        func(stack values.Stack, this *values.Instance, args []values.Value,
          ctx context.Context) (values.Value, error) {
          if err := args[0].EvalMethod(stack, []values.Value{NewInstance(NodeJS_mysql_Error, ctx)}, ctx); err != nil {
            return nil, err
          }

          return nil, nil
        }),
      "query": NewMethodLikeNormalFunction(&And{&Or{Object, String}, &And{&Opt{Array}, &Function{}}},
        func(stack values.Stack, this *values.Instance, args []values.Value,
          ctx context.Context) (values.Value, error) {
          // try and get a the literal query string
          queryStr := ""
          if args[0].IsInstanceOf(String) {
            queryStr_, ok := args[0].LiteralStringValue()
            if ok {
              queryStr = queryStr_
            }
          } else {
            // must be object
            if queryStrVal, err := args[0].GetMember(stack, "sql", true, ctx); err == nil {
              queryStr_, ok := queryStrVal.LiteralStringValue()
              if ok {
                queryStr = queryStr_
              }
            }
          }

          callbackArg0 := NewInstance(NodeJS_mysql_Error, ctx)
          var callbackArg1 values.Value
          queryStr = strings.ToLower(queryStr)
          if strings.HasPrefix(queryStr, "insert") || strings.HasPrefix(queryStr, "create") {
            callbackArg1 = NewObject(nil, ctx)
          } else {
            callbackArg1 = NewArray([]values.Value{NewObject(nil, ctx)}, ctx)
          }

          callbackArg2 := NewArray([]values.Value{NewInstance(NodeJS_mysql_FieldPacket, ctx)}, ctx)

          callback := args[1]
          if len(args) == 3 {
            callback = args[2]
          }

          if err := callback.EvalMethod(stack, []values.Value{callbackArg0, callbackArg1, callbackArg2}, callback.Context()); err != nil {
            return nil, err
          }

          return NewInstance(NodeJS_mysql_Query, ctx), nil
        }),
      "release": NewNormal(&None{}, nil),
		},
		nil,
	}

	return true
}

var _NodeJS_mysql_ConnectionOk = generateNodeJS_mysql_ConnectionPrototype()
