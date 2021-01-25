package html

import (
  "reflect"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Parens struct {
	values []Token
	alts   []Token // rhs of argDefaults for function and class, can contain nil in case of no defaults
	TokenData
}

func NewParensInterf(argNames []string, alts []Token, ctx context.Context) *Parens {
  values := make([]Token, len(argNames))
  for i, name := range argNames {
    values[i] = NewValueString(name, ctx)
  }

  return NewParens(values, alts, ctx)
}

func NewParens(values []Token, alts []Token, ctx context.Context) *Parens {
  if alts == nil {
    alts = make([]Token, len(values))
    for i, _ := range alts {
      alts[i] = nil
    }
  }

	if len(values) != len(alts) {
		panic("expected same lenghts")
	}

	return &Parens{values, alts, TokenData{ctx}}
}

func (t *Parens) Values() []Token {
	return t.values
}

func (t *Parens) Alts() []Token {
	return t.alts
}

func (t *Parens) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Parens(\n")

	for i, v := range t.values {
		b.WriteString(v.Dump(indent + "  "))
		b.WriteString("\n")
		if t.alts[i] != nil {
			b.WriteString(t.alts[i].Dump(indent + "= "))
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *Parens) Eval(scope Scope) (Token, error) {
	if len(t.values) != 1 || t.alts[0] != nil {
		errCtx := t.Context()
		err := errCtx.NewError("Error: bad parens (not a function or class declaration)")
		panic(err)
		return nil, err
	}

	return t.values[0].Eval(scope)
}

func (t *Parens) EvalLazy(tag FinalTag) (Token, error) {
	if len(t.values) != 1 || t.alts[0] != nil {
		errCtx := t.Context()
		err := errCtx.NewError("Error: bad parens (not a function or class declaration)")
		panic(err)
		return nil, err
	}

  return t.values[0].EvalLazy(tag)
}

func (t *Parens) EvalAsArgs(scope Scope) (*Parens, error) {
  values := make([]Token, len(t.values))
  alts := make([]Token, len(t.alts))

  hadAltBefore := false
  var err error
  for i, v := range t.values {
    alt := t.alts[i]

    if alt == nil {
      if hadAltBefore {
        errCtx := v.Context()
        return nil, errCtx.NewError("Error: after kwargs")
      }

      values[i], err = v.Eval(scope)
      if err != nil {
        return nil, err
      }

      alts[i] = nil
    } else {
      if w, err := AssertWord(v); err != nil {
        return nil, err
      } else {
        // check that kwarg wasn't already specified
        for j, vCheck := range t.values[0:i] {
          if t.alts[j] == nil {
            continue
          }

          wCheck, err := AssertWord(vCheck)
          if err != nil {
            panic(err)
          }

          if wCheck.Value() == w.Value() {
            errCtx := w.Context()
            return nil, errCtx.NewError("Error: duplicate kwarg")
          }
        }
      }

      values[i] = v

      alts[i], err = alt.Eval(scope)
      if err != nil {
        return nil, err
      }

      hadAltBefore = true
    }
  }

  return NewParens(values, alts, t.Context()), nil
}

func (t *Parens) EvalAsArgsLazy(tag FinalTag) (*Parens, error) {
  values := make([]Token, len(t.values))
  alts := make([]Token, len(t.alts))

  hadAltBefore := false
  var err error
  for i, v := range t.values {
    alt := t.alts[i]

    if alt == nil {
      if hadAltBefore {
        errCtx := v.Context()
        return nil, errCtx.NewError("Error: after kwargs")
      }

      values[i], err = v.EvalLazy(tag)
      if err != nil {
        return nil, err
      }

      alts[i] = nil
    } else {
      if w, err := AssertWord(v); err != nil {
        return nil, err
      } else {
        // check that kwarg wasn't already specified
        for j, vCheck := range t.values[0:i] {
          if t.alts[j] == nil {
            continue
          }

          wCheck, err := AssertWord(vCheck)
          if err != nil {
            panic(err)
          }

          if wCheck.Value() == w.Value() {
            errCtx := w.Context()
            return nil, errCtx.NewError("Error: duplicate kwarg")
          }
        }
      }

      values[i] = v

      alts[i], err = alt.EvalLazy(tag)
      if err != nil {
        return nil, err
      }

      hadAltBefore = true
    }
  }

  return NewParens(values, alts, t.Context()), nil
}

func (t *Parens) AssertUniqueNames() error {
  for i, v_ := range t.values {
    v, err := AssertWord(v_)
    if err != nil {
      return err
    }

    for j, w_ := range t.values {
      if i == j {
        continue
      }

      w, err := AssertWord(w_)
      if err != nil {
        return err
      }

      if w.Value() == v.Value() {
        errCtx := w.Context()
        return errCtx.NewError("Error: duplicate arg name")
      }

    }
  }

  return nil
}

func (t *Parens) Len() int {
	return len(t.values)
}

func (t *Parens) Loop(fn func(i int, value Token, alt Token) error) error {
	for i, v := range t.values {
		a := t.alts[i]

		if err := fn(i, v, a); err != nil {
			return err
		}
	}

	return nil
}

// only relevant for first token
func (t *Parens) IsSame(other_ Token) bool {
  other, ok := other_.(*Parens)
  if !ok {
    return false
  }

  if len(t.values) == len(other.values) {
    for i, _ := range t.values {
      if !t.values[i].IsSame(other.values[i]) {
        return false
      }

      if t.alts[i] == nil {
        if other.alts[i] != nil {
          return false
        }
      } else if other.alts[i] == nil {
        return false
      } else if !t.alts[i].IsSame(other.alts[i]) {
        return false
      }
    }

    return true
  } else {
    return false
  }
}

func (t *Parens) AnyLazy() bool {
	for i, arg := range t.values {
    alt := t.alts[i]

    if alt == nil {
      if IsLazy(arg) {
        return true
      }
    } else {
      if IsLazy(alt) {
        return true
      }
    }
  }

  return false
}

func IsParens(t Token) bool {
	_, ok := t.(*Parens)
	return ok
}

func AssertParens(t Token) (*Parens, error) {
	p, ok := t.(*Parens)
	if !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Parens, got " + reflect.TypeOf(t).String())
	}

	return p, nil
}
