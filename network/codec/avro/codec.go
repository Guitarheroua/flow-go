package avro

import (
	"bytes"
	"io"

	"github.com/onflow/flow-go/model/encoding/avro"
	"github.com/onflow/flow-go/network"
)

type Codec struct {
	codec *avro.Codec
}

var _ network.Codec = (*Codec)(nil)

func NewCodec(schema string) *Codec {
	return &Codec{codec: avro.NewCodec(schema)}
}

func (c *Codec) NewEncoder(w io.Writer) network.Encoder {
	return c.codec.NewEncoder(w)
}

func (c *Codec) NewDecoder(r io.Reader) network.Decoder {
	return &Decoder{dec: c.codec.NewDecoder(r)}
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	var data bytes.Buffer
	encoder := c.NewEncoder(&data)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	dataBytes := data.Bytes()

	return dataBytes, nil
}

func (c *Codec) Decode(data []byte) (interface{}, error) {
	dataBuf := bytes.NewBuffer(data)

	decoder := c.NewDecoder(dataBuf)
	return decoder.Decode()
}
