" for htmlpp dialect, made 
" Language: JavsScript
" Filenames: *.ts

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

syn match   Special "'\\.'"

syn match Constant "-\=\<\d\+L\=\>\|0[xX][0-9a-fA-F]\+\>"
syn match Constant '\<\zs\d\+\(\.\d\+\([e][-]\?\d\+\)\?\)\?\ze'

syn keyword Conditional if else switch break continue
syn keyword Repeat while for do in of
syn keyword Operator delete instanceof typeof
syn keyword Type Array Boolean Date Error Function Int Map Math Number Object Promise Set String RegExp 
" types that appear so often that it is worthwhile to highlight them
syn keyword Type Uint8Array Uint16Array Uint32Array Int32Array Float32Array Float64Array Blob
syn keyword Type Element HTMLElement HTMLImageElement HTMLInputElement HTMLTextAreaElement HTMLCanvasElement HTMLSelectElement
syn keyword Type Event MouseEvent WheelEvent KeyboardEvent
syn keyword Type WebGLRenderingContext WebGLProgram CanvasRenderingContext2D
syn keyword Type WebAssembly WebAssemblyEnv 
syn keyword Type IDBDatabase IDBRequest IDBKeyRange IDBCursorWithValue IDBVersionChangeEvent
syn keyword Type URL

syn keyword Statement return
syn keyword Boolean true false
syn keyword Constant null
syn keyword Keyword console document super this window
syn keyword Label case default
syn keyword Exception try catch finally throw
syn keyword Keyword abstract any as async await class const constructor enum export extends from function get implements import interface let new private set static universe var void

"All initial capitalized words are usually types
" syn match javaScriptType '::\zs[A-Z][a-zA-Z0-9_]*\ze'
" syn match javaScriptType '\<\zs[A-Z][a-zA-Z0-9_]*\ze'


let b:current_syntax = "jspp"
