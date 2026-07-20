# Product Requirements Document (PRD): Pharos Advanced Blocking (`pab`)

## 1. Executive Summary & Objective

**Pharos Advanced Blocking (`pab`)** is a statically linked command-line interface (CLI) tool written in Go, designed to manage, validate, and synchronize Technitium DNS Server **Advanced Blocking App** configurations (`dnsApp.config`). 

The tool is built **CLI-first** and structured as an **AI-agent-friendly** tool, allowing human operators to interact with it via a rich Terminal User Interface (TUI), while permitting AI coding assistants and automation scripts to use structured JSON payloads and non-interactive parameters.

---

## 2. Core Architecture & Philosophy

1.  **Decoupled & GitOps Ready**: The tool works entirely **on disk** when editing configurations. It has zero coupling with Git binaries, allowing users to choose their own Git deployment, staging, or pipeline engine.
2.  **View First, Edit Second**: The primary interactive user workflow is optimized for scanning, querying, and checking status before writing changes to disk or publishing to the servers.
3.  **Cross-Compilation Target**:
    *   `linux/amd64`
    *   `linux/arm64`
4.  **No Runtime Dependencies**: Statically compiled binary requiring no local runtime environment (like Node.js, Python, or external interpreters).

---

## 3. Functional Requirements

### 3.1 Workspace Detection & Initialization
When executed, the CLI will look for the presence of a configuration file in the working directory:
*   By default, it searches for `dnsApp.config` (Technitium's native filename).
*   If not found, it prompts the user with interactive options:
    1.  **Create New**: Generate a clean, validated `dnsApp.config` template with default structural groups on disk.
    2.  **Load Custom Path**: Specify a different config file path (e.g. `advanced-blocking-config.json`).
    3.  **Fetch from Server**: Connect to one or more configured Technitium API endpoints, download the active running configuration, and write it to disk.

### 3.2 View & Query Modes (Human vs. AI Agent)
*   **TUI Mode (Default for Interactive TTY)**:
    *   Renders a beautiful ASCII table using `lipgloss` showing client IPs mapped to their blocking groups.
    *   Provides search/filter capabilities to find which group a specific IP belongs to.
    *   Lists the active groups and their configured blocklists/regex filters.
*   **Machine-Readable Mode (AI Agent Friendly)**:
    *   Adding the `--json` flag prints clean, parsed JSON outputs to stdout.
    *   Silences all interactive terminal visual animations (like spinners or progress bars).
    *   Supports the `--quiet` and `--no-color` flags.

### 3.3 Edit & Modification Workflow
*   **Add/Update Mapping**: Interactive console commands or direct arguments to map/reassign a client IP or range (CIDR) to a specific blocking group:
    ```bash
    pab map --ip 192.168.86.16 --group "Default+Richard"
    ```
*   **Remove Mapping**: Deletes a client IP mapping from the configuration:
    ```bash
    pab unmap --ip 192.168.86.18
    ```

### 3.4 Schema Validation Engine
Before saving any modifications to disk or posting to an API, the CLI runs a validation pass:
1.  **IP & Subnet Check**: Ensures every key in the client network map is a valid IPv4, IPv6, or CIDR network block (automatically rejecting MAC addresses or hostname strings, which cause Technitium crashes).
2.  **Group Integrity**: Verifies that any client mapping targets a group that actually exists in the `groups` section.
3.  **JSON Schema Check**: Validates blocklist URL structures and regex patterns.

### 3.5 Sync & Deploy Engine
Once validated, the CLI can sync the disk configuration to one or more target Technitium installations:
*   **Command**: `pab deploy [flags]`
*   **Dry Run**: `--dry-run` performs a full structural diff between the disk configuration and the server API, listing what will be updated on each server without applying changes.
*   **Target Selection**: Deploy to all configured servers, or specify target nodes via `--node technitium-01`.
*   **Non-interactive**: The `--yes` flag bypasses confirmation prompts.

---

## 4. Technical Specifications & Stack

### 4.1 CLI Interface (Go)
*   **CLI Router & Parser**: `github.com/spf13/cobra` for handling subcommands, flags, and arguments.
*   **TUI Engine**: `github.com/charmbracelet/bubbletea` for rich terminal interactions.
*   **Styling**: `github.com/charmbracelet/lipgloss` for padding, borders, and color definitions.
*   **Interactive Prompts**: `github.com/AlecAivazis/survey` for simple wizard flows.

### 4.2 Security & Credential Store
To protect Technitium API tokens, the CLI will look for secrets in the following order:
1.  **Environment Variables**: `TECHNITIUM_TOKEN_01`, `TECHNITIUM_TOKEN_02`, etc.
2.  **OS Secure Config Path**: If environment variables are absent, reads credentials from a local configuration directory (e.g. `~/.config/pab/secrets.json`).
    *   **Requirements**: The CLI will refuse to run and print a warning if this configuration file does not have strict permissions (`chmod 600`), preventing other system users from reading the file contents.
3.  **Onboarding Wizard**: Prompts to securely paste tokens and saves them directly to `~/.config/pab/secrets.json` with permissions set automatically to `600`.

### 4.3 Plugin Extensibility
To support the future **Live Status Plugin** (e.g. displaying real-time queries and lease information across multiple nodes), the CLI core will implement:
*   A directory-based plugin loader (e.g., loading compiled Go plugins `.so` or Javascript WebAssembly extensions from a configured plug-in folder).
*   Dynamic command registration, allowing loaded plugins to register subcommands (like `pab status`) into Cobra's runtime registry.

---

## 5. Build, Release, & Installation

### 5.1 Compilation & Assembly
The project will use Go's native cross-compilation capability. We will integrate **GoReleaser** in a GitHub Actions workflow to build release binaries.

### 5.2 Release Artifacts
1.  **Static Binary**: Standalone compressed `.tar.gz` archive containing the binary.
2.  **Debian Package (`.deb`)**:
    *   Constructed via `goreleaser` or `nfpm`.
    *   Integrates with Debian-based systems' standard package registry.
    *   Exposes clean install and remove routes (`apt install ./pab.deb` / `apt remove pab`).

### 5.3 Installation Scripts
*   **Bash Installer**: An `install.sh` script hosted in the Git repository that:
    1.  Detects system CPU architecture (rejecting 32-bit x86 architectures).
    2.  Downloads the latest release binary matching the architecture from GitHub Releases.
    3.  Verifies the SHA256 checksum against the official release manifest.
    4.  Extracts and copies the binary to `/usr/local/bin/pab`.
