package addrdec

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
)

const (
	Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

var (
	NSG_mainnetAddressP2PKH         = addressEncoder.AddressType{EncodeType: "base58", Alphabet: Alphabet, ChecksumType: "doubleSHA256", HashType: "ripemd160", HashLen: 20, Prefix: nil, Suffix: nil}
	NSG_testnetAddressP2PKH         = addressEncoder.AddressType{EncodeType: "base58", Alphabet: Alphabet, ChecksumType: "doubleSHA256", HashType: "ripemd160", HashLen: 20, Prefix: nil, Suffix: nil}
	NSG_mainnetPrivateWIFCompressed = addressEncoder.AddressType{EncodeType: "base58", Alphabet: Alphabet, ChecksumType: "doubleSHA256", HashType: "", HashLen: 32, Prefix: []byte{}, Suffix: nil}
	NSG_testnetPrivateWIFCompressed = addressEncoder.AddressType{EncodeType: "base58", Alphabet: Alphabet, ChecksumType: "doubleSHA256", HashType: "", HashLen: 32, Prefix: []byte{}, Suffix: nil}

	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	IsTestNet bool
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {

	cfg := NSG_mainnetAddressP2PKH
	if dec.IsTestNet {
		cfg = NSG_testnetAddressP2PKH
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			if at, ok := opt.(addressEncoder.AddressType); ok {
				cfg = at
			}
		}
	}

	data := owcrypt.Hash(hash, 0, owcrypt.HASH_ALG_SHA256)
	address := addressEncoder.AddressEncode(data, cfg)
	return "N" + address, nil
}
