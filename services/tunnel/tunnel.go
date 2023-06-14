package tunnel

import (
	"io"
	"math"
	"math/rand"
)

type Tunnel struct {
	Writer   io.Writer
	DoneChan chan struct{}
}

// TODO: Make async safe (lock/mutex/etc)
var tunnels map[int]chan Tunnel

type TunnelService struct{}

func NewTunnelService() *TunnelService {
	return &TunnelService{}
}

func (t *TunnelService) Init() {
	tunnels = make(map[int]chan Tunnel)
}

func (t *TunnelService) GetTunnelChannel(i int) (chan Tunnel, bool) {
	tunnelchan, ok := tunnels[i]
	return tunnelchan, ok
}

func (t *TunnelService) MakeTunnelChannel() (chan Tunnel, int) {
	// TODO: Use some sort of UUID shortcode library
	id := rand.Intn(math.MaxInt)

	// Create a channel that the http handler will write to
	tunnels[id] = make(chan Tunnel)

	return tunnels[id], id
}

func (t *TunnelService) Cleanup(tnl *Tunnel, id int) {

	// Let the http handler know that we are done
	close(tnl.DoneChan)

	// Cleanup channels and tunnels
	close(tunnels[id])
	delete(tunnels, id)

}
