package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
  "sort"
	"strings"

	"../../pkg/cache"
	"../../pkg/directives"
	"../../pkg/files"
	"../../pkg/parsers"
	"../../pkg/tokens/context"
	tokens "../../pkg/tokens/html"
	"../../pkg/tokens/js"
	"../../pkg/tokens/js/macros"
	"../../pkg/tokens/js/values"
	"../../pkg/tree"
	"../../pkg/tree/scripts"
	"../../pkg/tree/styles"

	"./config"
)

var GitCommit string
var VERBOSITY = 0

type CmdArgs struct {
	config.CmdArgs // common for this transpiler and wt-search-index

	compactOutput bool
	forceBuild    bool
	noAliasing    bool
	autoHref      bool
	profFile      string
	xml           bool

	verbosity int // defaults to zero, every -v[v[v]] adds a level
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] configFile outputDir\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `	
Options:
  -?, -h, --help                Show this message
  -I, --include <dir>           Append a search directory to the HTMLPPPATH
  -c, --compact                 Compact output with minimal whitespace, newline etc.
  -f, --force                   Force a complete project build
	--auto-href                   Convert tags to <a> automatically if they have the 'href' attribute
  --no-aliasing                 Don't allow standard html tags to be aliased
  -D<name> <value>              Define a global variable with a value
  -B<name>                      Define a global flag (its value will be empty string though)
  --prof<file>                  Profile the transpiler, output written to file (analyzeable with go tool pprof)
	--js-url                      Override js-url in config
	--css-url                     Override css-url in config
	-i, --include-view            Include view group or view file. Cannot be combined with -x
	-x, --exclude-view            Exclude view group or view file. Cannot be combined with -i
	-j, --include-control         Include control group or control file. Cannot be combined with -y
	-y, --exclude-control         Exclude control group or control file. Cannot be combined with -j
	--math-font-url               Math font url (font name is always FreeSerifMath)
	-v[v[v[v...]]]                Verbosity
`)

	os.Exit(1)
}

func printMessageAndExit(msg string, printUsage bool) {
	config.PrintMessage(msg)
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
	// default args
	cmdArgs := CmdArgs{
		CmdArgs: config.NewDefaultCmdArgs(),

		compactOutput: false,
		forceBuild:    false,
		noAliasing:    false,
		autoHref:      false,
		profFile:      "",
		xml:           false,
		verbosity:     0,
	}

	positional := make([]string, 0)

	i := 1
	n := len(os.Args)

	for i < n {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "-D") {
			if i == n-1 {
				printMessageAndExit("Error: expected argument after "+arg, true)
			} else if len(arg) < 3 {
				printMessageAndExit("Error: expected -D<name>, not just -D", true)
			} else {
				name := arg[2:]
				value := os.Args[i+1]

				if _, ok := cmdArgs.GlobalVars[name]; ok {
					printMessageAndExit("Error: global var "+name+" already defined", true)
				} else {
					cmdArgs.GlobalVars[name] = value
					i++
				}
			}
		} else if strings.HasPrefix(arg, "-B") {
			if len(arg) < 3 {
				printMessageAndExit("Error: expected -B<name>, not just -B", true)
			} else {
				name := arg[2:]

				if _, ok := cmdArgs.GlobalVars[name]; ok {
					printMessageAndExit("Error: global var "+name+" already defined", true)
				} else {
					cmdArgs.GlobalVars[name] = ""
				}
			}
		} else if strings.HasPrefix(arg, "-v") {
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
			case "-f", "--force":
				if cmdArgs.forceBuild {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.forceBuild = true
				}
			case "-I", "--include":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else {
					cmdArgs.IncludeDirs = append(cmdArgs.IncludeDirs, os.Args[i+1])
					i++
				}
			case "--no-aliasing":
				if cmdArgs.noAliasing {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.noAliasing = true
				}
			case "--auto-href":
				if cmdArgs.autoHref {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.autoHref = true
				}
			case "--xml":
				if cmdArgs.xml {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.xml = true
				}
			case "--js-url":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)

				} else if cmdArgs.JsUrl != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.JsUrl = os.Args[i+1]
					i++
				}
			case "--css-url":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)

				} else if cmdArgs.CssUrl != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.CssUrl = os.Args[i+1]
					i++
				}
			case "-i", "--include-view":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)

				} else if len(cmdArgs.ExcludeViews) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --exclude-view (-x)", true)
				} else {
					cmdArgs.IncludeViews = append(cmdArgs.IncludeViews, os.Args[i+1])
					i++
				}
			case "-x", "--exclude-view":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if len(cmdArgs.IncludeViews) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --include-view (-i)", true)
				} else {
					cmdArgs.ExcludeViews = append(cmdArgs.ExcludeViews, os.Args[i+1])
					i++
				}
			case "-j", "--include-control":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)

				} else if len(cmdArgs.ExcludeControls) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --exclude-control (-y)", true)
				} else {
					cmdArgs.IncludeControls = append(cmdArgs.IncludeControls, os.Args[i+1])
					i++
				}
			case "-y", "--exclude-control":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if len(cmdArgs.IncludeControls) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --include-control (-j)", true)
				} else {
					cmdArgs.ExcludeControls = append(cmdArgs.ExcludeControls, os.Args[i+1])
					i++
				}
			case "--prof":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.profFile != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.profFile = os.Args[i+1]
					i++
				}
			case "--math-font-url":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if cmdArgs.MathFontUrl != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.MathFontUrl = os.Args[i+1]
					i++
				}
			default:
				printMessageAndExit("Error: unrecognized flag "+arg, true)
			}
		} else {
			positional = append(positional, arg)
		}

		i++
	}

	if len(positional) != 2 {
		printMessageAndExit("Error: expected 2 positional arguments", true)
	}

	if files.IsFile(positional[0]) {
		if files.IsFile(positional[1]) {
			printMessageAndExit("Error: got two files, expected 1 dir and 1 file", true)
		}
		cmdArgs.ConfigFile = positional[0]
		cmdArgs.OutputDir = positional[1]

	} else if files.IsDir(positional[0]) {
		if files.IsDir(positional[1]) {
			printMessageAndExit("Error: got two dirs, expected 1 dir and 1 file", true)
		}

		cmdArgs.ConfigFile = positional[1]
		cmdArgs.OutputDir = positional[0]
	} else {
		printMessageAndExit("Error: '"+positional[0]+"' or '"+positional[1]+"' doesn't exist", true)
	}

	for _, includeDir := range cmdArgs.IncludeDirs {
		if err := files.AssertDir(includeDir); err != nil {
			printMessageAndExit("Error: include dir '"+includeDir+"' "+err.Error(), true)
		}
	}

	if err := files.AssertFile(cmdArgs.ConfigFile); err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.ConfigFile+"' "+err.Error(), true)
	}

	configFile, err := filepath.Abs(cmdArgs.ConfigFile)
	if err != nil {
		printMessageAndExit("Error: configFile '"+cmdArgs.ConfigFile+"' "+err.Error(), true)
	} else {
		cmdArgs.ConfigFile = configFile
	}

	if err := files.AssertDir(cmdArgs.OutputDir); err != nil {
		printMessageAndExit("Error: output dir '"+cmdArgs.OutputDir+"' "+err.Error(), true)
	}

	absOutputDir, err := filepath.Abs(cmdArgs.OutputDir)
	if err != nil {
		panic(err)
	}

	/*pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if absOutputDir == pwd {
		printMessageAndExit("Error: output dir is same as current dir (must be different)", true)
	}*/

	cmdArgs.OutputDir = absOutputDir

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs, cfg *config.Config) {
	if cmdArgs.compactOutput {
		tree.NL = ""
		tree.TAB = ""
		tree.COMPRESS_NUMBERS = true

		styles.NL = ""
		styles.TAB = ""
		styles.LAST_SEMICOLON = ""
	}

	if cmdArgs.noAliasing {
		directives.NO_ALIASING = true
	}

	if cmdArgs.xml {
		files.XML_SYNTAX = true
	}

	if cmdArgs.autoHref {
		tree.AUTO_HREF = true
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
}

func buildHTMLFile(src, url, dst string, control string, cssUrl string, jsUrl string) error {
	cache.StartRootUpdate(src)

	directives.SetActiveURL(url)

	// must come before AddViewControl
	r, cssBundleRules, err := directives.NewRoot(src, url, control, cssUrl, jsUrl)

	directives.UnsetActiveURL()

	if err != nil {
		return err
	}

	// update the cache with the cssBundleRules
	for _, rules := range cssBundleRules {
		cache.AddCssEntry(rules, src)
	}

	output := r.Write("", tree.NL, tree.TAB)

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

	for _, src := range updatedViews {
		dst := cfg.GetViews()[src]

		control, ok := viewControls[src]
		if !ok {
			panic("should be present")
		}

		url := dst[len(cmdArgs.OutputDir):]

		err := buildHTMLFile(src, url, dst, control, cfg.CssUrl, cfg.JsUrl)

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
		if cmdArgs.compactOutput {
			js.NL = ""
			js.TAB = ""
			js.COMPACT_NAMING = true
			macros.COMPACT = true
		}

		js.TARGET = "browser"

		bundle := scripts.NewFileBundle(cmdArgs.GlobalVars)

		for _, control := range allControls {
      // each control acts as a separate entry point
      // so the cache differs from the js-project Cache
			cache.AddControl(control)

      // files.StartCacheUpdate() called internally when creating new ControlFileScript
			controlScript, err := scripts.NewControlFileScript(control, "")
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
		printMessageAndExit(err.Error()+"\n", false)
	}

	setUpEnv(cmdArgs, cfg)

	if cmdArgs.profFile != "" {
		fProf, err := os.Create(cmdArgs.profFile)
		if err != nil {
			printMessageAndExit(err.Error(), false)
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
			printMessageAndExit(err.Error(), false)
		}

		pprof.WriteHeapProfile(fMem)
		fMem.Close()
	}
}
