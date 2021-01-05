package macros

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
)

type Header interface {
	Dependencies() []Header

	Name() string
	GetVariable() js.Variable
	SetVariable(v js.Variable)

	UniqueNames(ns js.Namespace) error

	Write() string
}

type HeaderData struct {
	name string
	v    js.Variable
}

func newHeaderData(name string) HeaderData {
	return HeaderData{name, nil}
}

func (h *HeaderData) Name() string {
	return h.name
}

func (h *HeaderData) GetVariable() js.Variable {
	return h.v
}

func (h *HeaderData) SetVariable(v js.Variable) {
	h.v = v
}

func ResolveHeaderActivity(h Header, ctx context.Context) {
	for _, other := range h.Dependencies() {
		ResolveHeaderActivity(other, ctx)
	}

	if h.GetVariable() == nil {
		h.SetVariable(js.NewVariable(h.Name(), true, ctx))
	}
}

func (h *HeaderData) UniqueNames(ns js.Namespace) error {
	return ns.LibName(h.GetVariable(), h.Name())
}

func UniqueHeaderNames(h Header, ns js.Namespace) error {
	for _, other := range h.Dependencies() {
		if err := UniqueHeaderNames(other, ns); err != nil {
			return err
		}
	}

	return h.UniqueNames(ns)
}

func WriteHeaders() string {
	var b strings.Builder

	// order probably not important due to hoisting
	all := []Header{
    checkTypeHeader,
		objectFromInstanceHeader,
		objectToInstanceHeader,
		blobFromInstanceHeader,
		blobToInstanceHeader,
		sharedWorkerPostHeader,
		xmlPostHeader,
    rpcContextHeader,
    rpcClientHeader,
    rpcServerHeader,
		webAssemblyEnvHeader,
    webGLProgramHeader,
		searchIndexHeader,
    mathFontHeader,
	}

	for _, h := range all {
		if h.GetVariable() != nil {
			b.WriteString(h.Write())
		}
	}

	return b.String()
}
