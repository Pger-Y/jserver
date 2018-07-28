package floader

import (
	"sync"
	"testing"
)

func TestLoaderRun(t *testing.T) {
	fname := "test.json"
	fl := NewFloader(fname)
	fl.Run()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
