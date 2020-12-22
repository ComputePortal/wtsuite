package macros

type XMLPostHeader struct {
	HeaderData
}

func (h *XMLPostHeader) Dependencies() []Header {
	return []Header{objectToInstanceHeader, objectFromInstanceHeader}
}

func (h *XMLPostHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("function ", h.Name(), "(a, x){")
	b.tcn("return new Promise((resolve, reject)=>{")
	b.ttcn("var h=new XMLHttpRequest();")
	b.ttcn("h.open('POST',a);")
	b.ttcn("h.setRequestHeader('Content-Type','application/json');")
	b.ttcn("h.onerror=function(){")
	b.tttcn("reject(new Error('POST to '+a+' failed'));")
	b.ttcn("};")
	b.ttcn("h.onload=function(){")
	b.tttcn("if(h.status==200){")
	b.ttttcccn("let y=", objectToInstanceHeader.Name(), "(JSON.parse(h.responseText));")
	b.ttttcn("if(y.constructor.name=='Error'){reject(y);}else{resolve(y);}")
	b.tttcn("}else{")
	b.ttttcn("reject(new Error('POST to '+a+' failed'));")
	b.tttcn("}")
	b.ttcn("};")

	b.ttcccn("h.send(JSON.stringify(", objectFromInstanceHeader.Name(), "(x)));")
	b.tcn("});")
	b.c("}")

	return b.String()
}

var xmlPostHeader = &XMLPostHeader{newHeaderData("__xmlpost__")}
