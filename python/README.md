
# Legacy SSH Proxy for Python

A lightweight, cross-platform SSH compatibility proxy and Python wrapper. It terminates incoming SSH connections from modern Python clients (such as Paramiko) using modern cryptographic defaults and establishes a separate outbound connection to legacy SSH servers using selectively enabled legacy algorithms.
The precompiled Go binary for your operating system and CPU architecture is bundled inside the wheel—no Go toolchain or manual binary download is required.

## Installation

```bash
pip install legacy-ssh-proxy
```

If you plan to connect through the proxy using Paramiko, install paramiko alongside it:

```bash
pip install legacy-ssh-proxy paramiko
```

## Quick Start: Context Manager

Use the built-in legacy_ssh_proxy context manager to automatically start the bundled proxy on a free local TCP port, yield the (host, port) pair to Paramiko, and terminate the background process when the block exits.
```python
import paramiko
from legacy_ssh_proxy import legacy_ssh_proxy

with legacy_ssh_proxy("192.0.2.50:22", insecure_ignore_host_key=True, known_hosts="/home/me/.ssh/known_hosts") as (host, port):
    client = paramiko.SSHClient()
    client.load_system_host_keys()
    client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    try:
        client.connect(
            hostname=host, port=port, username="your-user",
            password="your-password", look_for_keys=False, allow_agent=False,
        )
        _, stdout, _ = client.exec_command("uname -a")
        print(stdout.read().decode())
    finally:
        client.close()
```

### Context Manager Parameters

* target (str): Destination legacy SSH server address (host:port).
* listen_host (str): Local bind address (defaults to "127.0.0.1").
* host_key_file (str | None): Path to a persistent SSH host key for the local proxy listener.
* known_hosts (str | None): Path to the known_hosts file used to verify the destination server.
* insecure_ignore_host_key (bool): Skip destination host-key verification. Use only in disposable lab environments.
* executable (str | None): Optional path to a custom legacy-ssh-proxy binary (overrides LEGACY_SSH_PROXY_BIN and the bundled executable).
* startup_timeout (float): Maximum seconds to wait for the local listener to accept connections (defaults to 10.0).

## CLI Usage

Installing the package also adds the legacy-ssh-proxy command to your environment:
```bash
legacy-ssh-proxy \
 -listen 127.0.0.1:2222 \
  -target 192.0.2.50:22 \
  -insecure-ignore-host-key \
  -known-hosts ~/.ssh/known_hosts
```
You can then point any standard SSH, SFTP, or Python Paramiko client at 127.0.0.1:2222.

## Security & Host Keys

* Destination verification: The proxy verifies the target server against the known_hosts file passed via known_hosts or -known-hosts. Avoid setting insecure_ignore_host_key=True in production, as it exposes the connection to man-in-the-middle attacks.
* Local proxy verification: Paramiko separately verifies the proxy's local-facing host key. Because the proxy generates an ephemeral host key on each start by default, pass host_key_file (or -host-key) and pin its public key in your Paramiko known_hosts database when using RejectPolicy().
* Avoid AutoAddPolicy(): Do not use paramiko.AutoAddPolicy() in production environment

