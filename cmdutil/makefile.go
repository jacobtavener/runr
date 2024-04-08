package cmdutil

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* FUNCTIONS */
func parseMakeFile(filename string, parsedFiles map[string]bool) (map[string]MetaData, error) {
	if parsedFiles[filename] {
		return nil, nil
	}

	parsedFiles[filename] = true

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	defer file.Close()

	dir := filepath.Dir(filename)
	targetMap := make(map[string]MetaData)

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		words := strings.Fields(line)
		if len(words) > 1 && words[0] == "include" {
			includedFile := filepath.Join(dir, words[1])

			includedResult, err := parseMakeFile(includedFile, parsedFiles)
			if err != nil {
				return nil, err
			}
			for target, metaData := range includedResult {
				targetMap[target] = metaData
			}
		}

		if len(words) > 0 && strings.HasSuffix(words[0], ":") {
			target := strings.TrimSuffix(words[0], ":")
			targetMap[target] = MetaData{
				FilePath: filename,
				Line:     lineNumber,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targetMap, nil
}

func ReadMakeFileConfig(filename string) (*Config, error) {
	targets, err := parseMakeFile(filename, make(map[string]bool))
	if err != nil {
		return nil, err
	}

	var config Config
	var commands []CommandConfig

	for command, metaData := range targets {
		commands = append(commands, CommandConfig{
			Name:     command,
			Short:    fmt.Sprintf("Runs the %s defined in %s", command, filename),
			Long:     fmt.Sprintf("Runs the %s defined in %s\n at line %d", command, filename, metaData.Line),
			Script:   []string{fmt.Sprintf("%s%s", "make ", command)},
			MetaData: metaData,
		})
	}
	config.Commands = commands
	return &config, nil
}
