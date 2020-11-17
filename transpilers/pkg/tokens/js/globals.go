package js

import (
  "strings"

  pr "./prototypes" // used a lot, so slightly shorter name
  "./values"

  "../context"
)

var TARGET = "nodejs"
// Possible values for TARGET:
//  * browser
//  * worker
//  * nodejs
//  * all (used by refactoring tools)

func handleRegistrationError(err error) {
  if err != nil {
    if TARGET != "all" {
      panic(err)
    }
  }
}

func reserveMacroName(scope Scope, name string) {
  ctx := context.NewDummyContext()

  handleRegistrationError(scope.SetVariable(name, NewVariable(name, true, ctx)))
}

func ReserveMacroNames(scope Scope) {
  reserveMacroName(scope, "cast")
  reserveMacroName(scope, "import")
  reserveMacroName(scope, "SyntaxTree") // to be replaced by __dump__
}

func registerPrototype(scope Scope, proto values.Prototype) {
  ctx := context.NewDummyContext()

  classValue, err := proto.GetClassValue()
  if err != nil {
    panic(err)
  }

  name := proto.Name()

  if strings.ContainsAny(name, "<>") {
    panic("registered prototype can't have type parameters")
  }

  variable := NewVariable(name, true, ctx)
  variable.SetObject(proto)
  variable.SetValue(classValue)

  handleRegistrationError(scope.SetVariable(name, variable))
}

func registerInterface(scope Scope, interfInstance values.Value) {
  ctx := context.NewDummyContext()

  interf := values.GetInterface(interfInstance)
  if interf == nil {
    panic("not an Instance with an Interface")
  }

  name := interf.Name()

  variable := NewVariable(name, true, ctx)
  variable.SetObject(interf)
  variable.SetValue(interfInstance)

  handleRegistrationError(scope.SetVariable(name, variable))
}

func registerValue(scope Scope, name string, v values.Value) {
  ctx := v.Context()
  variable := NewVariable(name, true, ctx)
  variable.SetValue(v)

  handleRegistrationError(scope.SetVariable(name, variable))
}

// scope provided by all js environments
func FillCoreScope(scope Scope) {
  ReserveMacroNames(scope)

  registerPrototype(scope, pr.NewArrayPrototype(nil))
  registerPrototype(scope, pr.NewArrayBufferPrototype())
  registerPrototype(scope, pr.NewBigIntPrototype())
  registerPrototype(scope, pr.NewBooleanPrototype())
  registerPrototype(scope, pr.NewDataViewPrototype())
  registerPrototype(scope, pr.NewDatePrototype())
  registerPrototype(scope, pr.NewErrorPrototype())
  registerPrototype(scope, pr.NewEventPrototype(nil))
  registerPrototype(scope, pr.NewFloat32ArrayPrototype())
  registerPrototype(scope, pr.NewFloat64ArrayPrototype())
  registerPrototype(scope, pr.NewIntPrototype())
  registerPrototype(scope, pr.NewInt8ArrayPrototype())
  registerPrototype(scope, pr.NewInt16ArrayPrototype())
  registerPrototype(scope, pr.NewInt32ArrayPrototype())
  registerPrototype(scope, pr.NewJSONPrototype())
  registerPrototype(scope, pr.NewMapPrototype(nil, nil))
  registerPrototype(scope, pr.NewMathPrototype())
  registerPrototype(scope, pr.NewNumberPrototype())
  registerPrototype(scope, pr.NewObjectPrototype(nil))
  registerPrototype(scope, pr.NewPromisePrototype(nil))
  registerPrototype(scope, pr.NewRegExpPrototype())
  registerPrototype(scope, pr.NewRegExpArrayPrototype())
  registerPrototype(scope, pr.NewSetPrototype(nil))
  registerPrototype(scope, pr.NewStringPrototype())
  registerPrototype(scope, pr.NewUint8ArrayPrototype())
  registerPrototype(scope, pr.NewUint16ArrayPrototype())
  registerPrototype(scope, pr.NewUint32ArrayPrototype())

  ctx := context.NewDummyContext()

  registerValue(scope, "console", pr.NewConsole(ctx))
  registerValue(scope, "setInterval", pr.NewSetTimeoutFunction(ctx))
  registerValue(scope, "setTimeout", pr.NewSetTimeoutFunction(ctx))
}

func FillBrowserAndWorkerCommonScope(scope Scope) {
  FillCoreScope(scope)

  registerPrototype(scope, pr.NewBlobPrototype())
  registerPrototype(scope, pr.NewEventTargetPrototype())
  registerPrototype(scope, pr.NewFileReaderPrototype())
  registerPrototype(scope, pr.NewIDBCursorPrototype())
  registerPrototype(scope, pr.NewIDBCursorWithValuePrototype())
  registerPrototype(scope, pr.NewIDBDatabasePrototype())
  registerPrototype(scope, pr.NewIDBFactoryPrototype())
  registerPrototype(scope, pr.NewIDBIndexPrototype())
  registerPrototype(scope, pr.NewIDBKeyRangePrototype())
  registerPrototype(scope, pr.NewIDBObjectStorePrototype())
  registerPrototype(scope, pr.NewIDBOpenDBRequestPrototype())
  registerPrototype(scope, pr.NewIDBRequestPrototype(nil))
  registerPrototype(scope, pr.NewIDBTransactionPrototype())
  registerPrototype(scope, pr.NewIDBVersionChangeEventPrototype())
  registerPrototype(scope, pr.NewLocationPrototype())
  registerPrototype(scope, pr.NewMessageEventPrototype())
  registerPrototype(scope, pr.NewMessagePortPrototype())
  registerPrototype(scope, pr.NewResponsePrototype())
  registerPrototype(scope, pr.NewTextDecoderPrototype())
  registerPrototype(scope, pr.NewTextEncoderPrototype())
  registerPrototype(scope, pr.NewWebAssemblyPrototype())
  registerPrototype(scope, pr.NewWebAssemblyEnvPrototype())

  ctx := context.NewDummyContext()

  registerInterface(scope, pr.NewWebAssemblyFS(ctx))

  registerValue(scope, "indexedDB", pr.NewIDBFactory(ctx))
  registerValue(scope, "fetch", pr.NewFetchFunction(ctx))
}

func FillWorkerScope(scope Scope) {
  FillBrowserAndWorkerCommonScope(scope)

  registerPrototype(scope, pr.NewDedicatedWorkerGlobalScopePrototype())
  registerPrototype(scope, pr.NewSharedWorkerGlobalScopePrototype())

  ctx := context.NewDummyContext()

  registerValue(scope, "postMessage", pr.NewPostMessageFunction(ctx))
}

func FillBrowserScope(scope Scope) {
  FillBrowserAndWorkerCommonScope(scope)

  registerPrototype(scope, pr.NewCanvasRenderingContext2DPrototype())
  registerPrototype(scope, pr.NewDOMMatrixPrototype())
  registerPrototype(scope, pr.NewDOMRectPrototype())
  registerPrototype(scope, pr.NewElementPrototype())
  registerPrototype(scope, pr.NewFontFaceSetPrototype())
  registerPrototype(scope, pr.NewGLEnumPrototype(""))
  registerPrototype(scope, pr.NewHashChangeEventPrototype())
  registerPrototype(scope, pr.NewHTMLCanvasElementPrototype())
  registerPrototype(scope, pr.NewHTMLCollectionPrototype())
  registerPrototype(scope, pr.NewHTMLElementPrototype())
  registerPrototype(scope, pr.NewHTMLIFrameElementPrototype())
  registerPrototype(scope, pr.NewHTMLImageElementPrototype())
  registerPrototype(scope, pr.NewHTMLInputElementPrototype())
  registerPrototype(scope, pr.NewHTMLLinkElementPrototype())
  registerPrototype(scope, pr.NewHTMLSelectElementPrototype())
  registerPrototype(scope, pr.NewHTMLTextAreaElementPrototype())
  registerPrototype(scope, pr.NewImagePrototype())
  registerPrototype(scope, pr.NewImageDataPrototype())
  registerPrototype(scope, pr.NewKeyboardEventPrototype(nil))
  registerPrototype(scope, pr.NewMouseEventPrototype(nil))
  registerPrototype(scope, pr.NewNodePrototype())
  registerPrototype(scope, pr.NewSearchIndexPrototype())
  registerPrototype(scope, pr.NewSharedWorkerPrototype())
  registerPrototype(scope, pr.NewTextPrototype())
  registerPrototype(scope, pr.NewURLPrototype())
  registerPrototype(scope, pr.NewURLSearchParamsPrototype())
  registerPrototype(scope, pr.NewWebGLBufferPrototype())
  registerPrototype(scope, pr.NewWebGLExtensionPrototype())
  registerPrototype(scope, pr.NewWebGLProgramPrototype())
  registerPrototype(scope, pr.NewWebGLRenderingContextPrototype())
  registerPrototype(scope, pr.NewWebGLShaderPrototype())
  registerPrototype(scope, pr.NewWebGLTexturePrototype())
  registerPrototype(scope, pr.NewWheelEventPrototype(nil))

  ctx := context.NewDummyContext()

  uriFn := values.NewFunction([]values.Value{pr.NewString(ctx), pr.NewString(ctx)}, ctx)

  registerValue(scope, "decodeURIComponent", uriFn)
  registerValue(scope, "encodeURIComponent", uriFn)
  registerValue(scope, "requestIdleCallback", pr.NewRequestIdleCallbackFunction(ctx))
}

func FillNodeJSScope(scope Scope) {
  FillCoreScope(scope)

  registerPrototype(scope, pr.NewNodeJS_BufferPrototype())
  registerPrototype(scope, pr.NewNodeJS_EventEmitterPrototype())

  // packages are added to scope by NodeJSImport statements
}

func FillAllScope(scope Scope) {
  FillWorkerScope(scope)
  FillBrowserScope(scope)
  FillNodeJSScope(scope)
}

func FillGlobalScope(scope Scope) {
  switch TARGET {
  case "nodejs":
    FillNodeJSScope(scope)
  case "browser":
    FillBrowserScope(scope)
  case "worker":
    FillWorkerScope(scope)
  case "all":
    FillAllScope(scope)
  default:
    panic("unrecognized target type")
  }
}

func NewFilledGlobalScope() *GlobalScopeData {
  scope := &GlobalScopeData{newScopeData(nil)}

  FillGlobalScope(scope)

  return scope
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
	b.WriteString(NL)

	return b.String()
}
