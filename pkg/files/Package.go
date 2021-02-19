package files

import (
  "encoding/json"
  "errors"
  "io/ioutil"
  "path/filepath"
  "os"
  "strings"
)

const PACKAGE_JSON = "package.json"

const USER_ENV_KEY = "WTPATH"

const BASE_DIR = ".local/share/wtsuite/"

// json structures

type DependencyConfig struct {
  MinVersion string `json:"minVersion"` // should be semver, empty == -infty
  MaxVersion string `json:"maxVersion"` // should be semver, empty == +infty
  URL string `json:"url"` // github.com/...
}

type PackageConfig struct {
  Dependencies map[string]DependencyConfig `json:"dependencies"`
  TemplateModules map[string]string `json:"templateModules"`
  ScriptModules map[string]string `json:"scriptModules"`
  ShaderModules map[string]string `json:"shaderModules"`
}

type Package struct {
  configPath string // for better error messages
  dependencies map[string]*Package
  templateModules map[string]string // resolved paths
  scriptModules map[string]string
  shaderModules map[string]string
}

var _packages map[string]*Package = nil

func NewEmptyPackageConfig() *PackageConfig {
  return &PackageConfig{
    Dependencies: make(map[string]DependencyConfig),
    TemplateModules: make(map[string]string),
    ScriptModules: make(map[string]string),
    ShaderModules: make(map[string]string),
  }
}

type FetchFunc func(url string, svr *SemVerRange) error

// dir assumed to be abs
func findPackageConfig(dir string, canMoveUp bool) string {
  fname := filepath.Join(dir, PACKAGE_JSON)

  if IsFile(fname) {
    return fname
  } else if canMoveUp {
    if dir == "/" {
      return ""
    } else {
      return findPackageConfig(filepath.Dir(dir), canMoveUp)
    }
  } else {
    return ""
  }
}

func readPackageConfig(dir string, canMoveUp bool) (*PackageConfig, string, error) {
  fname := findPackageConfig(dir, canMoveUp)
  if fname == "" {
    return nil, "", errors.New("Error: " + filepath.Join(dir, PACKAGE_JSON) + " not found\n")
  }

	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, "", errors.New("Error: problem reading the config file\n")
	}

  cfg := NewEmptyPackageConfig()
  if err := json.Unmarshal(b, &cfg); err != nil {
    return nil, "", errors.New("Error: bad " + PACKAGE_JSON + " file syntax (" + fname + ")\n")
  }

  return cfg, fname, nil
}

func readPackage(dir string, canMoveUp bool, fetcher FetchFunc) (*Package, error) {
  cfg, fname, err := readPackageConfig(dir, canMoveUp)
  if err != nil {
    return nil, err
  }

  deps := make(map[string]*Package)
  for k, depCfg := range cfg.Dependencies {
    deps[k], err = resolveDependency(depCfg, fetcher)
    if err != nil {
      return nil, err
    }
  }

  // might differ from input dir, dur to canMoveUp
  actualDir := filepath.Dir(fname)

  fn := func(relPath string) (string, error) {
    if !strings.HasPrefix(relPath, "./") {
      return "", errors.New("Error: " + relPath + " not relative to package root (see " + fname + ")\n")
    }

    absPath := filepath.Join(actualDir, relPath)
    if !IsFile(absPath) {
      return "", errors.New("Error: file " + relPath + " not found (see " + fname + ")\n")
    }

    return absPath, nil
  }

  templateModules := make(map[string]string)
  scriptModules := make(map[string]string)
  shaderModules := make(map[string]string)

  for k, relPath := range cfg.TemplateModules {
    if templateModules[k], err = fn(relPath); err != nil {
      return nil, err
    }
  }

  for k, relPath := range cfg.ScriptModules {
    if scriptModules[k], err = fn(relPath); err != nil {
      return nil, err
    }
  }

  for k, relPath := range cfg.ShaderModules {
    if shaderModules[k], err = fn(relPath); err != nil {
      return nil, err
    }
  }

  return &Package{
    fname,
    deps,
    templateModules,
    scriptModules,
    shaderModules,
  }, nil
}

func PkgInstallDst(url string) string {
  var base string
  if userBase := os.Getenv(USER_ENV_KEY); userBase != "" {
    base = userBase
  } else {
    base = filepath.Join(os.Getenv("HOME"), BASE_DIR)
  }

  return filepath.Join(base, url)
}

func validateURL(url string) error {
  if url == "" {
    return errors.New("Error: url can't be empty\n")
  }

  if strings.HasSuffix(url, ".git") {
    return errors.New("Error: url .git suffix must be omitted\n")
  }

  if strings.Contains(url, "+=^~`;<>,|:!?'\"&@%$#*(){}[]\\") {
    return errors.New("Error: url contains invalid chars (hint: schema must be omitted)\n")
  }

  return nil
}

// adds result to _packages too
func resolveDependency(depCfg DependencyConfig, fetcher FetchFunc) (*Package, error) {
  semVerMin, err := ParseSemVer(depCfg.MinVersion)
  if err != nil {
    return nil, errors.New("Error: bad minVersion semver\n")
  }

  var semVerMax *SemVer = nil 
  if !LATEST {
    semVerMax, err = ParseSemVer(depCfg.MaxVersion)
    if err != nil {
      return nil, errors.New("Error: bad maxVersion semver\n")
    }
  }

  semVerRange := NewSemVerRange(semVerMin, semVerMax)

  if err := validateURL(depCfg.URL); err != nil {
    return nil, err
  }

  pkgDir := PkgInstallDst(depCfg.URL)

  if fetcher != nil {
    if err := fetcher(depCfg.URL, semVerRange); err != nil {
      return nil, err
    }
  }

  if !IsDir(pkgDir) {
    if fetcher != nil {
      panic("fetcher failed")
    }

    if IsFile(pkgDir) {
      return nil, errors.New("Error: dependent package " + pkgDir + " is a file?\n")
    }

    return nil, errors.New("Error: dependent package " + pkgDir + " not found (hint: use wt-pkg-sync)\n")
  }

  semVerDir, err := semVerRange.FindBestVersion(pkgDir)
  if err != nil {
    return nil, err
  }
  
  pkg, err := readPackage(semVerDir, false, fetcher)
  if err != nil {
    return nil, err
  }

  _packages[filepath.Dir(pkg.configPath)] = pkg

  return pkg, nil
}

// must be called explicitly by cli tools so that packages become available for search
func resolvePackages(startFile string, fetcher FetchFunc) error {
  if _packages == nil {
    _packages = make(map[string]*Package)
  }

  dir := startFile

  if !filepath.IsAbs(dir) {
    return errors.New("Error: start path " + dir + " isn't absolute\n")
  }

  if IsFile(dir) {
    dir = filepath.Dir(dir)
  }

  if !IsDir(dir) {
    return errors.New("Error: " + dir + " is not a directory\n")
  }

  pkg, err := readPackage(dir, true, fetcher)
  if err != nil {
    return err
  }

  //fmt.Println("Resolved package.json in ", dir, " for ", startFile, filepath.Dir(pkg.configPath)
  _packages[filepath.Dir(pkg.configPath)] = pkg

  return nil
}

func ResolvePackages(startFile string) error {
  return resolvePackages(startFile, nil)
}

func SyncPackages(startFile string, fetcher FetchFunc) error {
  if fetcher == nil {
    panic("fetcher function can't be nil")
  }

  return resolvePackages(startFile, fetcher)
}

func findPackage(callerDir string) *Package {
  if pkg, ok := _packages[callerDir]; ok {
    return pkg
  } else if callerDir == "/" {
    return nil
  } else {
    pkg := findPackage(filepath.Dir(callerDir))
    if pkg != nil {
      _packages[callerDir] = pkg
    }

    return pkg
  }
}

func (pkg *Package) GetModule(moduleName string, lang Lang) (string, bool) {
  switch lang {
  case SCRIPT:
    modulePath, ok := pkg.scriptModules[moduleName]
    return modulePath, ok
  case TEMPLATE:
    modulePath, ok := pkg.templateModules[moduleName]
    return modulePath, ok
  case SHADER:
    modulePath, ok := pkg.shaderModules[moduleName]
    return modulePath, ok
  default:
    panic("unhandled")
  }

  return "", false
}

func SearchPackage(caller string, pkgPath string, lang Lang) (string, error) {
  currentPkg := findPackage(filepath.Dir(caller))
  if currentPkg == nil {
    return "", errors.New("Error: no " + PACKAGE_JSON + " found/loaded for " + caller + "\n")
  }

  if filepath.IsAbs(pkgPath) {
    err := errors.New("Error: package path can't be absolute (" + pkgPath + ")\n")
    panic(err)
    return "", err
  }

  // first try getting module from currentPkg
  if modulePath, ok := currentPkg.GetModule(pkgPath, lang); ok {
    return modulePath, nil
  }

  pkgParts := strings.Split(filepath.ToSlash(pkgPath), "/")

  if len(pkgParts) == 0 {
    return "", errors.New("Error: unable to determine package name\n")
  }
  
  pkgName := pkgParts[0]
  moduleName := ""
  if len(pkgParts) > 1 {
    moduleName = filepath.Join(pkgParts[1:]...)
  }

  pkg, ok := currentPkg.dependencies[pkgName]
  if !ok {
    return "", errors.New("Error: " + currentPkg.configPath + " doesn't reference a dependency called " + pkgPath + "\n")
  }

  modulePath, ok := pkg.GetModule(moduleName, lang)
  if !ok {
    return "", errors.New("Error: no " + strings.ToLower(string(lang)) + " module \"" + moduleName + "\" found in " + pkg.configPath + "\n")
  }

  return modulePath, nil
}

func SearchTemplate(caller string, path string) (string, error) {
  return SearchPackage(caller, path, TEMPLATE)
}

func SearchScript(caller string, path string) (string, error) {
  return SearchPackage(caller, path, SCRIPT)
}

func SearchShader(caller string, path string) (string, error) {
  return SearchPackage(caller, path, SHADER)
}
