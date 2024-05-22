# runr (WIP)

Go-based CLI tool that aims to improve visibility of available commands in projects that use package managers and/or build tools, such as Poetry and Make.

Currently Supported:
- `package.json` for yarn projects
- `pyproject.toml` for poetry projects
- `Makefile` for projects that use Make

Planned Support:
- `package.json` for npm projects

Users can also define their own commands in a project via a `runr.yaml` file or globally via a `~/.runr_global.yaml` file.


## Installation
Download the latest release from the [releases page](
    https://github.com/jacobtavener/runr/releases/
) and add it to your PATH.

You can also clone the repository and install the binary yourself.


[!NOTE]
The builds are not signed, so you may need to allow the binary to run in your system settings.

[!WARNING]
This project is still in development and may not be stable. Use at your own risk.
It has only been tested on MacOS.


## Usage
```
runr [command] [flags]
```

Running `runr` without any arguments will display the available commands in the current directory.

The available commands are determined by:
- Looking for the supported package managers and/or build tool files in the current directory.
- Looking for a `runr.yaml` file in the current directory.
- Looking for a `~/.runr_global.yaml` file in the home directory.

Commands defined in the `runr.yaml` file will take precedence over commands defined in the package managers and build files.

Commands defined in the package managers and build files will take precedence over commands defined in the `~/.runr_global.yaml` file.

If there are multiple commands with the same name in the package managers and build files, `runr` will run all of them, unless a specific command is specified using the `-t` flag.


### Flags
- `-h`, `--help`: Display help for the command.

- [command] `-h`, `--help`: Display help for a specific command.

- [command] `-e`, `--edit`: Open the file for the command at the specified line in a vim editor.

- [command] `-t`, `--tool`: If there are commands with the same name across the different package managers and build files, you can specify which one to run. By default, `runr` will run all the commands with the same name.

## Configuration
`runr` can be configured via a `runr.yaml` file in the root of your project or a `~/.runr_global.yaml` file in your home directory.

The configuration file should be in the following format:

```yaml
commands:
  - name: my-custom-command
    help: Short description of the command. This will be displayed when running `runr`.
    description: longer description of the command, which will be displayed when running `runr my-custom-command --help`. 
    env:
      - name: MY_ENV_VAR
        value: world
    script: 
      - echo "Hello, $MY_ENV_VAR!"
      - echo "This is a custom command." 
```
 Note: If multiple scripts are defined, they will be run in order, in different shell sessions. 