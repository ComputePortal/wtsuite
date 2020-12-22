package styles

// wraps rules used by js constructors
type UIDToUClassSheet struct {
	s DocSheet
}

func NewUIDToUClassSheet(s DocSheet) *UIDToUClassSheet {
	return &UIDToUClassSheet{s}
}

func (ss *UIDToUClassSheet) Append(r Rule) {
	switch rt := r.(type) {
	case *UIDRule:
		ss.s.Append(rt.ToUClassRule())
	case *HashClassRule:
		ss.s.Append(rt.ToUClassRule())
	default:
		ss.s.Append(r)
	}
}

func (ss *UIDToUClassSheet) Synchronize() error {
	return ss.s.Synchronize()
}

func (ss *UIDToUClassSheet) Write() string {
	return ss.s.Write()
}

func (ss *UIDToUClassSheet) IsEmpty() bool {
	return ss.s.IsEmpty()
}

func (ss *UIDToUClassSheet) ExpandNested() (Sheet, [][]Rule, error) {
	return ss.s.ExpandNested()
}
