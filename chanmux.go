package chanmux

import (
	"sync"
)

// ChanMux allows muxing multiple incoming channels into an
// outgoing channel (the "sink")
type ChanMux struct {
	sink chan interface{}
	wg   sync.WaitGroup
}

// AddChan accepts a channel to read from and write to the sink
func (c *ChanMux) AddChan(newChan chan interface{}) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()
		for p := range newChan {
			c.sink <- p
		}
	}()

}

// Finalize will close the sink when the muxes have all exited
func (c *ChanMux) Finalize() {
	go func() {
		// Block until everyone is done, and then close the sink
		c.wg.Wait()
		close(c.sink)
	}()
}

// NewChanMux accepts a channel to write to, and returns a *ChanMux
func NewChanMux(sink chan interface{}) *ChanMux {
	p := ChanMux{
		sink: sink,
	}
	return &p
}
