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

Initialize a repository

Use the init command to create a new GoVCS repository in a specified directory.

./govcs init --path <directory>


Example: Initialize in the current directory:

./govcs init --path .

View repository configuration

The config command displays the current repository’s configuration, including local and global settings.

./govcs config

Set configuration values

Use set-config to define user information or other config values. You can set values locally (for the repo) or globally (for all repos).

# Set user name in local repository config
./govcs set-config --key user.name --value "Your Name" --local

# Set user email in global config
./govcs set-config --key user.email --value "you@example.com" --global

Compute hash of a file

The hash-object command calculates the SHA-1 hash of a file, optionally writing it to the object database.

# Compute hash without writing to database
./govcs hash-object --file <file-path>

# Compute hash and write to the database
./govcs hash-object -w --file <file-path>

Display content of a repository object

Use cat-file to view the contents of a specific object using its SHA hash.

./govcs cat-file --sha <object-sha>

Add a file to the staging area

The add command stages a file for the next commit.

./govcs add --file <file-path>


Example: Stage main.go

./govcs add --file main.go

Commit staged changes

The commit command records the staged changes in the repository with a message describing the commit.

./govcs commit -m "Initial commit"
