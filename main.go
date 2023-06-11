package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gliderlabs/ssh"

	"github.com/chancesm/sendit-clone/services/tunnel"
)

const host = "http://localhost:1337"

func main() {

	ts := tunnel.NewTunnelService()
	ts.Init()

	go func() {
		http.HandleFunc("/", handleRequest)
		log.Fatal(http.ListenAndServe(":1337", nil))
	}()

	ssh.Handle(func(s ssh.Session) {

		tunnelchan, id := ts.MakeTunnelChannel()

		fmt.Printf("Tunnel ID: %d\n", id)
		s.Write([]byte(fmt.Sprintf("File accessed at: %s/?id=%d\n", host, id)))
		s.Write([]byte("Waiting for http connection...\n"))

		// Block until the http handler sends a Tunnel back
		tnl := <-tunnelchan
		fmt.Printf("Tunnel[%d] received by ssh server\n", id)

		// Write the bytes that were passed into the ssh session to the htttp response writer
		_, err := io.Copy(tnl.Writer, s)
		if err != nil {
			log.Fatal(err)
		}

		ts.Cleanup(&tnl, id)

		s.Write([]byte("Success!!!\n"))

	})

	log.Fatal(ssh.ListenAndServe(":22222", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ts := tunnel.NewTunnelService()

	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	tnlchan := ts.GetTunnel(id)
	// if !ok {
	// 	w.Write([]byte("Tunnel Not Found"))
	// 	return
	// }

	donech := make(chan struct{})
	tnlchan <- tunnel.Tunnel{
		Writer:   w,
		DoneChan: donech,
	}
	<-donech
}
