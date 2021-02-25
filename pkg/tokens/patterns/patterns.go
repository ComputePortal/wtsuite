package patterns

import (
	"regexp"
	"strings"
)

const (
	ESCAPED_NL = `\n`

	HASH       = "#"
	BOOL_TRUE  = "true"
	BOOL_FALSE = "false"
	NULL       = "null"

	PERIOD    = "."
	COLON     = ":"
	COMMA     = ","
	SEMICOLON = ";"
	EQUAL     = "="
  DOLLAR    = "$"

	SPLAT       = "..."
	DCOLON      = "::"
	COLON_EQUAL = ":="
	PLUS_PLUS   = "++"
	MINUS_MINUS = "--"
	ARROW       = "=>"
	PLUS_EQUAL  = "+="
	MINUS_EQUAL = "-="
	MUL_EQUAL   = "*="

	BRACES_START   = `{`
	BRACES_STOP    = `}`
	BRACKETS_START = `[`
	BRACKETS_STOP  = `]`
	PARENS_START   = `(`
	PARENS_STOP    = `)`
	ANGLED_START   = `<`
	ANGLED_STOP    = `>`
	ANGLED_STOP2   = `>>`
	ANGLED_STOP3   = `>>>`

	TAG_START            = `<`
	TAG_STOP_SELFCLOSING = `/>`

	NAMESPACE_SEPARATOR = `.`

	SQ_STRING_START  = `'`
	SQ_STRING_STOP   = SQ_STRING_START
	DQ_STRING_START  = `"`
	DQ_STRING_STOP   = DQ_STRING_START
	BT_FORMULA_START = "`"
	BT_FORMULA_STOP  = BT_FORMULA_START

	SL_COMMENT_START  = `//`
	ML_COMMENT_START  = `/*`
	ML_COMMENT_STOP   = `*/`
	XML_COMMENT_START = `<!--`
	XML_COMMENT_STOP  = `-->`

  HTML_SCRIPT_START = `<[\s]*script`
  HTML_SCRIPT_STOP = `</[\s]*script[\s]*>`
  HTML_STYLE_START = `<[\s]*style`
  HTML_STYLE_STOP = `</[\s]*style[\s]*>`
)

// internal names
const (
	//INTERNAL_STYLE_TREE = "__styleTree__"
)

var (
	XML_STRING_OR_COMMENT_REGEXP = compileRegexp(SQ_STRING_START,
		DQ_STRING_START, SL_COMMENT_START,
		ML_COMMENT_START, ML_COMMENT_STOP,
		XML_COMMENT_START, XML_COMMENT_STOP, BT_FORMULA_START)
	XML_STRING_REGEXP = compileRegexp(SQ_STRING_START, DQ_STRING_START)
	FORMULA_STRING_OR_COMMENT_REGEXP = XML_STRING_OR_COMMENT_REGEXP
	JS_STRING_OR_COMMENT_REGEXP      = compileRegexp(SQ_STRING_START,
		DQ_STRING_START, SL_COMMENT_START,
		ML_COMMENT_START, ML_COMMENT_STOP, BT_FORMULA_START)
  GLSL_STRING_OR_COMMENT_REGEXP    = JS_STRING_OR_COMMENT_REGEXP
	MATH_STRING_OR_COMMENT_REGEXP = compileRegexp(SL_COMMENT_START, ML_COMMENT_START, ML_COMMENT_STOP)
	TEMPLATE_STRING_OR_COMMENT_REGEXP   = compileRegexp(SQ_STRING_START,
		DQ_STRING_START, BT_FORMULA_START, SL_COMMENT_START,
		ML_COMMENT_START, ML_COMMENT_STOP)
  CSS_STRING_REGEXP = compileRegexp(SQ_STRING_START, DQ_STRING_START)

	BRACES_GROUP       = NewGroup(BRACES_START, BRACES_STOP)
	BRACKETS_GROUP     = NewGroup(BRACKETS_START, BRACKETS_STOP)
	PARENS_GROUP       = NewGroup(PARENS_START, PARENS_STOP)
	PARENS_OPEN_REGEXP = compileRegexp(PARENS_START)

	SQ_STRING_GROUP   = NewGroup(SQ_STRING_START, SQ_STRING_STOP)
	DQ_STRING_GROUP   = NewGroup(DQ_STRING_START, DQ_STRING_STOP)
	BT_FORMULA_GROUP  = NewGroup(BT_FORMULA_START, BT_FORMULA_STOP)
	SL_COMMENT_GROUP  = NewGroup(SL_COMMENT_START, "\n")
	ML_COMMENT_GROUP  = NewGroup(ML_COMMENT_START, ML_COMMENT_STOP)
	XML_COMMENT_GROUP = NewGroup(XML_COMMENT_START, XML_COMMENT_STOP)
  HTML_SCRIPT_GROUP = NewGroup(HTML_SCRIPT_START, HTML_SCRIPT_STOP)
  HTML_STYLE_GROUP = NewGroup(HTML_STYLE_START, HTML_STYLE_STOP)

	ALPHABET_REGEXP       = regexp.MustCompile(`\b[a-zA-Z]*\b`)
	SIMPLE_WORD_REGEXP    = regexp.MustCompile(`^[a-zA-Z]*$`)
	SVGPATH_LETTER_REGEXP = regexp.MustCompile(`[aAcClhHLmMqQsStTvVzZ]`)
	NL_REGEXP             = regexp.MustCompile(ESCAPED_NL)
	PERIOD_REGEXP         = compileRegexp(PERIOD)
	SVGPATH_MINUS_REGEXP  = regexp.MustCompile(`([^e])([\-])`)
	COMMA_REGEXP          = regexp.MustCompile(`[,]`)
	SQ_REGEXP             = compileRegexp(SQ_STRING_START)
	DQ_REGEXP             = compileRegexp(DQ_STRING_START)

	DIGIT_REGEXP = regexp.MustCompile(`[0-9]`)
	INT_REGEXP   = regexp.MustCompile(`^[\-]?[0-9]+$`)
	HEX_REGEXP   = regexp.MustCompile(`0x[0-9a-fA-F]+$`)
	FLOAT_REGEXP = regexp.MustCompile(`^[\-]?[0-9]+(\.[0-9]+)?(e[\-+]?[0-9]+)?([a-zA-Z%]*)?$`) // includes units
	FLOAT_UNITS  = []string{"s", "Q", "%", "cm", "mm", "in", "pc", "pt", "px", "em", "ch",
		"fr", "lh", "vw", "vh", "rem", "vmin", "vmax", "n"}
	PLAIN_FLOAT_REGEXP = regexp.MustCompile(`^[\-]?[0-9]+(\.[0-9]+)?(e[\-+]?[0-9]+)?$`) // doesnt include units

	TAG_START_REGEXP       = compileRegexp(TAG_START)
	TAG_NAME_REGEXP        = regexp.MustCompile(`(!\-\-)|([!?]?[#_a-zA-Z][0-9A-Za-z_\.]*)`)
	TAG_STOP_REGEXP        = regexp.MustCompile(`[/]?>`)
	XML_HEADER_STOP_REGEXP = regexp.MustCompile(`[?]>`)
	XML_COMMENT_STOP_REGEXP  = regexp.MustCompile(`-->`)
	DUMMY_TAG_NAME_REGEXP  = regexp.MustCompile(`^[\s]*>`)

	NAMESPACE_SEPARATOR_REGEXP = compileRegexp(NAMESPACE_SEPARATOR)
	XML_SYMBOLS_REGEXP        = regexp.MustCompile(`[=]`)
	//FORMULA_SYMBOLS_REGEXP     = regexp.MustCompile(`([=][=][=])|([<>=!:][=])|([&][&])|([|][|])|([!][!])|([?][?])|([!<>=:,;{}()[\]+*/\-?])`)
	JS_SYMBOLS_REGEXP          = regexp.MustCompile(`([>][>][>][=])|([=!][=][=])|([*][*][=])|([<][<][=])|([>][>][=])|([>][>][>])|([<>=!:+\-*/%&|^][=])|([*][*])|([&][&])|([<][<])|([>=][>])|([|][|])|([+][+])|([:][:])|([\-][\-])|([!<>=:,;{}()[\]+*/\-?%\.&|^~])`)
	MATH_SYMBOLS_REGEXP        = regexp.MustCompile(`([>][>])|([<][<])|([/][/])|([-=][>])|([!<>=~]?[=])|([{}()[\]+\-<>*/\.^_=,])`)
  GLSL_SYMBOLS_REGEXP        = regexp.MustCompile(`([+][+])|([-][-])|([&][&])|([|][|])|([<>!=*+\-][=])|([#:!<>;{}()[\]/\-\.+*=,])`)
  TEMPLATE_SYMBOLS_REGEXP          = regexp.MustCompile(`([=][=][=])|([|*~<>=!:^][=])|([&][&])|([|][|])|([!][!])|([?][?])|([!<>=:,;{}()[\]+*/\-?$@\.#])`)
  //CSS_SYMBOLS_REGEXP        = regexp.MustCompile(`([:][:])|([^$][=])|([:+>~()[\]*,=])`)
  CSS_SYMBOLS_REGEXP        = regexp.MustCompile(`([:][:])|([*|$~^][=])|([:+>~()[\]*,=])`)

	XML_WORD_REGEXP               = regexp.MustCompile(`[a-zA-Z_][0-9A-Za-z_\-.:]*\b`)
  MATH_WORD_REGEXP              = regexp.MustCompile(`[a-zA-Z_][0-9A-Za-z_\-.:]*\b`)
  TEMPLATE_WORD_REGEXP          = regexp.MustCompile(`[a-zA-Z_][0-9A-Za-z_\-.]*\b`)
	JS_WORD_REGEXP                = regexp.MustCompile(`^[a-zA-Z_][0-9A-Za-z_]*$`)
	XML_WORD_OR_LITERAL_REGEXP    = regexp.MustCompile(`[!]?[A-Za-z_]+[0-9A-Za-z_\-.:]*`)
	TEMPLATE_WORD_OR_LITERAL_REGEXP    = regexp.MustCompile(`([#][0-9a-fA-F]{8})|([#][0-9a-fA-F]{6})|([#][0-9a-fA-F]{4})|([#][0-9a-fA-F]{3})|([0-9A-Za-z_]+[0-9A-Za-z_\-%\.]*)`)
	//FORMULA_WORD_OR_LITERAL_REGEXP = regexp.MustCompile(`[!#$]?[0-9A-Za-z_]+[0-9A-Za-z_\-.%]*`)
	GLSL_WORD_REGEXP               = JS_WORD_REGEXP
  CSS_WORD_REGEXP                = regexp.MustCompile(`[a-zA-Z0-9\-#\.]+\b`)
  CSS_WORD_OR_LITERAL_REGEXP     = CSS_WORD_REGEXP

	// must match hex before number (because otherwise the '0' before the 'x' becomes a token by itself)
	JS_WORD_OR_LITERAL_REGEXP = regexp.MustCompile(`([A-Za-z_$]+[0-9A-Za-z_]*)|(0x[0-9a-fA-F]+)|([0-9]+(\.[0-9]+)?(e[\-+]?[0-9]+)?)`)
	GLSL_WORD_OR_LITERAL_REGEXP = JS_WORD_OR_LITERAL_REGEXP

	MATH_WORD_OR_LITERAL_REGEXP = regexp.MustCompile(`([A-Za-z]+[A-Za-z]*)|([0-9]+[0-9\-.]*[a-zA-Z]*)`)

	JS_STRING_TEMPLATE_START_REGEXP = regexp.MustCompile(`([$][{])`)
	JS_STRING_TEMPLATE_STOP_REGEXP  = regexp.MustCompile(`([}])`)
	// user variables can't contain namespace separators (i.e. dots)
	VALID_VAR_NAME_REGEXP = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_\-]*$`)

	JS_UNIVERSAL_CLASS_NAME_REGEXP = regexp.MustCompile(`^[A-Z][a-zA-Z]*$`)

	TEMPLATE_SUPER_REGEXP             = regexp.MustCompile(`\bsuper[\s]*[(]`)
	TEMPLATE_NAME_REGEXP              = regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)
	TEMPLATE_EXTENDS_TO_SUPER_REGEXP  = regexp.MustCompile(`\b(extends[\s]*[=].*)super`)

	CSS_PLAIN_CLASS_REGEXP = regexp.MustCompile(`^[\.][a-zA-Z_\-0-9]*$`)
)

func compileRegexp(x ...string) *regexp.Regexp {
	exp := ""
	if len(x) == 0 {
		panic("expected more than 0 arguments")
	} else if len(x) == 1 {
		exp = regexp.QuoteMeta(x[0])
	} else {
		for i, s := range x {
			if i > 0 {
				exp += "|"
			}

			exp += "(" + regexp.QuoteMeta(s) + ")"
		}
	}

	return regexp.MustCompile(exp)
}

func IsCompactSelfClosing(name string) bool {
  return name == "br" || name == "hr" || name == "!DOCTYPE" || name == "img" || name == "meta" || name == "input" || name == "?xml" || name == "link" || name == "base" || name == "col" || name == "param" || name == "source" || name == "track" || name == "wbr"
}

func IsSelfClosing(name string, close string) bool {
	return IsCompactSelfClosing(name) || close == TAG_STOP_SELFCLOSING
}

func StartsWithDigit(s string) bool {
	return DIGIT_REGEXP.MatchString(s[0:1])
}

func EndsWithDigit(s string) bool {
	n := len(s)
	return DIGIT_REGEXP.MatchString(s[n-1 : n])
}

func IsInt(s string) bool {
	return INT_REGEXP.MatchString(s)
}

func IsHex(s string) bool {
	return HEX_REGEXP.MatchString(s)
}

// includes units
func IsFloat(s string) bool {
	return FLOAT_REGEXP.MatchString(s)
}

func IsPlainFloat(s string) bool {
	return PLAIN_FLOAT_REGEXP.MatchString(s)
}

func IsColor(s string) bool {
	return strings.HasPrefix(s, HASH)
}

func IsSimpleWord(s string) bool {
	return SIMPLE_WORD_REGEXP.MatchString(s)
}

func ExtractUnit(s string) (string, bool) {
	suffix := ""
	found := false

	for _, u := range FLOAT_UNITS {
		if strings.HasSuffix(s, u) {
			suffix = u
			found = true
		}
	}

	return suffix, found
}

func IsBool(s string) bool {
	return s == BOOL_TRUE || s == BOOL_FALSE
}

func IsXMLWord(s string) bool {
	return XML_WORD_REGEXP.MatchString(s)
}

func IsCSSWord(s string) bool {
	return CSS_WORD_REGEXP.MatchString(s)
}

func IsMathWord(s string) bool {
	return MATH_WORD_REGEXP.MatchString(s)
}

func IsTemplateWord(s string) bool {
  return TEMPLATE_WORD_REGEXP.MatchString(s)
}

func IsJSWord(s string) bool {
	return JS_WORD_REGEXP.MatchString(s) || s == "$"
}

func IsGLSLWord(s string) bool {
	return GLSL_WORD_REGEXP.MatchString(s)
}

func IsNull(s string) bool {
	return s == NULL
}

func IsValidVar(s string) bool {
	return VALID_VAR_NAME_REGEXP.MatchString(s)
}

func IsValidFun(s string) bool {
	return VALID_VAR_NAME_REGEXP.MatchString(s)
}
