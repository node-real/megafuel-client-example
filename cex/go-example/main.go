package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const TokenContractAddress = "0x.."
const WithdrawRecipientAddress = "0x.."
const SponsorPolicyId = ".."
const HotwalletPrivateKey = ".."

const sponsorAPIEndpoint = "https://open-platform.nodereal.io/{Your_API_key}/megafuel"
const paymasterEndpoint = "https://bsc-megafuel.nodereal.io"

// testnet endpoint
// const sponsorAPIEndpoint = "https://open-platform.nodereal.io/{Your_API_key}/megafuel-testnet"
// const paymasterEndpoint = "https://bsc-megafuel-testnet.nodereal.io'"

func main() {
	sponsorSetUpPolicyRules()
	cexDoGaslessWithdrawl()
}

func sponsorSetUpPolicyRules() {
	sponsorClient, err := NewSponsorClient(sponsorAPIEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	// sponsor the tx that interact with the stable coin ERC20 contract
	success, err := sponsorClient.AddToWhitelist(context.Background(), WhitelistParams{
		PolicyUUID:    SponsorPolicyId,
		WhitelistType: ToAccountWhitelist,
		Values:        []string{TokenContractAddress},
	})
	if err != nil || !success {
		log.Fatal("failed to add token contract whitelist", err)
	}

	// sponsor the tx that from hotwallets
	fromAddress := getAddressFromPrivateKey(HotwalletPrivateKey)

	success, err = sponsorClient.AddToWhitelist(context.Background(), WhitelistParams{
		PolicyUUID:    SponsorPolicyId,
		WhitelistType: FromAccountWhitelist,
		Values:        []string{fromAddress.String()},
	})
	if err != nil || !success {
		log.Fatal("failed to add contract method whitelist")
	}
}

func cexDoGaslessWithdrawl() {
	withdrawAmount := big.NewInt(1e17)
	// Create a PaymasterClient (for transaction sending)
	paymasterClient, err := NewPaymasterClient(paymasterEndpoint)
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(HotwalletPrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	fromAddress := getAddressFromPrivateKey(HotwalletPrivateKey)

	// Token contract address
	tokenAddress := common.HexToAddress(TokenContractAddress)

	// Create ERC20 transfer data
	data, err := createERC20TransferData(common.HexToAddress(WithdrawRecipientAddress), withdrawAmount)
	if err != nil {
		log.Fatalf("Failed to create ERC20 transfer data: %v", err)
	}

	// Get the pending nonce for the from address, strongly suggest to fetch nonce from paymaster endpoint when
	// submitting multiple transactions in rapid succession, to ensure that the nonce are sequential.
	nonce, err := paymasterClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create the transaction
	gasPrice := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), 300000, gasPrice, data)

	// Get the chain ID
	chainID, err := paymasterClient.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Convert to Transaction struct for IsSponsorable check
	gasLimit := tx.Gas()
	sponsorableTx := Transaction{
		To:    &tokenAddress,
		From:  fromAddress,
		Value: (*hexutil.Big)(big.NewInt(0)),
		Gas:   (*hexutil.Uint64)(&gasLimit),
		Data:  (*hexutil.Bytes)(&data),
	}

	// Check if the transaction is sponsorable
	sponsorableInfo, err := paymasterClient.IsSponsorable(context.Background(), sponsorableTx)
	if err != nil {
		log.Fatalf("Error checking sponsorable status: %v", err)
	}

	jsonInfo, _ := json.MarshalIndent(sponsorableInfo, "", "  ")
	fmt.Printf("Sponsorable Information:\n%s\n", string(jsonInfo))

	if sponsorableInfo.Sponsorable {
		// Send the transaction using PaymasterClient
		err := paymasterClient.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatalf("Failed to send sponsorable transaction: %v", err)
		}
		fmt.Printf("Sponsorable transaction sent: %s\n", signedTx.Hash())
	} else {
		fmt.Println("Transaction is not sponsorable. You may need to send it as a regular transaction.")
	}
}

func getAddressFromPrivateKey(pk string) common.Address {
	// Load your private key
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}
