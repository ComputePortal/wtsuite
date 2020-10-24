package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

var (
	NO_ALIASING = false
	VERBOSITY   = 0
)

type Directive func(scope Scope, node Node, tag *tokens.Tag) error

var _directiveTable = make(map[string]Directive)

func registerDirective(key string, fn Directive) bool {
	_directiveTable[key] = fn

	return true
}

func IsDirective(key string) bool {
	_, ok := _directiveTable[key]
	return ok
}

func BuildDirective(scope Scope, node Node, tag *tokens.Tag) error {
	fn, ok := _directiveTable[tag.Name()]

	if !ok {
		panic("not found")
	}

	if err := fn(scope, node, tag); err != nil {
		context.AppendContextString(err, "Info: called here", tag.Context())
		return err
	}

	return nil
}

func registerNewLambdaScope() bool {
	functions.NewLambdaScope = func(fnScope_ tokens.Scope, callerScope_ tokens.Scope) functions.LambdaScope {

		fnScope, ok := fnScope_.(Scope)
		if !ok {
			panic("dont know how to sync")
		}

		callerScope, ok := callerScope_.(Scope)
		if !ok {
			panic("dont know how to sync")
		}

		//switch scope_ := scope.(type) {
		//case Scope:
		// use the same node
		//node := callerScope.GetNode()

		subScope := NewSubScope(fnScope, callerScope.GetNode())

		//if _, ok := node.(*RootNode); ok {
		//fmt.Printf("%p -> %p, root node in lambda scope\n", scope, subScope)
		//}
		return subScope
		//default:
		//panic("dont know how to sync")
		//}
	}

	return true
}

var newLambdaScopeRegistered = registerNewLambdaScope()
