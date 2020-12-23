package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/computeportal/wtsuite/pkg/cache"
	"github.com/computeportal/wtsuite/pkg/directives"
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
	"github.com/computeportal/wtsuite/pkg/tree/styles"

	"github.com/computeportal/wtsuite/cmd/wt-site/config"
)

var VERBOSITY = 0

type CmdArgs struct {
	config.CmdArgs

	searchIndexOutput string
	includeIndices    []string
	excludeIndices    []string

	noAliasing bool
	autoLink   bool

	verbosity int
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] configFile searchIndexOutput\n", os.Args[0])
	fmt.Fprintf(os.Stderr, `	
Options:
  -?, -h, --help         Show this message
  -I, --include <dir>    Append a search directory to the HTMLPPPATH
	--auto-link            Convert tags to <a> automatically if they have the 'href' attribute
  --no-aliasing          Don't allow standard html tags to be aliased
  -D<name> <value>       Define a global variable with a value
  -B<name>               Define a global flag (its value will be empty string though)
	-i, --include-view     Include view group or view file. Cannot be combined with -x
	-x, --exclude-view     Exclude view group or view file. Cannot be combined with -i
	--include-index        Include index group in final search index. Cannot be combined with --exclude-index
	--exclude-index        Exclude index group from final search index. Cannot be combined with --include-index
	-v[v[v[v...]]]         Verbosity
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

		searchIndexOutput: "",
		includeIndices:    make([]string, 0),
		excludeIndices:    make([]string, 0),

		noAliasing: false,
		autoLink:   false,

		verbosity: 0,
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
			case "--auto-link":
				if cmdArgs.autoLink {
					printMessageAndExit("Error: "+arg+" already specified", true)
				} else {
					cmdArgs.autoLink = true
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
			case "--include-index":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if len(cmdArgs.excludeIndices) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --exclude-index", true)
				} else {
					cmdArgs.includeIndices = append(cmdArgs.includeIndices, os.Args[i+1])
					i++
				}
			case "--exclude-index":
				if i == n-1 {
					printMessageAndExit("Error: expected argument after "+arg, true)
				} else if len(cmdArgs.includeIndices) != 0 {
					printMessageAndExit("Error: "+arg+" cannot be combined with --include-index", true)
				} else {
					cmdArgs.excludeIndices = append(cmdArgs.excludeIndices, os.Args[i+1])
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

	if !files.IsFile(positional[0]) {
		printMessageAndExit("Error: first argument is not a file", true)
	}

	cmdArgs.ConfigFile = positional[0]
	cmdArgs.OutputDir = "/tmp/wt-site"
	cmdArgs.searchIndexOutput = positional[1]

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

	absSearchIndexOutput, err := filepath.Abs(cmdArgs.searchIndexOutput)
	if err != nil {
		panic(err)
	}

	cmdArgs.searchIndexOutput = absSearchIndexOutput

	return cmdArgs
}

func setUpEnv(cmdArgs CmdArgs, cfg *config.Config) {
	if cmdArgs.noAliasing {
		directives.NO_ALIASING = true
	}

	if cmdArgs.autoLink {
		tree.AUTO_LINK = true
	}

	// TODO: disable math parsing

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
	parsers.VERBOSITY = cmdArgs.verbosity
	files.VERBOSITY = cmdArgs.verbosity
	cache.VERBOSITY = cmdArgs.verbosity
	tree.VERBOSITY = cmdArgs.verbosity
	styles.VERBOSITY = cmdArgs.verbosity
	scripts.VERBOSITY = cmdArgs.verbosity
}

type SearchIndexPage struct {
	Url     string   `json:"url"`     // used as key
	Title   string   `json:"title"`   // should be unique for each indexed page
	Content []string `json:"content"` // each string is a paragraph
}

type SearchIndex struct {
	Pages   []SearchIndexPage      `json:"pages"`   // sorted
	Ignore  map[string]string      `json:"ignore"`  // key is same as value
	Index   map[string]interface{} `json:"index"`   // nested character tree, leaves are indices into pages array
	Partial map[string]interface{} `json:"partial"` // nested character tree which doesn't start at beginning of word
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		Pages:   make([]SearchIndexPage, 0),
		Ignore:  make(map[string]string),
		Index:   make(map[string]interface{}), // XXX: can be left empty for initial test
		Partial: make(map[string]interface{}), // XXX: can be left empty for initial test
	}
}

func (si *SearchIndex) AddPage(url string, title string, content []string) {
	si.Pages = append(si.Pages, SearchIndexPage{url, title, content})
}

func findRootParagraph(xpath []tree.Tag) tree.Tag {
	for _, t := range xpath {
		if t.Name() == "p" {
			return t
		}
	}

	return nil
}

func indexConfigKeyIncluded(cmdArgs CmdArgs, key string) bool {
	for _, included := range cmdArgs.includeIndices {
		if included == key {
			return true
		}
	}

	if len(cmdArgs.includeIndices) > 0 {
		return false
	}

	for _, excluded := range cmdArgs.excludeIndices {
		if excluded == key {
			return false
		}
	}

	return true
}

func registerSearchableContent(cmdArgs CmdArgs, cfg *config.Config) (*SearchIndex, error) {
	viewControls := make(map[string]string)
	for view, _ := range cfg.GetViews() {
		viewControls[view] = "" // no controls
	}

	viewSearchIndicesConfig := make(map[string]*config.SearchIndexConfig)
	for view, _ := range cfg.GetViews() {
		viewSearchIndicesConfig[view] = nil // no search strategy by default
	}

	for key, indexConfig := range cfg.Search.Indices {
		if indexConfigKeyIncluded(cmdArgs, key) {
			for _, view := range indexConfig.Pages {
				if prev, ok := viewSearchIndicesConfig[view]; ok && prev != nil {
					return nil, errors.New("Error: " + view + " has more than one search strategy")
				} else if !ok {
					//fmt.Println(cfg.GetViews())
					//fmt.Println(viewSearchStrategies)
					//panic("unexpected")
				}

				viewSearchIndicesConfig[view] = &indexConfig
			}
		}
	}

	cache.LoadHTMLCache(cfg.GetViews(), viewControls,
		cfg.CssUrl, cfg.JsUrl, cfg.PxPerRem, cmdArgs.OutputDir, "",
		false, cmdArgs.GlobalVars, true)

	if cfg.MathFontUrl != "" {
		styles.MATH_FONT = "FreeSerifMath"
		styles.MATH_FONT_FAMILY = "FreeSerifMath, FreeSerif" // keep original FreeSerif as backup
		styles.MATH_FONT_URL = cfg.MathFontUrl
	}

	searchIndex := NewSearchIndex()

  fileSource := files.NewDefaultUIFileSource()
  c := directives.NewFileCache()

	for src, dst := range cfg.GetViews() {
		// TODO: do something with strategy
		if indexConfig, ok := viewSearchIndicesConfig[src]; ok && indexConfig != nil {
			// TODO: only the views that are mentioned in the config file
			url := dst[len(cmdArgs.OutputDir):]

			cache.StartRootUpdate(src) // XXX: is this really needed?

			directives.SetActiveURL(url)

			r, _, err := directives.NewRoot(fileSource, c, src, "", cfg.CssUrl, cfg.JsUrl)
			if err != nil {
				return nil, err
			}

			directives.UnsetActiveURL()

			var activeParagraph tree.Tag = nil

			title := ""
			content := []string{}
			if err := tree.WalkText(r, []tree.Tag{}, func(xpath []tree.Tag, s string) error {
				// TODO: track active paragraph

				if indexConfig.TitleMatch(xpath) {
					if title != "" {
						return errors.New("Error: non-unique title match in " + src)
					}
					title = s
				} else if indexConfig.ContentMatch(xpath) {
					rootParagraph := findRootParagraph(xpath)

					if rootParagraph != nil {
						if rootParagraph == activeParagraph {
							// append to last content entry
							content[len(content)-1] += s
						} else {
							activeParagraph = rootParagraph
							content = append(content, s)
						}
					} else {
						content = append(content, s)
					}
				}

				return nil
			}); err != nil {
				panic("unexpected")
			}

			// no title is not permitted
			if title == "" {
				return nil, errors.New("Error: no title found for " + src)
			}

			// no content is plausible in some cases though
			searchIndex.AddPage(url, title, content)
		}
	}

	return searchIndex, nil
}

func (si *SearchIndex) indexWord(m map[string]interface{}, pageID int, f string) error {
	chars := strings.Split(f, "")

	for _, char := range chars {
		if mInner, ok := m[char]; ok {
			m = mInner.(map[string]interface{})
		} else {
			mInner := make(map[string]interface{})
			m[char] = mInner
			m = mInner
		}
	}

	if pages_, ok := m["pages"]; ok {
		pages := pages_.([]float64)
		unique := true
		for _, page := range pages {
			if int(page) == pageID {
				unique = false
				break
			}
		}

		if unique {
			pages = append(pages, float64(pageID))
			m["pages"] = pages
		}
	} else {
		m["pages"] = []float64{float64(pageID)}
	}

	return nil
}

func (si *SearchIndex) IndexWord(pageID int, f string) error {
	if err := si.indexWord(si.Index, pageID, f); err != nil {
		return err
	}

	// partial versions of the word are also indexed (include the full word itself)
	for i := 0; i < len(f); i++ {
		fPart := f[i:]

		if err := si.indexWord(si.Partial, pageID, fPart); err != nil {
			return err
		}
	}

	return nil
}

func isIgnoredWord(cfg *config.Config, w string) bool {
	i := sort.SearchStrings(cfg.Search.Ignore, w)

	if i > -1 && i < len(cfg.Search.Ignore) {
		return cfg.Search.Ignore[i] == w
	} else {
		return false
	}
}

func indexSentence(cfg *config.Config, si *SearchIndex, pageID int, sentence string) error {
	fields := strings.FieldsFunc(strings.Trim(sentence, "."), func(r rune) bool {
		return r < 46 || // keep period as decimal separator
			r == 47 || // forward slash
			r == 58 || // :
			r == 59 || // ;
			r == 63 || // ?
			r == 95 // _
	})

	for _, field := range fields {
		f := strings.ToLower(strings.TrimSpace(field))
		if f != "" {
			if !isIgnoredWord(cfg, f) {
				if err := si.IndexWord(pageID, f); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// actually fill the index/partial nested trees
func buildSearchIndex(cfg *config.Config, searchIndex *SearchIndex) error {
	// loop each word of each page
	for i, page := range searchIndex.Pages {
		if err := indexSentence(cfg, searchIndex, i, page.Title); err != nil {
			return err
		}
		for _, paragraph := range page.Content {
			if err := indexSentence(cfg, searchIndex, i, paragraph); err != nil {
				return err
			}
		}
	}

	// add the ignored values
	for _, w := range cfg.Search.Ignore {
		searchIndex.Ignore[w] = w
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

	searchIndex, err := registerSearchableContent(cmdArgs, cfg)
	if err != nil {
		printMessageAndExit(err.Error()+"\n", false)
	}

	if err := buildSearchIndex(cfg, searchIndex); err != nil {
		printMessageAndExit(err.Error()+"\n", false)
	}

	b, err := json.Marshal(searchIndex)
	if err != nil {
		printMessageAndExit(err.Error()+"\n", false)
	}

	if err := ioutil.WriteFile(cmdArgs.searchIndexOutput, b, 0644); err != nil {
		printMessageAndExit(err.Error()+"\n", false)
	}
}
