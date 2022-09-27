package cryptoshares

import (
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/NikunjVaghasiya01/blockbook/bchain"
	"github.com/NikunjVaghasiya01/blockbook/bchain/coins/btc"
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
	// cryptoshares mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{63} // starting with 's'
	MainNetParams.ScriptHashAddrID = []byte{28} // starting with 'c'
	MainNetParams.PrivateKeyID = []byte{122}

}

// cryptosharesParser handle
type cryptosharesParser struct {
	*btc.BitcoinParser
	baseparser                         *bchain.BaseParser
}

// NewcryptosharesParser returns new cryptosharesParser instance
func NewcryptosharesParser(params *chaincfg.Params, c *btc.Configuration) *cryptosharesParser {
	p := &cryptosharesParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main cryptoshares network
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
func (p *cryptosharesParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *cryptosharesParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
