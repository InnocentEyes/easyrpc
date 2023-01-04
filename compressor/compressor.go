package compressor

type CompressType uint16

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

var Compressors = map[CompressType]Compressor{
	Raw:    RowCompressor{},
	Gzip:   GzipCompressor{},
	Snappy: SnappyCompressor{},
	Zlib:   ZlibCompressor{},
}

type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
