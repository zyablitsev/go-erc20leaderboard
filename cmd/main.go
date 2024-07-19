package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"leaderboard"
	"leaderboard/pkg/tlscfg"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("main: can't load configuration: %w", err))
	}

	// init http client
	tlsCfg := &tls.Config{
		RootCAs: tlscfg.ClientCertPool(),
	}
	roundtripper := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: cfg.dialTimeout,
		}).DialContext,
		TLSHandshakeTimeout: cfg.tlsHandshakeTimeout,
		TLSClientConfig:     tlsCfg,
	}
	client := &http.Client{
		Transport: roundtripper,
		Timeout:   cfg.httpClientTimeout,
	}

	// get most recent block number
	blockNumber, err := leaderboard.GetBlockNumber(client, cfg.rpcURL)
	if err != nil {
		log.Fatal(fmt.Errorf("main: can't get most recent block number: %w", err))
	}
	log.Printf("most resent block number: %d\n", blockNumber)

	fromBlockNumber := blockNumber
	if cfg.depth > 0 && cfg.depth < fromBlockNumber {
		fromBlockNumber -= cfg.depth
	} else if cfg.depth > 0 {
		fromBlockNumber = 0
	}
	log.Printf("fetching log records from %d block number\n", fromBlockNumber)

	// get logs for last blocks for depth specified with from block number
	records, err := leaderboard.GetLogs(client, cfg.rpcURL, fromBlockNumber)
	if err != nil {
		log.Fatal(fmt.Errorf("main: can't get logs: %w", err))
	}
	log.Printf("got %d log records\n", len(records))

	// get top5 active addresses from the log records
	top, err := leaderboard.GetTop5(records)
	if err != nil {
		log.Fatal(fmt.Errorf("main: can't get top5 from logs: %w", err))
	}

	log.Printf("top%d leaderboard:\n", len(top))
	for i := range top {
		log.Println(top[i].Address, top[i].Activity)
	}
}

const (
	defaultDialTimeout         = "15s"
	defaultTLSHandshakeTimeout = "15s"
	defaultHTTPClientTimeout   = "15s"
	defaultDepth               = "100"
)

type config struct {
	rpcURL              string
	dialTimeout         time.Duration
	tlsHandshakeTimeout time.Duration
	httpClientTimeout   time.Duration
	depth               int64
}

func loadConfig() (config, error) {
	rpcURLEnv := os.Getenv("RPC_URL")
	if rpcURLEnv == "" {
		err := errors.New("loadConfig: RPC_URL env value is required")
		return config{}, err
	}

	dialTimeoutEnv := os.Getenv("DIAL_TIMEOUT")
	if dialTimeoutEnv == "" {
		dialTimeoutEnv = defaultDialTimeout
	}
	dialTimeout, err := time.ParseDuration(dialTimeoutEnv)
	if err != nil {
		err = fmt.Errorf(
			"loadConfig: bad DIAL_TIMEOUT env value %q: %w",
			dialTimeoutEnv, err)
		return config{}, err
	}

	tlsHandshakeTimeoutEnv := os.Getenv("TLS_HANDSHAKE_TIMEOUT")
	if tlsHandshakeTimeoutEnv == "" {
		tlsHandshakeTimeoutEnv = defaultTLSHandshakeTimeout
	}
	tlsHandshakeTimeout, err := time.ParseDuration(tlsHandshakeTimeoutEnv)
	if err != nil {
		err = fmt.Errorf(
			"loadConfig: bad TLS_HANDSHAKE_TIMEOUT env value %q: %w",
			tlsHandshakeTimeoutEnv, err)
		return config{}, err
	}

	httpClientTimeoutEnv := os.Getenv("HTTP_CLIENT_TIMEOUT")
	if httpClientTimeoutEnv == "" {
		httpClientTimeoutEnv = defaultHTTPClientTimeout
	}
	httpClientTimeout, err := time.ParseDuration(httpClientTimeoutEnv)
	if err != nil {
		err = fmt.Errorf(
			"loadConfig: bad HTTP_CLIENT_TIMEOUT env value %q: %w",
			httpClientTimeoutEnv, err)
		return config{}, err
	}

	depthEnv := os.Getenv("DEPTH")
	if depthEnv == "" {
		depthEnv = defaultDepth
	}
	depth, err := strconv.ParseInt(depthEnv, 10, 64)
	if err != nil {
		err = fmt.Errorf(
			"loadConfig: bad DEPTH env value %q: %w",
			depthEnv, err)
		return config{}, err
	}
	if depth < 0 {
		depth = 0
	}

	cfg := config{
		rpcURL:              rpcURLEnv,
		dialTimeout:         dialTimeout,
		tlsHandshakeTimeout: tlsHandshakeTimeout,
		httpClientTimeout:   httpClientTimeout,
		depth:               depth,
	}

	return cfg, nil
}
