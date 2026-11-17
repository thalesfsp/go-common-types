# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Prevent collection callback reentry and nested-lock deadlocks with snapshots.
- Make binary collection operations and indexed set deletion safe under concurrent mutation.
- Copy `LastN` results and atomically replace decoded ordered-map state.
- Use value equality for sets, preserving distinct formatting-colliding values.
- Reject invalid percentile fractions while retaining the supported 0..1 scale.
- Encode and decode real BSON documents through the existing BSON methods.

### Changed
- Update dependencies and Go/lint tooling; require Go 1.26.11 or newer.
- Set membership uses linear scans with Go equality for comparable values and structural
  equality for noncomparable values; callbacks operate on container snapshots.

## [1.0.0] - 2023-02-08
### Added
- First release.
