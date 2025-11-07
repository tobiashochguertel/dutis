# Release Process

## Overview

This project uses a simplified release workflow:
- **CI Workflow**: Runs automatically on every push to validate code
- **Release Workflow**: Manual trigger to create releases

## Version Management

The version is defined in **one place**: `main.go`

```go
const (
    Version = "v0.2.1-fork"
    // ...
)
```

## How to Release

### Option 1: Using GitHub UI (Recommended)

1. **Update version in code**:
   ```bash
   # Edit main.go and update Version constant
   vim main.go
   
   # Update CHANGELOG.md
   vim CHANGELOG.md
   
   # Commit
   git add main.go CHANGELOG.md
   git commit -m "chore: Bump version to v0.3.0-fork"
   git push
   ```

2. **Trigger release via GitHub**:
   - Go to: https://github.com/tobiashochguertel/dutis/actions/workflows/release-manual.yml
   - Click "Run workflow"
   - Leave version empty (uses version from main.go)
   - Check "Create git tag"
   - Click "Run workflow"

### Option 2: Using gh CLI

```bash
# Update version in main.go first
vim main.go

# Commit changes
git add main.go CHANGELOG.md
git commit -m "chore: Bump version to v0.3.0-fork"
git push

# Trigger release (uses version from main.go)
gh workflow run release-manual.yml

# Or specify custom version
gh workflow run release-manual.yml -f version=v0.3.0-fork -f create_tag=true
```

### Option 3: Manual Tag (Advanced)

```bash
# Update version in main.go
vim main.go

# Commit and push
git add main.go
git commit -m "chore: Bump version to v0.3.0-fork"
git push

# Create and push tag
git tag -a v0.3.0-fork -m "Release v0.3.0-fork"
git push origin v0.3.0-fork

# This automatically triggers the release workflow
```

## Workflow Behavior

### CI Workflow (`.github/workflows/ci.yml`)
- **Triggers**: Every push to `main`, every PR
- **Purpose**: Build, test, and verify code
- **Does NOT**: Create releases or tags
- **Ignores**: Markdown and image changes

### Release Workflow (`.github/workflows/release-manual.yml`)
- **Triggers**: Manual (workflow_dispatch) or git tag push
- **Purpose**: Create GitHub release with binaries
- **Options**:
  - Use version from `main.go` (default)
  - Override with custom version
  - Optionally create git tag
- **Outputs**: GitHub release with binaries for macOS

## Version Numbering

We use Semantic Versioning with `-fork` suffix:

- `v0.2.1-fork` - Current version
- `v0.3.0-fork` - Next minor version (new features)
- `v0.2.2-fork` - Next patch version (bug fixes)
- `v1.0.0-fork` - Major version (breaking changes)

Format: `v<MAJOR>.<MINOR>.<PATCH>-fork`

## Checklist for Releases

- [ ] Update `Version` constant in `main.go`
- [ ] Update `CHANGELOG.md` with changes
- [ ] Update `README.md` if needed
- [ ] Commit changes with message: `chore: Bump version to vX.Y.Z-fork`
- [ ] Push to main
- [ ] Wait for CI to pass (green checkmark)
- [ ] Trigger release workflow via GitHub UI or `gh`
- [ ] Verify release assets are created
- [ ] Test installation: `go install github.com/tobiashochguertel/dutis@vX.Y.Z-fork`
- [ ] Update release notes on GitHub if needed

## Testing Releases

After release:

```bash
# Install specific version
go install github.com/tobiashochguertel/dutis@v0.3.0-fork

# Verify version
dutis version

# Test functionality
dutis help
```

## Troubleshooting

### "Tag already exists"
- Delete tag: `gh release delete vX.Y.Z-fork --yes && git tag -d vX.Y.Z-fork && git push origin :refs/tags/vX.Y.Z-fork`
- Try release again

### "Module path mismatch"
- Ensure `go.mod` has: `module github.com/tobiashochguertel/dutis`
- Ensure imports use: `github.com/tobiashochguertel/dutis/util`

### CI fails
- Check workflow logs: `gh run view --log-failed`
- Fix issues and push again
- Release only when CI is green ✓

## Benefits of This Approach

✅ **Single source of truth** - Version in one place (main.go)  
✅ **Clear separation** - CI validates, Release publishes  
✅ **Manual control** - No accidental releases  
✅ **Flexible** - Can override version if needed  
✅ **Safe** - Won't create release on every push  
✅ **Simple** - Easy to understand and use  
