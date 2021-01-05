" for opengl es shader language
" Language: GLSL
" Filenames: *.glsl

if exists("b:current_syntax")
  finish
endif

syn keyword Todo NOTE TODO FIXME XXX TBD contained

syn match Comment "\/\/.*" contains=Todo
syn match Comment "^[ \t]*\*\($\|[ \t]\+\)"
syn region Comment start="/\*"  end="\*/" contains=Todo

syn match Special "\\\d\d\d\|\\."
syn region String start=+"+  skip=+\\\\\|\\"+  end=+"\|$+	contains=Special
syn region String start=+'+  skip=+\\\\\|\\'+  end=+'\|$+	contains=Special

syn region PreProc	start="^\s*\zs\(%:\|#\)\s*\(version\>\|extension\>\)" skip="\\$" end="$" keepend contains=Comment

syn match Special "'\\.'"

syn match Constant "-\=\<\d\+L\=\>\|0[xX][0-9a-fA-F]\+\>"
syn match Constant '\<\zs\d\+\(\.\d\+\([e][-]\?\d\+\)\?\)\?\ze'
syn keyword Boolean true false

syn keyword Conditional if else
syn keyword Statement return
syn keyword Keyword export import from as
syn keyword Keyword attribute const highp in inout lowp mediump out precision struct uniform varying void

syn keyword Type bool int float
syn keyword Type vec2 vec3 vec4 ivec2 ivec3 ivec4 bvec2 bvec3 bvec4
syn keyword Type mat2 mat3 mat4

let b:current_syntax = "glsl"
