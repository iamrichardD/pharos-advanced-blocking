# Google Antigravity Rules: Pharos Advanced Blocking (`pab`)

This file is automatically loaded by the Antigravity CLI and IDE to align agent execution on workspace-specific constraints.

---

## 1. Containerized Development (Podman)
*   **No Host Installation**: Never attempt to run snap, apt, or download tools directly onto the host system.
*   **Compile & Execution**: Always execute Go compiler commands via the public ECR mirror image: `public.ecr.aws/docker/library/golang:1.22-bookworm`.
*   **Sandbox Exemption Flag**: Every `podman run` invocation MUST include the `--security-opt seccomp=unconfined` parameter to bypass OCI filter restrictions on the workstation.
*   **Workspace Mounting**: Ensure the working directory is mounted correctly with `:z` flags: `-v "$(pwd):/workspace:z" -w /workspace`.

---

## 2. Code & Repository Decoupling
*   **No Git Commands**: Do not write code containing execution of `git` commands, commit creation, or indexing hooks.
*   **Config Saving**: Write all configurations locally to disk (`dnsApp.config` or target path), allowing the developer to manage Git staging externally.

---

## 3. Technitium API & Security Rules
*   **API Schema**:
    *   Lease retrieval mappings: `address` (IP), `hardwareAddress` (MAC), `hostName` (Hostname), and `comments` (Description).
    *   Lease deletions: Use `/api/dhcp/scopes/removeReservedLease` and target `hardwareAddress`.
*   **Strict Credentials Guard**: Expose a startup validator that aborts execution with a high-priority warning if the keys file (`~/.config/pab/secrets.json`) has Unix permissions weaker than `600`.
