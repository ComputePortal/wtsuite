package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func (p *JSParser) buildInterfaceExtendsExpression(ts []raw.Token) ([]*js.VarExpression, []raw.Token, error) {
  parents := make([]*js.VarExpression, 0)

	if raw.IsWord(ts[0], "extends") {
    needExpr := true

    for needExpr {
      if len(ts) < 2 {
        errCtx := raw.MergeContexts(ts...)
        return nil, nil, errCtx.NewError("Error: bad interface extends")
      }

      var condensedNameToken *raw.Word
      var err error
      condensedNameToken, ts, err = p.condensePackagePeriods(ts[1:]) // shortens ts by 1 or more
      if err != nil {
        return nil, nil, err
      }

      parent, err := p.buildVarExpression(condensedNameToken)
      if err != nil {
        return nil, nil, err
      }

      parents = append(parents, parent)

      if raw.IsSymbol(ts[0], patterns.COMMA) {
        needExpr = true
      } else {
        needExpr = false
      }
    }
	} 

  return parents, ts, nil
}

func (p *JSParser) buildInterface(ts []raw.Token) (*js.Interface, error) {
	interfCtx := raw.MergeContexts(ts...)

	if len(ts) == 2 && raw.IsBracesGroup(ts[1]) {
		errCtx := interfCtx
		return nil, errCtx.NewError("Error: missing interface name")
	} else if len(ts) < 3 {
		errCtx := interfCtx
		return nil, errCtx.NewError("Error: bad interface definition")
	}

  isRPC := false
  if raw.IsWord(ts[0], "rpc") {
    isRPC = true
    ts = ts[1:]
  }

	clType, ts, err := p.buildClassOrExtendsTypeExpression(ts[1:])
	if err != nil {
		return nil, err
	}

  var parents []*js.VarExpression
  parents, ts, err = p.buildInterfaceExtendsExpression(ts)
  if err != nil {
    return nil, err
  }

	classInterface, err := js.NewInterface(clType, parents, isRPC, interfCtx)
	if err != nil {
		return nil, err
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[len(ts)-1])
	if err != nil {
		return nil, err
	}

	if bracesGroup.IsComma() {
		errCtx := bracesGroup.Context()
		return nil, errCtx.NewError("Error: interface uses semicolon separator")
	}

	for _, field := range bracesGroup.Fields {
		if len(field) == 0 {
			continue
		}

		fi, remaining, err := p.buildFunctionInterface(field, true, interfCtx)
		if err != nil {
			return nil, err
		}

		if fi.Role() != prototypes.NORMAL &&
			fi.Role() != prototypes.GETTER &&
			fi.Role() != prototypes.SETTER &&
			fi.Role() != prototypes.ASYNC {
			errCtx := fi.Context()
			return nil, errCtx.NewError("Error: illegal interface function role(s)")
		}

		if len(remaining) != 0 {
			errCtx := raw.MergeContexts(remaining...)
			return nil, errCtx.NewError("Error: unexpected tokens (hint: did forget a semicolon?)")
		}

		if err := classInterface.AddMember(fi); err != nil {
			return nil, err
		}
	}

	return classInterface, nil
}

func (p *JSParser) buildInterfaceStatement(ts []raw.Token) (*js.Interface, []raw.Token, error) {
	for i, t := range ts {
		if raw.IsBracesGroup(t) {
			statement, err := p.buildInterface(ts[0 : i+1])
			if err != nil {
				return nil, nil, err
			}

			remaining := stripSeparators(i+1, ts, patterns.SEMICOLON)

			return statement, remaining, nil
		}
	}

	errCtx := raw.MergeContexts(ts...)
	return nil, nil, errCtx.NewError("Error: no interface body found")
}
