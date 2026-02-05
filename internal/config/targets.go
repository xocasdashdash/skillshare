package config

import (
	_ "embed"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
	"skillshare/internal/utils"
)

type targetSpec struct {
	GlobalName  string `yaml:"global_name"`
	ProjectName string `yaml:"project_name"`
	GlobalPath  string `yaml:"global_path"`
	ProjectPath string `yaml:"project_path"`
}

type targetsFile struct {
	Targets []targetSpec `yaml:"targets"`
}

//go:embed targets.yaml
var defaultTargetsData []byte

var (
	loadedTargets   []targetSpec
	loadTargetsErr  error
	loadTargetsOnce sync.Once
)

func loadTargetSpecs() ([]targetSpec, error) {
	loadTargetsOnce.Do(func() {
		var file targetsFile
		if err := yaml.Unmarshal(defaultTargetsData, &file); err != nil {
			loadTargetsErr = err
			return
		}
		loadedTargets = file.Targets
	})

	return loadedTargets, loadTargetsErr
}

// DefaultTargets returns the well-known CLI skills directories for global mode.
func DefaultTargets() map[string]TargetConfig {
	specs, err := loadTargetSpecs()
	if err != nil {
		return map[string]TargetConfig{}
	}

	targets := make(map[string]TargetConfig)
	for _, spec := range specs {
		if spec.GlobalName == "" || spec.GlobalPath == "" {
			continue
		}
		path := normalizeTargetPath(spec.GlobalPath)
		targets[spec.GlobalName] = TargetConfig{Path: path}
	}

	return targets
}

// ProjectTargets returns the well-known CLI skills directories for project mode.
func ProjectTargets() map[string]TargetConfig {
	specs, err := loadTargetSpecs()
	if err != nil {
		return map[string]TargetConfig{}
	}

	targets := make(map[string]TargetConfig)
	for _, spec := range specs {
		if spec.ProjectName == "" || spec.ProjectPath == "" {
			continue
		}
		path := normalizeTargetPath(spec.ProjectPath)
		targets[spec.ProjectName] = TargetConfig{Path: path}
	}

	return targets
}

// LookupProjectTarget returns the known project target config for a name.
func LookupProjectTarget(name string) (TargetConfig, bool) {
	targets := ProjectTargets()
	target, ok := targets[name]
	return target, ok
}

// LookupGlobalTarget returns the known global target config for a name.
func LookupGlobalTarget(name string) (TargetConfig, bool) {
	targets := DefaultTargets()
	target, ok := targets[name]
	return target, ok
}

func normalizeTargetPath(path string) string {
	if path == "" {
		return path
	}
	if utils.HasTildePrefix(path) {
		path = expandPath(path)
	}
	return filepath.FromSlash(path)
}
