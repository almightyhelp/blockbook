package oduwausd

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
)

// oduwausdRPC is an interface to JSON-RPC bitcoind service.
type oduwausdRPC struct {
	*btc.BitcoinRPC
}

// NewoduwausdRPC returns new oduwausdRPC instance.
func NewoduwausdRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &oduwausdRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV1{}
	s.ChainConfig.SupportsEstimateFee = true
	s.ChainConfig.SupportsEstimateSmartFee = false

	return s, nil
}

// Initialize initializes oduwausdRPC instance.
func (b *oduwausdRPC) Initialize() error {
	ci, err := b.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	// always create parser
	b.Parser = NewoduwausdParser(params, b.ChainConfig)

	// parameters for getInfo request
	if params.Net == MainnetMagic {
		b.Testnet = false
		b.Network = "livenet"
	} else {
		b.Testnet = true
		b.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}