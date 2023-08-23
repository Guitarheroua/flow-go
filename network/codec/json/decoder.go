package json

import (
	"github.com/onflow/flow-go/model/encoding"
	"github.com/onflow/flow-go/network/codec"
)

type Decoder struct {
	decoder    encoding.Decoder
	marshaller encoding.Marshaler
}

func (d *Decoder) Decode() (interface{}, error) {
	var data []byte
	err := d.decoder.Decode(&data)
	if err != nil {
		return nil, codec.NewInvalidEncodingErr(err)
	}

	msgCode, err := codec.MessageCodeFromPayload(data)
	if err != nil {
		return nil, err
	}

	msgInterface, what, err := codec.InterfaceFromMessageCode(msgCode)
	if err != nil {
		return nil, err
	}

	// unmarshal the payload
	err = d.marshaller.Unmarshal(data[1:], msgInterface)
	if err != nil {
		return nil, codec.NewMsgUnmarshalErr(data[0], what, err)
	}
	return msgInterface, nil
}
