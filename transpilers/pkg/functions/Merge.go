package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func modifyList(a tokens.Token, b tokens.Token, ctx context.Context) (tokens.Token, error) {
	list, err := tokens.AssertList(a)
	if err != nil {
		panic(err)
	}

	dict, err := tokens.AssertIntDict(b)
	if err != nil {
		panic(err)
	}

	values := list.GetTokens()

	n := len(values)

	// negative indices append to end
	// indices that are larger than n, also append to the end (but before the implicit append, and the indices must be exact!)
	// holes are not allowed
	largestExplicitAppendIndex := 0  // from 0
	smallestImplicitAppendIndex := 0 // to -1

	explicitAppend := make(map[int]struct {
		key   *tokens.Int
		value tokens.Token
	})
	implicitAppend := make(map[int]struct {
		key   *tokens.Int
		value tokens.Token
	})

	if err := dict.Loop(func(k *tokens.Int, v tokens.Token, last bool) error {
		i := k.Value()
		if i >= 0 && i < n {
			values[i] = v
		} else if i < 0 {
			if i < smallestImplicitAppendIndex {
				smallestImplicitAppendIndex = i
			}
			implicitAppend[i] = struct {
				key   *tokens.Int
				value tokens.Token
			}{k, v}
		} else if i >= n {
			if i > largestExplicitAppendIndex {
				largestExplicitAppendIndex = i
			}
			explicitAppend[i] = struct {
				key   *tokens.Int
				value tokens.Token
			}{k, v}
		} else {
			panic("algo error")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// check for holes
	for i := n; i <= largestExplicitAppendIndex; i++ {
		if kv, ok := explicitAppend[i]; !ok {
			errCtx := kv.key.Context()
			return nil, errCtx.NewError("Error: can't have holes when modifying list with dict")
		} else {
			values = append(values, kv.value)
		}
	}

	for i := smallestImplicitAppendIndex; i <= -1; i++ {
		if kv, ok := implicitAppend[i]; !ok {
			errCtx := kv.key.Context()
			return nil, errCtx.NewError("Error: can't have holes when modifying list with dict")
		} else {
			values = append(values, kv.value)
		}
	}

	return tokens.NewValuesList(values, ctx), nil
}

func MergeStringDictsInplace(res *tokens.StringDict, new *tokens.StringDict, ctx context.Context) error {
	if err := new.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		oldValue, ok := res.Get(key)
		switch {
		case !ok:
			res.Set(key, value)
		case !tokens.IsDict(value):
			res.Set(key, value)
		default:
			switch {
			case tokens.IsDict(oldValue):
				newValue, err := Merge([]tokens.Token{oldValue, value}, ctx)
				if err != nil {
					return err
				}
				res.Set(key, newValue)
			case tokens.IsList(oldValue) && tokens.IsIntDict(value):
				newValue, err := modifyList(oldValue, value, ctx)
				if err != nil {
					return err
				}
				res.Set(key, newValue)
			default:
				res.Set(key, value)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func mergeStringDicts(old *tokens.StringDict, new *tokens.StringDict, ctx context.Context) (tokens.Token, error) {
	res, err := old.CopyStringDict(ctx)
	if err != nil {
		return nil, err
	}

	if err := MergeStringDictsInplace(res, new, ctx); err != nil {
		return nil, err
	}

	return res, nil
}

func mergeIntDicts(old *tokens.IntDict, new *tokens.IntDict, ctx context.Context) (tokens.Token, error) {
	res, err := old.CopyIntDict(ctx)
	if err != nil {
		return nil, err
	}

	if err := new.Loop(func(key *tokens.Int, value tokens.Token, last bool) error {
		oldValue, ok := res.Get(key)
		switch {
		case !ok:
			res.Set(key, value)
		case !tokens.IsDict(value):
			res.Set(key, value)
		default:
			switch {
			case tokens.IsDict(oldValue):
				newValue, err := Merge([]tokens.Token{oldValue, value}, ctx)
				if err != nil {
					return err
				}
				res.Set(key, newValue)
			case tokens.IsList(oldValue) && tokens.IsIntDict(value):
				newValue, err := modifyList(oldValue, value, ctx)
				if err != nil {
					return err
				}
				res.Set(key, newValue)
			default:
				res.Set(key, value)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func mergeTwo(a, b tokens.Token, ctx context.Context) (tokens.Token, error) {
	// must deep copy because internal values might change
	switch {
	case tokens.IsStringDict(a) && tokens.IsStringDict(b):
		old, err := tokens.AssertStringDict(a)
		if err != nil {
			panic(err)
		}

		new, err := tokens.AssertStringDict(b)
		if err != nil {
			panic(err)
		}
		return mergeStringDicts(old, new, ctx)
	case tokens.IsIntDict(a) && tokens.IsIntDict(b):
		old, err := tokens.AssertIntDict(a)
		if err != nil {
			panic(err)
		}

		new, err := tokens.AssertIntDict(b)
		if err != nil {
			panic(err)
		}
		return mergeIntDicts(old, new, ctx)
	default:
		return nil, ctx.NewError("Error: expected similar dicts")
	}
}

func Merge(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	// must deep copy because internal values might change
	if len(args) < 2 {
		return nil, ctx.NewError("Error: expected at least 2 arguments")
	}

	var result tokens.Token = nil

	for i := 1; i < len(args); i++ {
		var err error
		if i == 1 {
			result, err = mergeTwo(args[i-1], args[i], ctx)
		} else {
			result, err = mergeTwo(result, args[i], ctx)
		}

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
