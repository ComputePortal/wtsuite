package main

import (
  "fmt"
  "os"
  "strconv"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/git"
	"github.com/computeportal/wtsuite/pkg/parsers"
  tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
	//"github.com/computeportal/wtsuite/pkg/tree/styles"
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
  autoDownload bool

  // stylesheets and js is included inline

  compactOutput bool
  verbosity int
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
    autoDownload: false,
		compactOutput: false,
		verbosity:     0,
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("c", "compact"       , "-c, --compact          Compact output with minimal whitespace and short names", &(cmdArgs.compactOutput)),
      parsers.NewCLIUniqueFlag("", "auto-link"      , "--auto-link            Convert tags to <a> automatically if they have href attribute", &(cmdArgs.autoLink)),
      parsers.NewCLIUniqueFlag("", "auto-download"         , "--auto-download                   Automatically download missing packages (use wt-pkg-sync if you want to do this manually). Doesn't update packages!", &(cmdArgs.autoDownload)), 
      parsers.NewCLIUniqueFile("o", "output"        , "-o, --output <file>    Defaults to \"" + DEFAULT_OUTPUTFILE + "\" if not set", false, &(cmdArgs.outputFile)),
      parsers.NewCLIUniqueFile("", "control"        , "--control <file>       Optional control file", true, &(cmdArgs.control)),
      parsers.NewCLIUniqueString("", "math-font-url", "--math-font-url <url>  Defaults to \"" + DEFAULT_MATHFONTURL + "\"", &(cmdArgs.mathFontURL)),
      parsers.NewCLIUniqueInt("", "px-per-rem"      , "--px-per-rem <int>     Defaults to " + strconv.Itoa(DEFAULT_PX_PER_REM), &(cmdArgs.pxPerRem)),
      parsers.NewCLIUniqueFlag("l", "latest"        , "-l, --latest           Ignore max semver, use latest tagged versions of dependencies", &(files.LATEST)),
      parsers.NewCLICountFlag("v", ""               , "-v[v[v..]]             Verbosity", &(cmdArgs.verbosity)),
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

  if cmdArgs.autoDownload {
    git.RegisterFetchPublicOrPrivate()
  }

	if cmdArgs.pxPerRem > 0 {
		tokens.PX_PER_REM = cmdArgs.pxPerRem
	}

	directives.ForceNewViewFileScriptRegistration(directives.NewFileCache())

  directives.MATH_FONT = "FreeSerifMath"
  directives.MATH_FONT_FAMILY = "FreeSerifMath, FreeSerif" // keep original FreeSerif as backup
  directives.MATH_FONT_URL = cmdArgs.mathFontURL

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
  r, err := directives.NewRoot(c, src, "", "", "", nil)
  if err != nil {
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
