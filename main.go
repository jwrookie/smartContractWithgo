package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth_math "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jwrookie/smartContractWithgo/api"
	"math/big"
)

func getAccountAuth(client *ethclient.Client, private string) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//fetch the last use nonce of account
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(1000000)

	return auth
}

func main() {
	// address of etherum env
	client, err := ethclient.Dial("http://127.0.0.1:7545")
	if err != nil {
		panic(err)
	}

	// create auth and transaction package for deploying smart contract
	auth := getAccountAuth(client, "24b69cadade23f743c22eaf09fe732da5bbf1c6dd2b1df85a2cc701331470b42")

	//deploying smart contract
	deployedContractAddress, tx, instance, err := api.DeployApi(auth, client) //api is redirected from api directory from our contract go file
	if err != nil {
		panic(err)
	}

	fmt.Println(deployedContractAddress.Hex())
	fmt.Println(tx)
	fmt.Println(instance)

	// connect to smart contract
	conn, err := api.NewApi(common.HexToAddress(deployedContractAddress.Hex()), client)
	if err != nil {
		panic(err)
	}

	fmt.Println("--------------------")
	Admin(conn)
	Balance(conn)
	Deposite(client, conn, new(big.Int).Mul(eth_math.BigPow(10, 18), big.NewInt(2)))
	Withdrawl(client, conn, eth_math.BigPow(10, 18))
	Balance(conn)
	fmt.Println("++++++++++++++++++++")
}

func Admin(conn *api.Api) {
	reply, err := conn.Admin(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	s, _ := json.Marshal(reply)
	fmt.Println(string(s))
}

func Balance(conn *api.Api) {
	reply, err := conn.Balance(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}

	s, _ := json.Marshal(reply)
	fmt.Println(string(s))
}

func Deposite(client *ethclient.Client, conn *api.Api, amt *big.Int) {
	// tx opts 不能复用，因为 nonce 必须正确，每发生一笔交易，nonce都会加一。
	auth := getAccountAuth(client, "24b69cadade23f743c22eaf09fe732da5bbf1c6dd2b1df85a2cc701331470b42")
	reply, err := conn.Deposite(auth, amt)
	if err != nil {
		panic(err)
	}

	s, _ := json.Marshal(reply)
	fmt.Println(string(s))
}

func Withdrawl(client *ethclient.Client, conn *api.Api, amt *big.Int) {
	// tx opts 不能复用，因为 nonce 必须正确，每发生一笔交易，nonce都会加一。
	auth := getAccountAuth(client, "24b69cadade23f743c22eaf09fe732da5bbf1c6dd2b1df85a2cc701331470b42")

	reply, err := conn.Withdrawl(auth, amt)
	if err != nil {
		panic(err)
	}

	s, _ := json.Marshal(reply)
	fmt.Println(string(s))
}
