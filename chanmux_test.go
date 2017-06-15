package chanmux

import (
	"fmt"
	"testing"
)

func TestChanMux(t *testing.T) {
	msize := 200
	noise := make(chan interface{})

	m := NewChanMux(noise)

	for i := 0; i < msize; i++ {
		newChan := make(chan interface{})
		go m.AddChan(newChan)
		newChan <- fmt.Sprintf("OMG WHAT?! %d?!", i)
		close(newChan)
	}
	m.Finalize(msize)

	c := 0
	for _ = range noise {
		c++
	}

	if c != msize {
		t.Errorf("Size was %d but had %d outputs\n", msize, c)
	}

}
