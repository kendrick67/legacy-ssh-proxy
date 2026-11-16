package sshproxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Proxy is a running local SSH compatibility proxy.
type Proxy struct {
	cfg      Config
	listener net.Listener
	signer   ssh.Signer

	mu      sync.Mutex
	closed  bool
	conns   map[net.Conn]struct{}
	workers sync.WaitGroup
}

// Start starts accepting local SSH connections and returns immediately.
// Incoming clients authenticate with a username/password; those credentials
// are used to authenticate to TargetAddress.
func Start(cfg Config) (*Proxy, error) {
	if cfg.TargetAddress == "" {
		return nil, errors.New("TargetAddress is required")
	}
	if cfg.HostKeyCallback == nil {
		return nil, errors.New("HostKeyCallback is required; configure known_hosts verification")
	}
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = "127.0.0.1:0"
	}
	if cfg.HostKeyFile != "" && len(cfg.HostKeyPEM) != 0 {
		return nil, errors.New("HostKeyFile and HostKeyPEM cannot both be set")
	}

	signer, err := loadSigner(cfg)
	if err != nil {
		return nil, fmt.Errorf("load proxy host key: %w", err)
	}
	ln, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		return nil, fmt.Errorf("listen on %q: %w", cfg.ListenAddress, err)
	}

	p := &Proxy{
		cfg:      cfg,
		listener: ln,
		signer:   signer,
		conns:    make(map[net.Conn]struct{}),
	}
	p.workers.Add(1)
	go p.acceptLoop()
	return p, nil
}

// Addr returns the actual local listener address. When ListenAddress uses
// port 0, this contains the operating-system-assigned port.
func (p *Proxy) Addr() net.Addr {
	return p.listener.Addr()
}

// Close stops accepting new connections and closes active client connections.
// It is safe to call Close more than once.
func (p *Proxy) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	err := p.listener.Close()
	for conn := range p.conns {
		_ = conn.Close()
	}
	p.mu.Unlock()
	p.workers.Wait()
	if errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

func loadSigner(cfg Config) (ssh.Signer, error) {
	var key []byte
	switch {
	case len(cfg.HostKeyPEM) != 0:
		key = cfg.HostKeyPEM
	case cfg.HostKeyFile != "":
		b, err := os.ReadFile(cfg.HostKeyFile)
		if err != nil {
			return nil, err
		}
		key = b
	default:
		privateKey, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return nil, err
		}
		der := x509.MarshalPKCS1PrivateKey(privateKey)
		key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})
	}
	return ssh.ParsePrivateKey(key)
}

func (p *Proxy) acceptLoop() {
	defer p.workers.Done()
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			p.mu.Lock()
			closed := p.closed
			p.mu.Unlock()
			if closed || errors.Is(err, net.ErrClosed) {
				return
			}
			// A transient accept error should not terminate the proxy.
			time.Sleep(100 * time.Millisecond)
			continue
		}
		log.Printf("New connection from %s, Attempting a connection to >> %s", conn.RemoteAddr(), p.cfg.TargetAddress)

		p.mu.Lock()
		if p.closed {
			p.mu.Unlock()
			_ = conn.Close()
			return
		}
		p.conns[conn] = struct{}{}
		p.workers.Add(1)
		p.mu.Unlock()

		go func() {
			defer p.workers.Done()
			defer func() {
				p.mu.Lock()
				delete(p.conns, conn)
				p.mu.Unlock()
			}()
			p.handleConnection(conn)
		}()
	}
}
