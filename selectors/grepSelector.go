package selectors

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type GrepSelector struct {
	grepArgs []string `yaml:"args"`
}

func (g *GrepSelector) Args() []string {
    return g.grepArgs
}

func (g *GrepSelector) Run() ([]string, error) {
    args := prependArgument(g.grepArgs, "-l")
	cmd := exec.Command("grep", args...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            // Check if the exit status is 1, which means no matches found
            if exitError.ExitCode() == 1 {
                return []string{}, nil
            }
        }
        return nil, fmt.Errorf("command failed: %w, stderr: %s", err, output)
    }

	return extractFilepaths(output)
}

func(g *GrepSelector) Load(yamlRepr interface{}) error {
    obj, ok := yamlRepr.(map[interface{}]interface{})
    if !ok {
        return fmt.Errorf("invalid yaml representation: grep not a map")
    }

    args, ok := obj["args"]
    if !ok {
        return fmt.Errorf("invalid yaml representation: grep args not found")
    }

    argsSlice, ok := args.([]interface{})
    if !ok {
        return fmt.Errorf("invalid yaml representation: grep args not a list")
    }

    for _, arg := range argsSlice {
        argStr, ok := arg.(string)
        if !ok {
            return fmt.Errorf("invalid yaml representation: grep arg not a string")
        }
        g.grepArgs = append(g.grepArgs, argStr)
    }
    return nil
}

func extractFilepaths(output []byte) ([]string, error) {
    var filePaths []string
	
	lines := strings.Split(string(output), "\n")

    // Iterate over each line and convert to absolute path
    for _, line := range lines {
        if line != "" {
            //find first .yaml: and stop
            absPath, err := filepath.Abs(line)
            if err != nil {
                return nil, err
            }
            filePaths = append(filePaths, absPath)
        }
    }

    return filePaths, nil
}

func prependArgument(args []string, arg string) []string {
    for _, a := range args {
        if a == arg {
            return args
        }
    }
    return append([]string{arg}, args...)
}