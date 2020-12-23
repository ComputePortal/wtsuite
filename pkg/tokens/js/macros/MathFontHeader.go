package macros

import (
  "fmt"
  "strconv"

  "github.com/computeportal/wtsuite/pkg/tokens/math/serif"
)

type MathFontHeader struct {
  HeaderData
}

func (h *MathFontHeader) Dependencies() []Header {
  return []Header{}
}

func (h *MathFontHeader) writeSymbolToCodePointEntries(b *HeaderBuilder) {
  // italic
  b.tcn("Ait:0x1d434,")
  b.tcn("Bit:0x1d435,")
  b.tcn("Cit:0x1d436,")
  b.tcn("Dit:0x1d437,")
  b.tcn("Eit:0x1d438,")
  b.tcn("Fit:0x1d439,")
  b.tcn("Git:0x1d43a,")
  b.tcn("Hit:0x1d43b,")
  b.tcn("Iit:0x1d43c,")
  b.tcn("Jit:0x1d43d,")
  b.tcn("Kit:0x1d43e,")
  b.tcn("Lit:0x1d43f,")
  b.tcn("Mit:0x1d440,")
  b.tcn("Nit:0x1d441,")
  b.tcn("Oit:0x1d442,")
  b.tcn("Pit:0x1d443,")
  b.tcn("Qit:0x1d444,")
  b.tcn("Rit:0x1d445,")
  b.tcn("Sit:0x1d446,")
  b.tcn("Tit:0x1d447,")
  b.tcn("Uit:0x1d448,")
  b.tcn("Vit:0x1d449,")
  b.tcn("Wit:0x1d44a,")
  b.tcn("Xit:0x1d44b,")
  b.tcn("Yit:0x1d44c,")
  b.tcn("Zit:0x1d44d,")
  b.tcn("ait:0x1d44e,")
  b.tcn("bit:0x1d44f,")
  b.tcn("cit:0x1d450,")
  b.tcn("dit:0x1d451,")
  b.tcn("eit:0x1d452,")
  b.tcn("fit:0x1d453,")
  b.tcn("git:0x1d454,")
  b.tcn("hit:0x210e,")
  b.tcn("iit:0x1d456,")
  b.tcn("jit:0x1d457,")
  b.tcn("kit:0x1d458,")
  b.tcn("lit:0x1d459,")
  b.tcn("mit:0x1d45a,")
  b.tcn("nit:0x1d45b,")
  b.tcn("oit:0x1d45c,")
  b.tcn("pit:0x1d45d,")
  b.tcn("qit:0x1d45e,")
  b.tcn("rit:0x1d45f,")
  b.tcn("sit:0x1d460,")
  b.tcn("tit:0x1d461,")
  b.tcn("uit:0x1d462,")
  b.tcn("vit:0x1d463,")
  b.tcn("wit:0x1d464,")
  b.tcn("xit:0x1d465,")
  b.tcn("yit:0x1d466,")
  b.tcn("zit:0x1d467,")

  b.tcn("Gamma:0x393,")
  b.tcn("Delta:0x394,")
  b.tcn("Theta:0x398,")
  b.tcn("Lambda:0x39b,")
  b.tcn("Xi:0x39e,")
  b.tcn("Pi:0x3a0,")
  b.tcn("Sigma:0x3a3,")
  b.tcn("Upsilon:0x3a5,")
  b.tcn("Phi:0x3a6,")
  b.tcn("Psi:0x3a8,")
  b.tcn("Omega:0x3a9,")
  b.tcn("alpha:0x1d6fc,")
  b.tcn("beta:0x1d6fd,")
  b.tcn("gamma:0x1d6fe,")
  b.tcn("delta:0x1d6ff,")
  b.tcn("epsilon:0x1d700,")
  b.tcn("zeta:0x1d701,")
  b.tcn("eta:0x1d702,")
  b.tcn("theta:0x1d703,")
  b.tcn("iota:0x1d704,")
  b.tcn("kappa:0x1d705,")
  b.tcn("lambda:0x1d706,")
  b.tcn("mu:0x1d707,")
  b.tcn("nu:0x1d708,")
  b.tcn("xi:0x1d709,")
  b.tcn("pi:0x1d70b,")
  b.tcn("rho:0x1d70c,")
  b.tcn("sigma:0x1d70e,")
  b.tcn("tau:0x1d70f,")
  b.tcn("upsilon:0x1d710,")
  b.tcn("phi:0x1d711,")
  b.tcn("chi:0x1d712,")
  b.tcn("psi:0x1d713,")
  b.tcn("omega:0x1d714,")

  b.tcn("\"#\":0x23,")
  b.tcn("\",\":0x2c,")
  b.tcn("\"0\":0x30,")
  b.tcn("\"1\":0x31,")
  b.tcn("\"2\":0x32,")
  b.tcn("\"3\":0x33,")
  b.tcn("\"4\":0x34,")
  b.tcn("\"5\":0x35,")
  b.tcn("\"6\":0x36,")
  b.tcn("\"7\":0x37,")
  b.tcn("\"8\":0x38,")
  b.tcn("\"9\":0x39,")
  b.tcn("\";\":0x3b,")
  b.tcn("A:0x41,")
  b.tcn("B:0x42,")
  b.tcn("C:0x43,")
  b.tcn("D:0x44,")
  b.tcn("E:0x45,")
  b.tcn("F:0x46,")
  b.tcn("G:0x47,")
  b.tcn("H:0x48,")
  b.tcn("I:0x49,")
  b.tcn("J:0x4a,")
  b.tcn("K:0x4b,")
  b.tcn("L:0x4c,")
  b.tcn("M:0x4d,")
  b.tcn("N:0x4e,")
  b.tcn("O:0x4f,")
  b.tcn("P:0x50,")
  b.tcn("Q:0x51,")
  b.tcn("R:0x52,")
  b.tcn("S:0x53,")
  b.tcn("T:0x54,")
  b.tcn("U:0x55,")
  b.tcn("V:0x56,")
  b.tcn("W:0x57,")
  b.tcn("X:0x58,")
  b.tcn("Y:0x59,")
  b.tcn("Z:0x5a,")
  b.tcn("a:0x61,")
  b.tcn("b:0x62,")
  b.tcn("c:0x63,")
  b.tcn("d:0x64,")
  b.tcn("e:0x65,")
  b.tcn("f:0x66,")
  b.tcn("g:0x67,")
  b.tcn("h:0x68,")
  b.tcn("i:0x69,")
  b.tcn("j:0x6a,")
  b.tcn("k:0x6b,")
  b.tcn("l:0x6c,")
  b.tcn("m:0x6d,")
  b.tcn("n:0x6e,")
  b.tcn("o:0x6f,")
  b.tcn("p:0x70,")
  b.tcn("q:0x71,")
  b.tcn("r:0x72,")
  b.tcn("s:0x73,")
  b.tcn("t:0x74,")
  b.tcn("u:0x75,")
  b.tcn("v:0x76,")
  b.tcn("w:0x77,")
  b.tcn("x:0x78,")
  b.tcn("y:0x79,")
  b.tcn("z:0x7a,")

  // special symbols
  b.tcn("infty:0x221e,")
}

func (h *MathFontHeader) writeAdvanceWidthEntries(b *HeaderBuilder) {
  for i, aw := range serif.AdvanceWidths {
    b.t()
    b.c(strconv.Itoa(i))
    b.c(":")
    b.c(strconv.Itoa(aw))
    b.c(",")
    b.n()
  }
}

func (h *MathFontHeader) writeBoundingBoxEntries(b *HeaderBuilder) {
  for i, bb := range serif.Bounds {
    b.t()
    b.c(strconv.Itoa(i))
    b.c(":[")
    b.c(fmt.Sprintf("%g,%g,%g,%g", bb.Left(), bb.Right(), bb.Top(), bb.Bottom()))
    b.c("],")
    b.n()
  }
}

// TODO: a complete functional parser should be here
func (h *MathFontHeader) Write() string {
  b := NewHeaderBuilder()

  b.n()

  b.cccn("var ", h.Name(), "=(function(){")

  b.tcn("let s2cp={")
  h.writeSymbolToCodePointEntries(b)
  b.tcn("};")

  b.tcn("let aw={")
  h.writeAdvanceWidthEntries(b)
  b.tcn("};")

  b.tcn("let bb={")
  h.writeBoundingBoxEntries(b)
  b.tcn("};")

  b.tcn("return {")
  b.ttcn("symbolToCodePoint:function(s){let r=s2cp[s];if(r==undefined){return null}else{return r}},")
  b.ttcn("advanceWidth:function(i){let a=aw[i];if(a==undefined){return null}else{return a}},")
  b.ttcn("boundingBox:function(i){let b=bb[i];if(b==undefined){return null}else{return b}},")
  b.tcn("};")
  b.c("})();")
  b.n()

  return b.String()
}

var mathFontHeader = &MathFontHeader{newHeaderData("__mathFont__")}
