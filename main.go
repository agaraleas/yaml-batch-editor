package main

import (
	"log"
	"os"

	"github.com/agaraleas/yaml-batch-editor/config"
	"github.com/agaraleas/yaml-batch-editor/selectors"
)

func main() {
	config, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(config.WorkingDir)
    if err != nil {
        log.Fatal(err)
    }

	files, err := selectors.RunYamlSelectors(config.YamlSelectors)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Println(file)
	}
}