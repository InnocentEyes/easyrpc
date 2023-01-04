package codec

import (
	"encoding/binary"
	"io"
	"net"
)

//sendFrame 发送信息
func sendFrame(w io.Writer, data []byte) (err error) {
	var size [binary.MaxVarintLen64]byte

	if len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n]); err != nil {
			return
		}
	}

	n := binary.PutUvarint(size[:], uint64(len(data)))

	if err = write(w, size[:n]); err != nil {
		return
	}

	if err = write(w, data[:]); err != nil {
		return
	}
	return
}

//recvFrame 接受数据
func recvFrame(r io.Reader) (data []byte, err error) {
	size, err := binary.ReadUvarint(r.(io.ByteReader))
	if err != nil {
		return nil, err
	}
	data = make([]byte, size)
	if err = read(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

//write
func write(w io.Writer, data []byte) error {
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		index += n
	}
	return nil
}

//read
func read(r io.Reader, data []byte) error {
	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		index += n
	}

	return nil
}
