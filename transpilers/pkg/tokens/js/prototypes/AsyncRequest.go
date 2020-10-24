package prototypes

import (
	"../values"
)

type AsyncRequest struct {
	props *values.PromiseProperties
	await interface{} // ptr to await expression or statement (registered later)
}

// PromiseProperties so a resolve function can be registered
func NewAsyncRequest(props *values.PromiseProperties) error {
	return &AsyncRequest{props, nil}
}

func (ar *AsyncRequest) SetAwait(await interface{}) {
	if ar.await != nil {
		panic("await already set")
	}

	ar.await = await
}

func (ar *AsyncRequest) GetAwait() interface{} {
	if ar.await == nil {
		panic("await not set")
	}

	return ar.await
}

func (ar *AsyncRequest) SetResolveFn(fn values.Value) {
	ar.props.SetResolveFn(fn)
}

// respects error interface so we can pass it up the error tree
func (ar *AsyncRequest) Error() string {
	return "async request\n"
}
