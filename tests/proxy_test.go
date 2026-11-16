package tests

import (
	"net"
	"testing"
	"time"

	"github.com/kendrick67/legacy-ssh-proxy/sshproxy"
	"golang.org/x/crypto/ssh"
)

func TestStartRequiresTargetAndHostKeyCallback(t *testing.T) {
	_, err := sshproxy.Start(sshproxy.Config{})
	if err == nil {
		t.Fatal("expected missing target error")
	}

	_, err = sshproxy.Start(sshproxy.Config{TargetAddress: "127.0.0.1:22"})
	if err == nil {
		t.Fatal("expected missing host key callback error")
	}
}

func TestStartAndClose(t *testing.T) {
	p, err := sshproxy.Start(sshproxy.Config{
		ListenAddress:   "127.0.0.1:0",
		TargetAddress:   "127.0.0.1:22",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if p.Addr() == nil {
		t.Fatal("expected listener address")
	}
	conn, err := net.DialTimeout("tcp", p.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}
}
