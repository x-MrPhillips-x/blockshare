package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
)

// Pull the data from the block
// Create a counter nonce starts at 0
// Create hash of the data + the counter
// check to see if meets requirements

// Difficulty
// todo more research
// difficulty can increase with each state
const Difficulty = 1

// ProofOfStake contains the data required to validate block
type ProofOfStake struct {
	Block        *Block
	Requirements *Requirements
}

type Requirements struct {
	OnOffenderList      bool
	CarInsurance        bool
	ConfirmRequirements bool
	// is not already on the platform
}

func (pos *ProofOfStake) Run() (int, string) {
	var hashed []byte
	var nonce int
	for nonce < math.MaxInt64 {
		data := pos.InitData(nonce)
		record := fmt.Sprintf("%d%s%s%s", pos.Block.Nonce, pos.Block.Timestamp, data, pos.Block.PrevBlockHash)
		h := sha256.New()
		h.Write([]byte(record))
		hashed = h.Sum(nil)

		// if the requirements are met the break
		// else try again with another nonce
		if hasCarInsurance() {
			break
		} else {
			nonce++
		}

	}
	return nonce, hex.EncodeToString(hashed)
}

func NewProof(b *Block) *ProofOfStake {
	return &ProofOfStake{
		Block: b,
	}
}

func (pos *ProofOfStake) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pos.Block.PrevBlockHash),
			pos.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

func isOnOffenderList() bool {
	return false
}

// hasCarInsurance determines if the driver has the required car insurance for this month
// todo random true false for the proof of stake
func hasCarInsurance() bool {
	return true
}

func validatorHasConfirmedRequirements() bool {
	return true
}

func ToHex(i int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, i); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func (pos *ProofOfStake) Validate() bool {
	return !isOnOffenderList() && hasCarInsurance()
	// TODO this is from Tensor revisit this?
	// var intHash big.Int
	// data := pos.InitData(pos.Block.Nonce)

	// hash := sha256.Sum256(data)
	// intHash.SetBytes(hash[:])

	// return intHash.Cmp(pos.Target) == -1
}
