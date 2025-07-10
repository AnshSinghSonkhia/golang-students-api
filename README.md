# Golang Students REST API Project

**Tech-Stack:** Golang, SQLite

## To Create Go Project

```go
go mod init github.com/AnshSinghSonkhia/golang-students-api
```

## Project Structure

```text
Main-Folder
    |-> cmd
    |   |-> golang-students-api     <-- Your Project Name
    |       |-> main.go             <-- Main Entry File of Project
    |-> config
        |-> local.yaml
    |-> go.mod                      <-- run 'go mod init github.com/AnshSinghSonkhia/golang-students-api' command in cmd to get this file
```

# `Config` folder

The `config` folder in a Go project is essential for managing environment-specific configuration settings. By maintaining separate files such as `local.yaml` and `production.yaml`, you can:

- **Isolate Environment Variables:** Keep development (local) and production settings separate, reducing the risk of accidentally deploying development configurations to production.
- **Simplify Deployment:** Easily switch between environments by loading the appropriate configuration file, streamlining the deployment process.
- **Enhance Security:** Sensitive information (like API keys or database credentials) can be managed securely and kept out of the source code.
- **Improve Maintainability:** Centralizing configuration in dedicated files makes it easier to update and manage settings without modifying application code.
- **Facilitate Collaboration:** Team members can use their own `local.yaml` without affecting shared or production settings.

In summary, using a `config` folder with environment-specific YAML files promotes best practices for configuration management, security, and maintainability in Go projects.

# Storage - SQLite

# Golang Packages Used

```bash
go get -u github.com/ilyakaznacheev/cleanenv
```

# To run the server with config flag

```bash
go run cmd/golang-students-api/main.go -config config/local.yaml
```

## [ToDo] Create 2 more endpoints 
 
- PUT - Update Student
- DELETE - Delete a Student