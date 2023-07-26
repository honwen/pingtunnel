package common

import (
	"bytes"
	"compress/zlib"
	"io"
)

func CompressData(src []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(src)
	w.Close()
	return b.Bytes()
}

func DeCompressData(src []byte) ([]byte, error) {
	b := bytes.NewReader(src)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	io.Copy(&out, r)
	r.Close()
	return out.Bytes(), nil
}
