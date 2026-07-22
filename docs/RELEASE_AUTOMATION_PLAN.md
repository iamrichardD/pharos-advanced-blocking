# Release Automation Plan for v0.3.0+

## Overview
Automate the release workflow to trigger on git tag creation (format: `v*.*.*`), generating releases with versioning from commit messages, auto-generated release notes, and synced marketing artifacts.

---

## 1. Release Trigger & Versioning

**Trigger:** Git tag creation matching pattern `v*.*.*`

**Example:**
```bash
git tag -a v0.3.0 -m "Release message"
git push origin v0.3.0  # Triggers GitHub Actions workflow
```

**Workflow File:** `.github/workflows/release.yml`

**Actions:**
1. Detect version from tag (e.g., `v0.3.0` → `0.3.0`)
2. Validate tag format and version numbering
3. Trigger release pipeline

---

## 2. Release Notes Generation

**Source:** Commit messages using [Conventional Commits](https://www.conventionalcommits.org/)

**Format:**
```
feat: Add new feature           → Features section
fix: Bug fix                    → Fixes section
docs: Documentation            → Docs section
chore: Maintenance             → Chores section
BREAKING CHANGE: description   → Breaking Changes section (bold, top)
```

**Tools to Consider:**
- `conventional-changelog`: Auto-generate from git history
- GitHub API: Extract PR descriptions and labels
- Manual override: Option to provide custom release notes

**Automation:**
1. Parse git log since last release
2. Group commits by type
3. Format as markdown with sections
4. Include migration guide if breaking changes detected

---

## 3. Build Artifacts

**Current Setup:**
- GoReleaser: Generates linux/arm64 binaries, .deb, .tar.gz, checksums
- Cosign: Creates `.sigstore.json` signature bundles
- Manual installer script: At iamrichardd.com

**Questions for Spike #4:**
- Should we add Docker image builds?
- Homebrew formula auto-generation?
- AUR (Arch User Repository) package?
- Snapcraft snap package?
- Native Windows .exe builds?
- How to auto-update installer.sh on release?

**Current Workflow:**
```
Tag push → GitHub Actions → GoReleaser builds binaries → 
Cosign signs → Artifacts uploaded to GitHub Releases
```

---

## 4. Testing & Validation Gates

**Must Pass Before Release:**

1. **Test Suite** ✅
   ```bash
   go test ./... -v
   ```
   - Fail release if tests don't pass
   - Generate test coverage report

2. **Linting & Security** (Spike #3)
   ```bash
   golangci-lint run
   gosec ./...
   trivy fs .
   ```
   - Configurable severity thresholds
   - Report before release (pass/fail decision)

3. **Binary Validation**
   - Verify binary compiles
   - Run smoke tests (pab --help, pab --version)
   - Verify signatures validate

---

## 5. Marketing Sync Options

### Option A: Auto-Deploy Website on Release
**Pros:** Latest docs always available, no manual step
**Cons:** Risk of deploying broken marketing site, need to verify build

**Implementation:**
```
Tag push → Tests pass → Build binaries → Update website → Deploy
```

### Option B: Trigger Manual Marketing Review
**Pros:** Review before deploying, control deployment timing
**Cons:** Extra manual step, could delay release

**Implementation:**
```
Tag push → Tests pass → Build binaries → 
Create PR with updated docs → Await approval → Merge & deploy
```

### Option C: Create Release, Deploy Website Separately
**Pros:** Independent workflows, flexibility
**Cons:** Risk of drift between release and website

**Implementation:**
```
Tag push → Release workflow (code only)
Separate workflow: Manual trigger to update & deploy website
```

**Recommendation:** Option A (auto-deploy) if tests include website build verification

### Announcement Options:
1. **GitHub Release Notes** (auto-generated)
2. **Twitter/Social** (requires API key, manual or template)
3. **Email digest** (optional)
4. **CHANGELOG.md** (auto-updated)
5. **Blog post** (manual, high-effort)

---

## 6. Documentation Updates

### Auto-Update CHANGELOG.md

**Format:**
```markdown
## [0.3.0] - 2026-07-21

### Added
- pab list-nodes command
- Active-active deployment support

### Changed
- secrets.json schema: map → array (BREAKING)
- CLI flag: --yes → -f (BREAKING)

### Fixed
- Confirmation prompt now single instead of per-node

[Link to release](https://github.com/iamrichardD/pharos-advanced-blocking/releases/tag/v0.3.0)
```

**Tool:** `conventional-changelog-cli` or custom script

### Auto-Update README.md Version References

**Current refs:**
- Installation link to latest release
- Version badge
- Quickstart example commands

**Automation:**
1. Find all version references (grep)
2. Replace with new version number
3. Validate links work
4. Commit update

---

## 7. GitHub Release Management

### Auto-Create Release Notes
```bash
gh release create v0.3.0 \
  --title "v0.3.0: Active-active Deployment" \
  --notes "$(generate-release-notes)"
```

### Auto-Mark as Latest
```bash
gh release edit v0.3.0 --latest
```

### Pre-Release Handling: OPTIONS

**Option A: No pre-releases**
- Every tag is a full release
- Simple, no special handling

**Option B: Support rc (release candidate) tags**
- Pattern: `v0.3.0-rc.1` → marked as pre-release
- Allows testing before full release
- Requires documentation

**Option C: Manual pre-release flag**
- User provides `--prerelease` flag when tagging
- Flexible but requires discipline

**Questions for Panel:**
- Do you want release candidates before full releases?
- How should rc versions appear in installer/docs?
- What's the promotion path from rc to final?

---

## 8. Proposed Workflow

```
Developer: git tag -a v0.3.0 -m "feat: ..." && git push origin v0.3.0
                    ↓
GitHub Actions triggers on tag
                    ↓
1. Validate tag format (v*.*.*)
2. Run all tests (fail if any fail)
3. Run linting/security (report results)
4. Build with GoReleaser
5. Sign with Cosign
6. Generate release notes from commits
7. Create GitHub Release
8. Mark as Latest
9. Update CHANGELOG.md
10. Update README.md version refs
11. Rebuild marketing website
12. Deploy website (or create PR for review)
13. Post release announcement (optional)
                    ↓
Release published, website updated, GitHub Release created
```

---

## 9. GitHub Actions Workflow Skeleton

```yaml
name: Release

on:
  push:
    tags:
      - 'v[0-9]+.[0-9]+.[0-9]+'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test ./...
      - run: go fmt ./...
      - run: golangci-lint run  # After Spike #3
      
  build:
    needs: validate
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: goreleaser/goreleaser-action@v4
      - uses: sigstore/cosign-installer@v3
      
  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: gh release create ${{ github.ref_name }} ...
      - run: gh release edit ${{ github.ref_name }} --latest
      
  marketing:
    needs: release
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: npm run build (marketing site)
      - run: deploy to hosting
```

---

## 10. Questions for Review Panel

**Questions Requiring Panel Discussion:**

1. **Pre-release Strategy**
   - Should we use release candidates (v0.3.0-rc.1)?
   - How should rc versions be handled in docs/installer?
   - Promotion path from rc to final?

2. **Marketing Announcements**
   - Which channels? (GitHub Release, Twitter, Email, Blog?)
   - Manual or templated?
   - Who approves/publishes?

3. **Website Auto-Deploy Risk**
   - Accept auto-deploy on every tag?
   - Or require manual approval?
   - Include website tests as release gate?

4. **Installer Script Updates**
   - Should installer.sh auto-update on release?
   - How to notify users of breaking changes (secrets.json, flags)?
   - Version pinning option in installer?

5. **Breaking Change Protocol**
   - How should we communicate breaking changes?
   - Migration guides mandatory for breaking releases?
   - How long to support old versions?

6. **Release Workflow Complexity**
   - Is full automation too complex initially?
   - Should we phase it? (e.g., Step 1: release creation, Step 2: add testing gates, Step 3: add marketing sync)
   - Manual overrides for edge cases?

---

## Implementation Timeline

**Phase 1 (Essential):** Tag-based trigger + GoReleaser + GitHub Release creation
**Phase 2 (Important):** Linting/security gates + auto-update CHANGELOG
**Phase 3 (Nice-to-have):** Marketing sync + announcement automation
**Phase 4 (Future):** Additional artifact formats (Docker, Homebrew, etc.)

---

## Open Decisions

1. Pre-release handling strategy
2. Marketing announcement channels and frequency
3. Website auto-deploy approval requirements
4. Installer script update automation
5. Breaking change communication protocol
6. Implementation phasing approach

**Request:** Please submit to review panel for recommendations and solutions.
