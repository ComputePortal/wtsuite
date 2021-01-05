package js

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
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

func (t *ForBlock) writeStatementFooter(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	if len(t.statements) == 0 {
		b.WriteString(";")
	} else {
		b.WriteString("){")
		b.WriteString(nl)

		b.WriteString(t.writeBlockStatements(usage, indent+tab, nl, tab))

		b.WriteString(nl)
		b.WriteString(indent)
		b.WriteString("}")
	}

	return b.String()
}
