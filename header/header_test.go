package header

import (
	"easyrpc/compressor"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRequestHeader_Marshal(t *testing.T) {
	header := &RequestHeader{
		CompressType: 0,
		Method:       "Add",
		ID:           128,
		RequestLen:   512,
		Checksum:     1024520,
	}

	assert.Equal(t, []byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64,
		0x80, 0x1, 0x80, 0x4, 0x8, 0xa2, 0xf, 0x0}, header.Marshal())
}

func TestRequestHeader_Unmarshal(t *testing.T) {
	type expect struct {
		header *RequestHeader
		err    error
	}
	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test1",
			[]byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64,
				0x80, 0x1, 0x80, 0x4, 0x8, 0xa2, 0xf, 0x0},
			expect{
				header: &RequestHeader{
					CompressType: 0,
					Method:       "Add",
					ID:           128,
					RequestLen:   512,
					Checksum:     1024520,
				},
				err: nil,
			},
		},
		{
			"test2",
			nil,
			expect{
				header: &RequestHeader{},
				err:    UnmarshalError,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			header := &RequestHeader{}
			err := header.Unmarshal(c.data)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.header, header))
			assert.Equal(t, err, c.expect.err)
		})
	}
}

func TestRequestHeader_GetCompressType(t *testing.T) {
	header := &RequestHeader{CompressType: 0}
	assert.Equal(t, true, reflect.DeepEqual(compressor.CompressType(0), header.GetCompressType()))
}

func TestRequestHeader_ResetHeader(t *testing.T) {
	header := &RequestHeader{
		CompressType: 0,
		Method:       "Add",
		ID:           128,
		RequestLen:   512,
		Checksum:     1024520,
	}
	header.ResetHeader()
	assert.Equal(t, true, reflect.DeepEqual(header, &RequestHeader{}))
}

func TestResponseHeader_Marshal(t *testing.T) {
	res := &ResponseHeader{
		CompressType: 1,
		ID:           64,
		Error:        "error occured",
		ResponseLen:  128,
		Checksum:     256,
	}
	assert.Equal(t, []byte{0x1, 0x0, 0x40, 0xd, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x6f, 0x63,
		0x63, 0x75, 0x72, 0x65, 0x64, 0x80, 0x1, 0x0, 0x1, 0x0, 0x0}, res.Marshal())
}

func TestResponseHeader_Unmarshal(t *testing.T) {
	type expect struct {
		res *ResponseHeader
		err error
	}

	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test-1",
			[]byte{0x1, 0x0, 0x40, 0xd, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x20, 0x6f, 0x63,
				0x63, 0x75, 0x72, 0x65, 0x64, 0x80, 0x1, 0x0, 0x1, 0x0, 0x0},
			expect{
				res: &ResponseHeader{
					CompressType: 1,
					ID:           64,
					Error:        "error occured",
					ResponseLen:  128,
					Checksum:     256,
				},
				err: nil,
			},
		},
		{
			"test-2",
			nil,
			expect{
				res: &ResponseHeader{},
				err: UnmarshalError,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := &ResponseHeader{}
			err := res.Unmarshal(c.data)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.res, res))
			assert.Equal(t, err, c.expect.err)
		})
	}
}
