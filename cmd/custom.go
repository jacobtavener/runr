package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jacobtavener/runr/cmdutil"

	"github.com/spf13/cobra"
)

/* TYPES */
type projectCommand struct {
	Filepath   string
	Alias      string
	ConfigFunc func(string) (*cmdutil.Config, error)
}

type commandMap map[string]cobra.Command
type projectCommandMap map[string]map[string]cmdutil.CommandConfig

/* GLOBALS */
const (
	GlobalRunrConfig = "runr_global.yaml"
	LocalRunrConfig  = "runr.yaml"
	EditFlag         = "edit"
	ToolFlag         = "tool"
)

var ProjectCommands = [3]projectCommand{
	{Filepath: "pyproject.toml", Alias: "poetry", ConfigFunc: cmdutil.ReadPyProjectTomlConfig},
	{Filepath: "package.json", Alias: "yarn", ConfigFunc: cmdutil.ReadPackageJSONConfig}, // TODO: Should automatically detect yarn or npm
	{Filepath: "Makefile", Alias: "make", ConfigFunc: cmdutil.ReadMakeFileConfig},
}

/* FUNCTIONS */
func constructCommands(scripts []string, extraArgs []string, env []cmdutil.VarType) []cmdutil.Command {
	cmds := make([]cmdutil.Command, len(scripts))
	for i, command := range scripts {
		split := strings.Split(command, " ")

		if len(split) == 0 {
			fmt.Println("Invalid script format:", command)
			continue
		}

		name := split[0]
		args := split[1:]
		cmdArgs := make([]string, len(args)+len(extraArgs))
		copy(cmdArgs, args)
		copy(cmdArgs[len(args):], extraArgs)

		cmds[i] = cmdutil.Command{Name: name, Args: cmdArgs, Env: env}
	}
	return cmds
}

func constructProjectCommandMap() projectCommandMap {
	projectCommandMap := make(projectCommandMap)
	for _, projectCmd := range ProjectCommands {
		commands, err := projectCmd.ConfigFunc(projectCmd.Filepath)
		if err == nil && commands != nil {
			for _, cmdConfig := range commands.Commands {
				localCmdConfig := cmdConfig
				cmd, ok := projectCommandMap[localCmdConfig.Name]
				if !ok {
					cmd = make(map[string]cmdutil.CommandConfig)
					projectCommandMap[localCmdConfig.Name] = cmd
				}
				cmd[projectCmd.Alias] = localCmdConfig
			}
		}
	}
	return projectCommandMap
}

func setUpProjectCommands(commandMap commandMap) {
	projectCommandMap := constructProjectCommandMap()

	for name, cmds := range projectCommandMap {
		localCmds := cmds
		localName := name
		var localShort string
		var localLong string

		localShort = fmt.Sprintf("Runs the %s command using", localName)
		localLong = fmt.Sprintf("Runs the %s command using:\n", localName)
		available_tools := make([]string, 0)
		for tool, cmdConfig := range localCmds {
			localShort += fmt.Sprintf(" %s,", tool)
			localLong += fmt.Sprintf("   - %s from %s at line %d\n", tool, cmdConfig.MetaData.FilePath, cmdConfig.MetaData.Line)
			available_tools = append(available_tools, tool)
		}
		localShort = strings.TrimSuffix(localShort, ",")
		cmd := &cobra.Command{
			Use:   localName,
			Short: localShort,
			Long:  localLong,
			Run: func(cmd *cobra.Command, args []string) {
				if len(localCmds) > 1 {
					tool, _ := cmd.Flags().GetString(ToolFlag)
					editMode := cmd.Flag(EditFlag).Changed
					if _, ok := localCmds[tool]; !ok {
						if editMode {
							fmt.Println("You need to specify a tool to edit the command")
							os.Exit(1)
						}
						for tool, localCmdConfig := range localCmds {
							fmt.Println("Running command using", tool)
							cmdutil.RunCommands(constructCommands(localCmdConfig.Script, args, localCmdConfig.Env))
						}
					}
					if editMode {
						cmdutil.EditFile(localCmds[tool].MetaData.FilePath, localCmds[tool].MetaData.Line)
					} else {
						cmdutil.RunCommands(constructCommands(localCmds[tool].Script, args, localCmds[tool].Env))
					}
				} else {
					for _, cmdConfig := range localCmds {
						localCmdConfig := cmdConfig
						editMode := cmd.Flag(EditFlag).Changed
						if editMode {
							cmdutil.EditFile(localCmdConfig.MetaData.FilePath, localCmdConfig.MetaData.Line)
						} else {
							cmdutil.RunCommands(constructCommands(localCmdConfig.Script, args, localCmdConfig.Env))
						}
					}
				}
			},
		}
		cmdutil.AddEditFlag(cmd)
		if len(available_tools) > 1 {
			cmd.Flags().StringP(ToolFlag, "t", "", "Tool used to run the command (available: "+strings.Join(available_tools, ",")+")")
		}
		commandMap[localName] = *cmd
	}
}

func setUpRunrCommands(commandMap commandMap, filepath string) {
	config, err := cmdutil.ReadYamlConfig(filepath)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", filepath, err)
		return
	}

	if config == nil {
		return
	}
	for _, cmdConfig := range config.Commands {
		localCmdConfig := cmdConfig
		cmd := &cobra.Command{
			Use:   localCmdConfig.Name,
			Short: localCmdConfig.Short,
			Long:  localCmdConfig.Long,
			Run: func(cmd *cobra.Command, args []string) {
				editMode := cmd.Flag(EditFlag).Changed
				if editMode {
					cmdutil.EditFile(localCmdConfig.MetaData.FilePath, localCmdConfig.MetaData.Line)
				} else {
					cmdutil.RunCommands(constructCommands(localCmdConfig.Script, args, localCmdConfig.Env))
				}
			},
		}
		cmdutil.AddEditFlag(cmd)
		cmdutil.AddFlags(cmd, localCmdConfig.Flags) // TODO: Does this even do anything?
		cmdutil.AddCommandUnique(cmd, rootCmd)
	}
}

func addCommandsToRoot(commandMap commandMap) {
	for _, cmd := range commandMap {
		localCmd := cmd
		rootCmd.AddCommand(&localCmd)
	}
}

func SetUpCustomCommands() {
	commandMap := map[string]cobra.Command{}
	if homeDir, err := os.UserHomeDir(); err == nil {
		setUpRunrCommands(commandMap, filepath.Join(homeDir, GlobalRunrConfig))
	}
	setUpProjectCommands(commandMap)
	setUpRunrCommands(commandMap, LocalRunrConfig)
	addCommandsToRoot(commandMap)
}
