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
		m.AddChan(newChan)
		newChan <- fmt.Sprintf("OMG WHAT?! %d?!", i)
		close(newChan)
	}

	// If you don't care about ever closing noise,
	// or are going to close it elsewhere, safely,
	// you don't strictly need to Finalize
	m.Finalize()

	// Range over noise and print it
	for n := range noise {
		fmt.Printf("%s\n", n)
	}
}
```

## Benchmarks

There are some other ways of doing this, that I've since found, using reflection (See ReflectSelect), or more 
goros (See GoSelect). I added some benchmarks for those (sourced elsewhere, see source for attribution) and ChanMux 
is still much much faster (ChanMux1k is roughly the scale that the others are tested at). Of interest, is that if 
allocs are an issue (heap problems) the GoSelec issue manages a super small number of allocation, which may be worth 
the performance hit.


```
BenchmarkChanMux1k-8       	    2000	    739115 ns/op	  141138 B/op	    4070 allocs/op
BenchmarkChan1k-8          	    3000	    518314 ns/op	   40637 B/op	    3007 allocs/op
BenchmarkChanWG1k-8        	    2000	    558017 ns/op	   41282 B/op	    3017 allocs/op
BenchmarkChanMux1m-8       	       1	2810778448 ns/op	801630600 B/op	 7000792 allocs/op
BenchmarkChan1m-8          	       2	 751543633 ns/op	95962060 B/op	 3499705 allocs/op
BenchmarkChanWG1m-8        	       2	 789297487 ns/op	127953476 B/op	 3999592 allocs/op
BenchmarkReflectSelect-8   	       1	1908793458 ns/op	922682896 B/op	10107239 allocs/op
BenchmarkGoSelect-8        	      20	  68404865 ns/op	   11156 B/op	     108 allocs/op
```
