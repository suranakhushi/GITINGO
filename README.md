# GoVCS: A Version Control System in Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

GoVCS is a **Git-inspired version control system** built in **Go (Golang)**. It enables repository cloning, commit history tracking, remote synchronization, and efficient parallel task execution using Go’s concurrency model. GoVCS is designed to help developers **learn the inner workings of Git** and experiment with a lightweight VCS implementation.

---

## Features

- **Repository Management:** Initialize, clone, and manage repositories locally.  
- **Commit Tracking:** Track changes with commit history, similar to Git.  
- **Remote Synchronization:** Push and pull changes between repositories.  
- **Concurrency Optimizations:** Parallelized tasks for faster operations using Go routines.  
- **Lightweight & Educational:** Perfect for learning the mechanics of version control systems.  

---

## Tech Stack

- **Language:** Go (Golang)  
- **Core Concepts:** Concurrency (goroutines), File I/O, CLI tools, Data Serialization  
- **Version Control Principles:** Inspired by Git’s object storage and commit history mechanism  

---
## GoVCS Commands

### Clone the repository
Clone your GoVCS repository from GitHub:

```bash
git clone https://github.com/your-username/GoVCS.git
cd GoVCS
Build the GoVCS binary
Compile the GoVCS CLI tool:


go build -o govcs
Initialize a repository
Create a new GoVCS repository:

./govcs init --path <directory>
Initialize in the current directory:


./govcs init --path .
View repository configuration
Check the current repository configuration (local and global):


./govcs config
Set configuration values
Set your user name locally in the repository:


./govcs set-config --key user.name --value "Your Name" --local
Set your user email globally:


./govcs set-config --key user.email --value "you@example.com" --global
Compute hash of a file
Compute the hash of a file without storing it:


./govcs hash-object --file <file-path>
Compute the hash and store it in the object database:


./govcs hash-object -w --file <file-path>
Display content of a repository object
View the contents of an object by SHA:


./govcs cat-file --sha <object-sha>
Add a file to the staging area
Stage a file for commit:


./govcs add --file <file-path>
Example: Stage main.go:


./govcs add --file main.go
Commit staged changes
Commit all staged files with a message:


./govcs commit -m "Initial commit"
