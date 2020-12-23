package html

import (
	"errors"
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type List struct {
	values []Token // without brackets and commas
	TokenData
}

func NewValuesList(values interface{}, ctx context.Context) *List {
	switch vals := values.(type) {
	case []*Int:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*Float:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*Bool:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*Color:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*String:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*Null:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*List:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*StringDict:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*IntDict:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []*RawDict:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = v
		}
		return &List{result, TokenData{ctx}}
	case []string:
		result := make([]Token, len(vals))
		for i, v := range vals {
			result[i] = NewValueString(v, ctx)
		}
		return &List{result, TokenData{ctx}}
	case []Token:
		return &List{vals, TokenData{ctx}}
	default:
		panic("invalid list type")
	}
}

func NewEmptyList(ctx context.Context) *List {
	return &List{make([]Token, 0), TokenData{ctx}}
}

func NewNilList(n int, ctx context.Context) *List {
	return &List{make([]Token, n), TokenData{ctx}}
}

func interfaceToInt(x interface{}) int {
	v := 0
	switch x_ := x.(type) {
	case int:
		v = x_
	case *Int:
		v = x_.Value()
	default:
		panic("expected int or tokens.Int")
	}

	return v
}

func (t *List) Get(x interface{}) (Token, error) {
	index := interfaceToInt(x)

	if index < 0 || index >= t.Len() {
		return nil, errors.New(fmt.Sprintf("out of range, %d doesn't lie in [0:%d)", index, t.Len()))
	}

	return t.values[index], nil
}

func (t *List) Append(v Token) {
	t.values = append(t.values, v)
}

func (t *List) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent + "List\n")

	for _, value := range t.values {
		b.WriteString(value.Dump(indent + "| "))
	}

	return b.String()
}

func (t *List) EvalList(scope Scope) (*List, error) {
	result := make([]Token, 0)

	for _, value := range t.values {
		a, err := value.Eval(scope)
		if err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return &List{result, TokenData{t.Context()}}, nil
}

func (t *List) Eval(scope Scope) (Token, error) {
	return t.EvalList(scope)
}

func (t *List) Len() int {
	return len(t.values)
}

func (t *List) Loop(fn func(i int, value Token, last bool) error) error {
	n := len(t.values)

	for i, value := range t.values {
		if err := fn(i, value, i == n-1); err != nil {
			return err
		}
	}

	return nil
}

// loop in which indices are not accessable
func (t *List) LoopValues(fn func(t Token) error) error {
	for _, v_ := range t.values {
		switch v := v_.(type) {
		case *List: // uncertain if List respects Container interface
			if err := v.LoopValues(fn); err != nil {
				return nil
			}
		case Container:
			if err := v.LoopValues(fn); err != nil {
				return nil
			}
		default:
			if err := fn(v_); err != nil {
				return nil
			}
		}
	}

	return nil
}

func (t *List) CopyList(ctx context.Context) (*List, error) {
	res := NewEmptyList(ctx)

	for _, value := range t.values {
		if IsContainer(value) {
			value_, err := AssertContainer(value)
			if err != nil {
				panic(err)
			}

			copy, err := value_.Copy(value_.Context())
			if err != nil {
				return nil, err
			}

			res.values = append(res.values, copy)
		} else {
			res.values = append(res.values, value)
		}
	}

	return res, nil
}

func (t *List) Copy(ctx context.Context) (Token, error) {
	return t.CopyList(ctx)
}

func IsList(t Token) bool {
	_, ok := t.(*List)
	return ok
}

func IsStringList(t Token) bool {
	l, ok := t.(*List)
	if !ok {
		return false
	}

	for _, value := range l.values {
		b := IsString(value)
		if !b {
			return false
		}
	}

	return true
}

func IsIntList(t Token) bool {
	l, ok := t.(*List)
	if !ok {
		return false
	}

	for _, value := range l.values {
		b := IsInt(value)
		if !b {
			return false
		}
	}

	return true
}

func AssertList(t Token) (*List, error) {
	l, ok := t.(*List)
	if !ok {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected list")
		return nil, err
	}

	return l, nil
}

// ignore first null to accomodate enum attrs
func (t *List) GetStrings() ([]string, error) {
	res := make([]string, 0)

	for i, value := range t.values {
		if IsNull(value) && i == 0 {
			continue
		}
		t, err := AssertPrimitive(value)
		if err != nil {
			return res, err
		}

		res = append(res, t.Write())
	}

	return res, nil
}

func (t *List) GetTokens() []Token {
	return t.values[:]
}

func (a *List) IsSame(other Token) bool {
	if b, ok := other.(*List); ok {
		if a.Len() == b.Len() {
			for i, _ := range a.values {
				if !a.values[i].IsSame(b.values[i]) {
					return false
				}
			}

			return true
		}
	}

	return false
}

func ToStringList(t_ Token) (*List, error) {
	t, ok := t_.(*List)
	if !ok {
		errCtx := t_.Context()
		return nil, errCtx.NewError("Error: expected list")
	}

	for _, v := range t.values {
		if _, err := AssertString(v); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func GolangSliceToList(x []interface{}, ctx context.Context) (*List, error) {
  res := NewNilList(len(x), ctx)

  for i := 0; i < len(x); i++ {
    item, err := GolangToToken(x[i], ctx)
    if err != nil {
      return nil, err
    }

    res.values[i] = item
  }

  return res, nil
}
