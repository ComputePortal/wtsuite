package macros

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SearchIndexHeader struct {
	HeaderData
}

func (h *SearchIndexHeader) Dependencies() []Header {
	return []Header{}
}

func (h *SearchIndexHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()
	b.cccn("class ", h.Name(), "{")
	b.tcn("constructor(n,opt){")

	// fetch the index and parse the json
	b.ttcn("this._onready=[];") // pending callbacks
	b.ttcn("var data={pages:[null]};fetch(n).then((r)=>{return r.json()}).catch((e)=>{console.log('unable to fetch ' + n)}).then((d)=>{data=d;this._onready.forEach((fn)=>{fn()});this._onready=[];});")

	b.ttcn("this.page=function(i){")
	b.tttcn("return data.pages[i];")
	b.ttcn("};")

	b.ttcn("this.ignore=function(w){")
	b.tttcn("return data.ignore[w]!==undefined;")
	b.ttcn("};")

	b.ttcn("var _collect=function(m,s){")
	b.tttcn("for(let k in m){")
	b.ttttcn("if(k=='pages'){")
	b.tttttcn("m.pages.forEach((i)=>{s.add(i)});")
	b.ttttcn("}else{")
	b.tttttcn("_collect(m[k],s)")
	b.ttttcn("}")
	b.tttcn("}")
	b.ttcn("};")

	// recursive private function
	b.ttcn("var _search=function(m,w,sw,s){")
	b.tttcn("if(w.length==0){")
	b.ttttcn("if(sw){")
	b.tttttcn("_collect(m,s)")
	b.ttttcn("}else if(m.pages!==undefined){")
	b.tttttcn("m.pages.forEach((i)=>{s.add(i)});")
	b.tttttcn("return")
	b.ttttcn("}")
	b.tttcn("}")
	b.tttcn("let c=w.slice(0,1);")
	b.tttcn("let m_=m[c];")
	b.tttcn("if(m_!==undefined){")
	b.ttttcn("_search(m_,w.slice(1),sw,s)")
	b.tttcn("}")
	b.ttcn("};")

	// whole word match
	b.ttcn("this.match=function(w){")
	b.tttcn("let s=new Set();")
	b.tttcn("_search(data.index,w,false,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.matchPrefix=function(w){")
	b.tttcn("let s=new Set();")
	b.tttcn("_search(data.index,w,true,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.matchSuffix=function(w){")
	b.tttcn("let s=new Set();")
	b.tttcn("_search(data.partial,w,false,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.matchSubstring=function(w){")
	b.tttcn("let s=new Set();")
	b.tttcn("_search(data.partial,w,true,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	// m: map
	// w: search word
	// r: prev edit cost row
	// t: result string, built 1 char at a time
	// l: limit
	// s: result set
	b.ttcn("var _fuzzy=function(m,w,r,t,l,sw,s){")
	b.tttcn("let n=w.length;") // the row is one longer than the length of searched word
	b.tttcn("let d=t.length;") // recursion depth
	b.tttcn("let i=t.length;") // row index

	// start by building the next row
	b.tttcn("let R=Array(n+1);")
	b.tttcn("R[0]=i;")
	b.tttcn("for(let j=1;j<=n;j++){")
	b.ttttcn("if(i==0){")
	b.tttttcn("R[j]=j")
	b.ttttcn("}else{")
	b.tttttcn("let dc=r[j]+1;")
	b.tttttcn("let ic=R[j-1]+1;")
	b.tttttcn("let sc=r[j-1]+((w.codePointAt(j-1)==t.codePointAt(i-1))?0:1);")
	b.tttttcn("R[j]=Math.min(dc,Math.min(ic, sc));")
	b.ttttcn("}")
	b.tttcn("}")

	// current best case edit distance
	b.tttcn("let ll=R[(i<n+1)?i:n];")

	// current worst case edit distance
	b.tttcn("let ul=R[n];")

	// quit if best case already surpasses limit
	b.tttcn("if(ll>l){return}")

	// if we are doing a prefix search then an adequate ul can already be used to collect
	b.tttcn("if(sw&&ul<=l){")
	b.ttttcn("_collect(m,s[ul]);")
	b.ttttcn("return;")
	b.tttcn("}")

	b.tttcn("for(let k in m){")
	b.ttttcn("if(k=='pages'){")
	b.tttttcn("if(ul<=l){")
	b.ttttttcn("m.pages.forEach((i)=>{s[ul].add(i)});")
	// otherwise ignore the pages of this node
	b.tttttcn("}")
	b.ttttcn("}else{")
	// deeper recursion
	b.tttttcn("_fuzzy(m[k],w,R,t+k,l,sw,s);")
	b.ttttcn("}")
	b.tttcn("}")
	b.ttcn("};")

	b.ttcn("var _initFuzzy=function(l){")
	b.tttcn("if(l<0){throw new Error('fuzzy limit must be positive')}")
	b.tttcn("let s=new Array(l+1);")
	b.tttcn("for(let i=0;i<=l;i++){")
	b.ttttcn("s[i]=new Set();")
	b.tttcn("};")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.fuzzy=function(w,l){")
	b.tttcn("let s=_initFuzzy(l);")
	b.tttcn("_fuzzy(data.index,w,null,'',l,false,s);") // first row is filled automatically internally
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.fuzzyPrefix=function(w,l){")
	b.tttcn("let s=_initFuzzy(l);")
	b.tttcn("_fuzzy(data.index,w,null,'',l,true,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.fuzzySuffix=function(w,l){")
	b.tttcn("let s=_initFuzzy(l);")
	b.tttcn("_fuzzy(data.partial,w,null,'',l,false,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	b.ttcn("this.fuzzySubstring=function(w,l){")
	b.tttcn("let s=_initFuzzy(l);")
	b.tttcn("_fuzzy(data.partial,w,null,'',l,true,s);")
	b.tttcn("return s;")
	b.ttcn("};")

	// end of constructor
	b.tcn("}")

	b.tcn("set onready(fn){")
	b.ttcn("if(this.page(0)!==null){")
	b.tttcn("fn()")
	b.ttcn("}else{")
	b.tttcn("this._onready.push(fn)")
	b.ttcn("}")
	b.tcn("}")
	b.c("}")
	b.n()

	// TODO: add other searches

	return b.String()
}

var searchIndexHeader = &SearchIndexHeader{newHeaderData("SearchIndex")}

func ActivateSearchIndexHeader() {
	ResolveHeaderActivity(searchIndexHeader, context.NewDummyContext())
}
