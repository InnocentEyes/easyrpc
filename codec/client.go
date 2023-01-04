package codec

import (
	"bufio"
	"easyrpc/compressor"
	"easyrpc/header"
	"easyrpc/serializer"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	compressor compressor.CompressType
	serializer serializer.Serializer
	response   header.ResponseHeader
	mutex      sync.Mutex
	pending    map[uint64]string
}

func NewClientCodec(conn io.ReadWriteCloser, compressType compressor.CompressType,
	serializer serializer.Serializer) rpc.ClientCodec {
	return &clientCodec{
		r: bufio.NewReader(conn),
		w: bufio.NewWriter(conn),
		c: conn,

		compressor: compressType,
		serializer: serializer,
		response:   header.ResponseHeader{},
		mutex:      sync.Mutex{},
		pending:    make(map[uint64]string),
	}
}

func (c *clientCodec) WriteRequest(rq *rpc.Request, params interface{}) error {
	c.mutex.Lock()
	c.pending[rq.Seq] = rq.ServiceMethod
	c.mutex.Unlock()

	if _, ok := compressor.Compressors[c.compressor]; !ok {
		return NotFoundCompressorError
	}
	resbody, err := c.serializer.Marshal(params)
	if err != nil {
		return err
	}
	compressedReqBody, err := compressor.Compressors[c.compressor].Zip(resbody)
	if err != nil {
		return nil
	}
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.ResponsePool.Put(h)
	}()
	h.ID = rq.Seq
	h.Method = rq.ServiceMethod
	h.RequestLen = uint32(len(compressedReqBody))
	h.CompressType = c.compressor
	h.Checksum = crc32.ChecksumIEEE(compressedReqBody)
	if err := sendFrame(c.w, h.Marshal()); err != nil {
		return nil
	}
	if err := write(c.w, compressedReqBody); err != nil {
		return nil
	}
	c.w.(*bufio.Writer).Flush()
	return nil
}

func (c *clientCodec) ReadResponseHeader(res *rpc.Response) error {
	c.response.ResetHeader()
	data, err := recvFrame(c.r)
	if err != nil {
		return err
	}

	err = c.response.Unmarshal(data)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	res.Seq = c.response.ID
	res.ServiceMethod = c.pending[c.response.ID]
	res.Error = c.response.Error
	delete(c.pending, res.Seq)
	c.mutex.Unlock()
	return nil
}

func (c *clientCodec) ReadResponseBody(params interface{}) error {
	//如果没有参数
	if params == nil {
		if err := read(c.r, make([]byte, c.response.ResponseLen)); err != nil {
			return err
		}
		return nil
	}

	resbody := make([]byte, c.response.ResponseLen)

	err := read(c.r, resbody)

	if err != nil {
		return err
	}

	if c.response.Checksum != 0 {
		if crc32.ChecksumIEEE(resbody) != c.response.Checksum {
			return UnexpectedChecksumError
		}
	}

	if c.response.GetCompressType() != c.compressor {
		return CompressorTypeMismatchError
	}

	resp, err := compressor.Compressors[c.response.GetCompressType()].Unzip(resbody)

	if err != nil {
		return err
	}

	return c.serializer.UnMarshal(resp, params)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
