package directives

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/svg"
)

type Scope interface {
	Parent() Scope
  Permissive() bool

	HasVar(key string) bool
	GetVar(key string) functions.Var
	SetVar(key string, v functions.Var) error

	HasTemplate(key string) bool
	GetTemplate(key string) Template
	SetTemplate(key string, d Template) error

	listValidVarNames() string
  //setBlockTarget(block *tokens.Tag, target string)
  //getBlockTarget(block *tokens.Tag)

	// implements tokens.Scope
	Eval(key string, args *tokens.Parens, ctx context.Context) (tokens.Token, error)

  GetCache() *FileCache
}

func NewSubScope(parent Scope) Scope {
	if parent == nil {
		panic("parent can't be nil")
	}

  scope_ := newScopeData(parent)
  scope := &scope_

	return scope
}

func IsTopLevel(scope Scope) bool {
	return scope.Parent() == nil
}

func buildAttributes(scope Scope, tag *tokens.Tag,
	pos2opt []string) (*tokens.StringDict, error) {
	attr, err := tag.Attributes(pos2opt)
	if err != nil {
		return nil, err
	}

  attr, err = attr.EvalStringDict(scope)
  if err != nil {
    return nil, err
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

func removeForcedSuffix(attr *tokens.StringDict) (*tokens.StringDict, error) {
  return attr.MapStringKeys(func (k string) string {
    if strings.HasSuffix(k, "!") {
      return k[0:len(k)-1]
    } else {
      return k
    }
  })
}

// NodeType can change from parentNode to this node
// collectDefaultOps==true in case Template extends this tag
func buildTree(parent Scope, parentNode Node, nt NodeType,
	tagToken *tokens.Tag, opName string) error {

	scope := NewSubScope(parent) // the enumNode absorbs intermediate enum declarations

	attr, err := buildAttributes(scope, tagToken, []string{})
	if err != nil {
		return err
	}

  attr, err = removeForcedSuffix(attr)
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

	var newNode Node
	switch nt {
	case SVG:
		newNode = NewSVGNode(tag, parentNode)
	case HTML:
		newNode = NewNode(tag, parentNode)
	default:
		panic("unrecognized node type")
	}

	var op Operation
	if opName != "" {
		op, err = parentNode.PopOp(opName)
		if err != nil {
			return err
		}
	}

  if err := parentNode.AppendChild(tag); err != nil {
    return err
  }

	if op != nil {
		if err := op.Apply(scope, newNode, tagToken.Children()); err != nil {
			return err
		}
	} else {
    // there is no way to estimate the number of children before hand?
		for _, child := range tagToken.Children() {
			if err := BuildTag(scope, newNode, child); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildText(node Node, tag *tokens.Tag) error {
	return node.AppendChild(tree.NewText(tag.Text(), tag.Context()))
}

// TODO: append, replace, prepend scope based
func buildDeferred(scope Scope, node *TemplateNode, tag *tokens.Tag) error {
  key := tag.Name()
  switch {
  case tag.IsText() || scope.HasTemplate(key) || key == "block" || key == "print":
    return node.AppendToDefault(scope, tag)
  case key == "append":
    // append directive is not directly registered
    return AppendToBlock(scope, node, tag)
  case key == "replace":
    // replace directive is not directly registered
    return ReplaceBlockChildren(scope, node, tag)
  //case IsDirective(key) && tag.IsDirective():
		//return BuildDirective(scope, node, tag)
  default:
    return node.AppendToDefault(scope, tag)
  }
}

func BuildTag(scope Scope, node Node, tag *tokens.Tag) error {
	key := tag.Name()

	switch {
  case IsDeferringTemplateNode(node):
    return buildDeferred(scope, node.(*TemplateNode), tag)
	case tag.IsText():
		return buildText(node, tag)
	case scope.HasTemplate(key):
		return BuildTemplate(scope, node, tag)
	case IsDirective(key) && ((key != "template" && key != "var") || tag.IsDirective()):
		return BuildDirective(scope, node, tag)
	case node.Type() == SVG && key == "path":
		return buildSVGPath(scope, node, tag)
	case node.Type() == SVG && key == "arrow":
		return buildSVGArrow(scope, node, tag)
	case key == "path" || key == "arrow":
		panic("node type is bad")
	default:
		if err := buildTree(scope, node, node.Type(), tag, ""); err != nil {
			// some error hints
			if key == "else" || key == "elseif" {
				context.AppendString(err, "Hint: did you forget to wrap in ifelse tag?")
			} else if key == "case" || key == "default" {
				context.AppendString(err, "Hint: did you forget to wrap in switch tag?")
			} else if key == "replace" || key == "append" || key == "preprend" {
				context.AppendString(err, "Hint: are you trying to instantiate a templated tag?")
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
func eval(scope Scope, key string, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
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
	//case key == "search-style":
		//return evalSearchStyle(scope, args, ctx)
	case key == "var":
		return evalVar(scope, args, ctx)
  case key == "issymbol":
    return evalIsSymbol(scope, args, ctx)
  case key == "import":
    return evalDynamicImport(scope, args, ctx)
	case key == "get":
		if args.Len() > 0 && tokens.IsString(args.Values()[0]) {
			return evalGet(scope, args, ctx)
		}
		fallthrough
	default:
		return functions.Eval(scope, key, args, ctx)
	}
}
