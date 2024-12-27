package main

import (
	"flag"
	"fmt"
	"gopract/commands"
	"gopract/config"
	"gopract/objects"
	"gopract/repository"
	"gopract/staging"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Ensure global config is initialized
	if err := config.InitializeGlobalConfig(); err != nil {
		fmt.Printf("Failed to initialize global config: %v\n", err)
		return
	}

	// Check if a command is provided
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	fmt.Printf("Received command: %s\n", command)
	switch command {
	case "init":
		handleInit(os.Args[2:])
	case "config":
		handleConfig()
	case "set-config":
		handleSetConfig(os.Args[2:])
	case "hash-object":
		handleHashObject(os.Args[2:])
	case "cat-file":
		handleCatFile(os.Args[2:])
	case "add":
		handleAdd(os.Args[2:])
	case "commit":
		handleCommit(os.Args[2:])
	case "log":
		handleLog()
	case "status":
		handleStatus(os.Args[2:])
	case "checkout":
		handleCheckout(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: gopract <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  init          Initialize a new repository")
	fmt.Println("  config        Show repository configuration")
	fmt.Println("  set-config    Set configuration values")
	fmt.Println("  cat-file      Show content of a repository object")
	fmt.Println("  add           Add files to the staging area")
	fmt.Println("  commit        Commit staged changes to the repository")
	fmt.Println("  log           Show commit history")

}
func handleCheckout(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: gitty checkout <commit-hash> [path]")
		return
	}

	commitHash := args[0]
	var path string
	if len(args) >= 2 {
		path = args[1]
	} else {
		path = "."
	}

	fmt.Printf("Checking out commit: %s\n", commitHash)
	fmt.Printf("Path for checkout: %s\n", path)

	if path != "." {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			fmt.Println("Directory already exists or is not empty!")
			return
		}
	}

	repo, err := repository.Find(".", true)
	if err != nil {
		fmt.Printf("Error locating repository: %v\n", err)
		return
	}

	commit, err := commands.GetCommit(repo, commitHash)
	if err != nil {
		fmt.Printf("Error getting commit %s: %v\n", commitHash, err)
		return
	}

	treeHash := commit.Tree
	treeObj, err := objects.ReadObject(repo.Worktree, treeHash)
	if err != nil {
		fmt.Printf("Error reading tree object: %v\n", err)
		return
	}
	if path != "." {
		err = checkoutTree(repo, treeObj.(*objects.Tree), path)
		if err != nil {
			fmt.Printf("Error checking out tree: %v\n", err)
		} else {
			fmt.Printf("Checked out commit %s into directory %s\n", commitHash, path)
		}
	} else {
		err = checkoutTree(repo, treeObj.(*objects.Tree), repo.Worktree)
		if err != nil {
			fmt.Printf("Error checking out tree: %v\n", err)
		} else {
			fmt.Printf("Checked out commit %s in the current directory\n", commitHash)
		}
	}
}

func checkoutTree(repo *repository.Repository, tree *objects.Tree, path string) error {

	for _, entry := range tree.Entries {
		destPath := filepath.Join(path, entry.Name)
		obj, err := objects.ReadObject(repo.Worktree, entry.Hash)
		if err != nil {
			return fmt.Errorf("Error reading object: %v", err)
		}

		switch obj := obj.(type) {
		case *objects.Blob:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("Error creating directories for %s: %v", destPath, err)
			}
			if err := os.WriteFile(destPath, obj.Data, 0644); err != nil {
				return fmt.Errorf("Error writing file %s: %v", destPath, err)
			}
		case *objects.Tree:
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("Error creating directory %s: %v", destPath, err)
			}
			if err := checkoutTree(repo, obj, destPath); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown object type: %T", obj)
		}
	}
	return nil
}

func handleLog() {
	repo, err := repository.Find(".", true)
	if err != nil {
		fmt.Printf("Error locating repository: %v\n", err)
		return
	}

	headPath := filepath.Join(repo.Gitdir, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		fmt.Printf("Error reading HEAD file: %v\n", err)
		return
	}

	refLine := strings.TrimSpace(string(data))
	if strings.HasPrefix(refLine, "ref: ") {
		ref := strings.TrimPrefix(refLine, "ref: ")
		refPath := filepath.Join(repo.Gitdir, ref)

		commitHashBytes, err := os.ReadFile(refPath)
		if err != nil {
			fmt.Printf("Error reading ref file: %v\n", err)
			return
		}
		commitHashStr := strings.TrimSpace(string(commitHashBytes))

		fmt.Printf("Starting log from commit %s\n", commitHashStr)

		for commitHashStr != "" {
			commit, err := commands.GetCommit(repo, commitHashStr)
			if err != nil {
				fmt.Printf("Error getting commit %s: %v\n", commitHashStr, err)
				return
			}

			fmt.Printf("commit %s\n", commitHashStr)
			fmt.Printf("Author: %s\n", commit.Author)
			fmt.Printf("\n    %s\n\n", commit.Message)

			if len(commit.Parents) > 0 {
				commitHashStr = commit.Parents[0]
			} else {
				commitHashStr = ""
			}
		}
	} else {
		fmt.Println("HEAD is not pointing to a branch, unable to log.")
	}
}

func getHeadCommitHash(repo *repository.Repository) string {
	headPath := filepath.Join(repo.Gitdir, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		fmt.Printf("Error reading HEAD file: %v\n", err)
		return ""
	}

	fmt.Printf("Read HEAD: %s\n", string(data))

	refLine := string(data)
	if len(refLine) > 4 && refLine[:4] == "ref:" {
		refPath := filepath.Join(repo.Gitdir, refLine[5:])
		fmt.Printf("Reading ref from: %s\n", refPath)

		refData, err := os.ReadFile(refPath)
		if err != nil {
			fmt.Printf("Error reading ref file: %v\n", err)
			return ""
		}
		return strings.TrimSpace(string(refData))
	}

	fmt.Println("HEAD does not point to a reference (e.g., refs/heads/master).")
	return ""
}

func handleConfig() {

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

	var localCfg, globalCfg *config.Config
	if localConfigPath != "" {
		localCfg, _ = config.LoadConfig(localConfigPath)
	}
	globalCfg, err = config.LoadConfig(globalConfigPath)
	if err != nil {
		fmt.Printf("Error loading global config: %v\n", err)
		return
	}

	// Basically yahan par local overrides global
	finalCfg := mergeConfigs(localCfg, globalCfg)

	fmt.Printf("Configuration (merged):\n")
	fmt.Printf("User Name: %s\n", finalCfg.User.Name)
	fmt.Printf("User Email: %s\n", finalCfg.User.Email)
	fmt.Printf("Repository Format Version: %d\n", finalCfg.Core.RepositoryFormatVersion)
	fmt.Printf("File Mode: %t\n", finalCfg.Core.FileMode)
	fmt.Printf("Bare Repository: %t\n", finalCfg.Core.Bare)
}
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

	err := config.SetConfigValue(configPath, *key, *value)
	if err != nil {
		fmt.Printf("Error setting config value: %v\n", err)
		return
	}

	fmt.Printf("Configuration updated in %s: [%s] = %s\n", configPath, *key, *value)
}

func handleHashObject(args []string) {
	hashFlags := flag.NewFlagSet("hash-object", flag.ExitOnError)
	write := hashFlags.Bool("w", false, "Write the object to the database")
	filePath := hashFlags.String("file", "", "File path to hash")
	hashFlags.Parse(args)

	if *filePath == "" {
		fmt.Println("File path is required")
		return
	}

	err := commands.HashObject(".", *filePath, *write)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func handleCatFile(args []string) {
	catFlags := flag.NewFlagSet("cat-file", flag.ExitOnError)
	sha := catFlags.String("sha", "", "SHA of the object to read")
	catFlags.Parse(args)

	if *sha == "" {
		fmt.Println("SHA is required")
		return
	}

	err := commands.CatFile(".", *sha)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func handleAdd(args []string) {
	addFlags := flag.NewFlagSet("add", flag.ExitOnError)
	filePath := addFlags.String("file", "", "File to add to the staging area")
	addFlags.Parse(args)

	if *filePath == "" {
		fmt.Println("File path is required")
		return
	}

	err := commands.Add(".", *filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func handleCommit(args []string) {
	commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
	message := commitFlags.String("m", "", "Commit message")
	commitFlags.Parse(args)

	if *message == "" {
		fmt.Println("Commit message is required")
		return
	}

	repo, err := repository.Find(".", true)
	if err != nil {
		fmt.Printf("Error locating repository: %v\n", err)
		return
	}

	index, err := staging.ReadIndex(repo.Worktree)
	if err != nil {
		fmt.Printf("Error reading index: %v\n", err)
		return
	}

	tree := &objects.Tree{}
	for filePath, blobHash := range index {
		tree.Entries = append(tree.Entries, objects.TreeEntry{
			Mode: "100644",
			Hash: blobHash,
			Name: filePath,
		})
	}

	treeHash, err := objects.WriteObject(tree, repo.Worktree)
	if err != nil {
		fmt.Printf("Error writing tree object: %v\n", err)
		return
	}

	var parents []string
	headPath := filepath.Join(repo.Gitdir, "HEAD")
	head, err := os.ReadFile(headPath)
	if err == nil {
		ref := string(head)
		refPath := filepath.Join(repo.Gitdir, ref)
		if parentHash, err := os.ReadFile(refPath); err == nil {
			parents = append(parents, string(parentHash))
		}
	}
	commit := &objects.Commit{
		Tree:    treeHash,
		Parents: parents,
		Author:  "Khushi Surana<suranakhushi17@gmail.com>",
		Message: *message,
		Date:    time.Now().Format(time.RFC3339),
	}

	commitHash, err := objects.WriteObject(commit, repo.Worktree)
	if err != nil {
		fmt.Printf("Error writing commit object: %v\n", err)
		return
	}

	refPath := filepath.Join(repo.Gitdir, "refs/heads/master")
	if err := os.WriteFile(refPath, []byte(commitHash), 0644); err != nil {
		fmt.Printf("Error updating HEAD: %v\n", err)
		return
	}

	fmt.Printf("Committed with hash %s\n", commitHash)
}
func handleStatus(args []string) {
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}
	err := commands.Status(repoPath)
	if err != nil {
		fmt.Printf("Error checking status: %v\n", err)
	}
}
