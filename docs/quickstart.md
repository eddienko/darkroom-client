# Darkroom CLI - Quickstart

## Acessing the storage

1. Authenticate. You will need your username, password and authenticator code.

```bash
darkroom login
```

2. Inspect available commandds.

```bash
darkroom storage --help
```

3. List buckets

```bash
darkroom storage ls
```

4. List directories

```bash
darkroom storage ls scratch.space/users/eglez
```

4. Copy a file

```bash
darkroom storage cp localfile.txt scratch.space/users/eglez/
```

> 💡 **Tip:** You can alias the `darkroom storage` command to something different if you like, `alias darks="darkroom storage"` and then use `darks ls`

## Submit jobs

1. Authenticate. You will need your username, password and authenticator code.

```bash
darkroom login
```

2. Inspect available commands.

```bash
darkroom jobs --help
```

3. Submit a job to run remotely.

```bash
darkroom job submit --name test-job-1 --script "sleep 600"
```

The script is a command to execute. If it points to a script, the script has to exist in the remote server and have all the necessary lines to e.g. load a virtual environment if needed.

4. List submitted jobs

```bash
darkroom job list
```

5. Retrieve a more detailed status about a job.

```bash
darkroom job status test-job-1
```

6. Display output logs.

```bash
dadkroom job log test-job-1
```

7. Execute a shell inside the running job. This opens a shell in the remote container running the job.

```bash
darkroom job shell test-job-1
```

8. Cancel/remove a job.

```bash
darkroom job cancel test-job-1
```

> 💡 **Tip:** You can alias the `darkroom job` command to something different if you like, `alias darkj="darkroom job"` and then use `darkj list`