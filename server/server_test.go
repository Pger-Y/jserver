package server

import (
	"testing"
)

func TestServe(t *testing.T) {
	d := make(map[string]interface{})
	d["/sss"] = []int{1, 2, 3}
	s := NewServer(":8080", d)
	s.Restart()
	select {}
}
