package client

import (
	"context"
	"testing"

	"github.com/smallnest/rpcxlite/protocol"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func TestClient_Res_Reset(t *testing.T) {
	res := protocol.NewMessage()
	res.Payload = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	data := res.Payload
	res.Reset()

	if len(data) == 0 {
		t.Fatalf("data has been set to empty after response has been reset: %v", data)
	}
}
