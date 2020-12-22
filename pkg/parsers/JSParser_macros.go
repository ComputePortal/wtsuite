package parsers

import (
	"../files"
	"../tokens/context"
	"../tokens/js"
)

func (p *JSParser) buildImportDefaultMacro(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if path, ok := args[0].(*js.LiteralString); ok {
    // add path as invisible pathplaceholder, so refactoring code can change it

    if _, err := files.Search(ctx.Caller(), path.Value()); err != nil {
			errCtx := path.Context()
			return nil, errCtx.NewError("Error: file not found")
		}

		// give the import a dummy name so it can never actually be referred to
		name := "." + path.Value() + " default"
		if err := p.module.AddImportedName(name, "default", path, ctx); err != nil {
			return nil, err
		}

		return js.NewVarExpression(name, ctx), nil
	} else {
		return nil, ctx.NewError("Error: expected literal string")
	}
}
