package macros

type BlobFromInstanceHeader struct {
	HeaderData
}

func (h *BlobFromInstanceHeader) Dependencies() []Header {
	return []Header{}
}

func (h *BlobFromInstanceHeader) writeHeaderCreationAndAllocation(b *HeaderBuilder) {
	b.tcn("function a(x,incr){")
	b.ttcn("if(Array.isArray(x)){")
	b.tttcn("var y=new Array(x.length);")
	b.tttcn("for(var i=0;i<x.length;i++){")
	b.ttttcn("y[i]=a(x[i],incr);")
	b.tttcn("}")
	b.tttcn("return y;")
	b.ttcn("}else if(x===undefined){")
	b.tttcn("throw new Error('cannot count undefined');")
	b.ttcn("}else if(x===null){")
	b.tttcn("return x;")
	b.ttcn("}else if(typeof x==='function'){")
	b.tttcn("throw new Error('cannot count function');")
	b.ttcn("}else if(typeof x!=='object'){")
	b.tttcn("return x;")
	b.ttcn("}else if(x.constructor.name==='Int8Array'){")
	b.tttcn("return {__type__:'Int8Array',start:incr(x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Int16Array'){")
	b.tttcn("return {__type__:'Int16Array',start:incr(2*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Int32Array'){")
	b.tttcn("return {__type__:'Int32Array',start:incr(4*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Uint8Array'){")
	b.tttcn("return {__type__:'Uint8Array',start:incr(x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Uint16Array'){")
	b.tttcn("return {__type__:'Uint16Array',start:incr(2*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Uint32Array'){")
	b.tttcn("return {__type__:'Uint32Array',start:incr(4*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Float32Array'){")
	b.tttcn("return {__type__:'Float32Array',start:incr(4*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Float64Array'){")
	b.tttcn("return {__type__:'Float64Array',start:incr(8*x.length),length:x.length};")
	b.ttcn("}else if(x.constructor.name==='Date'){")
	b.tttcn("return {__type__:'Date',time:x.getTime()};")
	b.ttcn("}else if(x.constructor.name==='Error'){")
	b.tttcn("return {__type__:'Error',message:x.message};")
	b.ttcn("}else{")
	b.tttcn("var y={};")
	b.tttcn("for(var k in x){")
	b.ttttcn("y[k]=a(x[k],incr);")
	b.tttcn("}")
	b.tttcn("y.__type__=x.constructor.name;")
	b.tttcn("return y;")
	b.ttcn("}")
	b.tcn("}")
}

func (h *BlobFromInstanceHeader) writeFilling(b *HeaderBuilder) {
	b.tcn("function f(x,y,views,o){")
	b.ttcn("if(Array.isArray(x)){")
	b.tttcn("for(var i=0;i<x.length;i++){")
	b.ttttcn("f(x[i],y[i],views,o);")
	b.tttcn("}")
	b.ttcn("}else if(typeof x!=='object'){")
	b.ttcn("}else if(x===null){")
	b.ttcn("}else if(x.constructor.name==='Int8Array'){")
	b.tttcn("views.int8.set(x, y.start+o);")
	b.ttcn("}else if(x.constructor.name==='Int16Array'){")
	b.tttcn("views.int16.set(x, (y.start+o)/2);")
	b.ttcn("}else if(x.constructor.name==='Int32Array'){")
	b.tttcn("views.int32.set(x, (y.start+o)/4);")
	b.ttcn("}else if(x.constructor.name==='Uint8Array'){")
	b.tttcn("views.uint8.set(x, y.start+o);")
	b.ttcn("}else if(x.constructor.name==='Uint16Array'){")
	b.tttcn("views.uint16.set(x, (y.start+o)/2);")
	b.ttcn("}else if(x.constructor.name==='Uint32Array'){")
	b.tttcn("views.uint32.set(x, (y.start+o)/4);")
	b.ttcn("}else if(x.constructor.name==='Float32Array'){")
	b.tttcn("views.float32.set(x, (y.start+o)/4);")
	b.ttcn("}else if(x.constructor.name==='Float64Array'){")
	b.tttcn("views.float64.set(x, (y.start+o)/8);")
	b.ttcn("}else if(x.constructor.name==='Date'){")
	b.ttcn("}else if(x.constructor.name==='Error'){")
	b.ttcn("}else{")
	b.tttcn("for(var k in x){")
	b.ttttcn("f(x[k],y[k],views,o);")
	b.tttcn("}")
	b.ttcn("}")
	b.tcn("}")
}

func (h *BlobFromInstanceHeader) Write() string {
	b := NewHeaderBuilder()

	// first count the size of the necessary ArrayBuffer (because resizing is not directly possible
	b.n()
	b.cccn("function ", h.Name(), "(x){")

	h.writeHeaderCreationAndAllocation(b)

	b.tcn("var c=0;")
	// round s to nearest 8 byte boundary
	b.tcn("var incr=function(l){var s=c+(8-c%8)%8;c=s+l;return s;};")
	b.tcn("var y=a(x,incr);")
	b.tcn("c=c+(8-c%8)%8;") // final c must also round to 8 byte boundary
	b.tcn("var j=JSON.stringify(y);")
	b.tcn("var o=4+j.length;")
	b.tcn("var q=o+(8-o%8)%8;") // round to 8 byte boundary
	b.tcn("var b=new ArrayBuffer(q+c);")
	// first four bytes contain Uint32Array offset of typed data
	b.tcn("var views={int8:new Int8Array(b),int16:new Int16Array,int32:new Int32Array,uint8:new Uint8Array(b),uint16:new Uint16Array(b),uint32:new Uint32Array(b),float32:new Float32Array(b),float64:new Float64Array(b)};")
	b.tcn("views.uint32[0]=o;")
	b.tcn("views.uint8.set(new TextEncoder('utf-8').encode(j), 4);")

	h.writeFilling(b)
	b.tcn("f(x,y,views,q);")

	// return the blob
	b.tcn("return new Blob([b]);") // TODO: set the MIME type?

	b.c("}")
	b.n()

	return b.String()
}

var blobFromInstanceHeader = &BlobFromInstanceHeader{newHeaderData("__blobFromInstance__")}
