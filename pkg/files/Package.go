package files

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "path/filepath"
  "os"
  "strings"
)

const (
  PACKAGE_JSON = "package.json"
  USER_DIR_ENV_KEY = "WTPATH"
  USER_REL_DIR = ".local/share/wtsuite/private"

  SHARE_DIR_ENV_KEY = "WTSHARE"
  SHARE_REL_DIR = ".local/share/wtsuite/public"

  PRIVSSH_ENV_KEY = "WTSSHKEY"
  PRIVSSH_REL_PATH = ".ssh/id_rsa"
)

var (
  _packages map[string]*Package = nil
  CACHE_PACKAGES = true // the wtaas server should set this to false though
)

// json structures

type DependencyConfig struct {
  MinVersion string `json:"minVersion"` // should be semver, empty == -infty
  MaxVersion string `json:"maxVersion"` // should be semver, empty == +infty
  URL string `json:"url"` // github.com/...
}

type SSGConfig struct {
  MinVersion string `json:"minVersion"`
  MaxVersion string `json:"maxVersion"`
}

type PackageConfig struct {
  Dependencies map[string]DependencyConfig `json:"dependencies"`
  TemplateModules map[string]string `json:"templateModules"`
  ScriptModules map[string]string `json:"scriptModules"`
  ShaderModules map[string]string `json:"shaderModules"`
  SSG SSGConfig `json:"ssg"`
}

type Package struct {
  configPath string // for better error messages
  dependencies map[string]*Package
  templateModules map[string]string // resolved paths
  scriptModules map[string]string
  shaderModules map[string]string
  ssgSemVerRange *SemVerRange
}

func NewEmptyPackageConfig() *PackageConfig {
  return &PackageConfig{
    Dependencies: make(map[string]DependencyConfig),
    TemplateModules: make(map[string]string),
    ScriptModules: make(map[string]string),
    ShaderModules: make(map[string]string),
    SSG: SSGConfig{},
  }
}

// returns the directory of the installed package
type FetchFunc func(url string, svr *SemVerRange) (string, error)

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
    fmt.Fprintf(os.Stderr, "Warning: " + filepath.Join(dir, PACKAGE_JSON) + " not found\n")

    return NewEmptyPackageConfig(), dir, nil
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

func LoadPackage(dir string, canMoveUp bool, fetcher FetchFunc) (*Package, error) {
  return loadPackage(dir, canMoveUp, fetcher, []string{})
}

func loadPackage(dir string, canMoveUp bool, fetcher FetchFunc, prevDeps []string) (*Package, error) {
  cfg, fname, err := readPackageConfig(dir, canMoveUp)
  if err != nil {
    return nil, err
  }

  // TODO: detect circular dependencies
  deps := make(map[string]*Package)
  for k, depCfg := range cfg.Dependencies {
    deps[k], err = resolveDependency(depCfg, fetcher, prevDeps)
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

  var ssgMinVersion *SemVer = nil
  if cfg.SSG.MinVersion != "" {
    ssgMinVersion, err = ParseSemVer(cfg.SSG.MinVersion)
    if err != nil {
      return nil, err
    }
  }

  var ssgMaxVersion *SemVer = nil
  if cfg.SSG.MaxVersion != "" {
    ssgMaxVersion, err = ParseSemVer(cfg.SSG.MaxVersion)
    if err != nil {
      return nil, err
    }
  }

  return &Package{
    fname,
    deps,
    templateModules,
    scriptModules,
    shaderModules,
    NewSemVerRange(ssgMinVersion, ssgMaxVersion),
  }, nil
}

func PublicPkgInstallDir(url string) string {
  var base string
  if shareBase := os.Getenv(SHARE_DIR_ENV_KEY); shareBase != "" {
    base = shareBase
  } else {
    base = filepath.Join(os.Getenv("HOME"), SHARE_REL_DIR)
  }
  
  return filepath.Join(base, url)
}

func PrivatePkgInstallDir(url string) string {
  var base string
  if userBase := os.Getenv(USER_DIR_ENV_KEY); userBase != "" {
    base = userBase
  } else {
    base = filepath.Join(os.Getenv("HOME"), USER_REL_DIR)
  }

  return filepath.Join(base, "pkg", url)
}

func PkgInstallDir(url string) string {
  privateDir := PrivatePkgInstallDir(url)
  publicDir := PublicPkgInstallDir(url)

  if IsDir(publicDir) {
    return publicDir
  } else if IsDir(privateDir) {
    return privateDir
  } else {
    return publicDir
  }
}

func ReadPrivateSSHKey() (string, error) {
  var path string
  if privSSHPath := os.Getenv(PRIVSSH_ENV_KEY); privSSHPath != "" {
  
    path = privSSHPath
  } else {
    path = filepath.Join(os.Getenv("HOME"), PRIVSSH_REL_PATH)
  }

  b, err := ioutil.ReadFile(path)
  if err != nil {
    return "", err
  }

  return string(b), nil
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

func resolveDependency(depCfg DependencyConfig, fetcher FetchFunc, prevDeps []string) (*Package, error) {
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

  svr := NewSemVerRange(semVerMin, semVerMax)

  url := strings.ToLower(strings.TrimSpace(depCfg.URL))

  for _, prevURL := range prevDeps {
    if url == prevURL {
      return nil, errors.New("Error: circular dependencies (" + strings.Join(prevDeps, ", ") + ", " + url + ")")
    }
  }

  prevDeps = append(prevDeps, url)

  if err := validateURL(url); err != nil {
    return nil, err
  }

  pkgDir, err := fetcher(url, svr)
  if err != nil {
    return nil, err
  }

  if !IsDir(pkgDir) {
    if IsFile(pkgDir) {
      return nil, errors.New("Error: dependent package " + pkgDir + " is a file?\n")
    }

    if FetchPublicOrPrivate == nil {
      return nil, errors.New("Error: dependent package " + pkgDir + " not found (hint: use wt-pkg-sync)\n")
    } else {
      // auto download at least one version
      pkgDir, err = FetchPublicOrPrivate(url, svr)
      if err != nil {
        return nil, err
      }
    }
  }

  version, err := svr.FindBestVersion(pkgDir)
  if err != nil {
    return nil, err
  }

  if version == "" {
    return nil, errors.New("Error: no valid package versions found for " + pkgDir)
  }

  semVerDir := filepath.Join(pkgDir, version)
  
  pkg, err := loadPackage(semVerDir, false, fetcher, prevDeps)
  if err != nil {
    return nil, err
  }

  return pkg, nil
}

// must be called explicitly by cli tools so that packages become available for search
func resolvePackages(startFile string, fetcher FetchFunc) error {
  if _packages == nil && CACHE_PACKAGES {
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

  pkg, err := LoadPackage(dir, true, fetcher)
  if err != nil {
    return err
  }

  if CACHE_PACKAGES {
    _packages[filepath.Dir(pkg.configPath)] = pkg
  }

  return nil
}

func ResolvePackages(startFile string) error {
  return resolvePackages(startFile, func(url string, semVer *SemVerRange) (string, error) {
    return PkgInstallDir(url), nil
  })
}

func SyncPackages(startFile string, fetcher FetchFunc) error {
  if fetcher == nil {
    panic("fetcher function can't be nil")
  }

  return resolvePackages(startFile, fetcher)
}

func findPackage(callerDir string) *Package {
  if _packages != nil {
    if pkg, ok := _packages[callerDir]; ok {
      return pkg
    }
  }

  if callerDir == "/" {
    return nil
  } else {
    pkg := findPackage(filepath.Dir(callerDir))

    if CACHE_PACKAGES {
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

func (pkg *Package) SSGSemVerRange() *SemVerRange {
  return pkg.ssgSemVerRange
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
