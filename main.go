package main

import (
	"github.com/chancesm/sendit-clone/services/httpservice"
	"github.com/chancesm/sendit-clone/services/sshservice"
	"github.com/chancesm/sendit-clone/services/tunnel"
)

func main() {

	ts := tunnel.NewTunnelService()
	ts.Init()

	sshservice := sshservice.NewSSHService(ts)
	httpservice := httpservice.NewHttpService(ts)

	go func() {
		httpservice.Run()
	}()

	sshservice.Run()
}
