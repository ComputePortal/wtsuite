package main

import (
  "fmt"
  "os"
  "path/filepath"
  "regexp"
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
  tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
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
)

type CmdArgs struct {
  inputFile string
  outputFile string
  includeDirs []string

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
	-I, --include <dir>    Append a search directory to HTMLPPPATH
	-o, --output <file>    Defaults to "a.js" if not set
  --control <file>       Optional control file
  --math-font-url <url>  Defaults to "FreeSerifMath.woff2"
  --px-per-rem <int>     Defaults to 16
  --auto-link            Convert tags to <a> automatically if they have href attribute
	-v[v[v[v...]]]         Verbosity
`)

	os.Exit(1)
}

func printMessageAndExit(msg string, printUsage bool) {
	fmt.Fprintf(os.Stderr, "\u001b[1m"+msg+"\u001b[0m\n\n")

	if printUsage {
		printUsageAndExit()
	} else {
		os.Exit(1)
	}
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	cmdArgs := CmdArgs{
		inputFile:     "",
		outputFile:    DEFAULT_OUTPUTFILE,
		includeDirs:   make([]string, 0),
		control:        "",
    mathFontURL: DEFAULT_MATHFONTURL,
    pxPerRem: DEFAULT_PX_PER_REM,
		compactOutput: false,
		verbosity:     0,
	}

	positional := make([]string, 0)

	i := 1
	n := len(os.Args)

	for i < n {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-v") {
			re := regexp.MustCompile(`^[\-][v]+$`)
			if !re.MatchString(arg) {
				printMessageAndExit("Error: bad verbosity option "+arg, true)
			} else {
				cmdArgs.verbosity += len(arg) - 1
			}
		} else if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-?", "-h", "--help":
				printUsageAndExit()
			case "-c", "--compact":
				if cmdArgs.compactOutput {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.compactOutput = true
				}
			case "-I", "--include":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					cmdArgs.includeDirs = append(cmdArgs.includeDirs, os.Args[i+1])
					i++
				}
			case "-o", "--output":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.outputFile != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.outputFile = os.Args[i+1]
					i++
				}
			case "--control":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.control != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.control = os.Args[i+1]
					i++
				}
			case "--auto-link":
				if cmdArgs.autoLink {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.autoLink = true
				}
			case "--math-font-url":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.mathFontURL != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.mathFontURL = os.Args[i+1]
					i++
				}
			case "--px-per-rem":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.pxPerRem != -1 {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
          pxPerRem, err := strconv.ParseInt(os.Args[i+1], 10, 64)
          if err != nil {
            printMessageAndExit(err.Error(), true)
          }

          if pxPerRem <= 0 {
            printMessageAndExit("Error: invalid px-per-rem value " + os.Args[i+1], true)
          }

					cmdArgs.pxPerRem = int(pxPerRem)
					i++
				}
			}
		} else {
			positional = append(positional, arg)
		}

		i++
	}

	if len(positional) != 1 {
		printMessageAndExit("Error: expected 1 positional argument", true)
	}

	cmdArgs.inputFile = positional[0]
	if err := files.AssertFile(cmdArgs.inputFile); err != nil {
		printMessageAndExit("Error: input file '"+cmdArgs.inputFile+"' "+err.Error(), false)
	}

	inputFile, err := filepath.Abs(cmdArgs.inputFile)
	if err != nil {
		printMessageAndExit("Error: input file '"+cmdArgs.inputFile+"' "+err.Error(), false)
	} else {
		cmdArgs.inputFile = inputFile
	}

	outputFile, err := filepath.Abs(cmdArgs.outputFile)
	if err != nil {
		printMessageAndExit("Error: output file '"+cmdArgs.outputFile+"' "+err.Error(), false)
	} else {
		cmdArgs.outputFile = outputFile
	}

	for _, includeDir := range cmdArgs.includeDirs {
		if err := files.AssertDir(includeDir); err != nil {
			printMessageAndExit("Error: include dir '"+includeDir+"' "+err.Error(), false)
		}
	}

  if cmdArgs.control != "" {
    if err := files.AssertFile(cmdArgs.control); err != nil {
      printMessageAndExit("Error: control file '"+cmdArgs.control+"' "+err.Error(), false)
    }

    control, err := filepath.Abs(cmdArgs.control)
    if err != nil {
      printMessageAndExit("Error: control file '"+cmdArgs.control+"' "+err.Error(), false)
    } else {
      cmdArgs.control = control
    }
  }

  return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs) {
  if cmdArgs.compactOutput {
		tree.NL = ""
		tree.TAB = ""
		tree.COMPRESS_NUMBERS = true

		styles.NL = ""
		styles.TAB = ""
		styles.LAST_SEMICOLON = ""

    js.NL = ""
    js.TAB = ""
    js.COMPACT_NAMING = true
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

	files.AppendIncludeDirs(cmdArgs.includeDirs)
}

func buildHTMLFile(fileSource files.Source, c *directives.FileCache, src string, dst string, control string, compactOutput bool) error {
  r, cssBundleRules, err := directives.NewRoot(fileSource, c, src, "", "", "")
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

    entryScript, err := scripts.NewInitFileScript(control, "")
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

    if err := r.IncludeControl(content); err != nil {
      return err
    }
  }

	output := r.Write("", tree.NL, tree.TAB)

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


  fileSource := files.NewDefaultUIFileSource()
  c := directives.NewFileCache()

	if err := buildHTMLFile(fileSource, c, cmdArgs.inputFile, cmdArgs.outputFile, cmdArgs.control, cmdArgs.compactOutput); err != nil {
    return err
  }

  return nil
}

func main() {
  cmdArgs := parseArgs()

  setUpEnv(cmdArgs)

  if err := buildFile(cmdArgs); err != nil {
    printSyntaxErrorAndExit(err)
  }
}
