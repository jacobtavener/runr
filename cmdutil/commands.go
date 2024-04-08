package cmdutil

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

/* TYPES */
type Config struct {
	Commands []CommandConfig
}

type Flag struct {
	Flag        string
	Description string
	Type        string
	Default     string
}

type VarType struct {
	Name  string
	Value string
}

type MetaData struct {
	FilePath string
	Line     int
}

type CommandConfig struct {
	Name     string
	MetaData MetaData
	Short    string
	Long     string
	Script   []string
	Flags    []Flag
	Env      []VarType
}
type Command struct {
	Name string
	Args []string
	Env  []VarType
}

/* FUNCTIONS */
func runCommand(name string, env []VarType, args ...string) (int, error) {
	cmd := exec.Command(name, args...)

	additionalEnv := os.Environ()
	for _, e := range env {
		additionalEnv = append(additionalEnv, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	cmd.Env = additionalEnv

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	err := cmd.Start()
	if err != nil {
		return -1, err
	}

	go func() {
		<-sigCh
		fmt.Println("Received interrupt signal. Stopping the command...")
		cmd.Process.Signal(os.Interrupt)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error running command:", err)
	}

	return cmd.ProcessState.ExitCode(), nil
}

func RunCommands(commands []Command) {
	for _, cmd := range commands {
		exitCode, err := runCommand(cmd.Name, cmd.Env, cmd.Args...)
		if err != nil {
			fmt.Println("Error running command:", err)
		}
		if exitCode != 0 {
			break
		}
	}
}

func AddCommandUnique(cmd *cobra.Command, parent *cobra.Command) {
	existingCmd, _, _ := parent.Find([]string{cmd.Name()})
	if existingCmd != nil {
		parent.RemoveCommand(existingCmd)
	}
	parent.AddCommand(cmd)
}

func EditFile(filename string, lineNumber int) {
	// TODO: Add support for other editors
	editor := "vim"
	cmd := exec.Command(editor, fmt.Sprintf("+%d", lineNumber), filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening editor:", err)
	}
}
