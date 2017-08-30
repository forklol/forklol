package bitcoin

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

var rpcClient *rpcclient.Client

const BITSATO = 100000000

type txDetails struct {
	Type       string
	Size       int32
	SizeV      int32
	FeeSize    uint32
	FeeSizeV   uint32
	NumInputs  int
	NumOutputs int
	FeeSat     uint32
	Fee        float64
	AmountIn   float64
	AmountOut  float64
}

type blockTxDetails struct {
	TxRegular   uint32
	TxSegwit    uint32
	TxTotal     uint32
	SizeTotal   uint32
	SizeVTotal  uint32
	SizeMedian  uint32
	SizeVMedian uint32
	FeeTotal	uint32
	FeeMedian	uint32
	FeeVMedian	uint32
	NumInputs	int
	NumOutputs	int
}

func RPCConnect() error {
	connCfg := &rpcclient.ConnConfig{
		Host:         "172.17.0.2:8332",
		User:         "forklol",
		Pass:         "forklol",
		HTTPPostMode: true,
		DisableTLS:   true,
	}
	client, err := rpcclient.New(connCfg, nil)
	rpcClient = client

	return err
}

func RPCDisconnect() {
	rpcClient.Shutdown()
}

func GetBlockTxStats(blockhash string) error {
	hash, err := chainhash.NewHashFromStr(blockhash)
	if err != nil {
		log.Printf("Hash %s invalid: %s\n", blockhash, err.Error())
		return err
	}

	blk, err := rpcClient.GetBlockVerbose(hash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	types := map[string]int{
		"segwit":   0,
		"regular":  0,
		"error":    0,
		"coinbase": 0,
	}

	for _, txid := range blk.Tx {
		d := getTxDetails(txid)
		if d != nil {
			types[d.Type]++
		} else {
			types["error"]++
		}

		if d.Type == "segwit" || d.Type == "coinbase" {
			log.Println(d)
		}
	}

	log.Printf("Block %s had %d txs\n", blk.Hash, len(blk.Tx))
	log.Println(types)

	return nil
}

func getTxDetails(txid string) *txDetails {
	hash, _ := chainhash.NewHashFromStr(txid)
	tx, err := rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		log.Printf("%s: %s\n", txid, err.Error())
		return nil
	}

	d := txDetails{}
	d.Size = tx.Size
	d.SizeV = tx.Vsize
	d.Type = ""
	d.NumInputs = len(tx.Vin)
	d.NumOutputs = len(tx.Vout)

	for _, vin := range tx.Vin {
		if vin.Coinbase != "" {
			d.Type = "coinbase"
		}
		if vin.HasWitness() && d.Type == "" {
			d.Type = "segwit"
		}

		prevOut := getTxOutputAmount(vin.Txid, vin.Vout)
		d.AmountIn += prevOut
	}

	for _, vout := range tx.Vout {
		d.AmountOut += vout.Value
	}

	if d.Type != "coinbase" {
		fee := d.AmountIn - d.AmountOut

		if d.Type == "" {
			d.Type = "regular"
		}

		d.FeeSize = uint32(fee / float64(d.Size) * BITSATO)
		d.FeeSizeV = uint32(fee / float64(d.SizeV) * BITSATO)
		d.Fee = fee
		d.FeeSat = uint32(fee * BITSATO)
	}

	return &d
}

func getTxOutputAmount(txid string, n uint32) float64 {
	hash, _ := chainhash.NewHashFromStr(txid)
	tx, err := rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return 0.0
	}

	return tx.Vout[n].Value
}
