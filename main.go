package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gliderlabs/ssh"
)

type Tunnel struct {
	writer io.Writer
	donech chan struct{}
}

// TODO: Make async safe (lock/mutex/etc)
var tunnels = map[int]chan Tunnel{}

const host = "http://localhost:1337"

func main() {
	go func() {
		http.HandleFunc("/", handleRequest)
		log.Fatal(http.ListenAndServe(":1337", nil))
	}()

	ssh.Handle(func(s ssh.Session) {
		// TODO: Use some sort of UUID shortcode library
		id := rand.Intn(math.MaxInt)

		// Create a channel that the http handler will write to
		tunnels[id] = make(chan Tunnel)

		fmt.Printf("Tunnel ID: %d\n", id)
		s.Write([]byte(fmt.Sprintf("File accessed at: %s/?id=%d\n", host, id)))
		s.Write([]byte("Waiting for http connection...\n"))

		// Block until the http handler sends a Tunnel back
		tunnel := <-tunnels[id]
		fmt.Printf("Tunnel[%d] received by ssh server\n", id)

		// Write the bytes that were passed into the ssh session to the htttp response writer
		_, err := io.Copy(tunnel.writer, s)
		if err != nil {
			log.Fatal(err)
		}

		// Let the http handler know that we are done
		close(tunnel.donech)

		// Cleanup channels and tunnels
		close(tunnels[id])
		delete(tunnels, id)

		s.Write([]byte("Success!!!\n"))

	})

	log.Fatal(ssh.ListenAndServe(":22222", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idstr)

	tunnel, ok := tunnels[id]
	if !ok {
		w.Write([]byte("Tunnel Not Found"))
		return
	}

	donech := make(chan struct{})
	tunnel <- Tunnel{
		writer: w,
		donech: donech,
	}
	<-donech
}
