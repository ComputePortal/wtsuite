package main

import (
  "strings"

	"../../pkg/tokens/context"
	"../../pkg/tokens/js"
	"../../pkg/tokens/js/values"
)

type Graph struct {
  radial bool
  files map[string]string // only nodes in these files are considered, edges aren't allowed to point outward if files isn't nil
  nodes []Node
  edges []Edge
}

// files can be nil in order to graph more greedily
func NewGraph(files map[string]string) *Graph {
  return &Graph{
    true,
    files,
    []Node{},
    []Edge{},
  }
}

func (g *Graph) AppendNode(node Node) {
  g.nodes = append(g.nodes, node)
}

func (g *Graph) AppendEdge(edge Edge) {
  g.edges = append(g.edges, edge)
}

func (g *Graph) Clean() {
  // remove nodes with the same name
  // the first ones get priority

  uniqueNodes := make(map[string]Node)
  for _, node := range g.nodes {
    name := node.Name()

    if _, ok := uniqueNodes[name]; !ok {
      uniqueNodes[name] = node
    }
  }

  g.nodes = make([]Node, 0)

  // order is lost
  // XXX: should we preserve order?
  for _, node := range uniqueNodes {
    g.nodes = append(g.nodes, node)
  }
}

func (g *Graph) Write() string {
  var b strings.Builder

  b.WriteString("digraph {")
  b.WriteString("\n")

  if g.radial {
    b.WriteString("  ranksep=3;\nratio=auto;\n")
  }

  for _, node := range g.nodes {
    b.WriteString(node.Write("  "))
    b.WriteString("\n")
  }

  b.WriteString("\n")

  for _, edge := range g.edges {
    b.WriteString(edge.Write("  "))
    b.WriteString("\n")
  }

  b.WriteString("}\n")

  return b.String()
}

func (g *Graph) IsIncluded(ctx context.Context) bool {
  if g.files == nil {
    return true
  }

  path := ctx.Path()

  if _, ok := g.files[path]; ok {
    return true
  } else {
    return false
  }
}

func (g *Graph) addClassParentAndInterface(class *js.Class) error {
  parent_ := class.GetParent()

  if parent_ != nil {
    if parent, ok := parent_.(*js.Class); ok {
      if g.IsIncluded(parent.Context()) {
        g.AppendEdge(NewInheritsEdge(class, parent))

        g.AppendNode(NewClassNode(parent))
      }
    }
  }

  interface_ := class.GetInterface()

  if interface_ != nil {
    if interf, ok := interface_.(*js.ClassInterface); ok {
      if g.IsIncluded(interf.Context()) {
        // TODO: special interface node for internal nodes?
        g.AppendNode(NewInterfaceNode(interf))

        g.AppendEdge(NewImplementsEdge(class, interf))
      }
    }
  }

  return nil
}

func (g *Graph) HasClass(name string) bool {
  for _, node_ := range g.nodes {
    if node, ok := node_.(*ClassNode); ok {
      if node.Name() == name {
        return true
      }
    }
  }

  return false
}

// TODO: some options?
func (g *Graph) AddClass(class *js.Class) error {
  if !g.HasClass(class.Name()) {
    g.AppendNode(NewClassNode(class))

    return g.addClassParentAndInterface(class)
  }

  return nil
}

func (g *Graph) HasInstance(name string) bool {
  for _, node_ := range g.nodes {
    if node, ok := node_.(*InstanceNode); ok {
      if node.Name() == name {
        return true
      }
    }
  }

  return false
}

// circular dependencies are possible
func (g *Graph) AddInstance(instance *values.Instance) error {
  // use the class name (i.e. the type name)
  if !g.HasInstance(instance.TypeName()) {
    g.AppendNode(NewInstanceNode(instance))

    // loop the properties
    props := instance.Properties().InstanceProperties()

    for label, prop := range props {
      propProto, _ := prop.GetInstancePrototype()
      // only user defined classes are used
      if propClass, ok := propProto.(*js.Class); ok {

        if g.IsIncluded(propClass.Context()) {
          g.AppendEdge(NewPropertyEdge(instance, label, prop))

          if err := g.AddInstance(prop); err != nil {
            return err
          }
        }
      }
    }
  }
  
  return nil
}

func (g *Graph) AddClassInstance(class *js.Class, instance *values.Instance) error {
  if err := g.AddInstance(instance); err != nil {
    return err
  }

  if err := g.addClassParentAndInterface(class); err != nil {
    return err
  }

  return nil
}
