package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kendrick67/legacy-ssh-proxy/sshproxy"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	hostKeyCallback, err := knownhosts.New(os.ExpandEnv("$HOME/.ssh/known_hosts"))
	if err != nil {
		log.Fatal(err)
	}

	proxy, err := sshproxy.Start(sshproxy.Config{
		ListenAddress:    "127.0.0.1:0",
		TargetAddress:    "192.0.2.50:22",
		HostKeyCallback:  hostKeyCallback,
		LegacyAlgorithms: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer proxy.Close()

	log.Printf("proxy listening on %s", proxy.Addr())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
