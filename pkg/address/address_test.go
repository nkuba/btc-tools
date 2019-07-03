package address

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestWIFtoPublicKey(t *testing.T) {
	wif := "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"
	expectedPublicKey := publicKey()

	publicKey, err := WIFtoPublicKey(wif)
	if err != nil {
		t.Errorf("unexpected error: [%s]", err)
	}

	if !reflect.DeepEqual(publicKey, expectedPublicKey) {
		t.Errorf(
			"unexpected address\nexpected: [%v]\nactual:   [%v]",
			expectedPublicKey,
			publicKey,
		)
	}
}

func TestPublicKeyToWPKHaddress(t *testing.T) {
	publicKey := publicKey()
	expectedAddress := "bc1qmy63mjadtw8nhzl69ukdepwzsyvv4yex5qlmkd"

	networkParams := &chaincfg.MainNetParams

	address, err := PublicKeyToWPKHaddress(publicKey, networkParams)
	if err != nil {
		t.Fatalf("unexpected error: [%s]", err)
	}

	if address != expectedAddress {
		t.Fatalf(
			"unexpected address\nexpected: [%s]\nactual:   [%s]",
			expectedAddress,
			address,
		)
	}
}

func publicKey() *btcec.PublicKey {
	x, _ := new(big.Int).SetString("d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645c", 16)
	y, _ := new(big.Int).SetString("d85228a6fb29940e858e7e55842ae2bd115d1ed7cc0e82d934e929c97648cb0a", 16)
	return &btcec.PublicKey{
		X:     x,
		Y:     y,
		Curve: btcec.S256(),
	}
}
