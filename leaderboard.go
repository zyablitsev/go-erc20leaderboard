package leaderboard

import (
	"bytes"
	"encoding/hex"
	"hash"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"

	"leaderboard/proto"
)

const (
	hashLength    = 32
	addressLength = 20

	transferTopicsLength = 3

	methodTopicIdx = 0
	fromTopicIdx   = 1
	toTopicIdx     = 2

	sigTransfer = "Transfer(address,address,uint256)"
	rpcVersion  = "2.0"
)

var (
	sigTransferBytes []byte = nil
	sigTransferHEX   string = ""
)

func init() {
	// init transfer method signature
	sha := sha3.NewLegacyKeccak256().(keccakState)
	sha.Write([]byte(sigTransfer))
	sigTransferBytes = make([]byte, hashLength)
	sha.Read(sigTransferBytes)

	enc := make([]byte, len(sigTransferBytes)*2+2)
	copy(enc, "0x")
	hex.Encode(enc[2:], sigTransferBytes)
	sigTransferHEX = string(enc)
}

// GetTop5 returns list of top active addresses based on log records
func GetTop5(records []proto.EthLogRecord) ([]*AddressActivity, error) {
	m := make(map[string]*AddressActivity, len(records)*2)
	s := make([]*AddressActivity, 0, len(records)*2)

	incActivity := func(h string) error {
		h = removePrefix(h)
		b, err := hex.DecodeString(h)
		if err != nil {
			return err
		}
		h = bytesToAddressHex(b)
		if v, ok := m[h]; ok {
			v.Activity++
		} else {
			item := &AddressActivity{
				Address:  h,
				Activity: 1,
			}
			m[h] = item
			s = append(s, item)
		}

		return nil
	}

	for i := 0; i < len(records); i++ {
		if len(records[i].Topics) < transferTopicsLength {
			continue
		}

		// check method topic
		methodHex := records[i].Topics[methodTopicIdx]
		methodHex = removePrefix(methodHex)
		methodBytes, err := hex.DecodeString(methodHex)
		if err != nil {
			return nil, err
		}
		if bytes.Compare(sigTransferBytes, methodBytes) != 0 {
			continue
		}

		// handle from
		err = incActivity(records[i].Topics[fromTopicIdx])
		if err != nil {
			return nil, err
		}

		// handle to
		err = incActivity(records[i].Topics[toTopicIdx])
		if err != nil {
			return nil, err
		}
	}

	// sort slice descend and keep up to 5 most active addresses
	sort.Slice(s, func(i, j int) bool { return s[i].Activity > s[j].Activity })
	if len(s) > 4 {
		s = s[:5]
	}

	return s, nil
}

// AddressActivity object
type AddressActivity struct {
	Address  string
	Activity int
}

func bytesToAddressHex(b []byte) string {
	addr := [addressLength]byte{}
	if len(b) > len(addr) {
		b = b[len(b)-addressLength:]
	}
	copy(addr[addressLength-len(b):], b)

	var buf [len(addr)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], addr[:])
	b = buf[:]

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(b[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(b); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if b[i] > '9' && hashByte > 7 {
			b[i] -= 32
		}
	}

	return string(b[:])
}

func removePrefix(s string) string {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return s
}

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}
