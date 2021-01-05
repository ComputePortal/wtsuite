package scripts

import (
	"github.com/computeportal/wtsuite/pkg/files"
	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/js"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
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

var NewViewFileScript func(absPath string) (FileScript, error) = nil

func SetNewViewFileScript(fn func(absPath string) (FileScript, error)) bool {
	NewViewFileScript = fn

	return true
}

func newFileScriptData(absPath string) (FileScriptData, error) {
	// for caching
	files.StartCacheUpdate(absPath)

	p, err := parsers.NewJSParser(absPath)
	if err != nil {
		return FileScriptData{}, err
	}

	m, err := p.BuildModule()
	if err != nil {
		return FileScriptData{}, err
	}

	return FileScriptData{absPath, m}, nil
}

func NewFileScript(absPath string, lang files.Lang) (FileScript, error) {
	if lang == files.TEMPLATE {
		if NewViewFileScript == nil {
			panic("NewViewFileScript not yet registered")
		}

		return NewViewFileScript(absPath)
	} else {
    // called is not needed
		s, err := newFileScriptData(absPath)
		if err != nil {
			return nil, err
		}

		return &s, nil
	}
}

func (s *FileScriptData) Hash() string {
	return js.HashControl(s.path)
}

func (s *FileScriptData) Dependencies() []files.PathLang {
  return s.module.Dependencies()
}

func (s *FileScriptData) Write() (string, error) {
	return s.module.Write(nil, patterns.NL, patterns.TAB)
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
