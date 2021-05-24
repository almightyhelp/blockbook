package oduwausd

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/big"

	"github.com/juju/errors"
	"github.com/martinboehm/btcd/blockchain"
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
	"github.com/almightyhelp/blockbook/bchain/coins/utils"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0xad6cb7e3
	TestnetMagic wire.BitcoinNet = 0xba657645

	// Zerocoin op codes
	OP_ZEROCOINMINT  = 0xc1
	OP_ZEROCOINSPEND = 0xc2
)

// chain parameters
var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

func init() {
	// oduwausd mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{115} // starting with 'D'
	MainNetParams.ScriptHashAddrID = []byte{68}
	MainNetParams.PrivateKeyID = []byte{183}

	// oduwausd testnet Address encoding magics
	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic
	TestNetParams.PubKeyHashAddrID = []byte{139} // starting with 'x' or 'y'
	TestNetParams.ScriptHashAddrID = []byte{19}
	TestNetParams.PrivateKeyID = []byte{239}
}

// oduwausdParser handle
type oduwausdParser struct {
	*btc.BitcoinParser
}

// NewoduwausdParser returns new oduwausdParser instance
func NewoduwausdParser(params *chaincfg.Params, c *btc.Configuration) *oduwausdParser {
	p := &oduwausdParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main oduwausd network
func GetChainParams(chain string) *chaincfg.Params {
	// register bitcoin parameters in addition to litecoin parameters
	// litecoin has dual standard of addresses and we want to be able to
	// parse both standards
	if !chaincfg.IsRegistered(&chaincfg.MainNetParams) {
		chaincfg.RegisterBitcoinParams()
	}
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err == nil {
			err = chaincfg.Register(&TestNetParams)
		}
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	case "test":
		return &TestNetParams
	default:
		return &MainNetParams
	}
}
