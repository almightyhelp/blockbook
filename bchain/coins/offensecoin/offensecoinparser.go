package offensecoin

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
	// offensecoin mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{115} // starting with 'o'
	MainNetParams.ScriptHashAddrID = []byte{13}
	MainNetParams.PrivateKeyID = []byte{212}

}

// offensecoinParser handle
type offensecoinParser struct {
	*btc.BitcoinParser
	baseparser                         *bchain.BaseParser
}

// NewoffensecoinParser returns new offensecoinParser instance
func NewoffensecoinParser(params *chaincfg.Params, c *btc.Configuration) *offensecoinParser {
	p := &offensecoinParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main offensecoin network
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
func (p *offensecoinParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *offensecoinParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
