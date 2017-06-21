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

	noise := make(chan interface{})

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
	noise := make(chan interface{})

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

	noise := make(chan interface{})

	m := NewChanMux(noise)

	for i := 0; i < msize; i++ {
		newChan := make(chan interface{})
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
