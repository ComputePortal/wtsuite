package prototypes

import (
  "../values"

  "../../context"
)

func NewLocaleOptions(ctx context.Context) values.Value {
  b := NewBoolean(ctx)
  i := NewInt(ctx)
  s := NewString(ctx)

  return NewObject(map[string]values.Value{
    "compactDisplay": s,
    "currency": s,
    "currencyDisplay": s,
    "currencySign": s,
    "localeMatcher": s,
    "notation": s,
    "numberingSystem": s,
    "signDisplay": s,
    "style": s,
    "unit": s,
    "unitDisplay": s,
    "useGrouping": b,
    "minimumIntegerDigits": i,
    "minimumFractionDigits": i,
    "maximumFractionDigits": i,
    "minimumSignificantDigits": i,
    "maximumSignificantDigits": i,
  }, ctx)
}
