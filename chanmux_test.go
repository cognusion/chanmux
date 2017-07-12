package chanmux

import (
	"fmt"
	"sync"
	"testing"
)

func TestChanMux(t *testing.T) {

	for i := 1; i < 1000; i++ {
		c := testChanMux(i)

		if c != i {
			t.Errorf("Size was %d but had %d outputs\n", i, c)
		}
	}

}

func TestChan(t *testing.T) {

	for i := 1; i < 1000; i++ {
		c := testChan(i)

		if c != i {
			t.Errorf("Size was %d but had %d outputs\n", i, c)
		}
	}

}

func TestChanWG(t *testing.T) {

	for i := 1; i < 1000; i++ {
		c := testChanWG(i)

		if c != i {
			t.Errorf("Size was %d but had %d outputs\n", i, c)
		}
	}

}

// Creates 200 channels, adds them to the ChanMux,
// sends them unique strings, and ranges over the
// aggregate
func ExampleChanMux() {
	msize := 200
	noise := make(chan interface{})

	m := NewChanMux(noise)

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

func BenchmarkChanMux1k(b *testing.B) {

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChanMux(1000)
	}
}

func BenchmarkChan1k(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChan(1000)
	}
}

func BenchmarkChanWG1k(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChanWG(1000)
	}
}

func BenchmarkChanMux1m(b *testing.B) {

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChanMux(1000000)
	}
}

func BenchmarkChan1m(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChan(1000000)
	}
}

func BenchmarkChanWG1m(b *testing.B) {
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testChanWG(1000000)
	}
}

func testChan(msize int) int {
	if msize < 1 {
		// facepalm
		return msize
	}

	noise := make(chan interface{}, msize)

	for i := 0; i < msize; i++ {
		go func(x int, c chan interface{}) {
			c <- fmt.Sprintf("OMG WHAT?! %d?!", x)
		}(i, noise)
	}

	c := 0
	for _ = range noise {
		c++
		if c >= msize {
			close(noise)
		}
	}
	return c
}

func testChanWG(msize int) int {
	if msize < 1 {
		// facepalm
		return msize
	}

	var wg sync.WaitGroup
	noise := make(chan interface{}, msize)

	wg.Add(msize)

	go func(c chan interface{}) {
		wg.Wait()
		close(c)
	}(noise)

	for i := 0; i < msize; i++ {
		go func(x int, c chan interface{}) {
			defer wg.Done()
			c <- fmt.Sprintf("OMG WHAT?! %d?!", x)
		}(i, noise)
	}

	c := 0
	for _ = range noise {
		c++
	}

	return c
}

func testChanMux(msize int) int {
	if msize < 1 {
		// facepalm
		return msize
	}

	noise := make(chan interface{}, msize)

	m := NewChanMux(noise)

	for i := 0; i < msize; i++ {
		newChan := make(chan interface{}, 1)
		m.AddChan(newChan)
		go func(x int, c chan interface{}) {
			c <- fmt.Sprintf("OMG WHAT?! %d?!", x)
			close(c)
		}(i, newChan)
	}
	m.Finalize()

	c := 0
	for _ = range noise {
		c++
	}
	return c
}
