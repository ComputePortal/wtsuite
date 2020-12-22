package math

import (
	"./boundingbox"

	"../context"
)

type Word struct {
	value string
	TokenData
}

// letters ARE NOT symbols
var letterMap = map[string]int{
	"(": 0x00028,
	")": 0x00028,
	"*": 0x0002a,
	"+": 0x0002b,
	"-": 0x0002d,
	".": 0x0002e,
	"/": 0x0002f,
	"0": 0x00030,
	"1": 0x00031,
	"2": 0x00032,
	"3": 0x00033,
	"4": 0x00034,
	"5": 0x00035,
	"6": 0x00036,
	"7": 0x00037,
	"8": 0x00038,
	"9": 0x00039,

	"A": 0x00041,
	"B": 0x00042,
	"C": 0x00043,
	"D": 0x00044,
	"E": 0x00045,
	"F": 0x00046,
	"G": 0x00047,
	"H": 0x00048,
	"I": 0x00049,
	"J": 0x0004a,
	"K": 0x0004b,
	"L": 0x0004c,
	"M": 0x0004d,
	"N": 0x0004e,
	"O": 0x0004f,
	"P": 0x00050,
	"Q": 0x00051,
	"R": 0x00052,
	"S": 0x00053,
	"T": 0x00054,
	"U": 0x00055,
	"V": 0x00056,
	"W": 0x00057,
	"X": 0x00058,
	"Y": 0x00059,
	"Z": 0x0005a,

	"a": 0x00061,
	"b": 0x00062,
	"c": 0x00063,
	"d": 0x00064,
	"e": 0x00065,
	"f": 0x00066,
	"g": 0x00067,
	"h": 0x00068,
	"i": 0x00069,
	"j": 0x0006a,
	"k": 0x0006b,
	"l": 0x0006c,
	"m": 0x0006d,
	"n": 0x0006e,
	"o": 0x0006f,
	"p": 0x00070,
	"q": 0x00071,
	"r": 0x00072,
	"s": 0x00073,
	"t": 0x00074,
	"u": 0x00075,
	"v": 0x00076,
	"w": 0x00077,
	"x": 0x00078,
	"y": 0x00079,
	"z": 0x0007a,
}

func newWord(value string, ctx context.Context) Token {
	return &Word{value, newTokenData(ctx)}
}

func NewValueWord(value string, ctx context.Context) (Token, error) {
	return newWord(value, ctx), nil
}

func NewWord(value string, ctx context.Context) (Token, error) {
	if unicode, ok := symbolMap[value]; ok {
		return NewUnicodeSymbol(value, unicode, ctx)
	} else {
		return NewValueWord(value, ctx)
	}
}

func (t *Word) Value() string {
	return t.value
}

func (t *Word) Dump(indent string) string {
	return indent + "Word(" + t.value + ")\n"
}

func (t *Word) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	bb := boundingbox.NewBB(x, y, x, y) // start at current start, dont exactly left-align like Symbol

	x_ := x
	// calculate the bb before generating the tag
	for _, c := range t.value {
		if unicode, ok := letterMap[string(c)]; ok {
			bbc := unicodeBB(unicode).Translate(x_, y)

			bb = boundingbox.Merge(bbc, bb)

			x_ += unicodeAdvanceWidth(unicode)
		} else {
			panic("character " + string(c) + " not yet handled")
		}
	}

	if err := scope.BuildMathText(x, y, 1.0, t.value, t.Context()); err != nil {
		return nil, err
	}

	return bb, nil
}

func IsWord(t Token) bool {
	_, ok := t.(*Word)
	return ok
}
