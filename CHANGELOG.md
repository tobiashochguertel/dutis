# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.2.0-fork] - 2024-11-07

### Added
- **Version Display**: Shows fork version and repository URL on startup
- **Cache Status Indicator**: Displays cache readiness at startup (✓ or ○)
- **Persistent Application Cache**: 
  - Caches scanned applications for 24 hours in `~/.cache/dutis/uti_cache.gob`
  - Reduces startup time from ~5-10s to <100ms
  - Lazy-loaded on first use for instant startup
- **Recommended Apps Cache**:
  - Per-suffix caching of recommended applications
  - 900x performance improvement (~470ms → ~0.5ms)
  - Stored in `~/.cache/dutis/recommended_apps_cache.gob`
- **YAML Configuration System**:
  - Tracks all file associations in `~/.dutis/config.yaml`
  - Auto-saves associations during interactive mode
  - Human-readable YAML format with timestamps
- **CLI Commands**:
  - `dutis list` - List all configured associations
  - `dutis apply` - Bulk apply all associations (great for new machine setup)
  - `dutis remove <suffix>` - Remove specific association
  - `dutis --refresh-cache` - Force refresh application cache
  - `dutis help` - Show help message
- **Double CTRL+C Exit**: Press CTRL+C twice quickly to exit

### Improved
- **Recommended Applications Display**:
  - Beautiful colored output with alternating cyan/blue bullets
  - Unicode box-drawing characters (─) for headers
  - Application count display
  - URL decoded paths (no more %20)
  - Removed trailing slashes from app names
  - Clean path display (removes system prefixes)
  - Preserves utility folder context (e.g., `Utilities/Script Editor.app`)
  - Automatic deduplication of apps
- **Error Handling**: Better error messages and graceful fallbacks
- **User Experience**: Visual feedback for all operations with colored status messages

### Fixed
- Double CTRL+C implementation (moved counter to global scope)
- Long application paths now cleaned and shortened
- URL-encoded characters in app names properly decoded
- Duplicate applications in recommended list removed

### Technical
- **Dependencies Added**:
  - `gopkg.in/yaml.v3` - YAML config file support
- **New Files**:
  - `util/cache.go` - Caching system for UTI map and recommended apps
  - `util/config.go` - YAML configuration management
- **Modified Files**:
  - `main.go` - Added CLI commands, version display, improved UI
  - `util/uti.go` - Path cleaning, caching integration, deduplication
  - `go.mod`, `go.sum` - Updated dependencies

### Performance
- **Startup Time**: ~5-10s → <100ms (instant with cache)
- **Recommended Apps**: ~470ms → ~0.5ms with cache (900x faster)
- **Memory Usage**: Reduced via lazy loading

### Notes
This is a fork of [mrtkrcm/dutis](https://github.com/mrtkrcm/dutis) with significant enhancements for performance, usability, and configuration management.

---

## [Original] - Prior versions

See [upstream repository](https://github.com/mrtkrcm/dutis) for original version history.
