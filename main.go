package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/charmbracelet/huh"
	"gopkg.in/yaml.v2"

	config_lib "main/config_lib"
)

type Order struct {
	PackageManager string
	tags           []string
}

func main() {

	var order Order
	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error getting home directory: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(homeDir, ".config", "yugo", "config.yaml"))
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config config_lib.Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	availablePackageManagers := config.GetAvailablePackageManagers()
	if len(availablePackageManagers) == 0 {
		log.Fatalf("No package managers available for the current OS identified in .config")
	}

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("YugoScript").
			Description("Welcome to _YugoScript_.\n\nAssistant for batch installing programs on a clean install ...\n\n").
			Next(true).
			NextLabel("Next"),
		),

		// Choose a package manager.
		huh.NewGroup(

			huh.NewSelect[string]().
				Options(huh.NewOptions(availablePackageManagers...)...).
				Title("Choose your package manager").
				Description("Select a package manager available for your OS.").
				Validate(func(t string) error {
					if t == "" {
						return fmt.Errorf("a package manager must be selected")
					}
					return nil
				}).
				Value(&order.PackageManager),
		),
	).WithAccessible(accessible)

	err = form.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	order.tags = config.GetProgramTagsByPackageManager(order.PackageManager)

	var tag_options []huh.Option[string]
	for _, tag := range order.tags {
		// second argument ("fruit") is the value,
		// first argument is the user-facing label
		tag_options = append(tag_options, huh.NewOption(tag, tag))
	}

	////////////////////////////////////////////////////////////////////////

	form_1 := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Tags").
				Description("Choose your tags").
				Options(tag_options...).
				Validate(func(t []string) error {
					if len(t) <= 0 {
						return fmt.Errorf("at least one tag is required")
					}
					return nil
				}).
				Value(&order.tags).
				Filterable(true),
		),
	).WithAccessible(accessible)

	err = form_1.Run()
	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	for _, cmd_str := range config.GetInstallCommands(order.PackageManager, order.tags) {
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", cmd_str)
		} else {
			cmd = exec.Command("sh", "-c", cmd_str)
		}

		// Capture the output of the command
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Error executing command %s: %v", cmd_str, err)
		}
		fmt.Printf("Output of command %s:\n%s\n", cmd_str, string(output))
	}
}
