package avro

import (
	"github.com/linkedin/goavro/v2"
	"github.com/onflow/flow-go/model/encoding"
	"io"
)

var _ encoding.Marshaler = (*Marshaler)(nil)

type Marshaler struct{}

func (m Marshaler) Marshal(i interface{}) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m Marshaler) Unmarshal(bytes []byte, i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewMarshaler() *Marshaler {
	return &Marshaler{}
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
	schema string
}

func NewCodec(schema string) *Codec {
	return &Codec{schema: schema}
}

func (c *Codec) NewEncoder(w io.Writer) encoding.Encoder {
	codec, err := goavro.NewCodec(c.schema)
	if err != nil {
		panic(err)
	}

	return &Encoder{
		codec: codec,
		w:     w,
	}
}

func (c *Codec) NewDecoder(r io.Reader) encoding.Decoder {
	codec, err := goavro.NewCodec(c.schema)
	if err != nil {
		panic(err)
	}

	return &Decoder{
		codec: codec,
		r:     r,
	}
}

type Encoder struct {
	codec *goavro.Codec
	w     io.Writer
}

func (e *Encoder) Encode(v interface{}) error {
	var data []byte
	data, err := e.codec.BinaryFromNative(data, v)
	if err != nil {
		return err
	}

	_, err = e.w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

type Decoder struct {
	codec *goavro.Codec
	r     io.Reader
}

func (e *Decoder) Decode(v interface{}) error {
	data := make([]byte, 512)
	totalRead := 0

	for totalRead < len(data) {
		n, err := e.r.Read(data[totalRead:])
		if err != nil {
			return err
		}

		totalRead += n

		// If n is less than len(data) - totalRead, there is more data to read
		if totalRead >= len(data) {
			// Grow the data buffer and continue reading
			newData := make([]byte, len(data)*2) // Double the buffer size
			copy(newData, data)
			data = newData
		} else {
			// Trim data to the actual length of the content
			data = data[:totalRead]
		}
	}

	// Decode the data using the Avro codec
	v, _, err := e.codec.NativeFromBinary(data)
	if err != nil {
		return err
	}

	return nil
}
