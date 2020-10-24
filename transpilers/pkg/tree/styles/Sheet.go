package styles

import (
	"strings"
)

type Sheet interface {
	Append(r Rule)
	IsEmpty() bool
	Synchronize() error
	Write() string
	ExpandNested() (Sheet, [][]Rule, error)
}

type SheetData struct {
	rules []Rule
}

func newSheetData() SheetData {
	return SheetData{make([]Rule, 0)}
}

func (ss *SheetData) Append(r Rule) {
	ss.rules = append(ss.rules, r)
}

func (ss *SheetData) Synchronize() error {
	for _, r_ := range ss.rules {
		if rule, ok := r_.(SyncRule); ok {
			if err := rule.Synchronize(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ss *SheetData) Write() string {
	var b strings.Builder

	// assume ids and classes have been synchronized
	for _, rule := range ss.rules {
		b.WriteString(rule.Write(""))
	}

	return b.String()
}

func (ss *SheetData) IsEmpty() bool {
	return len(ss.rules) == 0
}

func (ss *SheetData) ExpandNested() (Sheet, [][]Rule, error) {
	newRules := make([]Rule, 0)
	bundleableRules := make([][]Rule, 0) // the inner order needs to be kept

	for _, r_ := range ss.rules {
		var expanded []Rule
		var err error
		bundleable := false

		// determine basic selectors
		switch rule := r_.(type) {
		case *AllRule:
			expanded, err = rule.ExpandNested(NewSelector("*"))
			bundleable = true
		case *UClassRule:
			expanded, err = rule.ExpandNested(NewSelector("." + rule.class))
		case *HashClassRule:
			expanded, err = rule.ExpandNested(NewSelector("." + rule.hash))
			bundleable = true
		case *TagRule:
			expanded, err = rule.ExpandNested(NewSelector(rule.name))
			if rule.name == "a" {
				bundleable = true
			}
		case *UIDRule:
			expanded, err = rule.ExpandNested(NewSelector("#" + rule.tag.GetID()))
		default:
			panic("unhandled")
		}

		if err != nil {
			return nil, nil, err
		}

		if bundleable {
			bundleableRules = append(bundleableRules, expanded)
		} else {
			newRules = append(newRules, expanded...)
		}
	}

	return &SheetData{newRules}, bundleableRules, nil
}

// for css bundle caching
func WriteBundleRules(br [][]Rule) [][]string {
	res := make([][]string, 0)

	for _, rs := range br {
		subLst := make([]string, 0)

		for _, r := range rs {
			subLst = append(subLst, r.Write(""))
		}

		res = append(res, subLst)
	}

	return res
}
