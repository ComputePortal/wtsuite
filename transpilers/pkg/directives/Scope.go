package directives

import (
	"../functions"
	"../parsers"
	"../tokens/context"
	tokens "../tokens/html"
	"../tree"
	"../tree/svg"
)

type Scope interface {
	Parent() Scope
	GetNode() Node

	HasVar(key string) bool
	GetVar(key string) functions.Var
	SetVar(key string, v functions.Var)

	HasClass(key string) bool
	GetClass(key string) Class
	SetClass(key string, d Class)

	listValidVarNames() string

	// implements tokens.Scope
	Eval(key string, args []tokens.Token, ctx context.Context) (tokens.Token, error)
}

func IsTopLevel(scope Scope) bool {
	return scope.Parent() == nil
}

func buildAttributes(scope Scope, enumNode *ClassNode, tag *tokens.Tag,
	pos2opt []string) (*tokens.StringDict, error) {
	attr, err := tag.Attributes(pos2opt)
	if err != nil {
		return nil, err
	}

	if enumNode == nil {
		attr, err = attr.EvalStringDict(scope)
		if err != nil {
			return nil, err
		}
	} else {
		scanFn := func(key *tokens.String, lst_ tokens.Token) error {
			if tokens.IsAttrEnumList(lst_) {
				lst := lst_.(*tokens.List)
				enumNode.addAttrEnum(key.Value(), lst)
			}
			return nil
		}

		attr, err = attr.EvalStringDictScan(scope, scanFn)
		if err != nil {
			return nil, err
		}
	}

	if style_, ok := attr.Get("style"); ok && tokens.IsString(style_) {
		styleStr, err := tokens.AssertString(style_)
		if err != nil {
			panic(err)
		}

		style, err := parsers.ParseInlineDict(styleStr.Value(), styleStr.InnerContext())
		if err != nil {
			return nil, err
		}

		attr.Set("style", style)
	}

	return attr, nil
}

// NodeType can change from parentNode to this node
func buildTree(parent Scope, parentNode Node, nt NodeType,
	tagToken *tokens.Tag, collectDefaultOps bool) error {

	enumNode, _ := NewClassNode(parentNode, []Operation{})
	scope := NewSubScope(parent, enumNode) // the enumNode absorbs intermediate enum declarations

	attr, err := buildAttributes(scope, enumNode, tagToken, []string{})
	if err != nil {
		return err
	}

	var tag tree.Tag
	switch parentNode.Type() {
	case SVG:
		if !svg.IsTag(tagToken.Name()) {
			errCtx := tagToken.Context()
			return errCtx.NewError("Error: '" + tagToken.Name() + "' is not a valid svg tag")
		}

		tag, err = svg.BuildTag(tagToken.Name(), attr, tagToken.Context())
	case HTML:
		if !tree.IsTag(tagToken.Name()) {
			errCtx := tagToken.Context()
			return errCtx.NewError("Error: '" + tagToken.Name() + "' is not a valid html tag")
		}

		tag, err = tree.BuildTag(tagToken.Name(), attr, tagToken.Context())
	default:
		panic("unrecognized node type")
	}
	if err != nil {
		return err
	}

	id := tag.GetID()
	hasId := id != "" || collectDefaultOps
	var op Operation
	hasOp := false
	if hasId {
		op, hasOp, err = parentNode.PopOp(id)
		if err != nil {
			return err
		}
	}

	var newNode Node
	switch nt {
	case SVG:
		newNode = NewSVGNode(tag, parentNode)
	case HTML:
		newNode = NewNode(tag, parentNode)
	default:
		panic("unrecognized node type")
	}

	if hasOp {
		if err := op.Apply(scope, parentNode, newNode, tag, tagToken.Children()); err != nil {
			return err
		}
	} else {
		if err := parentNode.AppendChild(tag); err != nil {
			return err
		}

		for _, child := range tagToken.Children() {
			if err := BuildTag(scope, newNode, child); err != nil {
				return err
			}
		}
	}

	if collectDefaultOps && tag.GetID() != "" {
		// can only be append op
		op, hasOp, err := parentNode.PopOp("")
		if err != nil {
			return err
		}

		if hasOp {
			if err := op.Apply(scope, parentNode, newNode, nil, []*tokens.Tag{}); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildText(node Node, tag *tokens.Tag) error {
	return node.AppendChild(tree.NewText(tag.Text(), tag.Context()))
}

func BuildTag(scope Scope, node Node, tag *tokens.Tag) error {
	key := tag.Name()

	switch {
	case tag.IsText():
		return buildText(node, tag)
	case scope.HasClass(key):
		return BuildClass(scope, node, tag)
	case IsDirective(key):
		return BuildDirective(scope, node, tag)
	case node.Type() == SVG && key == "path":
		return buildSVGPath(scope, node, tag)
	case node.Type() == SVG && key == "arrow":
		return buildSVGArrow(scope, node, tag)
	case key == "path" || key == "arrow":
		panic("node type is bad")
	default:
		if err := buildTree(scope, node, node.Type(), tag, false); err != nil {
			// some error hints
			if key == "else" || key == "elseif" {
				context.AppendString(err, "Hint: did you forget to wrap in ifelse tag?")
			} else if key == "case" || key == "default" {
				context.AppendString(err, "Hint: did you forget to wrap in switch tag?")
			} else if key == "replace" || key == "append" || key == "preprend" {
				context.AppendString(err, "Hint: are you trying to instantiate a class'ed tag?")
			} else if node.Type() != SVG && svg.IsTag(key) {
				context.AppendString(err, "Hint: are you trying to use an svg tag?")
			}

			return err
		} else {
			return nil
		}
	}
}

// TODO: need access to node here
func eval(scope Scope, key string, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	switch {
	case scope.HasVar(key):
		fn, err := functions.AssertFun(scope.GetVar(key).Value)
		if err != nil {
			context.AppendContextString(err, "Info: called here", ctx)
			return nil, err
		}
		res, err := fn.EvalFun(scope, args, ctx)
		if err != nil {
			return nil, err
		}
		return res, nil
	case key == "svg-uri":
		return evalSVGURI(scope, args, ctx)
	case key == "file-url":
		return evalFileURL(scope, args, ctx)
	case key == "math-uri":
		return evalMathURI(scope, args, ctx)
	case key == "new":
		return evalNew(scope, args, ctx)
	case key == "block-id":
		return evalBlockID(scope, args, ctx)
	case key == "search-style":
		return evalSearchStyle(scope, args, ctx)
	/*case key == "attr": // 20200509: impossible to check with this()
		return evalAttrEnum(scope, args, ctx)
	case key == "attrs":
		return evalAttrEnumS(scope, args, ctx)
	case key == "attre":
		return evalAttrEnumS(scope, args, ctx)
	case key == "attrc":
		return evalAttrEnumC(scope, args, ctx)
	case key == "attra":
		return evalAttrEnumA(scope, args, ctx)*/
	case key == "idx":
		return evalElementCount(scope, args, ctx)
	case key == "var":
		return evalVar(scope, args, ctx)
	case key == "get":
		if len(args) > 0 && tokens.IsString(args[0]) {
			return evalGet(scope, args, ctx)
		}
		fallthrough
	default:
		return functions.Eval(scope, key, args, ctx)
	}
}
