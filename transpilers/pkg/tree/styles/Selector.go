package styles

import (
	"errors"
	"strings"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type Selector struct {
	key            string
	parents        []string // will be joined by a space (note: parents can only be set once)
	pseudoClasses  []string
	pseudoElements []string
	filters        []string // attribute filters
}

func NewSelector(key string) Selector {
	// start from scratch
	return Selector{key, []string{}, make([]string, 0), make([]string, 0), make([]string, 0)}
}

func (sel *Selector) Copy() Selector {
	parentsCpy := make([]string, len(sel.parents))
	pseudoClassesCpy := make([]string, len(sel.pseudoClasses))
	pseudoElementsCpy := make([]string, len(sel.pseudoElements))
	filtersCpy := make([]string, len(sel.filters))

	copy(parentsCpy, sel.parents)
	copy(pseudoClassesCpy, sel.pseudoClasses)
	copy(pseudoElementsCpy, sel.pseudoElements)
	copy(filtersCpy, sel.filters)

	return Selector{sel.key, parentsCpy, pseudoClassesCpy, pseudoElementsCpy, filtersCpy}
}

// self becomes the next parent
func (sel *Selector) SetChildren(cs []string) error {
	n := len(cs)
	key := cs[n-1]

	ps := append(sel.parents, sel.key)
	ps = append(ps, cs[0:n-1]...)

	sel.key = key
	sel.parents = ps
	return nil
}

func (sel *Selector) AppendPseudoClass(cl string) error {
	for _, pcl := range sel.pseudoClasses {
		if pcl == cl {
			return errors.New("pseudo class already set")
		}
	}

	sel.pseudoClasses = append(sel.pseudoClasses, cl)

	return nil
}

func (sel *Selector) AppendPseudoElement(el string) error {
	for _, pel := range sel.pseudoElements {
		if pel == el {
			return errors.New("pseudo element already set")
		}
	}

	sel.pseudoElements = append(sel.pseudoElements, el)

	return nil
}

func (sel *Selector) AppendFilter(f string) error {
	// TODO: validity checks
	sel.filters = append(sel.filters, f)
	return nil
}

func (sel *Selector) Write() string {
	var b strings.Builder

	b.WriteString(strings.Join(sel.parents, " "))

	if len(sel.parents) > 0 {
		b.WriteString(" ")
	}
	b.WriteString(sel.key)

	for _, f := range sel.filters {
		b.WriteString("[")
		b.WriteString(f)
		b.WriteString("]")
	}

	for _, cl := range sel.pseudoClasses {
		b.WriteString(":")
		b.WriteString(cl)
	}

	for _, el := range sel.pseudoElements {
		b.WriteString("::")
		b.WriteString(el)
	}

	return b.String()
}

/*
 * Selector functions
 */

func expandNested_(sel Selector, v tokens.Token) ([]Rule, error) {
	if tokens.IsNull(v) {
		return []Rule{}, nil
	}

	attr, err := tokens.AssertStringDict(v)
	if err != nil {
		return nil, err
	}

	leafAttr, rules, err := expandNested(attr, sel)
	if err != nil {
		return nil, err
	}

	if len(leafAttr) > 0 {
		rules = append([]Rule{NewSelectorRule(sel, leafAttr)}, rules...)
	}
	return rules, nil
}

func appendPseudoClass(name string, sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	if err := sel.AppendPseudoClass(name); err != nil {
		return nil, ctx.NewError("Error: " + err.Error())
	}

	return expandNested_(sel, v)
}

func Hover(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("hover", sel, args, v, ctx)
}

func Visited(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("visited", sel, args, v, ctx)
}

func Checked(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("checked", sel, args, v, ctx)
}

func Active(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("active", sel, args, v, ctx)
}

func Focus(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("focus", sel, args, v, ctx)
}

func Invalid(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("invalid", sel, args, v, ctx)
}

func Disabled(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("disabled", sel, args, v, ctx)
}

func Valid(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("valid", sel, args, v, ctx)
}

func FirstChild(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("first-child", sel, args, v, ctx)
}

func LastChild(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("last-child", sel, args, v, ctx)
}

func OnlyChild(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoClass("only-child", sel, args, v, ctx)
}

func NthChild(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) < 1 {
		return nil, ctx.NewError("Error: expected at least 1 argument")
	}

	var b strings.Builder
	b.WriteString("nth-child(")
	for _, arg := range args {
		b.WriteString(arg)
	}
	b.WriteString(")")

	return appendPseudoClass(b.String(), sel, []string{}, v, ctx)
}

func Filter(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if err := sel.AppendFilter(args[0]); err != nil {
		return nil, ctx.NewError("Error: " + err.Error())
	}

	return expandNested_(sel, v)
}

func appendPseudoElement(name string, sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	if err := sel.AppendPseudoElement(name); err != nil {
		return nil, ctx.NewError("Error: " + err.Error())
	}

	return expandNested_(sel, v)
}

func After(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoElement("after", sel, args, v, ctx)
}

func Placeholder(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoElement("placeholder", sel, args, v, ctx)
}

func Before(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoElement("before", sel, args, v, ctx)
}

func WebkitScrollbar(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	return appendPseudoElement("-webkit-scrollbar", sel, args, v, ctx)
}

func Parents(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) == 0 {
		return nil, ctx.NewError("Error: expected at least 1 argument")
	}

	allRules := []Rule{}

	// inject parents at each possible location
	for i := 0; i <= len(sel.parents); i++ {
		cpy := sel.Copy()

		bef := cpy.parents[0:i]
		newParents := make([]string, len(bef))
		copy(newParents, bef) // otherwise something strange happens to the cpy.parents[0:i] slice
		aft := cpy.parents[i:]

		for _, arg := range args {
			newParents = append(newParents, arg)
		}

		for _, af := range aft {
			newParents = append(newParents, af)
		}

		cpy.parents = newParents

		cpyRules, err := expandNested_(cpy, v)
		if err != nil {
			return nil, err
		}

		allRules = append(allRules, cpyRules...)
	}

	return allRules, nil
}

func Child(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) == 0 {
		return nil, ctx.NewError("Error: expected at least 1 argument")
	}

	//if len(sel.parents) != 0 {
	//return nil, ctx.NewError("Error: cannot combine @child with another @child/@parents")
	//}

	if err := sel.SetChildren(args); err != nil {
		return nil, ctx.NewError("Error: " + err.Error())
	}

	return expandNested_(sel, v)
}

// dispatch by preceding '.', doesnt have any args, can't contain spaces
// XXX: are we sure that no css attribute names begin with a '.'?
func NamedClass(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	// args[0] includes the dot !
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected exactly one token")
	}

	if err := sel.SetChildren(args); err != nil {
		return nil, ctx.NewError("Error: " + err.Error())
	}

	return expandNested_(sel, v)
}

var _activeOk = registerAtFunction("active", Active)
var _afterOk = registerAtFunction("after", After)
var _beforeOk = registerAtFunction("before", Before)
var _checkedOk = registerAtFunction("checked", Checked)
var _childOk = registerAtFunction("child", Child)
var _filterOk = registerAtFunction("filter", Filter)
var _firstChildOk = registerAtFunction("first-child", FirstChild)
var _focusOk = registerAtFunction("focus", Focus)
var _hoverOk = registerAtFunction("hover", Hover)
var _invalidOk = registerAtFunction("invalid", Invalid)
var _disabledOk = registerAtFunction("disabled", Disabled)
var _lastChildOk = registerAtFunction("last-child", LastChild)
var _nthChildOk = registerAtFunction("nth-child", NthChild)
var _onlyChildOk = registerAtFunction("only-child", OnlyChild)
var _parentsOk = registerAtFunction("parents", Parents)
var _placeholderOk = registerAtFunction("placeholder", Placeholder)
var _validOk = registerAtFunction("valid", Valid)
var _visitedOk = registerAtFunction("visited", Visited)
var _webkitScrollbarOk = registerAtFunction("-webkit-scrollbar", WebkitScrollbar)
