# Darkroom CLI

`darkroom` is a command-line tool for interacting with the Darkroom platform.  
It provides an easy way to **submit and manage jobs on Kubernetes clusters** and to **store, sync, and share files via S3**.

---

## Features

- 🔐 Secure login with encrypted local config.  
- 🚀 Submit, monitor, cancel, and fetch logs from user jobs.  
- 📦 Manage S3 storage (list, copy, remove, stat, presign, sync).  
- 🛠 Simple YAML-based configuration, automatically encrypted.  
- 🧩 Extensible design with Cobra commands.

---

## 📦 Releases

You can always find the latest binaries and release notes here: 👉 [Latest Release Notes](https://github.com/eddienko/darkroom-client/releases/latest)

Once downloaded, make the binary executable:

```bash
chmod +x darkroom
```   

and make it available system-wide, e.g.:

```bash
mv darkroom /usr/local/bin/
```

or in any directory in your `$PATH`.

---

## Quick Start

1. **Login** to Darkroom:

   ```bash
   darkroom login
   ```

   This will securely store your credentials and kubeconfig in
   `~/.darkroom/config.yaml.enc`.

2. **Submit a job**:

   ```bash
   darkroom job submit test-job \
     --image docker.io/6darkroom/jh-darkroom:latest \
     --script "echo hello && sleep 60" \
     --cpu 1 --memory 1Gi
   ```

3. **List jobs**:

   ```bash
   darkroom job list
   ```

4. **Use storage** (S3 backend):

   ```bash
   darkroom storage ls        # list buckets
   darkroom storage ls mybucket/  # list objects
   ```

---

## Cheat Sheet

| Command                                     | Description                                 |
| ------------------------------------------- | ------------------------------------------- |
| `darkroom login`                            | Authenticate and fetch credentials          |
| `darkroom job submit <name> --image ...`    | Submit a job                                |
| `darkroom job list`                         | List submitted jobs                         |
| `darkroom job status <name>`                | Show detailed job status                    |
| `darkroom job cancel <name>`                | Cancel a job                                |
| `darkroom job log <name> [-f --tail N]`     | View or follow job logs                     |
| `darkroom storage ls [path]`                | List buckets or objects                     |
| `darkroom storage cp <src> <dst> [-r]`      | Copy local↔remote files                     |
| `darkroom storage rm <path> [-r]`           | Remove remote files/folders                 |
| `darkroom storage mb <bucket>`              | Create a new bucket                         |
| `darkroom storage stat <path>`              | Show file metadata                          |
| `darkroom storage presign <path>`           | Generate download/upload URL                |
| `darkroom storage sync <localdir> <remote>` | Sync directory to remote (add `--checksum`) |
| `darkroom config show`                      | Show decrypted config (redacted secrets)    |
| `darkroom config set myVar=<value>`         | Set configuration value                     |
| `darkroom version`                          | Show CLI version, Git commit, and build date|

Note that not all commands and flags are currently in production use.
See the latest release notes for details.

---

## Full Command Reference

For detailed usage, arguments, and examples see:
👉 [docs/commands.md](docs/commands.md)

---

## Development

Clone and build:

```bash
git clone https://github.com/your-org/darkroom.git
cd darkroom
go build -o darkroom ./main.go
````

Run with debug enabled:

```bash
DARKROOM_DEBUG=true darkroom job list
```

---

## License

MIT License. See [LICENSE](LICENSE).

---
