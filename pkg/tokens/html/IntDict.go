package html

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// IntDict is always the result of an evaluation
type IntDict struct {
	RawDict
}

func IsIntDict(t Token) bool {
	_, ok := t.(*IntDict)
	return ok
}

func AssertIntDict(t Token) (*IntDict, error) {
	if d, ok := t.(*IntDict); ok {
		return d, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected int dict")
	}
}

func NewEmptyIntDict(ctx context.Context) *IntDict {
	return &IntDict{RawDict{make([]struct{ key, value Token }, 0), TokenData{ctx}}}
}

func (t *IntDict) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (a *IntDict) IsSame(other Token) bool {
	if b, ok := other.(*IntDict); ok {
		if a.Len() != b.Len() {
			return false
		}

		for _, item := range a.items {
			akey, err := AssertInt(item.key)
			if err != nil {
				panic(err)
			}

			avalue := item.value

			if bvalue, ok := b.Get(akey); !ok {
				return false
			} else {
				if !avalue.IsSame(bvalue) {
					return false
				}
			}
		}

		return true
	} else {
		return false
	}
}

func (t *IntDict) CopyIntDict(ctx context.Context) (*IntDict, error) {
	res := NewEmptyIntDict(ctx)

	for _, item := range t.items {
		tkey, err := AssertInt(item.key)
		if err != nil {
			return nil, err
		}

		if IsContainer(item.value) {
			value, err := AssertContainer(item.value)
			if err != nil {
				panic(err)
			}
			copy, err := value.Copy(value.Context())
			if err != nil {
				return nil, err
			}
			res.Set(tkey, copy)
		} else {
			res.Set(tkey, item.value)
		}
	}

	return res, nil
}

func (t *IntDict) Copy(ctx context.Context) (Token, error) {
	return t.CopyIntDict(ctx)
}

func (t *IntDict) convertKey(x interface{}) int {
	key := 0
	switch x_ := x.(type) {
	case int:
		key = x_
	case *Int:
		key = x_.Value()
	default:
		panic("expected int or tokens.Int")
	}

	return key
}

func (t *IntDict) GetKeyValue(x interface{}) (*Int, Token, bool) {
	key := t.convertKey(x)

	for i := len(t.items) - 1; i >= 0; i-- {
		item := t.items[i]
		if tkey, ok := item.key.(*Int); ok && tkey.Value() == key {
			return tkey, item.value, true
		}
	}

	return nil, nil, false
}

func (t *IntDict) Get(x interface{}) (Token, bool) {
	_, val, ok := t.GetKeyValue(x)
	return val, ok
}

func (t *IntDict) Set(x_ interface{}, value Token) {
	var key *Int
	switch x := x_.(type) {
	case int:
		key = NewValueInt(x, value.Context())
	case *Int:
		key = x
	default:
		panic("expected string")
	}

	t.items = append(t.items, struct{ key, value Token }{key, value})
}

func (t *IntDict) Delete(x interface{}) {
	key := t.convertKey(x)

	tmpItems := make([]struct{ key, value Token }, 0)
	for _, item := range t.items {
		if check, ok := item.key.(*Int); ok && check.Value() != key {
			tmpItems = append(tmpItems, item)
		} else if !ok {
			panic("bad key")
		}
	}

	t.items = tmpItems
}

func (t *IntDict) Loop(fn func(key *Int, value Token, last bool) error) error {
	count := 0
	n := len(t.items)
	for _, item := range t.items {
		count++
		key, err := AssertInt(item.key)
		if err != nil {
			panic(err)
		}

		if err := fn(key, item.value, count == n); err != nil {
			return err
		}
	}

	return nil
}
