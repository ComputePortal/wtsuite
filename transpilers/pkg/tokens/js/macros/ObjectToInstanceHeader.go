package macros

type ObjectToInstanceHeader struct {
	HeaderData
}

func (h *ObjectToInstanceHeader) Dependencies() []Header {
	return []Header{}
}

// TODO: create untrusted variant for server-side use
func (h *ObjectToInstanceHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("function ", h.Name(), "(x,p){")
	b.tcn("if(x===null){")
	b.ttcn("return x")
	b.tcn("}else if(x===undefined){")
	b.ttcn("throw new Error('cannot apply proto to undefined');")
	b.tcn("}else if(typeof x==='function'){")
	b.ttcn("throw new Error('didnt expect a function');")
	b.tcn("}else if(typeof x!=='object'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Int8Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Int16Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Int32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Uint8Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Uint16Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Uint32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Float32Array'){")
	b.ttcn("return x")
	b.tcn("}else if(x.constructor.name=='Float64Array'){")
	b.ttcn("return x")
	b.tcn("}else if(Array.isArray(x)){")
	b.ttcn("var y=new Array(x.length);")
	b.ttcn("for(var i=0;i<x.length;i++){")
	b.tttcccn("y[i]=", h.Name(), "(x[i]);")
	b.ttcn("}")
	b.ttcn("return y;")
	b.tcn("}else if(x.__type__=='Date'){")
	b.ttcn("return new Date(x.time);")
	b.tcn("}else if(x.__type__=='Error'){")
	b.ttcn("return new Error(x.message);")
	b.tcn("}else{")
	b.ttcn("var y={};")
	b.ttcn("for(var k in x){")
	b.tttcn("if(k==='__type__'){continue;}")
	b.tttcccn("y[k]=", h.Name(), "(x[k]);")
	b.ttcn("}")
	b.ttcn("var t=x.__type__;")
	b.ttcn("if(t==undefined){")
	b.tttcn("if(p!==undefined){")
	b.ttttcn("return Object.assign(Object.create(p.prototype),y);")
	b.tttcn("}else{")
	b.ttttcn("throw new Error('object __type__ not set');")
	b.tttcn("}")
	b.ttcn("}else if(!/[A-Z][a-zA-Z]*/.test(t)){")
	b.tttcn("throw new Error('invalid __type__ '+t);")
	b.ttcn("}else if (t=='Object'){")
	b.tttcn("return y;")
	b.ttcn("}else{")
	b.tttcn("var protoClass=eval(t);")
	b.tttcn("if(protoClass === undefined){")
	b.ttttcn("return null;")
	b.tttcn("}else if(protoClass.toInstance!==undefined){")
	b.ttttcn("return protoClass.toInstance(y);")
	b.tttcn("}else{")
	b.ttttcn("return Object.assign(Object.create(protoClass.prototype),y);")
	b.tttcn("}")
	b.ttcn("}")
	b.tcn("}")
	b.c("}")

	return b.String()
}

var objectToInstanceHeader = &ObjectToInstanceHeader{newHeaderData("__objectToInstance__")}
