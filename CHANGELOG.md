# Changelog

## [1.6.1] - 2025-07-16

### Changed

* Update modules

## [1.6.0] - 2025-07-13

### Added

* Support Zstandard Seekable Format format for data file.
    * https://github.com/facebook/zstd/blob/dev/contrib/seekable_format/zstd_seekable_compression_format.md

## [1.5.0] - 2025-06-11

### Added

* Support Environment Variables in DSN.

## [1.4.0] - 2025-03-03

### Added

* Support RDS IAM auth.

## [1.3.3] - 2024-11-02

### Added

* Add empty test data check.

## [1.3.2] - 2024-11-02

### Changed

* Use "JSON Lines" instead of "NDJSON".

## [1.3.1] - 2024-11-02

### Changed

* Update database URL example link.

## [1.3.0] - 2024-11-02

### Added

* Allow empty lines.
* Supports comment out.

### Changed

* Close DB,File on initialization error.

## [1.2.1] - 2024-11-01

### Fixed

- Fix unlimited rate (`-r 0`) option.

## [1.2.0] - 2024-10-31

### Changed

- Supports multiple data files.

## [1.1.0] - 2024-08-29

### Added

- Add "--color" option.

## [1.0.4] - 2023-11-27

### Changed

- Update help message.

## [1.0.3] - 2023-11-27

### Changed

- Check stdout tty instead of stdin tty.

## [1.0.2] - 2023-11-27

### Changed

- refactor: Change agent loop flow.
- Add progress.go test.

## [1.0.1] - 2023-11-25

### Changed

- refactor: Change Recorder.errorQueryCount type.

## [1.0.0] - 2023-11-24

### Added

- First stable release.

<!-- cf. https://keepachangelog.com/ -->
