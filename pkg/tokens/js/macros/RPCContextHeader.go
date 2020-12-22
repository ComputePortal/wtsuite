package macros

type RPCContextHeader struct {
  HeaderData
}

func (h *RPCContextHeader) Dependencies() []Header {
  return []Header{objectFromInstanceHeader, objectToInstanceHeader}
}

func (h *RPCContextHeader) Write() string {
  b := NewHeaderBuilder()

  b.n()
  b.cccn("class ", h.Name(), "{")

  b.tcn("constructor(){")

  // queu of raw objects and callbacks 
  b.ttcn("this.queu=[];") 

  b.ttcn("this.channels={};")

  b.ttcn("this.packetCount=0;") // incremented forever
  b.ttcn("this.channelCount=0;") // incremented forever

  b.ttcn("this.notified=null;")
  b.ttcn("this.notify=null;")

  b.ttcn("this.resetNotifier();")
  b.tcn("}")

  b.tcn("resetNotifier(){")
  b.ttcn("this.notified=new Promise((resolve,reject)=>{")
  b.tttcn("this.notify=()=>{resolve();this.resetNotifier()};")
  b.ttcn("});")
  b.tcn("}")

  b.tcn("register(interf,value){")
  b.ttcn("let channel=this.channelCount++;")
  b.ttcn("this.channels[channel]=[interf,value];")
  b.ttcn("return channel;")
  b.tcn("}")

  b.tcn("append(packet,callback,reject){")
  b.ttcn("this.queu.push([packet,callback,reject]);")
  b.ttcn("this.notify()");
  b.tcn("}")

  // called by rpc interface.new instances
  b.tcn("request(packet,channels){")
  b.ttcn("return new Promise((callback,reject)=>{")
  b.tttcn("packet.id=this.packetCount++;")
  b.tttcn("packet.type='request';")
  b.tttcn("this.append(packet,(r)=>{for(let c of channels){delete this.channels[c]};callback(r)},reject);")
  b.ttcn("});")
  b.tcn("}")

  b.tcn("respond(packet){")
  b.ttcn("return new Promise((callback,reject)=>{")
  b.tttcn("packet.type='response';")
  b.tttcn("this.append(packet,callback,reject);")
  b.ttcn("});")
  b.tcn("}")

  // not static because it depends on the options
  b.tcn("serialize(obj){")
  b.ttcccn("return JSON.stringify(", objectFromInstanceHeader.Name(), "(obj));")
  b.tcn("}")

  // not static because it depends on the options
  b.tcn("deserialize(str){")
  b.ttcccn("return ", objectToInstanceHeader.Name(), "(JSON.parse(str));")
  b.tcn("}")

  b.c("}")
  b.n()

  return b.String()

}

var rpcContextHeader = &RPCContextHeader{newHeaderData("__rpcContext__")}
