package tunnel

import (
	"io"

	shortid "github.com/teris-io/shortid"
)

type Tunnel struct {
	Writer   io.Writer
	DoneChan chan struct{}
}

// TODO: Make async safe (lock/mutex/etc)
var tunnels map[string]chan Tunnel

type TunnelService struct {
	sid *shortid.Shortid
}

func NewTunnelService() *TunnelService {
	sid, _ := shortid.New(1, shortid.DefaultABC, 2468)
	return &TunnelService{
		sid: sid,
	}
}

func (t *TunnelService) Init() {
	tunnels = make(map[string]chan Tunnel)
}

func (t *TunnelService) GetTunnelChannel(id string) (chan Tunnel, bool) {
	tunnelchan, ok := tunnels[id]
	return tunnelchan, ok
}

func (t *TunnelService) MakeTunnelChannel() (chan Tunnel, string) {
	// TODO: Use some sort of UUID shortcode library

	id, _ := t.sid.Generate()

	// Create a channel that the http handler will write to
	tunnels[id] = make(chan Tunnel)

	return tunnels[id], id
}

func (t *TunnelService) Cleanup(tnl *Tunnel, id string) {

	// Let the http handler know that we are done
	close(tnl.DoneChan)

	// Cleanup channels and tunnels
	close(tunnels[id])
	delete(tunnels, id)

}
