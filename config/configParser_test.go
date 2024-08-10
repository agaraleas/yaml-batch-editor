package config

import (
	"testing"

	"github.com/agaraleas/yaml-batch-editor/selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const grepKeyword = "grep"
const argsKeyword = "args"
const yamlSelectorsKeyword = "yamlSelectors"
const workingDirKeyword = "workingDir"
const dryRunKeyword = "dryRun"

// -------------------------------------------------------
type mockConfigReader struct {
	content []byte
}

func (r *mockConfigReader) read() ([]byte, error) {
	return r.content, nil
}
// -------------------------------------------------------

func buildWorkingDirContentForCurrentDir() []byte {
	return []byte(workingDirKeyword + ": \".\"")
}

func TestParseConfigWorkingDir_Current(t *testing.T) {
    configReader := &mockConfigReader{buildWorkingDirContentForCurrentDir()}
	configParser := &configParser{configFileReader: configReader}
	config, err := configParser.parse()
	assert.Nil(t, err)
	assert.Equal(t, ".", config.WorkingDir)
}

func buildWorkingDirContentForHomeDir() []byte {
	return []byte(workingDirKeyword + ": \"~\"")
}
func TestParseConfigWorkingDir_Home(t *testing.T) {
    configReader := &mockConfigReader{buildWorkingDirContentForHomeDir()}
	configParser := &configParser{configFileReader: configReader}
	config, err := configParser.parse()
	assert.Nil(t, err)
	assert.Equal(t, "~", config.WorkingDir)
}

func buildWorkingDirContentForRelativeDir() []byte {
	return []byte(workingDirKeyword + ": \"./test\"")
}

func TestParseConfigWorkingDir_Relative(t *testing.T) {
    configReader := &mockConfigReader{buildWorkingDirContentForRelativeDir()}
	configParser := &configParser{configFileReader: configReader}
	config, err := configParser.parse()
	assert.Nil(t, err)
	assert.Equal(t, "./test", config.WorkingDir)
}

func buildGrepYamlSelectorContentWithArgs() []byte {
	var content []byte
	content = append(content, []byte(yamlSelectorsKeyword+":\n")...)
	content = append(content, []byte("  - "+grepKeyword+":\n")...)
	content = append(content, []byte("      "+argsKeyword+":\n")...)
	content = append(content, []byte("        - --exclude-dir=\"dir\"\n")...)
	content = append(content, []byte("        - -rnw \n")...)
	content = append(content, []byte("        - '.'\n")...)
	content = append(content, []byte("        - -e\n")...)
	content = append(content, []byte("        - '\"namespace: lala-lala-lala\"'\n")...)
	return content
}

func TestParseConfigGrepSelectorWithArgs(t *testing.T) {
	configReader := &mockConfigReader{buildGrepYamlSelectorContentWithArgs()}
	configParser := &configParser{configFileReader: configReader}
	config, err := configParser.parse()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(config.YamlSelectors))
	grepSelector, ok := config.YamlSelectors[0].(*selectors.GrepSelector)
	require.True(t, ok)
	args := grepSelector.Args()
	require.Equal(t, 5, len(args))
	assert.Equal(t, "--exclude-dir=\"dir\"", args[0])
	assert.Equal(t, "-rnw", args[1])
	assert.Equal(t, ".", args[2])
	assert.Equal(t, "-e", args[3])
	assert.Equal(t, "\"namespace: lala-lala-lala\"", args[4])
}