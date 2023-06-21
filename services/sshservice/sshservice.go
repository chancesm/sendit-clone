package sshservice

import (
	"fmt"
	"io"
	"log"

	"github.com/chancesm/sendit-clone/services/tunnel"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
)

type SSHService struct {
	ts   *tunnel.TunnelService
	s    *ssh.Server
	host string
}

func NewSSHService(t *tunnel.TunnelService, host string) *SSHService {
	ss := &SSHService{
		ts:   t,
		host: host,
	}

	port := 22222
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf(":%d", port)),
		wish.WithHostKeyPath("ssh_example"),
		wish.WithMiddleware(ss.sshHandler, lm.Middleware()),
	)
	if err != nil {
		log.Fatal(err)
	}

	ss.s = s

	return ss
}

func (ss *SSHService) Run() {
	log.Fatal(ss.s.ListenAndServe())
}

func (ss *SSHService) sshHandler(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {

		// Get client's output.
		clientOutput := outputFromSession(s)

		pty, _, active := s.Pty()
		if !active {
			next(s)
			return
		}
		width := pty.Window.Width
		_ = width

		// Initialize new renderer for the client.
		renderer := lipgloss.NewRenderer(s)
		renderer.SetOutput(clientOutput)

		tunnelchan, id := ss.ts.MakeTunnelChannel()

		fmt.Printf("Tunnel ID: %s\n", id)
		// Create a renderer for the client.

		s.Write([]byte(fmt.Sprintf("File accessed at: %s/f/%s\n", ss.host, id)))
		s.Write([]byte("Waiting for http connection...\n"))

		// Block until the http handler sends a Tunnel back
		tnl := <-tunnelchan
		fmt.Printf("Tunnel[%s] received by ssh server\n", id)

		// Write the bytes that were passed into the ssh session to the htttp response writer
		_, err := io.Copy(tnl.Writer, s)
		if err != nil {
			log.Fatal(err)
		}

		ss.ts.Cleanup(&tnl, id)

		s.Write([]byte("Success!!!\n"))

		next(s)
	}
}
