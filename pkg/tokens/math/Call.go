package math

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"./boundingbox"

	"../context"
)

const (
	extraCallSpacing = 0.05
)

type Call struct {
	name Token
	args Token
	TokenData
}

func NewCall(name string, args []Token, ctx context.Context) (Token, error) {
	switch {
	case name == "text":
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected a single argument to the text() builtin function")
		}

		switch arg := args[0].(type) {
		case *Symbol:
			// convert to word
			return NewValueWord(arg.Value(), arg.Context())
		case *Word:
			return arg, nil
		default:
			errCtx := args[0].Context()
			return nil, errCtx.NewError("Error: expected word or symbol, cant convert complex tokens to text")
		}
	case name == "diff":
		switch len(args) {
		case 1:
			panic("not yet implemented (requires apostrophe)")
		case 2:
			return newFrac(
				newMul(newSymbol("d", ctx), args[0], ctx),
				newMul(newSymbol("d", ctx), args[1], ctx),
				ctx,
			), nil
		case 3:
			if args[1].Dump("") == args[2].Dump("") {
				return newFrac(
					newMul(newPow(newSymbol("d", ctx), newWord("2", ctx), ctx), args[0], ctx),
					newMul(newSymbol("d", ctx), newPow(args[1], newWord("2", ctx), ctx), ctx),
					ctx), nil
			} else {
				return newFrac(
					newMul(newPow(newSymbol("d", ctx), newWord("2", ctx), ctx), args[0], ctx),
					newMul(
						newMul(newSymbol("d", ctx), args[1], ctx),
						newMul(newSymbol("d", ctx), args[2], ctx),
						ctx),
					ctx), nil
			}
		default:
			return nil, ctx.NewError("Error: expected 1 or 2 arguments to the diff() builtin function")
		}
	case name == "sqrt":
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected a single argument to the sqrt() builtin function")
		}

		return NewRoot(args[0], nil, ctx)
	case name == "dot":
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected a single argument to the dot() builtin function")
		}

		return NewDot(args[0], ctx)
	case name == "align": // intended for debugging
		if len(args) < 2 {
			return nil, ctx.NewError("Error: expected at least 1 arguments")
		}

		last_ := args[len(args)-1]

		f, ok := last_.(*Float)
		if !ok {
			fmt.Println(reflect.TypeOf(last_).String())
			return nil, ctx.NewError("Error: expected literal int for last argument")
		}

		n := int(f.v)

		if n < 1 {
			return nil, ctx.NewError("Error: literal int must be >= 1")
		}

		if (len(args)-1)%n != 0 {
			return nil, ctx.NewError("Error: len(args)-1 not divisible by " + strconv.Itoa(n))
		}

		eqs := make([][]Token, 0)

		for i := 0; i < len(args)-1; i += n {
			row := []Token{}

			for j := 0; j < n; j++ {
				row = append(row, args[i+j])
			}

			eqs = append(eqs, row)
		}

		return NewAlign(genericMinHorSpacing, genericMinVerSpacing, false, true, eqs, ctx)
	case name == "ifelse":
		if len(args) < 2 {
			return nil, ctx.NewError("Error: expected at least 2 arguments")
		}

		if len(args)%2 != 0 {
			return nil, ctx.NewError("Error: expected an even number of arguments")
		}

		conds := []Token{}
		exprs := []Token{}
		for i, arg := range args {
			if i%2 == 0 {
				conds = append(conds, arg)
			} else {
				exprs = append(exprs, arg)
			}
		}

		return NewIfElse(conds, exprs, ctx)
	case name == "int":
		switch len(args) {
		case 2:
			return NewIntegral(args[0], nil, nil, args[1], ctx)
		case 3:
			return NewIntegral(args[0], args[1], nil, args[2], ctx)
		case 4:
			return NewIntegral(args[0], args[1], args[2], args[3], ctx)
		default:
			return nil, ctx.NewError("Error: expected 2, 3 or 4 arguments")
		}
	case name == "oint":
		if len(args) != 2 {
			return nil, ctx.NewError("Error: expected 2 arguments")
		}
		return NewContourIntegral(args[0], args[1], ctx)
	case name == "sum":
		switch len(args) {
		case 1:
			return NewSum(args[0], nil, nil, ctx)
		case 2:
			return NewSum(args[0], args[1], nil, ctx)
		case 3:
			return NewSum(args[0], args[1], args[2], ctx)
		default:
			return nil, ctx.NewError("Error: expected 1, 2 or 3 arguments")
		}
	default:
		nameToken, err := NewWord(name, ctx)
		if err != nil {
			return nil, err
		}

		pContent, _ := NewCSV(args, ctx)
		parens, _ := NewParens(pContent, ctx)

		return &Call{nameToken, parens, newTokenData(ctx)}, nil
	}
}

func (t *Call) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Call ")
	b.WriteString(t.name.Dump(""))
	b.WriteString("\n")
	b.WriteString(t.args.Dump(indent + "  "))

	return b.String()
}

func (t *Call) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bbName, err := t.name.GenerateTags(scope, x, y)
	if err != nil {
		return nil, err
	}

	bbArgs, err := t.args.GenerateTags(scope, bbName.Right()+extraCallSpacing, y)
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bbName, bbArgs), nil
}
