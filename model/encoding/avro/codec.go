package capnp

import (
	"capnproto.org/go/capnp/v3"
	"errors"
	"github.com/onflow/flow-go/model/encoding"
	"io"
)

var _ encoding.Marshaler = (*Marshaler)(nil)

type Marshaler struct{}

func NewMarshaler() *Marshaler {
	return &Marshaler{}
}

func (m Marshaler) Marshal(i interface{}) ([]byte, error) {
	msg, err := convertToCapnpMessage(i)
	if err != nil {
		return nil, err
	}

	bytes, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (m Marshaler) Unmarshal(bytes []byte, i interface{}) error {
	if _, ok := i.(*capnp.Message); !ok {
		return errors.New("conversion to *capnp.Message failed")
	}

	i, err := capnp.Unmarshal(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (m Marshaler) MustMarshal(i interface{}) []byte {
	b, err := m.Marshal(i)
	if err != nil {
		panic(err)
	}

	return b
}

func (m Marshaler) MustUnmarshal(bytes []byte, i interface{}) {
	err := m.Unmarshal(bytes, i)
	if err != nil {
		panic(err)
	}
}

var _ encoding.Codec = (*Codec)(nil)

type Codec struct {
	IsPacked bool
}

func (c *Codec) NewEncoder(w io.Writer) encoding.Encoder {
	if c.IsPacked {
		return &Encoder{capnp.NewPackedEncoder(w)}
	}

	return &Encoder{capnp.NewEncoder(w)}
}

func (c *Codec) NewDecoder(r io.Reader) encoding.Decoder {
	if c.IsPacked {
		return &Decoder{capnp.NewPackedDecoder(r)}
	}

	return &Decoder{capnp.NewDecoder(r)}
}

type Encoder struct {
	encoder *capnp.Encoder
}

func (e *Encoder) Encode(v interface{}) error {
	msg, err := convertToCapnpMessage(v)
	if err != nil {
		return err
	}

	return e.encoder.Encode(msg)
}

type Decoder struct {
	decoder *capnp.Decoder
}

func (e *Decoder) Decode(v interface{}) error {
	_, err := convertToCapnpMessage(v)
	if err != nil {
		return err
	}

	v, err = e.decoder.Decode()
	if err != nil {
		return err
	}

	return nil
}

func convertToCapnpMessage(data interface{}) (*capnp.Message, error) {
	if msg, ok := data.(*capnp.Message); ok {
		return msg, nil
	}

	return nil, errors.New("conversion of interface to *capnp.Message failed")
}
