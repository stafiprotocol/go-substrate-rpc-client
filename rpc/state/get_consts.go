package state

import (
	"fmt"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

const (
	metaV10 = 10
	metaV11 = 11
	metaV12 = 12
)

func (s *State) GetConst(prefix, name string, res interface{}) error {
	meta, err := s.GetMetadataLatest()
	if err != nil {
		return err
	}

	return s.GetConstWithMetadata(meta, prefix, name, res)
}

func (s *State) GetConstWithMetadata(meta *types.Metadata, prefix, name string, res interface{}) error {
	switch meta.Version {
	case metaV12:
		return meta.AsMetadataV12.GetConst(prefix, name, res)
	case metaV11:
		return meta.AsMetadataV11.GetConst(prefix, name, res)
	case metaV10:
		return meta.AsMetadataV10.GetConst(prefix, name, res)
	default:
		return fmt.Errorf("unsupported metadata version")
	}
}
