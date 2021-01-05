package macros

type WebGLProgramHeader struct {
  HeaderData
}

func (h *WebGLProgramHeader) Dependencies() []Header {
  return []Header{}
}

func (h *WebGLProgramHeader) Write() string {
  b := NewHeaderBuilder()

  b.n()

  b.cccn("function ", h.Name(), "(gl,v,f){")
  b.tcn("let vs=gl.createShader(gl.VERTEX_SHADER);")

  b.tcn("gl.shaderSource(vs,v);")
  b.tcn("gl.compileShader(vs);")
  b.tcn("let ve=gl.getShaderInfoLog(vs);")
  b.tcn("if(ve.length!=0){if(gl.getError()==gl.NO_ERROR){console.log(ve)}else{throw new Error('VertexShader '+ve);}}")

  b.tcn("gl.shaderSource(fs,f);")
  b.tcn("gl.compileShader(fs);")
  b.tcn("let fe=gl.getShaderInfoLog(fs);")
  b.tcn("if(fe.length!=0){if(gl.getError()==gl.NO_ERROR){console.log(fe)}else{throw new Error('FragmentShader '+fe);}}")

  b.tcn("let p=gl.createProgram();")
  b.tcn("gl.attachShader(p,vs);")
  b.tcn("gl.attachShader(p,fs);")
  b.tcn("gl.linkProgram(p);")
  b.tcn("return p;")
  b.c("}")
  b.n()

  return b.String()
}

var webGLProgramHeader = &WebGLProgramHeader{newHeaderData("__newWebGLProgram__")}
