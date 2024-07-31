package easycache

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

type Response struct {
	Response []byte
	Header   http.Header
}

func (resp Response) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(resp)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromBytes(data []byte) (Response, error) {
	var resp Response
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	err := dec.Decode(&resp)
	return resp, err
}
