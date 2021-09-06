package client

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/smallnest/rpcxlite/server"
)

func TestXClient_IT(t *testing.T) {
	s := server.NewServer()
	s.RegisterName("Arith", new(Arith), "")
	go s.Serve("tcp", "127.0.0.1:0")
	defer s.Close()
	time.Sleep(500 * time.Millisecond)

	addr := s.Address().String()

	d, err := NewPeer2PeerDiscovery("tcp@"+addr, "desc=a test service")
	if err != nil {
		t.Fatalf("failed to NewPeer2PeerDiscovery: %v", err)
	}

	xclient := NewXClient("Arith", Failtry, RandomSelect, d, DefaultOption)

	defer xclient.Close()

	args := &Args{
		A: 10,
		B: 20,
	}

	reply := &Reply{}
	err = xclient.Call(context.Background(), "Mul", args, reply)
	if err != nil {
		t.Fatalf("failed to call: %v", err)
	}

	if reply.C != 200 {
		t.Fatalf("expect 200 but got %d", reply.C)
	}
}

func TestXClient_filterByStateAndGroup(t *testing.T) {
	servers := map[string]string{"a": "", "b": "state=inactive&ops=10", "c": "ops=20", "d": "group=test&ops=20"}
	filterByStateAndGroup("test", servers)
	if _, ok := servers["b"]; ok {
		t.Error("has not remove inactive node")
	}
	if _, ok := servers["a"]; ok {
		t.Error("has not remove inactive node")
	}
	if _, ok := servers["c"]; ok {
		t.Error("has not remove inactive node")
	}
	if _, ok := servers["d"]; !ok {
		t.Error("node must be removed")
	}
}

func TestUncoverError(t *testing.T) {
	var e error = ServiceError("error")
	if uncoverError(e) {
		t.Fatalf("expect false but get true")
	}

	if uncoverError(context.DeadlineExceeded) {
		t.Fatalf("expect false but get true")
	}

	if uncoverError(context.Canceled) {
		t.Fatalf("expect false but get true")
	}

	e = errors.New("error")
	if !uncoverError(e) {
		t.Fatalf("expect true but get false")
	}
}
