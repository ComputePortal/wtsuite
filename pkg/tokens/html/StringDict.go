package html

import (
	"reflect"
  "sort"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type StringDict struct {
	RawDict
}

func IsStringDict(t Token) bool {
	_, ok := t.(*StringDict)
	return ok
}

func AssertStringDict(t Token) (*StringDict, error) {
	if d, ok := t.(*StringDict); ok {
		return d, nil
	} else if rd, ok := t.(*RawDict); ok {
		return rd.ToStringDict()
	} else {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected string dict (got " + reflect.TypeOf(t).String() + ")")
		panic(err)
		return nil, err
	}
}

func NewEmptyStringDict(ctx context.Context) *StringDict {
	return &StringDict{RawDict{make([]struct{ key, value Token }, 0), TokenData{ctx}}}
}

func (t *StringDict) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *StringDict) ToRaw() *RawDict {
	return &RawDict{t.items, TokenData{t.Context()}}
}

func (a *StringDict) IsSame(other Token) bool {
	if b, ok := other.(*StringDict); ok {
		if a.Len() != b.Len() {
			return false
		}

		for _, item := range a.items {
			akey, err := AssertString(item.key)
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

func (t *StringDict) CopyStringDict(ctx context.Context) (*StringDict, error) {
	res := NewEmptyStringDict(ctx)

	for _, item := range t.items {
		tkey, err := AssertString(item.key)
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

func (t *StringDict) Copy(ctx context.Context) (Token, error) {
	return t.CopyStringDict(ctx)
}

func (t *StringDict) GetKeyValue(x interface{}) (*String, Token, bool) {
	key := t.convertKey(x)

	// scan in reverse order (last set is first out)
	for i := len(t.items) - 1; i >= 0; i-- {
		item := t.items[i]
		if tkey, ok := item.key.(*String); ok && tkey.Value() == key {
			return tkey, item.value, true
		}
	}

	return nil, nil, false
}

func (t *StringDict) Get(x interface{}) (Token, bool) {
	_, val, ok := t.GetKeyValue(x)
	return val, ok
}

// this is a raw append method, duplication isnt checked, last entries are all that matters anyway
func (t *StringDict) Set(x_ interface{}, value Token) {
	var key *String
	switch x := x_.(type) {
	case string:
		key = NewValueString(x, value.Context())
	case *String:
		key = x
	default:
		panic("expected string")
	}

	for i, item := range t.items {
		if check, ok := item.key.(*String); ok && check.Value() == key.Value() {
			t.items[i] = struct{ key, value Token }{key, value}
			return
		}
	}

	t.items = append(t.items, struct{ key, value Token }{key, value})
}

// try to call this as little as possible
func (t *StringDict) Delete(x interface{}) {
	key := t.convertKey(x)

	tmpItems := make([]struct{ key, value Token }, 0)
	for _, item := range t.items {
		if check, ok := item.key.(*String); ok && check.Value() != key {
			tmpItems = append(tmpItems, item)
		} else if !ok {
			panic("bad key")
		}
	}

	t.items = tmpItems
}

func (t *StringDict) Loop(fn func(key *String, value Token, last bool) error) error {
	count := 0
	n := len(t.items)
	for _, item := range t.items {
		count++
		key, err := AssertString(item.key)
		if err != nil {
			panic(err)
		}

		if err := fn(key, item.value, count == n); err != nil {
			return err
		}
	}

	return nil
}

func (t *StringDict) AssertOnlyValidKeys(validKeys []string) error {
	for _, item := range t.items {
		key, err := AssertString(item.key)
		if err != nil {
			panic(err)
		}

    found := false
    for _, test := range validKeys {
      if test == key.Value() {
        found = true
        break
      }
    }

    if !found {
      errCtx := key.Context()
      return errCtx.NewError("Error: invalid attribute")
    }
	}

  return nil
}

func GolangStringMapToStringDict(m map[string]interface{}, ctx context.Context) (*StringDict, error) {
  // sort the keys alphabetically!
  keys := make([]string, 0)
  for k, _ := range m {
    keys = append(keys, k)
  }

  sort.Strings(keys)

  res := NewEmptyStringDict(ctx)

  for _, k_ := range keys {
    item_ := m[k_]

    item, err := GolangToToken(item_, ctx)
    if err != nil {
      return nil, err
    }

    k := NewValueString(k_, ctx)

    res.Set(k, item)
  }

  return res, nil
}
