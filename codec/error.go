package codec

import "errors"

var (
	InvalidSequenceError        = errors.New("invalid sequence number in response")
	UnexpectedChecksumError     = errors.New("unexpected checksum")
	NotFoundCompressorError     = errors.New("not found compressor")
	CompressorTypeMismatchError = errors.New("request and response Compressor type mismatch")
)
