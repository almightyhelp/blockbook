package tnotes

import (
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0x6ae31bc4
)

// chain parameters
var (
	MainNetParams chaincfg.Params
)

func init() {
	// tnotes mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{103} // starting with 't'
	MainNetParams.ScriptHashAddrID = []byte{13}
	MainNetParams.PrivateKeyID = []byte{212}

}

// tnotesParser handle
type tnotesParser struct {
	*btc.BitcoinParser
	baseparser                         *bchain.BaseParser
}

// NewtnotesParser returns new tnotesParser instance
func NewtnotesParser(params *chaincfg.Params, c *btc.Configuration) *tnotesParser {
	p := &tnotesParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main tnotes network
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	default:
		return &MainNetParams
	}
}

// PackTx packs transaction to byte array using protobuf
func (p *tnotesParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *tnotesParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
