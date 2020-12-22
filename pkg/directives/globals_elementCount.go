package directives

import (
	"../functions"
	"../tokens/context"
	tokens "../tokens/html"
)

// actually not globals, as they are only available inside the template scope
const ELEMENT_COUNT = "__idx__"
const ELEMENT_COUNT_FOLDED = "__idxf__"

func setElementCount(scope Scope, node Node, ctx context.Context) {
  idx := node.getElementCount()
  idxVar := functions.Var{tokens.NewValueInt(idx, ctx), true, true, false, false, ctx}
  scope.SetVar(ELEMENT_COUNT, idxVar)

  idxf := node.getElementCountFolded()
  idxfVar := functions.Var{tokens.NewValueInt(idxf, ctx), true, true, false, false, ctx}
  scope.SetVar(ELEMENT_COUNT_FOLDED, idxfVar)
}
