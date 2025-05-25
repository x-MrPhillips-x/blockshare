package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var now = time.Now

type Block struct {
	Timestamp     time.Time
	Data          []byte // TODO will hold a transaction
	PrevBlockHash string
	Hash          string
	Nonce         int
	Validators    []Validator
}

type BlockChain struct {
	Blocks []Block
}

// Validator stake will increase with each ride and/or driver transaction
type Validator struct {
	UUID  uuid.UUID
	Stake int
}

func NewBlockChain() *BlockChain {
	return &BlockChain{Blocks: []Block{createGenesis([]byte("each one teach one"))}}
}

func (bc *BlockChain) AddBlock(data string, validator Validator) { // for now adding block with simple string
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Timestamp:     now(),
		PrevBlockHash: prevBlock.Hash,
		Hash:          "",
		Data:          []byte(data),
	}

	// todo proof of stake
	pos := NewProof(newBlock)
	pos.Run()

	// todo validate
	if pos.Requirements.OnOffenderList || !pos.Requirements.CarInsurance {
		return
	}

	if !pos.Requirements.ConfirmRequirements {
		return
	}

	newBlock.Hash = newBlock.calculateHash(validator)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func createGenesis(data []byte) Block {
	b := Block{
		Timestamp:     now(),
		Data:          data,
		PrevBlockHash: "",
		Hash:          "",
		Nonce:         0,
	}

	b.Hash = b.calculateHash(Validator{
		UUID: uuid.Nil,
	})

	return b
}

func (b *Block) calculateHash(validator Validator) string {
	var record string
	if validator.UUID.String() == "00000000-0000-0000-0000-000000000000" {
		record = fmt.Sprintf("%d%s%s%s", b.Nonce, b.Timestamp, b.Data, b.PrevBlockHash)

	} else {
		record = fmt.Sprintf("%d%s%s%s%s", b.Nonce, b.Timestamp, b.Data, b.PrevBlockHash, validator.UUID)

	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
