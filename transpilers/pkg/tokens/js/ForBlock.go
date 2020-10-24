package js

import (
	"strings"

	"../context"
)

// common for For, ForIn, ForOf
// all for statements must specify a new variable using const, let or var
type ForBlock struct {
	varType VarType
	Block
}

func newForBlock(varType VarType, ctx context.Context) ForBlock {
	return ForBlock{varType, newBlock(ctx)}
}

// extra string between 'for' and first '(' (eg. ' await')
func (t *ForBlock) writeStatementHeader(indent string, extra string,
	writeVarType bool) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("for")
	b.WriteString(extra)
	b.WriteString("(")

	if writeVarType {
		b.WriteString(VarTypeToString(t.varType))
		b.WriteString(" ")
	}

	return b.String()
}

func (t *ForBlock) writeStatementFooter(indent string) string {
	var b strings.Builder

	if len(t.statements) == 0 {
		b.WriteString(";")
	} else {
		b.WriteString("){")
		b.WriteString(NL)

		b.WriteString(t.writeBlockStatements(indent+TAB, NL))

		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString("}")
	}

	return b.String()
}
