package js

import (
	"./values"

	"../context"
)

type BranchStack struct {
	typeGuards map[interface{}]values.Interface
	values.StackData
}

func NewBranchStack(parent values.Stack) *BranchStack {
	return &BranchStack{nil, values.NewStackData(parent)}
}

func NewTypeGuardBranchStack(typeGuards map[interface{}]values.Interface, parent values.Stack) *BranchStack {
	return &BranchStack{typeGuards, values.NewStackData(parent)}
}

func (s *BranchStack) GetValue(ptr interface{}, ctx context.Context) (values.Value, error) {
	v, err := s.StackData.GetValue(ptr, ctx)
	if err != nil {
		return nil, err
	}

	// change InstanceInterface if necessary
	if s.typeGuards != nil {
		if interf, ok := s.typeGuards[ptr]; ok {
			vPos := values.UnpackMulti([]values.Value{v})

			// first filter out everything that doesn't match
			vRes := make([]values.Value, 0)
			for _, vp := range vPos {
				if vr, ok := vp.ChangeInstanceInterface(interf, false, false); ok {
					vRes = append(vRes, vr)
				}
			}

			switch len(vRes) {
			case 0:
				// null isnt actually used, because branch is never taken
				v = values.NewAllNull(v.Context())
			case 1:
				v = vRes[0]
			default:
				v = values.NewMulti(vRes, v.Context())
			}
		}
	}

	return v, nil
}

func (s *BranchStack) SetValue(ptr interface{}, v values.Value,
	allowBranching bool, ctx context.Context) error {
	if allowBranching && s.HasValue(ptr) {
		// prevOwner stays prevOwner and this owner is ignored
		prev, err := s.GetValue(ptr, ctx)
		if err != nil {
			panic(err)
		}

		if prev != nil {
			v = values.NewMulti([]values.Value{prev, v}, v.Context())
		}
	}

	// delete typeGuard
	if s.typeGuards != nil {
		if _, ok := s.typeGuards[ptr]; ok {
			delete(s.typeGuards, ptr)
		}
	}

	// allowBranching set to false because it would be redundant
	return s.StackData.SetValue(ptr, v, false, ctx)
}
