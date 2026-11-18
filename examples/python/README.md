# Python / Paramiko examples

These examples use the Go proxy as a local SSH endpoint. Paramiko connects to
the local proxy using modern algorithms; the Go process handles the legacy
negotiation with the destination.

## Requirements

```bash
python -m pip install -r requirements.txt
```

Build or download the `legacy-ssh-proxy` executable and put it on `PATH`, or set
`LEGACY_SSH_PROXY_BIN` to its full path.

## Option 1: Run the proxy separately

Start the proxy:

```bash
legacy-ssh-proxy \
  -listen 127.0.0.1:2222 \
  -target 192.0.2.50:22 \
  -legacy-algorithms=true \
  -insecure-ignore-host-key 
```

Then run [`traditional.py`](traditional.py):

```bash
export SSH_USERNAME='your-user'
export SSH_PASSWORD='your-password'
python traditional.py
```

## Option 2: Manage the proxy from Python

See [`with_context_manager.py`](with_context_manager.py). The helper starts the
proxy, waits for the local TCP listener, yields its address, and stops the
process when the context exits.

```python
from context_manager import legacy_ssh_proxy
import paramiko

with legacy_ssh_proxy(
    "192.0.2.50:22",
    insecure_ignore_host_key=True,
    known_hosts="/home/me/.ssh/known_hosts",
) as (host, port):
    client = paramiko.SSHClient()
    client.load_system_host_keys()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(
            host,
            port=port,
            username="your-user",
            password="your-password",
            look_for_keys=False,
            allow_agent=False,
        )
        _, stdout, _ = client.exec_command("uname -a")
        print(stdout.read().decode())
    finally:
        client.close()
```

Set `insecure_ignore_host_key=True` only in a disposable lab. It disables
verification of the destination server and exposes the connection to
man-in-the-middle attacks.

## Host keys

The proxy verifies the destination host against the `known_hosts` file passed
to its CLI. Paramiko separately verifies the proxy's local-facing host key.
Because the proxy uses a generated host key by default, a fresh Paramiko client
may need an explicitly managed/pinned local proxy key if strict host-key
verification is enabled. For a persistent proxy key, pass `-host-key` to the
CLI and manage its public key in the Paramiko known-hosts database.

Do not use `AutoAddPolicy()` in production: it trusts an unknown host key
without verifying its identity.
