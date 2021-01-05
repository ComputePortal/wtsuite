package patterns

import ()

const (
	COMPACT_LETTER = "abcdefghijklmnopqrstuvwxyz"
)

type NameGenerator struct {
	allowCompactNaming bool
	i                  int
	name               string
}

func NewCustomStartNameGenerator(allowCompactNaming bool, start int, name string) *NameGenerator {
	return &NameGenerator{allowCompactNaming, start, name}
}

func NewNameGenerator(allowCompactNaming bool, name string) *NameGenerator {
	return NewCustomStartNameGenerator(allowCompactNaming, 0, name)
}

func (ng *NameGenerator) GenName() string {
	if COMPACT_NAMING && ng.allowCompactNaming {

		n := ""

		for n == "" ||
			(len(ng.name) != 1 && len(n) == 1) ||
			n == "if" ||
			n == "of" ||
			n == "in" ||
			n == "do" ||
			n == "as" {
			i := ng.i

			if i == 0 {
				n = "A"
				if len(ng.name) == 1 {
					n = ng.name
				}
			} else {
				for i > 0 {
					rem := (i - 1) % 26
					n = n + COMPACT_LETTER[rem:rem+1]
					i = (i - 1) / 26
				}
			}

			ng.i++
		}

		return n
	} else {
		var n string
		if ng.i > 0 {
			ng.name += "_"

			n = ng.name
		} else {
			n = ng.name
		}

		ng.i++

		return n
	}
}
