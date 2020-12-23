package parsers

import (
	"regexp"
	"sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

type TypeBit int

const (
	SING TypeBit = 1 << 1
	PRE          = 1 << 2
	POST         = 1 << 3
	BIN          = 1 << 4
	L2R          = 1 << 5
	TER          = 1 << 6
)

type quotedGroupSettings struct {
	maskType        RuneMask
	groupPattern    patterns.Group
	assertStopMatch bool
	info            string
	trackStarts     bool
}

type quotedGroupsSettings struct {
	pattern *regexp.Regexp
	groups  []quotedGroupSettings
}

type formulasSettings struct {
	tokenizer func(string, context.Context) ([]raw.Token, error)
}

type wordsAndLiteralsSettings struct {
	maskType  RuneMask
	pattern   *regexp.Regexp
	tokenizer func(string, context.Context) (raw.Token, error)
}

type symbolsSettings struct {
	maskType RuneMask
	pattern  *regexp.Regexp
}

type operatorSettings struct {
	precedence int    // lower -> sooner, same precedence: right to left
	symbol     string // can also be KeyWords!
	typeBits   TypeBit
}

type operatorSettingsMap struct { // all with same precedence
	m map[string]operatorSettings
}

type operatorsSettings struct {
	sortedOperators []operatorSettingsMap // sorted by precedence, key is symbol
}

type ParserSettings struct {
	quotedGroups             quotedGroupsSettings
	formulas                 formulasSettings
	wordsAndLiterals         wordsAndLiteralsSettings
	symbols                  symbolsSettings
	operators                operatorsSettings
	tmpGroupWords            bool
	tmpGroupPeriods          bool
	tmpGroupArrows           bool
	tmpGroupDColons          bool
	tmpGroupAngled           bool
	recursivelyNestOperators bool
}

func newOperatorsSettings(operators []operatorSettings) operatorsSettings {
	// build the cache
	sort.Slice(operators, func(i, j int) bool {
		return operators[i].precedence > operators[j].precedence // higher precedence is treated sooner
	})

	cachedOperators := make([]operatorSettingsMap, 0)

	idst := -1
	prevPrecedence := -1
	for _, op := range operators {
		if op.precedence < 0 {
			panic("negative precedence is illegal")
		} else if op.precedence == prevPrecedence {
			cachedOperators[idst].m[op.symbol] = op
		} else {
			idst += 1
			if idst > len(cachedOperators)-1 {
				cachedOperators = append(cachedOperators,
					operatorSettingsMap{make(map[string]operatorSettings)})
			}

			symbols := strings.Fields(op.symbol)
			if len(symbols) == 0 {
				panic("unexpected")
			}

			// '?' for ternary operator L2R (or ':' if R2L)
			firstSymbol := symbols[len(symbols)-1]
			if op.typeBits&L2R > 0 {
				firstSymbol = symbols[0]
			}

			cachedOperators[idst].m[firstSymbol] = op
			prevPrecedence = op.precedence
		}
	}

	return operatorsSettings{cachedOperators}
}

func (osm *operatorSettingsMap) LeftToRight() bool {
	l2r := true

	someLeftToRight := false
	for _, v := range osm.m {
		if (v.typeBits & L2R) > 0 {
			someLeftToRight = true
		} else {
			l2r = false
		}
	}

	if someLeftToRight && (!l2r) {
		panic("bad operator settings, all operators in group must have LeftToRight associacivity or vice-versa")
	}

	return l2r
}

func (osm *operatorSettingsMap) Has(t raw.Token) bool {
	if s, ok := t.(*raw.Symbol); ok {
		if _, ok := osm.m[s.Value()]; ok {
			return true
		}
	}

	return false
}

func (osm *operatorSettingsMap) getName(t raw.Token, fn func(ns *operatorSettings, symbol string) string) (string, bool) {
	if s, ok := t.(*raw.Symbol); ok {
		if ns, ok := osm.m[s.Value()]; ok {
			n := fn(&ns, s.Value())
			if n != "" {
				return n, true
			}
		}
	}

	return "", false
}

func (osm *operatorSettingsMap) GetSingular(t raw.Token) (string, bool) {
	return osm.getName(t, func(ns *operatorSettings, symbol string) string {
		if (ns.typeBits & SING) > 0 {
			return "sing" + symbol
		} else {
			return ""
		}
	})
}

func (osm *operatorSettingsMap) GetPreUnary(t raw.Token) (string, bool) {
	return osm.getName(t, func(ns *operatorSettings, symbol string) string {
		if (ns.typeBits & PRE) > 0 {
			return "pre" + symbol
		} else {
			return ""
		}
	})
}

func (osm *operatorSettingsMap) GetPostUnary(t raw.Token) (string, bool) {
	return osm.getName(t, func(ns *operatorSettings, symbol string) string {
		if (ns.typeBits & POST) > 0 {
			return "post" + symbol
		} else {
			return ""
		}
	})
}

func (osm *operatorSettingsMap) GetBinary(t raw.Token) (string, bool) {
	return osm.getName(t, func(ns *operatorSettings, symbol string) string {
		if (ns.typeBits & BIN) > 0 {
			return "bin" + symbol
		} else {
			return ""
		}
	})
}

func (osm *operatorSettingsMap) getTernary(t0, t1 raw.Token, l2r bool) (string, bool) {
	s0, ok0 := t0.(*raw.Symbol)
	s1, ok1 := t1.(*raw.Symbol)

	if ok0 && ok1 {
		firstSymbol := s1.Value()
		if l2r {
			firstSymbol = s0.Value()
		}

		if ns, ok := osm.m[firstSymbol]; ok {
			if (ns.typeBits&TER) > 0 && (ns.symbol == s0.Value()+" "+s1.Value()) {
				return "ter" + ns.symbol, true
			}
		}
	}

	return "", false
}

func (osm *operatorSettingsMap) GetTernaryR2L(t0, t1 raw.Token) (string, bool) {
	return osm.getTernary(t0, t1, false)
}

func (osm *operatorSettingsMap) GetTernaryL2R(t0, t1 raw.Token) (string, bool) {
	return osm.getTernary(t0, t1, true)
}
