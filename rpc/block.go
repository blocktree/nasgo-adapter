package rpc

import (
	"encoding/json"
	"strconv"

	"github.com/go-errors/errors"
	"gopkg.in/resty.v1"
)

type Block struct {
	bk *BaseClient
}

func newBlockClient(bk *BaseClient) *Block {
	return &Block{
		bk: bk,
	}
}

type BlockResponse struct {
	Success bool    `json:"success"`
	Block   *Header `json:"block"` // block
}

type Header struct {
	ID                   string `json:"id"`                   // Hash
	Version              uint32 `json:"version"`              // Block version information (note, this is signed)
	Timestamp            int64  `json:"timestamp"`            // A timestamp recording when this block was created (Will overflow in 2106[2])
	Height               uint64 `json:"height"`               // Height
	PrevBlock            string `json:"previousBlock"`        // The hash value of the previous block this particular block references
	NumberOfTransactions uint32 `json:"numberOfTransactions"` // Number Of Transactions
}

type BlockHeightResponse struct {
	Success bool   `json:"success"`
	Height  uint64 `json:"height"`
}

// GetBlockHeight get height
func (blk *Block) GetBlockHeight() (uint64, error) {
	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/blocks/getHeight")
	if err != nil {
		return 0, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return 0, err
	}
	bhResp := BlockHeightResponse{}
	if err := json.Unmarshal(body, &bhResp); err != nil {
		return 0, errors.New(err)
	}
	return bhResp.Height, nil
}

// GetByHash by hash
func (blk *Block) GetByHash(hash string) (*Header, error) {
	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/blocks/get?hash=" + hash)
	if err != nil {
		return nil, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	blockResponse := BlockResponse{}
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		return nil, errors.New(err)
	}
	return blockResponse.Block, nil
}

// GetByHeight by height
func (blk *Block) GetByHeight(height uint64) (*Header, error) {
	h := strconv.FormatInt(int64(height), 10)
	resp, err := resty.
		R().
		Get(blk.bk.baseAddress + "/api/blocks/get?height=" + h)
	if err != nil {
		return nil, err
	}
	body, err := blk.bk.ReadResponse(resp)
	if err != nil {
		return nil, err
	}
	blockResponse := BlockResponse{}
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		return nil, errors.New(err)
	}
	return blockResponse.Block, nil
}
