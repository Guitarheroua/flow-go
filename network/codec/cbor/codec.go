// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package cbor

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"time"

	"github.com/fxamacker/cbor/v2"

	cborcodec "github.com/onflow/flow-go/model/encoding/cbor"
	"github.com/onflow/flow-go/network"
	"github.com/onflow/flow-go/network/codec"
	_ "github.com/onflow/flow-go/utils/binstat"
)

var defaultDecMode, _ = cbor.DecOptions{ExtraReturnErrors: cbor.ExtraDecErrorUnknownField}.DecMode()

// Codec represents a CBOR codec for our network.
type Codec struct {
	l *zerolog.Logger
}

// NewCodec creates a new CBOR codec.
func NewCodec() *Codec {
	c := &Codec{}
	return c
}

func (c *Codec) WithLogger(l *zerolog.Logger) *Codec {
	c.l = l
	return c
}

// NewEncoder creates a new CBOR encoder with the given underlying writer.
func (c *Codec) NewEncoder(w io.Writer) network.Encoder {
	enc := cborcodec.EncMode.NewEncoder(w)
	return &Encoder{enc: enc}
}

// NewDecoder creates a new CBOR decoder with the given underlying reader.
func (c *Codec) NewDecoder(r io.Reader) network.Decoder {
	dec := defaultDecMode.NewDecoder(r)
	return &Decoder{dec: dec}
}

// Encode will, given a Golang interface 'v', return a []byte 'envelope'.
// Return an error if packing the envelope fails.
// NOTE: 'v' is the network message payload in unserialized form.
// NOTE: 'code' is the message type.
// NOTE: 'what' is the 'code' name for debugging / instrumentation.
// NOTE: 'envelope' contains 'code' & serialized / encoded 'v'.
// i.e.  1st byte is 'code' and remaining bytes are CBOR encoded 'v'.
func (c *Codec) Encode(v interface{}) ([]byte, error) {
	var start time.Time
	if c.l != nil {
		start = time.Now()
	}
	// encode the value
	code, what, err := codec.MessageCodeFromInterface(v)
	if err != nil {
		return nil, fmt.Errorf("could not determine envelope code: %w", err)
	}

	// NOTE: benchmarking shows that prepending the code and then using
	//       .NewEncoder() to .Encode() is the fastest.

	// encode / append the envelope code
	//bs1 := binstat.EnterTime(binstat.BinNet + ":wire<1(cbor)envelope2payload")
	var data bytes.Buffer
	data.WriteByte(code.Uint8())
	//binstat.LeaveVal(bs1, int64(data.Len()))

	// encode the payload
	//bs2 := binstat.EnterTime(fmt.Sprintf("%s%s%s:%d", binstat.BinNet, ":wire<2(cbor)", what, code)) // e.g. ~3net::wire<1(cbor)CodeEntityRequest:23
	encoder := cborcodec.EncMode.NewEncoder(&data)
	err = encoder.Encode(v)
	//binstat.LeaveVal(bs2, int64(data.Len()))
	if err != nil {
		return nil, fmt.Errorf("could not encode CBOR payload with envelope code %d AKA %s: %w", code, what, err) // e.g. 2, "CodeBlockProposal", <CBOR error>
	}

	dataBytes := data.Bytes()

	var duration time.Duration
	if c.l != nil {
		duration = time.Since(start)
	}

	byteUUID := generateOneByteUUID()
	dataBytes = append(dataBytes, byteUUID)

	if c.l != nil {
		c.l.Debug().Msg(fmt.Sprintf("CBOR: Execution of Encode took: %s for %s with size: %d bytes.", duration, fmt.Sprintf("1-byte UUID: 0x%02X\n", byteUUID), data.Len()))
	}

	return dataBytes, nil
}

// Decode will, given a []byte 'envelope', return a Golang interface 'v'.
// Return an error if unpacking the envelope fails.
// NOTE: 'v' is the network message payload in un-serialized form.
// NOTE: 'code' is the message type.
// NOTE: 'what' is the 'code' name for debugging / instrumentation.
// NOTE: 'envelope' contains 'code' & serialized / encoded 'v'.
// i.e.  1st byte is 'code' and remaining bytes are CBOR encoded 'v'.
// Expected error returns during normal operations:
//   - codec.ErrInvalidEncoding if message encoding is invalid.
//   - codec.ErrUnknownMsgCode if message code byte does not match any of the configured message codes.
//   - codec.ErrMsgUnmarshal if the codec fails to unmarshal the data to the message type denoted by the message code.
func (c *Codec) Decode(data []byte) (interface{}, error) {
	var start time.Time
	if c.l != nil {
		start = time.Now()
	}

	msgCode, err := codec.MessageCodeFromPayload(data)
	if err != nil {
		return nil, err
	}
	// decode the envelope
	//bs1 := binstat.EnterTime(binstat.BinNet + ":wire>3(cbor)payload2envelope")

	//binstat.LeaveVal(bs1, int64(len(data)))

	msgInterface, what, err := codec.InterfaceFromMessageCode(msgCode)
	if err != nil {
		return nil, err
	}

	// unmarshal the payload
	//bs2 := binstat.EnterTimeVal(fmt.Sprintf("%s%s%s:%d", binstat.BinNet, ":wire>4(cbor)", what, code), int64(len(data))) // e.g. ~3net:wire>4(cbor)CodeEntityRequest:23
	err = defaultDecMode.Unmarshal(data[1:len(data)-1], msgInterface) // all but first byte
	//binstat.Leave(bs2)
	if err != nil {
		return nil, codec.NewMsgUnmarshalErr(data[0], what, err)
	}

	var duration time.Duration
	if c.l != nil {
		duration = time.Since(start)
	}

	byteUUID := data[len(data)-1]
	data = data[:len(data)-1]

	if c.l != nil {
		c.l.Debug().Msg(fmt.Sprintf("CBOR: Execution of Decode took: %s for %s with size: %d bytes.", duration, fmt.Sprintf("1-byte UUID: 0x%02X\n", byteUUID), len(data)))
	}

	return msgInterface, nil
}

func generateOneByteUUID() byte {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random byte value
	return byte(r.Intn(256))
}
