package share

import (
	"testing"

	"github.com/smallnest/rpcxlite/protocol"
)

type MockCodec struct{}

func (codec MockCodec) Encode(i interface{}) ([]byte, error) {
	return nil, nil
}

func (codec MockCodec) Decode(data []byte, i interface{}) error {
	return nil
}

func TestShare(t *testing.T) {
	registeredCodecNum := len(Codecs)
	codec := MockCodec{}

	mockCodecType := 127
	RegisterCodec(protocol.SerializeType(mockCodecType), codec)
	if registeredCodecNum+1 != len(Codecs) {
		t.Fatalf("expect %d but got %d", registeredCodecNum+1, len(Codecs))
	}
}
