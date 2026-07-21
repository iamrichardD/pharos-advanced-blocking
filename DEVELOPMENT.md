# Developer & AI Agent Guidelines: Pharos Advanced Blocking (`pab`)

This document outlines key technical constraints, workflows, and implementation directives for developers and AI Coding Assistants working on this repository.

---

## 🤖 CRITICAL DIRECTIVES FOR AI AGENTS (Do Not Bypass)

1.  **No Host Tool Installation**: Do NOT attempt to install development tools (such as `golang-go` or `snap`) on the host system. All compilations, dependency updates (`go get`), and executions MUST be run inside the Podman container environment using the AWS ECR public mirror.
2.  **Container Execution Flag**: Every `podman run` invocation MUST include the `--security-opt seccomp=unconfined` flag to bypass workstation sandbox restrictions (bdflush OCI permission errors).
3.  **Strict Git Decoupling**: Do NOT write Go code that interacts with the `git` binary or attempts to commit files. All configuration changes must be written directly to disk (`dnsApp.config`), leaving Git management exclusively to the developer's external workflow.
4.  **Technitium API Schema Integrity**: 
    *   When fetching/updating DHCP leases, map the parameters to the correct Technitium keys: `address` (IP), `hardwareAddress` (MAC), `hostName` (Hostname), and `comments` (Description).
    *   When deleting/releasing leases, use `/api/dhcp/scopes/removeReservedLease` and pass the `hardwareAddress` (formatted with colons) rather than the IP.
5.  **Strict Security Checks**: The CLI must refuse to run and exit with a high-priority security error if the credentials file (`~/.config/pab/secrets.json`) is configured with permissions weaker than `chmod 600` (Unix user read/write only).

---

## 🛠️ Step-by-Step Implementation Roadmap

To develop this project in future sessions, implement modules in this sequence:

### Phase 1: Structs & Disk Manager (`internal/config/`)
1.  Define the Go structs representing the Advanced Blocking configuration schema in `internal/config/models.go`.
2.  Implement validation rules in `internal/config/validator.go` to reject hostnames or MAC addresses in the client mapping keys, ensuring only valid IPs and CIDR blocks are parsed.

### Phase 2: API Client (`internal/client/`)
1.  Implement the Technitium REST HTTP client in `internal/client/client.go`.
2.  Implement methods: `AddReservation`, `RemoveReservation`, `SetAppConfig`, and `FetchCurrentScope`.

### Phase 3: CLI Commands (`cmd/pab/`)
1.  Implement Cobra command routing under `cmd/pab/`.
2.  Define subcommands:
    *   `pab map --ip <IP> --group <GROUP>`
    *   `pab unmap --ip <IP>`
    *   `pab deploy` (with `--dry-run` and `--yes` overrides).

### Phase 4: Terminal User Interface (`internal/tui/`)
1.  Implement a Bubble Tea model-view loop in `internal/tui/` to render client lists, groups, and search filters.
2.  Style all terminal outputs using Lip Gloss to match Pharos brand aesthetics.

---

## ⌨️ TUI Features & Interactions

### Tab Completion (v0.2.0+)

**Behavior:** IDE-style auto-complete for slash commands.

- **Single match**: Pressing Tab completes the command to `<command-name> ` (with trailing space)
- **Multi-match**: Tab completes to the currently highlighted candidate + space; use Up/Down arrows to navigate before Tab if needed
- **After completion**: User can immediately type subcommand arguments (e.g., after Tab-completing `/view ` to see `/view `, typing `groups` produces `/view groups`)
- **Outside command mode**: Tab is a no-op (safe fallback)

**Example flows:**
```
User: /v[Tab]           → /view 
User: /[Down][Down][Tab] → /exit  (if /exit was second in list)
User: /view groups[Enter] → Execute view groups command
```

**Implementation note**: Exit typeahead mode after Tab completion (`inTypeaheadMode = false`) to allow subcommand argument typing without re-filtering against command names.

### Slash Commands

Available commands (typed with `/` prefix):
- `/help` — Show available slash commands
- `/clear` — Reset search and clear filters
- `/exit` — Exit TUI
- `/view groups` — List all groups with device counts
- `/view group <name>` — Show details for a specific group (blocklists, allowed/blocked domains)
- `/view networkGroupMap` — Show IP-to-group mappings

---

## 🚀 Common Dev Commands (Cheat Sheet)

```bash
# Initialize/Update dependencies
podman run --rm --security-opt seccomp=unconfined -v "$(pwd):/workspace:z" -w /workspace public.ecr.aws/docker/library/golang:1.22-bookworm go mod tidy

# Compile the binary natively to host disk
podman run --rm --security-opt seccomp=unconfined -v "$(pwd):/workspace:z" -w /workspace public.ecr.aws/docker/library/golang:1.22-bookworm go build -o pab cmd/pab/main.go
```
