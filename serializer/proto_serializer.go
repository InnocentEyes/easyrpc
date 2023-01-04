package serializer

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

var NotImplementProtoMessageError = errors.New("param does not implement proto.Message")

var Proto = ProtoSerializer{}

type ProtoSerializer struct {
}

func (_ ProtoSerializer) Marshal(message interface{}) ([]byte, error) {
	var body proto.Message
	if message == nil {
		return []byte{}, nil
	}
	var ok bool
	if body, ok = message.(proto.Message); !ok {
		return nil, NotImplementProtoMessageError
	}
	return proto.Marshal(body)
}

func (_ ProtoSerializer) UnMarshal(data []byte, message interface{}) error {
	var body proto.Message
	if message == nil {
		return nil
	}
	var ok bool
	body, ok = message.(proto.Message)
	if !ok {
		return NotImplementProtoMessageError
	}
	return proto.Unmarshal(data, body)
}
