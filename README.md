# Introduction
*wtsuite* is collection of transpilers and other tools for html5 technologies
* wt-script
* wt-script-grapher
* wt-script-refactor
* wt-script-syntax-tree
* wt-search-index
* wt-site
* wt-svg-minify
* wt-template-syntax-tree

Documentation can be found [here](https://computeportal.github.io/wtsuite-doc)

# Build dependencies
* shell
* make
* libfreetype-dev (v2)
* python3-fonttools
* fonts-freefont-ttf

There should be no runtime dependencies

# Compiling and installing on linux
Running `make` builds everything
Running `make install` builds and installs into /usr/local/bin

# VIM syntax/indentation
In your `~/.vimrc` or `~/.config/nvim/init.vim` file:

```
au BufNewFile,BufRead *.wtt set filetype=wtt
au BufNewFile,BufRead *.wts set filetype=wts
```

Run `make install-vim` to install the relevant syntax/indentation files.

# License
[GPLv3](./LICENSE.txt)

# TODO
* Simplify js values
* Can interfaces that refer to prototypes be removed?
* This means that all members must be registered on the value
* How to do inheritance n 
