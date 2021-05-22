package oduwacoin

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
)

// oduwacoinRPC is an interface to JSON-RPC bitcoind service.
type oduwacoinRPC struct {
	*btc.BitcoinRPC
}

// NewoduwacoinRPC returns new oduwacoinRPC instance.
func NewoduwacoinRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &oduwacoinRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV1{}
	s.ChainConfig.SupportsEstimateFee = true
	s.ChainConfig.SupportsEstimateSmartFee = false

	return s, nil
}

// Initialize initializes oduwacoinRPC instance.
func (b *oduwacoinRPC) Initialize() error {
	ci, err := b.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain

	glog.Info("Chain name ", chainName)
	params := GetChainParams(chainName)

	// always create parser
	b.Parser = NewoduwacoinParser(params, b.ChainConfig)

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
