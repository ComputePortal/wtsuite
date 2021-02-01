" Vim syntax file
" Language: htmlpp-ui
" Filenames: *.ui

if exists("b:current_syntax")
  finish
endif

syn match Constant '[#][0-9abcdef]\{3}'
syn match Constant '[#][0-9abcdef]\{6}'
syn match Constant '[#][0-9abcdef]\{8}'

" syn region IndentifierLine start='^\zs\s*[#!a-zA-Z_][a-zA-Z0-9_\.\-]*\ze\s*' end='$' oneline keepend contains=Identifier

" tag identifiers match anywhere except in containers
syn match Identifier '\s*\zs[#!a-zA-Z_][a-zA-Z0-9_\.\-]*\ze\s*'
" syn match Statement '\s*\zs\(import\|print\|for\|switch\|case\|default\|if\|else\|elseif\|function\)\ze\s\+'
" syn match Statement '\s*\zs\(import\|print\|for\|switch\|case\|default\|if\|else\|elseif\|function\)\ze\s\+'

syn match DecreaseIndent '\s*\zs[\<]\ze\s*'


" syn match Indentifier '\s*\zs[#!a-zA-Z_][a-zA-Z0-9_\.\-]*\ze\s*' containedin=IdentifierLine

" statement keywords that are only sensible on the beginning of a line
syn match Statement '^\zspermissive\ze\s*'
syn match Statement '^\zsparameters\ze\s*'
syn match Statement '^\s*\zsexport\ze\s\+'


" statements can span multiple lines, so these keywords can be anywhere
" syn keyword Statement from
" syn region StatementAttr start='^\zs\s*\ze\(print\|var\|import\|switch\|case\|default\|export\|for\|if\|elseif\)' end='$' oneline contains=ALLBUT,Identifier,DecreaseIndent keepend transparent

"syn match Statement contained '^\s\*\zs\(import\|print\|for\|switch\|case\|default\|if\|else\|elseif\|var\|function\)\ze\s\+'
" syn keyword StatementKeywords from in 

syn keyword TemplateKeywords contained template extends blocks super
syn region Template start='template' end='super' contains=TemplateKeywords,Block,VarAction,Constant,String, Comment keepend
syn keyword TemplateKeywords this

syn keyword VarKeywords contained var 
syn region Var start='var' end='=[^{]*' oneline contains=VarKeywords,Block,Action,Constant,String, Comment

syn keyword FunctionKeywords contained function
syn region Function_ start='function' end='[^(]*' oneline contains=FunctionKeywords

syn keyword AsKeyword contained as
syn keyword ImportKeywords contained  import from
syn region Import start='import' end='$' oneline contains=ImportKeywords,Block,Action,Constant,String, Comment, AsKeyword

syn keyword ExportKeywords contained export from
syn region Export start='export\s*[\{*]' end='$' oneline contains=ExportKeywords,Block,Action,Constant,String, Comment, AsKeyword

syn keyword ForKeywords contained  for in
syn region For start='^\s*for' end='$' oneline contains=ForKeywords,Block,Action,Constant,String, Comment

syn keyword SwitchKeywords contained switch
syn region Switch start='switch' end='$' oneline contains=SwitchKeywords,Block,Action,Constant,String, Comment

syn keyword CaseKeywords contained case
syn region Case start='case' end='$' oneline contains=CaseKeywords,Block,Action,Constant,String, Comment

syn keyword DefaultKeywords contained default
syn region Default start='default' end='$' oneline contains=DefaultKeywords,Block,Action,Constant,String, Comment

syn keyword PrintKeywords contained print
syn region Print start='print' end='$' oneline contains=PrintKeywords,Block,Action,Constant,String, Comment

syn keyword PrintKeywords contained print
syn region Print start='print' end='$' oneline contains=PrintKeywords,Block,Action,Constant,String, Comment

syn keyword IfKeywords contained if
syn region If start='if' end='$' oneline contains=IfKeywords,Block,Action,Constant,String, Comment

syn keyword ElseifKeywords contained elseif
syn region Elseif start='elseif' end='$' oneline contains=ElseifKeywords,Block,Action,Constant,String, Comment

syn keyword ElseKeywords contained else
syn region Else start='else' end='$' oneline contains=ElseKeywords,Block,Action,Constant,String, Comment

syn keyword DeclBlockKeywords contained block
syn region DeclBlock start='block' end='$' oneline contains=DeclBlockKeywords,Comment

syn keyword ReplaceKeywords contained replace
syn region Replace start='replace' end='$' oneline contains=ReplaceKeywords,Block,Action,Constant,String, Comment

syn keyword AppendKeywords contained append
syn region Append start='append' end='$' oneline contains=AppendKeywords,Block,Action,Constant,String, Comment

syn keyword PrependKeywords contained prepend
syn region Prepend start='prepend' end='$' oneline contains=PrependKeywords,Block,Action,Constant,String, Comment

syn region Block start='(' end=')' contains=Block, Action, Comment, String, Constant
syn region Block start='{' end='}' contains=Block, Action, Comment, String, Constant, AsKeyword
syn region Block start='\[' end='\]' contains=Block, Action, Comment, String, Constant

syn match Action '\zs[\$][a-zA-Z_][a-zA-Z0-9_\.\-]*\ze'
syn match Action '\zs[\$]\?[a-zA-Z_][a-zA-Z0-9_\.\-]*\ze[(\[]' contained
syn match VarAction '\zs[\$][a-zA-Z_][a-zA-Z0-9_\.\-]*\ze' contained

"syn match Statement '^\zs\(export\|import\)\ze\s\+'
"syn keyword Statement class var if ifelse else for in switch case default dummy
"
"syn match Special '\zs[\\]\ze'
"syn region LineContinuation start='\\' end='.*[^\s]\+.*$' contains=ALLBUT,Statement,Identifier keepend
"
"
"syn match StatementAttr '[ ]*\zsextends[ ]*\ze='
"syn match StatementAttr '[ ]*\zsas[ ]*\ze='
"syn match StatementAttr '[ ]*\zsfrom[ ]*\ze='
"syn match StatementAttr '[ ]*\zsaliases[ ]*\ze='
"
"syn region FunctionStatement start='\(^export[ ]\+\)\?\zsfunction\ze' end='[(]' contains=FunctionDef1, FunctionDef2
"syn keyword FunctionDef1 function contained
"syn keyword ExportFunction export contained
"
"syn match Statement '^[ ]*\zs[_]\ze[ ]*'
"
syn keyword Constant true false null
syn region String start='\'' end='\''
syn region String start='\"' end='\"'
syn match Constant '\<\zs\d\+\(\.\d\+\)\?\(px\|[%]\|rem\|fr\|em\|s\|vh\|vw\)\?\ze'
"
"syn match FunctionDef2 '[\$]\?[a-zA-Z_][a-zA-Z0-9_\.\-]*' contained
"
"
syn keyword Todo contained TODO XXX NOTE
"
syn match	Comment	"\/\/.*" contains=Todo
syn region	Comment	start="/\*" end="\*/" extend contains=Todo

"hi def link StatementAttr Statement
"hi def link FunctionDef1 Statement
"hi def link FunctionDef2 Normal

hi def link Action PreProc
hi def link VarAction PreProc

hi def link TemplateKeywords Statement
hi def link VarKeywords Statement
hi def link ForKeywords Statement
hi def link AsKeyword Statement
hi def link ImportKeywords Statement
hi def link ExportKeywords Statement
hi def link SwitchKeywords Statement
hi def link CaseKeywords Statement
hi def link DefaultKeywords Statement
hi def link PrintKeywords Statement
hi def link IfKeywords Statement
hi def link ElseifKeywords Statement
hi def link ElseKeywords Statement
hi def link ReplaceKeywords Statement
hi def link DeclBlockKeywords Statement
hi def link AppendKeywords Statement
hi def link PrependKeywords Statement
hi def link FunctionKeywords Statement

hi def link StatementKeywords Statement
hi def link DecreaseIndent Statement

let b:current_syntax = "ui"
