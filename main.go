package main

import (
	"fmt"
	"os"

	"github.com/chancesm/sendit-clone/services/httpservice"
	"github.com/chancesm/sendit-clone/services/sshservice"
	"github.com/chancesm/sendit-clone/services/tunnel"
)

var host string

func main() {

	codespace := os.Getenv("CODESPACES")
	if codespace == "true" {
		host = fmt.Sprintf("https://%s-%s.%s", os.Getenv("CODESPACE_NAME"), "1337", os.Getenv("GITHUB_CODESPACES_PORT_FORWARDING_DOMAIN"))
	}
	fmt.Println("HOST: ", host)
	ts := tunnel.NewTunnelService()
	ts.Init()

	sshservice := sshservice.NewSSHService(ts, host)
	httpservice := httpservice.NewHttpService(ts)

	go func() {
		httpservice.Run()
	}()

	sshservice.Run()
}
