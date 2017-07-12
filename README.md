# chanmux
Multiplex incoming channels into an outgoing channel. 
You're probably doing something wrong if you need this. 
But if not, here it is.

The reason you _might_ not be doing it wrong, is that you've decided that giving an unknown number of goros their own 
individual channels allows them to bettern communicate when they're done (by closing their channel) than by passing multiple
channels, or having some kind of channel protocol, or assigning names to your goros and keeping score. Benchmarks (included) show 
it's hella fast and, again, doesn't require advanced knowledge of goro counts.

[![GoDoc](https://godoc.org/github.com/cognusion/chanmux?status.svg)](https://godoc.org/github.com/cognusion/chanmux)

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

Benchmarking 1thousand and 1million muxed channels (prefix _*_) compared to the static methods of raw channel iteration (prefix _-_)
and raw channel iteration using Waitgroups (prefix _^_) shows that the mux outscales raw channel iteration, and isn't too shabby vs.
the similar using WaitGroups.

There are some other ways of doing this, that I've since found (prefix _+_ below), using reflection (See ReflectSelect), or more 
goros (See GoSelect). I added some benchmarks for those (sourced elsewhere, see source for attribution) and ChanMux 
is still much much faster (ChanMux1k is roughly the scale that the others are tested at), and more dynamic regardless. 

Of interest, is that if allocs are an issue (heap problems) the GoSelect scheme manages a super small number of 
allocations, which may be worth the performance hit and initial sizing requirements.


```
*BenchmarkChanMux1k-8       	    2000	    657259 ns/op	  169019 B/op	    5006 allocs/op
-BenchmarkChan1k-8         	    5000	    302023 ns/op	   56760 B/op	    3004 allocs/op
^BenchmarkChanWG1k-8        	    5000	    315627 ns/op	   57215 B/op	    3012 allocs/op
*BenchmarkChanMux1m-8       	       2	 694727670 ns/op	184342996 B/op	 5002773 allocs/op
-BenchmarkChan1m-8          	       1	2051025698 ns/op	578326880 B/op	 4014742 allocs/op
^BenchmarkChanWG1m-8        	       5	 375152299 ns/op	80566953 B/op	 3121603 allocs/op
+BenchmarkReflectSelect-8   	       1	1725111094 ns/op	920304784 B/op	10068110 allocs/op
+BenchmarkGoSelect-8        	      20	  62888297 ns/op	   10964 B/op	     106 allocs/op

```
