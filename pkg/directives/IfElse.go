package directives

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func buildConditional(scope Scope, node Node, tag *tokens.Tag) (bool, error) {
	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{"cond"})
	if err != nil {
		return false, err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return false, err
	}

	cond, err := tokens.DictBool(attr, "cond")
	if err != nil {
		return false, err
	}

	if cond.Value() {
		for _, child := range tag.Children() {
			if err := BuildTag(subScope, node, child); err != nil {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}

func buildElse(scope Scope, node Node, tag *tokens.Tag) error {
	subScope := NewSubScope(scope)

	if err := tag.AssertNoAttributes(); err != nil {
		return err
	}

	for _, child := range tag.Children() {
		if err := BuildTag(subScope, node, child); err != nil {
			return err
		}
	}

	return nil
}

func buildIfElse(scope Scope, node Node, tags []*tokens.Tag) error {
	var ifFound *context.Context = nil   // no if gives an error
	var elseFound *context.Context = nil // two or more elses give an error

	var err error
	done := false

	for _, tag := range tags {
		key := tag.Name()
		ctx := tag.Context()

		switch key {
		case "if":
			if ifFound != nil {
				errCtx := context.MergeContexts(ctx, *ifFound)
				return errCtx.NewError("Error: if already defined")
			} else {
				ifFound = &ctx
				if !done {
					done, err = buildConditional(scope, node, tag)
					if err != nil {
						return err
					}
				}
			}
		case "elseif":
			if ifFound == nil {
				return ctx.NewError("Error: if not yet defined")
			} else {
				if !done {
					done, err = buildConditional(scope, node, tag)
					if err != nil {
						return err
					}
				}
			}
		case "else":
			if ifFound == nil {
				return ctx.NewError("Error: if not yet defined")
			} else if elseFound != nil {
				errCtx := context.MergeContexts(ctx, *elseFound)
				return errCtx.NewError("Error: else already defined")
			} else {
				elseFound = &ctx

				if !done {
					err = buildElse(scope, node, tag)
					if err != nil {
						return err
					}
					done = true
				}
			}
		default:
			return ctx.NewError("Error: invalid ifelse directive")
		}
	}

	return nil
}

func IfElse(scope Scope, node Node, tag *tokens.Tag) error {
	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{})
	if err != nil {
		return err
	}

	if attr.Len() > 0 {
		errCtx := attr.Context()
		return errCtx.NewError("Error: unexpected attributes")
	}

	return buildIfElse(subScope, node, tag.Children())
}

var _ifElseOk = registerDirective("ifelse", IfElse)
