package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"../../../pkg/files"
)

type AnimationScene struct {
	Control string `json:"control"` // these files will be relative upon import, must be made absolute
	View    string `json:"view"`
}

type Animation struct {
	Scenes []AnimationScene `json:"scenes"`
}

func relToAbsAnimationControlsAndViews(fpath string, cfg *Animation) error {

	for i, scene := range cfg.Scenes {
		controlAbs, err := files.Search(fpath, scene.Control)
		if err != nil {
			return err
		}

		viewAbs, err := files.Search(fpath, scene.View)
		if err != nil {
			return err
		}

		scene.Control = controlAbs
		scene.View = viewAbs

		cfg.Scenes[i] = scene
	}

	return nil
}

//fpath is absolute
func ReadAnimationFile(fpath string) (*Animation, error) {
	cfg := &Animation{
		Scenes: make([]AnimationScene, 0),
	}

	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return cfg, errors.New("Error: problem reading the animation config file")
	}

	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, errors.New("Error: bad animation config file syntax (" + err.Error() + ")")
	}

	// IncludeDirs will already have been appened in files package during
	if err := relToAbsAnimationControlsAndViews(fpath, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// input view is abspath
func (cfg *Animation) GetViewScenes(view string) []int {
	res := make([]int, 0)

	for i, scene := range cfg.Scenes {
		if scene.View == view {
			res = append(res, i)
		}
	}

	return res
}
