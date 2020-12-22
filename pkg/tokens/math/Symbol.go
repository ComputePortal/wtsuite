package math

import (
	"fmt"

	"./boundingbox"

	"../context"
)

type Symbol struct {
	unicode int
	symbol  string
	Word
}

var symbolMap = map[string]int{
	// group symbols,
	"(": 0x0028,
	")": 0x0029,
	"[": 0x005b,
	"]": 0x005d,
	"{": 0x007b,
	"}": 0x007d,

	// separators
	",": 0x0002c,
	";": 0x0003b,

	"/": 0x0002f,
	"<": 0x0003c,
	">": 0x0003e,

	// operators
	"nabla": 0x02207,
	"-":     0x02212,

	// lating upper case
	"A": 0x1d434,
	"B": 0x1d435,
	"C": 0x1d436,
	"D": 0x1d437,
	"E": 0x1d438,
	"F": 0x1d439,
	"G": 0x1d43a,
	"H": 0x1d43b,
	"I": 0x1d43c,
	"J": 0x1d43d,
	"K": 0x1d43e,
	"L": 0x1d43f,
	"M": 0x1d440,
	"N": 0x1d441,
	"O": 0x1d442,
	"P": 0x1d443,
	"Q": 0x1d444,
	"R": 0x1d445,
	"S": 0x1d446,
	"T": 0x1d447,
	"U": 0x1d448,
	"V": 0x1d449,
	"W": 0x1d44a,
	"X": 0x1d44b,
	"Y": 0x1d44c,
	"Z": 0x1d44d,

	// latin lower case
	"a": 0x1d44e,
	"b": 0x1d44f,
	"c": 0x1d450,
	"d": 0x1d451,
	"e": 0x1d452,
	"f": 0x1d453,
	"g": 0x1d454,
	"h": 0x0210e,
	"i": 0x1d456,
	"j": 0x1d457,
	"k": 0x1d458,
	"l": 0x1d459,
	"m": 0x1d45a,
	"n": 0x1d45b,
	"o": 0x1d45c,
	"p": 0x1d45d,
	"q": 0x1d45e,
	"r": 0x1d45f,
	"s": 0x1d460,
	"t": 0x1d461,
	"u": 0x1d462,
	"v": 0x1d463,
	"w": 0x1d464,
	"x": 0x1d465,
	"y": 0x1d466,
	"z": 0x1d467,

	// greek upper case (not italic!)
	// "Alpha":      0x00391, // too similar to latin
	// "Beta":       0x00392, // too similar to latin
	"Gamma":   0x00393,
	"Delta":   0x00394,
	"Theta":   0x00398,
	"Lambda":  0x0039b,
	"Xi":      0x0039e,
	"Pi":      0x003a0,
	"Sigma":   0x003a3,
	"Upsilon": 0x003a5,
	"Phi":     0x003a6,
	"Psi":     0x003a8,
	"Omega":   0x003a9,

	// greek lower case (use italic versions)
	"alpha":      0x1d6fc,
	"beta":       0x1d6fd,
	"gamma":      0x1d6fe,
	"delta":      0x1d6ff,
	"varepsilon": 0x1d700,
	"zeta":       0x1d701,
	"eta":        0x1d702,
	"theta":      0x1d703,
	"iota":       0x1d704,
	"kappa":      0x1d705, // very similar to k
	"lambda":     0x1d706,
	"mu":         0x1d707,
	"nu":         0x1d708,
	"xi":         0x1d709,
	//"omicron": 0x1d70a, // identical to latin o, and indeed not available in tex
	"pi":       0x1d70b,
	"rho":      0x1d70c,
	"varsigma": 0x1d70d,
	"sigma":    0x1d70e,
	"tau":      0x1d70f,
	"upsilon":  0x1d710,
	"varphi":   0x1d711,
	"chi":      0x1d712,
	"psi":      0x1d713,
	"omega":    0x1d714,
	"epsilon":  0x1d716,
	"varkappa": 0x1d718,
	"phi":      0x1d719,
	"varrho":   0x1d71a,

	"infty": 0x0221e,
}

var boldMap = map[int]int{
	0x1d44e: 0x1d482, // a
	0x1d44f: 0x1d483, // b
	0x1d450: 0x1d484, // c
	0x1d451: 0x1d485, // d
	0x1d452: 0x1d486, // e
	0x1d453: 0x1d487, // f
	0x1d454: 0x1d488, // g
	0x0210e: 0x1d489, // h
	0x1d456: 0x1d48a, // i
	0x1d457: 0x1d48b, // j
	0x1d458: 0x1d48c, // k
	0x1d459: 0x1d48d, // l
	0x1d45a: 0x1d48e, // m
	0x1d45b: 0x1d48f, // n
	0x1d45c: 0x1d490, // o
	0x1d45d: 0x1d491, // p
	0x1d45e: 0x1d492, // q
	0x1d45f: 0x1d493, // r
	0x1d460: 0x1d494, // s
	0x1d461: 0x1d495, // t
	0x1d462: 0x1d496, // u
	0x1d463: 0x1d497, // v
	0x1d464: 0x1d498, // w
	0x1d465: 0x1d499, // x
	0x1d466: 0x1d49a, // y
	0x1d467: 0x1d49b, // z
}

func NewUnicodeSymbol(inputValue string, unicode int, ctx context.Context) (*Symbol, error) {
	value := fmt.Sprintf("&#x%x;", unicode)
	return &Symbol{unicode, inputValue, Word{value, newTokenData(ctx)}}, nil
}

func NewSymbol(value string, ctx context.Context) (*Symbol, error) {
	if unicode, ok := symbolMap[value]; ok {
		return NewUnicodeSymbol(value, unicode, ctx)
	} else {
		panic("not a symbol") // use NewWord instead
	}
}

func newSymbol(value string, ctx context.Context) *Symbol {
	s, err := NewSymbol(value, ctx)
	if err != nil {
		panic(err)
	}

	return s
}

func newUnicodeSymbol(inputValue string, unicode int, ctx context.Context) *Symbol {
	s, err := NewUnicodeSymbol(inputValue, unicode, ctx)
	if err != nil {
		panic(err)
	}

	return s
}

func (t *Symbol) Value() string {
	return t.symbol
}

func (t *Symbol) Dump(indent string) string {
	return indent + "Symbol(" + t.value + ")\n"
}

// dont take advance width into account (rely on operator spacing)
func (t *Symbol) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bb := unicodeBB(t.unicode)

	if err := scope.BuildMathText(x-bb.Left(), y, 1.0, t.value, t.Context()); err != nil {
		return nil, err
	}

	bb = bb.Translate(x-bb.Left(), y)
	return bb, nil
}

func IsSymbol(t Token) bool {
	_, ok := t.(*Symbol)
	return ok
}
