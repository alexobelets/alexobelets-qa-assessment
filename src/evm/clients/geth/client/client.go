package client

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"main/src/contracts/getter_setter"
	"main/src/evm/clients/geth/transactions"
	"math/big"
)

func ConnectClient(rpcURL string) *ethclient.Client {
	log.Println("Connecting to the client...")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal("Could not connect to the client:", err)
	}
	log.Println("Client connected")
	return client
}

func DeployContract(auth *bind.TransactOpts, client *ethclient.Client, timeout int) common.Address {
	address, transaction, _, err := getter_setter.DeployGetterSetter(auth, client)
	if err != nil {
		log.Fatalf("Failed to deploy new GetterSetter contract: %v", err)
	}

	log.Printf("Waiting for pending contract deployment with transaction hash: 0x%x, for contract address: 0x%x", transaction.Hash(), address)
	deployedContractAddress := transactions.WaitDeployed(client, transaction, timeout)
	log.Printf("Contract address: %s. Write it down for future (re-)usage.", deployedContractAddress)

	return deployedContractAddress
}

// GetTransactor returns the transactor (signer) for the given private key
func GetTransactor(privateKeyECDSA *ecdsa.PrivateKey, chainID *big.Int) *bind.TransactOpts {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	log.Println("Authorized transactor for address", auth.From)
	return auth
}

func AttachToContract(contractAddress string, client *ethclient.Client) *getter_setter.GetterSetter {
	getterSetterContract, err := getter_setter.NewGetterSetter(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatalf("Failed to attach to the contract: %v", err)
	}
	log.Println("Successfully attached to the contract:", contractAddress)
	return getterSetterContract
}
