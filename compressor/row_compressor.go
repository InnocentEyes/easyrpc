package compressor

type RowCompressor struct {
}

func (_ RowCompressor) Zip(data []byte) ([]byte, error) {
	return data, nil
}

func (_ RowCompressor) Unzip(data []byte) ([]byte, error) {
	return data, nil
}
