package html

import (
	"fmt"
	"reflect"
  "sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type DictType int

const (
	ANY DictType = iota
	STRING
	INT
)

// unevaluated dict, when evaluated becomes IntDict or StringDict
type RawDict struct {
	items []struct{ key, value Token }
	TokenData
}

func NewValuesRawDict(keys []Token, values []Token, ctx context.Context) *RawDict {
	items := make([]struct{ key, value Token }, 0)

	if len(keys) != len(values) {
		panic("inconsistent lengths")
	}

	for i, k := range keys {
		v := values[i]

		items = append(items, struct{ key, value Token }{k, v})
	}

	return &RawDict{items, TokenData{ctx}}
}

func NewEmptyRawDict(ctx context.Context) *RawDict {
	return NewValuesRawDict([]Token{}, []Token{}, ctx)
}

func (t *RawDict) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	if len(t.items) == 0 {
		b.WriteString("Dict\n")
	} else {
		switch t.items[0].key.(type) {
		case *Int:
			b.WriteString("IntDict\n")
		case *String:
			b.WriteString("StringDict\n")
		default:
			b.WriteString("RawDict\n")
		}
	}

	for _, item := range t.items {
		b.WriteString(item.key.Dump(indent + "  "))
		b.WriteString(item.value.Dump(indent + "   :"))
	}

	return b.String()
}

// only eval if scope != nil
func (t *RawDict) toKeyDict(scope Scope, dtype DictType, scanFn func(Token, Token) error) (Dict, error) {
	result := make([]struct{ key, value Token }, 0)

	strTable := make(map[string]Token)
	intTable := make(map[int]Token)

	for _, item := range t.items {
		key := item.key
		if scope != nil {
			var err error
			key, err = key.Eval(scope)
			if err != nil {
				return nil, err
			}
		}

		if dtype == ANY {
			switch key.(type) {
      case *Lazy:
			case *String:
				dtype = STRING
			case *Int:
				dtype = INT
			default:
				errCtx := key.Context()
				return nil, errCtx.NewError("Error: expected int or string keys, got " + reflect.TypeOf(key).String())
			}
		}

		value := item.value
		if scope != nil {
			var err error
			value, err = item.value.Eval(scope)
			if err != nil {
				return nil, err
			}
		}

		if scanFn != nil {
			if err := scanFn(key, value); err != nil {
				return nil, err
			}
		}

    if IsLazy(key) {
      result = append(result, struct{ key, value Token }{key, value}) // keep evaluated key for context

    } else {
      switch dtype {
      case ANY:
        panic("should'be been set")
      case STRING:
        strKey, ok := key.(*String)
        if !ok {
          errCtx := key.Context()
          return nil, errCtx.NewError("Error: expected string, got " + reflect.TypeOf(key).String())
        }

        // is key already defined?
        if other, ok := strTable[strKey.Value()]; ok {
          errCtx := strKey.Context()
          err := errCtx.NewError("Error: key \"" + strKey.Value() + "\" already defined")
          err.AppendContextString("Info: defined here", other.Context())

          return nil, err
        }

        // single underscore attributes are always ignored and can be used as a scope for tag variables (eg. __elementCount__)
        if strKey.Value() != "_" {
          strTable[strKey.Value()] = value
          result = append(result, struct{ key, value Token }{strKey, value}) // keep evaluated key for context
        }
      case INT:
        intKey, ok := key.(*Int)
        if !ok {
          errCtx := key.Context()
          return nil, errCtx.NewError("Error: expected string, got " + reflect.TypeOf(key).String())
        }

        // is key already defined?
        if other, ok := intTable[intKey.Value()]; ok {
          errCtx := intKey.Context()
          err := errCtx.NewError("Error: key \"" + fmt.Sprintf("%d", intKey.Value()) + "\" already defined")
          err.AppendContextString("Info: defined here", other.Context())

          return nil, err
        }

        intTable[intKey.Value()] = value

        result = append(result, struct{ key, value Token }{intKey, value}) // keep evaluated key for context
      default:
        panic("unhandled")
      }
    }
	}

	switch dtype {
	case ANY:
		// empty dict, assume string dict
		return &StringDict{RawDict{result, TokenData{t.Context()}}}, nil
	case STRING:
		return &StringDict{RawDict{result, TokenData{t.Context()}}}, nil
	case INT:
		return &IntDict{RawDict{result, TokenData{t.Context()}}}, nil
	default:
		panic("unhandled")
	}
}

func (t *RawDict) EvalRawDict(scope Scope) (*RawDict, error) {
  return t.EvalRawDictScan(scope, nil)
}

func (t *RawDict) EvalRawDictScan(scope Scope,
	scanFn func(Token, Token) error) (*RawDict, error) {
	// we know that attr enums

	result := make([]struct{ key, value Token }, 0)
	// do a test build first
	for _, item := range t.items {
		key, err := item.key.Eval(scope)
		if err != nil {
			return nil, err
		}

		val, err := item.value.Eval(scope)
		if err != nil {
			return nil, err
		}

    if scanFn != nil {
      if err := scanFn(key, val); err != nil {
        return nil, err
      }
    }

		result = append(result, struct{ key, value Token }{key, val})
	}

	return &RawDict{result, TokenData{t.Context()}}, nil
}

// could be IntDict, so code differs from EvalDict
func (t *RawDict) Eval(scope Scope) (Token, error) {
	d, err := t.toKeyDict(scope, ANY, nil)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (t *RawDict) evalLazy(tag FinalTag, dType DictType) error {
  // inplace evaluation of lazy
  for i, item := range t.items {
    keyPre := item.key
    valPre := item.value

    keyPost := keyPre
    if IsLazy(keyPre) {
      var err error
      keyPost, err = keyPre.EvalLazy(tag)
      if err != nil {
        return err
      }

      switch dType {
      case STRING:
        // TODO: check uniqueness
        if _, err := AssertString(keyPost); err != nil {
          return err
        }
      case INT:
        // TODO: check uniqueness
        if _, err := AssertInt(keyPost); err != nil {
          return err
        }
      }
    }

    valPost, err := valPre.EvalLazy(tag)
    if err != nil {
      return err
    }

    t.items[i] = struct{key, value Token}{
      keyPost,
      valPost,
    }
  }

  return nil
}

func (t *RawDict) EvalLazy(tag FinalTag) (Token, error) {
  if err := t.evalLazy(tag, ANY); err != nil {
    return nil, err
  }

  return t, nil
}

func (t *RawDict) convertKey(x interface{}) string {
	key := ""
	switch x_ := x.(type) {
	case string:
		key = x_
	case *String:
		key = x_.Value()
	default:
		panic("expected string or String got " + reflect.TypeOf(x).String())
	}

	return key
}

func (t *RawDict) GetKeyValue(x interface{}) (*String, Token, bool) {
	xKey := t.convertKey(x)

	for _, item := range t.items {
		key_ := item.key

		if IsString(key_) {
			key, err := AssertString(key_)
			if err != nil {
				panic(err)
			}

			if key.Value() == xKey {
				return key, item.value, true
			}
		}
	}

	return nil, nil, false
}

func (t *RawDict) Get(x interface{}) (Token, bool) {
	_, v, ok := t.GetKeyValue(x)
	return v, ok
}

func (t *RawDict) toStringDict(scope Scope, scanFn func(*String, Token) error) (*StringDict, error) {

	var innerScanFn func(Token, Token) error = nil
	if scanFn != nil {
		innerScanFn = func(key_ Token, val Token) error {
			key, err := AssertString(key_)
			if err != nil {
				return err
			}

			return scanFn(key, val)
		}
	}

	d, err := t.toKeyDict(scope, STRING, innerScanFn)
	if err != nil {
		return nil, err
	}

	strDict, _ := d.(*StringDict)
	return strDict, nil
}

func (t *RawDict) toIntDict(scope Scope) (*IntDict, error) {
	d, err := t.toKeyDict(scope, INT, nil)
	if err != nil {
		return nil, err
	}

	intDict, _ := d.(*IntDict)
	return intDict, nil
}

func (t *RawDict) Set(key Token, value Token) {
	t.items = append(t.items, struct{ key, value Token }{key, value})
}

func (t *RawDict) EvalStringDict(scope Scope) (*StringDict, error) {
	return t.toStringDict(scope, nil)
}

func (t *RawDict) EvalStringDictScan(scope Scope, scanFn func(*String, Token) error) (*StringDict, error) {
	return t.toStringDict(scope, scanFn)
}

func (t *RawDict) EvalIntDict(scope Scope) (*IntDict, error) {
	return t.toIntDict(scope)
}

// without evaluation
func (t *RawDict) ToStringDict() (*StringDict, error) {
	return t.toStringDict(nil, nil)
}

func ToStringDict(t Token) (*StringDict, error) {
	switch res := t.(type) {
	case *RawDict:
		return res.ToStringDict()
	case *StringDict:
		return res, nil
	default:
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected dict, got " + reflect.TypeOf(t).String())
		return nil, err
	}
}

func (t *RawDict) ToIntDict() (*IntDict, error) {
	return t.toIntDict(nil)
}

func (t *RawDict) Len() int {
	return len(t.items)
}

func (t *RawDict) IsEmpty() bool {
	return t.Len() == 0
}

func (t *RawDict) Copy(ctx context.Context) (Token, error) {
	items := make([]struct{ key, value Token }, len(t.items))

	copyFn := func(t Token) (Token, error) {
		if IsContainer(t) {
			c, err := AssertContainer(t)
			if err != nil {
				panic(err)
			}

			return c.Copy(ctx)
		} else {
			return t, nil
		}
	}

	for i, item := range t.items {
		key, err := copyFn(item.key)
		if err != nil {
			return nil, err
		}

		value, err := copyFn(item.value)
		if err != nil {
			return nil, err
		}

		items[i] = struct{ key, value Token }{key, value}
	}

	return &RawDict{items, TokenData{ctx}}, nil
}

// order is important (which is not the case for KeyDict)
func (a *RawDict) IsSame(other Token) bool {
	if b, ok := other.(*RawDict); ok {
		if a.Len() != b.Len() {
			return false
		}

		for i, item := range a.items {
			akey := item.key
			avalue := item.value

			bkey := b.items[i].key
			bvalue := b.items[i].value

			if !akey.IsSame(bkey) || !avalue.IsSame(bvalue) {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func (t *RawDict) Loop(fn func(key Token, value Token, last bool) error) error {
	count := 0
	n := len(t.items)
	for _, item := range t.items {
		count++
		if err := fn(item.key, item.value, count == n); err != nil {
			return err
		}
	}

	return nil
}

func (t *RawDict) ContainsLazy() bool {
	for _, item := range t.items {
    if IsLazy(item.key) || IsLazy(item.value) {
      return true
    }

    switch {
      case IsRawDict(item.value):
        val, err := AssertRawDict(item.value)
        if err != nil {
          panic(err)
        }
        if val.ContainsLazy() {
          return true
        }
      case IsStringDict(item.value):
        val, err := AssertStringDict(item.value)
        if err != nil {
          panic(err)
        }

        if val.ContainsLazy() {
          return true
        }
    }
	}

  return false
}

func (t *RawDict) LoopValues(fn func(t Token) error) error {
	for _, item := range t.items {
		switch v := item.value.(type) {
		case *List:
			if err := v.LoopValues(fn); err != nil {
				return err
			}
		case Container:
			if err := v.LoopValues(fn); err != nil {
				return err
			}
		default:
			if err := fn(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func IsRawDict(t Token) bool {
	_, ok := t.(*RawDict)
	return ok
}

func AssertRawDict(t_ Token) (*RawDict, error) {
	if t, ok := t_.(*RawDict); ok {
		return t, nil
	} else {
		errCtx := t_.Context()
		return nil, errCtx.NewError("Internal Error: expected a raw dict")
	}
}

func GolangStringMapToRawDict(m map[string]interface{}, ctx context.Context) (*RawDict, error) {
  // sort the keys alphabetically!
  keys := make([]string, 0)
  for k, _ := range m {
    keys = append(keys, k)
  }

  sort.Strings(keys)

  res := NewEmptyRawDict(ctx)

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
