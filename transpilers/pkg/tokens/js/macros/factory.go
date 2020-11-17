package macros

import (
	"../"

	"../../context"
)

var _classMacros = map[string]MacroGroup{
	"SyntaxTree": MacroGroup{
		macros: map[string]MacroConstructor{
      // TODO: probably better as underscore call macros
			"info":  NewSyntaxTreeInfo,
		},
	},

	"Math": MacroGroup{
		macros: map[string]MacroConstructor{
			"advanceWidth":  NewMathAdvanceWidth,
			"boundingBox":   NewMathBoundingBox,
			"degToRad":      NewDegToRad,
			"radToDeg":      NewRadToDeg,
			"formatMetrics": NewMathFormatMetrics,
		},
	},

	"Blob": MacroGroup{
		macros: map[string]MacroConstructor{
			"toInstance":   NewBlobToInstance,
			"fromInstance": NewBlobFromInstance,
		},
	},

	"Object": MacroGroup{
		macros: map[string]MacroConstructor{
			"toInstance":   NewObjectToInstance,
			"fromInstance": NewObjectFromInstance,
			"isUndefined":  NewIsUndefined,
		},
	},

	"SharedWorker": MacroGroup{
		macros: map[string]MacroConstructor{
			"post": NewSharedWorkerPost,
		},
	},

	"URL": MacroGroup{
		macros: map[string]MacroConstructor{
			"current": NewURLCurrent,
		},
	},

	"WebAssembly": MacroGroup{
		macros: map[string]MacroConstructor{
			"exec": NewWebAssemblyExec,
			// "load": NewWebAssemblyLoad, // TODO
		},
	},

	"XMLHttpRequest": MacroGroup{
		macros: map[string]MacroConstructor{
			"post": NewXMLHttpRequestPost,
		},
	},
}

var _callMacros = map[string]MacroConstructor{
  "cast":   NewCastCall,
	"BigInt": NewBigIntCall,
	//"WebAssemblyEnv": NewWebAssemblyEnvCall,
}

func IsClassMacroGroup(gname string) bool {
	_, ok := _classMacros[gname]
	return ok
}

func IsClassMacro(gname string, name string) bool {
	if mg, ok := _classMacros[gname]; ok {
		_, ok = mg.macros[name]
		return ok
	} else {
		return false
	}
}

func IsCallMacro(name string) bool {
	_, ok := _callMacros[name]
	return ok
}

func MemberIsClassMacro(m *js.Member) bool {
	if name, key := m.ObjectNameAndKey(); name != "" {
		return IsClassMacro(name, key)
	}

	return false
}

func CallIsCallMacro(call *js.Call) bool {
	if name := call.Name(); name != "" {
		return IsCallMacro(name)
	}

	return false
}

func NewParseTime(args []js.Expression, ctx context.Context) (js.Expression, error) {
	panic(ctx.NewError("Internal Error: should be absorbed at parse time"))
}

func NewClassMacro(gname string, name string, args []js.Expression,
	ctx context.Context) (js.Expression, error) {
	return _classMacros[gname].macros[name](args, ctx)
}

func NewClassMacroFromMember(m *js.Member, args []js.Expression,
	ctx context.Context) (js.Expression, error) {
	if name, key := m.ObjectNameAndKey(); name != "" {
		return NewClassMacro(name, key, args, ctx)
	} else {
		panic("unhandled")
	}
}

func NewCallMacro(name string, args []js.Expression,
	ctx context.Context) (js.Expression, error) {
	return _callMacros[name](args, ctx)
}

func NewCallMacroFromCall(call *js.Call,
	ctx context.Context) (js.Expression, error) {
	name := call.Name()
	if name == "" {
		panic("should've been handled before")
	}

	args := call.Args()

	return NewCallMacro(name, args, ctx)
}

func RegisterActivateMacroHeadersCallback() bool {
	js.ActivateMacroHeaders = func(name string) {
		switch name {
		case "WebAssemblyEnv":
			ActivateWebAssemblyEnvHeader()
		case "SearchIndex":
			ActivateSearchIndexHeader()
		}
	}

	return true
}

var _activateMacroHeadersCallbackOk = RegisterActivateMacroHeadersCallback()
