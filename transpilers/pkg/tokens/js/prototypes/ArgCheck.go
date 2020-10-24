package prototypes

import (
	"fmt"
	"reflect"
	"strconv"

	"../values"

	"../../context"
)

type ArgCheck interface {
	// uses IsInstanceOf
	Check(args []interface{}, pos int, ctx context.Context) (int, error)
}

func CheckPrototype(p values.Prototype, args []interface{}, pos int,
	ctx context.Context) (int, error) {
	if pos >= len(args) {
		return 0, ctx.NewError(fmt.Sprintf("Error: expected an argument at %d", pos+1))
	}

	arg_ := args[pos]

	switch arg := arg_.(type) {
	case values.Value:
		if !arg.IsInstanceOf(p) {
			errCtx := arg.Context()
			err := errCtx.NewError("Error: expected " + p.Name() +
				", got " + arg.TypeName())
			return 0, err
		}
	case string:
		if arg != p.Name() {
			errCtx := ctx
			err := errCtx.NewError("Error: expected " + p.Name() +
				", got " + arg)
			return 0, err
		}
	default:
		panic("expected Value or string")
	}

	return pos + 1, nil
}

func CheckInterface(interf values.Interface, args []interface{}, pos int,
	ctx context.Context) (int, error) {
	if pos >= len(args) {
		return 0, ctx.NewError(fmt.Sprintf("Error: expected an argument at %d", pos+1))
	}

	arg_ := args[pos]

	switch arg := arg_.(type) {
	case values.Value:
		argInstances := values.UnpackMulti([]values.Value{arg})
		for _, instance := range argInstances {
			proto, ok := instance.GetInstancePrototype()
			if !ok {
				errCtx := arg.Context()
				err := errCtx.NewError("Error: not an instance " + reflect.TypeOf(instance).String())
				return 0, err
			}

			if msg, ok := interf.IsImplementedBy(proto); !ok {
				errCtx := arg.Context()
				err := errCtx.NewError("Error: " + arg.TypeName() + " doesn't implement " +
					interf.Name() + ", " + msg)
				return 0, err
			}
		}
	case string:
		if interf.Name() != arg {
			errCtx := ctx
			err := errCtx.NewError("Error: got " + arg + ", expected " +
				interf.Name())
			return 0, err
		}
	default:
		panic("expected Value or string")
	}

	return pos + 1, nil
}

type None struct {
}

func (c *None) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	if len(args) > pos {

		if arg, ok := args[len(args)-1].(values.Value); ok {
			errCtx := arg.Context()
			return 0, errCtx.NewError("Error: unexpected argument")
		} else {
			errCtx := ctx
			return 0, errCtx.NewError("Error: unexpected argument at pos " + strconv.Itoa(pos))
		}
	}

	return pos, nil
}

type Opt struct {
	a ArgCheck
}

func (c *Opt) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	if pos >= len(args) {
		return pos, nil
	} else {
		return c.a.Check(args, pos, ctx)
	}
}

type And struct {
	a, b ArgCheck
}

func (c *And) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	apos, err := c.a.Check(args, pos, ctx)
	if err != nil {
		return 0, err
	}

	bpos, err := c.b.Check(args, apos, ctx)
	if err != nil {
		return 0, err
	}

	return bpos, nil
}

type Rest struct {
	a ArgCheck
}

func (c *Rest) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	for pos < len(args) {
		var err error
		pos, err = c.a.Check(args, pos, ctx)
		if err != nil {
			return 0, err
		}
	}

	return pos, nil
}

type Many struct {
	n int
	a ArgCheck
}

func (c *Many) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	for i := 0; i < c.n; i++ {
		var err error
		pos, err = c.a.Check(args, pos, ctx)
		if err != nil {
			return 0, err
		}
	}

	return pos, nil
}

type AtLeast struct {
	atLeast int
	a       ArgCheck
}

func (c *AtLeast) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	count := 0
	var err error
	for pos < len(args) {
		pos, err = c.a.Check(args, pos, ctx)
		if err != nil {
			return 0, nil
		}

		count++
	}

	if count < c.atLeast {
		return 0, ctx.NewError("Error: expected at least " + strconv.Itoa(c.atLeast) + " more arguments, but got " + strconv.Itoa(count))
	}

	return len(args), nil
}

type Any struct {
}

func (c *Any) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	if pos >= len(args) {
		return 0, ctx.NewError("Error: expected an argument at " + strconv.Itoa(pos+1))
	}

	return pos + 1, nil
}

type Or struct {
	a, b ArgCheck
}

func (c *Or) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	apos, aerr := c.a.Check(args, pos, ctx)
	bpos, berr := c.b.Check(args, pos, ctx)

	if aerr != nil {
		if berr != nil {
			return 0, aerr
		} else {
			return bpos, nil
		}
	} else {
		if berr != nil {
			return apos, nil
		} else if apos == bpos {
			return apos, nil
		} else {
			panic("unexpected Or input check (both paths are valid, but both paths returns different positions)")
		}
	}
}

type Function struct {
}

func (c *Function) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	if pos >= len(args) {
		return 0, ctx.NewError("Error: expected an argument at " + strconv.Itoa(pos+1))
	}

	arg, ok := args[pos].(values.Value)
	if !ok {
		panic("not a value")
	}

	if !arg.IsFunction() {
		err := ctx.NewError("Error: expected a function at " + strconv.Itoa(pos+1))
		return 0, err
	}

	return pos + 1, nil
}

// use this if only values.Prototype is available (and not *BuiltinPrototype)
type PrototypeCheck struct {
	proto values.Prototype
}

func (c *PrototypeCheck) Check(args []interface{}, pos int,
	ctx context.Context) (int, error) {
	if pos >= len(args) {
		return 0, ctx.NewError(fmt.Sprintf("Error: expected an argument at %d", pos+1))
	}

	arg_ := args[pos]

	switch arg := arg_.(type) {
	case values.Value:
		if !arg.IsInstanceOf(c.proto) {
			errCtx := arg.Context()
			err := errCtx.NewError("Error: expected " + c.proto.Name() +
				", got " + arg.TypeName())
			return 0, err
		}
	case string:
		if arg != c.proto.Name() {
			errCtx := ctx
			err := errCtx.NewError("Error: expected " + c.proto.Name() +
				", got " + arg)
			return 0, err
		}
	default:
		panic("expected Value or string")
	}

	return pos + 1, nil
}

func CheckInputs(c ArgCheck, args []values.Value, ctx context.Context) error {
	pos := 0

	args_ := make([]interface{}, len(args))
	for i, arg := range args {
		args_[i] = arg
	}

	var err error
	pos, err = c.Check(args_, pos, ctx)
	if err != nil {
		return err
	}

	if pos != len(args) {
		err := ctx.NewError(fmt.Sprintf("Error: expected at most %d arguments, got %d", pos, len(args)))
		return err
	}

	return nil
}
