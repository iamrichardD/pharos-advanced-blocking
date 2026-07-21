# Quick Start: Pharos Advanced Blocking Local Development

Get up and running with pab development in 60 seconds.

## Prerequisites

- Podman installed and running
- Claude Code CLI with this project as working directory
- Git configured

## The Three Essential Commands

All Go operations run in a containerized environment (`golang:1.22-bookworm`). You don't need Go installed locally.

### 1. Run Unit Tests (3-5 seconds)

```bash
claude pab test
```

Runs `go test ./... -v` inside Podman. Fails fast if tests break. **Run this first.**

### 2. Compile Binary (5-10 seconds)

```bash
claude pab compile
```

Builds the pab binary for both `linux/amd64` and `linux/arm64` using Podman. Produces `./pab` binary ready for testing.

### 3. Full Release Snapshot (15-20 seconds)

```bash
claude pab snapshot
```

Executes the complete GoReleaser snapshot build pipeline. Validates:
- License metadata (AGPL-3.0)
- Checksum integrity
- Cross-platform binary generation
- Release artifacts structure

## Development Workflow

**Typical cycle while coding:**

```bash
# Edit your code
vim internal/config/models.go

# Run tests
claude pab test          # ✅ Pass? Great. ❌ Fail? Fix and repeat.

# Build binary
claude pab compile       # ✅ Compiles? Test it. ❌ Error? Fix.

# When done with feature, validate full release pipeline
claude pab snapshot      # ✅ Release valid? Ready to review.
```

## Code Review Gate

When your code is ready and all tests pass, invoke the mob programming review panel:

```bash
claude pab review
```

This gates your code changes through 7 expert perspectives:
- **Go Engineering** — Implementation quality, Podman patterns
- **DevSecOps** — Security, credential handling, release integrity
- **Kent Beck** — Simplicity, feedback loops, test coverage
- **Robert Martin** — SOLID principles, clean code
- **Martin Fowler** — Domain modeling, architecture
- **Kathy Sierra** — CLI ergonomics, user experience
- **Seth Godin** — Brand alignment, product positioning

The mob provides feedback:
- ✅ **Approved** — Ready to push
- 🟡 **Approved with Changes** — Specific fixes requested
- ❌ **Rejected** — Fundamental issues; major revision needed

Repeat: `claude pab review` after addressing feedback until approved.

## Push & Release

After mob approval:

```bash
git push origin <branch>   # Triggers CI/CD pipeline
```

To release a new version:

```bash
git tag v1.2.3
git push origin v1.2.3     # Triggers GitHub Actions release workflow
```

GitHub Actions handles Cosign signing and GoReleaser publishing—no local intervention needed.

## Keyboard Shortcuts

Speed up your workflow:

| Shortcut | Command | Purpose |
|----------|---------|---------|
| `Ctrl+Shift+T` | `claude pab test` | Run tests |
| `Ctrl+Shift+B` | `claude pab compile` | Build binary |
| `Ctrl+Shift+S` | `claude pab snapshot` | Release snapshot |
| `Ctrl+Shift+R` | `claude pab review` | Mob review |

## Common Issues

**"Command not found: podman"**
- Install Podman: `sudo apt install podman` (or your package manager)
- Ensure Podman socket is accessible: `podman ps`

**"First run is slow (pulling golang:1.22-bookworm)"**
- One-time download from AWS ECR (~300MB). Subsequent runs use cached image.
- Check: `podman images | grep golang`

**Tests pass locally but fail in CI**
- Verify Podman image version matches CI: `golang:1.22-bookworm`
- Check for uncommitted changes: `git status`

**Mob review taking > 2 minutes**
- Large diffs take longer. Break into smaller commits if needed.
- Mob gates code changes only; expected time at merge boundary.

## Next Steps

1. **Read DEVELOPMENT.md** for architecture and implementation roadmap
2. **Read PRD.md** for product requirements and feature specs
3. **Review .agents/workflows/pab-feature-loop.md** for the full workflow diagram
4. **Start coding** — Use `claude pab test/compile/review` commands in your IDE

---

**Questions?** Check `.agents/` for detailed agent and workflow documentation.
