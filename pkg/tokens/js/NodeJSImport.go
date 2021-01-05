package js

import (
  "github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

var nodeJSPackages map[string]FillPackageFunction = map[string]FillPackageFunction{
  "crypto": prototypes.FillNodeJS_cryptoPackage,
  "fs": prototypes.FillNodeJS_fsPackage,
  "http": prototypes.FillNodeJS_httpPackage,
  "mysql": prototypes.FillNodeJS_mysqlPackage,
  "nodemailer": prototypes.FillNodeJS_nodemailerPackage,
  "path": prototypes.FillNodeJS_pathPackage,
  "process": prototypes.FillNodeJS_processPackage,
  "stream": prototypes.FillNodeJS_streamPackage,
}

func IsNodeJSPackage(name string) bool {
  _, ok := nodeJSPackages[name]
  return ok
}

type NodeJSImport struct {
  name   string
	expr   *VarExpression
	TokenData
}

func NewNodeJSImport(name string, expr *VarExpression, ctx context.Context) *NodeJSImport {
	return &NodeJSImport{
    name,
		expr,
		newTokenData(ctx),
	}
}

func (m *NodeJSImport) Dump(indent string) string {
  return indent + "NodeJSImport(" + m.expr.Name() + ")\n"
}

func (m *NodeJSImport) AddStatement(st Statement) {
	panic("not a block statement")
}

func (m *NodeJSImport) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	return "const " + m.expr.Name() + "=require('" + m.name + "')"
}

func (m *NodeJSImport) HoistNames(scope Scope) error {
	return nil
}

func (m *NodeJSImport) ResolveStatementNames(scope Scope) error {
  pkg := NewBuiltinPackage(m.expr.Name())

  if pkgFiller, ok := nodeJSPackages[m.name]; ok {
    pkgFiller(pkg)
  } else {
    panic("should've been caught before")
  }

  m.expr.variable = pkg
	if err := scope.SetVariable(m.expr.Name(), pkg); err != nil {
		return err
	}

	return nil
}

func (m *NodeJSImport) EvalStatement() error {
	return nil
}

func (m *NodeJSImport) ResolveStatementActivity(usage Usage) error {
	return m.expr.ResolveExpressionActivity(usage)
}

func (m *NodeJSImport) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (m *NodeJSImport) UniqueStatementNames(ns Namespace) error {
	return m.expr.uniqueDeclarationName(ns, CONST)
}

func (m *NodeJSImport) Walk(fn WalkFunc) error {
  if err := m.expr.Walk(fn); err != nil {
    return err
  }

  return fn(m)
}
