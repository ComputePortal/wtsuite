package macros

import (
  "strings"

  pat "github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

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
	b.b.WriteString(pat.NL)
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
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) cccccn(s1, s2, s3, s4, s5 string) {
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(s4)
	b.b.WriteString(s5)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) t() {
  b.b.WriteString(pat.TAB)
}

func (b *HeaderBuilder) tcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) tcccn(s1, s2, s3 string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) ttcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) ttcccn(s1, s2, s3 string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) tttcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) ttttcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) tttcccn(s1, s2, s3 string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) ttttcccn(s1, s2, s3 string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s1)
	b.b.WriteString(s2)
	b.b.WriteString(s3)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) tttttcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}

func (b *HeaderBuilder) ttttttcn(s string) {
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(pat.TAB)
	b.b.WriteString(s)
	b.b.WriteString(pat.NL)
}
