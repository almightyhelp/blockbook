package voltnote

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
	// voltnote mainnet Address encoding magics
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic
	MainNetParams.PubKeyHashAddrID = []byte{103} // starting with 'v'
	MainNetParams.ScriptHashAddrID = []byte{13}
	MainNetParams.PrivateKeyID = []byte{212}

}

// voltnoteParser handle
type voltnoteParser struct {
	*btc.BitcoinParser
	baseparser                         *bchain.BaseParser
}

// NewvoltnoteParser returns new voltnoteParser instance
func NewvoltnoteParser(params *chaincfg.Params, c *btc.Configuration) *voltnoteParser {
	p := &voltnoteParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	return p
}

// GetChainParams contains network parameters for the main voltnote network
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
func (p *voltnoteParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *voltnoteParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}
