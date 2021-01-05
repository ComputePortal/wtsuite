# lists of all the htmlpp command-line tools 
cmds = wt-site wt-search-index wt-template wt-template-syntax-tree wt-script wt-script-syntax-tree wt-svg-minify wt-script-refactor wt-script-grapher wt-glsl wt-glsl-syntax-tree

# temporary destination directory of the commands, also used by temporary files from math-font
export build = $(abspath ./build)

# installation directory of the commands
prefix = /usr/local/bin

.PHONY: math-font
# package files on which all the commands depend
pkg = $(shell find ./pkg/ -name \*.go)

GIT_COMMIT=$(shell git rev-list -1 HEAD)

dsts = $(addprefix $(build)/,$(cmds))

all: math-font $(dsts)

math-font:
	make -C $@

$(build)/htmlpp: extra_flags=-ldflags "-X main.GitCommit="$(GIT_COMMIT)

.SECONDEXPANSION:

$(dsts): $$(shell find ./cmd/$$(notdir $$@) -name \*.go) $(pkg) | $(build)
	cd $(dir $<); \
	go build -o $(abspath $@) $(extra_flags)

$(build):
	mkdir -p $@

clean:
	rm $(build)/*

install: all
	sudo cp -t $(prefix) $(dsts)

install-vim: install-vim-syntax install-vim-indent

install-vim-%:
ifneq ($(wildcard $(HOME)/.vim/$*),)
	cp ./syntax/wtt_$*.vim $(HOME)/.vim/$*/wtt.vim
	cp ./syntax/wts_$*.vim $(HOME)/.vim/$*/wts.vim
	cp ./syntax/glsl_$*.vim $(HOME)/.vim/$*/glsl.vim
else ifneq ($(wildcard $(HOME)/.config/nvim/$*),)
	cp ./syntax/wtt_$*.vim $(HOME)/.config/nvim/$*/wtt.vim
	cp ./syntax/wts_$*.vim $(HOME)/.config/nvim/$*/wts.vim
	cp ./syntax/glsl_$*.vim $(HOME)/.config/nvim/$*/glsl.vim
else
	$(warning no vim $* directory found)
endif

uninstall:
	sudo rm -f $(addprefix $(prefix)/,$(cmds))
