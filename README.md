# chanmux
Multiplex incoming channels into an outgoing channel. 
You're probably doing something wrong if you need this. 
But if not, here it is.

The reason you might not be doing it wrong, is that you've decided that giving an unknown number of goros their own 
individual channels allows them to bettern communicate when they're done (by closing their channel) than by passing multiple
channels, or having some kind of channel protocol, or assigning names to your goros and keeping score. Benchmarks (included) show 
it's hella fast and, again, doesn't require advanced knowledge of goro counts.


```go
package main

import (
	"fmt"
	"github.com/cognusion/chanmux"
)

func main() {
	msize := 200
	noise := make(chan interface{})

	m := chanmux.NewChanMux(noise)

	// Create msize number of channels, adding them
	// to the mux, sending them a string, and closing
	// them thereafter
	for i := 0; i < msize; i++ {
		newChan := make(chan interface{})
		go m.AddChan(newChan)
		newChan <- fmt.Sprintf("OMG WHAT?! %d?!", i)
		close(newChan)
	}

	// If you don't care about ever closing noise,
	// or are going to close it elsewhere, safely,
	// you don't strictly need to Finalize
	m.Finalize(msize)

	// Range over noise and print it
	for n := range noise {
		fmt.Printf("%s\n", n)
	}
}
```

