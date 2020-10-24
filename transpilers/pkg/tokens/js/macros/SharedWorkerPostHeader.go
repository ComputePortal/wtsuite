package macros

type SharedWorkerPostHeader struct {
	HeaderData
}

func (h *SharedWorkerPostHeader) Dependencies() []Header {
	return []Header{objectFromInstanceHeader, objectToInstanceHeader}
}

func (h *SharedWorkerPostHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("function ", h.Name(), "(w, x, m){")
	b.tcn("return new Promise((resolve, reject)=>{")
	b.ttcn("if(w.__pending__==undefined){")
	b.tttcn("w.__pending__=[];")
	b.tttcn("w.port.onmessage=(e)=>{")
	b.ttttcn("let pending=w.__pending__.pop();")
	b.ttttcn("pending(e.data);")
	b.tttcn("}")
	b.ttcn("}")

	b.ttcn("w.__pending__.push(function(data){")
	b.tttcccn("let y=", objectToInstanceHeader.Name(), "(data, m);")
	b.tttcn("if(y.constructor.name=='Error'){")
	b.ttttcn("reject(y);")
	b.tttcn("}else{")
	b.ttttcn("resolve(y);")
	b.tttcn("}")
	b.ttcn("});")

	b.ttcccn("w.port.postMessage(", objectFromInstanceHeader.Name(), "(x);")
	b.tcn("}")
	b.c("}")

	return b.String()
}

var sharedWorkerPostHeader = &SharedWorkerPostHeader{newHeaderData("__swpost__")}
