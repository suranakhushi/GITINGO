package repository

import (
	"errors"
	"fmt"
	"gopract/config" // Import the config package to load the global config
	"os"
	"path/filepath"
)

// Repository represents a Git repository.
type Repository struct {
	Worktree string
	Gitdir   string
}

// NewRepository creates a new repository object and validates its structure.
func NewRepository(path string, force bool) (*Repository, error) {
	worktree, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	gitdir := filepath.Join(worktree, ".git")
	if !force {
		info, err := os.Stat(gitdir)
		if os.IsNotExist(err) || !info.IsDir() {
			return nil, fmt.Errorf("not a Git repository: %s", path)
		}
	}

	return &Repository{Worktree: worktree, Gitdir: gitdir}, nil
}

// Create initializes a new Git repository.
func (r *Repository) Create() error {
	// Create the .git directory
	if err := os.MkdirAll(r.Gitdir, 0755); err != nil {
		return fmt.Errorf("failed to create .git directory: %w", err)
	}

	// Create subdirectories
	subDirs := []string{"branches", "objects", "refs/heads", "refs/tags"}
	for _, dir := range subDirs {
		path := filepath.Join(r.Gitdir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}

	// Create description file
	description := filepath.Join(r.Gitdir, "description")
	if err := os.WriteFile(description, []byte("Unnamed repository; edit this file to name the repository.\n"), 0644); err != nil {
		return fmt.Errorf("failed to write description: %w", err)
	}

	// Create HEAD file
	head := filepath.Join(r.Gitdir, "HEAD")
	if err := os.WriteFile(head, []byte("ref: refs/heads/master\n"), 0644); err != nil {
		return fmt.Errorf("failed to write HEAD: %w", err)
	}

	// Create config file and merge with global config
	configPath := filepath.Join(r.Gitdir, "config")
	globalConfigPath, err := config.GetGlobalConfigPath()
	if err != nil {
		return fmt.Errorf("failed to locate global config: %w", err)
	}

	globalConfig, err := config.LoadConfig(globalConfigPath)
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	localConfigContent := `[core]
repositoryformatversion = 0
filemode = true
bare = false
`

	// Add user details from global config if available
	if globalConfig.User.Name != "" {
		localConfigContent += fmt.Sprintf("[user]\nname = %s\n", globalConfig.User.Name)
	}
	if globalConfig.User.Email != "" {
		localConfigContent += fmt.Sprintf("email = %s\n", globalConfig.User.Email)
	}

	if err := os.WriteFile(configPath, []byte(localConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Initialized repository with config at %s\n", configPath)
	return nil
}

// Find locates the repository root starting from a given path.
func Find(path string, required bool) (*Repository, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	gitdir := filepath.Join(absPath, ".git")
	if info, err := os.Stat(gitdir); err == nil && info.IsDir() {
		return NewRepository(absPath, false)
	}

	parent := filepath.Dir(absPath)
	if parent == absPath {
		if required {
			return nil, errors.New("no Git repository found")
		}
		return nil, nil
	}

	return Find(parent, required)
}
