dst = ../pkg/tokens/math/serif

build ?= $(abspath ../build)

font = /usr/share/fonts/truetype/freefont/FreeSerif.ttf
font_maker = ./mathFontMaker.py
font_reader = $(build)/mathFontReader

# the two final go files
woff2 = $(dst)/woff2.go
dimensions = $(dst)/dimensions.go

all: $(woff2) $(dimensions)

$(woff2) $(dimensions) &: $(font) $(font_maker) $(font_reader)
	$(font_maker) $(font) $(font_reader) $(dimensions) $(woff2)

$(font_reader): ./mathFontReader.c | $(build)
	gcc $< -std=c99 -I /usr/include/freetype2 -lfreetype -o $@

$(build): 
	mkdir -p $(build)
