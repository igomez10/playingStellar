package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	fmt.Println()
	fmt.Println("-------------START----------------")
	// Generate a new randomly generated address
	keypair0, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Seed 0:", keypair0.Seed())
	log.Println("Address 0:", keypair0.Address())

	// Create and fund the address on TestNet, using friendbot
	client := horizonclient.DefaultTestNetClient
	_, err = client.Fund(keypair0.Address())
	if err != nil {
		fmt.Printf("error funding test account: %s", err)
		return
	}

	//fmt.Printf("Successfully funded %+v : %+v", keypair0.Address(), trSuccess)

	requestForAccount := horizonclient.AccountRequest{AccountID: keypair0.Address()}

	responseAccount, err := client.AccountDetail(requestForAccount)
	if err != nil {
		fmt.Println("unable to retrieve created account", err)
		return
	}

	fmt.Println("account:")

	spew.Dump(responseAccount)

	// Generate a second randomly generated address

	keypair1, err := keypair.Random()
	if err != nil {
		fmt.Println("unable to create second keypair", err)
		return
	}

	fmt.Println("Seed 1:", keypair1.Seed())
	fmt.Println("Address 1:", keypair1.Address())

	// Construct the operation
	moveMoney1To2Operation := txnbuild.CreateAccount{
		Destination: keypair1.Address(),
		Amount:      "10",
	}

	transaction := txnbuild.Transaction{
		SourceAccount: &responseAccount,
		Operations:    []txnbuild.Operation{&moveMoney1To2Operation},
		Timebounds:    txnbuild.NewTimeout(300),
		Network:       network.TestNetworkPassphrase,
	}

	transactionBase64, err := transaction.BuildSignEncode(keypair0)
	if err != nil {
		fmt.Println("error signing transaction to base64", err)
		return
	}

	resp, err := client.SubmitTransactionXDR(transactionBase64)
	if err != nil {
		fmt.Printf("error submitting transaction %+v", err)
	}

	fmt.Println("transaction response")
	spew.Dump(resp)

	fmt.Println("-------------END----------------")
	fmt.Println()
}
