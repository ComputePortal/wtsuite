package functions

import (
	"math"

	"../tokens/context"
	tokens "../tokens/html"
)

func floatToFloatMath(args []tokens.Token, fn func(val float64) float64, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	x, err := tokens.AssertAnyIntOrFloat(args[0])
	if err != nil {
		return nil, err
	}

	return tokens.NewValueFloat(fn(x.Value()), ctx), nil
}

func twoFloatsToFloatMath(args []tokens.Token, fn func(x float64, y float64) float64, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	x, err := tokens.AssertAnyIntOrFloat(args[0])
	if err != nil {
		return nil, err
	}

	y, err := tokens.AssertAnyIntOrFloat(args[1])
	if err != nil {
		return nil, err
	}

	return tokens.NewValueFloat(fn(x.Value(), y.Value()), ctx), nil
}

func Sqrt(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return floatToFloatMath(args, math.Sqrt, ctx)
}

func Sin(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return floatToFloatMath(args, math.Sin, ctx)
}

func Cos(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return floatToFloatMath(args, math.Cos, ctx)
}

func Tan(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return floatToFloatMath(args, math.Tan, ctx)
}

func Rad(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return floatToFloatMath(args, func(val float64) float64 {
		return val * math.Pi / 180.0
	}, ctx)
}

func Pow(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return twoFloatsToFloatMath(args, math.Pow, ctx)
}

func Pi(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: unexpected arguments")
	}

	return tokens.NewValueFloat(math.Pi, ctx), nil
}

func round(args []tokens.Token, fn func(val float64) float64, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	switch {
	case tokens.IsFloat(args[0]):
		fl, err := tokens.AssertAnyIntOrFloat(args[0])
		if err != nil {
			panic(err)
		}

		val := fn(fl.Value())

		if fl.Unit() == "" {
			return tokens.NewValueInt(int(val), ctx), nil
		} else {
			return tokens.NewValueUnitFloat(val, fl.Unit(), ctx), nil
		}
	case tokens.IsInt(args[0]):
		i, err := tokens.AssertInt(args[0])
		if err != nil {
			panic(err)
		}

		return tokens.NewValueInt(i.Value(), ctx), nil
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected int or float as argument")
	}
}

func Round(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return round(args, math.Round, ctx)
}

func Floor(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return round(args, math.Floor, ctx)
}

func Ceil(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return round(args, math.Ceil, ctx)
}
