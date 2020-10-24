package scripts

var (
	VERBOSITY = 0
)

type Script interface {
	Write() (string, error)
	Dependencies() []string // src fields in script or call
}
