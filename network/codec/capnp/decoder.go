// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package capnp

import (
	"capnproto.org/go/capnp/v3"
	"github.com/onflow/flow-go/model/encoding"
	"github.com/onflow/flow-go/network/codec"
)

type Decoder struct {
	dec encoding.Decoder
}

func (d *Decoder) Decode() (interface{}, error) {
	msg, _, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}

	err = d.dec.Decode(msg)
	if err != nil {
		return nil, codec.NewInvalidEncodingErr(err)
	}

	return msg, nil
}
