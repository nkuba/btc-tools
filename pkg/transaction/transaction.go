package transaction

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/nkuba/btc-tools/pkg/chain"
)

type TransactionData struct {
	// Source transaction
	SourceTxHash         string `json:"sourceTxHash"`
	SourceTxOutputIndex  uint32 `json:"sourceTxOutputIndex"`
	SourceTxOutputAmount uint64 `json:"sourceTxOutputAmount"`
	SourceTxOutputScript string `json:"sourceTxOutputScript"`
	// Output 1
	DestinationAddress1 string `json:"output1Address"`
	FundingAmount       uint64 `json:"output1Amount"`
	// Output 2
	DestinationAddress2 string `json:"output2Address"`
	Fee                 uint64 `json:"fee"`
}

func CreateAndPublish(
	txData *TransactionData,
	chain *chain.BTC,
	signer *Signer,
	networkParams *chaincfg.Params,
) (string, error) {
	// Calculate value for output 2.
	outputAmount2 := txData.SourceTxOutputAmount -
		txData.FundingAmount -
		txData.Fee

	// Create unsigned transaction
	msgTx, err := CreateUnsigned(
		txData.SourceTxHash,
		txData.SourceTxOutputIndex,
		txData.SourceTxOutputAmount,
		txData.DestinationAddress1,
		txData.DestinationAddress2,
		txData.FundingAmount,
		outputAmount2,
		networkParams,
	)
	if err != nil {
		log.Fatalf("cannot create unsigned transaction: [%s]", err)
	}

	// Sign transaction.
	inputToSignIndex := uint32(0) // we support only transactions with one input

	sourceOutputScript, err := hex.DecodeString(txData.SourceTxOutputScript)
	if err != nil {
		log.Fatalf("cannot decode subscript: [%s]", err)
	}

	if txscript.IsWitnessProgram(sourceOutputScript) {
		signer.SignWitness(
			msgTx,
			sourceOutputScript,
			inputToSignIndex,
			txData.SourceTxOutputAmount,
		)
	} else {
		signer.SignNoWitness(
			msgTx,
			sourceOutputScript,
			inputToSignIndex,
		)
	}

	// Publish transaction.
	rawTransaction, err := Serialize(msgTx)
	if err != nil {
		return "", fmt.Errorf("transaction serialization failed: [%s]", err)
	}

	transactionHash, err := chain.PublishTransaction(rawTransaction)
	if err != nil {
		log.Fatalf("transaction publication failed: [%v]", err)
	}

	return transactionHash, nil
}

func ReadTransactionData(filePath string) (*TransactionData, error) {
	var txData *TransactionData

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &txData)

	return txData, nil
}

func CreateUnsigned(
	sourceTxHashString string,
	sourceTxOutputIndex uint32,
	sourceTxOutputAmount uint64,
	outputAddress1 string,
	outputAddress2 string,
	outputAmount1 uint64,
	outputAmount2 uint64,
	networkParams *chaincfg.Params,
) (*wire.MsgTx, error) {
	var msgTx *wire.MsgTx

	addTxOutput := func(address string, amount uint64) error {
		addr, err := btcutil.DecodeAddress(address, networkParams)
		if err != nil {
			return fmt.Errorf("address decoding failed [%s]", err)
		}

		outScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return fmt.Errorf("script to pay creation failed [%s]", err)
		}

		msgTx.AddTxOut(wire.NewTxOut(int64(amount), outScript))

		return nil
	}

	// Initialize transaction message
	msgTx = wire.NewMsgTx(wire.TxVersion)

	// Create Transaction input based on the output of historic transaction
	sourceTxHash, err := chainhash.NewHashFromStr(sourceTxHashString)
	if err != nil {
		return nil, fmt.Errorf("cannot create source transaction hash [%s]", err)
	}
	txIn := wire.NewTxIn(
		wire.NewOutPoint(sourceTxHash, sourceTxOutputIndex),
		nil,
		nil)
	msgTx.AddTxIn(txIn)

	// Transaction Output 1.
	if err := addTxOutput(outputAddress1, outputAmount1); err != nil {
		return nil, fmt.Errorf("output 1 creation failed [%s]", err)
	}

	// Transaction Output 2.
	if err := addTxOutput(outputAddress2, outputAmount2); err != nil {
		return nil, fmt.Errorf("output 2 creation failed [%s]", err)
	}

	return msgTx, nil
}

// Serialize encodes a bitcoin transaction message to a hexadecimal format.
func Serialize(msgTx *wire.MsgTx) ([]byte, error) {
	var buffer bytes.Buffer

	err := msgTx.Serialize(&buffer)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize transaction [%s]", err)
	}

	return buffer.Bytes(), nil
}
