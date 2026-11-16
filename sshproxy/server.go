package sshproxy

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

type sshCredentials struct {
	username string
	password string
}

func (p *Proxy) handleConnection(raw net.Conn) {
	var credentials sshCredentials
	serverConfig := &ssh.ServerConfig{
		PasswordCallback: func(meta ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if len(meta.User()) == 0 || len(password) == 0 {
				return nil, fmt.Errorf("username and password are required")
			}
			credentials = sshCredentials{username: meta.User(), password: string(password)}
			return &ssh.Permissions{}, nil
		},
	}
	serverConfig.AddHostKey(p.signer)

	serverConn, chans, reqs, err := ssh.NewServerConn(raw, serverConfig)
	if err != nil {
		return
	}
	defer serverConn.Close()

	backend, err := p.dialBackend(credentials)
	if err != nil {
		log.Printf(
			"ERROR: backend SSH connection failed: target=%s user=%s error=%v",
			p.cfg.TargetAddress,
			credentials.username,
			err,
		)
		return
	}

	log.Printf("Connected to backend SSH server: %s", p.cfg.TargetAddress)
	defer backend.Close()

	// golang.org/x/crypto/ssh exposes incoming global requests on the
	// server side, but ssh.Client consumes backend global requests internally
	// and does not expose them for forwarding. Reject client global requests.
	go ssh.DiscardRequests(reqs)

	var channelWorkers = make(chan struct{}, 64)
	for newChannel := range chans {
		channelWorkers <- struct{}{}
		go func(ch ssh.NewChannel) {
			defer func() { <-channelWorkers }()
			p.forwardChannel(ch, backend)
		}(newChannel)
	}
}

func (p *Proxy) dialBackend(credentials sshCredentials) (*ssh.Client, error) {
	timeout := time.Duration(p.cfg.DialTimeout) * time.Second
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	raw, err := net.DialTimeout("tcp", p.cfg.TargetAddress, timeout)
	if err != nil {
		return nil, fmt.Errorf("connect to target: %w", err)
	}
	_ = raw.SetDeadline(time.Now().Add(timeout))

	cfg := &ssh.ClientConfig{
		User:            credentials.username,
		Auth:            []ssh.AuthMethod{ssh.Password(credentials.password)},
		HostKeyCallback: p.cfg.HostKeyCallback,
		Timeout:         timeout,
	}
	if p.cfg.LegacyAlgorithms {
		cfg.Config.KeyExchanges = appendUnique(
			cfg.Config.KeyExchanges, "diffie-hellman-group14-sha1", "diffie-hellman-group1-sha1",
		)
		cfg.Config.Ciphers = appendUnique(
			cfg.Config.Ciphers, "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-cbc",
			"aes192-cbc", "aes256-cbc", "3des-cbc",
		)
		cfg.Config.MACs = appendUnique(
			cfg.Config.MACs, "hmac-sha1", "hmac-sha1-96",
		)
		cfg.HostKeyAlgorithms = appendUnique(
			cfg.HostKeyAlgorithms, "ssh-rsa",
		)
	}
	cfg.Config.KeyExchanges = appendUnique(cfg.Config.KeyExchanges, p.cfg.KeyExchanges...)
	cfg.Config.Ciphers = appendUnique(cfg.Config.Ciphers, p.cfg.Ciphers...)
	cfg.Config.MACs = appendUnique(cfg.Config.MACs, p.cfg.MACs...)
	cfg.HostKeyAlgorithms = appendUnique(cfg.HostKeyAlgorithms, p.cfg.HostKeyAlgorithms...)

	conn, chans, reqs, err := ssh.NewClientConn(raw, p.cfg.TargetAddress, cfg)
	if err != nil {
		_ = raw.Close()
		return nil, fmt.Errorf("SSH handshake with target: %w", err)
	}
	_ = raw.SetDeadline(time.Time{})
	return ssh.NewClient(conn, chans, reqs), nil
}

func appendUnique(dst []string, values ...string) []string {
	seen := make(map[string]struct{}, len(dst)+len(values))
	for _, v := range dst {
		seen[v] = struct{}{}
	}
	for _, v := range values {
		if _, ok := seen[v]; !ok {
			dst = append(dst, v)
			seen[v] = struct{}{}
		}
	}
	return dst
}
