package directives

import ()

func HasGlobal(k string) bool {
	if _, ok := _defines[k]; ok {
		return true
	}

	return k == URL || k == FILE || k == ELEMENT_COUNT
}
