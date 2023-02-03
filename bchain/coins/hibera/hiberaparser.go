package hibera

import (
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0xcba4b73e
)

// chain parameters
var (
	MainNetParams chaincfg.Params
)

func init() {
	// hibera mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{40} // starting with 'H'
	MainNetParams.ScriptHashAddrID = []byte{60} // starting with 'R'
	MainNetParams.PrivateKeyID = []byte{122}

}

// hiberaParser handle
type hiberaParser struct {
	*btc.BitcoinParser
	baseparser                         *bchain.BaseParser
}

// NewhiberaParser returns new hiberaParser instance
func NewhiberaParser(params *chaincfg.Params, c *btc.Configuration) *hiberaParser {
	p := &hiberaParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main hibera network
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
func (p *hiberaParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *hiberaParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
