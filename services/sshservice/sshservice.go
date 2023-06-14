package sshservice

import (
	"fmt"
	"io"
	"log"

	"github.com/chancesm/sendit-clone/services/tunnel"
	"github.com/gliderlabs/ssh"
)

type SSHService struct {
	ts *tunnel.TunnelService
	s  *ssh.Server
}

func NewSSHService(t *tunnel.TunnelService) *SSHService {
	s := &ssh.Server{
		Addr: ":22222",
	}

	ss := &SSHService{
		ts: t,
	}

	// Register ssh handler and add server to service
	s.Handle(ss.sshHandler)
	ss.s = s

	return ss
}

func (ss *SSHService) Run() {
	log.Fatal(ss.s.ListenAndServe())
}

func (ss *SSHService) sshHandler(s ssh.Session) {
	tunnelchan, id := ss.ts.MakeTunnelChannel()

	fmt.Printf("Tunnel ID: %d\n", id)
	s.Write([]byte(fmt.Sprintf("File accessed at: %s/%d\n", "example.com", id)))
	s.Write([]byte("Waiting for http connection...\n"))

	// Block until the http handler sends a Tunnel back
	tnl := <-tunnelchan
	fmt.Printf("Tunnel[%d] received by ssh server\n", id)

	// Write the bytes that were passed into the ssh session to the htttp response writer
	_, err := io.Copy(tnl.Writer, s)
	if err != nil {
		log.Fatal(err)
	}

	ss.ts.Cleanup(&tnl, id)

	s.Write([]byte("Success!!!\n"))

}
