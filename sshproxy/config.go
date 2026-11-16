package sshproxy

import (
	"golang.org/x/crypto/ssh"
)

// Config controls the local SSH server and the legacy SSH destination.
//
// HostKeyCallback is mandatory. Use knownhosts.New(...) for normal operation.
// ssh.InsecureIgnoreHostKey() is available for isolated labs only.
type Config struct {
	// ListenAddress is the local address, for example "127.0.0.1:0".
	// An empty value defaults to "127.0.0.1:0".
	ListenAddress string

	// TargetAddress is the legacy SSH server, for example "192.0.2.10:22".
	TargetAddress string

	// HostKeyCallback verifies the SSH host key presented by the target.
	HostKeyCallback ssh.HostKeyCallback

	// HostKeyFile is an optional PEM private key used by this proxy when
	// accepting local SSH clients. If empty, an ephemeral RSA key is generated.
	HostKeyFile string

	// HostKeyPEM optionally supplies the proxy's PEM private key directly.
	// HostKeyFile and HostKeyPEM are mutually exclusive.
	HostKeyPEM []byte

	// LegacyAlgorithms enables older algorithms for the connection to the
	// target. Enable only when the target requires them.
	LegacyAlgorithms bool

	// These lists are appended to the defaults when LegacyAlgorithms is true.
	// They allow deployments to add narrowly-scoped compatibility algorithms.
	KeyExchanges      []string
	Ciphers           []string
	MACs              []string
	HostKeyAlgorithms []string

	// DialTimeout applies to TCP connect and SSH handshake with the target.
	// Values <= 0 use 10 seconds.
	DialTimeout int // seconds
}
