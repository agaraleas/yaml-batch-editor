package config

import (
	"fmt"
	"io"
	"os"

	"github.com/agaraleas/yaml-batch-editor/selectors"
	"gopkg.in/yaml.v2"
)

func ParseConfig() (*Config, error) {
	configParser := &configParser{
		configFileReader: &concreteConfigFileReader{ConfigFile},
	}
	return configParser.parse()
}

type fileReader interface {
	read() ([]byte, error)
}

type concreteConfigFileReader struct {
	filepath string
}

func (r *concreteConfigFileReader) read() ([]byte, error) {
	return os.ReadFile(r.filepath)
}

type configParser struct {
	configFileReader fileReader
	WorkingDir string `yaml:"workingDir"`
	DryRun	 bool   `yaml:"dryRun"`
	Selectors []interface{} `yaml:"yamlSelectors"`
}

func (p *configParser) parse() (*Config, error) {
	data, err := p.configFileReader.read()
    if err != nil {
        return nil, err
    }

    err = yaml.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}

	var config Config
	config.WorkingDir = p.WorkingDir
	config.FileWriter = parseFileWriter(p.DryRun)
	config.YamlSelectors, err = parseYamlSelectors(p.Selectors)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func parseFileWriter(dryRun bool) io.Writer {
	if dryRun {
		return os.Stdout
	}
	//temp
	return os.Stdout
}

func parseYamlSelectors(yamlEntries []interface{}) ([]selectors.YamlSelector, error) {
	var yamlSelectors []selectors.YamlSelector
	for _, selector := range yamlEntries {

		switch selector.(type) {
		case map[interface{}]interface{}:
			selector, err := parseYamlSelectorObject(selector.(map[interface{}]interface{}))
			if err != nil {
				return nil, err
			}
			yamlSelectors = append(yamlSelectors, selector)
		default:
			return nil, fmt.Errorf("invalid selector type")
		}
	}

	return yamlSelectors, nil
}

func parseYamlSelectorObject(selector map[interface{}]interface{}) (selectors.YamlSelector, error) {
	grepSelector, ok := selector["grep"]
	if ok {
		return parseGrepSelector(grepSelector)
	}
	return nil, fmt.Errorf("selector type not found")
}

func parseGrepSelector(yamlRepr interface{}) (selectors.YamlSelector, error) {
	var grepSelector selectors.GrepSelector
	err := grepSelector.Load(yamlRepr)
	return &grepSelector, err
}