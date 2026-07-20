# Pharos Advanced Blocking (`pab`)

`pab` is a statically linked command-line interface (CLI) tool written in Go, designed to manage, validate, and synchronize Technitium DNS Server **Advanced Blocking App** configurations (`dnsApp.config`). 

This project is an open-source tool developed by **Pharos Systems (LLC)** to build home-lab developer utility, establish brand equity, and support the broader Pharos ecosystem.

---

## Features
*   **CLI-First & AI-Agent Friendly**: Optimized for human terminal usage via a rich TUI and scripts/AI assistants via standard JSON outputs (`--json` flag).
*   **GitOps & Decoupled Architecture**: Performs configuration edits entirely on disk (`dnsApp.config`), leaving repository commits/pushes to the user's preferred Git workflow.
*   **Schema Validation**: Validates client IP and subnet formats before writing, blocking invalid formats (such as hostnames or MAC addresses) that cause Technitium DNS failures.
*   **Multi-Node Sync Engine**: Syncs and deploys local configuration rules directly to one or more Technitium API endpoints.
*   **Secure Credential Handling**: Locks session keys and tokens in a user configuration directory (`~/.config/pab/secrets.json`) protected with `chmod 600` file permissions.

---

## Getting Started (Development via Podman)

To run or build the application without installing Go on your host workstation, you can run all compiler operations inside a Podman container using the public ECR mirror registry.

### 1. Initialize the Modules & Dependencies
```bash
# Verify the module is initialized
podman run --rm --security-opt seccomp=unconfined -v "$(pwd):/workspace:z" -w /workspace public.ecr.aws/docker/library/golang:1.22-bookworm go mod tidy
```

### 2. Run the Application in Development Mode
To execute the app in real-time:
```bash
podman run --rm --security-opt seccomp=unconfined -v "$(pwd):/workspace:z" -w /workspace public.ecr.aws/docker/library/golang:1.22-bookworm go run cmd/pab/main.go
```

### 3. Compile a Native Executable for the Host
Compile the statically linked Go binary. It will write the output binary `pab` directly to your workstation directory, allowing you to run it natively on your system:
```bash
# Build the binary
podman run --rm --security-opt seccomp=unconfined -v "$(pwd):/workspace:z" -w /workspace public.ecr.aws/docker/library/golang:1.22-bookworm go build -o pab cmd/pab/main.go

# Execute natively on your host workstation
./pab --help
```

---

## Product Specifications
For detailed architecture, technical stack selections, functional workflows, and release guidelines, please refer to the [Product Requirements Document (PRD)](docs/06_technitium_blocking_cli_prd.md).
