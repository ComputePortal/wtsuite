package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	//"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

func AssertValidVar(nameToken *tokens.String) error {
	errCtx := nameToken.Context()

	if !patterns.IsValidVar(nameToken.Value()) {
		return errCtx.NewError("Error: invalid var name")
	} else {
		return nil
	}
}

// doesnt change the node
func AddVar(scope Scope, node Node, tag *tokens.Tag) error {
	if !tag.IsEmpty() {
		errCtx := tag.Context()
		return errCtx.NewError("Error: unexpected child tags of var directive")
	}

	subScope := NewSubScope(scope)

	attr, err := tag.Attributes([]string{})
	if err != nil {
		return err
	}

	attr, err = attr.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	n := attr.Len()
	if n != 1 && n != 2 {
		errCtx := attr.Context()
		return errCtx.NewError("Error: expected 1 or 2 var attributes")
	}

	isExported := false
	if n == 2 {
		isExported, err = tokens.DictHasFlag(attr, "export")
		if err != nil {
			return err
		}

		if !isExported {
			errCtx := attr.Context()
			return errCtx.NewError("Error: if 2 attributes, expected export flag")
		}
	}

	var nameToken *tokens.String = nil
	var valueToken tokens.Token = nil
	if err := attr.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		if key.Value() != "export" {
			nameToken = key
			valueToken = value
		}

		return nil
	}); err != nil {
		panic(err)
	}

	if err := AssertValidVar(nameToken); err != nil {
		return err
	}

	constant, err := tokens.DictHasFlag(attr, "constant")
	if err != nil {
		return err
	}

	key := nameToken.Value()
	switch {
	case HasGlobal(key):
		errCtx := nameToken.InnerContext()
		err := errCtx.NewError("Error: cant redefine global")
		return err
	case scope.HasVar(key):
		v := scope.GetVar(key)
		errCtx := nameToken.InnerContext() // unless context is identical, in which case we are probably inside a for loop
		if !(!errCtx.Less(&v.Ctx) && !v.Ctx.Less(&errCtx)) {
			err := errCtx.NewError("Error: cant redefine var")
			err.AppendContextString("Info: defined here", v.Ctx)
			return err
		}

		fallthrough
	default:
		v := functions.Var{valueToken, constant, false, false, isExported,
			nameToken.InnerContext()}
    if err := scope.SetVar(key, v); err != nil {
      return err
    }
	}

	return nil
}

var _addVarOk = registerDirective("var", AddVar)
