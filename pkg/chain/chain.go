// Package chain contains implementation of the chain interface communicating
// with [Block Cypher API](https://www.blockcypher.com/dev/bitcoin/).
package chain

import (
	"encoding/hex"

	"github.com/blockcypher/gobcy"
)

type BTC struct {
	api gobcy.API
}

// Config contains configuration for Block Cypher API.
type Config struct {
	// Token is Block Cypher's user token required for access to POST and DELETE
	// calls on the API.
	Token string
	Coin  string // Options: "btc", "bcy", "ltc", "doge"
	Chain string // Options: "main", "test3", "test"
}

// PublishTransaction sends a raw transaction to Block Cypher's API. It returns
// a transaction hash as a hexadecimal string.
func (b *BTC) PublishTransaction(rawTransaction []byte) (string, error) {
	tx, err := b.api.PushTX(hex.EncodeToString(rawTransaction))
	if err != nil {
		return "", err
	}

	return tx.Trans.Hash, nil
}

// Connect performs initialization for communication with Block Cypher based on
// provided config.
func Connect(config *Config) (*BTC, error) {
	blockCypherAPI := gobcy.API{
		Token: config.Token,
		Coin:  config.Coin,
		Chain: config.Chain,
	}

	return &BTC{api: blockCypherAPI}, nil
}
