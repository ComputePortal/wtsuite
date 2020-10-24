package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strconv"
	"strings"

	"../../pkg/cache"
	"../../pkg/directives"
	"../../pkg/files"
	"../../pkg/parsers"
	"../../pkg/tokens/context"
	tokens "../../pkg/tokens/html"
	"../../pkg/tokens/js"
	"../../pkg/tokens/js/macros"
	"../../pkg/tokens/js/prototypes"
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
	serial        bool
	noCaching     bool
	animationFile string

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
	--serial                      Don't parallelize the eval types step
	--no-caching                  Don't cache js class and function results
	--animation <animation-config>  Apply page animation scripts to a list of views (activated during browsing by pressing PrintScreen)
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
		serial:        false,
		noCaching:     false,
		animationFile: "",
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
			case "--serial":
				if cmdArgs.serial {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.serial = true
				}
			case "--no-caching":
				if cmdArgs.noCaching {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.noCaching = true
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
			case "--animation":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)

				} else if cmdArgs.animationFile != "" {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.animationFile = os.Args[i+1]
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

	if cmdArgs.animationFile != "" {
		if err := files.AssertFile(cmdArgs.animationFile); err != nil {
			printMessageAndExit("Error: animation file '"+cmdArgs.animationFile+"' "+err.Error(), true)
		}

		animationFile, err := filepath.Abs(cmdArgs.animationFile)
		if err != nil {
			printMessageAndExit("Error: animation file '"+cmdArgs.animationFile+"' "+err.Error(), true)
		} else {
			cmdArgs.animationFile = animationFile
		}
	}

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

	if cmdArgs.serial {
		js.SERIAL = true
	}

	if cmdArgs.noCaching {
		values.ALLOW_CACHING = false
	}

	VERBOSITY = cmdArgs.verbosity
	directives.VERBOSITY = cmdArgs.verbosity
	tokens.VERBOSITY = cmdArgs.verbosity
	js.VERBOSITY = cmdArgs.verbosity
	prototypes.VERBOSITY = cmdArgs.verbosity
	values.VERBOSITY = cmdArgs.verbosity
	parsers.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	tree.VERBOSITY = cmdArgs.verbosity
	styles.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity
}

func buildHTMLFile(src, url, dst string, control string, animationScenes []int, cssUrl string, jsUrl string) (*js.ViewInterface, error) {
	cache.StartRootUpdate(src)

	directives.SetActiveURL(url)

	// must come before AddViewControl
	// viewInterface cannot contain auto uids
	r, cssBundleRules, viewInterface, err := directives.NewRoot(src, url, control, cssUrl, jsUrl)

	directives.UnsetActiveURL()

	if err != nil {
		return nil, err
	}

	// update the cache with the cssBundleRules
	for _, rules := range cssBundleRules {
		cache.AddCssEntry(rules, src)
	}

	if len(animationScenes) > 0 {
		if err := r.ApplyAnimation(animationScenes); err != nil {
			return nil, err
		}
	}

	output := r.Write("", tree.NL, tree.TAB)

	// src is just for info
	if err := files.WriteFile(src, dst, []byte(output)); err != nil {
		return nil, err
	}

	return viewInterface, nil
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

func buildProjectViews(cfg *config.Config, cmdArgs CmdArgs, animationCfg *config.Animation) (map[string]*js.ViewInterface, []string, error) {
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

	viewAnimationScenes := make(map[string][]int)
	for view, _ := range cfg.GetViews() {
		// start with no scenes
		viewAnimationScenes[view] = []int{}
	}

	if animationCfg != nil {
		for view, _ := range cfg.GetViews() {
			// returns empty int list if view doesnt appear in the animationCfg
			viewAnimationScenes[view] = animationCfg.GetViewScenes(view)
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

	updatedViews := make([]string, 0)

	cache.SyncHTMLLastModifiedTimes()

	// TODO: view should be sorted
	for src, _ := range cfg.GetViews() {
		if cache.RequiresUpdate(src) {
			updatedViews = append(updatedViews, src)
		}
	}

	for _, src := range updatedViews {
		dst := cfg.GetViews()[src]

		control, ok := viewControls[src]
		if !ok {
			panic("should be present")
		}

		url := dst[len(cmdArgs.OutputDir):]

		animationScenes := viewAnimationScenes[src]

		viewInterface, err := buildHTMLFile(src, url, dst, control,
			animationScenes, cfg.CssUrl, cfg.JsUrl)

		if err != nil {
			context.AppendString(err, "Info: error encountered in \""+src+"\"")
			panic(err)

			// remove src from the cache and write the cache up till that point
			cache.RollbackUpdate(src)
			cache.SaveHTMLCache(cmdArgs.OutputDir)

			return nil, nil, err
		}

		if viewInterface != nil {
			cache.SetHTMLViewInterface(src, viewInterface)
		} else if control != "" {
			panic("viewinterface should be set if control!=''")
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

	return cache.GetHTMLViewInterfaces(), updatedViews, nil
}

func buildProjectControls(cfg *config.Config, cmdArgs CmdArgs, viewInterfaces map[string]*js.ViewInterface, updatedViews []string) error {
	cache.LoadControlCache(cfg.GetControls(), cfg.GetJsDst(), viewInterfaces, cmdArgs.compactOutput, cmdArgs.forceBuild)

	updatedControls := make([]string, 0)
	for control, _ := range cfg.GetControls() {
		if cache.RequiresUpdate(control) {
			updatedControls = append(updatedControls, control)
		}
	}

	// whole bundle is updated or none of the bundle
	// but only some of the bundle needs EvalType checking
	if len(updatedControls) > 0 {
		if cmdArgs.compactOutput {
			js.NL = ""
			js.TAB = ""
			js.COMPACT_NAMING = true
			macros.COMPACT = true
		}

		js.TARGET = "browser"

		bundle := scripts.NewFileBundle(cmdArgs.GlobalVars)

		for control, views := range cfg.GetControls() {
			cache.AddControl(control, views, viewInterfaces)

			controlScript, err := scripts.NewControlFileScript(control, "", views)
			if err != nil {
				return err
			}

			bundle.Append(controlScript)
		}

		if err := bundle.FinalizeControls(viewInterfaces, updatedViews, updatedControls); err != nil {
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

func buildAnimationControls(cmdArgs CmdArgs, viewInterfaces map[string]*js.ViewInterface, animationCfg *config.Animation) error {
	if cmdArgs.compactOutput {
		// compact output might make the load times a little better for screen casts
		js.NL = ""
		js.TAB = ""
		js.COMPACT_NAMING = true
		macros.COMPACT = true
	}

	js.TARGET = "browser"

	for i, scene := range animationCfg.Scenes {
		view := scene.View
		control := scene.Control

		controlScript, err := scripts.NewControlFileScript(control, "", []string{view})
		if err != nil {
			return err
		}

		bundle := scripts.NewFileBundle(cmdArgs.GlobalVars)
		bundle.Append(controlScript)

		if err := bundle.FinalizeControls(viewInterfaces, []string{view}, []string{control}); err != nil {
			return err
		}

		content, err := bundle.Write()
		if err != nil {
			return err
		}

		content += controlScript.Hash() + "();"

		// TODO: name prefix specified in config file
		dstName := "scene" + strconv.Itoa(i) + ".js"
		dstFile := filepath.Join(cmdArgs.OutputDir, dstName)

		if VERBOSITY >= 2 {
			fmt.Fprintf(os.Stdout, "writing js scene %s\n", dstFile)
		}

		if err := ioutil.WriteFile(dstFile, []byte(content), 0644); err != nil {
			return errors.New("Error: " + err.Error())
		}
	}

	return nil
}

func buildProject(cmdArgs CmdArgs, cfg *config.Config, animationCfg *config.Animation) error {
	if err := buildProjectFiles(cfg, cmdArgs); err != nil {
		return err
	}

	viewInterfaces, updatedViews, err := buildProjectViews(cfg, cmdArgs, animationCfg)
	if err != nil {
		return err
	}

	files.JS_MODE = true

	if VERBOSITY >= 2 {
		fmt.Println("Info: views built")
	}

	// build animation controls first, so as to catch those errors sooner
	if animationCfg != nil {
		if err := buildAnimationControls(cmdArgs, viewInterfaces, animationCfg); err != nil {
			return err
		}
	}

	if err := buildProjectControls(cfg, cmdArgs, viewInterfaces, updatedViews); err != nil {
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

	var animationCfg *config.Animation = nil
	if cmdArgs.animationFile != "" {
		animationCfg, err = config.ReadAnimationFile(cmdArgs.animationFile)
		if err != nil {
			printMessageAndExit(err.Error()+"\n", false)
		}
	}

	if cmdArgs.profFile != "" {
		fProf, err := os.Create(cmdArgs.profFile)
		if err != nil {
			printMessageAndExit(err.Error(), false)
		}

		pprof.StartCPUProfile(fProf)
		defer pprof.StopCPUProfile()
	}

	if err := buildProject(cmdArgs, cfg, animationCfg); err != nil {
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
