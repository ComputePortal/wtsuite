package html

import (
	"../context"
)

type KeyDict interface {
	Dict
	Get(key interface{}) (Token, bool)
	Set(key interface{}, value Token)
	Delete(key interface{})
}

type RawKeyDict struct {
	RawDict
}

func IsKeyDict(t Token) bool {
	_, ok := t.(KeyDict)
	return ok
}

func AssertKeyDict(t Token) (KeyDict, error) {
	if d, ok := t.(KeyDict); ok {
		return d, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected KeyDict (StringDict or IntDict)")
	}
}

func DictHasFlag(d KeyDict, x interface{}) (bool, error) {
	val, ok := d.Get(x)
	if !ok {
		return false, nil
	}

	if err := AssertFlag(val); err != nil {
		return false, err
	}

	return true, nil
}

func DictOptionString(d KeyDict, x interface{}, def string) (string, error) {
	val_, ok := d.Get(x)
	if !ok {
		return def, nil
	}

	val, err := AssertString(val_)
	if err != nil {
		return def, err
	}

	if val.Value() == "" {
		errCtx := val.Context()
		return def, errCtx.NewError("Error: expected non-empty string")
	}

	return val.Value(), nil
}

func makeNotFoundError(x_ interface{}, ctx context.Context) error {
	switch x := x_.(type) {
	case string:
		return ctx.NewError("Error: " + x + " not found")
	case *String:
		errCtx := x.Context()
		return errCtx.NewError("Error: " + x.Value() + " not found")
	case int:
		return ctx.NewError("Error: not found")
	default:
		panic("unexpected")
	}
}

func DictString(d KeyDict, x interface{}) (*String, error) {
	val_, ok := d.Get(x)
	if !ok {
		err := makeNotFoundError(x, d.Context())
		return nil, err
	}

	val, err := AssertString(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictBool(d KeyDict, x interface{}) (*Bool, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertBool(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictFloat(d KeyDict, x interface{}) (*Float, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertIntOrFloat(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictInt(d KeyDict, x interface{}) (*Int, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertInt(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictPrimitive(d KeyDict, x interface{}) (Primitive, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertPrimitive(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictKeyDict(d KeyDict, x interface{}) (KeyDict, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertKeyDict(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictStringDict(d KeyDict, x interface{}) (*StringDict, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertStringDict(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictRawDict(d KeyDict, x interface{}) (*RawDict, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertRawDict(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictIntDict(d KeyDict, x interface{}) (*IntDict, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertIntDict(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictList(d KeyDict, x interface{}) (*List, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertList(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func DictColor(d KeyDict, x interface{}) (*Color, error) {
	val_, ok := d.Get(x)
	if !ok {
		return nil, makeNotFoundError(x, d.Context())
	}

	val, err := AssertColor(val_)
	if err != nil {
		return nil, err
	}

	return val, nil
}
