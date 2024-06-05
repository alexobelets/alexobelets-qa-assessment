package types

import "math/big"

type ContractGetterSetterInformation struct {
	ContractAddress string   `json:"contractAddress"`
	DeployerAddress string   `json:"deployerAddress"`
	UintValue       *big.Int `json:",omitempty"`
	Byte32Value     [32]byte `json:",omitempty"`
	BytesValue      []byte   `json:",omitempty"`
}
