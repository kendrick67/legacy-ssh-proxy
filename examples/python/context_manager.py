"""Run the Go compatibility proxy as a context manager.

The proxy executable must be installed or supplied with LEGACY_SSH_PROXY_BIN.
"""
from __future__ import annotations

import os
import socket
import subprocess
import time
from contextlib import contextmanager
from typing import Iterator


def find_free_tcp_port() -> int:
    # There is a small bind-close-rebind race. For production orchestration,
    # prefer starting the Go library with ListenAddress="127.0.0.1:0".
    with socket.socket() as sock:
        sock.bind(("127.0.0.1", 0))
        return int(sock.getsockname()[1])


@contextmanager
def legacy_ssh_proxy(
    target: str,
    *,
    listen_host: str = "127.0.0.1",
    executable: str | None = None,
    host_key_file: str | None = None,
    known_hosts: str | None = None,
    insecure_ignore_host_key: bool = False,
    startup_timeout: float = 10.0,
) -> Iterator[tuple[str, int]]:
    binary = executable or os.environ.get("LEGACY_SSH_PROXY_BIN", "legacy-ssh-proxy")
    port = find_free_tcp_port()
    command = [binary, "-listen", f"{listen_host}:{port}", "-target", target]
    if host_key_file:
        command += ["-host-key", host_key_file]
    if known_hosts:
        command += ["-known-hosts", known_hosts]
    if insecure_ignore_host_key:
        command += ["-insecure-ignore-host-key"]

    process = subprocess.Popen(command, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE, text=True)
    try:
        deadline = time.monotonic() + startup_timeout
        while time.monotonic() < deadline:
            if process.poll() is not None:
                error = process.stderr.read() if process.stderr else ""
                raise RuntimeError(f"SSH proxy exited during startup: {error}")
            try:
                with socket.create_connection((listen_host, port), timeout=0.2):
                    break
            except OSError:
                time.sleep(0.05)
        else:
            raise TimeoutError("SSH proxy did not become ready")
        yield listen_host, port
    finally:
        process.terminate()
        try:
            process.wait(timeout=5)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait(timeout=5)
