package directives

import (
	"path/filepath"

	"../files"
	//"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

const URL = "__url__"
var IGNORE_UNSET_URLS = false

var _fileURLs map[string]string = nil
var _activeURL *tokens.String = nil

func RegisterURL(path string, url string) {
	if _fileURLs == nil {
		_fileURLs = make(map[string]string)
	}

	_fileURLs[path] = url
}

/*func SetURL(scope Scope, path string, ctx context.Context) {
	if url, ok := _fileURLs[path]; ok {
		urlToken := tokens.NewValueString(url, ctx)

		v := functions.Var{
			urlToken,
			true,
			true,
			false,
			false,
			ctx,
		}

		scope.SetVar(URL, v)

		//_activeURL = urlToken
	}
}*/

func GetActiveURL(ctx context.Context) (*tokens.String, error) {
	if _activeURL == nil {
		return nil, ctx.NewError("Error: __url__ not set here")
	}

	return _activeURL, nil
}

func SetActiveURL(url string) {
	_activeURL = tokens.NewValueString(url, context.NewDummyContext())
}

func UnsetActiveURL() {
	_activeURL = nil
}

/*func UnsetURL(scope *TagScope) {
	if _, ok := scope.vars[URL]; ok {
		delete(scope.vars, URL)
	}

}*/

func evalFileURL(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	arg0, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	pathToken, err := tokens.AssertString(arg0)
	if err != nil {
		return nil, err
	}

	path := pathToken.Value()
	if path == "" {
		return nil, ctx.NewError("Error: path can't be empty (hint: use \".\" to refer to this file)")
	}

	// XXX: alternatively we could use path in the ctx
	current_ := scope.GetVar(FILE)

	current, err := tokens.AssertString(current_.Value)
	if err != nil {
		panic(err)
	}

	if path == "." {
		path = current.Value()
	} else if !filepath.IsAbs(path) {
		currentDir := filepath.Dir(current.Value())

		path = filepath.Join(currentDir, path)
	}

	if err := files.AssertFile(path); err != nil {
		return nil, ctx.NewError("Error: file '" + path + "' not found")
	}

	if url, ok := _fileURLs[path]; ok {
		return tokens.NewValueString(url, ctx), nil
	} else {
    if !IGNORE_UNSET_URLS {
      return nil, ctx.NewError("Error: url for '" + path + "' not set")
    } else {
      // used when doing refactorings, where the url doesnt matter
      return tokens.NewValueString("", ctx), nil
    }
	}
}
