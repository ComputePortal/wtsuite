package config

import (
	"errors"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tree"
)

type XQuery interface {
	Match([]tree.Tag) bool
	Length() int
}

type XQueryData struct {
	immediate bool
	name      string
	child     *XQueryData
}

func ParseXQuery(source string) (*XQueryData, error) {
	fields := strings.Fields(source)

	n := len(fields)
	if n == 0 {
		return nil, errors.New("Error: xquery is empty")
	}

	var node *XQueryData = nil
	for i := n - 1; i >= 0; i-- {
		f := fields[i]
		if f == ">" {
			return nil, errors.New("Error: bad >")
		}

		immediate := false
		if i > 0 && fields[i-1] == ">" {
			i--
			immediate = true
		}

		node = &XQueryData{immediate, f, node}
	}

	return node, nil
}

func (q *XQueryData) Length() int {
	if q.child == nil {
		return 1
	} else {
		return 1 + q.child.Length()
	}
}

func (q *XQueryData) Match(tags []tree.Tag) bool {
	l := q.Length()

	if len(tags) < l {
		return false
	}

	if q.immediate && q.name != tags[0].Name() {
		return false
	}

	for i, tag := range tags {
		if tag.Name() == q.name || q.name == "*" {
			if q.child != nil {

				if q.child.Match(tags[i+1:]) {
					return true
				}
			} else {
				if len(tags)-i == 1 {
					return true
				}
			}
		}
	}

	if !q.immediate && len(tags) > 1 {
		return q.Match(tags[1:])
	}

	return false
}

func DumpXPath(tags []tree.Tag) string {
	var b strings.Builder

	for i, tag := range tags {
		b.WriteString(tag.Name())
		if i < len(tags)-1 {
			b.WriteString(" ")
		}
	}

	return b.String()
}
