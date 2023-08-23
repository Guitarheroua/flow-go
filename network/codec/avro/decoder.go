package avro

import (
	"github.com/onflow/flow-go/model/encoding"
	"github.com/onflow/flow-go/network/codec"
)

type Decoder struct {
	dec encoding.Decoder
}

func (d *Decoder) Decode() (interface{}, error) {
	data := make(map[string]interface{})
	err := d.dec.Decode(data)
	if err != nil {
		return nil, codec.NewInvalidEncodingErr(err)
	}

	return data, nil
}
