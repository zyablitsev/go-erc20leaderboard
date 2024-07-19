package proto

import (
	"encoding/json"
)

// EthGetLogsParams object
type EthGetLogsParams struct {
	FromBlock string   `json:"fromBlock,omitempty"`
	ToBlock   string   `json:"toBlock,omitempty"`
	Address   []string `json:"address,omitempty"`
	Topics    []string `json:"topics,omitempty"`
	BlockHash string   `json:"blockhash,omitempty"`
}

// Marshall returns the json encoding of EthGetLogsParams
func (s EthGetLogsParams) Marshal() json.RawMessage {
	v := []EthGetLogsParams{s}
	b, _ := json.Marshal(v)
	return json.RawMessage(b)
}

// EthLogRecord object
type EthLogRecord struct {
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}
