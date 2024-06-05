package dto

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// EthereumDTO is the struct that holds the data
type EthereumDTO struct {
	Bytes32 [32]byte
	Uint256 *big.Int
	Bytes   []byte
}

// EthereumDTOBuilder is the builder struct
type EthereumDTOBuilder struct {
	bytes32 [32]byte
	uint256 *big.Int
	bytes   []byte
}

// NewEthereumDTOBuilder returns a new instance of the builder
func NewEthereumDTOBuilder() *EthereumDTOBuilder {
	return &EthereumDTOBuilder{}
}

// SetBytes32 sets the Bytes32 field from a hex string
func (b *EthereumDTOBuilder) SetBytes32(value string) *EthereumDTOBuilder {
	byteArray := []byte(value)
	b.bytes32 = common.BytesToHash(byteArray)
	return b
}

// SetUint256 sets the Uint256 field
func (b *EthereumDTOBuilder) SetUint256(uint256 *big.Int) *EthereumDTOBuilder {
	b.uint256 = uint256
	return b
}

// SetBytes sets the Bytes field
func (b *EthereumDTOBuilder) SetBytes(bytes []byte) *EthereumDTOBuilder {
	b.bytes = bytes
	return b
}

// Build constructs the EthereumDTO
func (b *EthereumDTOBuilder) Build() (EthereumDTO, error) {
	return EthereumDTO{
		Bytes32: b.bytes32,
		Uint256: b.uint256,
		Bytes:   b.bytes,
	}, nil
}
