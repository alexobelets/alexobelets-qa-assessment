package account

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"main/src/utils"
	"math/big"
)

func GetDeployerAddressFromPrivateKey(privateKey string) common.Address {
	privateKeyECDSA := PrivateToECDSA(privateKey)
	deployerAddress := GetOwnerAddress(GetEOAPublicKey(privateKeyECDSA))
	return deployerAddress
}

func GetEOAPublicKey(privateKeyECDSA *ecdsa.PrivateKey) *ecdsa.PublicKey {
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	log.Println("Public key obtained")
	return publicKeyECDSA
}

func PrivateToECDSA(privateKey string) *ecdsa.PrivateKey {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal("Error converting private key to ECDSA: ", err)
	}
	return privateKeyECDSA
}

func GetOwnerAddress(publicKeyECDSA *ecdsa.PublicKey) common.Address {
	deployerAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("Owner address:", deployerAddress)
	return deployerAddress
}

func ValidateBalanceFunded(deployerAddress common.Address, client *ethclient.Client) {
	account := common.HexToAddress(deployerAddress.String())
	balanceInWei, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal("Could not get balance:", err)
	}

	if isUnfunded(balanceInWei) {
		log.Fatalf("Balance of the account %s is 0. Fund the account and try again.", deployerAddress)
	}

	log.Println("Current balance in ETH:", utils.WeiToEther(balanceInWei))
}

func isUnfunded(balance *big.Int) bool {
	return balance.Cmp(big.NewInt(0)) == 0
}
