package macros

import (
  "../../js"
)

type RPCClientHeader struct {
  HeaderData
}

func (h *RPCClientHeader) Dependencies() []Header {
  return []Header{rpcContextHeader}
}

func (h *RPCClientHeader) Write() string {
  b := NewHeaderBuilder()

  // 'this' is irrelevant
  b.n()
  b.cccn("var ", h.Name(), "=(function(interf,fn){")
  b.tcccn("let ctx=new ", rpcContextHeader.Name(), "();")

  // infinite loop function
  b.tcn("(async function(){")
  b.ttcn("let callbacks={};") // save callbacks, because replies dont necessarily correspond to preceding requests
  b.ttcn("while(true){")
  b.tttcn("await ctx.notified;")
  b.tttcn("if(ctx.queu.length===0){continue};")

  // fifo queu
  b.tttcn("let tuple=ctx.queu.shift();")
  b.tttcn("let message=tuple[0];")
  b.tttcn("let callback=tuple[1];")
  b.tttcn("let reject=tuple[2];")
  b.tttcn("if(message==undefined||message.type==undefined){")
  b.ttttcn("reject(new Error('malformed message'))")
  b.tttcn("}else if(message.type==='response'){")
  b.ttttcn("callback();") // shoot and forget
  b.tttcn("}else if(message.type==='request'){")
  b.ttttcn("if(message.id==undefined||!Number.isInteger(message.id)){reject(new Error('malformed request message'))};")
  b.ttttcn("callbacks[message.id]=callback;") // save the callback for an expected later response
  b.tttcn("}else{")
  b.ttttcn("reject(new Error('invalid message type: '+message.type));")
  b.tttcn("}")

  b.tttcn("let reply;")
  b.tttcn("try{")
  b.ttttcn("reply=ctx.deserialize(await fn(ctx.serialize(message)));")
  b.tttcn("}catch(e){")
  b.ttttcn("reject(new Error('malformed response'));")
  b.tttcn("}")

  b.tttcn("if(reply==undefined||reply.type==undefined||reply.id==undefined||!Number.isInteger(reply.id)){")
  b.ttttcn("reject(new Error('malformed response'))")
  b.tttcn("}else if(reply.type==='response'){")
  b.ttttcn("let callback=callbacks[reply.id];")
  b.ttttcn("if(callback===undefined){reject(new Error('bad message id'));}")
  b.ttttcn("callback(reply);") // should delete the subChannels
  b.ttttcn("delete callbacks[reply.id];")
  b.tttcn("}else if(reply.type==='request'){")
  b.ttttcn("if(reply.channel==undefined||!Number.isInteger(reply.channel)){reject(new Error('malformed incoming request'))};")
  b.ttttcn("let pair=ctx.channels[reply.channel];")
  b.ttttcn("if(pair==undefined){reject(new Error('invalid request channel'))};")
  b.ttttcn("let interf=pair[0];")
  b.ttttcn("let value=pair[1];")
  b.ttttcn("interf.rpc(value,reply,ctx);") // shoot and forget (handle by queu anyway)
  b.tttcn("}else{")
  b.ttttcn("reject(new Error('invalid reply type: '+reply.type));")
  b.tttcn("}")
  b.ttcn("}")
  b.tcn("}());")

  b.tcccn("return interf.", js.NewRPCClientMemberName, "(0,ctx);")
  b.c("});")
  b.n()

  return b.String()
}

var rpcClientHeader = &RPCClientHeader{newHeaderData("RPCClient")}
