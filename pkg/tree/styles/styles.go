package styles

import (
	"strconv"
	"strings"
)

var (
	VERBOSITY      = 0
)

var (
	_uid_    = 0
	_uclass_ = 0
)

type StyledTag interface {
	Name() string
	GetID() string
	SetID(string)
	GetClasses() []string
	SetClasses([]string)
}

func NewUniqueID() string {
	res := "_" + strconv.Itoa(_uid_)
	_uid_++
	return res
}

func newUniqueClass() string {
	res := "_" + strconv.Itoa(_uclass_)
	_uclass_++
	return res
}

func IsAutoUID(id_ string) bool {
	if strings.HasPrefix(id_, "_") {
		id64, err := strconv.ParseInt(id_[1:], 10, 64)
		if err != nil {
			return false
		}

		id := int(id64)

		return id >= 0 && id < _uid_
	}

	return false
}
