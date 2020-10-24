" Vim indent file
" Language:		ui
if exists("b:did_indent")
  finish
endif
let b:did_indent = 1

setlocal nosmartindent
setlocal autoindent " usefull when going to newline (I think that our rules are then appliad immediately after)

" Now, set up our indentation expression and keys that trigger it.
setlocal indentexpr=GetUIIndent()
setlocal indentkeys=0{,0},0),0[,0],!^F,o,O,e

" Only define the function once.
if exists("*GetUIIndent")
  finish
endif

let s:cpo_save = &cpo
set cpo&vim


function s:IsInStringOrComment(lnum, cnum)
  let n=synIDattr(synID(a:lnum, a:cnum, 1), 'name')
  return (n=='Comment' || n == 'String')
endfunction

function s:IsInBlock(lnum, cnum)
  let n=synIDattr(synID(a:lnum, a:cnum, 0), 'name')
  echo n
  return (n=='Block')
endfunction

function s:CurrentIndent(lnum)
  let line = getline(a:lnum)
  let ind = matchend(line, '^\s*')
  if ind < 0 
    return 0
  else
    return ind
  endif
endfunction

function s:LineBracketCount(lnum)
  let open_count = 0
  let line = getline(a:lnum)
  let pos = match(line, '[][(){}]', 0)
  while pos != -1
    if (s:IsInStringOrComment(a:lnum, pos+1) == 0) 
      let idx = stridx('(){}[]', line[pos])
      if (idx % 2) == 0
        let open_count = open_count + 1
      else
        let open_count = open_count - 1
      endif

      if open_count < 0 
        let open_count = 0
      endif
    endif

    let pos = match(line, '[][(){}]', pos + 1)
  endwhile

  return open_count
endfunction

function GetUIIndent()
  " get indent of current line (based on first non-blank line
  let reflnum = prevnonblank(v:lnum)
  if (reflnum == 0)
    return 0
  endif

  let refline = getline(reflnum)
  let thisline = getline(v:lnum)

  " ignore multiline strings/comments in first version
  let ind = -1

  " if lnum is this line, go one back
  let prevlnum = reflnum
  if prevlnum == v:lnum
    let prevlnum = prevnonblank(v:lnum - 1)
  endif


  " increase the indent if the prev line has opening brackets
  if (s:LineBracketCount(prevlnum) > 0)
    if (thisline =~ '^\s*[}\])]')
      let ind = s:CurrentIndent(prevlnum)
    else
      let ind = s:CurrentIndent(prevlnum) + &sw
    endif
  else
    if (thisline =~ '^\s*[}\])]')
      let ind = s:CurrentIndent(prevlnum) - &sw
    else
      if (s:IsInBlock(reflnum, 1) == 1)
        let ind = s:CurrentIndent(prevlnum)
      else
        let ind = s:CurrentIndent(reflnum)
      endif
    endif
  endif

  return ind
endfunction

let &cpo = s:cpo_save
unlet s:cpo_save
