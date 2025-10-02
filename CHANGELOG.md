# Changelog

All notable changes to this project will be documented in this file.


## v0.1

Initial release of Darkroom CLI for storage management and authentication.

### Authentication
- `darkroom login` – authenticate with the remote API

### Storage Management
- `darkroom storage ls <path>` – list buckets or folder contents
- `darkroom storage cp <local> <remote>` – copy files to/from remote storage
- `darkroom storage rm <path>` – delete remote file or directory (`--recursive` supported)
- `darkroom storage stat <path>` – show metadata for remote object

### Utilities
- `darkroom version` – show Darkroom version, Git commit, and build date
- Configuration stored in `~/.darkroom/config.yaml` (encrypted)

