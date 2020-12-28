package scripts

import (
	"errors"
	"strings"

	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
)

type FileScript interface {
	Script
	ResolveNames(scope js.GlobalScope) error
	EvalTypes() error
	ResolveActivity(usage js.Usage) error
	UniqueEntryPointNames(ns js.Namespace) error
	UniversalNames(ns js.Namespace) error
	UniqueNames(ns js.Namespace) error
  Walk(fn func(p string, obj interface{}) error) error
	Module() js.Module
	Path() string
}

type FileScriptData struct {
	path   string
	module *js.ModuleData
}

var NewViewFileScript func(absPath string, caller string) (FileScript, error) = nil

func SetNewViewFileScript(fn func(absPath string, caller string) (FileScript, error)) bool {
	NewViewFileScript = fn

	return true
}

// if relPath is already absolute, then caller can be left empty
func newFileScriptData(relPath string, caller string) (FileScriptData, error) {
	path, err := files.Search(caller, relPath)
	if err != nil {
		return FileScriptData{}, errors.New("Error: " + relPath + " not found\n")
	}

	// for caching
	files.StartCacheUpdate(path)

	p, err := parsers.NewJSParser(path)
	if err != nil {
		return FileScriptData{}, err
	}

	m, err := p.BuildModule()
	if err != nil {
		return FileScriptData{}, err
	}

	return FileScriptData{path, m}, nil
}

func NewFileScript(relPath string, caller string) (FileScript, error) {
	if strings.HasSuffix(relPath, files.UIFILE_EXT) {
		path, err := files.Search(caller, relPath)
		if err != nil {
			return &FileScriptData{}, errors.New("Error: " + relPath + " not found\n")
		}

		if NewViewFileScript == nil {
			panic("NewViewFileScript not yet registered")
		}

		return NewViewFileScript(path, caller)
	} else {
		s, err := newFileScriptData(relPath, caller)
		if err != nil {
			return nil, err
		}

		return &s, nil
	}
}

func (s *FileScriptData) Hash() string {
	return js.HashControl(s.path)
}

func (s *FileScriptData) Dependencies() []string {
	relPaths := s.module.Dependencies()

	absPaths := make([]string, len(relPaths))

	for i, relPath := range relPaths {
		absPath, err := files.Search(s.Path(), relPath)
		if err != nil {
			panic(err)
		}

		absPaths[i] = absPath
	}

	return absPaths
}

func (s *FileScriptData) Write() (string, error) {
	return s.module.Write()
}

func (s *FileScriptData) ResolveNames(scope js.GlobalScope) error {
	return s.module.ResolveNames(scope)
}

func (s *FileScriptData) EvalTypes() error {
	return s.module.EvalTypes()
}

func (s *FileScriptData) ResolveActivity(usage js.Usage) error {
	return s.module.ResolveActivity(usage)
}

func (s *FileScriptData) UniqueEntryPointNames(ns js.Namespace) error {
	return nil
}

func (s *FileScriptData) UniversalNames(ns js.Namespace) error {
	return s.module.UniversalNames(ns)
}

func (s *FileScriptData) UniqueNames(ns js.Namespace) error {
	return s.module.UniqueNames(ns)
}

func (s *FileScriptData) Walk(fn func(scriptPath string, obj interface{}) error) error {
  return s.module.Walk(func(obj_ interface{}) error {
    return fn(s.path, obj_)
  })
}

func (s *FileScriptData) Module() js.Module {
	return s.module
}

func (s *FileScriptData) Path() string {
	return s.path
}
