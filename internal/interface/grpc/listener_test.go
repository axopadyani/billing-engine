package grpc

import "testing"

func TestInitListener(t *testing.T) {
	listener, err := InitListener()
	if err != nil {
		t.Fatal(err)
	}
	_ = listener.Close()
}
