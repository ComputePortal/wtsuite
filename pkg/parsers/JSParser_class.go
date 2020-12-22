package parsers

import (
	"../tokens/js"
	"../tokens/patterns"
	"../tokens/raw"
)

func (p *JSParser) buildClassExtendsExpression(ts []raw.Token) (*js.TypeExpression, []raw.Token, error) {
	var extends *js.TypeExpression = nil

	if raw.IsWord(ts[0], "extends") {
		if len(ts) < 2 {
			errCtx := raw.MergeContexts(ts...)
			return nil, nil, errCtx.NewError("Error: bad class extends definition")
		}

		var err error
		extends, ts, err = p.buildClassOrExtendsTypeExpression(ts[1:])
		if err != nil {
			panic(err)
		}
	}

	return extends, ts, nil
}

func (p *JSParser) buildClassImplementsExpression(ts []raw.Token) ([]*js.VarExpression, []raw.Token, error) {
  impls := make([]*js.VarExpression, 0)

	if raw.IsWord(ts[0], "implements") {
    needImpl := true

    for needImpl {
      if len(ts) < 2 {
        errCtx := raw.MergeContexts(ts...)
        return nil, nil, errCtx.NewError("Error: bad class implements interface")
      }

      var condensedNameToken *raw.Word
      var err error
      condensedNameToken, ts, err = p.condensePackagePeriods(ts[1:]) // shortens ts by 1 or more
      if err != nil {
        return nil, nil, err
      }

      implements, err := p.buildVarExpression(condensedNameToken)
      if err != nil {
        return nil, nil, err
      }

      impls = append(impls, implements)

      if raw.IsSymbol(ts[0], patterns.COMMA) {
        needImpl = true
      } else {
        needImpl = false
      }
    }
	} 

  return impls, ts, nil
}

func (p *JSParser) buildClassUniversalName(ts []raw.Token) (string, []raw.Token, error) {
	if raw.IsWord(ts[0], "universe") {
		if len(ts) < 2 {
			errCtx := raw.MergeContexts(ts...)
			return "", nil, errCtx.NewError("Error: bad class universe definition")
		}

		nameToken, err := raw.AssertWord(ts[1])
		if err != nil {
			return "", nil, err
		}

		name := nameToken.Value()
		return name, ts[2:], nil
	}

	return "", ts, nil
}

func (p *JSParser) buildClass(ts []raw.Token) (*js.Class, error) {
	clCtx := raw.MergeContexts(ts...)

	if len(ts) < 2 {
		errCtx := clCtx
		return nil, errCtx.NewError("Error: bad class definition")
	}

  isAbstract := false
  isFinal := false
  if raw.IsWord(ts[0], "abstract") {
    isAbstract = true
    ts = ts[1:]
  } else if raw.IsWord(ts[0], "final") {
    isFinal = true
    ts = ts[1:]
  }

  if !raw.IsWord(ts[0], "class") {
    errCtx := ts[0].Context()
    // should've actually been caught earlier
    return nil, errCtx.NewError("Error: expected class keyword")
  }

	// special, because classes dont necessarily have a name
	var clType *js.TypeExpression
	var err error = nil
	if raw.IsAnyWord(ts[1]) && len(ts) > 2 && raw.IsAngledGroup(ts[2]) {
		clType, err = p.buildTypeExpression(ts[1:3])
		ts = ts[3:]
	} else if raw.IsAnyWord(ts[1]) {
		clType, err = p.buildTypeExpression(ts[1:2])
		ts = ts[2:]
	}
	if err != nil {
		return nil, err
	}

	// TODO: add container typing parsing here

	if len(ts) < 1 {
		errCtx := clCtx
		return nil, errCtx.NewError("Error: bad class definition")
	}

	extends, ts, err := p.buildClassExtendsExpression(ts)
	if err != nil {
		return nil, err
	}

	implements, ts, err := p.buildClassImplementsExpression(ts)
	if err != nil {
		return nil, err
	}

	universalName, ts, err := p.buildClassUniversalName(ts)
	if err != nil {
		return nil, err
	}

  class, err := js.NewClass(clType, extends, implements, isAbstract, isFinal, universalName, clCtx)
	if err != nil {
		return nil, err
	}

	if len(ts) != 1 {
		errCtx := raw.MergeContexts(ts...)
		return nil, errCtx.NewError("Error: unexpected tokens")
	}

	bracesGroup, err := raw.AssertBracesGroup(ts[0])
	if err != nil {
		return nil, err
	}

	// cant use buildBlockStatements, because classes have special syntax

	for _, field := range bracesGroup.Fields {
		remaining := field

		if len(field) == 0 {
			continue
		}

	Outer:
		for len(remaining) > 0 {
      encounteredAbstract := false
			for i, t := range remaining {
        if raw.IsWord(t, "abstract") {
          encounteredAbstract = true
        } else if raw.IsBracesGroup(t) {
          if encounteredAbstract {
            errCtx := t.Context()
            return nil, errCtx.NewError("Error: abstract member can't have function body")
          }

					switch i {
					case 0, 1:
						errCtx := t.Context()
						return nil, errCtx.NewError("Error: bad class member function definition")
					default:
						function, innerRemaining, err := p.buildFunction(remaining[0:i+1], true, false)
						if err != nil {
							return nil, err
						}

						if len(innerRemaining) != 0 {
							errCtx := raw.MergeContexts(innerRemaining...)
							return nil, errCtx.NewError("Error: unexpected tokens after member function")
						}

						if err := class.AddFunction(function); err != nil {
							return nil, err
						}

						remaining = remaining[i+1:]
						continue Outer
					}
				} else if encounteredAbstract && i == len(remaining) -1 && bracesGroup.IsSemiColon() {
          // append an empty braces group
          function, innerRemaining, err := p.buildFunction(append(remaining, raw.NewEmptyBracesGroup(t.Context())), true, false)
          if err != nil {
            return nil, err
          }

          if len(innerRemaining) != 0 {
            errCtx := raw.MergeContexts(innerRemaining...)
            return nil, errCtx.NewError("Error: unexpected tokens after abstract member function")
          }

          if err := class.AddFunction(function); err != nil {
            return nil, err
          }

          break Outer
        }
			} // loop over remaining

      // could be a property
      if raw.IsAnyWord(remaining[0]) {
        if p.isFunctionRoleKeyword(remaining[0]) {
          errCtx := remaining[0].Context()
          return nil, errCtx.NewError("Error: unexpected keyword")
        }

        propNameToken, err := raw.AssertWord(remaining[0])
        if err != nil {
          panic(err)
        }

        propName := js.NewWord(propNameToken.Value(), propNameToken.Context())
        var typeExpr *js.TypeExpression = nil
        if len(remaining) > 1 {
          typeExpr, err = p.buildTypeExpression(remaining[1:])
          if err != nil {
            return nil, err
          }
        }

        if err := class.AddProperty(propName, typeExpr); err != nil {
          return nil, err
        }

        // set remaining to zero length, so the loop quits
        remaining = []raw.Token{}
      } else {
        errCtx := raw.MergeContexts(remaining...)
        return nil, errCtx.NewError("Error: invalid class content")
      }
		}
	}

	return class, nil
}

func (p *JSParser) buildClassExpression(ts []raw.Token) (js.Expression, error) {
	return p.buildClass(ts)
}

func (p *JSParser) buildClassStatement(ts []raw.Token) (*js.Class, []raw.Token, error) {
	for i, t := range ts {
		if raw.IsBracesGroup(t) {
			statement, err := p.buildClass(ts[0 : i+1])
			if err != nil {
				return nil, nil, err
			}

			remaining := p.stripSeparators(i+1, ts, patterns.SEMICOLON)

			return statement, remaining, nil
		}
	}

	errCtx := raw.MergeContexts(ts...)
	return nil, nil, errCtx.NewError("Error: no class body found")
}
