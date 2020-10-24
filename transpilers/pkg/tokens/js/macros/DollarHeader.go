package macros

import ()

// cache the literal queries
type DollarHeader struct {
	HeaderData
}

func (h *DollarHeader) Dependencies() []Header {
	return []Header{}
}

func (h *DollarHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("class ", h.Name(), "{")
	b.tcn("constructor(e,s){")
	b.ttcn("this.element=e;")
	b.ttcn("for(var k of s){")
	b.tttcn("let l=k;")
	b.tttcn("this[l]=function(x){")
	b.ttttcn("if(x===undefined){return e.getAttribute(l)}")
	b.ttttcn("else if(x==null){e.removeAttribute(l)}")
	b.ttttcn("else{e.setAttribute(l,x)}")
	b.tttcn("}")
	b.ttcn("}")
	b.ttcn("let d='';")
	b.ttcn("this.hide=()=>{")
	b.tttcn("if(this.element.style.display!='none'){")
	b.ttttcn("d=this.element.style.display;")
	b.ttttcn("this.element.style.display='none';")
	b.tttcn("}")
	b.tttcn("return this")
	b.ttcn("}")
	b.ttcn("this.show=()=>{")
	b.tttcn("this.element.style.display=d;")
	b.tttcn("return this")
	b.ttcn("}")
	b.tcn("}")
	b.tcn("append(e){")
	b.ttcn("if(e instanceof Array){for(let c of e){this.append(c)}}")
	b.ttcccn("else if(e instanceof ", h.Name(), "){this.element.appendChild(e.element)}")
	b.ttcn("else{this.element.appendChild(e)}")
	b.ttcn("return this")
	b.tcn("}")
	b.tcn("on(s,f){")
	b.ttcn("this.element.addEventListener(s,f);return this")
	b.tcn("}")
	b.tcn("html(s){")
	b.ttcn("if(s===undefined){return this.element.innerHTML}")
	b.ttcn("else{this.element.innerHTML=s;return this}")
	b.tcn("}")
	b.tcn("clear(){")
	b.ttcn("while(this.element.firstChild){") // it might be faster to check firstChild
	b.tttcn("this.element.removeChild(this.element.lastChild)")
	b.ttcn("}")
	b.tcn("}")
	b.tcn("get value(){")
	b.ttcn("return this.element.value;")
	b.tcn("}")
	b.tcn("set value(x){")
	b.ttcn("this.element.value=x;")
	b.tcn("}")
	b.tcn("remove(){")
	b.ttcn("this.element.parentNode.removeChild(this.element);")
	b.tcn("}")
	b.tcn("height(h){")
	b.ttcn("if(h===undefined){return this.element.clientHeight}")
	b.ttcn("else{this.element.style.height=h+'px';return this}")
	b.tcn("}")
	b.tcn("left(l){")
	b.ttcn("if(l===undefined){return this.element.clientLeft}")
	b.ttcn("else{this.element.style.left=l+'px';return this}")
	b.tcn("}")
	b.tcn("top(t){")
	b.ttcn("if(t===undefined){return this.element.clientTop}")
	b.ttcn("else{this.element.style.top=t+'px';return this}")
	b.tcn("}")
	b.tcn("width(w){")
	b.ttcn("if(w===undefined){return this.element.clientWidth}")
	b.ttcn("else{this.element.style.width=w+'px';return this}")
	b.tcn("}")
	b.tcn("attr(k,v){")
	b.ttcn("if(v===undefined){return this.element.getAttribute(k)}")
	b.ttcn("else if(v==null){this.element.removeAttribute(k);return this}")
	b.ttcn("else{this.element.setAttribute(k,v);return this}")
	b.tcn("}")
	b.c("}")
	b.n()

	return b.String()
}

var dollarHeader = &DollarHeader{newHeaderData("$")}
