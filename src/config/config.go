package config

import "math/big"

type Config struct {
	RPC      RPC
	Client   Client
	Account  Account
	Contract Contract
}

type RPC struct {
	Url string
}

type Client struct {
	GasLimit       uint64
	WaitingTimeout int
}

type Account struct {
	Key string
}

// Contract configuration
type Contract struct {
	Mode    string
	Address string
	Values  Values
}

// Values to be set in the contract
type Values struct {
	Uint256 *big.Int
	Bytes32 string
	Bytes   string
}
