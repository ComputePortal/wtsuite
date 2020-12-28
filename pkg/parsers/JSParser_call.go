package parsers

import (
	"reflect"

	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildCallArgs(t raw.Token) ([]js.Expression, error) {
	args := make([]js.Expression, 0)

	group, err := raw.AssertParensGroup(t)
	if err != nil {
		return nil, err
	}

	for _, field := range group.Fields {
		arg, err := p.buildExpression(field)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}

	return args, nil
}

func (p *JSParser) buildCastArgs(t raw.Token) ([]js.Expression, error) {
	args := make([]js.Expression, 0)

	group, err := raw.AssertParensGroup(t)
	if err != nil {
		return nil, err
	}

	for i, field := range group.Fields {

    if i < len(group.Fields) -1 {
      arg, err := p.buildExpression(field)
      if err != nil {
        return nil, err
      }
      args = append(args, arg)
    } else {
      arg, err := p.buildTypeExpression(field)
      if err != nil {
        return nil, err
      }
      args = append(args, arg)
    }
	}

	return args, nil
}

func (p *JSParser) buildCallExpression(ts []raw.Token) (js.Expression, error) {
	n := len(ts)

	lhs, err := p.buildExpression(ts[0 : n-1])
	if err != nil {
		return nil, err
	}

	if err := js.AssertCallable(lhs); err != nil {
		return nil, err
	}

  var args []js.Expression = nil

  if ve, ok := lhs.(*js.VarExpression); ok && ve.Name() == "cast" {
    args, err = p.buildCastArgs(ts[n-1])
    if err != nil {
      return nil, err
    }
  } else {
    args, err = p.buildCallArgs(ts[n-1])
    if err != nil {
      return nil, err
    }
  }

	if lhsMember, ok := lhs.(*js.Member); ok {
		if macros.MemberIsClassMacro(lhsMember) {
			return macros.NewClassMacroFromMember(lhsMember, args, lhs.Context())
		}
	}

	if ve, ok := lhs.(*js.VarExpression); ok {
		switch {
		case ve.Name() == "import":
			return p.buildImportDefaultMacro(args, lhs.Context())
		case macros.IsCallMacro(ve.Name()):
			return macros.NewCallMacro(ve.Name(), args, lhs.Context())
		}
	}

	return js.NewCall(lhs, args, lhs.Context()), nil
}

// method call
func (p *JSParser) buildCallStatement(ts_ []raw.Token) (js.Statement, []raw.Token, error) {
	ts, remainingTokens := splitByNextSeparator(ts_, patterns.SEMICOLON)

	if raw.IsWord(ts[0], "void") {
		call, err := p.buildExpression(ts[1:])
		if err != nil {
			return nil, nil, err
		}

		voidStatement := js.NewVoidStatement(call, ts[0].Context())
		if err != nil {
			return nil, nil, err
		}

		return voidStatement, remainingTokens, nil
	} else {
		call_, err := p.buildExpression(ts)
		if err != nil {
			return nil, nil, err
		}

		switch call := call_.(type) {
		case *js.Call:
			return call, remainingTokens, nil
		case *js.Await:
			return call, remainingTokens, nil
		default:
			errCtx := call_.Context()
			err := errCtx.NewError("Error: expected a method call (" + reflect.TypeOf(call_).String() + ")")
			panic(err)
			return nil, nil, err
		}
	}
}
