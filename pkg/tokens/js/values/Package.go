package values

type Package interface {
  AddPrototype(proto Prototype)
  AddValue(name string, v Value)
}
