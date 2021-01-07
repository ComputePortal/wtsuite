package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
  "sort"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/js/macros"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
	"github.com/computeportal/wtsuite/pkg/tree/styles"

	"github.com/computeportal/wtsuite/cmd/wt-site/config"
)

var (
  GitCommit string
  VERBOSITY = 0
  cmdParser *parsers.CLIParser = nil
)

type CmdArgs struct {
	config.CmdArgs // common for this transpiler and wt-search-index

	compactOutput bool
	forceBuild    bool
	noAliasing    bool
	autoLink      bool
	profFile      string

	verbosity int // defaults to zero, every -v[v[v]] adds a level
}

func printMessageAndExit(msg string) {
	config.PrintMessage(msg)
  os.Exit(1)
}

func printSyntaxErrorAndExit(err error) {
	os.Stderr.WriteString(err.Error())
	os.Exit(1)
}

func parseArgs() CmdArgs {
	// default args
	cmdArgs := CmdArgs{
		CmdArgs: config.NewDefaultCmdArgs(),

		compactOutput: false,
		forceBuild:    false,
		noAliasing:    false,
		autoLink:      false,
		profFile:      "",
		verbosity:     0,
	}

	var positional []string = nil

  cmdParser = parsers.NewCLIParser(
    fmt.Sprintf("Usage: %s [options] <config-file> <output-dir>\n", os.Args[0]),
    "",
    []parsers.CLIOption{
      parsers.NewCLIUniqueFlag("c", "compact", "-c, --compact                 Compact output with minimal whitespace, newline etc.", &(cmdArgs.compactOutput)),
      parsers.NewCLIUniqueFlag("f", "force", "-f, --force                   Force a complete project build", &(cmdArgs.forceBuild)),
      parsers.NewCLIUniqueFlag("", "auto-link", "--auto-link                   Convert tags to <a> automatically if they have the 'href' attribute", &(cmdArgs.autoLink)), 
      parsers.NewCLIUniqueFlag("", "no-aliasing", "--no-aliasing                 Don't allow standard html tags to be aliased", &(cmdArgs.noAliasing)),
      parsers.NewCLIUniqueKeyValue("D", "-D<name> <value>              Define a global variable with a value", cmdArgs.GlobalVars),
      parsers.NewCLIUniqueKey("B", "-B<name>               Define a global flag (its value is an empty string)", cmdArgs.GlobalVars),
      parsers.NewCLIUniqueFile("", "prof", "--prof<file>                  Profile the transpiler, output written to file (analyzeable with go tool pprof)", false, &(cmdArgs.profFile)),
      parsers.NewCLIUniqueString("", "js-url", "--js-url                      Override js-url in config", &(cmdArgs.JsUrl)),
      parsers.NewCLIUniqueString("", "css-url", "--css-url                      Override css-url in config", &(cmdArgs.CssUrl)),
      parsers.NewCLIAppendString("i", "include-view", "-i, --include-view <view-group>|<view-file>   Can't be combined with -x", &(cmdArgs.IncludeViews)),
      parsers.NewCLIAppendString("x", "exclude-view", "-x, --exclude-view <view-group>|<view-file>   Can't be combined with -i", &(cmdArgs.ExcludeViews)),
      parsers.NewCLIAppendString("j", "include-control", "-j, --include-control <control-group>|<control-file>   Can't be combined with -y", &(cmdArgs.IncludeControls)),
      parsers.NewCLIAppendString("y", "exclude-control", "-y, --exclude-control <control-group>|<control-file>   Can't be combined with -j", &(cmdArgs.ExcludeControls)),
      parsers.NewCLIUniqueString("", "math-font-url", "--math-font-url                      Math font url (font name is always FreeSerifMath)", &(cmdArgs.MathFontUrl)),
      parsers.NewCLICountFlag("-v", "", "Verbosity", &(cmdArgs.verbosity)),
    },
    parsers.NewCLIRemaining(&positional),
  )

  if err := cmdParser.Parse(os.Args[1:]); err != nil {
    printMessageAndExit(err.Error())
  }

  if len(cmdArgs.IncludeViews) != 0 && len(cmdArgs.ExcludeViews) != 0 {
    printMessageAndExit("Error: --include-view can't be combined with --exclude-view")
  }

  if len(cmdArgs.IncludeControls) != 0 && len(cmdArgs.ExcludeControls) != 0 {
    printMessageAndExit("Error: --include-control can't be combined with --exclude-control")
  }

	if len(positional) != 2 {
		printMessageAndExit("Error: expected 2 positional arguments")
	}

	if files.IsFile(positional[0]) {
		if files.IsFile(positional[1]) {
			printMessageAndExit("Error: got two files, expected 1 dir and 1 file")
		}
		cmdArgs.ConfigFile = positional[0]
		cmdArgs.OutputDir = positional[1]

	} else if files.IsDir(positional[0]) {
		if files.IsDir(positional[1]) {
			printMessageAndExit("Error: got two dirs, expected 1 dir and 1 file")
		}

		cmdArgs.ConfigFile = positional[1]
		cmdArgs.OutputDir = positional[0]
	} else {
		printMessageAndExit("Error: '"+positional[0]+"' or '"+positional[1]+"' doesn't exist")
	}

	if err := files.AssertFile(cmdArgs.ConfigFile); err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.ConfigFile+"' "+err.Error())
	}

	configFile, err := filepath.Abs(cmdArgs.ConfigFile)
	if err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.ConfigFile+"' "+err.Error())
	} else {
		cmdArgs.ConfigFile = configFile
	}

	if err := files.AssertDir(cmdArgs.OutputDir); err != nil {
		printMessageAndExit("Error: output dir '"+cmdArgs.OutputDir+"' "+err.Error())
	}

	absOutputDir, err := filepath.Abs(cmdArgs.OutputDir)
	if err != nil {
		panic(err)
	}

	cmdArgs.OutputDir = absOutputDir

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs, cfg *config.Config) error {
	if cmdArgs.compactOutput {
		patterns.NL = ""
		patterns.TAB = ""
		patterns.LAST_SEMICOLON = ""
    patterns.COMPACT_NAMING = true
    macros.COMPACT = true
		tree.COMPRESS_NUMBERS = true
	}

	if cmdArgs.noAliasing {
		directives.NO_ALIASING = true
	}

	if cmdArgs.autoLink {
		tree.AUTO_LINK = true
	}

	if cfg.PxPerRem != 0 {
		tokens.PX_PER_REM = cfg.PxPerRem
	}

	for k, v := range cmdArgs.GlobalVars {
		directives.RegisterDefine(k, v)
	}

	VERBOSITY = cmdArgs.verbosity
	directives.VERBOSITY = cmdArgs.verbosity
	tokens.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	tree.VERBOSITY = cmdArgs.verbosity
	styles.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity

  return files.ResolvePackages(cmdArgs.ConfigFile)
}

func buildHTMLFile(c *directives.FileCache, src, url, dst string, control string, cssUrl string, jsUrl string) error {
	cache.StartRootUpdate(src)

	directives.SetActiveURL(url)

	// must come before AddViewControl
	r, cssBundleRules, err := directives.NewRoot(c, src, control, cssUrl, jsUrl)

	directives.UnsetActiveURL()

	if err != nil {
		return err
	}

	// update the cache with the cssBundleRules
	for _, rules := range cssBundleRules {
		cache.AddCssEntry(rules, src)
	}

	output := r.Write("", patterns.NL, patterns.TAB)

	// src is just for info
	if err := files.WriteFile(src, dst, []byte(output)); err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// src is just for info
	if err := files.WriteFile(src, dst, content); err != nil {
		return err
	}

	return nil
}

func buildProjectFiles(cfg *config.Config, cmdArgs CmdArgs) error {
	cache.LoadFileCache(cfg.Files, cmdArgs.OutputDir, cmdArgs.forceBuild)

	anyUpdated := false
	for src, dst := range cfg.Files {
		if cache.RequiresUpdate(src) {
			anyUpdated = true

			files.StartCacheUpdate(src)
			if err := copyFile(src, dst); err != nil {
				return err
			}
		}
	}

	if anyUpdated {
		cache.SaveFileCache(cmdArgs.OutputDir)
	}

	return nil
}

func buildProjectViews(cfg *config.Config, cmdArgs CmdArgs) error {
	// collect the controls for each view
	viewControls := make(map[string]string)
	for view, _ := range cfg.GetViews() {
		viewControls[view] = "" // start with no controls
	}

	for control, controlViews := range cfg.GetControls() {
		for _, view := range controlViews {
			if _, ok := viewControls[view]; !ok {
				panic("should be present")
			}

			if viewControls[view] == "" {
				viewControls[view] = control
			} else {
				panic("view can only have one control, should've been set before")
			}
		}
	}

	cache.LoadHTMLCache(cfg.GetViews(), viewControls,
		cfg.CssUrl, cfg.JsUrl, cfg.PxPerRem, cmdArgs.OutputDir, GitCommit,
		cmdArgs.compactOutput, cmdArgs.GlobalVars, cmdArgs.forceBuild)

	if cfg.MathFontUrl != "" {
		styles.MATH_FONT = "FreeSerifMath"
		styles.MATH_FONT_FAMILY = "FreeSerifMath, FreeSerif" // keep original FreeSerif as backup
		styles.MATH_FONT_URL = cfg.MathFontUrl
	}


	cache.SyncHTMLLastModifiedTimes()

  // sort views for consistent behaviour
	updatedViews := make([]string, 0)
	for src, _ := range cfg.GetViews() {
		if cache.RequiresUpdate(src) {
			updatedViews = append(updatedViews, src)
		}
	}

  sort.Strings(updatedViews)

  c := directives.NewFileCache()

	for _, src := range updatedViews {
		dst := cfg.GetViews()[src]

		control, ok := viewControls[src]
		if !ok {
			panic("should be present")
		}

		url := dst[len(cmdArgs.OutputDir):]

		err := buildHTMLFile(c, src, url, dst, control, cfg.CssUrl, cfg.JsUrl)
		if err != nil {
			context.AppendString(err, "Info: error encountered in \""+src+"\"")

			// remove src from the cache and write the cache up till that point
			cache.RollbackUpdate(src)
			cache.SaveHTMLCache(cmdArgs.OutputDir)

			return err
		}
	}

	// all views, not just updated views
	for src, _ := range cfg.GetViews() {
		control, ok := viewControls[src]
		if !ok {
			panic("should be present")
		}

		// so the cache invalidates if the control changes next time
		cache.SetViewControl(src, control)
	}

	if len(updatedViews) > 0 {
		cache.SaveHTMLCache(cmdArgs.OutputDir) // also cleans

		if VERBOSITY >= 2 {
			fmt.Println("writing css bundle", cfg.GetCssDst())
		}

		cache.SaveCSSBundle(cfg.GetCssDst(), cfg.MathFontUrl, cfg.GetMathFontDst())
	}

	return nil
}

func buildProjectControls(cfg *config.Config, cmdArgs CmdArgs) error {
	allControls := make([]string, 0)
	for control, _ := range cfg.GetControls() { // we don't need the info of which views are handled by which controls here
    allControls = append(allControls, control)
	}

  sort.Strings(allControls)

	cache.LoadControlCache(allControls, cfg.GetJsDst(), cmdArgs.compactOutput, cmdArgs.forceBuild)

  // sort controls for consistent behaviour
  anyUpdated := false
	for _, control := range allControls { 
		if cache.RequiresUpdate(control) {
      anyUpdated = true
		}
	}

	// whole bundle is updated or none of the bundle
	if anyUpdated {
		js.TARGET = "browser"

		bundle := scripts.NewFileBundle(cmdArgs.GlobalVars)

		for _, control := range allControls {
      // each control acts as a separate entry point
      // so the cache differs from the js-project Cache
			cache.AddControl(control)

      // files.StartCacheUpdate() called internally when creating new ControlFileScript
			controlScript, err := scripts.NewControlFileScript(control)
			if err != nil {
				return err
			}

			bundle.Append(controlScript)
		}

		if err := bundle.Finalize(); err != nil {
			return err
		}

		content, err := bundle.Write()
		if err != nil {
			return err
		}

		if VERBOSITY >= 2 {
			fmt.Fprintf(os.Stdout, "writing js bundle %s\n", cfg.GetJsDst())
		}

		if err := ioutil.WriteFile(cfg.GetJsDst(), []byte(content), 0644); err != nil {
			return errors.New("Error: " + err.Error())
		}

		cache.SaveCache(cfg.GetJsDst())

		return nil
	}

	return nil
}

func buildProject(cmdArgs CmdArgs, cfg *config.Config) error {
	if err := buildProjectFiles(cfg, cmdArgs); err != nil {
		return err
	}

	if err := buildProjectViews(cfg, cmdArgs); err != nil {
		return err
	}

	files.JS_MODE = true

	if err := buildProjectControls(cfg, cmdArgs); err != nil {
		return err
	}

	return nil
}

func main() {
	cmdArgs := parseArgs()

	// age of the configFile doesn't matter
	cfg, err := config.ReadConfigFile(&(cmdArgs.CmdArgs))
	if err != nil {
		printMessageAndExit(err.Error()+"\n")
	}

  if err := setUpEnv(cmdArgs, cfg); err != nil {
		printMessageAndExit(err.Error()+"\n")
  }

	if cmdArgs.profFile != "" {
		fProf, err := os.Create(cmdArgs.profFile)
		if err != nil {
			printMessageAndExit(err.Error())
		}

		pprof.StartCPUProfile(fProf)
		defer pprof.StopCPUProfile()
	}

	if err := buildProject(cmdArgs, cfg); err != nil {
		printSyntaxErrorAndExit(err)
	}

	if cmdArgs.profFile != "" {
		fMem, err := os.Create(cmdArgs.profFile + ".mprof")
		if err != nil {
			printMessageAndExit(err.Error())
		}

		pprof.WriteHeapProfile(fMem)
		fMem.Close()
	}
}
