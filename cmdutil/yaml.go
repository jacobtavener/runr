package cmdutil

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

/* TYPES */
type ConfigYaml struct {
	Commands  []CommandConfigYaml `yaml:"commands"`
	Variables []VarType           `yaml:"variables"`
}

type CommandConfigYaml struct {
	Name        yaml.Node
	Hint        string    `yaml:"hint"`
	Description string    `yaml:"description"`
	Script      []string  `yaml:"script"`
	Flags       []Flag    `yaml:"flags"`
	Env         []VarType `yaml:"env"`
}

/* GLOBALS */
var variableRegex = regexp.MustCompile(`\\$([a-zA-Z0-9_]+)`)

/* FUNCTIONS */
func replaceVariables(input string, variables map[string]string) string {
	result := variableRegex.ReplaceAllStringFunc(input, func(match string) string {
		varName := match[1:]
		if value, ok := variables[varName]; ok {
			return value
		}
		return match // If variable not found, return the original match
	})
	return result
}

func replaceVariablesInScript(script []string, variables map[string]string) []string {
	result := make([]string, len(script))
	for i, line := range script {
		result[i] = replaceVariables(line, variables)
	}
	return result
}

func ReadYamlConfig(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// In this context, a missing file is not an error
			return nil, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var configYaml ConfigYaml
	err = yaml.Unmarshal(yamlFile, &configYaml)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	variables := make(map[string]string)
	for _, v := range configYaml.Variables {
		variables[v.Name] = v.Value
	}

	var config Config
	for _, cmd := range configYaml.Commands {
		var name string
		var metaData MetaData
		switch cmd.Name.Kind {
		case yaml.ScalarNode:
			name = cmd.Name.Value
			metaData = MetaData{
				FilePath: filename,
				Line:     cmd.Name.Line,
			}
		default:
			log.Fatalf("Error parsing YAML file: command name is not a scalar")
			continue
		}
		config.Commands = append(config.Commands, CommandConfig{
			Name:     name,
			MetaData: metaData,
			Short:    cmd.Hint,
			Long:     fmt.Sprintf("%s\n  - %s at line %d", cmd.Description, filename, cmd.Name.Line),
			Script:   replaceVariablesInScript(cmd.Script, variables),
			Flags:    cmd.Flags,
			Env:      cmd.Env,
		})
	}
	return &config, nil
}
