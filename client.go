package leaderboard

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"leaderboard/proto"
)

// GetBlockNumber returns the number of most recent block
func GetBlockNumber(client *http.Client, url string) (int64, error) {
	request := proto.Request{
		Version: rpcVersion,
		Method:  "eth_blockNumber",
		ID:      rand.Intn(100),
	}

	response, err := doRequest(client, url, request)
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}

	// unpack response payload
	s := ""
	err = json.Unmarshal(response.Result, &s)
	if err != nil {
		return 0, err
	}

	s = removePrefix(s) // remove 0x prefix

	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// GetLogs return an array of logs for erc20 tokens transfers
// from the specified block to latest
func GetLogs(
	client *http.Client,
	url string,
	fromBlockNumber int64,
) ([]proto.EthLogRecord, error) {
	fromBlock := "0x" + strconv.FormatInt(fromBlockNumber, 16)
	params := proto.EthGetLogsParams{
		FromBlock: fromBlock,
		Topics:    []string{sigTransferHEX},
	}
	request := proto.Request{
		Version: rpcVersion,
		Method:  "eth_getLogs",
		Params:  params.Marshal(),
		ID:      rand.Intn(100),
	}

	response, err := doRequest(client, url, request)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	// unpack response payload
	result := []proto.EthLogRecord{}
	err = json.Unmarshal(response.Result, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func doRequest(
	client *http.Client,
	url string,
	request proto.Request,
) (proto.Response, error) {
	req, err := http.NewRequest(
		"POST", url, bytes.NewBuffer(request.Marshal()))
	if err != nil {
		return proto.Response{}, err
	}
	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return proto.Response{}, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return proto.Response{}, err
	}

	response := proto.Response{}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return proto.Response{}, err
	}

	return response, nil
}
