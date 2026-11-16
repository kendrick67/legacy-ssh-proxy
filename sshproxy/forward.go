package sshproxy

import (
	"io"

	"golang.org/x/crypto/ssh"
)

type requestSender interface {
	SendRequest(name string, wantReply bool, payload []byte) (bool, []byte, error)
}

func (p *Proxy) forwardGlobalRequests(dst requestSender, requests <-chan *ssh.Request) {
	for req := range requests {
		ok, payload, err := dst.SendRequest(req.Type, req.WantReply, req.Payload)
		if req.WantReply {
			_ = req.Reply(err == nil && ok, payload)
		}
	}
}

func (p *Proxy) forwardChannel(in ssh.NewChannel, backend *ssh.Client) {
	out, requests, err := backend.OpenChannel(in.ChannelType(), in.ExtraData())
	if err != nil {
		_ = in.Reject(ssh.ConnectionFailed, "unable to open channel on target")
		return
	}
	local, localRequests, err := in.Accept()
	if err != nil {
		_ = out.Close()
		return
	}

	go relayRequests(localRequests, out)
	go relayRequests(requests, local)
	go func() {
		_, _ = io.Copy(out, local)
		_ = out.CloseWrite()
	}()
	_, _ = io.Copy(local, out)
	_ = local.CloseWrite()
	_ = local.Close()
	_ = out.Close()
}

func relayRequests(src <-chan *ssh.Request, dst ssh.Channel) {
	for req := range src {
		ok, err := dst.SendRequest(
			req.Type,
			req.WantReply,
			req.Payload,
		)
		if req.WantReply {
			_ = req.Reply(err == nil && ok, nil)
		}
	}
}
