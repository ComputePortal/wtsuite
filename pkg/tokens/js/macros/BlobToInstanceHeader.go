package macros

type BlobToInstanceHeader struct {
	HeaderData
}

func (h *BlobToInstanceHeader) Dependencies() []Header {
	return []Header{}
}

// TODO: create untrusted variant, for server side use
func (h *BlobToInstanceHeader) writeUnpacker(b *HeaderBuilder) {
	b.tcn("function u(y,views,o,p){")
	b.ttcn("if(y===null){")
	b.tttcn("return null;")
	b.ttcn("}else if(y===undefined){")
	b.tttcn("throw new Error('cannot apply proto to null');")
	b.ttcn("}else if(typeof y==='function'){")
	b.tttcn("throw new Error('didnt expect function');")
	b.ttcn("}else if(typeof y!=='object'){")
	b.tttcn("return y;")
	b.ttcn("}else if(Array.isArray(y)){")
	b.tttcn("var x=new Array(y.length);")
	b.tttcn("for(var i=0;i<x.length;i++){")
	b.tttcn("x[i]=u(y[i],views,o);")
	b.tttcn("}")
	b.tttcn("return x;")
	b.ttcn("}else if(y.__type__===undefined){")
	b.tttcn("if(p===undefined){")
	b.ttttcn("throw new Error('object __type__ not set');")
	b.tttcn("}else{")
	b.ttttcn("return Object.assign(Object.create(p.prototype),y);")
	b.tttcn("}")
	b.ttcn("}else if(y.__type__==='Uint8Array'){")
	b.tttcn("return views.uint8.slice(y.start+o,y.start+o+y.length);")
	b.ttcn("}else if(y.__type__==='Uint16Array'){")
	b.tttcn("return views.uint16.slice((y.start+o)/2,(y.start+o)/2+y.length);")
	b.ttcn("}else if(y.__type__==='Uint32Array'){")
	b.tttcn("return views.uint32.slice((y.start+o)/4,(y.start+o)/4+y.length);")
	b.ttcn("}else if(y.__type__==='Int8Array'){")
	b.tttcn("return views.int8.slice((y.start+o),(y.start+o)+y.length);")
	b.ttcn("}else if(y.__type__==='Int16Array'){")
	b.tttcn("return views.int16.slice((y.start+o)/2,(y.start+o)/2+y.length);")
	b.ttcn("}else if(y.__type__==='Int32Array'){")
	b.tttcn("return views.int32.slice((y.start+o)/4,(y.start+o)/4+y.length);")
	b.ttcn("}else if(y.__type__==='Float32Array'){")
	b.tttcn("return views.float32.slice((y.start+o)/4,(y.start+o)/4+y.length);")
	b.ttcn("}else if(y.__type__==='Float64Array'){")
	b.tttcn("return views.float64.slice((y.start+o)/8,(y.start+o)/8+y.length);")
	b.ttcn("}else if(y.__type__==='Date'){")
	b.tttcn("return new Date(y.time);")
	b.ttcn("}else if(y.__type__==='Error'){")
	b.tttcn("return new Error(y.message);")
	b.ttcn("}else if(!/[A-Z][a-zA-Z]*/.test(y.__type__)){")
	b.tttcn("throw new Error('invalid __type__ '+y.__type__);")
	b.ttcn("}else{")
	b.tttcn("var x={};")
	b.tttcn("for(var k in y){")
	b.ttttcn("if(k==='__type__'){continue;}")
	b.ttttcn("x[k]=u(y[k],views,o);")
	b.tttcn("}")
	b.tttcn("var protoClass=eval(y.__type__);")
	b.tttcn("if(protoClass===undefined){")
	b.ttttcn("return null;")
	b.tttcn("}else if(protoClass.toInstance!==undefined){")
	b.ttttcn("return protoClass.toInstance(x);")
	b.tttcn("}else{")
	b.ttttcn("return Object.assign(Object.create(protoClass.prototype),x);")
	b.tttcn("}")
	b.ttcn("}")
	b.tcn("}")
}

func (h *BlobToInstanceHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	// returns a promise actually
	b.cccn("function ", h.Name(), "(bl,p){")
	h.writeUnpacker(b)

	b.tcn("return new Promise((resolve,reject)=>{")
	b.tcn("let reader=new FileReader();")
	b.tcn("reader.addEventListener('loadend',()=>{")
	b.ttcn("let b=reader.result;")
	b.ttcn("var views={uint8:new Uint8Array(b),uint16:new Uint16Array(b),uint32:new Uint32Array(b),int8:new Int8Array(b),int16:new Int16Array(b),int32:new Int32Array(b),float32:new Float32Array(b),float64:new Float64Array(b)};")
	b.ttcn("var o=views.uint32[0];")
	b.ttcn("var y=JSON.parse(new TextDecoder('utf-8').decode(views.uint8.slice(4, o)));")
	b.ttcn("var q=o+(8-o%8)%8;") // round to 8 byte boundary
	b.ttcn("resolve(u(y,views,q,p));")
	b.tcn("});")
	b.tcn("reader.readAsArrayBuffer(bl);")
	b.tcn("});")
	b.c("}")
	b.n()

	return b.String()
}

var blobToInstanceHeader = &BlobToInstanceHeader{newHeaderData("__blobToInstance__")}
