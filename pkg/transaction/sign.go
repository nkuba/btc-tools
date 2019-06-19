package transaction

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

type Signer struct {
	privateKey *btcec.PrivateKey
}

// NewSigner initializes a new signer with a private key decoded from WIF.
func NewSigner(wif string) (*Signer, error) {
	privateKey, err := privateKeyFromWIF(wif)
	if err != nil {
		return nil, err
	}

	return &Signer{

		privateKey: privateKey,
	}, nil
}

// privateKeyFromWIF decodes WIF and returns a private key from this WIF.
func privateKeyFromWIF(wif string) (*btcec.PrivateKey, error) {
	decodedWIF, err := btcutil.DecodeWIF(wif)
	if err != nil {
		return nil, fmt.Errorf("cannot decode WIF [%s]", err)
	}
	return decodedWIF.PrivKey, nil
}

func (s *Signer) SignNoWitness(
	msgTx *wire.MsgTx,
	subscript []byte,
	inputToSignIndex uint32,
) error {
	sigScript, err := txscript.SignatureScript(
		msgTx,
		int(inputToSignIndex),
		subscript,
		txscript.SigHashAll,
		s.privateKey,
		true,
	)
	if err != nil {
		return fmt.Errorf("signature script calculation failed: [%s]", err)
	}

	msgTx.TxIn[inputToSignIndex].SignatureScript = sigScript

	return nil
}

func (s *Signer) SignWitness(
	msgTx *wire.MsgTx,
	sourceOutputScript []byte,
	inputToSignIndex uint32,
	sourceTxOutputAmount uint64,
) error {
	witness, err := txscript.WitnessSignature(
		msgTx,
		txscript.NewTxSigHashes(msgTx),
		int(inputToSignIndex),
		int64(sourceTxOutputAmount),
		sourceOutputScript,
		txscript.SigHashAll,
		s.privateKey,
		false,
	)
	if err != nil {
		return fmt.Errorf("witness signature calculation failed: [%s]", err)
	}

	msgTx.TxIn[inputToSignIndex].Witness = witness

	return nil
}
