package tree

import (
	"strconv"
)

var (
	_uid_    = 0
	_uclass_ = 0
)

func NewUniqueID() string {
	res := "_" + strconv.Itoa(_uid_)
	_uid_++
	return res
}

func NewUniqueClass() string {
	res := "_" + strconv.Itoa(_uclass_)
	_uclass_++
	return res
}
