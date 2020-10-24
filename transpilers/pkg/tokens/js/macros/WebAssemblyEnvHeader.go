package macros

import (
	"../../context"
)

type WebAssemblyEnvHeader struct {
	HeaderData
}

func (h *WebAssemblyEnvHeader) Dependencies() []Header {
	return []Header{}
}

func (h *WebAssemblyEnvHeader) Write() string {
	b := NewHeaderBuilder()

	b.n()

	b.cccn("class ", h.Name(), "{constructor(fs=null){")
	b.tcn("var scope=this;")
	b.tcn("this.memory=new WebAssembly.Memory({initial:2});")
	b.tcn("this.heapOffset=null;") // set later

	b.tcn("var uint8View,uint32View,int32View,float32View,float64View;")
	b.tcn("var attachViews=function(){")
	b.ttcn("uint8View=new Uint8Array(scope.memory.buffer);")
	b.ttcn("uint32View=new Uint32Array(scope.memory.buffer);")
	b.ttcn("int32View=new Int32Array(scope.memory.buffer);")
	b.ttcn("float32View=new Float32Array(scope.memory.buffer);")
	b.ttcn("float64View=new Float64Array(scope.memory.buffer);")
	b.tcn("};")
	b.tcn("attachViews();")

	// these are 'virtual' pointers (i.e. just linear indices in the memory
	// translation using __heap_base must be done in c
	b.tcn("var blocks=[{f:true,l:-1,p:0}];")

	b.tcn("var d=new TextDecoder('ascii');")

	// TODO: remove the unused constants
	b.tcn("const PAGE_SIZE=64*1024;") // page size is 64 KiB

	b.tcn("const PROT_NONE=0x0;")
	b.tcn("const PROT_READ=0x1;")
	b.tcn("const PROT_WRITE=0x2;")
	b.tcn("const PROT_EXEC=0x4;")

	b.tcn("const MAP_SHARED=0x01;")
	b.tcn("const MAP_PRIVATE=0x02;")
	b.tcn("const MAP_SHARED_VALIDATE=0x03;")

	b.tcn("const O_RDONLY=0o0;")
	b.tcn("const O_WRONLY=0o1;")
	b.tcn("const O_RDWR=0o2;")
	b.tcn("const O_CREAT=0o100;")
	b.tcn("const O_EXCL=0o200;")
	b.tcn("const O_NOCTTY=0o400;")
	b.tcn("const O_TRUNC=0o1000;")
	b.tcn("const O_APPEND=0o2000;")
	b.tcn("const O_NONBLOCK=0o4000;")
	b.tcn("const O_NDELAY=O_NONBLOCK;")
	b.tcn("const O_SYNC=0o4010000;")
	b.tcn("const O_FSYNC=O_SYNC;")
	b.tcn("const O_ASYNC=0o20000;")

	b.tcn("const SEEK_SET=0;")
	b.tcn("const SEEK_CUR=1;")
	b.tcn("const SEEK_END=2;")

	// private utility functions
	b.tcn("var getString=function(p){")
	b.ttcn("var q=p;while(uint8View[q]!=0){q++};")
	b.ttcn("return d.decode(uint8View.subarray(p,q))")
	b.tcn("};")

	b.tcn("var setString=function(p,s){")
	b.ttcn("for(var i=0;i<s.length;i++){")
	b.tttcn("uint8View[p]=s.codePointAt(i);")
	b.ttcn("}")
	b.ttcn("uint8View[s.length]=0;")
	b.ttcn("return s.length;")
	b.tcn("};")

	b.tcn("var setBytes=function(p,b){")
	b.ttcn("for(var i=0;i<b.length;i++){")
	b.tttcn("uint8View[p+i]=b[i]")
	b.ttcn("}")
	b.ttcn("return b.length")
	b.tcn("};")

	b.tcn("var getBytes=function(p,n){")
	b.ttcn("return uint8View.slice(p,p+n)")
	b.tcn("};")

	b.tcn("var getInt=function(p){return int32View[p/4]};")

	b.tcn("var getFloat=function(p){return float32View[p/4]};")
	b.tcn("var getDouble=function(p){return float64View[p/8]};")

	b.tcn("var format=function(s,f){")
	// find all %(d|g|...) strings and replace with value on heap
	// XXX: do variadic functions give heap conflicts?
	b.ttcn("var r=/\\%[dfg]/;")
	b.ttcn("return s.replace(r,function(match){")
	b.tttcn("if (f.length==0){throw new Error('too few var args parameters')};")
	b.tttcn("var a=f.shift();")
	b.tttcn("switch(match){")
	b.ttttcn("case '%d':return getInt(a).toString();")
	b.ttttcn("case '%f':return getFloat(a).toString();")
	b.ttttcn("case '%g':return getDouble(a).toString();")
	b.ttttcn("default:throw new Error('unrecognized ' + match);")
	b.tttcn("}")
	b.ttcn("})")
	b.tcn("};")

	b.tcn("var findBlock=function(p){")
	b.ttcn("var i=blocks.findIndex((b)=>{return b.p==p});")
	b.ttcn("if (i==-1){throw new Error('attempting to free/realloc non-existent block')}")
	b.ttcn("return i")
	b.tcn("};")

	b.tcn("var assertFS=function(){")
	b.ttcn("if(fs==null){throw new Error('no FS defined')}")
	b.tcn("};")

	b.tcn("var assertInRange=function(fd,p){")
	b.ttcn("if(p<0||p>fs.size(fd)-1){throw new Error('file position ' + p.toString() + ' out of range ')}")
	b.tcn("};")

	// clib surrogates
	b.tcn("this.printf=function(s,...f){s=format(getString(s),f);console.log(s);return s.length};")
	b.tcn("this.puts=function(s){console.log(getString(s));return 0};")
	b.tcn("this.errorf=function(s,...f){s=format(getString(s),f);console.error(s);return s.length};")
	b.tcn("this.fprintf=function(fd,s,...f){if(fd==1){return scope.printf(s,...f)}else if(fd==2){return scope.errorf(s, ...f)}else{throw new Error('formatted print only available for stdout and stderr')}};")
	b.tcn("this.exit=function(c){throw new Error(c)};")

	b.tcn("this.sprintf=function(d,s,...f){")
	b.ttcn("var t=format(getString(s),f);")
	b.ttcn("setString(d,t);")
	b.ttcn("return t.length;")
	b.tcn("};")

	b.tcn("this.memcpy=function(d,s,n){uint8View.copyWithin(d,s,s+n)};")
	b.tcn("var memcpy=function(d,s,n){scope.memcpy(d+scope.heapOffset,s+scope.heapOffset,n)};")

	b.tcn("var grow=function(p){")
	b.ttcn("if(p>scope.memory.buffer.byteLength){")
	b.tttcn("scope.memory.grow(Math.ceil((p-scope.memory.buffer.byteLength)/PAGE_SIZE));")
	b.tttcn("attachViews()")
	b.ttcn("}")
	b.tcn("};")

	b.tcn("var free=function(p){")
	b.ttcn("var i=findBlock(p);")
	b.ttcn("var b=blocks[i];")
	b.ttcn("if (i>0&&i<(blocks.length-1)&&blocks[i-1].f&&i+1==blocks.length-1&&blocks[i+1].f){")
	b.tttcn("blocks.splice(i-1,3,{f:false,l:-1,p:blocks[i-1].p})")
	b.ttcn("}else if (i>0&&i<(blocks.length-1)&&blocks[i-1].f&&blocks[i+1].f){")
	b.tttcn("blocks.splice(i-1,3,{f:false,l:blocks[i-1].l+b.l+blocks[i+1].l,p:blocks[i-1].p})")
	b.ttcn("}else if(i>0&&blocks[i-1].f&&i==(blocks.length-1)){")
	b.tttcn("blocks.splice(i-1,2,{f:false,l:-1,p:blocks[i-1].p})")
	b.ttcn("}else if(i>0&&blocks[i-1].f){")
	b.tttcn("blocks.splice(i-1,2,{f:false,l:blocks[i-1].l+b.l,p:blocks[i-1].p})")
	b.ttcn("}else if(i<(blocks.length-1)&&blocks[i+1].f&&i+1==blocks.length-1){")
	b.tttcn("blocks.splice(i,2,{f:false,l:-1,p:b.p})")
	b.ttcn("}else if(i<(blocks.length-1)&&blocks[i+1].f){")
	b.tttcn("blocks.splice(i,2,{f:false,l:b.l+blocks[i+1].l,p:b.p})")
	b.ttcn("}else if(i==blocks.length-1){")
	b.tttcn("blocks[i]={f:false,l:-1,p:b.p}")
	b.ttcn("}else{")
	b.tttcn("blocks[i]={f:false,l:b.l,p:b.p}")
	b.ttcn("}")
	b.tcn("};")

	b.tcn("var malloc=function(n){")
	b.ttcn("for(var i=0;i<blocks.length-1;i++){")
	b.tttcn("var b=blocks[i];")
	b.tttcn("if(b.f&&(b.l>=n)){")
	b.ttttcn("if(b.l==n){blocks[i]={f:false,l:b.l,p:b.p}}")
	b.ttttcn("else{blocks.splice(i,1,{f:false,l:n,p:b.p},{f:true,l:b.l-n,p:b.p+n})};")
	b.ttttcn("return b.p")
	b.tttcn("}")
	b.ttcn("}")
	b.ttcn("var b=blocks[blocks.length-1];")
	b.ttcn("blocks.splice(blocks.length-1,1,{f:false,l:n,p:b.p},{f:true,l:-1,p:b.p+n});")
	b.ttcn("return b.p")
	b.tcn("};")

	b.tcn("var realloc=function(p,n){")
	b.ttcn("if(n==0){free(p);return 0}")
	b.ttcn("var i=findBlock(p);")
	b.ttcn("var b=blocks[i];")
	b.ttcn("if(b.l<n){")
	b.tttcn("if(i<(blocks.length-1)&&blocks[i+1].l==-1){")
	b.ttttcn("blocks[i]={f:false,l:n,p:b.p};")
	b.ttttcn("blocks[i+1]={f:true,l:-1,p:b.p+n};")
	b.ttttcn("return b.p")
	b.tttcn("}else if(i<(blocks.length-1)&&blocks[i+1].f&&blocks[i+1].l>(n-b.l)){")
	b.ttttcn("blocks[i]={f:false,l:n,p:b.p};")
	b.ttttcn("blocks[i+1]={f:true,l:blocks[i+1].l-(n-b.l),p:b.p+n};")
	b.ttttcn("return b.p")
	b.tttcn("}else{")
	b.ttttcn("var q=malloc(n);memcpy(q,p,b.l);free(p);return q")
	b.tttcn("}")
	b.ttcn("}else if(b.l>n){")
	b.tttcn("if(i<(blocks.length-1)&&blocks[i+1].f){")
	b.ttttcn("blocks[i]={f:false,l:n,p:b.p};")
	b.ttttcn("blocks[i+1]={f:true,l:(blocks[i+1].l==-1)?-1:blocks[i+1].l+(b.l-n),p:b.p+n};")
	b.tttcn("}else{")
	b.ttttcn("blocks.splice(i,1,{f:false,l:n,p:b.p},{f:true,l:b.l-n,p:b.p+n})")
	b.tttcn("};")
	b.tttcn("return b.p")
	b.ttcn("}else{")
	b.tttcn("return p")
	b.ttcn("}")
	b.tcn("};")

	b.tcn("this.realloc=function(p,n){")
	b.ttcn("var q=realloc(p-scope.heapOffset,n)+scope.heapOffset;")
	b.ttcn("grow(q+n);")
	b.ttcn("return q")
	b.tcn("};")

	b.tcn("this.free=function(p){")
	b.ttcn("if(p!=0){") // null ptr is ignored (like stdlib version of free)
	b.tttcn("free(p-scope.heapOffset)")
	b.ttcn("}")
	b.tcn("};")

	b.tcn("this.malloc=function(n){")
	b.ttcn("var p=malloc(n)+scope.heapOffset;")
	b.ttcn("grow(p+n);")
	b.ttcn("return p")
	b.tcn("};")

	b.tcn("this.open=function(c,f){assertFS();")
	b.ttcn("var s=getString(c);")
	b.ttcn("if(f==O_RDONLY){")
	b.tttcn("if(!fs.exists(s)){throw new Error('file ' + s + ' not found')};")
	b.tttcn("return fs.open(s)")
	b.ttcn("}else if(f==O_WRONLY){")
	b.tttcn("if(fs.exists(s)){throw new Error('file ' + s + ' already exists')};")
	b.tttcn("return fs.create(s)")
	b.ttcn("}else{throw new Error('unsupported mode')}")
	b.tcn("};")

	b.tcn("this.close=function(fd){assertFS();")
	b.ttcn("fs.close(fd);")
	b.ttcn("return 0")
	b.tcn("};")

	b.tcn("this.fopen=function(c,f){assertFS();")
	b.ttcn("var s=getString(c);")
	b.ttcn("var m=getString(f);")
	b.ttcn("if(m=='rb'){")
	b.tttcn("if(!fs.exists(s)){throw new Error('file '+s+' not found (flags: '+m+')')};")
	b.tttcn("return fs.open(s)")
	b.ttcn("}else if(m=='wb'){")
	b.tttcn("if(fs.exists(s)){throw new Error('file '+s+' already exists (flags: '+m+')')};")
	b.tttcn("return fs.create(s)")
	b.ttcn("}else{throw new Error('unsupported mode ' + m)}")
	b.tcn("};")

	b.tcn("this.fclose=function(fd){assertFS();")
	b.ttcn("fs.close(fd);")
	b.ttcn("return 0")
	b.tcn("};")

	b.tcn("this.fseek=function(fd,o,w){assertFS();")
	b.ttcn("var p=o;")
	b.ttcn("if(w==SEEK_CUR){")
	b.tttcn("p+=fs.tell(fd)")
	b.ttcn("}else if(w==SEEK_END){")
	b.tttcn("p=fs.size(fd)-o")
	b.ttcn("}else if(w!=SEEK_SET){")
	b.tttcn("throw new Error('unrecognized whence specifier')")
	b.ttcn("}")
	b.ttcn("assertInRange(fd,p);fs.seek(fd,p);")
	b.ttcn("return 0")
	b.tcn("};")

	// XXX: are these assertions too much overhead?
	b.tcn("this.read=function(fd,p,n){assertFS();")
	b.ttcn("assertInRange(fd,fs.tell(fd)+n-1);")
	b.ttcn("var b=fs.read(fd,n);")
	b.ttcn("setBytes(p,b);")
	b.ttcn("return n")
	b.tcn("};")

	b.tcn("this.fread=function(p,s,m,fd){")
	b.ttcn("return scope.read(fd,p,s*m)/s")
	b.tcn("};")

	b.tcn("this.ftell=function(fd){")
	b.ttcn("return fs.tell(fd)")
	b.tcn("};")

	b.tcn("this.write=function(fd,p,n){assertFS();")
	b.ttcn("var b=getBytes(p,n);")
	b.ttcn("fs.write(fd,b);")
	b.ttcn("return n")
	b.tcn("};")

	b.tcn("this.fwrite=function(p,s,m,fd){")
	b.ttcn("return scope.write(fd,p,s*m)/s")
	b.tcn("};")

	// XXX: input addr is ignored, just uses malloc
	b.tcn("this.mmap=function(a,n,f,g,fd,o){assertFS();")
	b.ttcn("if(f!=PROT_READ||g!=MAP_SHARED){throw new Error('unsupported protection or map mode')}")
	b.ttcn("var t=fs.tell(fd);")
	b.ttcn("fs.seek(fd,o);")
	b.ttcn("var b=fs.read(fd,n);") // read the whole at once!
	b.ttcn("fs.seek(fd,t);")
	b.ttcn("var p=scope.malloc(n);")
	b.ttcn("setBytes(p,b);")
	b.ttcn("return p")
	b.tcn("};")

	// length is ignored, simply calls free
	b.tcn("this.munmap=function(p,n){")
	b.ttcn("scope.free(p);")
	b.ttcn("return 0")
	b.tcn("};")

	b.c("}}")
	b.n()

	return b.String()
}

var webAssemblyEnvHeader = &WebAssemblyEnvHeader{newHeaderData("WebAssemblyEnv")}

func ActivateWebAssemblyEnvHeader() {
	ResolveHeaderActivity(webAssemblyEnvHeader, context.NewDummyContext())
}
