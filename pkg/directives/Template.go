package directives

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type Template struct {
	name        string
	extends     string
	scope       Scope
	args        *tokens.List // works a little different from Parens
	argDefaults *tokens.List // works a little different from Parens
	superAttr   *tokens.RawDict // passed on to super
	children    []*tokens.Tag
	imported    bool
	exported    bool
  final       bool
	ctx         context.Context
}

func newTemplate(name string, extends string, scope Scope, args *tokens.List, argDefaults *tokens.List,
	superAttr *tokens.RawDict,
	children []*tokens.Tag,
	exported bool, final bool, ctx context.Context) Template {

	// copy the scope, in order to take a snapshot of its state
	subScope := NewSubScope(scope)

	return Template{
		name,
		extends,
		subScope,
		args,
		argDefaults,
		superAttr,
		children,
		false,
		exported,
    final,
		ctx,
	}
}

func assertValidTag(nameToken *tokens.String) error {
	errCtx := nameToken.InnerContext()
	name := nameToken.Value()
	if patterns.NAMESPACE_SEPARATOR_REGEXP.MatchString(name) {
		return errCtx.NewError("Error: invalid tag name, can't contain namespace separator '" +
			patterns.NAMESPACE_SEPARATOR + "'")
	} else if name == "template" || name == "for" || name == "if" || name == "ifelse" ||
		name == "import" || name == "print" || name == "script" || name == "style" ||
		name == "switch" || name == "var" || name == "else" || name == "elseif" ||
		name == "case" || name == "default" ||
		name == "replace" || name == "append" || name == "prepend" || name == "block" {
		return errCtx.NewError("Error: invalid tag name, is already a directive")
	} else if tree.IsTag(name) && NO_ALIASING {
		err := errCtx.NewError("Error: invalid tag name, is already a tag")
		return err
	} else {
		return nil
	}
}

// args for contexts
func assertArgDefaultsLast(argDefaults *tokens.List) error {
  prevDefault := -1
  if err := argDefaults.Loop(func(i int, value tokens.Token, last bool) error {
    if value == nil {
      if prevDefault >= 0 {
        prev, err := argDefaults.Get(prevDefault)
        if err != nil {
          panic(err)
        }
        errCtx := prev.Context()
        return errCtx.NewError("Error: defaults must come last")
      } 
    } else {
      prevDefault = i
    }

    return nil
  }); err != nil {
    return err
  }

  return nil
}

// doesnt change the node
func AddTemplate(scope Scope, node Node, tag *tokens.Tag) error {
	attr, err := tag.Attributes([]string{"args"})
	if err != nil {
		return err
	}

	nameToken, err := tokens.DictString(attr, "name")
	if err != nil {
		return err
	}

	if err := assertValidTag(nameToken); err != nil {
		return err
	}

	// extends is allowed to be evaluated
	subScope := NewSubScope(scope)

	extendsToken_, ok := attr.Get("extends")
	if !ok {
		errCtx := tag.Context()
		return errCtx.NewError("Error: extends not found")
	}

	// problem: surrounding scope can be modified?
	extendsToken_, err = extendsToken_.Eval(subScope) // TODO: variables could be set here but wont be available anywhere: this should throw an error
	if err != nil {
		return err
	}

	extendsToken, err := tokens.AssertString(extendsToken_)
	if err != nil {
		return err
	}

	var args *tokens.List = nil
	var argDefaults *tokens.List = nil
	if args_, ok := attr.Get("args"); ok {
		if tokens.IsList(args_) {
			// dont evaluate!, but make sure we have only strings
			args, err = tokens.ToStringList(args_)
			if err != nil {
				return err
			}

			argDefaults = tokens.NewNilList(args.Len(), attr.Context())
		} else if tokens.IsParens(args_) {
			argParens, err := tokens.AssertParens(args_)
			if err != nil {
				panic(err)
			}

			args = tokens.NewValuesList(argParens.Values(), argParens.Context())
			argDefaults = tokens.NewValuesList(argParens.Alts(), argParens.Context())

      if err := assertArgDefaultsLast(argDefaults); err != nil {
        return err
      }
		} else {
			errCtx := args_.Context()
			return errCtx.NewError("Error: expected list or parens")
		}
	} else {
		args = tokens.NewEmptyList(attr.Context())
		argDefaults = tokens.NewNilList(args.Len(), attr.Context())
	}

	exported, err := tokens.DictHasFlag(attr, "export")
	if err != nil {
		return err
	}

	superAttr, err := tokens.DictRawDict(attr, "super")
	if err != nil {
		return err
	}

	extends := extendsToken.Value()

  final := false
  if finalToken, ok := attr.Get(".final"); ok {
    if tokens.IsTrueBool(finalToken) {
      final = true
    }
  }

  // check that attr has no other args
  // (TODO: also for other directives)
  if err := attr.AssertOnlyValidKeys([]string{"export", "name", "args", "extends", "super", ".final"}); err != nil {
    return err
  }

	key := nameToken.Value()

	switch {
	case scope.HasTemplate(key):
		errCtx := nameToken.InnerContext()
		err := errCtx.NewError("Error: can't redefine tag")
		err.AppendContextString("Info: defined here", scope.GetTemplate(key).ctx)
		return err
	default:
    if err := scope.SetTemplate(key, newTemplate(key, extends, scope, args, argDefaults, superAttr, tag.Children(), exported, final, tag.Context())); err != nil {
      return err
    }
	}

	return nil
}

// first return value: ok
// second return value: can be passed to parent
func (c Template) hasArg(key string) bool {
	args := c.args.GetTokens()

	for _, arg_ := range args {
		arg, err := tokens.AssertString(arg_)
		if err != nil {
			panic("should've been caught before")
		}

		test := arg.Value()

    if test == key {
      return true
    }
	}

	return false
}

func (c Template) argsStringList() ([]string, error) {
	res := make([]string, 0)

	for _, v := range c.args.GetTokens() {
		arg, err := tokens.AssertString(v)
		if err != nil {
      return nil, err
		}

		res = append(res, arg.Value())
	}

	return res, nil
}

func (c Template) argsWithoutDefaults() []string {
	res := make([]string, 0)

	for i, v := range c.args.GetTokens() {
    argDefault, err := c.argDefaults.Get(i)
    if err != nil {
      panic(err)
    }

    if argDefault == nil {
      arg, err := tokens.AssertString(v)
      if err != nil {
        panic(err)
      }

      res = append(res, arg.Value())
    }
	}

	return res
}

func (c Template) listValidArgNames() string {
	var b strings.Builder

	for _, v := range c.args.GetTokens() {
		arg, err := tokens.AssertString(v)
		if err != nil {
			panic(err)
		}

		b.WriteString(arg.Value())
		b.WriteString("\n")
	}

	return b.String()
}

func (c Template) instantiate(node *TemplateNode, args *tokens.StringDict,
	ctx context.Context) error {
	subScope := NewSubScope(c.scope)
  setElementCount(subScope, node, c.ctx)
  setLazyTagVars(subScope, c.ctx) // __nchildren__ is only usable in the attributes
  setParentStyle(subScope, node, c.ctx)

	// loop incoming attr and check if it is in c.args
	if err := args.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
    kVal := k.Value()
    force := false
    if strings.HasSuffix(kVal, "!") {
      force = true
      kVal = kVal[0:len(kVal)-1]
    }

		if ok := c.hasArg(kVal); !ok && !force {
			errCtx := k.Context()
			err := errCtx.NewError("Error: invalid tag attribute")
			context.AppendString(err, "Info: available args for "+c.name+
				"\n"+c.listValidArgNames())
			return err
		} else if ok {
      // dont set if forced but not actually available
      vVar := functions.Var{v, true, true, false, false, v.Context()}
      if err := subScope.SetVar(kVal, vVar); err != nil {
        return err
      }
    }
		return nil
	}); err != nil {
		return err
	}

  // check that the argsWithout defaults are all available in the incoming arg lists
  for _, argName := range c.argsWithoutDefaults() {
    if _, ok := args.Get(argName); !ok {
      errCtx := args.Context()
      return errCtx.NewError("Error: arg " + argName + " not specified")
    }
  }

	// cut off the exclamation marks
	templateArgNames, err := c.argsStringList()
  if err != nil {
    return err
  }

	// now loop the defaults, and instantiate those that are not in incoming args (using the same subScope
	if err := c.argDefaults.Loop(func(i int, t tokens.Token, last bool) error {
		if t == nil {
			// continue
			return nil
		}

		argName := templateArgNames[i]

		if _, ok1 := args.Get(argName); !ok1 {
      if _, ok2 := args.Get(argName + "!"); !ok2 {
        v, err := t.Eval(subScope)
        if err != nil {
          return err
        }

        vVar := functions.Var{v, true, true, false, false, v.Context()}
        if err := subScope.SetVar(argName, vVar); err != nil {
          return err
        }
      }
		}

		return nil
	}); err != nil {
		return err
	}

	templateSuperAttr, err := c.superAttr.EvalRawDict(subScope)
	if err != nil {
		return err
	}

	templateCtx := ctx

	if subScope.HasTemplate(c.extends) {
    // check that extends is not final
    if subScope.GetTemplate(c.extends).final {
      errCtx := c.ctx
      return errCtx.NewError("Error: can't extend " + c.extends + " (is final class)")
    }

    subTag := tokens.NewTag(c.extends, templateSuperAttr, c.children, templateCtx)
		if err := BuildTemplate(subScope, node, subTag); err != nil {
			return err
		}
	} else {
		nType := node.Type()
		if c.extends == "svg" {
			nType = SVG
		}

    tNode, err := prepareOperations(subScope, node, c.children, c.ctx)
    if err != nil {
      return err
    }

    subTag := tokens.NewTag(c.extends, templateSuperAttr, []*tokens.Tag{}, templateCtx)
		if err := buildTree(NewSubScope(subScope), tNode, nType, subTag, "default"); err != nil {
			return err
		}

    if err := tNode.AssertAllOperationsDone(templateCtx); err != nil {
      return err
    }
	}

  // insert the forced attributes
  child := node.getLastChild()
  if child == nil {
    panic("shouldn't be nil")
  }

  childAttr := child.Attributes()

	if err := args.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
    kVal := k.Value()
    if strings.HasSuffix(kVal, "!") {
      kVal = kVal[0:len(kVal)-1]

      childAttr.Set(tokens.NewValueString(kVal, k.Context()), v)
    }

    return nil
  }); err != nil {
    panic(err)
  }

	return nil
}

// pop all operations referenced by 'block' directives, and give them a new unique target name
// no scope needed because nothing is evaluated
func prepareBlocks(node *TemplateNode, tags []*tokens.Tag, insideBranch bool, newOpNames map[string]string) error {
  for _, tag := range tags {
    switch tag.Name() {
    case "ifelse", "if", "elseif", "else", "for", "switch", "case", "default":
      if err := prepareBlocks(node, tag.Children(), true, newOpNames); err != nil {
        return err
      }
    case "block":
      if node.GetBlockTarget(tag) == "" {
        name, err := getOpNameTarget("name", tag)
        if err != nil {
          return err
        }

        // blocks can exist multiple times

        op, err := node.PopOp(name)
        if err != nil {
          return err
        }

        if op != nil {
          newName := NewUniqueOpTargetName()
          newOpNames[op.Target()] = newName
          op.SetTarget(newName)
          if err := node.PushOp(op); err != nil {
            return err
          }

          node.SetBlockTarget(tag, newName)
        } else if newName, ok := newOpNames[name]; ok {
          node.SetBlockTarget(tag, newName)
        } else {
          node.SetBlockTarget(tag, name)
        }
      }

      if err := prepareBlocks(node, tag.Children(), insideBranch, newOpNames); err != nil {
        return err
      }
    case "var", "template":
      // dont do anything
    default:
      if err := prepareBlocks(node, tag.Children(), insideBranch, newOpNames); err != nil {
        return err
      }
    }
  }

  return nil
}

func prepareOperations(scope Scope, node Node, tags []*tokens.Tag, ctx context.Context) (*TemplateNode, error) {

  tNode := NewTemplateNode(node, ctx)

	subScope := NewSubScope(scope)

  if err := prepareBlocks(tNode, tags, false, make(map[string]string)); err != nil {
    return nil, err
  }

  tNode.StartDeferral() // needed for tags nested inside directives

	for _, tag := range tags {
    if err := BuildTag(subScope, tNode, tag); err != nil {
      return nil, err
    }
  }

  tNode.StopDeferral()

  return tNode, nil
}

func BuildTemplate(scope Scope, node Node, tag *tokens.Tag) error {
	templateName := tag.Name()
	template := scope.GetTemplate(templateName)

	// evaluate the attributes
	attrScope := NewSubScope(scope)
  argsStringList, err := template.argsStringList()
  if err != nil {
    return err
  }

	attr, err := buildAttributes(attrScope, tag, argsStringList)
	if err != nil {
		return err
	}

	tNode, err := prepareOperations(attrScope, node, tag.Children(), template.ctx)
	if err != nil {
		return err
	}

	if err := template.instantiate(tNode, attr, tag.Context()); err != nil {
		return err
	}

  if err := tNode.AssertAllOperationsDone(tag.Context()); err != nil {
    return err
  }

	return nil
}

var _addTemplateOk = registerDirective("template", AddTemplate)
