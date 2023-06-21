package ingest

import (
	"github.com/golang/protobuf/proto"
	pb "github.com/xxr3376/gtboard/tensorboard_pb"
)

type ScalarEventCallback func(string, float64, int64, float32)

// Parser is an interface for parsing TFRecord files
type parser struct {
	buffer pb.Event
}

func NewParser() *parser {
	return &parser{}
}

// ParseRecord parses a single TFRecord event
// ParseRecord is not thread-safe due to reusing the same buffer
func (p *parser) ParseRecord(data []byte, scalar_callback ScalarEventCallback) error {
	p.buffer.Reset()
	err := proto.Unmarshal(data, &p.buffer)
	if err != nil {
		return err
	}
	if p.buffer.What == nil {
		return nil
	}
	summary := p.buffer.GetSummary()
	if summary == nil {
		return nil
	}
	for _, v := range summary.Value {
		simple_value, ok := v.Value.(*pb.Summary_Value_SimpleValue)
		if !ok {
			continue
		}
		scalar_callback(v.Tag, float64(p.buffer.WallTime), p.buffer.Step, simple_value.SimpleValue)
	}
	return nil
}
