package macros

import (
	"strings"

	"../"

	"../../context"
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

type HeaderBuilder struct {
	b strings.Builder
}

func NewHeaderBuilder() *HeaderBuilder {
	return &HeaderBuilder{strings.Builder{}}
}

func (b *HeaderBuilder) String() string {
	return b.b.String()
}

func (b *HeaderBuilder) n() {
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) c(s string) {
	b.b.WriteString(s)
}

func (b *HeaderBuilder) ccc(s1, s2, s3 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
}

func (b *HeaderBuilder) ccccc(s1, s2, s3, s4, s5 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(s4)
	b.b.WriteString(s5)
}

func (b *HeaderBuilder) ccccccc(s1, s2, s3, s4, s5, s6, s7 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(s4)
	b.b.WriteString(s5)
	b.b.WriteString(s6)
	b.b.WriteString(s7)
}

func (b *HeaderBuilder) cccn(s1, s2, s3 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) cccccn(s1, s2, s3, s4, s5 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(s4)
	b.b.WriteString(s5)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) tcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) tcccn(s1, s2, s3 string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) ttcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) ttcccn(s1, s2, s3 string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) tttcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) ttttcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) tttcccn(s1, s2, s3 string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) ttttcccn(s1, s2, s3 string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) tttttcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func (b *HeaderBuilder) ttttttcn(s string) {
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(js.TAB)
	b.b.WriteString(s)
	b.b.WriteString(js.NL)
}

func WriteHeaders() string {
	var b strings.Builder

	// order probably not important due to hoisting
	all := []Header{
		objectFromInstanceHeader,
		objectToInstanceHeader,
		blobFromInstanceHeader,
		blobToInstanceHeader,
		sharedWorkerPostHeader,
		xmlPostHeader,
		dollarHeader,
		webAssemblyEnvHeader,
		searchIndexHeader,
	}

	//ResolveHeaderActivity(SearchIndexHeader, context.NewDummyContext())

	for _, h := range all {
		if h.GetVariable() != nil {
			b.WriteString(h.Write())
		}
	}

	return b.String()
}
