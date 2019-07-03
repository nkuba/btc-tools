package address

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

// WIFtoPublicKey decodes WIF from provided string value and returns a public
// key associated with the WIF.
func WIFtoPublicKey(wifString string) (*btcec.PublicKey, error) {
	wif, err := btcutil.DecodeWIF(wifString)
	if err != nil {
		return nil, fmt.Errorf("cannot decode WIF [%s]", err)
	}
	return (*btcec.PublicKey)(&wif.PrivKey.PublicKey), nil
}

// PublicKeyToWPKHaddress calculates Witness Public Key Hash address for the
// given public key. It returns the address encoded to a string.
func PublicKeyToWPKHaddress(
	publicKey *btcec.PublicKey,
	networkParams *chaincfg.Params,
) (string, error) {
	compressedPublicKey := publicKey.SerializeCompressed()

	// Hash ripemd160(sha256(compressedPublicKey)).
	witnessProgram := btcutil.Hash160(compressedPublicKey)

	address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, networkParams)
	if err != nil {
		return "", err
	}

	return address.EncodeAddress(), nil
}
