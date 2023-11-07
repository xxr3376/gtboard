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

func extractScalarValue(v *pb.Summary_Value) (float32, bool) {
	simple_value, ok := v.Value.(*pb.Summary_Value_SimpleValue)
	if ok {
		return simple_value.SimpleValue, true
	}

	var value float32
	if v.Metadata == nil || v.Metadata.PluginData == nil {
		return 0, false
	}
	if v.Metadata.PluginData.PluginName != "scalars" {
		return 0, false
	}
	t, ok := v.Value.(*pb.Summary_Value_Tensor)
	if !ok {
		return 0, false
	}
	switch t.Tensor.Dtype {
	case pb.DataType_DT_HALF:
		// TODO: support half-precision floats
		return 0, false
	case pb.DataType_DT_FLOAT:
		if len(t.Tensor.FloatVal) != 1 {
			return 0, false
		}
		value = t.Tensor.FloatVal[0]
	case pb.DataType_DT_DOUBLE:
		if len(t.Tensor.DoubleVal) != 1 {
			return 0, false
		}
		value = float32(t.Tensor.DoubleVal[0])
	case pb.DataType_DT_INT32, pb.DataType_DT_UINT8, pb.DataType_DT_UINT16, pb.DataType_DT_INT8, pb.DataType_DT_INT16:
		if len(t.Tensor.IntVal) != 1 {
			return 0, false
		}
		value = float32(t.Tensor.IntVal[0])
	case pb.DataType_DT_INT64:
		if len(t.Tensor.Int64Val) != 1 {
			return 0, false
		}
		value = float32(t.Tensor.Int64Val[0])
	case pb.DataType_DT_BOOL:
		if len(t.Tensor.BoolVal) != 1 {
			return 0, false
		}
		if t.Tensor.BoolVal[0] {
			value = 1
		} else {
			value = 0
		}
	case pb.DataType_DT_UINT32:
		if len(t.Tensor.Uint32Val) != 1 {
			return 0, false
		}
		value = float32(t.Tensor.Uint32Val[0])
	}
	return value, true
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
		value, ok := extractScalarValue(v)
		if !ok {
			continue
		}
		scalar_callback(v.Tag, float64(p.buffer.WallTime), p.buffer.Step, value)
	}
	return nil
}
