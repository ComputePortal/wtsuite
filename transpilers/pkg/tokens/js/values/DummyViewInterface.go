package values

import (
  "../../context"
)

// used during graphing, to avoid eval errors
type DummyViewInterface struct {
}

func NewDummyViewInterface() ViewInterface {
  return &DummyViewInterface{}
}

func (vi *DummyViewInterface) GetVarTypeInstance(stack Stack, name string, ctx context.Context) (Value, error) {
  return NewAllNull(ctx), nil
}

func (vi *DummyViewInterface) GetElemTypeInstance(stack Stack, name string, ctx context.Context) (Value, error) {
  return NewAllNull(ctx), nil
}

func (vi *DummyViewInterface) GetDefTypeInstance(stack Stack, name string, ctx context.Context) (Value, error) {
  return NewAllNull(ctx), nil
}

func (vi *DummyViewInterface) GetURL() string {
  return ""
}

func (vi *DummyViewInterface) GetHTML(stack Stack, id string, ctx context.Context) (string, error) {
  return "", nil
}

func (vi *DummyViewInterface) GetElemStates(id string) map[string][]string {
  return make(map[string][]string)
}

func (vi *DummyViewInterface) IsElem(key string) bool {
  // we dont know during graphing, so just return true to avoid an error
  return true
}
