package codec

import (
	"testing"
)

type ColorGroup struct {
	Id     int      `json:"id" xml:"id,attr" msg:"id"`
	Name   string   `json:"name" xml:"name" msg:"name"`
	Colors []string `json:"colors" xml:"colors" msg:"colors"`
}

var group = ColorGroup{
	Id:     1,
	Name:   "Reds",
	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
}

func BenchmarkJSONCodec_Encode(b *testing.B) {
	raw := make([]byte, 0, 1024)
	serializer := JSONCodec{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw, _ = serializer.Encode(group)
	}
	b.ReportMetric(float64(len(raw)), "bytes")
}

func BenchmarkMsgpackCodec_Encode(b *testing.B) {
	raw := make([]byte, 0, 1024)
	serializer := MsgpackCodec{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		raw, _ = serializer.Encode(group)
	}
	b.ReportMetric(float64(len(raw)), "bytes")
}

func BenchmarkJSONCodec_Decode(b *testing.B) {
	serializer := JSONCodec{}
	bytes, _ := serializer.Encode(group)
	result := ColorGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = serializer.Decode(bytes, &result)
	}
}

func BenchmarkMsgpackCodec_Decode(b *testing.B) {
	serializer := MsgpackCodec{}
	bytes, _ := serializer.Encode(group)
	result := ColorGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = serializer.Decode(bytes, &result)
	}
}
