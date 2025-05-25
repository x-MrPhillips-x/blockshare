package blockchain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlock_calculateHash(t *testing.T) {
	firstValidatorUUID, _ := uuid.Parse("00000000-0000-0000-0000-00000000000")
	secondValidatorUUID, _ := uuid.Parse("43d03424-693a-4f90-a97d-fdeda3f23df1")
	type fields struct {
		Timestamp     time.Time
		Data          []byte
		PrevBlockHash string
		Hash          string
	}
	tests := []struct {
		name      string
		fields    fields
		Validator Validator
		want      string
	}{
		{
			name: "genesis block hash calculation",
			fields: fields{
				Timestamp:     time.Date(1911, 1, 5, 0, 0, 0, 0, time.UTC),
				Data:          []byte(`testing 1 2 3`),
				PrevBlockHash: "",
				Hash:          "235c07f768300fd8e42690c69172cec9e2125dced3cb26148b78ccf0f2ccc805",
			},
		},
		{
			name: "genesis +1 block hash calculation",
			fields: fields{
				Timestamp:     time.Date(1911, 1, 5, 0, 0, 0, 0, time.UTC),
				Data:          []byte(`testing 1 2 3`),
				PrevBlockHash: "235c07f768300fd8e42690c69172cec9e2125dced3cb26148b78ccf0f2ccc805",
				Hash:          "45364154e052ac703823f79a9ee9b337cdd6e28b99c4c482868cc0bbc061e684",
			},
			Validator: Validator{
				UUID:  firstValidatorUUID,
				Stake: 0,
			},
		},
		{
			name: "genesis +2 block hash calculation",
			fields: fields{
				Timestamp:     time.Date(1911, 1, 5, 0, 0, 0, 0, time.UTC),
				Data:          []byte(`testing 1 2 3`),
				PrevBlockHash: "235c07f768300fd8e42690c69172cec9e2125dced3cb26148b78ccf0f2ccc805",
				Hash:          "78de985367689eed626990f3be52b99aac235423820a3a8c7ce1c408c7952a85",
			},
			Validator: Validator{
				UUID:  secondValidatorUUID,
				Stake: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				Timestamp:     tt.fields.Timestamp,
				Data:          tt.fields.Data,
				PrevBlockHash: tt.fields.PrevBlockHash,
				Hash:          tt.fields.Hash,
			}
			got := b.calculateHash(tt.Validator)
			assert.Equal(t, tt.fields.Hash, got)
		})
	}
}

func Test_createGenesis(t *testing.T) {
	tests := []struct {
		name string
		want Block
		mock time.Time
	}{
		{
			name: "successfully created genesis block",
			want: Block{
				Timestamp:     time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
				Data:          []byte("each one teach one"),
				PrevBlockHash: "",
				Hash:          "f1dd07239268617d5d30c105f48b965a6256d5d53d61b0770b55ef5f1037eabb",
				Nonce:         0,
			},
			mock: time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now = func() time.Time { return tt.mock }
			got := createGenesis()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBlockChain_AddBlock(t *testing.T) {
	now = func() time.Time { return time.Date(1983, 10, 25, 0, 0, 0, 0, time.UTC) }
	genesisBlock := createGenesis()

	tests := []struct {
		name      string
		Blocks    []Block
		data      string
		validator Validator
	}{
		{
			name: "add first block",
			Blocks: []Block{
				genesisBlock,
			},
			data: "this is the second block",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := &BlockChain{
				Blocks: tt.Blocks,
			}
			bc.AddBlock(tt.data, tt.validator)

			assert.Equal(t, tt.data, string(bc.Blocks[len(bc.Blocks)-1].Data))
			assert.Equal(t, genesisBlock.Hash, string(bc.Blocks[len(bc.Blocks)-1].PrevBlockHash))
		})
	}
}

func TestNewBlockChain(t *testing.T) {
	now = func() time.Time { return time.Date(1983, 10, 25, 0, 0, 0, 0, time.UTC) }

	tests := []struct {
		name string
		want *BlockChain
	}{
		{
			name: "success creating blockchain with genesis block",
			want: &BlockChain{
				Blocks: []Block{
					{
						Timestamp:     time.Date(1983, 10, 25, 0, 0, 0, 0, time.UTC),
						Data:          []byte("each one teach one"),
						PrevBlockHash: "",
						Hash:          "4bd8b9906c934b21dd866f02c0a9238f9c851f9624e7f18b47c19c5c5304fb5c",
						Nonce:         0,
						Validators:    nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBlockChain()
			assert.Equal(t, tt.want, got)
		})
	}
}
