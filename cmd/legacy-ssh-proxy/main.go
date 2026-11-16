package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kendrick67/legacy-ssh-proxy/sshproxy"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	listen := flag.String("listen", "", "local listen address (ie. '127.0.0.1:2222')")
	target := flag.String("target", "", "legacy SSH target host:port")
	hostKey := flag.String("host-key", "", "proxy host private key PEM file (optional)")
	knownHosts := flag.String("known-hosts", defaultKnownHosts(), "known_hosts file for target host verification")
	insecure := flag.Bool("insecure-ignore-host-key", false, "disable target host-key verification (INSECURE; lab use only)")
	legacy := flag.Bool("legacy-algorithms", true, "enable legacy algorithms for the target connection")
	flag.Parse()

	if *listen == "" {
		flag.Usage()
		os.Exit(2)
	}

	if *target == "" {
		flag.Usage()
		os.Exit(2)
	}

	var callback ssh.HostKeyCallback
	if *insecure {
		log.Print("WARNING: target host-key verification is DISABLED")
		callback = ssh.InsecureIgnoreHostKey()
	} else {
		cb, err := knownhosts.New(*knownHosts)
		if err != nil {
			log.Fatalf("load known_hosts %q: %v (use -insecure-ignore-host-key only for isolated testing)", *knownHosts, err)
		}
		callback = cb
	}

	proxy, err := sshproxy.Start(sshproxy.Config{
		ListenAddress:    *listen,
		TargetAddress:    *target,
		HostKeyFile:      *hostKey,
		HostKeyCallback:  callback,
		LegacyAlgorithms: *legacy,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("SSH compatibility proxy listening on %s -> %s", proxy.Addr(), *target)
	log.Printf("Press Ctrl+C to stop")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	if err := proxy.Close(); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func defaultKnownHosts() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".ssh", "known_hosts")
	}
	return filepath.Join(home, ".ssh", "known_hosts")
}
