package main

import (
  "fmt"
  "os"
  "strconv"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
  tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
	"github.com/computeportal/wtsuite/pkg/tree/styles"
)

const (
  DEFAULT_OUTPUTFILE = "a.html"
  DEFAULT_MATHFONTURL = "FreeSerifMath.woff2" // mathfont is always included in the result
  DEFAULT_PX_PER_REM = 16
)

var (
  VERBOSITY = 0
  cmdParser *parsers.CLIParser = nil
)

type CmdArgs struct {
  inputFile string
  outputFile string

  control string // optional control to be built along with view
  mathFontURL string
  pxPerRem int
  autoLink bool

  // stylesheets and js is included inline

  compactOutput bool
  verbosity int
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `
Options:
  -?, -h, --help         Show this message, other options are ignored
	-c, --compact          Compact output with minimal whitespace and short names
	-o, --output <file>    Defaults to "a.js" if not set
  --control <file>       Optional control file
  --math-font-url <url>  Defaults to "FreeSerifMath.woff2"
  --px-per-rem <int>     Defaults to 16
  --auto-link            Convert tags to <a> automatically if they have href attribute
	-v[v[v[v...]]]         Verbosity
`)

	os.Exit(1)
}

func printMessageAndExit(msg string) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")
  os.Exit(1)
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		inputFile:     "",
		outputFile:    DEFAULT_OUTPUTFILE,
		control:        "",
    mathFontURL: DEFAULT_MATHFONTURL,
    pxPerRem: DEFAULT_PX_PER_REM,
    autoLink: false,
		compactOutput: false,
		verbosity:     0,
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("c", "compact", "-c, --compact          Compact output with minimal whitespace and short names", &(cmdArgs.compactOutput)),
      parsers.NewCLIUniqueFlag("", "auto-link", "--auto-link            Convert tags to <a> automatically if they have href attribute", &(cmdArgs.autoLink)),
  
      parsers.NewCLIUniqueFile("o", "output", "-o, --output <file>    Defaults to \"" + DEFAULT_OUTPUTFILE + "\" if not set", false, &(cmdArgs.outputFile)),
      parsers.NewCLIUniqueFile("", "control", "--control <file>    Optional control file", true, &(cmdArgs.control)),
      parsers.NewCLIUniqueString("", "math-font-url", "--math-font-url <url>  Defaults to \"" + DEFAULT_MATHFONTURL + "\"", &(cmdArgs.mathFontURL)),
      parsers.NewCLIUniqueInt("", "px-per-rem", "--px-per-rem <int>     Defaults to " + strconv.Itoa(DEFAULT_PX_PER_REM), &(cmdArgs.pxPerRem)),
      parsers.NewCLICountFlag("v", "", "Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIFile("", "", "", true, &(cmdArgs.inputFile)),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

  if cmdArgs.pxPerRem <= 0 {
    printMessageAndExit("Error: invalid px-per-rem value " + strconv.Itoa(cmdArgs.pxPerRem))
  }

  return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) error {
  if cmdArgs.compactOutput {
		tree.COMPRESS_NUMBERS = true
		patterns.NL = ""
		patterns.TAB = ""
		patterns.LAST_SEMICOLON = ""
    patterns.COMPACT_NAMING = true
    macros.COMPACT = true
  }

	if cmdArgs.autoLink {
		tree.AUTO_LINK = true
	}

	if cmdArgs.pxPerRem > 0 {
		tokens.PX_PER_REM = cmdArgs.pxPerRem
	}

  styles.MATH_FONT = "FreeSerifMath"
  styles.MATH_FONT_FAMILY = "FreeSerifMath, FreeSerif" // keep original FreeSerif as backup
  styles.MATH_FONT_URL = cmdArgs.mathFontURL

  js.TARGET = "browser"

	VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

  return files.ResolvePackages(cmdArgs.inputFile)
}

func buildHTMLFile(c *directives.FileCache, src string, dst string, control string, compactOutput bool) error {
  r, cssBundleRules, err := directives.NewRoot(c, src, "", "", "")
  if err != nil {
    return err
  }

	// update the cache with the cssBundleRules
	for _, rules := range cssBundleRules { // added to file later
		cache.AddCssEntry(rules, src)
	}

  cssContent := cache.WriteCSSBundle(styles.MATH_FONT_URL)

  // add to the Root
  if err := r.IncludeStyle(cssContent); err != nil {
    return err
  }

  if control != "" {
    files.JS_MODE = true

    cache.LoadJSCache("", true)

    entryScript, err := scripts.NewInitFileScript(control)
    if err != nil {
      return err
    }

		bundle := scripts.NewFileBundle(map[string]string{})

		bundle.Append(entryScript)

		if err := bundle.Finalize(); err != nil {
			return err
		}

		content, err := bundle.Write()
		if err != nil {
			return err
		}

    fmt.Println("including content: ", content)
    if err := r.IncludeControl(content); err != nil {
      return err
    }
  }

	output := r.Write("", patterns.NL, patterns.TAB)

	// src is just for info
	if err := files.WriteFile(src, dst, []byte(output)); err != nil {
		return err
	}

  return nil
}

func buildFile(cmdArgs CmdArgs) error {
  views := make(map[string]string)
  views[cmdArgs.inputFile] = cmdArgs.outputFile

  viewControls := make(map[string]string)
  if cmdArgs.control != "" {
    viewControls[cmdArgs.inputFile] = cmdArgs.control
  }

  // stick as close to the way it is done in wt-site as possible
	cache.LoadHTMLCache(views, viewControls,
		"", "", cmdArgs.pxPerRem, cmdArgs.outputFile, "",
		cmdArgs.compactOutput, make(map[string]string), true)


  c := directives.NewFileCache()

	if err := buildHTMLFile(c, cmdArgs.inputFile, cmdArgs.outputFile, cmdArgs.control, cmdArgs.compactOutput); err != nil {
    return err
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()

  if err := setUpEnv(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }

  if err := buildFile(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
