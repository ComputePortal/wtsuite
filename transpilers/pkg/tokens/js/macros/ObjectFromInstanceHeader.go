package macros

type ObjectFromInstanceHeader struct {
	HeaderData
}

func (h *ObjectFromInstanceHeader) Dependencies() []Header {
	return []Header{}
}

func (h *ObjectFromInstanceHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("function ", h.Name(), "(x){")
	b.tcn("if(Array.isArray(x)){")
	b.ttcn("var y=new Array(x.length);")
	b.ttcn("for(var i=0;i<x.length;i++){")
	b.tttcccn("y[i]=", h.Name(), "(x[i]);")
	b.ttcn("}")
	b.ttcn("return y;")
	b.tcn("}else if(x===undefined){")
	b.ttcn("throw new Error('cannot add proto name to undefined');")
	b.tcn("}else if(typeof x==='function'){")
	b.ttcn("throw new Error('cannot add proto name to function');")
	b.tcn("}else if(typeof x!=='object'){")
	b.ttcn("return x")
	b.tcn("}else if(x===null){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Int8Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Int16Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Int32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Uint8Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Uint16Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Uint32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Float32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Float64Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name==='Date'){")
	b.ttcn("var y={")
	b.tttcn("__type__:'Date',")
	b.tttcn("time:x.getTime()")
	b.ttcn("};")
	b.ttcn("return y;")
	b.tcn("}else if(x.constructor.name==='Error'){")
	b.ttcn("var y={")
	b.tttcn("__type__:'Error',")
	b.tttcn("message:x.message")
	b.ttcn("};")
	b.ttcn("return y;")
	b.tcn("}else{")
	b.ttcn("var y={};")
	b.ttcn("for(var k in x){")
	b.tttcccn("y[k]=", h.Name(), "(x[k]);")
	b.ttcn("}")
	b.ttcn("y.__type__=x.constructor.name;")
	b.ttcn("return y;")
	b.tcn("}")
	b.c("}")
	b.n()

	return b.String()
}

var objectFromInstanceHeader = &ObjectFromInstanceHeader{newHeaderData("__objectFromInstance__")}
