package styles

import (
	"strconv"
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
