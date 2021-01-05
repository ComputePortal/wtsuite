package styles

import (
	"encoding/base64"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/math/serif"
)

var MATH_FONT = "FreeSerif"
var MATH_FONT_FAMILY = "FreeSerif"
var MATH_FONT_URL = ""

func listToString(t tokens.Token) (string, error) {
	lst, err := tokens.AssertList(t)
	if err != nil {
		return "", err
	}

	var b strings.Builder

	// separated by spaces
	if err := lst.Loop(func(i int, v_ tokens.Token, last bool) error {
		v, err := tokens.AssertPrimitive(v_)
		if err != nil {
			return err
		}

		b.WriteString(v.Write())

		if !last {
			b.WriteString(" ")
		}

		return nil
	}); err != nil {
		return "", nil
	}

	return b.String(), nil
}

func dictEntryToStringMapEntry(k *tokens.String, v tokens.Token, dst map[string]string) error {
	// null values are ignored in final output
	if tokens.IsNull(v) {
		return nil
	}

	if tokens.IsList(v) {
		str, err := listToString(v)
		if err != nil {
			return err
		}

		dst[k.Value()] = str
	} else {
		value, err := tokens.AssertPrimitive(v)
		if err != nil {
			return err
		}

		dst[k.Value()] = value.Write()
	}

	return nil
}

func dictToStringMap(d *tokens.StringDict,
	ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)

	if err := d.Loop(func(key *tokens.String, val tokens.Token, last bool) error {
		return dictEntryToStringMapEntry(key, val, result)
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func stringMapToString(m map[string]string, nl, indent string) string {
	var b strings.Builder

	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	// sort
	sort.Strings(keys)

	for i, k := range keys {
		v := m[k]

		if k == "display" && v == "grid" {
			b.WriteString(indent)
			b.WriteString("display:-ms-grid;")
			b.WriteString(nl)
		}

		b.WriteString(indent)
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString(v)
		if i == len(keys)-1 {
			b.WriteString(patterns.LAST_SEMICOLON)
		} else {
			b.WriteString(";")
		}
		b.WriteString(nl)
	}

	return b.String()
}

// for __style__ tags
func DictToString(d *tokens.StringDict) (string, error) {
	m, err := dictToStringMap(d, d.Context())
	if err != nil {
		return "", err
	}

	res := stringMapToString(m, "", "")
	return res, nil
}

func SaveMathFont(dst string) error {
	data, err := base64.StdEncoding.DecodeString(serif.Woff2Blob)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, 0644)
}
