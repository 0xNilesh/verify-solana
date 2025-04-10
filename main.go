package main

import (
	"context"
	"encoding/hex"
	"fmt"

	// "github.com/davecgh/go-spew/spew"
	binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	// Input variables
	txSignature := "5HuzJVLE7PYfJxTs2JYdEK9EarQfH5HhJf9GuJKQ9ngbNRQEXV9KDbrppVzGnLY9QfBg8naoAJVcqJw2cjpkQ3PR" // transaction signature (txHash)
	address := "4HwuvaEVT4qnvb5TkSvMyYPrprqpnJjz8LSL6TPnuJ2U"                                                 // address of the signer

	// Setup
	endpoint := rpc.MainNetBeta_RPC
	client := rpc.New(endpoint)
	sig := solana.MustSignatureFromBase58(txSignature)
	pubKey := solana.MustPublicKeyFromBase58(address)

	// Fetch transaction
	out, err := client.GetTransaction(
		context.TODO(),
		sig,
		&rpc.GetTransactionOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		panic(err)
	}

	// Decode transaction
	decodedTx, err := solana.TransactionFromDecoder(binary.NewBinDecoder(out.Transaction.GetBinary()))
	if err != nil {
		panic(err)
	}

	// Get message bytes
	msgBytes, err := decodedTx.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}

	// Get first signature
	signature := decodedTx.Signatures[0]

	// Output info
	fmt.Println("Signature (Base58):", signature.String())
	fmt.Println("Message bytes (hex):", hex.EncodeToString(msgBytes))
	// spew.Dump(decodedTx.Message)

	// Verify signature against provided pubkey
	isValid := signature.Verify(pubKey, msgBytes)
	fmt.Printf("Was the transaction signed by %s? %v\n", address, isValid)
}
