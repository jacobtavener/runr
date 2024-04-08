package cmdutil

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func IsYarnProject() bool {
	return fileExists("yarn.lock")
}

func IsNpmProject() bool {
	return fileExists("package-lock.json")
}

func findLineNumber(parentKey string, key string, filename string) (int, error) {
	sedCmd := fmt.Sprintf(`/"scripts": {/,/},/ { /"%s"/=;}`, key)
	cmd := exec.Command("sed", "-n", sedCmd, filename)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	lineNumStr := strings.TrimSpace(string(output))
	lineNumber, err := strconv.Atoi(lineNumStr)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	return lineNumber, nil
}

func ReadPackageJSONConfig(filename string) (*Config, error) {
	packageJsonFile, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var packageJson map[string]interface{}
	err = json.Unmarshal(packageJsonFile, &packageJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON file: %v", err)
	}

	var runCommandPrefix string
	switch {
	case IsYarnProject():
		runCommandPrefix = "yarn run "
	case IsNpmProject():
		runCommandPrefix = "npm run "
	default:
		return nil, fmt.Errorf("error: no package manager found")
	}

	if commandMap, ok := packageJson["scripts"].(map[string]interface{}); ok {
		var config Config
		var commands []CommandConfig

		for command := range commandMap {
			lineNumber, _ := findLineNumber("scripts", command, filename)

			commands = append(commands, CommandConfig{
				Name:   command,
				Short:  fmt.Sprintf("Runs the %s defined in %s", command, filename),
				Long:   fmt.Sprintf("Runs the %s defined in %s\n at line %d", command, filename, lineNumber),
				Script: []string{fmt.Sprintf("%s%s", runCommandPrefix, command)},
				MetaData: MetaData{
					FilePath: filename,
					Line:     lineNumber,
				},
			})
		}
		config.Commands = commands
		return &config, nil
	} else {
		return nil, nil
	}
}
