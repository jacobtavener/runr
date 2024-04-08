package main

import (
	"github.com/jacobtavener/runr/cmd"
)

func main() {
	cmd.SetUpCustomCommands()
	cmd.Execute()
}
