package macros

type RPCServerHeader struct {
  HeaderData
}

func (h *RPCServerHeader) Dependencies() []Header {
  return []Header{rpcContextHeader}
}

func (h *RPCServerHeader) Write() string {
  b := NewHeaderBuilder()

  b.n()
  b.cccn("var ", h.Name(), "=(function(interf,value){")
  b.tcccn("let ctx=new ", rpcContextHeader.Name(), "();")
  b.tcn("ctx.register(interf,value);")
  b.tcn("let callbacks={};")
  b.tcn("return {")
  b.ttcn("handle: async function(str){")

  b.tttcn("try{") // try catch all

  b.tttcn("let message=ctx.deserialize(str);")

  b.tttcn("if(message.type==='request'){")
  b.ttttcn("let pair=ctx.channels[message.channel];")
  b.ttttcn("let interf=pair[0];")
  b.ttttcn("let value=pair[1];")
  b.ttttcn("interf.rpc(value,message,ctx);")
  b.tttcn("}else if(message.type==='response'){")
  b.ttttcn("let callback=callbacks[message.id];")
  b.ttttcn("if(callback===undefined){throw new Error('bad message id');};")
  b.ttttcn("callback(message);")
  b.ttttcn("delete callbacks[message.id];")
  b.tttcn("}else{")
  b.ttttcn("throw new Error('invalid type: '+message.type);")
  b.tttcn("}")

  b.tttcn("while(ctx.queu.length===0){")
  b.ttttcn("await ctx.notified;")
  b.tttcn("}")

  b.tttcn("let pair=ctx.queu.shift();")
  b.tttcn("let reply=pair[0];")
  b.tttcn("let callback=pair[1];")
  b.tttcn("if(reply.type==='response'){")
  b.ttttcn("callback();")
  b.tttcn("}else if(reply.type==='request'){")
  b.ttttcn("callbacks[reply.id]=callback;")
  b.tttcn("}else{")
  b.ttttcn("throw new Error('invalid type: '+reply.type);")
  b.tttcn("}")

  b.tttcn("return ctx.serialize(reply);")

  b.tttcn("}catch(e){")
  b.ttttcn("return ctx.serialize(new Error(e.message));")
  b.tttcn("}")
  b.ttcn("}")
  b.tcn("}")
  b.c("});")
  b.n()

  return b.String()
}

var rpcServerHeader = &RPCServerHeader{newHeaderData("RPCServer")}
