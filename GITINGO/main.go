package main

import (
	"flag"
	"fmt"
	"gopract/config"
	"gopract/repository"
	"os"
	"path/filepath"
)

func main() {
	// Ensure global config is initialized
	if err := config.InitializeGlobalConfig(); err != nil {
		fmt.Printf("Failed to initialize global config: %v\n", err)
		return
	}

	// Check if a command is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: gopract <command> [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		handleInit(os.Args[2:])
	case "config":
		handleConfig()
	case "set-config":
		handleSetConfig(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// handleConfig processes the `config` command to display configuration details.
func handleConfig() {
	// Locate local and global configs
	repo, err := repository.Find(".", false)
	var localConfigPath, globalConfigPath string

	if err == nil && repo != nil {
		localConfigPath = filepath.Join(repo.Gitdir, "config")
		fmt.Printf("Using local config at: %s\n", localConfigPath)
	}

	globalConfigPath, err = config.GetGlobalConfigPath()
	if err != nil {
		fmt.Printf("Error locating global config file: %v\n", err)
		return
	}
	fmt.Printf("Using global config at: %s\n", globalConfigPath)

	// Load configurations
	var localCfg, globalCfg *config.Config
	if localConfigPath != "" {
		localCfg, _ = config.LoadConfig(localConfigPath)
	}
	globalCfg, err = config.LoadConfig(globalConfigPath)
	if err != nil {
		fmt.Printf("Error loading global config: %v\n", err)
		return
	}

	// Merge configurations (local overrides global)
	finalCfg := mergeConfigs(localCfg, globalCfg)

	// Display merged configuration
	fmt.Printf("Configuration (merged):\n")
	fmt.Printf("User Name: %s\n", finalCfg.User.Name)
	fmt.Printf("User Email: %s\n", finalCfg.User.Email)
	fmt.Printf("Repository Format Version: %d\n", finalCfg.Core.RepositoryFormatVersion)
	fmt.Printf("File Mode: %t\n", finalCfg.Core.FileMode)
	fmt.Printf("Bare Repository: %t\n", finalCfg.Core.Bare)
}

// mergeConfigs merges local and global configurations, prioritizing local values.
func mergeConfigs(local, global *config.Config) *config.Config {
	final := &config.Config{}

	if local != nil {
		final.User.Name = local.User.Name
		final.User.Email = local.User.Email
		final.Core.RepositoryFormatVersion = local.Core.RepositoryFormatVersion
		final.Core.FileMode = local.Core.FileMode
		final.Core.Bare = local.Core.Bare
	}

	if global != nil {
		if final.User.Name == "" {
			final.User.Name = global.User.Name
		}
		if final.User.Email == "" {
			final.User.Email = global.User.Email
		}
		if final.Core.RepositoryFormatVersion == 0 {
			final.Core.RepositoryFormatVersion = global.Core.RepositoryFormatVersion
		}
		if !final.Core.FileMode {
			final.Core.FileMode = global.Core.FileMode
		}
		if !final.Core.Bare {
			final.Core.Bare = global.Core.Bare
		}
	}

	return final
}

// handleInit processes the `init` command to initialize a repository.
func handleInit(args []string) {
	initFlags := flag.NewFlagSet("init", flag.ExitOnError)
	repoPath := initFlags.String("path", ".", "Path where the repository should be created")
	initFlags.Parse(args)

	repo, err := repository.NewRepository(*repoPath, true)
	if err != nil {
		fmt.Printf("Error creating repository object for path %s: %v\n", *repoPath, err)
		return
	}

	err = repo.Create()
	if err != nil {
		fmt.Printf("Error initializing repository in path %s: %v\n", *repoPath, err)
		return
	}

	fmt.Printf("Initialized empty Git repository in %s\n", repo.Worktree)
}

// handleSetConfig updates configuration values.
func handleSetConfig(args []string) {
	setConfigFlags := flag.NewFlagSet("set-config", flag.ExitOnError)
	key := setConfigFlags.String("key", "", "Configuration key to set (e.g., user.name)")
	value := setConfigFlags.String("value", "", "Value to set for the configuration key")
	isLocal := setConfigFlags.Bool("local", false, "Set value in local config (.git/config)")
	isGlobal := setConfigFlags.Bool("global", false, "Set value in global config (.mygitconfig)")
	setConfigFlags.Parse(args)

	if *key == "" || *value == "" {
		fmt.Println("Both --key and --value must be provided")
		return
	}

	// Determine target config file
	var configPath string
	if *isLocal {
		repo, err := repository.Find(".", true)
		if err != nil {
			fmt.Printf("Error locating local repository: %v\n", err)
			return
		}
		configPath = filepath.Join(repo.Gitdir, "config")
	} else if *isGlobal {
		var err error
		configPath, err = config.GetGlobalConfigPath()
		if err != nil {
			fmt.Printf("Error locating global config file: %v\n", err)
			return
		}
	} else {
		fmt.Println("Please specify either --local or --global")
		return
	}

	// Update the config file
	err := config.SetConfigValue(configPath, *key, *value)
	if err != nil {
		fmt.Printf("Error setting config value: %v\n", err)
		return
	}

	fmt.Printf("Configuration updated in %s: [%s] = %s\n", configPath, *key, *value)
}
