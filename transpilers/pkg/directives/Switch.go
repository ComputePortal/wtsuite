package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

func buildCase(scope Scope, node Node, swValue tokens.Token,
	tag *tokens.Tag) (bool, error) {
	subScope := NewStatementScope(scope)

	attr, err := tag.Attributes([]string{"value"})
	if err != nil {
		return false, err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return false, err
	}

	value, ok := attr.Get("value")
	if !ok {
		errCtx := attr.Context()
		return false, errCtx.NewError("Error: switch case value not found")
	}

	cond := false
	if swValue != nil {
		eqToken_, eqErr := functions.EQ([]tokens.Token{swValue, value}, tag.Context())
		eqToken, err := tokens.AssertBool(eqToken_)
		if err != nil {
			panic(err)
		}

		if eqErr == nil && eqToken.Value() {
			cond = true
		}
	} else {
		condToken, err := tokens.AssertBool(value)
		if err != nil {
			return false, err
		}

		if condToken.Value() {
			cond = true
		}
	}

	if cond {
		for _, child := range tag.Children() {
			if err := BuildTag(subScope, node, child); err != nil {
				return false, nil
			}
		}
	}

	ft, err := tokens.DictHasFlag(attr, "fallthrough")
	if err != nil {
		return false, err
	}

	if cond && ft {
		cond = false
	}

	return cond, nil
}

func buildDefault(scope Scope, node Node, tag *tokens.Tag) error {
	subScope := NewStatementScope(scope)

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

func buildSwitch(scope Scope, node Node, swValue tokens.Token, tags []*tokens.Tag) error {
	var defaultFound *context.Context = nil // 2 defaults gives an error, case after default gives an error

	var err error
	done := false
	for _, tag := range tags {
		key := tag.Name()
		ctx := tag.Context()

		switch key {
		case "case":
			if defaultFound != nil {
				errCtx := context.MergeContexts(ctx, *defaultFound)
				return errCtx.NewError("Error: default defined before case")
			} else {
				if !done {
					done, err = buildCase(scope, node, swValue, tag)
					if err != nil {
						return err
					}
				}
			}
		case "default":
			if defaultFound != nil {
				return ctx.NewError("Error: default defined more than once")
			} else {
				defaultFound = &ctx
				if !done {
					err = buildDefault(scope, node, tag)
					if err != nil {
						return err
					}
					done = true
				}
			}
		default:
			return ctx.NewError("Error: invalid switch directive")
		}
	}

	return nil
}

func Switch(scope Scope, node Node, tag *tokens.Tag) error {
	subScope := NewStatementScope(scope)

	// value is optional, in which condition needs to be generated for each case
	attr, err := tag.Attributes([]string{"value"})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	value, ok := attr.Get("value")
	if !ok {
		value = nil
	}

	if err := buildSwitch(subScope, node, value, tag.Children()); err != nil {
		return err
	}

	return nil
}

var _switchOk = registerDirective("switch", Switch)
