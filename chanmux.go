package chanmux

import (
	"github.com/cognusion/semaphore"
)

// ChanMux allows muxing multiple incoming channels into an
// outgoing channel (the "sink")
type ChanMux struct {
	sink    chan interface{}
	counter semaphore.Semaphore
}

// AddChan accepts a channel to read from and write to the sink
func (c *ChanMux) AddChan(newChan chan interface{}) {

	for p := range newChan {
		c.sink <- p
	}

	// Free up a lock
	c.counter.Unlock()
}

// Finalize accepts the peak number of channels being muxed,
// and will close the sink when they've all exited
func (c *ChanMux) Finalize(final int) {
	go func() {
		// Block until everyone is done, and then close the sink
		c.counter.Add(final)
		close(c.sink)
	}()
}

// NewChanMux accepts a channel to write to, and returns a *ChanMux
func NewChanMux(sink chan interface{}) *ChanMux {
	p := ChanMux{
		sink:    sink,
		counter: semaphore.NewSemaphore(0),
	}
	return &p
}
