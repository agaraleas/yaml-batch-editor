package selectors

type YamlSelector interface {
	Load(yamlRepr interface{}) error
	Run() ([]string, error)
}

func RunYamlSelectors(selectors []YamlSelector) ([]string, error) {
	var files []string
	for _, selector := range selectors {
		batch, err := selector.Run()
		if err != nil {
			return nil, err
		}
		files = append(files, batch...)
	}
	return files, nil
}