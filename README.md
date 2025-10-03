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

## Installation

1. **Clone the repository**  
```bash
git clone https://github.com/<your-username>/GoVCS.git
cd GoVCS
## Usage / Commands

After cloning and building the GoVCS binary:

```bash
go build -o govcs
You can run the following commands:

Initialize a repository
bash
Copy code
./govcs init --path <directory>
# Example: Initialize in the current directory
./govcs init --path .
View repository configuration
bash
Copy code
./govcs config
Set configuration values
bash
Copy code
# Set user name in local repository config
./govcs set-config --key user.name --value "Your Name" --local

# Set user email in global config
./govcs set-config --key user.email --value "you@example.com" --global
Compute hash of a file
bash
Copy code
# Compute hash without writing to object database
./govcs hash-object --file <file-path>

# Compute hash and write to the database
./govcs hash-object -w --file <file-path>
Display content of a repository object
bash
Copy code
./govcs cat-file --sha <object-sha>
Add a file to the staging area
bash
Copy code
./govcs add --file <file-path>
# Example: Stage main.go
./govcs add --file main.go
Commit staged changes
bash
Copy code
./govcs commit -m "Initial commit"
Notes / Tips
Run all commands from the repository root directory.

Flags like --local and --global apply only to set-config.

Always build the binary first using go build -o govcs.

You can use relative or absolute file paths in commands like add or hash-object.

This CLI mimics Git, so commands like add and commit work similarly.
