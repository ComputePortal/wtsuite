package values

import (
  "fmt"

  "../../context"
)

func checkOverload(overload []Value, ts []Value, ctx context.Context) error {
  if len(overload) == len(ts) {
    for j, arg := range overload {
      if j == len(overload) -1 && arg == nil && ts[j] != nil {
        errCtx := ts[j].Context()
        return errCtx.NewError("Error: expected void return value")
      } else if j == len(overload) - 1 && arg != nil && ts[j] == nil {
        errCtx := ctx
        return errCtx.NewError("Error: unexpected void return value")
      } else if err := arg.Check(ts[j], ctx); err != nil {
        return err
      }
    }

    return nil
  } else {
    return ctx.NewError("Error: incompatible function interface")
  }
}

func checkAnyOverload(overloads [][]Value, ts []Value, ctx context.Context) (int, error) {
  for i, overload := range overloads {
    if len(overload) == len(ts) {
      ok := true
      for j, arg := range overload {
        var detectedError error = nil
        if j == len(overload) - 1 && arg == nil && ts[j] != nil {
          detectedError = ctx.NewError("Error: expected void return value")
        } else if j == len(overload) - 1 && arg != nil && ts[j] == nil {
          detectedError = ctx.NewError("Error: unexpected void return value")
        } else if err := arg.Check(ts[j], ctx); err != nil {
          detectedError = err
        }

        if detectedError != nil {
          if len(overloads) == 1 {
            return 0, detectedError
          }

          ok = false
          break
        }
      }

      if ok {
        return i, nil
      }
    }
  }

  return 0, ctx.NewError("Error: incompatible function interface")
}

func checkAllOverloads(overloads [][]Value, tss [][]Value, ctx context.Context) error {
  if len(overloads) > len(tss) {
    return ctx.NewError("Error: missing function overloads")
  }

  for i, overload := range overloads {
    overloadFound := false

    for _, ts := range tss {
      if len(ts) == len(overload) {
        ok := true

        for j, arg := range overload {
          // return value can be nil (i.e. void)
          if j == len(ts) - 1 && arg == nil && ts[j] != nil {
            if len(tss) == 1 {
              return ctx.NewError("Error: expected void return value, got " + ts[j].TypeName())
            }

            ok = false
            break
          } else if j == len(ts) - 1 && arg != nil && ts[j] == nil {
            if len(tss) == 1 {
              return ctx.NewError("Error: expected non-void return value")
            }

            ok = false
            break
          } else if err := arg.Check(ts[j], ctx); err != nil {
            if len(tss) == 1 {
              return err
            }

            ok = false
            break
          }
        }

        if ok {
          overloadFound = true
          break
        }
      } else if len(tss) == 1 {
        return ctx.NewError("Error: function has different number of arguments")
      }
    }

    if !overloadFound {
      return ctx.NewError(fmt.Sprintf("Error: function overload %d not found", i+1))
    }
  }

  return nil
}

