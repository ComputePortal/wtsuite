package macros

import (
  "../../context"
)

type CheckTypeHeader struct {
  HeaderData
}

func (h *CheckTypeHeader) Dependencies() []Header {
  return []Header{}
}

func (h *CheckTypeHeader) Write() string {
  b := NewHeaderBuilder()

  b.n()
  b.cccn("function ", h.Name(), "(x,t){")
  b.tcn("if(x===null){")
  b.ttcn("return x")
  b.tcn("}else if(x==undefined){")
  b.ttcn("throw new Error('undefined');")
  b.tcn("}else if(typeof x==='function'){")
  b.ttcn("throw new Error('didnt expect a function');")
  b.tcn("}else if(Array.isArray(t)){")
  b.ttcn("if(!Array.isArray(x)){")
  b.tttcn("throw new Error('expected array');")
  b.ttcn("}else if(t.length==1){") // regular array
  b.tttcn("for(let c of x){")
  b.ttttcn(h.Name() + "(c,t[0])")
  b.tttcn("}")
  b.ttcn("}else if(x.length!=t.length){") // tuple
  b.tttcn("throw new Error('expected tuple of length '+t.length.toString());")
  b.ttcn("}else{")
  b.tttcn("for(let i=0;i<x.length;i++){")
  b.ttttcn(h.Name()+"(x[i],t[i]);")
  b.tttcn("}")
  b.ttcn("}")
  b.tcn("}else if(t.__implementations__!=undefined){")
  b.ttcn("let ok=false;")
  b.ttcn("for(let p of t.__implementations__){")
  b.tttcccn("try{", h.Name(), "(x,p);ok=true}catch(e){}")
  b.ttcn("}")
  b.ttcn("if(!ok){throw new Error('expected '+t.prototype.constructor.name);}")
  b.tcn("}else if(t.prototype!=undefined){")
  b.ttcn("let xn='';")
  b.ttcn("if(x.constructor!=undefined){")
  b.tttcn("xn=x.constructor.name;")
  b.ttcn("}")
  b.ttcn("switch(t.prototype.constructor.name){")
  b.ttcn("case 'String':")
  b.tttcn("if(typeof(x)!='string'&&xn!='String'){")
  b.ttttcn("throw new Error('expected String');")
  b.tttcn("};break;")
  b.ttcn("case 'Number':")
  b.tttcn("if(typeof(x)!='number'&&xn!='Number'){")
  b.ttttcn("throw new Error('expected Number');")
  b.tttcn("};break;")
  b.ttcn("case 'Int':")
  b.tttcn("if(!Number.isInteger(x)&&xn!='Int'){")
  b.ttttcn("throw new Error('expected Int');")
  b.tttcn("};break;")
  b.ttcn("case 'Boolean':")
  b.tttcn("if(typeof(x)!='boolean'&&xn!='Boolean'){")
  b.ttttcn("throw new Error('expected Boolean');")
  b.tttcn("};break;")
  b.ttcn("default:")
  b.tttcn("if(x.constructor==undefined){")
  b.ttttcn("throw new Error('expected '+t.prototype.constructor.name+', got non-instance');")
  b.tttcn("}else if(x.constructor!=undefined&&t.prototype!==x.constructor.prototype){")

  b.ttttcn("throw new Error('expected '+t.prototype.constructor.name+', got '+x.constructor.name);")
  b.tttcn("}else if(t.__propertyTypes__!=undefined){")
  b.ttttcn("for(let k in t.__propertyTypes__){")
  b.tttttcn(h.Name() + "(x[k], t.__propertyTypes__[k]);")
  b.ttttcn("}")
  b.tttcn("}")
  b.ttcn("}")
  b.tcn("}else if(t.constructor!=undefined&&t.constructor.name=='Object'){")
  b.ttcn("for(let k in t){")  // anything that x contains as surplus is ignored
  b.tttcn("if(k==''){for(let m in x){" + h.Name() + "(x[m],t[k])};return;")
  b.tttcn("}else{")
  b.ttttcn("if(Object.keys(x).length!==Object.keys(t).length){throw new Error('different number of object properties');}")
  b.ttttcn(h.Name() + "(x[k],t[k])")
  b.tttcn("}")
  b.ttcn("}") 
  b.tcn("}else{")
  b.ttcn("throw new Error('bad type check');")
  b.tcn("}")
  b.tcn("return x;")
  b.c("}")
  b.n()

  return b.String()
}

var checkTypeHeader = &CheckTypeHeader{newHeaderData("__checkType__")}

func ActivateCheckTypeHeader() {
  ResolveHeaderActivity(checkTypeHeader, context.NewDummyContext())
}
