# Darkroom CLI - Command Reference

This document provides a complete reference for all `darkroom` commands, flags, and usage examples.

---

## Table of Contents

- [Darkroom CLI - Command Reference](#darkroom-cli---command-reference)
  - [Table of Contents](#table-of-contents)
  - [Authentication](#authentication)
    - [Login](#login)
  - [Job Commands](#job-commands)
    - [Submit](#submit)
    - [List](#list)
    - [Status](#status)
    - [Cancel](#cancel)
    - [Logs](#logs)
  - [Storage Commands](#storage-commands)
    - [List](#list-1)
    - [Copy](#copy)
    - [Remove](#remove)
    - [Stat](#stat)
    - [Presign](#presign)
    - [Sync](#sync)
    - [Make Bucket](#make-bucket)
  - [Other Commands](#other-commands)
    - [Show Config](#show-config)
    - [Set configuration value](#set-configuration-value)
    - [Version](#version)
  - [Global Flags](#global-flags)

---

## Authentication

### Login

Authenticate with the API and fetch credentials:

```bash
darkroom login
````

* Prompts for username and password.
* Stores encrypted kubeconfig and auth token in `~/.darkroom/config.yaml.enc`.

---

## Job Commands

### Submit

```bash
darkroom job submit <jobName> --image <image> --script "<script>" --cpu <n> --memory <mem>
```

* **Arguments**:

  * `<jobName>`: Name of the job.

* **Flags**:

  * `--image` (required): Docker image for the job.
  * `--script` (required): Command/script to run inside the job.
  * `--cpu` (default 1): Number of CPUs.
  * `--memory` (default 1Gi): Memory for the job.

**Example**:

```bash
darkroom job submit pi-job --image docker.io/6darkroom/jh-darkroom:latest --script "sleep 3600" --cpu 1 --memory 1Gi
```

---

### List

```bash
darkroom job list
```

* Lists jobs submitted by the current user.

* **Flags**:

  * `--completed`: list only completed jobs
  * `--running`: list only running jobs
  * `--failed`: list only failed jobs
  * `--pending`: list only pending jobs


---

### Status

```bash
darkroom job status <jobName>
```

* Shows detailed status and metadata for the job.

---

### Cancel

```bash
darkroom job cancel <jobName>
```

* Cancels a submitted job.

---

### Logs

```bash
darkroom job log <jobName> [--follow|-f] [--tail <N>]
```

* **Flags**:

  * `--follow`, `-f`: Stream logs live.
  * `--tail <N>`: Show last N lines.

**Examples**:

```bash
darkroom job log pi-job --tail 50
darkroom job log pi-job --follow
darkroom job log pi-job --tail 100 --follow
```

---

## Storage Commands

### List

```bash
darkroom storage ls [<path>]
```

* Lists S3 buckets or objects under a path.
* No path: lists buckets.
* With path: lists objects/folders in the bucket.

---

### Copy

```bash
darkroom storage cp <src> <dest> [--recursive]
```

* Copies files between local and remote.
* `--recursive` required for directories.

**Examples**:

```bash
# Local -> Remote
darkroom storage cp ./file.txt eglez/data/file.txt

# Remote -> Local
darkroom storage cp eglez/data/file.txt ./file.txt

# Recursive directory copy
darkroom storage cp ./data eglez/data --recursive
```

---

### Remove

```bash
darkroom storage rm <remotePath> [--recursive]
```

* Deletes a file or folder from remote storage.
* `--recursive` required for directories.

---

### Stat

```bash
darkroom storage stat <remotePath>
```

* Displays metadata for a remote file: size, last modified, bucket, etc.

---

### Presign

```bash
darkroom storage presign <remotePath>
```

* Generates a presigned URL for download or upload.
* Download URLs allow sharing files with external users.
* Upload URLs can be used to send data to remote storage securely.

---

### Sync

**This command in in development.**

```bash
darkroom storage sync <localDir> <remotePath> [--checksum]
```

* Synchronizes local directory to remote.
* `--checksum`: compare checksums to avoid unnecessary uploads.

### Make Bucket

In general users do not have permissions to create buckets.

```bash
darkroom storage mb <bucketName>
```

* Creates a new S3 bucket with the specified name.

---

## Other Commands

### Show Config

```bash
darkroom config show
```

* Prints the current configuration (decrypted).
* Sensitive fields (like `KubeConfig` and `S3AccessToken`) are redacted.

### Set configuration value

```bash
darkroom config set myVar=<value>
```

* Sets a configuration variable to a new value.

### Version

```bash
darkroom version
```

* Displays the current version of the Darkroom CLI, along with Git commit hash and build date.

---

## Global Flags

* `--api-endpoint <URL>`: Override default API endpoint.
* `DARKROOM_DEBUG=true`: Enable debug logging.

---

