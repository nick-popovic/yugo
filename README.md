# Yu(go)

Yugo is an assistant for batch installing programs on a clean OS install written in Go. It helps select and install various programs using different package managers based on your operating system and the programs you add to your configuration file.

## Features

- Detects available package managers for your OS.
- Provides a list of program tags supported by the selected package manager.
- Generates and executes installation commands for the selected programs.



## Configuration

1. Create a configuration file `config.yaml` in the `$HOME/.config/yugo`. Below is an example:
 
2. The configuration file should define the available package managers and programs. Each program should include a description, tags, and installation commands for different package managers and operating systems.

## Sample Configuration:
```yaml
package_managers:
  darwin:
    - name: brew
      check_cmd: brew -v
  windows:
    - name: choco
      check_cmd: choco -v
  linux:
    - name: apt
      check_cmd: apt --version
      
programs:
  git:
    description: "Git Version Control System"
    tags: ["version control"]
    installs:
      linux:
        apt: sudo apt install -y git
      darwin:
        brew: brew install git
      windows:
        choco: choco install git