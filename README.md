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
