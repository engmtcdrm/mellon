# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added options to create new credential without having to interactive with minno. When the options `-n`, `--cred-name` and `-f`, `--file` are used with the `create` command, it will run without requiring any user input. (#6)

## [v1.1.0] - 2025-03-22

### Added

- Added an option to output a decrypted credential to a file. This only works when the `-n`, `--cred-name` option is used when viewing a credential. If no output is specified, the output will be to stdout. (#3)

## [v1.0.0] - 2024-11-03

Initial release
