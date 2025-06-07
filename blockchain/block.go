package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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

// Validator stake will increase with each ride and/or driver transaction
type Validator struct {
	UUID  uuid.UUID
	Stake int
}

func CreateBlock(data []byte, prevHash string) *Block {
	block := &Block{
		Timestamp:     time.Now(),
		Hash:          "",
		Data:          data,
		PrevBlockHash: prevHash,
		Nonce:         0,
	}

	pos := NewProof(block)
	nonce, hash := pos.Run()
	block.Hash = string(hash)
	block.Nonce = nonce
	return block
}

func Genesis(data []byte) *Block {
	return CreateBlock(data, "")
}

func (b *Block) calculateHash() string {
	var record string
	record = fmt.Sprintf("%d%s%s%s", b.Nonce, b.Timestamp, b.Data, b.PrevBlockHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Serialize block for badgerdb
// returns a byte representation of the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func Deserialize(data []byte) *Block {
	var block *Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(block); err != nil {
		log.Panic(err)
	}

	return block
}
