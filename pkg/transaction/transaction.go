package transaction

import (
	"encoding/json"
	"io/ioutil"
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

func ReadTransactionData(filePath string) (*TransactionData, error) {
	var txData *TransactionData

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &txData)

	return txData, nil
}
