package capnp

import (
	"bytes"
	"io"

	capnpCodec "github.com/onflow/flow-go/model/encoding/capnp"
	"github.com/onflow/flow-go/network"
)

type CapnpCodec struct {
	codec capnpCodec.Codec
}

var _ network.Codec = (*CapnpCodec)(nil)

func NewCapnpCodec(isPacked bool) *CapnpCodec {
	c := &CapnpCodec{
		codec: capnpCodec.Codec{
			IsPacked: isPacked,
		},
	}
	return c
}

func (c *CapnpCodec) NewEncoder(w io.Writer) network.Encoder {
	return c.codec.NewEncoder(w)
}

func (c *CapnpCodec) NewDecoder(r io.Reader) network.Decoder {
	return &Decoder{dec: c.codec.NewDecoder(r)}
}

func (c *CapnpCodec) Encode(v interface{}) ([]byte, error) {
	var data bytes.Buffer
	encoder := c.NewEncoder(&data)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}

	dataBytes := data.Bytes()

	return dataBytes, nil
}

func (c *CapnpCodec) Decode(data []byte) (interface{}, error) {
	dataBuf := bytes.NewBuffer(data)

	decoder := c.NewDecoder(dataBuf)
	return decoder.Decode()

	// In case of unmarshaling
	//capnpMarshaler := capnpCodec.NewMarshaler()
	//msg, _, err := capnp.NewMessage(capnp.SingleSegment(nil))
	//if err != nil {
	//	return nil, err
	//}
	//err = capnpMarshaler.Unmarshal(data, msg)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return msg
}
