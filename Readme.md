# Dutis

A command-line tool to select default applications. It is a wrapper around [duti](https://github.com/moretension/duti).

> **Note**: This is a fork of [mrtkrcm/dutis](https://github.com/mrtkrcm/dutis) with performance improvements, caching, and configuration management features.

## ✨ Features

- 🚀 **Fast Performance**: Instant startup with smart caching (~900x faster)
- 💾 **Configuration Management**: YAML-based config for tracking and bulk-applying associations
- 🎨 **Beautiful UI**: Colored output with clean formatting
- ⚡ **CLI Commands**: Manage associations without interactive mode
- 🔄 **Persistent Cache**: 24-hour cache for applications and recommendations
- 🎯 **Smart Path Cleaning**: Clean, readable application names

## Installation

### From Source (Recommended for this fork)

```shell
go install github.com/tobiashochguertel/dutis@latest
```

### Original Installation Methods

#### Using HomeBrew

```shell
brew tap mrtkrcm/dutis https://github.com/mrtkrcm/dutis
brew install dutis
```

#### Using Go (Original)

```shell
go install github.com/mrtkrcm/dutis@latest
```

## Usage

### Interactive Mode

```shell
dutis
```

Launches interactive TUI to set file associations. All selections are automatically saved to `~/.dutis/config.yaml`.

### CLI Commands

```shell
# List all configured associations
dutis list

# Apply all configured associations (bulk restore)
dutis apply

# Remove a specific association
dutis remove .txt

# Refresh application cache
dutis --refresh-cache

# Show help
dutis help
```

## Configuration

All file associations are stored in `~/.dutis/config.yaml`:

```yaml
version: "1.0"
associations:
  .txt:
    suffix: .txt
    application: Visual Studio Code.app
    bundle_id: com.microsoft.VSCode
    set_at: 2024-11-07T20:00:00Z
```

### Workflows

**Backup your associations**:
```shell
cp ~/.dutis/config.yaml ~/Dropbox/dotfiles/
```

**Restore on new machine**:
```shell
cp ~/Dropbox/dotfiles/config.yaml ~/.dutis/
dutis apply
```

## Performance

| Feature | Before | After | Improvement |
|---------|--------|-------|-------------|
| Startup | ~5-10s | <100ms | 50-100x faster |
| Recommended apps | ~470ms | ~0.5ms | 900x faster |
| Autocomplete | Slow | Instant | Cached |

## Cache

Application data is cached in `~/.cache/dutis/`:
- `uti_cache.gob` - Application list (24h expiry)
- `recommended_apps_cache.gob` - Recommended apps per suffix (24h expiry)

## What's New in This Fork

See [CHANGELOG.md](./CHANGELOG.md) for detailed changes.

**Highlights**:
- ⚡ Instant startup with lazy-loaded caching
- 📝 YAML configuration tracking
- 🎨 Beautiful colored output
- 🔧 CLI commands for automation
- 🚀 900x performance improvement for recommended apps
- 🧹 Clean path display (no more %20 or long paths)
- ✨ Version info display

## Screenshots

1. Waiting for environment checking

![env-check](./images/env-check.png)

1. Selecting suffix

![choose-suffix](./images/choose-suffix.png)

1. Checking recommended applications

![recommend](./images/recommend.png)

1. Selecting application UTI

![choose-uti](./images/choose-uti.png)

1. Finished

![finish](./images/finish.png)

## Stargazers over time

[![Stargazers over time](https://starchart.cc/tobiashochguertel/dutis.svg?variant=adaptive)](https://starchart.cc/tobiashochguertel/dutis)

## Original Project

This is a fork of [mrtkrcm/dutis](https://github.com/mrtkrcm/dutis). See the original project for the base implementation.
