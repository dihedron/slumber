# Changelog

All notable changes to this project will be documented in this file.

## [0.0.1] - 2025-12-21

### Added
- Dynamic logging configuration via environment variables (`<APP>_LOG_LEVEL`, `<APP>_LOG_STREAM`).
- CPU and memory profiling control via environment variables.

### Changed
- Improved security by using `path.Clean` for file path operations.
- Updated `monitor` command to handle signals and exit gracefully.

### Fixed
- Potential security vulnerabilities (G304: file path as variable).

## [0.0.0] - 2025-12-20

### Added
- Initial project structure following `netcheck` patterns (Makefile, `.mk` files, `.goreleaser.yaml`).
- `init.go` for logging (`slog`) and profiling (`pprof`) support.
- `monitor` command implementation using `jessevdk/go-flags`.
- `version` command to show application metadata.
- Process detection logic for VS Code, Antigravity, Windsurf, Cursor, and Zed.
- Support for detecting editor servers running via interpreters (node, python, sh, etc.).
- Active SSH connection check via `/proc/net/tcp` to exclude "stuck" editor processes.
- Power management logic (shutdown/hibernate) via DBus (`systemd-logind`).
- Comprehensive unit tests for detection logic with mocked `/proc` and network files.

### Changed
- Refactored `IsAnyEditorActive` to return only `[]string` for simplicity.
- Refactored process scanning to avoid false positives from command-line flag values.
- Updated `isPID` to use `strconv.Atoi` for more robust validation.
- Replaced deprecated `ioutil.ReadDir` with `os.ReadDir`.
