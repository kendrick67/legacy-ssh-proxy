"""Paramiko usage with the proxy lifecycle managed by Python."""
import os
import paramiko

from context_manager import legacy_ssh_proxy

TARGET = os.environ.get("LEGACY_SSH_TARGET", "192.0.2.50:22")

with legacy_ssh_proxy(TARGET, known_hosts=os.path.expanduser("~/.ssh/known_hosts")) as (host, port):
    client = paramiko.SSHClient()
    client.load_system_host_keys()
    client.set_missing_host_key_policy(paramiko.RejectPolicy())
    try:
        client.connect(
            hostname=host,
            port=port,
            username=os.environ["SSH_USERNAME"],
            password=os.environ["SSH_PASSWORD"],
            look_for_keys=False,
            allow_agent=False,
        )
        _, stdout, _ = client.exec_command("uname -a")
        print(stdout.read().decode(), end="")
    finally:
        client.close()
