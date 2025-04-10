package main

import (
	"context"
	"crypto/ed25519"
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

	// Method 1: Using solana-go verify
	ok1 := verifyUsingSolanaLib(signature, pubKey, msgBytes)
	fmt.Printf("[solana-go] Was the transaction signed by %s? %v\n", address, ok1)

	// Method 2: Using standard crypto/ed25519
	ok2 := verifyUsingEd25519(signature[:], pubKey[:], msgBytes)
	fmt.Printf("[ed25519]   Was the transaction signed by %s? %v\n", address, ok2)
}

// Method using solana-go's built-in Verify
func verifyUsingSolanaLib(sig solana.Signature, pub solana.PublicKey, msg []byte) bool {
	return sig.Verify(pub, msg)
}

// Method using standard crypto/ed25519
func verifyUsingEd25519(sigBytes, pubKeyBytes, msg []byte) bool {
	if len(pubKeyBytes) != ed25519.PublicKeySize || len(sigBytes) != ed25519.SignatureSize {
		return false
	}

	return ed25519.Verify(pubKeyBytes, msg, sigBytes)
}
