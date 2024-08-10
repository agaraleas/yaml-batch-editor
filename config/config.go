package config

import (
	"io"

	"github.com/agaraleas/yaml-batch-editor/selectors"
)

var ConfigFile = "./config.yaml"

type Config struct {
	WorkingDir string
	FileWriter io.Writer
	YamlSelectors []selectors.YamlSelector
	//ProcessEngine
}
