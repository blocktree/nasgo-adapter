package rpc

type Client struct {
	baseAddress string
	Wallet      *Wallet
	Tx          *Tx
	Block       *Block
	bk          *BaseClient
}

func NewClient(baseAddress string) *Client {
	bk := newBaseClient(baseAddress)
	return &Client{
		baseAddress: baseAddress,
		bk:          newBaseClient(baseAddress),
		Wallet:      newWalletClient(bk),
		Tx:          newTxClient(bk),
		Block:       newBlockClient(bk),
	}
}
