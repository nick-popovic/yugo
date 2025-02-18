package config_lib

import (
	"io/ioutil"
	"os/exec"
	"runtime"

	"gopkg.in/yaml.v2"
)

type PackageManager struct {
	Name     string `yaml:"name"`
	CheckCmd string `yaml:"check_cmd"`
}

type Program struct {
	Description string                       `yaml:"description"`
	Tags        []string                     `yaml:"tags"`
	Installs    map[string]map[string]string `yaml:"installs"`
}

type Config struct {
	PackageManagers map[string][]PackageManager `yaml:"package_managers"`
	Programs        map[string]Program          `yaml:"programs"`
}

// IsPackageManagerAvailable checks if a package manager is available for the current OS
func (c *Config) IsPackageManagerAvailable() bool {
	var osKey string
	switch runtime.GOOS {
	case "darwin":
		osKey = "darwin"
	case "windows":
		osKey = "windows"
	case "linux":
		osKey = "linux"
	default:
		return false
	}

	if _, ok := c.PackageManagers[osKey]; ok {
		return true
	}
	return false
}

// GetAvailablePackageManagers returns the names of available package managers for the current OS
func (c *Config) GetAvailablePackageManagers() []string {
	var osKey string
	switch runtime.GOOS {
	case "darwin":
		osKey = "darwin"
	case "windows":
		osKey = "windows"
	case "linux":
		osKey = "linux"
	default:
		return nil
	}

	var availableManagers []string
	if managers, ok := c.PackageManagers[osKey]; ok {
		for _, manager := range managers {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("cmd", "/C", manager.CheckCmd)
			} else {
				cmd = exec.Command("sh", "-c", manager.CheckCmd)
			}
			if err := cmd.Run(); err == nil {
				availableManagers = append(availableManagers, manager.Name)
			}
		}
	}
	return availableManagers
}

// GetProgramTagsByPackageManager returns a list of program tags supported by the given package manager and system OS
func (c *Config) GetProgramTagsByPackageManager(packageManager string) []string {
	var osKey string
	switch runtime.GOOS {
	case "darwin":
		osKey = "darwin"
	case "windows":
		osKey = "windows"
	case "linux":
		osKey = "linux"
	default:
		return nil
	}

	var tags []string
	for _, program := range c.Programs {
		if installs, ok := program.Installs[osKey]; ok {
			if _, ok := installs[packageManager]; ok {
				tags = append(tags, program.Tags...)
			}
		}
	}

	// Remove duplicate tags
	tagMap := make(map[string]bool)
	var uniqueTags []string
	for _, tag := range tags {
		if _, ok := tagMap[tag]; !ok {
			tagMap[tag] = true
			uniqueTags = append(uniqueTags, tag)
		}
	}

	return uniqueTags
}

// GetInstallCommands generates a list of installation commands for the current OS, package manager, and tags
func (c *Config) GetInstallCommands(packageManager string, tags []string) []string {
	var osKey string
	switch runtime.GOOS {
	case "darwin":
		osKey = "darwin"
	case "windows":
		osKey = "windows"
	case "linux":
		osKey = "linux"
	default:
		return nil
	}

	var commands []string

	for _, program := range c.Programs {
		// Check if the program has any of the specified tags
		hasTag := false
		for _, tag := range tags {
			for _, programTag := range program.Tags {
				if tag == programTag {
					hasTag = true
					break
				}
			}
			if hasTag {
				break
			}
		}

		// If the program has the specified tags and an install command for the given OS and package manager, add it to the list
		if hasTag {
			if installs, ok := program.Installs[osKey]; ok {
				if cmd, ok := installs[packageManager]; ok {
					commands = append(commands, cmd)
				}
			}
		}
	}

	return commands
}

func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
