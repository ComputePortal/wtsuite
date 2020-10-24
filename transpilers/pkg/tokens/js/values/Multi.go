package values

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"../../context"
)

type Multi struct {
	values []Value
	ValueData
}

func UnpackMulti(vs []Value) []Value {
	res := make([]Value, 0)

	for _, v_ := range vs {
		v_ = UnpackContextValue(v_)
		switch v := v_.(type) {
		case *Multi:
			// only append the unique
			for _, val := range v.values {
				unique := true
				for i, test := range res {
					// XXX: is this actually worth it?
					// just test for similarity?
					if vres := test.Merge(val); vres != nil {
						unique = false
						res[i] = vres
						break
					}
				}
				if unique {
					res = append(res, val)
				}
			}
			//res = append(res, v.values...)
		default:
			if len(vs) == 1 {
				// to avoid a merge from removing literal data
				return []Value{v}
			}

			unique := true
			for i, test := range res {
				if vres := test.Merge(v); vres != nil {
					unique = false
					res[i] = vres
					break
				}
			}
			if unique {
				res = append(res, v)
			}
		}
	}

	return res
}

func NewMulti(vs []Value, ctx context.Context) Value {
	vs = UnpackMulti(vs)

	// remove anything that is Null first
	hasAllNull := false
	otherNulls := make([]Value, 0)
	nonNulls := make([]Value, 0)

	addNull := func(v *Null) {
		for _, otherNull := range otherNulls {
			if otherNull.Merge(v) != nil {
				return
			}
		}

		otherNulls = append(otherNulls, v)
	}

	addNonNull := func(v Value) {
		/*for _, otherNonNull := range nonNulls {
			if otherNonNull.Merge(v) != nil {
				return
			}
		}*/

		nonNulls = append(nonNulls, v)
	}

	for _, v_ := range vs {
		switch v := v_.(type) {
		case *Null:
			if v.proto == nil {
				hasAllNull = true
			} else {
				addNull(v)
			}
		default:
			addNonNull(v)
		}
	}

	// remove the otherNulls that have a proto in the nonNulls
	filteredOtherNulls := make([]Value, 0)
Outer:
	for _, otherNull := range otherNulls {
		otherNullProto, _ := otherNull.GetNullPrototype()
		for _, nonNull := range nonNulls {
			if nonNullProto, ok := nonNull.GetInstancePrototype(); ok {
				if otherNullProto == nonNullProto {
					continue Outer
				}
			}
		}

		filteredOtherNulls = append(filteredOtherNulls, otherNull)
	}

	otherNulls = filteredOtherNulls

	if hasAllNull && len(otherNulls) == 0 && len(nonNulls) == 0 {
		return NewAllNull(ctx)
	} else if len(otherNulls) == 1 && len(nonNulls) == 0 {
		return otherNulls[0]
	} else if len(otherNulls) == 0 && len(nonNulls) == 1 {
		return nonNulls[0]
	} else {
		nonNulls = append(nonNulls, otherNulls...)
		if len(nonNulls) == 0 {
			err := ctx.NewError("Internal Error: empty multi shouldn't be possible")
			panic(err)
		}
		return &Multi{nonNulls, ValueData{ctx}}
	}
}

func (v *Multi) TypeName() string {
	typeNames := make(map[string]string)

	for _, val := range v.values {
		tn := val.TypeName()
		typeNames[tn] = tn
	}

	var b strings.Builder
	i := 0
	for k, _ := range typeNames {
		b.WriteString(k)

		if i < len(typeNames)-1 {
			b.WriteString("|") // or other
		}

		i += 1
	}

	return b.String()
}

func (v *Multi) Copy(cache CopyCache) Value {
	values := copyValueList(v.values, cache)

	return &Multi{values, ValueData{v.Context()}}
}

func (v *Multi) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	values := make([]Value, len(v.values))

	for i, val := range v.values {
		castVal, err := val.Cast(ntype, ctx)
		if err != nil {
			return nil, err
		}

		values[i] = castVal
	}

	return NewMulti(values, ctx), nil
}

func (v *Multi) IsInstanceOf(ps ...Prototype) bool {
	for _, val := range v.values {
		if !val.IsInstanceOf(ps...) {
			return false
		}
	}

	return true
}

func (v *Multi) MaybeInstanceOf(p Prototype) bool {
	for _, val := range v.values {
		if val.MaybeInstanceOf(p) {
			return true
		}
	}

	return false
}

func (v *Multi) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*Multi)
	if !ok {
		return nil
	}

	vs := mergeValueLists(v.values, other.values)
	if vs == nil {
		return nil
	} else {
		return NewMulti(vs, v.Context())
	}
}

func (v *Multi) LoopNestedPrototypes(fn func(Prototype)) {
	for _, v := range v.values {
		v.LoopNestedPrototypes(fn)
	}
}

func (v *Multi) RemoveLiteralness(all bool) Value {
	result := make([]Value, len(v.values))

	for i, val := range v.values {
		result[i] = val.RemoveLiteralness(all)
	}

	// literals would'be been merged with non-literals, so it is pointless to call the multi constructor
	return &Multi{result, ValueData{v.Context()}}
}

func (v *Multi) EvalFunction(stack Stack, args []Value, ctx context.Context) (Value, error) {
	result := make([]Value, len(v.values))

	for i, val := range v.values {
		res, err := val.EvalFunction(stack, args, ctx)
		if err != nil {
			return nil, err
		}

		result[i] = res
	}

	return NewMulti(result, ctx), nil
}

// all nil, or all values
func (v *Multi) EvalFunctionNoReturn(stack Stack, args []Value, ctx context.Context) (Value, error) {
	result := make([]Value, len(v.values))

	allNil := false
	for i, val := range v.values {
		res, err := val.EvalFunctionNoReturn(stack, args, ctx)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			allNil = (res == nil)
		} else if allNil && (res != nil) {
			return nil, ctx.NewError("Error: some return values are void, others are non-void")
		}

		result[i] = res
	}

	return NewMulti(result, ctx), nil
}

func (v *Multi) EvalMethod(stack Stack, args []Value, ctx context.Context) error {
	if (VERBOSITY >= 3 && len(v.values) > 6) || (VERBOSITY >= 2 && len(v.values) > 20) {
		warningCtx := ctx
		fmt.Fprintf(os.Stderr, warningCtx.NewError("Warning: evaluating multi method "+strconv.Itoa(len(v.values))+" times").Error())
	}

	for _, val := range v.values {
		if err := val.EvalMethod(stack, args, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *Multi) EvalConstructor(stack Stack, args []Value, ctx context.Context) (Value, error) {
	result := make([]Value, len(v.values))

	for i, val := range v.values {
		res, err := val.EvalConstructor(stack, args, ctx)
		if err != nil {
			return nil, err
		}

		result[i] = res
	}

	return NewMulti(result, ctx), nil
}

func (v *Multi) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	for _, val := range v.values {
		err := val.EvalAsEntryPoint(stack, ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v *Multi) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	result := make([]Value, len(v.values))

	for i, val := range v.values {
		res, err := val.GetMember(stack, key, includePrivate, ctx)
		if err != nil {
			return nil, err
		}

		result[i] = res
	}

	return NewMulti(result, ctx), nil
}

func (v *Multi) SetMember(stack Stack, key string, value Value, includePrivate bool,
	ctx context.Context) error {
	for _, val := range v.values {
		if err := val.SetMember(stack, key, value, includePrivate, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *Multi) GetIndex(stack Stack, index Value, ctx context.Context) (Value, error) {
	result := make([]Value, len(v.values))

	for i, val := range v.values {
		res, err := val.GetIndex(stack, index, ctx)
		if err != nil {
			return nil, err
		}

		result[i] = res
	}

	return NewMulti(result, ctx), nil
}

func (v *Multi) SetIndex(stack Stack, index Value, value Value, ctx context.Context) error {
	for _, val := range v.values {
		if err := val.SetIndex(stack, index, value, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (v *Multi) LoopForOf(fn func(Value) error, ctx context.Context) error {
	vs := make([]Value, 0)

	fnCollect := func(v_ Value) error {
		vs = append(vs, v_)
		return nil
	}

	for _, v_ := range v.values {
		if err := v_.LoopForOf(fnCollect, ctx); err != nil {
			return err
		}
	}

	return fn(NewMulti(vs, v.Context()))
}

func (v *Multi) LoopForIn(fn func(Value) error, ctx context.Context) error {
	vs := make([]Value, 0)

	fnCollect := func(v_ Value) error {
		vs = append(vs, v_)
		return nil
	}

	for _, v_ := range v.values {
		if err := v_.LoopForIn(fnCollect, ctx); err != nil {
			return err
		}
	}

	return fn(NewMulti(vs, v.Context()))
}

func (v *Multi) IsClass() bool {
	for _, val := range v.values {
		if !val.IsClass() {
			return false
		}
	}

	return true
}

func (v *Multi) IsFunction() bool {
	for _, val := range v.values {
		if !val.IsFunction() {
			return false
		}
	}

	return true
}

func (v *Multi) IsInstance() bool {
	for _, val := range v.values {
		if !val.IsInstance() {
			return false
		}
	}

	return true
}

func (v *Multi) IsNull() bool {
	for _, val := range v.values {
		if !val.IsNull() {
			return false
		}
	}

	return true
}

func (v *Multi) IsVoid() bool {
	for _, val := range v.values {
		if !val.IsVoid() {
			return false
		}
	}

	return len(v.values) != 0
}

func (v *Multi) IsInterface() bool {
	for _, val := range v.values {
		if !val.IsInterface() {
			return false
		}
	}

	return true
}

// only same prototype, not common ancestor
func (v *Multi) GetClassPrototype() (Prototype, bool) {
	var p Prototype = nil

	for _, val := range v.values {
		p_, ok := val.GetClassPrototype()
		if !ok {
			return nil, false
		} else {
			if p == nil {
				p = p_
			} else if p != p_ {
				return nil, false
			}
		}
	}

	if p == nil {
		return nil, false
	} else {
		return p, true
	}
}

func (v *Multi) GetClassInterface() (Interface, bool) {
	var interf Interface = nil

	for _, val := range v.values {
		interf_, ok := val.GetClassInterface()
		if !ok {
			return nil, false
		} else {
			if interf == nil {
				interf = interf_
			} else if interf != interf_ {
				return nil, false
			}
		}
	}

	if interf == nil {
		return nil, false
	} else {
		return interf, true
	}
}

func (v *Multi) GetInstancePrototype() (Prototype, bool) {
	var p Prototype = nil

	for _, val := range v.values {
		p_, ok := val.GetInstancePrototype()
		if !ok {
			return nil, false
		} else {
			if p == nil {
				p = p_
			} else if p != p_ {
				return nil, false
			}
		}
	}

	return p, true
}

func (v *Multi) GetNullPrototype() (Prototype, bool) {
	var p Prototype = nil

	for _, val := range v.values {
		p_, ok := val.GetNullPrototype()
		if !ok {
			return nil, false
		} else {
			if p == nil {
				p = p_
			} else if p != p_ {
				return nil, false
			}
		}
	}

	return p, true
}

func (v *Multi) ChangeInstancePrototype(p Prototype, inPlace bool) (Value, bool) {
	vs := make([]Value, 0)

	for _, v_ := range v.values {
		newV, ok := v_.ChangeInstancePrototype(p, inPlace)
		if !ok {
			return nil, false
		}

		vs = append(vs, newV)
	}

	return NewMulti(vs, v.Context()), true
}

func (v *Multi) ChangeInstanceInterface(interf Interface, inPlace bool, checkOuter bool) (Value, bool) {
	vs := make([]Value, 0)

	for _, v_ := range v.values {
		newV, ok := v_.ChangeInstanceInterface(interf, inPlace, checkOuter)
		if !ok {
			return nil, false
		}

		vs = append(vs, newV)
	}

	return NewMulti(vs, v.Context()), true
}

func (v *Multi) Dump() string {
	var b strings.Builder

	for _, v_ := range v.values {
		b.WriteString(reflect.TypeOf(v_).String())
		b.WriteString(", ")
		b.WriteString(v_.TypeName())
		b.WriteString("\n")
	}

	return b.String()
}
