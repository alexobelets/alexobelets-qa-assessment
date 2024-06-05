package utils

import (
	"github.com/BurntSushi/toml"
	"log"
	"main/src/config"
)

const (
	DEMO_MODE      = "demo"
	DEPLOY_MODE    = "deploy-contract"
	CALL_MODE      = "call-contract"
	READ_ONLY_MODE = "read-only-contract"
)

func GetConfig() config.Config {
	var fetchedConfig config.Config

	if _, err := toml.DecodeFile("config.toml", &fetchedConfig); err != nil {
		log.Fatalf("Failed to parse fetchedConfig.toml: %v", err)
	}

	ValidateConfig(fetchedConfig)

	return fetchedConfig
}

func ValidateConfig(config config.Config) {
	if config.RPC.Url == "" {
		log.Fatal("config.toml: RPC.URL is required")
	}
	if config.Account.Key == "" {
		log.Fatal("config.toml: Account.key is required")
	}
	if config.Client.GasLimit <= 0 {
		defaultGasLimit := 3000000
		log.Printf("config.toml: Client.gasLimit is 0 or less, defaulting to '%v'", defaultGasLimit)
		config.Client.GasLimit = uint64(defaultGasLimit)
	}
	if !isValidMode(config.Contract.Mode) {
		log.Fatal("config.toml: Contract.mode is required, acceptable values: ", DEMO_MODE, DEPLOY_MODE, CALL_MODE, READ_ONLY_MODE)
	}
	log.Printf("config.toml: Script will run in mode: '%s'", config.Contract.Mode)
	if config.Contract.Mode == "read-only-contract" {
		hasContractAddress(config)
	}
	if config.Contract.Mode == "call-contract" || config.Contract.Mode == "demo" {
		hasContractAddress(config)
		hasValuesToSet(config)
	}
	if config.Contract.Address == "" { // optional
		log.Print("config.toml: Contract.address is not provided. Either deploy a new or use an existing contract.")
	}
}

func hasContractAddress(config config.Config) {
	if config.Contract.Address == "" {
		log.Fatal("config.toml: Contract.address is required to be set for the Contract.Mode: ", config.Contract.Mode)
	}
}

func isValidMode(mode string) bool {
	// allowedModes contains the list of valid values'
	var allowedModes = []string{DEMO_MODE, DEPLOY_MODE, CALL_MODE, READ_ONLY_MODE}

	// Check if mode is in allowedModes list
	for _, allowedMode := range allowedModes {
		if allowedMode == mode {
			return true
		}
	}
	return false
}

func hasValuesToSet(config config.Config) bool {
	// Check if at least one of the values to set in the contracts
	// otherwise, it is pointless to call the existing contract without setting any of the values
	if config.Contract.Values.Uint256 == nil && config.Contract.Values.Bytes32 == "" && config.Contract.Values.Bytes == "" {
		log.Fatal("config.toml: Ensure you have all or one of the possible values to set in the contract as follows: Contract.Uint256, Contract.Bytes32, Contract.Bytes")
	}
	return true
}
