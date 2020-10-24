package js

import (
	"fmt"
	"reflect"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type GlobalInstance struct {
	ref   Variable
	proto values.Prototype
	prop  map[string]values.Value
}

type GlobalFunction struct {
	ref Variable
	fn  func(values.Stack, *values.Instance, []values.Value,
		context.Context) (values.Value, error)
}

type GlobalClass struct {
	ref   Variable
	proto values.Prototype
}

type GlobalInterface struct {
	ref    Variable
	interf values.Interface
}

var TARGET = "nodejs"

var globalBrowserInstances map[string]GlobalInstance = make(map[string]GlobalInstance)
var globalBrowserFunctions map[string]GlobalFunction = make(map[string]GlobalFunction)
var globalBrowserClasses map[string]GlobalClass = make(map[string]GlobalClass)
var globalBrowserInterfaces map[string]GlobalInterface = make(map[string]GlobalInterface)

var globalWorkerInstances map[string]GlobalInstance = make(map[string]GlobalInstance)
var globalWorkerFunctions map[string]GlobalFunction = make(map[string]GlobalFunction)
var globalWorkerClasses map[string]GlobalClass = make(map[string]GlobalClass)
var globalWorkerInterfaces map[string]GlobalInterface = make(map[string]GlobalInterface)

var globalNodeJSInstances map[string]GlobalInstance = make(map[string]GlobalInstance)
var globalNodeJSFunctions map[string]GlobalFunction = make(map[string]GlobalFunction)
var globalNodeJSClasses map[string]GlobalClass = make(map[string]GlobalClass)
var globalNodeJSModules map[string]GlobalClass = make(map[string]GlobalClass)
var globalNodeJSInterfaces map[string]GlobalInterface = make(map[string]GlobalInterface)

var globalAllInstances map[string]GlobalInstance = make(map[string]GlobalInstance)
var globalAllFunctions map[string]GlobalFunction = make(map[string]GlobalFunction)
var globalAllClasses map[string]GlobalClass = make(map[string]GlobalClass)
var globalAllInterfaces map[string]GlobalInterface = make(map[string]GlobalInterface)

func registerBrowserGlobals() bool {
	ctx := context.NewDummyContext()

	registerInstance := func(name string, proto values.Prototype,
		prop map[string]values.Value) {
		globalBrowserInstances[name] = GlobalInstance{
			NewVariable(name, true, ctx),
			proto,
			prop,
		}
	}

	registerFunction := func(name string, fn func(values.Stack, *values.Instance,
		[]values.Value, context.Context) (values.Value, error)) {
		globalBrowserFunctions[name] = GlobalFunction{
			NewVariable(name, true, ctx),
			fn,
		}
	}

	registerClass := func(proto values.Prototype) {
		name := proto.Name()
		if name == "" {
			panic("empty builtin class name")
		}

		v := NewVariable(name, true, ctx)
		v.SetObject(proto)
		globalBrowserClasses[name] = GlobalClass{
			v,
			proto,
		}
	}

	registerInstance("console", prototypes.Console, map[string]values.Value{})
	registerInstance("document", prototypes.Document, map[string]values.Value{
		".main": prototypes.NewLiteralBoolean(true, ctx), // contentDocument of iframe is non-main for example
	})
	registerInstance("window", prototypes.Window, map[string]values.Value{})
	registerInstance("indexedDB", prototypes.IDBFactory, map[string]values.Value{})

	registerClass(prototypes.Array)
	registerClass(prototypes.ArrayBuffer)
	registerClass(prototypes.BigInt)
	registerClass(prototypes.Blob)
	registerClass(prototypes.Boolean)
	registerClass(prototypes.CanvasRenderingContext2D)
	registerClass(prototypes.DataView)
	registerClass(prototypes.DOMRect)
	registerClass(prototypes.Element)
	registerClass(prototypes.Error)
	registerClass(prototypes.Event)
	registerClass(prototypes.EventTarget)
	registerClass(prototypes.FileReader)
	registerClass(prototypes.Float32Array)
	registerClass(prototypes.Float64Array)
	registerClass(prototypes.FontFaceSet)
	registerClass(prototypes.GLEnum)
	registerClass(prototypes.HashChangeEvent)
	registerClass(prototypes.HTMLCanvasElement)
	registerClass(prototypes.HTMLCollection)
	registerClass(prototypes.HTMLElement)
	registerClass(prototypes.HTMLIFrameElement)
	registerClass(prototypes.HTMLImageElement)
	registerClass(prototypes.HTMLInputElement)
	registerClass(prototypes.HTMLLinkElement)
	registerClass(prototypes.HTMLSelectElement)
	registerClass(prototypes.HTMLTextAreaElement)
	registerClass(prototypes.IDBDatabase)
	registerClass(prototypes.IDBFactory)
	registerClass(prototypes.IDBIndex)
	registerClass(prototypes.IDBKeyRange)
	registerClass(prototypes.IDBObjectStore)
	registerClass(prototypes.IDBOpenDBRequest)
	registerClass(prototypes.IDBRequest)
	registerClass(prototypes.IDBCursor)
	registerClass(prototypes.IDBCursorWithValue)
	registerClass(prototypes.IDBTransaction)
	registerClass(prototypes.IDBVersionChangeEvent)
	registerClass(prototypes.Image)
	registerClass(prototypes.ImageData)
	registerClass(prototypes.Int)
	registerClass(prototypes.Int8Array)
	registerClass(prototypes.Int16Array)
	registerClass(prototypes.Int32Array)
	registerClass(prototypes.JSON)
	registerClass(prototypes.KeyboardEvent)
	registerClass(prototypes.Location)
	registerClass(prototypes.Map)
	registerClass(prototypes.Math)
	registerClass(prototypes.MessageEvent)
	registerClass(prototypes.MessagePort)
	registerClass(prototypes.MouseEvent)
	registerClass(prototypes.Node)
	registerClass(prototypes.Number)
	registerClass(prototypes.Date)
	registerClass(prototypes.Object)
	registerClass(prototypes.Promise)
	registerClass(prototypes.RegExp)
	registerClass(prototypes.RegExpArray)
	registerClass(prototypes.Response)
	registerClass(prototypes.SearchIndex)
	registerClass(prototypes.Set)
	registerClass(prototypes.SharedWorker)
	registerClass(prototypes.Storage)
	registerClass(prototypes.String)
	registerClass(prototypes.Text)
	registerClass(prototypes.TextDecoder)
	registerClass(prototypes.TextEncoder)
	registerClass(prototypes.TypedArray)
	registerClass(prototypes.Uint8Array)
	registerClass(prototypes.Uint16Array)
	registerClass(prototypes.Uint32Array)
	registerClass(prototypes.URL)
	registerClass(prototypes.URLSearchParams)
	registerClass(prototypes.WebAssembly)
	registerClass(prototypes.WebAssemblyEnv)
	registerClass(prototypes.WebGLBuffer)
	registerClass(prototypes.WebGLProgram)
	registerClass(prototypes.WebGLRenderingContext)
	registerClass(prototypes.WebGLShader)
	registerClass(prototypes.WebGLTexture)
	registerClass(prototypes.WheelEvent)
	registerClass(prototypes.Worker)

	registerFunction("decodeURIComponent", func(stack values.Stack, this *values.Instance,
		args []values.Value, ctx_ context.Context) (values.Value, error) {
		// anything goes?
		if err := prototypes.CheckInputs(prototypes.String, args, ctx_); err != nil {
			return nil, err
		}
		return prototypes.NewInstance(prototypes.String, ctx_), nil
	})

	registerFunction("encodeURIComponent", func(stack values.Stack, this *values.Instance,
		args []values.Value, ctx_ context.Context) (values.Value, error) {
		// anything goes?
		if err := prototypes.CheckInputs(prototypes.String, args, ctx_); err != nil {
			return nil, err
		}
		return prototypes.NewInstance(prototypes.String, ctx_), nil
	})

	// for convenience
	registerFunction("fetch", prototypes.WindowFetch)
	registerFunction("requestIdleCallback", prototypes.WindowRequestIdleCallback)
	registerFunction("setTimeout", prototypes.WindowSetTimeout)
	registerFunction("setInterval", prototypes.WindowSetTimeout)

	return true
}

func registerWorkerGlobals() bool {
	ctx := context.NewDummyContext()

	registerInstance := func(name string, proto values.Prototype,
		prop map[string]values.Value) {
		globalWorkerInstances[name] = GlobalInstance{
			NewVariable(name, true, ctx),
			proto,
			prop,
		}
	}

	registerFunction := func(name string, fn func(values.Stack, *values.Instance,
		[]values.Value, context.Context) (values.Value, error)) {
		globalWorkerFunctions[name] = GlobalFunction{
			NewVariable(name, true, ctx),
			fn,
		}
	}

	registerClass := func(proto values.Prototype) {
		name := proto.Name()
		if name == "" {
			panic("empty builtin class name")
		}

		v := NewVariable(name, true, ctx)
		v.SetObject(proto)
		globalWorkerClasses[name] = GlobalClass{
			v,
			proto,
		}
	}

	registerInterface := func(interf values.Interface) {
		name := interf.Name()
		if name == "" {
			panic("empty builtin interface name")
		}

		v := NewVariable(name, true, ctx)
		v.SetObject(interf)
		globalWorkerInterfaces[name] = GlobalInterface{
			v,
			interf,
		}
	}

	registerInstance("console", prototypes.Console, map[string]values.Value{})
	registerInstance("indexedDB", prototypes.IDBFactory, map[string]values.Value{})

	registerClass(prototypes.Array)
	registerClass(prototypes.ArrayBuffer)
	registerClass(prototypes.BigInt)
	registerClass(prototypes.Blob)
	registerClass(prototypes.Boolean)
	registerClass(prototypes.DataView)
	registerClass(prototypes.DedicatedWorkerGlobalScope)
	registerClass(prototypes.Error)
	registerClass(prototypes.Event)
	registerClass(prototypes.EventTarget)
	registerClass(prototypes.FileReader)
	registerClass(prototypes.Float32Array)
	registerClass(prototypes.Float64Array)
	registerClass(prototypes.IDBDatabase)
	registerClass(prototypes.IDBFactory)
	registerClass(prototypes.IDBIndex)
	registerClass(prototypes.IDBKeyRange)
	registerClass(prototypes.IDBObjectStore)
	registerClass(prototypes.IDBOpenDBRequest)
	registerClass(prototypes.IDBRequest)
	registerClass(prototypes.IDBCursor)
	registerClass(prototypes.IDBCursorWithValue)
	registerClass(prototypes.IDBTransaction)
	registerClass(prototypes.IDBVersionChangeEvent)
	registerClass(prototypes.Int)
	registerClass(prototypes.Int8Array)
	registerClass(prototypes.Int16Array)
	registerClass(prototypes.Int32Array)
	registerClass(prototypes.JSON)
	registerClass(prototypes.Location)
	registerClass(prototypes.Map)
	registerClass(prototypes.Math)
	registerClass(prototypes.MessageEvent)
	registerClass(prototypes.MessagePort)
	registerClass(prototypes.Number)
	registerClass(prototypes.Date)
	registerClass(prototypes.Object)
	registerClass(prototypes.Promise)
	registerClass(prototypes.RegExp)
	registerClass(prototypes.RegExpArray)
	registerClass(prototypes.Set)
	registerClass(prototypes.SharedWorkerGlobalScope)
	registerClass(prototypes.String)
	registerClass(prototypes.Text)
	registerClass(prototypes.TextDecoder)
	registerClass(prototypes.TextEncoder)
	registerClass(prototypes.TypedArray)
	registerClass(prototypes.Uint8Array)
	registerClass(prototypes.Uint16Array)
	registerClass(prototypes.Uint32Array)
	registerClass(prototypes.WebAssembly)
	registerClass(prototypes.WebAssemblyEnv)

	registerInterface(prototypes.WebAssemblyFS)

	registerFunction("postMessage", func(stack values.Stack, this *values.Instance,
		args []values.Value, ctx_ context.Context) (values.Value, error) {
		// anything goes?
		if err := prototypes.CheckInputs(&prototypes.Any{}, args, ctx_); err != nil {
			return nil, err
		}
		return nil, nil
	})

	registerFunction("setTimeout", prototypes.WindowSetTimeout)
	registerFunction("setInterval", prototypes.WindowSetTimeout)
	registerFunction("fetch", prototypes.WindowFetch)

	return true
}

func registerNodeJSGlobals() bool {
	ctx := context.NewDummyContext()

	registerInstance := func(name string, proto values.Prototype,
		prop map[string]values.Value) {
		globalNodeJSInstances[name] = GlobalInstance{
			NewVariable(name, true, ctx),
			proto,
			prop,
		}
	}

	registerFunction := func(name string, fn func(values.Stack, *values.Instance,
		[]values.Value, context.Context) (values.Value, error)) {
		globalNodeJSFunctions[name] = GlobalFunction{
			NewVariable(name, true, ctx),
			fn,
		}
	}

	registerClass := func(proto values.Prototype) {
		name := proto.Name()
		if name == "" {
			fmt.Println(reflect.TypeOf(proto).String())
			panic("empty builtin class name")
		}

		v := NewVariable(name, true, ctx)
		v.SetObject(proto)
		globalNodeJSClasses[name] = GlobalClass{
			v,
			proto,
		}
	}

	registerModule := func(proto values.Prototype) {
		name := proto.Name()
		if name == "" {
			fmt.Println(reflect.TypeOf(proto).String())
			panic("empty builtin module name")
		}

		v := NewVariable(name, true, ctx)
		v.SetObject(proto)
		globalNodeJSModules[name] = GlobalClass{
			v,
			proto,
		}
	}

	registerInstance("console", prototypes.Console, map[string]values.Value{})

	registerClass(prototypes.Array)
	registerClass(prototypes.ArrayBuffer)
	registerClass(prototypes.BigInt)
	registerClass(prototypes.Boolean)
	registerClass(prototypes.Buffer)
	registerClass(prototypes.DataView)
	registerClass(prototypes.Error)
	registerClass(prototypes.Event)
	registerClass(prototypes.Float32Array)
	registerClass(prototypes.Float64Array)
	registerClass(prototypes.Int)
	registerClass(prototypes.JSON)
	registerClass(prototypes.Map)
	registerClass(prototypes.Math)
	registerClass(prototypes.Number)
	registerClass(prototypes.Date)
	registerClass(prototypes.Object)
	registerClass(prototypes.Path)
	registerClass(prototypes.Promise)
	registerClass(prototypes.RegExp)
	registerClass(prototypes.RegExpArray)
	registerClass(prototypes.Set)
	registerClass(prototypes.String)
	registerClass(prototypes.TypedArray)
	registerClass(prototypes.Uint16Array)
	registerClass(prototypes.Uint32Array)
	registerClass(prototypes.Uint8Array)

	// needed so they are available in the stack, for the scope it doesn't matter because reference names will differ anyway
	registerClass(prototypes.NodeJS_crypto_Cipher)
	registerClass(prototypes.NodeJS_crypto_Decipher)
	registerClass(prototypes.NodeJS_EventEmitter)
	registerClass(prototypes.NodeJS_http_IncomingMessage)
	registerClass(prototypes.NodeJS_http_Server)
	registerClass(prototypes.NodeJS_http_ServerResponse)
  registerClass(prototypes.NodeJS_mysql_Connection)
  registerClass(prototypes.NodeJS_mysql_Error)
  registerClass(prototypes.NodeJS_mysql_FieldPacket)
  registerClass(prototypes.NodeJS_mysql_Pool)
  registerClass(prototypes.NodeJS_mysql_Query)
	registerClass(prototypes.NodeJS_nodemailer_SMTPTransport)
	registerClass(prototypes.NodeJS_stream_Readable)

	registerModule(prototypes.NodeJS_crypto)
	registerModule(prototypes.NodeJS_fs)
	registerModule(prototypes.NodeJS_http)
	registerModule(prototypes.NodeJS_nodemailer)
	registerModule(prototypes.NodeJS_mysql)
	registerModule(prototypes.NodeJS_process)
	registerModule(prototypes.NodeJS_stream)

	registerFunction("setTimeout", prototypes.WindowSetTimeout)
	registerFunction("setInterval", prototypes.WindowSetTimeout)

	return true
}

// used when refactoring (so there is no need to specify the target, and we can refactor mixed targets at once)
func registerAllGlobals() bool {
  // prefer NodeJS, then Browser, (Worker is least important)
  for key, obj := range globalWorkerInstances {
    globalAllInstances[key] = obj
  }

  for key, obj := range globalWorkerFunctions {
    globalAllFunctions[key] = obj
  }

  for key, obj := range globalWorkerClasses {
    globalAllClasses[key] = obj
  }

  for key, obj := range globalWorkerInterfaces {
    globalAllInterfaces[key] = obj
  }

  for key, obj := range globalBrowserInstances {
    globalAllInstances[key] = obj
  }

  for key, obj := range globalBrowserFunctions {
    globalAllFunctions[key] = obj
  }

  for key, obj := range globalBrowserClasses {
    globalAllClasses[key] = obj
  }

  for key, obj := range globalBrowserInterfaces {
    globalAllInterfaces[key] = obj
  }

  for key, obj := range globalNodeJSInstances {
    globalAllInstances[key] = obj
  }

  for key, obj := range globalNodeJSFunctions {
    globalAllFunctions[key] = obj
  }

  for key, obj := range globalNodeJSClasses {
    globalAllClasses[key] = obj
  }

  for key, obj := range globalNodeJSInterfaces {
    globalAllInterfaces[key] = obj
  }

  return true
}

var _BrowserGlobalsOk bool = registerBrowserGlobals()

var _WorkerGlobalsOk bool = registerWorkerGlobals()

var _NodeJSGlobalsOk bool = registerNodeJSGlobals()

var _AllGlobalsOk bool = registerAllGlobals()

func newFilledGlobalScope(
	globalInstances map[string]GlobalInstance,
	globalFunctions map[string]GlobalFunction,
	globalClasses map[string]GlobalClass,
	globalInterfaces map[string]GlobalInterface) *GlobalScopeData {
	scope := &GlobalScopeData{newScopeData(nil)}

	for name, instance := range globalInstances {
		if err := scope.SetVariable(name, instance.ref); err != nil {
			panic(err)
		}
	}

	for name, fn := range globalFunctions {
		if err := scope.SetVariable(name, fn.ref); err != nil {
			panic(err)
		}
	}

	for name, class := range globalClasses {
		if err := scope.SetVariable(name, class.ref); err != nil {
			panic(err)
		}
	}

	for name, interf := range globalInterfaces {
		if err := scope.SetVariable(name, interf.ref); err != nil {
			panic(err)
		}
	}

	return scope
}

func NewFilledGlobalScope() *GlobalScopeData {
	switch TARGET {
  case "all":
		return newFilledGlobalScope(
			globalAllInstances,
			globalAllFunctions,
			globalAllClasses,
			globalAllInterfaces,
		)
	case "nodejs":
		return newFilledGlobalScope(
			globalNodeJSInstances,
			globalNodeJSFunctions,
			globalNodeJSClasses,
			globalNodeJSInterfaces,
		)
	case "browser":
		return newFilledGlobalScope(
			globalBrowserInstances,
			globalBrowserFunctions,
			globalBrowserClasses,
			globalBrowserInterfaces,
		)
	case "worker":
		return newFilledGlobalScope(
			globalWorkerInstances,
			globalWorkerFunctions,
			globalWorkerClasses,
			globalWorkerInterfaces,
		)
	default:
		panic("unrecognized TARGET " + TARGET)
	}
}

func fillGlobalStack(stack values.Stack,
	globalInstances map[string]GlobalInstance,
	globalFunctions map[string]GlobalFunction,
	globalClasses map[string]GlobalClass,
	globalInterfaces map[string]GlobalInterface,
) {
	ctx := context.NewDummyContext()

	for _, instance := range globalInstances {
		val := values.NewInstance(instance.proto,
			values.NewPropertiesWithContent(instance.prop, ctx), ctx)
		if err := stack.SetValue(instance.ref, val, false, ctx); err != nil {
			panic(err)
		}
	}

	for _, fn := range globalFunctions {
		val := values.NewFunctionFunction(fn.fn, stack, nil, ctx)
		if err := stack.SetValue(fn.ref, val, false, ctx); err != nil {
			panic(err)
		}
	}

	for _, class := range globalClasses {
		val := values.NewClass(class.proto, ctx)
		if err := stack.SetValue(class.ref, val, false, ctx); err != nil {
			panic(err)
		}
	}

	for _, interf := range globalInterfaces {
		val := values.NewClassInterface(interf.interf, ctx)
		if err := stack.SetValue(interf.ref, val, false, ctx); err != nil {
			panic(err)
		}
	}
}

func NewFilledGlobalStack(cs *CacheStack) *GlobalStack {
	stack := NewGlobalStack(cs)

	switch TARGET {
  case "all":
		fillGlobalStack(
			stack,
			globalAllInstances,
			globalAllFunctions,
			globalAllClasses,
			globalAllInterfaces,
		)
	case "nodejs":
		fillGlobalStack(
			stack,
			globalNodeJSInstances,
			globalNodeJSFunctions,
			globalNodeJSClasses,
			globalNodeJSInterfaces,
		)
	case "browser":
		fillGlobalStack(
			stack,
			globalBrowserInstances,
			globalBrowserFunctions,
			globalBrowserClasses,
			globalBrowserInterfaces,
		)
	case "worker":
		fillGlobalStack(
			stack,
			globalWorkerInstances,
			globalWorkerFunctions,
			globalWorkerClasses,
			globalWorkerInterfaces,
		)
	default:
		panic("unrecognized TARGET " + TARGET)
	}

  return stack
}

func GetNodeJSModule(name string) (*prototypes.BuiltinPrototype, bool) {
	if TARGET == "nodejs" || TARGET == "all" {
		if gModule, ok := globalNodeJSModules[name]; ok {
			proto_ := gModule.proto
			if proto, ok := proto_.(*prototypes.BuiltinPrototype); ok {
				return proto, true
			}
		}
	}

	return nil, false
}

// used to emulate nodejs modules
func GetNodeJSClassVariable(name string) (Variable, bool) {
	if TARGET == "nodejs" || TARGET == "all" {
		if gClass, ok := globalNodeJSClasses[name]; ok {
			return gClass.ref, true
		}
	}

	return nil, false
}

func IsBuiltinName(name string) bool {
	fn := func(cm map[string]GlobalClass, fm map[string]GlobalFunction, im map[string]GlobalInstance) bool {
		if _, ok := cm[name]; ok {
			return true
		} else if _, ok := fm[name]; ok {
			return true
		} else if _, ok := im[name]; ok {
			return true
		} else {
			return false
		}
	}

	// interface names dont matter
	switch TARGET {
  case "all":
		return fn(globalAllClasses, globalAllFunctions, globalAllInstances)
	case "nodejs":
		return fn(globalNodeJSClasses, globalNodeJSFunctions, globalNodeJSInstances)
	case "browser":
		return fn(globalBrowserClasses, globalBrowserFunctions, globalBrowserInstances)
	case "worker":
		return fn(globalWorkerClasses, globalWorkerFunctions, globalWorkerInstances)
	default:
		panic("unrecognized TARGET " + TARGET)
	}
}

func GetBuiltinPrototype(key string) values.Prototype {
	fn := func(m map[string]GlobalClass) values.Prototype {
		if cl, ok := m[key]; !ok {
			return nil
		} else {
			return cl.proto
		}
	}

	switch TARGET {
  case "all":
		return fn(globalAllClasses)
	case "nodejs":
		return fn(globalNodeJSClasses)
	case "browser":
		return fn(globalBrowserClasses)
	case "worker":
		return fn(globalWorkerClasses)
	default:
		panic("unrecognized TARGET " + TARGET)
	}
}

func WriteGlobalHeaders() string {
	var b strings.Builder

  if TARGET == "all" {
    panic("js.TARGET can't be used for printing")
  }
	if TARGET == "nodejs" {
		b.WriteString("'use strict'\n")
	}

	b.WriteString("class Int extends Number{")
	b.WriteString(NL)
	b.WriteString(TAB)
	b.WriteString("constructor(x){super(parseInt(x))}")
	b.WriteString(NL)
	b.WriteString("}")

	return b.String()
}
