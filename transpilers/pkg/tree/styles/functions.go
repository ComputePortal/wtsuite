package styles

import (
	"strings"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type AtFunction func(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error)

var _atFunctions = make(map[string]AtFunction)

// to avoid circular initialization
func registerAtFunction(key string, fn AtFunction) bool {
	_atFunctions[key] = fn

	return true
}

func BuildNested(sel Selector, k *tokens.String, v tokens.Token) ([]Rule, error) {
	args := strings.Fields(strings.TrimLeft(k.Value(), "@"))
	if len(args) == 0 {
		errCtx := k.Context()
		return nil, errCtx.NewError("Error: expected something after the @/.")
	}

	fnKey := args[0]
	args = args[1:]

	if fn, ok := _atFunctions[fnKey]; ok {
		return fn(sel, args, v, k.Context())
	} else {
		errCtx := k.Context()
		return nil, errCtx.NewError("Error: '" + fnKey + "' at-function not recognized")
	}
}

func expandNested(attr *tokens.StringDict, sel Selector) (map[string]string, []Rule, error) {
	leafAttributes := make(map[string]string)
	nestedRules := make([]Rule, 0)

	if err := attr.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
		if strings.HasPrefix(k.Value(), "@") {
			fnRules, err := BuildNested(sel, k, v)
			if err != nil {
				return err
			}

			nestedRules = append(nestedRules, fnRules...)
		} else if strings.HasPrefix(k.Value(), ".") {
			args := strings.Fields(k.Value())
			fnRules, err := NamedClass(sel, args, v, k.Context())
			if err != nil {
				return err
			}

			nestedRules = append(nestedRules, fnRules...)
		} else {
			if err := dictEntryToStringMapEntry(k, v, leafAttributes); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return leafAttributes, nil, err
	}

	return leafAttributes, nestedRules, nil
}
