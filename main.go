package main

import (
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/nkuba/btc-tools/pkg/chain"
	"github.com/nkuba/btc-tools/pkg/config"
	"github.com/nkuba/btc-tools/pkg/transaction"
)

const configFilePath = "./configs/config.toml"
const transactionDataFilePath = "./configs/transaction.json"

func main() {
	config, err := config.ReadConfig(configFilePath)
	if err != nil {
		log.Fatalf("cannot read config: [%s]", err)
	}

	var networkParams *chaincfg.Params
	switch config.BlockCypher.Chain {
	case "mainnet":
		networkParams = &chaincfg.MainNetParams
	default:
		networkParams = &chaincfg.TestNet3Params
	}

	// Initialize new signer.
	signer, err := transaction.NewSigner(config.Signer.WIF)
	if err != nil {
		log.Fatalf("signer initialization failed: [%s]", err)
	}

	// Initialize connection to the bitcoin chain.
	chain, err := chain.Connect(&config.BlockCypher)
	if err != nil {
		log.Fatalf("chain initialization failed: [%s]", err)
	}

	// Read transaction details from file.
	transactionData, err := transaction.ReadTransactionData(transactionDataFilePath)
	if err != nil {
		log.Fatalf("cannot read the transaction details: [%s]", err)
	}

	// Create and publish transaction.
	transactionHash, err := transaction.CreateAndPublish(
		transactionData,
		chain,
		signer,
		networkParams,
	)
	if err != nil {
		log.Fatalf("transaction creation or publication failed: [%s]", err)
	}

	log.Printf("Published transaction hash: %v\n", transactionHash)
}
