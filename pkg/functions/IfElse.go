package functions

import (
	"math"
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

// 3, or 5, or 7 etc. (last must always be else block
func IfElse(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if int(math.Mod(float64(len(args)-1), 2.0)) != 0 {
		return nil, ctx.NewError("Error: expected 3, 5, 7... arguments")
	}

	for i := 0; i < len(args); i += 2 {
		if i < len(args)-1 {
			argCond, err := args[i].Eval(scope)
			if err != nil {
				return nil, err
			}

      if tokens.IsLazy(argCond) {
        remArgs := args[i:]
        return tokens.NewLazy(func(tag tokens.FinalTag) (tokens.Token, error) {
          for i_ := 0; i_ < len(remArgs); i_ += 2 {
            if i_ < len(remArgs)-1 {
              remArg, err := remArgs[i_].Eval(scope)
              if err != nil {
                return nil, err
              }

              if tokens.IsLazy(remArg) {
                remArg, err = remArg.EvalLazy(tag)
                if err != nil {
                  return nil, err
                }
              } 

              cond, err := tokens.AssertBool(remArg)
              if err != nil {
                return nil, ctx.NewError("Error: expected bool for arg " + strconv.Itoa(i_))
              }

              if cond.Value() {
                res, err := remArgs[i_+1].Eval(scope)
                if err != nil {
                  return nil, err
                }

                if tokens.IsLazy(res) {
                  res, err = res.EvalLazy(tag)
                  if err != nil {
                    return nil, err
                  }
                }

                return res, nil
              }
            }
          }

          res, err := remArgs[len(remArgs)-1].Eval(scope)
          if err != nil {
            return nil, err
          }
          if tokens.IsLazy(res) {
            res, err = res.EvalLazy(tag)
            if err != nil {
              return nil, err
            }
          }
          return res, nil
        }, ctx), nil
      }

			cond, err := tokens.AssertBool(argCond)
			if err != nil {
				return nil, err
			}

			if cond.Value() {
				return args[i+1].Eval(scope)
			}
		}
	}

	return args[len(args)-1].Eval(scope)
}
