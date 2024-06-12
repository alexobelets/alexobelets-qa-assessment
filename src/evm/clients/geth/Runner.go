package geth

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"main/src/config"
	"main/src/contracts/getter_setter"
	"main/src/evm/clients/geth/account"
	"main/src/evm/clients/geth/client"
	"main/src/evm/clients/geth/dto"
	"main/src/evm/clients/geth/transactions"
	"main/src/evm/clients/geth/types"
	"main/src/utils"
	"math/big"
)

// Main values used by the application
var (
	tomlConfig           = utils.GetConfig()
	rpcURL               = tomlConfig.RPC.Url
	privateKey           = tomlConfig.Account.Key
	contractAddress      = tomlConfig.Contract.Address
	mode                 = tomlConfig.Contract.Mode
	values               = tomlConfig.Contract.Values
	gasLimit             = tomlConfig.Client.GasLimit
	timeout              = tomlConfig.Client.WaitingTimeout
	ethClient            = client.ConnectClient(rpcURL)
	deployerAddress      = account.GetDeployerAddressFromPrivateKey(privateKey)
	getterSetterContract *getter_setter.GetterSetter
)

func Run() {
	// Validate if the account is sufficiently funded
	account.ValidateBalanceFunded(deployerAddress, ethClient)

	switch mode {
	case utils.DEPLOY_MODE:
		auth := GetSigner(ethClient, deployerAddress, privateKey)
		deployedContractAddress := client.DeployContract(auth, ethClient, timeout)
		contractAddress := deployedContractAddress.Hex()
		getterSetterContract = client.AttachToContract(contractAddress, ethClient)

	case utils.CALL_MODE:
		// Attach to the contract
		getterSetterContract = client.AttachToContract(contractAddress, ethClient)
		ExecuteSetterGetterContractFunction(values, deployerAddress, privateKey, getterSetterContract, ethClient, timeout)

	case utils.DEMO_MODE: // executes full e2e scenario "deploy-contract" + "call-contract" + "read-only-contract"
		auth := GetSigner(ethClient, deployerAddress, privateKey)
		deployedContractAddress := client.DeployContract(auth, ethClient, timeout)
		contractAddress := deployedContractAddress.Hex()
		getterSetterContract = client.AttachToContract(contractAddress, ethClient)
		ExecuteSetterGetterContractFunction(values, deployerAddress, privateKey, getterSetterContract, ethClient, timeout)

	case utils.READ_ONLY_MODE:
		log.Printf("Reading the contract: %s, owned by: %s", contractAddress, deployerAddress)
		getterSetterContract = client.AttachToContract(contractAddress, ethClient)
		// N/B: There is no need to add reading logic here as it takes place after the `switch-case` execution.
		// If the common logic after `switch-case` changes, update this case body too.
	}
	// Write the contract information to a JSON file
	output := ReadGetterSetterContract(getterSetterContract, contractAddress, deployerAddress)
	pathToJson := "output/contractOutputInformation.json" // paste your desired path
	utils.JsonWriter(output, pathToJson)
}

// !!! Below are the steps/actions for runner

// ExecuteSetterGetterContractFunction executes functions to set values in the contract based on the provided input values.
// If a value is not provided, the corresponding setter function will be gracefully skipped.
//
// Parameters
// - values: configuration values for contract interaction (config.Values)
// - deployerAddress: the address of the deployer (common.Address)
// - privateKey: the private key for authentication (string)
// - getterSetterContract: the contract instance for setting and getting values (*getter_setter.GetterSetter)
// - ethClient: the Ethereum client for interaction (*ethclient.Client)
// - timeout: the timeout duration for the contract interaction (int)
func ExecuteSetterGetterContractFunction(values config.Values, deployerAddress common.Address, privateKey string, getterSetterContract *getter_setter.GetterSetter, ethClient *ethclient.Client, timeout int) {
	getterSetterDto := SetGetterSetterDTO(values)

	uint256Value := values.Uint256
	if uint256Value != nil && uint256Value.Cmp(big.NewInt(0)) >= 0 {
		auth := GetSigner(ethClient, deployerAddress, privateKey)
		SetUintInGetterSetterContract(getterSetterContract, auth, getterSetterDto, ethClient, timeout)
	}
	bytes32Value := values.Bytes32
	if bytes32Value != "" {
		auth := GetSigner(ethClient, deployerAddress, privateKey)
		SetBytes32InGetterSetterContract(getterSetterContract, auth, getterSetterDto, ethClient, timeout)
	}
	bytesValue := values.Bytes
	if bytesValue != "" {
		auth := GetSigner(ethClient, deployerAddress, privateKey)
		SetBytesInGetterSetterContract(getterSetterContract, auth, getterSetterDto, ethClient, timeout)
	}
}

// ReadGetterSetterContract retrieves values from the contract, using its getters.
//
// Parameters:
// - getterSetterContract: the contract instance for getting values (*getter_setter.GetterSetter)
// - contractAddress: the address of the contract (string)
// - deployerAddress: the address of the deployer (common.Address)
// Return type:
// - types.ContractGetterSetterInformation
func ReadGetterSetterContract(getterSetterContract *getter_setter.GetterSetter, contractAddress string, deployerAddress common.Address) types.ContractGetterSetterInformation {
	uintResponse, err := getterSetterContract.GetUint256(&bind.CallOpts{})
	if err != nil {
		log.Fatal("Uint256 value could not be fetched from the contract", err)
	}
	log.Println("Uint value is: ", uintResponse)

	bytes32Response, err := getterSetterContract.GetBytes32(&bind.CallOpts{})
	if err != nil {
		log.Fatal("Bytes32 value could not be fetched from the contract", err)
	}
	log.Println("Bytes32 value is:", bytes32Response)

	bytesResponse, err := getterSetterContract.GetBytes(&bind.CallOpts{})
	if err != nil {
		log.Fatal("Bytes value could not be fetched from the contract", err)
	}
	log.Println("Bytes value is: "+string(bytesResponse)+" -> Bytes raw value:", bytesResponse)

	output := types.ContractGetterSetterInformation{
		ContractAddress: contractAddress,
		DeployerAddress: deployerAddress.String(),
		UintValue:       uintResponse,
		Byte32Value:     bytes32Response,
		BytesValue:      bytesResponse,
	}
	return output
}

func SetUintInGetterSetterContract(getterSetterContract *getter_setter.GetterSetter, auth *bind.TransactOpts, getterSetterDto dto.EthereumDTO, ethClient *ethclient.Client, timeout int) {
	transaction, err := getterSetterContract.SetUint256(auth, getterSetterDto.Uint256)
	if err != nil {
		log.Fatalf("Failed to set uint256: %v", err)
	}

	log.Printf("Waiting for transaction for setUint256: %s\n", transaction.Hash().Hex())
	transactions.WaitMined(ethClient, transaction, timeout)
}

func SetBytes32InGetterSetterContract(getterSetterContract *getter_setter.GetterSetter, auth *bind.TransactOpts, getterSetterDto dto.EthereumDTO, ethClient *ethclient.Client, timeout int) {
	transaction, err := getterSetterContract.SetBytes32(auth, getterSetterDto.Bytes32)
	if err != nil {
		log.Fatalf("Failed to set bytes32: %v", err)
	}
	log.Printf("Waiting for transaction for setBytes32: %s\n", transaction.Hash().Hex())
	transactions.WaitMined(ethClient, transaction, timeout)
}

func SetBytesInGetterSetterContract(getterSetterContract *getter_setter.GetterSetter, auth *bind.TransactOpts, getterSetterDto dto.EthereumDTO, ethClient *ethclient.Client, timeout int) {
	transaction, err := getterSetterContract.SetBytes(auth, getterSetterDto.Bytes)
	if err != nil {
		log.Fatalf("Failed to set bytes: %v", err)
	}
	log.Printf("Waiting for transaction for setBytes: %s\n", transaction.Hash().Hex())
	transactions.WaitMined(ethClient, transaction, timeout)
}

// SetGetterSetterDTO creates a new DTO based on the provided values.
//
// Parameters:
// - values: the values to be set in the contract (config.Values)
// Returns:
// - DTO
func SetGetterSetterDTO(values config.Values) dto.EthereumDTO {
	log.Printf("Creating DTO with values: Uint256: '%d', Bytes32: '%s', Bytes: '%s'", values.Uint256, values.Bytes32, values.Bytes)
	getterSetterDto, err := dto.NewEthereumDTOBuilder().
		SetUint256(values.Uint256).
		SetBytes32(values.Bytes32).
		SetBytes([]byte(values.Bytes)).
		Build()
	if err != nil {
		log.Fatalf("Failed to create DTO: %v", err)
	}
	log.Println("DTO successfully created:", getterSetterDto)
	return getterSetterDto
}

// GetSigner retrieves latest account information and sets the signer options for a transaction,
// including the chain ID, the nonce, gas price, and gas limit.
//
// Parameters:
// - ethClient: the Ethereum client (*ethclient.Client)
// - deployerAddress: the address of the deployer (common.Address)
// - privateKey: the private key for authentication (string)
// Returns:
// - TransactOpts
func GetSigner(ethClient *ethclient.Client, deployerAddress common.Address, privateKey string) *bind.TransactOpts {
	log.Println("Getting signer options...")
	chainID := transactions.GetChainId(ethClient)
	nonce := transactions.GetNonce(ethClient, deployerAddress)
	gasPrice := transactions.GetTransactionGasPrice(ethClient)
	privateKeyECDSA := account.PrivateToECDSA(privateKey)

	auth := client.GetTransactor(privateKeyECDSA, chainID)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	log.Println("Signer options successfully obtained")
	return auth
}
