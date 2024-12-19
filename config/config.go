package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	User struct {
		Name  string `ini:"name"`
		Email string `ini:"email"`
	}
	Core struct {
		RepositoryFormatVersion int  `ini:"repositoryformatversion"`
		FileMode                bool `ini:"filemode"`
		Bare                    bool `ini:"bare"`
	}
}

// InitializeGlobalConfig ensures the global config file exists with default values.
// InitializeGlobalConfig ensures the global config file exists with default values.
func InitializeGlobalConfig() error {
	// Create the global config file in the current directory
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	globalPath := filepath.Join(usr.HomeDir, ".mygitconfig")

	if _, err := os.Stat(globalPath); os.IsNotExist(err) {
		fmt.Printf("Creating global config file at: %s\n", globalPath)
		iniFile := ini.Empty()
		iniFile.Section("user").Key("name").SetValue("Default User")
		iniFile.Section("user").Key("email").SetValue("default@example.com")
		err = iniFile.SaveTo(globalPath)
		if err != nil {
			return fmt.Errorf("failed to create global config file: %w", err)
		}
	}

	return nil
}

// GetGlobalConfigPath retrieves the path to the global `.mygitconfig` file.
func GetGlobalConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return filepath.Join(usr.HomeDir, ".mygitconfig"), nil
}

// LoadConfig loads a configuration file.
func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to determine absolute path: %w", err)
	}
	fmt.Println("Loading config from:", absPath)

	iniFile, err := ini.Load(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	cfg := new(Config)
	cfg.User.Name = iniFile.Section("user").Key("name").String()
	cfg.User.Email = iniFile.Section("user").Key("email").String()
	cfg.Core.RepositoryFormatVersion = iniFile.Section("core").Key("repositoryformatversion").MustInt(0)
	cfg.Core.FileMode = iniFile.Section("core").Key("filemode").MustBool(false)
	cfg.Core.Bare = iniFile.Section("core").Key("bare").MustBool(false)

	// Debugging Output
	fmt.Println("Loaded Configuration:")
	fmt.Printf("User Name: %s\n", cfg.User.Name)
	fmt.Printf("User Email: %s\n", cfg.User.Email)
	fmt.Printf("Repository Format Version: %d\n", cfg.Core.RepositoryFormatVersion)
	fmt.Printf("File Mode: %t\n", cfg.Core.FileMode)
	fmt.Printf("Bare Repository: %t\n", cfg.Core.Bare)

	return cfg, nil
}

// SetConfigValue sets a key-value pair in the configuration file.
func SetConfigValue(path, key, value string) error {
	iniFile, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Split the key into section and field (e.g., "user.name")
	parts := strings.SplitN(key, ".", 2) // Correctly splits on the first '.'
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format, expected 'section.key'")
	}

	section, field := parts[0], parts[1]
	iniFile.Section(section).Key(field).SetValue(value)

	err = iniFile.SaveTo(path)
	if err != nil {
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}
