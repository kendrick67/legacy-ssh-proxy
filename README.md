# Legacy SSH Proxy

A lightweight, cross-platform SSH compatibility proxy and Go library.

**Written in Go. Built for modern applications that need to communicate with legacy SSH servers.**

Legacy SSH Proxy bridges modern SSH clients and older SSH servers that require deprecated cryptographic algorithms that may no longer be enabled by default in modern SSH libraries.

It terminates an incoming SSH connection using modern cryptographic defaults and establishes a separate outbound connection to the legacy server using selectively enabled legacy algorithms.

**Compile it yourself, download a precompiled executable, or import it directly into your Go application.**

## Why Legacy SSH Proxy?

Modern SSH libraries and cryptography packages continue to remove or disable older algorithms due to security concerns.

Unfortunately, many network devices, embedded systems, and legacy servers still depend on these algorithms.

This can create compatibility problems for applications using Python Paramiko, OpenSSH, Go SSH libraries, and other modern SSH clients.

Legacy SSH Proxy provides a compatibility layer without requiring applications to weaken their own SSH cryptographic settings.

## How It Works

```text
+-----------------------------+
|      Application            |
|                             |
|  Python / Paramiko          |
|  Go Applications            |
|  OpenSSH / Other SSH Clients|
+--------------+--------------+
               |
               | Modern SSH
               |
               v
+-----------------------------+
|      Legacy SSH Proxy       |
|                             |
|         Written in Go       |
|                             |
|  Local SSH Server           |
|  Authentication Forwarding |
|  SSH Channel Forwarding     |
|  Legacy Algorithm Support   |
+--------------+--------------+
               |
               | Legacy SSH
               |
               v
+-----------------------------+
|      Legacy SSH Server      |
|                             |
|  Older RSA / SHA-1          |
|  Legacy Key Exchanges       |
|  Older Ciphers and MACs     |
+-----------------------------+
```

The proxy terminates both SSH connections independently.

- Applications connect using modern SSH cryptography.
- The proxy establishes a separate connection to the legacy server.
- Legacy algorithms are enabled only on the outbound connection.
- Applications do not need to modify their SSH library configuration.

## Features

- Written in Go.
- Standalone compiled executable.
- Importable Go library.
- Cross-platform support.
- Linux, Windows, and macOS.
- No Python runtime dependency.
- Compatible with Python Paramiko.
- Compatible with OpenSSH clients.
- Local TCP listener.
- Dynamically assigned TCP ports.
- SSH password authentication forwarding.
- Optional legacy private key authentication.
- SSH command execution forwarding.
- Interactive shell forwarding.
- SFTP channel forwarding.
- Configurable legacy cryptographic algorithms.
- Suitable for automation and network engineering applications.

# Installation

There are three ways to use Legacy SSH Proxy:

1. Download a precompiled executable.
2. Compile the application yourself.
3. Import the Go library into your own application.

## 1. Download a Precompiled Executable

You do not need Go installed to use the executable.

Visit the GitHub Releases page and download the binary matching your operating system and CPU architecture.

**[Download Legacy SSH Proxy Releases](https://github.com/kendrick67/legacy-ssh-proxy/releases)**

### Supported Platforms

| Operating System | Architecture | Release Binary |
|---|---|---|
| Linux | AMD64 | `legacy-ssh-proxy-linux-amd64` |
| Linux | ARM64 | `legacy-ssh-proxy-linux-arm64` |
| Windows | AMD64 | `legacy-ssh-proxy-windows-amd64.exe` |
| Windows | ARM64 | `legacy-ssh-proxy-windows-arm64.exe` |
| macOS | AMD64 | `legacy-ssh-proxy-darwin-amd64` |
| macOS | ARM64 | `legacy-ssh-proxy-darwin-arm64` |

### Linux Installation

Download the appropriate binary from Releases.

Install it into `/usr/local/bin`:

```bash
sudo install -m 755 legacy-ssh-proxy-linux-amd64 \
    /usr/local/bin/legacy-ssh-proxy
```

Run:

```bash
legacy-ssh-proxy \
    -listen 127.0.0.1:2222 \
    -target 192.168.1.50:22
```

### Windows Installation

Download the Windows executable and place it in a directory included in your system PATH.

Run from PowerShell:

```powershell
.\legacy-ssh-proxy-windows-amd64.exe `
    -listen 127.0.0.1:2222 `
    -target 192.168.1.50:22
```

### macOS Installation

For Apple Silicon:

```bash
chmod +x legacy-ssh-proxy-darwin-arm64

sudo mv legacy-ssh-proxy-darwin-arm64 \
    /usr/local/bin/legacy-ssh-proxy
```

Download the matching architecture from Releases.

## 2. Compile From Source

Legacy SSH Proxy is written in Go and can be compiled for your current operating system or cross-compiled for other platforms.

### Requirements

- Go toolchain compatible with the repository's `go.mod`.
- Git.

### Clone the Repository

```bash
git clone https://github.com/kendrick67/legacy-ssh-proxy.git

cd legacy-ssh-proxy
```

### Install Dependencies

```bash
go mod tidy
```

### Compile for Your Current Operating System

```bash
go build -o legacy-ssh-proxy ./cmd/legacy-ssh-proxy
```

On Windows:

```powershell
go build -o legacy-ssh-proxy.exe ./cmd/legacy-ssh-proxy
```

### Cross-Compile

Linux AMD64:

```bash
GOOS=linux GOARCH=amd64 go build \
    -o legacy-ssh-proxy-linux-amd64 \
    ./cmd/legacy-ssh-proxy
```

Linux ARM64:

```bash
GOOS=linux GOARCH=arm64 go build \
    -o legacy-ssh-proxy-linux-arm64 \
    ./cmd/legacy-ssh-proxy
```

Windows AMD64:

```bash
GOOS=windows GOARCH=amd64 go build \
    -o legacy-ssh-proxy-windows-amd64.exe \
    ./cmd/legacy-ssh-proxy
```

Windows ARM64:

```bash
GOOS=windows GOARCH=arm64 go build \
    -o legacy-ssh-proxy-windows-arm64.exe \
    ./cmd/legacy-ssh-proxy
```

macOS AMD64:

```bash
GOOS=darwin GOARCH=amd64 go build \
    -o legacy-ssh-proxy-darwin-amd64 \
    ./cmd/legacy-ssh-proxy
```

macOS ARM64:

```bash
GOOS=darwin GOARCH=arm64 go build \
    -o legacy-ssh-proxy-darwin-arm64 \
    ./cmd/legacy-ssh-proxy
```

# Using the Command-Line Application

Start the proxy:

```bash
legacy-ssh-proxy \
    -listen 127.0.0.1:2222 \
    -target 192.168.1.50:22
```

### Command-Line Parameters

| Parameter | Description |
|---|---|
| `-listen` | Local TCP listening address and port |
| `-target` | Legacy SSH server address and port |
| `-key` | Optional private key for backend authentication |

Example using a private key:

```bash
legacy-ssh-proxy \
    -listen 127.0.0.1:2222 \
    -target 192.168.1.50:22 \
    -key /home/user/.ssh/legacy_key
```

### Connect Using OpenSSH

```bash
ssh -p 2222 admin@127.0.0.1
```

Execute a command:

```bash
ssh -p 2222 admin@127.0.0.1 "show version"
```

Connect using SFTP:

```bash
sftp -P 2222 admin@127.0.0.1
```

Transfer files:

```bash
scp -P 2222 file.txt admin@127.0.0.1:/tmp/
```

# Go Library Integration

Legacy SSH Proxy is also available as an importable Go package.

This allows developers to embed the SSH compatibility proxy directly into their own applications without starting an external process.

## Install the Go Library

```bash
go get github.com/kendrick67/legacy-ssh-proxy
```

## Import the Package

```go
import sshproxy "github.com/kendrick67/legacy-ssh-proxy/sshproxy"
```

## Basic Example

Start the proxy directly from your Go application:

```go
package main

import (
    "fmt"
    "log"

    sshproxy "github.com/kendrick67/legacy-ssh-proxy/sshproxy"
)

func main() {

    proxy, err := sshproxy.Start(sshproxy.Config{
        ListenAddress: "127.0.0.1:0",
        TargetAddress: "192.168.1.50:22",
    })

    if err != nil {
        log.Fatal(err)
    }

    defer proxy.Close()

    fmt.Println("SSH Proxy listening on:", proxy.Addr())

    select {}
}
```

Using `127.0.0.1:0` allows the operating system to allocate an available TCP port automatically.

The actual listening address is available through `proxy.Addr()`.

## Public Library API

| API | Description |
|---|---|
| `sshproxy.Config` | Proxy configuration |
| `sshproxy.Start()` | Starts the proxy |
| `proxy.Addr()` | Returns the listening address |
| `proxy.Close()` | Stops the proxy and releases resources |

The library and command-line executable use the same underlying implementation.

See `examples/go/` for additional Go integration examples.

# Python Integration

Python applications can use Legacy SSH Proxy through Paramiko.

A context manager is provided to automatically start and stop the Go executable around an SSH session.

Traditional subprocess integration is also supported.

See the dedicated Python documentation:

**[Python Integration Guide](examples/python/README.md)**

Examples include:

- Python context manager.
- Automatic proxy startup and shutdown.
- Random TCP port allocation.
- Traditional subprocess management.
- Paramiko SSH connections.
- Remote command execution.

# Security Considerations

Legacy SSH Proxy is intended to provide compatibility with systems that cannot immediately be upgraded.

Legacy cryptographic algorithms may have known security weaknesses.

Recommended practices:

- Bind the proxy to `127.0.0.1` whenever possible.
- Enable only the algorithms required by the target server.
- Verify backend SSH host keys.
- Protect private keys and credentials.
- Use a persistent proxy host key.
- Avoid exposing the proxy directly to untrusted networks.
- Never disable host-key verification in production.

Legacy algorithm support should be treated as a compatibility exception, not a replacement for modern SSH security.

# Testing

Run the Go test suite:

```bash
go test ./...
```

# Releases

Precompiled executables are published through GitHub Releases.

**[View All Releases](https://github.com/kendrick67/legacy-ssh-proxy/releases)**

Each release provides binaries for supported operating systems and architectures.

# Contributing

Contributions, bug reports, and improvements are welcome.

If you encounter a legacy SSH server that requires a specific algorithm combination, please consider submitting an issue describing:

- Server operating system or device type.
- SSH server implementation and version, if known.
- Required host key algorithm.
- Required key exchange algorithm.
- Required cipher and MAC.
- Relevant connection error messages.

Never include passwords, private keys, or other sensitive information in an issue.

# License

MIT License.

See [LICENSE](LICENSE) for details.
