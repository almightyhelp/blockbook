package oduwausd

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
//	"io"
	"math/big"

	"github.com/juju/errors"
	"github.com/martinboehm/btcd/blockchain"
	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/almightyhelp/blockbook/bchain"
	"github.com/almightyhelp/blockbook/bchain/coins/btc"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0xad6cb7e3
	TestnetMagic wire.BitcoinNet = 0xba657645

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
	baseparser                         *bchain.BaseParser
	BitcoinOutputScriptToAddressesFunc btc.OutputScriptToAddressesFunc
}

// NewoduwausdParser returns new oduwausdParser instance
func NewoduwausdParser(params *chaincfg.Params, c *btc.Configuration) *oduwausdParser {
	p := &oduwausdParser{
		BitcoinParser: btc.NewBitcoinParser(params, c),
		baseparser:    &bchain.BaseParser{},
	}
	p.BitcoinOutputScriptToAddressesFunc = p.OutputScriptToAddressesFunc
	p.OutputScriptToAddressesFunc = p.outputScriptToAddresses
	return p
}

// GetChainParams contains network parameters for the main oduwausd network
func GetChainParams(chain string) *chaincfg.Params {
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

// PackTx packs transaction to byte array using protobuf
func (p *oduwausdParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

// UnpackTx unpacks transaction from protobuf byte array
func (p *oduwausdParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}

func (p *oduwausdParser) ParseBlock(b []byte) (*bchain.Block, error) {
	r := bytes.NewReader(b)
	w := wire.MsgBlock{}
	h := wire.BlockHeader{}
	err := h.Deserialize(r)
	if err != nil {
		return nil, errors.Annotatef(err, "Deserialize")
	}

//	if h.Version > 3 && h.Version < 7 {
		// Skip past AccumulatorCheckpoint (block version 4, 5 and 6)
//		r.Seek(32, io.SeekCurrent)
//	}

//	err = utils.DecodeTransactions(r, 0, wire.WitnessEncoding, &w)
//	if err != nil {
//		return nil, errors.Annotatef(err, "DecodeTransactions")
//	}

	txs := make([]bchain.Tx, len(w.Transactions))
	for ti, t := range w.Transactions {
		txs[ti] = p.TxFromMsgTx(t, false)
	}

	return &bchain.Block{
		BlockHeader: bchain.BlockHeader{
			Size: len(b),
			Time: h.Timestamp.Unix(),
		},
		Txs: txs,
	}, nil
}


// ParseTx parses byte array containing transaction and returns Tx struct
func (p *oduwausdParser) ParseTx(b []byte) (*bchain.Tx, error) {
	t := wire.MsgTx{}
	r := bytes.NewReader(b)
	if err := t.Deserialize(r); err != nil {
		return nil, err
	}
	tx := p.TxFromMsgTx(&t, true)
	tx.Hex = hex.EncodeToString(b)
	return &tx, nil
}

// TxFromMsgTx parses tx and adds handling for OP_ZEROCOINSPEND inputs
func (p *oduwausdParser) TxFromMsgTx(t *wire.MsgTx, parseAddresses bool) bchain.Tx {
	vin := make([]bchain.Vin, len(t.TxIn))
	for i, in := range t.TxIn {

		// extra check to not confuse Tx with single OP_ZEROCOINSPEND input as a coinbase Tx
		if blockchain.IsCoinBaseTx(t) {
			vin[i] = bchain.Vin{
				Coinbase: hex.EncodeToString(in.SignatureScript),
				Sequence: in.Sequence,
			}
			break
		}

		s := bchain.ScriptSig{
			Hex: hex.EncodeToString(in.SignatureScript),
			// missing: Asm,
		}

		txid := in.PreviousOutPoint.Hash.String()

		vin[i] = bchain.Vin{
			Txid:      txid,
			Vout:      in.PreviousOutPoint.Index,
			Sequence:  in.Sequence,
			ScriptSig: s,
		}
	}
	vout := make([]bchain.Vout, len(t.TxOut))
	for i, out := range t.TxOut {
		addrs := []string{}
		if parseAddresses {
			addrs, _, _ = p.OutputScriptToAddressesFunc(out.PkScript)
		}
		s := bchain.ScriptPubKey{
			Hex:       hex.EncodeToString(out.PkScript),
			Addresses: addrs,
			// missing: Asm,
			// missing: Type,
		}
		var vs big.Int
		vs.SetInt64(out.Value)
		vout[i] = bchain.Vout{
			ValueSat:     vs,
			N:            uint32(i),
			ScriptPubKey: s,
		}
	}
	tx := bchain.Tx{
		Txid:     t.TxHash().String(),
		Version:  t.Version,
		LockTime: t.LockTime,
		Vin:      vin,
		Vout:     vout,
		// skip: BlockHash,
		// skip: Confirmations,
		// skip: Time,
		// skip: Blocktime,
	}
	return tx
}

// ParseTxFromJson parses JSON message containing transaction and returns Tx struct
func (p *oduwausdParser) ParseTxFromJson(msg json.RawMessage) (*bchain.Tx, error) {
	var tx bchain.Tx
	err := json.Unmarshal(msg, &tx)
	if err != nil {
		return nil, err
	}

	for i := range tx.Vout {
		vout := &tx.Vout[i]
		// convert vout.JsonValue to big.Int and clear it, it is only temporary value used for unmarshal
		vout.ValueSat, err = p.AmountToBigInt(vout.JsonValue)
		if err != nil {
			return nil, err
		}
		vout.JsonValue = ""

		if vout.ScriptPubKey.Addresses == nil {
			vout.ScriptPubKey.Addresses = []string{}
		}
	}

	return &tx, nil
}

// outputScriptToAddresses converts ScriptPubKey to bitcoin addresses
func (p *oduwausdParser) outputScriptToAddresses(script []byte) ([]string, bool, error) {
	rv, s, _ := p.BitcoinOutputScriptToAddressesFunc(script)
	return rv, s, nil
}

func (p *oduwausdParser) GetAddrDescForUnknownInput(tx *bchain.Tx, input int) bchain.AddressDescriptor {
	if len(tx.Vin) > input {
		scriptHex := tx.Vin[input].ScriptSig.Hex

		if scriptHex != "" {
			script, _ := hex.DecodeString(scriptHex)
			return script
		}
	}

	s := make([]byte, 10)
	return s
}