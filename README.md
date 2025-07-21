# Red Hat Developer Hub (RHDH) Installer PoC

A CLI tool for deploying Red Hat Developer Hub (RHDH) with a set of default configuration options set using the new RHDH profile workflow.

## Overview

This repository contains:
- **`/profile`**: Kustomize manifests for deploying the RHDH operator using the external profile workflow.
- **`/presets`**: Preset configurations for full RHDH installation.
- **`rhdh-cli`**: A Go CLI tool that wraps the Kustomize deployments.

## Prerequisites

- Go 1.21 or later
- `kubectl` configured to access your Kubernetes cluster (with built-in Kustomize support)
- Kubernetes cluster with appropriate permissions

## Building the CLI

```bash
make build
```

### Development:
```bash
# Download dependencies
make deps

## Using the CLI

### Basic Usage

### Deploy the RHDH operator

./rhdh-cli deploy operator

### Deploy the RHDH preset

./rhdh-cli deploy presets

### Deploy presets with a custom .env file:

./rhdh-cli deploy presets --env-file local.env
```

### Command Options

- `--env-file`: Path to environment file (default: ".env")
- `--verbose, -v`: Enable verbose output
- `--dry-run, -d`: Perform dry run without actual deployment
- `--timeout`: Timeout in seconds for deployment operations (default: 600)

### Environment Configuration

Create a `.env` file to customize your deployment:

```bash
cp example.env local.env
# Edit .env with your configuration
```

The CLI supports template variable substitution in manifests using environment variables from the `.env` file. Variables can be referenced as `${VAR_NAME}` or `$VAR_NAME` in your Kustomize manifests.

You will need to create a `rhdh-secrets.yaml` file under `/rhdh` based on the template `example-secrets.yaml`

```bash
cp ./presets/rhdh-complete/rhdh/example-secrets.yaml ./presets/rhdh-complete/rhdh/rhdh-secrets.yaml
```

### Examples

Deploy presets with custom environment file:
```bash
./rhdh-cli deploy presets \
  --env-file production.env \
  --verbose
```

Perform a dry run to see what would be deployed:
```bash
./rhdh-cli deploy presets --dry-run --verbose
```

## Makefile Commands

The repository includes both traditional kubectl/kustomize commands and new CLI-based commands:

| Command | Description |
| ------- | ----------- |
| **make deploy-operator-cli** | Deploy the RHDH Operator |
| **make deploy-presets-cli** | Deploy the RHDHPAI Preset for RHDH |
| **make build** | Make the Go binary |
| **make deps** | Tidy and verify go mod |
