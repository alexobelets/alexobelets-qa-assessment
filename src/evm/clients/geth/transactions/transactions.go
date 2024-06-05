package transactions

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func GetChainId(client *ethclient.Client) *big.Int {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}
	log.Println("Chain ID:", chainID)
	return chainID
}

func GetNonce(client *ethclient.Client, deployerAddress common.Address) uint64 {
	nonce, err := client.PendingNonceAt(context.Background(), deployerAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Nonce:", nonce)
	return nonce
}

func GetTransactionGasPrice(client *ethclient.Client) *big.Int {
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Gas price:", gasPrice)
	return gasPrice
}

func WaitMined(client *ethclient.Client, transaction *types.Transaction, timeout int) *types.Receipt {
	timeToWait := SetTimeToWait(timeout)
	backgroundContext, cancel := context.WithTimeout(context.Background(), timeToWait)
	defer cancel()

	receipt, err := bind.WaitMined(backgroundContext, client, transaction)
	if err != nil {
		if errors.Is(backgroundContext.Err(), context.DeadlineExceeded) {
			log.Fatalf("Cancelled waiting for transaction status by timeout: %v. Verify your transaction manually or increase the timeout.", timeToWait)
		} else {
			log.Fatalf("Failed to wait and get receipt from transaction: %v", err)
		}
	}

	log.Println("Transaction mined in block", receipt.BlockNumber)
	return receipt
}

func WaitDeployed(err error, client *ethclient.Client, transaction *types.Transaction, timeout int) common.Address {
	timeToWait := SetTimeToWait(timeout)
	backgroundContext, cancel := context.WithTimeout(context.Background(), timeToWait)
	defer cancel()

	deployed, err := bind.WaitDeployed(backgroundContext, client, transaction)
	if err != nil {
		if errors.Is(backgroundContext.Err(), context.DeadlineExceeded) {
			log.Fatalf("Cancelled waiting for transaction status by timeout: %s. Verify your transaction manually or increase the timeout.", timeToWait)
		} else {
			log.Fatal("Failed to wait and get receipt from transaction:", err)
		}
	}
	return deployed
}

func SetTimeToWait(timeout int) time.Duration {
	defaultTimeout := 300
	defaultTimeToWait := time.Duration(defaultTimeout) * time.Second
	timeToWait := time.Duration(timeout) * time.Second

	// if timeout is less than default value, use default timeout
	if timeout < defaultTimeout {
		log.Println("Waiting for transaction status update by default timeout:", defaultTimeToWait)
		timeToWait = defaultTimeToWait
	} else {
		log.Println("Waiting for transaction status update by custom timeout: ", timeToWait)
	}
	return timeToWait
}
