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
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/styles"
)

const (
  DEFAULT_OUTPUTFILE = "a.css"
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

  mathFontURL string
  pxPerRem int

  compactOutput bool
  autoDownload bool
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
    mathFontURL: DEFAULT_MATHFONTURL,
    pxPerRem: DEFAULT_PX_PER_REM,
		compactOutput: false,
    autoDownload: false,
		verbosity:     0,
	}

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s <input-file> [-o <output-file>] [options]\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("c", "compact"       , "-c, --compact          Compact output with minimal whitespace and short names", &(cmdArgs.compactOutput)),
  
      parsers.NewCLIUniqueFile("o", "output"        , "-o, --output <file>    Defaults to \"" + DEFAULT_OUTPUTFILE + "\" if not set", false, &(cmdArgs.outputFile)),
      parsers.NewCLIUniqueString("", "math-font-url", "--math-font-url <url>  Defaults to \"" + DEFAULT_MATHFONTURL + "\"", &(cmdArgs.mathFontURL)),
      parsers.NewCLIUniqueInt("", "px-per-rem"      , "--px-per-rem <int>     Defaults to " + strconv.Itoa(DEFAULT_PX_PER_REM), &(cmdArgs.pxPerRem)),
      parsers.NewCLIUniqueFlag("", "auto-download"         , "--auto-download                   Automatically download missing packages (use wt-pkg-sync if you want to do this manually). Doesn't update packages!", &(cmdArgs.autoDownload)), 
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
  }

	if cmdArgs.pxPerRem > 0 {
		tokens.PX_PER_REM = cmdArgs.pxPerRem
	}

  if cmdArgs.autoDownload {
    git.RegisterFetchPublicOrPrivate()
  }

  directives.MATH_FONT = "FreeSerifMath"
  directives.MATH_FONT_FAMILY = "FreeSerifMath, FreeSerif" // keep original FreeSerif as backup
  directives.MATH_FONT_URL = cmdArgs.mathFontURL

	VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity

  return files.ResolvePackages(cmdArgs.inputFile)
}

func buildFile(cmdArgs CmdArgs) error {
  // stick as close to the way it is done in wt-site as possible
  views := make(map[string]string)
  views[cmdArgs.inputFile] = cmdArgs.outputFile
	cache.LoadHTMLCache(views, map[string]string{},
		"", "", cmdArgs.pxPerRem, cmdArgs.outputFile, "",
		cmdArgs.compactOutput, make(map[string]string), true)

	if err := styles.BuildFile(cmdArgs.inputFile, cmdArgs.outputFile); err != nil {
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
