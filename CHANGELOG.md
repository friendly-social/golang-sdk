# Changelog

All notable changes to this project will be documented in this file.

The format is based on [**Keep a Changelog**](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [**Semantic Versioning**](https://semver.org/spec/v2.0.0.html).

## [0.3.1] - 22-02-2026

### Added
- Email Locale

## [0.3.0] - 16-02-2026

### Added
- Enhanced typing for primitives
- Typed APIError
- Builder pattern for Client
### Removed
- NewProductionClient() and NewLocalhostClient() (their functionality is now handled via WithBaseURL() builder method)

## [0.2.0] - 10-02-2026

### Added
- Profile editing
- E-mail linking/unlinking
- Login via e-mail
- Replaced `Interests` to `[]Interest`
### Changed
- Renamed Generate to Register

## [0.1.0] - 04-02-2026

### Added
- The whole project!
- This is the starting point.
