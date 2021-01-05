package parsers

import (
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

func splitByNextSeparator(ts []raw.Token, sep string) ([]raw.Token,
	[]raw.Token) {
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			if i < len(ts)-1 {
				return ts[0:i], ts[i+1:]
			} else {
				return ts[0:i], []raw.Token{}
			}
		}
	}

	return ts, []raw.Token{}
}

func splitBySeparator(ts []raw.Token, sep string) [][]raw.Token {
	result := make([][]raw.Token, 0)

	prev := 0
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			result = append(result, ts[prev:i])
			prev = i + 1
		}
	}

	if prev < len(ts) {
		result = append(result, ts[prev:])
	}

	return result
}

func stripSeparators(iStart int, ts []raw.Token, symbol string) []raw.Token {
	if iStart > len(ts)-1 {
		return []raw.Token{}
	}

	iRemaining := iStart
	for i := iStart; i < len(ts); i++ {
		if !raw.IsSymbol(ts[i], symbol) {
			iRemaining = i
			break
		}
	}

	return ts[iRemaining:]
}

func nextSeparatorPosition(ts []raw.Token, sep string) int {
	for i, t := range ts {
		if raw.IsSymbol(t, sep) {
			return i
		}
	}

	return len(ts)
}

func nextSymbolPositionThatEndsWith(ts []raw.Token, sep string) int {
	for i, t := range ts {
		if raw.IsSymbolThatEndsWith(t, sep) {
			return i
		}
	}

	return len(ts)
}

// turn [Word(aaa), Symbol(.), Word(bbb), ....] into Word(aaa.bbb...)
func condensePackagePeriods(ts []raw.Token) (*raw.Word, []raw.Token, error) {
	nameToken, err := raw.AssertWord(ts[0])
	if err != nil {
		return nil, nil, err
	}

	i := 1
	for i < len(ts) {
		if raw.IsSymbol(ts[i], patterns.PERIOD) {
			if i == len(ts)-1 {
				errCtx := ts[i].Context()
				return nil, nil, errCtx.NewError("Error: expected tokens after .")
			}

			nextWord, err := raw.AssertWord(ts[i+1])
			if err != nil {
				return nil, nil, err
			}

			nameToken = raw.NewValueWord(nameToken.Value()+"."+nextWord.Value(), raw.MergeContexts(nameToken, ts[i], ts[i+1]))
			i += 2
		} else {
			break
		}
	}

	return nameToken, ts[i:], nil
}
