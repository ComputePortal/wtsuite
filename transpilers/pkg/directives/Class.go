package directives

import (
	"strings"

	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
	"../tokens/patterns"
	"../tree"
)

type Class struct {
	name        string
	extends     string
	scope       Scope
	args        *tokens.List
	argDefaults *tokens.List
	superAttr   *tokens.RawDict // passed on to super
	thisAttr    *tokens.RawDict // attributes not passed to super
	blocks      *tokens.StringDict
	children    []*tokens.Tag
	imported    bool
	exported    bool
	ctx         context.Context
}

func newClass(name string, extends string, scope Scope, args *tokens.List, argDefaults *tokens.List,
	superAttr *tokens.RawDict,
	thisAttr *tokens.RawDict,
	blocks *tokens.StringDict,
	children []*tokens.Tag,
	exported bool, ctx context.Context) Class {

	// copy the scope, in order to take a snapshot of its state
	subScope := NewSubScope(scope, scope.GetNode())

	if blocks == nil {
		panic("cant be nil")
	}

	return Class{
		name,
		extends,
		subScope,
		args,
		argDefaults,
		superAttr,
		thisAttr,
		blocks,
		children,
		false,
		exported,
		ctx,
	}
}

func assertValidTag(nameToken *tokens.String) error {
	errCtx := nameToken.InnerContext()
	name := nameToken.Value()
	if patterns.NAMESPACE_SEPARATOR_REGEXP.MatchString(name) {
		return errCtx.NewError("Error: invalid tag name, can't contain namespace separator '" +
			patterns.NAMESPACE_SEPARATOR + "'")
	} else if name == "class" || name == "for" || name == "if" || name == "ifelse" ||
		name == "import" || name == "print" || name == "script" || name == "style" ||
		name == "switch" || name == "var" || name == "else" || name == "elseif" ||
		name == "case" || name == "default" ||
		name == "replace" || name == "append" || name == "prepend" {
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
func AddClass(scope Scope, node Node, tag *tokens.Tag) error {
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
	subScope := NewSubScope(scope, node)

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


	var blocks *tokens.StringDict = nil
	if blocks_, ok := attr.Get("blocks"); ok {
		if tokens.IsStringDict(blocks_) {
			blocks, err = tokens.AssertStringDict(blocks_)
			if err != nil {
				panic(err)
			}
		} else if tokens.IsList(blocks_) {
			blocksLst, err := tokens.AssertList(blocks_)
			if err != nil {
				return err
			}

			blocks = tokens.NewEmptyStringDict(blocksLst.Context())

			if err := blocksLst.Loop(func(i int, val_ tokens.Token, last bool) error {
				val, err := tokens.AssertString(val_)
				if err != nil {
					return err
				}

				rhs := tokens.NewValueFunction("uid", []tokens.Token{}, val.Context())
				blocks.Set(val, rhs)

				return nil
			}); err != nil {
				return err
			}
		} else if tokens.IsRawDict(blocks_) {
			rd, err := tokens.AssertRawDict(blocks_)
			if err != nil {
				panic(err)
			}

			blocks, err = rd.ToStringDict()
			if err != nil {
				return err
			}
		} else {
			errCtx := blocks_.Context()
			return errCtx.NewError("Error: expected list or dict")
		}
	} else {
		blocks = tokens.NewEmptyStringDict(tag.Context())
	}

	exported, err := tokens.DictHasFlag(attr, "export")
	if err != nil {
		return err
	}

	superAttr, err := tokens.DictRawDict(attr, "super")
	if err != nil {
		return err
	}

	var thisAttr *tokens.RawDict = nil
	if this_, ok := attr.Get("this"); ok {
		thisAttr, err = tokens.AssertRawDict(this_)
		if err != nil {
			return err
		}
	} else {
		thisAttr = tokens.NewEmptyRawDict(tag.Context())
	}

	extends := extendsToken.Value()

	key := nameToken.Value()

	switch {
	case scope.HasClass(key):
		errCtx := nameToken.InnerContext()
		err := errCtx.NewError("Error: can't redefine tag")
		err.AppendContextString("Info: defined here", scope.GetClass(key).ctx)
		return err
	default:
		scope.SetClass(key, newClass(key, extends, scope, args, argDefaults, superAttr, thisAttr, blocks, tag.Children(), exported, tag.Context()))
	}

	return nil
}

// first return value: ok
// second return value: can be passed to parent
func (c Class) hasArg(key string) (bool, bool) {
	// should list be customizable via command line args?
	if key == "id" || (tree.AUTO_HREF && key == "href") || key == "style" || key == "inactive" { // id/href/style are always passed on for convenience
		return true, true
	}

	args := c.args.GetTokens()

	for _, arg_ := range args {
		arg, err := tokens.AssertString(arg_)
		if err != nil {
			panic("should've been caught before")
		}

		test := arg.Value()

		if strings.HasSuffix(test, "!") {
			if test[0:len(test)-1] == key {
				return true, true
			}
		} else {
			if test == key {
				return true, false
			}
		}
	}

	return false, false
}

func (c Class) argsStringList() []string {
	res := make([]string, 0)

	for _, v := range c.args.GetTokens() {
		arg, err := tokens.AssertString(v)
		if err != nil {
			panic(err)
		}

		res = append(res, strings.TrimRight(arg.Value(), "!"))
	}

	return res
}

func (c Class) listValidArgNames() string {
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

func (c Class) instantiate(node *ClassNode, args *tokens.StringDict,
	ctx context.Context) error {
	subScope := NewSubScope(c.scope, node)

	// loop incoming attr and check if it is in c.args
	if err := args.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
    kVal := k.Value()
    force := false
    if strings.HasSuffix(kVal, "!") {
      force = true
      kVal = kVal[0:len(kVal)-1]
    }

		if ok, _ := c.hasArg(kVal); !ok && !force {
			errCtx := k.Context()
			err := errCtx.NewError("Error: invalid tag attribute")
			context.AppendString(err, "Info: available args for "+c.name+
				"\n"+c.listValidArgNames())
			return err
		} else if ok {
      // dont set if forced but not actually available
      vVar := functions.Var{v, true, true, false, false, v.Context()}
      subScope.SetVar(kVal, vVar)
    }
		return nil
	}); err != nil {
		return err
	}

	// cut off the exclamation marks
	classArgNames := c.argsStringList()

	// now loop the defaults, and instantiate those that are not in incoming args (using the same subScope
	if err := c.argDefaults.Loop(func(i int, t tokens.Token, last bool) error {
		if t == nil {
			// continue
			return nil
		}

		argName := classArgNames[i]

		if _, ok1 := args.Get(argName); !ok1 {
      if _, ok2 := args.Get(argName + "!"); !ok2 {
        v, err := t.Eval(subScope)
        if err != nil {
          return err
        }

        vVar := functions.Var{v, true, true, false, false, v.Context()}
        subScope.SetVar(argName, vVar)
      }
		}

		return nil
	}); err != nil {
		return err
	}

	// superAttr, are built in a special way
	//  it is scanned a first time to find literal enums (which are added to the ClassNode,
	//  then the other pairs are evaluated
	// superAttr
	scanFn := func(key_ tokens.Token, lst_ tokens.Token) error {
		if tokens.IsString(key_) && tokens.IsAttrEnumList(lst_) {
			key := key_.(*tokens.String)
			lst := lst_.(*tokens.List)
			node.addAttrEnum(key.Value(), lst)
		}
		return nil
	}

	// dont bother scanning, enum consistency is detected during
	classSuperAttr, err := c.superAttr.EvalRawDictScan(subScope, scanFn)
	if err != nil {
		return err
	}

	if c.blocks == nil {
		panic("cant be nil")
	}
	classBlocks, err := c.blocks.EvalStringDict(subScope)
	if err != nil {
		return err
	}

	if err := node.MapBlocks(classBlocks); err != nil {
		return err
	}

	// TODO: this must be done differently, so we can catch the operations
	// set the __blocks__ variable in the subScope
	blocksVar := functions.Var{classBlocks, true, true, false, false, classBlocks.Context()}
	subScope.SetVar("__blocks__", blocksVar)

	// pass attr to dattr if they are not already there and have a ! suffix
  // TODO: get rid of this
	if err := args.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
    if !strings.HasSuffix(k.Value(), "!") {
      ok, canBePassedOn := c.hasArg(k.Value())
      if !ok {
        panic("should've been caught in previous d.attr.Loop")
      }

      _, hasAlready := classSuperAttr.Get(k.Value())
      if canBePassedOn && !hasAlready {
        classSuperAttr.Set(k, v)
      }
    }

		return nil
	}); err != nil {
		return err
	}

	classCtx := c.ctx.ChangeCaller(ctx.Caller())

	if subScope.HasClass(c.extends) {
		subTag := tokens.NewTag(c.extends, classSuperAttr, []*tokens.Tag{}, classCtx)
		cNode, err := prepareOperations(subScope, node, c.children)
		if err != nil {
			return err
		}
		if err := BuildTag(subScope, cNode, subTag); err != nil {
			return err
		}

		remainingOps := cNode.GetOperations()
		if len(remainingOps) != 0 {
			errCtx := ctx
			err := errCtx.NewError("Error: unapplied ops")
			for _, op := range remainingOps {
				err.AppendContextString("Info: not applied to "+op.ID(), op.Context())
			}
			return err
		}
	} else {
		nType := node.Type()
		if c.extends == "svg" {
			nType = SVG
		}

		subTag := tokens.NewTag(c.extends, classSuperAttr, c.children, classCtx)
		if err := buildTree(subScope, node, nType, subTag, true); err != nil {
			return err
		}
	}

	// add this attr to the new child of node
  child := node.getLastChild()
  if child == nil {
    panic("shouldn't be nil")
  }

  childAttr := child.Attributes()
	if !c.thisAttr.IsEmpty() {
		thisAttr, err := c.thisAttr.EvalStringDict(subScope)
		if err != nil {
			return err
		}

		if err := thisAttr.Loop(func(key *tokens.String, val tokens.Token, last bool) error {
			// XXX: should we display warnings if we are overwritting attributes?
			childAttr.Set(key, val)
			return nil
		}); err != nil {
			panic(err)
		}
	}

  // also insert the forced attributes
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

func prepareOperations(scope Scope, node Node, tags []*tokens.Tag) (*ClassNode, error) {
	subScope := NewSubScope(scope, node)

	operations := make([]Operation, 0)
	toDefault := make([]*tokens.Tag, 0)

	for _, tag := range tags {
		key := tag.Name()
		switch {
		case key == "var" || key == "class" || key == "import":
			if err := BuildDirective(subScope, node, tag); err != nil {
				return nil, err
			}
		case strings.HasPrefix(key, "#"):
			op, err := NewReplaceChildrenOp(subScope, tag)
			if err != nil {
				return nil, err
			}
			operations = append(operations, op)
		case key == "append":
			op, err := NewAppendOp(subScope, tag)
			if err != nil {
				return nil, err
			}
			operations = append(operations, op)
		case key == "prepend":
			errCtx := tag.Context()
			return nil, errCtx.NewError("Error: prepend no longer supported")
			/*op, err := NewPrependOp(subScope, tag)
			if err != nil {
				return nil, err
			}
			operations = append(operations, op)*/
		case key == "replace":
			op, err := NewReplaceOp(subScope, tag)
			if err != nil {
				return nil, err
			}
			operations = append(operations, op)
		default:
			toDefault = append(toDefault, tag)
		}
	}

	if len(toDefault) != 0 {
		op, err := NewAppendToDefaultOp(subScope, toDefault)
		if err != nil {
			return nil, err
		}
		operations = append(operations, op)
	}

	return NewClassNode(node, operations)
}

func BuildClass(scope Scope, node Node, tag *tokens.Tag) error {
	className := tag.Name()
	class := scope.GetClass(className)

	// evaluate the attributes
	attrScope := NewSubScope(scope, node)
	attr, err := buildAttributes(attrScope, nil, tag, class.argsStringList())
	if err != nil {
		return err
	}

	cNode, err := prepareOperations(attrScope, node, tag.Children())
	if err != nil {
		return err
	}

	if err := class.instantiate(cNode, attr, tag.Context()); err != nil {
		return err
	}

	remainingOps := cNode.GetOperations()
	if len(remainingOps) != 0 {
		errCtx := tag.Context()
		err := errCtx.NewError("Error: unapplied ops")
		for _, op := range remainingOps {
			err.AppendContextString("Info: not applied to "+op.ID(), op.Context())
		}
		return err
	}

	return nil
}

var _addClassOk = registerDirective("class", AddClass)
