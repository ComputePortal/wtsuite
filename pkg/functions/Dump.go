package functions

import (
	"fmt"

	"../tokens/context"
	tokens "../tokens/html"
)

func Dump(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	infoErr := ctx.NewError("Info: variable dump")

	doDump := true
	i := 0
	if len(args) == 2 {
		// first should be bool
		b, err := tokens.AssertBool(args[0])
		if err != nil {
			return nil, err
		}

		doDump = b.Value()
		i = 1
	} else if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}

	if doDump {
		fmt.Printf(infoErr.Error())
		fmt.Println(args[i].Dump("#DUMP: "))
	}

	return args[i], nil
}
