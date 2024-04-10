package cmdutil

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

/* FUNCTIONS */
func readTOMLConfig(filename string, scriptPath string, runCommandPrefix string) (*Config, error) {
	tomlFile, err := toml.LoadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// In this context, a missing file is not an error
			return nil, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	script := tomlFile.Get(scriptPath)
	if script == nil {
		// No scripts found, return nil
		return nil, nil
	}

	tree, ok := script.(*toml.Tree)
	if !ok {
		return nil, fmt.Errorf("error parsing TOML file: %s is not a TOML tree", scriptPath)
	}

	var config Config
	var commands []CommandConfig

	for _, command := range tree.Keys() {
		position := tree.GetPosition(command)
		commands = append(commands, CommandConfig{
			Name:   command,
			Short:  fmt.Sprintf("Runs the %s defined in %s", command, filename),
			Long:   fmt.Sprintf("Runs the %s defined in %s\n at line %d", command, filename, position.Line),
			Script: []string{fmt.Sprintf("%s%s", runCommandPrefix, command)},
			MetaData: MetaData{
				FilePath: filename,
				Line:     position.Line,
			},
		})
	}
	config.Commands = commands
	return &config, nil
}

func ReadPyProjectTomlConfig(filename string) (*Config, error) {
	scriptPath := "tool.poetry.scripts"
	runCommandPrefix := "poetry run "
	return readTOMLConfig(filename, scriptPath, runCommandPrefix)
}
