package serial

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// MarshalAny converts a proto Message into any.Any.
func MarshalAny(message proto.Message) *any.Any {
	a, _ := ptypes.MarshalAny(message)
	return a
}

// UnmarshalAny converts given any into a proto Message.
func UnmarshalAny(any *any.Any) (proto.Message, error) {
	protoMessage, err := ptypes.Empty(any)
	if err != nil {
		return nil, err
	}
	if err := ptypes.UnmarshalAny(any, protoMessage); err != nil {
		return nil, err
	}
	return protoMessage, nil
}
