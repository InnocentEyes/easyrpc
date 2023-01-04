package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{
		New: func() interface{} {
			return &RequestHeader{}
		},
	}
	ResponsePool = sync.Pool{
		New: func() interface{} {
			return &ResponseHeader{}
		},
	}
}
