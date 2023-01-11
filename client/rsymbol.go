package client

import (
	"fmt"

	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/scale"
)

type RSymbol string

const (
	RFIS   = RSymbol("RFIS")
	RDOT   = RSymbol("RDOT")
	RKSM   = RSymbol("RKSM")
	RATOM  = RSymbol("RATOM")
	RSOL   = RSymbol("RSOL")
	RMATIC = RSymbol("RMATIC")
	RBNB   = RSymbol("RBNB")
	RETH   = RSymbol("RETH")
)

func (r *RSymbol) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*r = RFIS
		// return fmt.Errorf("RSymbol decode error: %d", b)
	case 1:
		*r = RDOT
	case 2:
		*r = RKSM
	case 3:
		*r = RATOM
	case 4:
		*r = RSOL
	case 5:
		*r = RMATIC
	case 6:
		*r = RBNB
	case 7:
		*r = RETH
	default:
		return fmt.Errorf("RSymbol decode error: %d", b)
	}

	return nil
}

func (r RSymbol) Encode(encoder scale.Encoder) error {
	switch r {
	case RFIS:
		return encoder.PushByte(0)
		// return fmt.Errorf("RFIS not supported")
	case RDOT:
		return encoder.PushByte(1)
	case RKSM:
		return encoder.PushByte(2)
	case RATOM:
		return encoder.PushByte(3)
	case RSOL:
		return encoder.PushByte(4)
	case RMATIC:
		return encoder.PushByte(5)
	case RBNB:
		return encoder.PushByte(6)
	case RETH:
		return encoder.PushByte(7)
	default:
		return fmt.Errorf("RSymbol %s not supported", r)
	}
}

// used in db of rtoken-info
func (r RSymbol) ToRtokenType() int8 {
	switch r {
	case RFIS:
		return (0)
	case RDOT:
		return (1)
	case RKSM:
		return (2)
	case RATOM:
		return (3)
	case RSOL:
		return (4)
	case RMATIC:
		return (5)
	case RBNB:
		return (6)
	case RETH:
		return (-1)
	default:
		return -2
	}
}

func (r RSymbol) ToString() string {
	return string(r)
}
