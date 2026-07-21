---
name: pab-feature-loop
description: Feature development loop with builder inner loop and mob review gate
metadata:
  type: workflow
---

# pab Feature Development Loop

Orchestrates the complete flow for developing, testing, and reviewing code changes to Pharos Advanced Blocking.

## Workflow Stages

### Stage 1: Fast Inner Loop (Developer, No Mob)

Developer codes locally and runs fast feedback commands repeatedly:

```bash
claude pab test      # go test ./... -v in Podman (3-5 sec, fail fast)
claude pab compile   # go build in Podman (5-10 sec, binary ready)
claude pab snapshot  # goreleaser snapshot in Podman (15-20 sec, full validation)
```

**No mob review happens during this stage.** The builder agent provides immediate feedback:
- ✅ Tests pass or ❌ Tests fail with error details
- ✅ Binary compiles or ❌ Compile error with context
- ✅ Snapshot valid or ❌ Release validation failed

**Repeat**: Developer iterates on code, re-runs these commands until satisfied.

### Stage 2: Mob Review Gate (Before Push)

Once code changes are complete and tests pass locally, invoke mob review:

```bash
claude pab review
```

The mob programming review panel evaluates:
- Architecture and domain modeling
- SOLID principles and clean code
- Simplicity and feedback loops (Kent Beck)
- Security and credential handling (DevSecOps)
- CLI ergonomics and user experience
- Brand alignment and product positioning
- Domain-specific rules (Cosign v2, AGPL, API schema, DHCP scope validation)

**Mob consensus gates merge**: Code is ready for push only after mob approval.

### Stage 3: CI/Release (GitHub Actions, No Local Intervention)

After push, GitHub Actions automatically runs:
1. **.github/workflows/ci.yml** (main branch push or PR)
   - Checkout, setup Go, unit tests, Cosign install, GoReleaser snapshot
2. **.github/workflows/release.yml** (on version tags v*)
   - Checkout, setup Go, unit tests, Cosign install, GoReleaser full release

Developer does not run these locally; CI/CD handles release pipeline.

---

## Workflow Loop Diagram

```
Developer codes
    ↓
───────────────── Fast Inner Loop (Repeat until working) ───────────────
│                                                                         │
├─ claude pab test → ✅ Tests pass / ❌ Fail? → Fix & repeat            │
├─ claude pab compile → ✅ Binary built / ❌ Error? → Fix & repeat      │
└─ claude pab snapshot → ✅ Snapshot valid / ❌ Error? → Fix & repeat   │
    ↓
Code changes complete + Tests pass
    ↓
claude pab review (Mob Review Gate)
    ↓
Mob consensus on approval?
    ├─ ✅ Approved → Ready to push
    ├─ 🟡 Approved with Changes → Developer fixes, runs `claude pab review` again
    └─ ❌ Rejected → Major revision needed, return to coding
    ↓
git push (or git tag for release)
    ↓
GitHub Actions CI/Release (automated, no local intervention)
    ↓
Release published
```

---

## Keyboard Shortcuts (Quick Access)

- **Ctrl+Shift+T**: `claude pab test`
- **Ctrl+Shift+B**: `claude pab compile`
- **Ctrl+Shift+S**: `claude pab snapshot`
- **Ctrl+Shift+R**: `claude pab review`

---

## Key Rules

1. **Builder feedback only during fast loop** — No mob gates on build cycles (tests, compile, snapshot)
2. **Mob gates code changes only** — Invoked before push, after developer is satisfied with tests/compile
3. **Secrets not needed locally** — GitHub Actions handles release secrets (GITHUB_TOKEN, Cosign keyless OIDC)
4. **Podman always** — All Go operations containerized; nothing runs on host system
5. **No git in code** — Builder agent never executes git commands; config changes written to disk only
6. **Domain validation enforced** — Builder validates AGPL metadata, Cosign v2 format, checksum integrity, API schema

---

## Example: Adding a New Feature

1. Create branch: `git checkout -b feat/dynamic-config`
2. Code and test in fast loop:
   ```bash
   # Edit internal/config/models.go
   claude pab test      # ❌ Test fails? Fix
   claude pab compile   # ✅ Binary built
   # Edit internal/tui/ui.go
   claude pab test      # ✅ Tests pass
   claude pab compile   # ✅ Binary built
   claude pab snapshot  # ✅ Snapshot valid
   ```
3. Ready for review: `claude pab review`
   - Mob reviews code changes (architecture, domain, SOLID, UX, branding)
   - Mob requests changes OR approves
4. If changes requested: repeat fast loop + review until approved
5. Push: `git push origin feat/dynamic-config`
6. CI runs automatically; release happens if tagged

---

## Troubleshooting

**Fast loop commands slow (> 20 sec)?**
- Check Podman image pulled locally: `podman images | grep golang`
- If missing, first run will pull from ECR (1-2 min one-time)

**Mob review taking > 3 min?**
- Large diffs take longer; break into smaller commits if possible
- Mob only gates code changes, not build cycles, so time investment is at merge boundary (expected)

**Tests pass locally but fail in CI?**
- CI uses `go test ./...` same as local; check Git diff (uncommitted changes?)
- Verify Podman image version matches CI (golang:1.22-bookworm in both)
