package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"github.com/node-real/megafuel-go-sdk/pkg/paymasterclient"
	"github.com/node-real/megafuel-go-sdk/pkg/sponsorclient"
)

var (
	PaymasterURL string
	SponsorURL   string
	ChainID      string

	PolicyUUID        uuid.UUID
	PrivatePolicyUUID uuid.UUID

	TokenContractAddress     common.Address
	WithdrawRecipientAddress common.Address
	HotwalletPrivateKey      string
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PaymasterURL = os.Getenv("PAYMASTER_URL")
	SponsorURL = os.Getenv("SPONSOR_URL")
	ChainID = os.Getenv("CHAIN_ID")

	PolicyUUID, err = uuid.FromString(os.Getenv("POLICY_UUID"))
	if err != nil {
		log.Fatalf("Error parsing POLICY_UUID")
	}
	PrivatePolicyUUID, err = uuid.FromString(os.Getenv("PRIVATE_POLICY_UUID"))
	if err != nil {
		log.Fatalf("Error parsing POLICY_UUID")
	}

	TokenContractAddress = common.HexToAddress(os.Getenv("TOKEN_CONTRACT_ADDRESS"))
	WithdrawRecipientAddress = common.HexToAddress(os.Getenv("WITHDRAW_RECIPIENT_ADDRESS"))
	HotwalletPrivateKey = os.Getenv("HOTWALLET_PRIVATE_KEY")
}

func main() {
	sponsorSetUpPolicyRules()
	cexDoGaslessWithdrawl()
	// wait for nonce to get updated
	time.Sleep(8 * time.Second)
	cexDoPrivatePolicyGaslessWithdrawl()
}

func sponsorSetUpPolicyRules() {
	sponsorClient, err := sponsorclient.New(context.Background(), SponsorURL)
	if err != nil {
		log.Fatal(err)
	}

	// sponsor the tx that interact with the stable coin ERC20 contract
	success, err := sponsorClient.AddToWhitelist(context.Background(), sponsorclient.WhiteListArgs{
		PolicyUUID:    PolicyUUID,
		WhitelistType: sponsorclient.ToAccountWhitelist,
		Values:        []string{TokenContractAddress.String()},
	})
	if err != nil || !success {
		log.Fatal("failed to add token contract whitelist", err)
	}

	// sponsor the tx that from hotwallets
	fromAddress := getAddressFromPrivateKey(HotwalletPrivateKey)

	success, err = sponsorClient.AddToWhitelist(context.Background(), sponsorclient.WhiteListArgs{
		PolicyUUID:    PolicyUUID,
		WhitelistType: sponsorclient.FromAccountWhitelist,
		Values:        []string{fromAddress.String()},
	})
	if err != nil || !success {
		log.Fatal("failed to add contract method whitelist")
	}
}

func cexDoGaslessWithdrawl() {
	withdrawAmount := big.NewInt(1e17)

	// Create a PaymasterClient (for transaction sending)
	paymasterClient, err := paymasterclient.New(context.Background(), PaymasterURL)
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(HotwalletPrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	fromAddress := getAddressFromPrivateKey(HotwalletPrivateKey)

	// Create ERC20 transfer data
	data, err := createERC20TransferData(WithdrawRecipientAddress, withdrawAmount)
	if err != nil {
		log.Fatalf("Failed to create ERC20 transfer data: %v", err)
	}

	// Get the latest nonce for the from address
	nonce, err := paymasterClient.GetTransactionCount(context.Background(), fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create the transaction
	gasPrice := big.NewInt(0)
	tx := types.NewTransaction(nonce, TokenContractAddress, big.NewInt(0), 300000, gasPrice, data)

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

	txInput, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}

	// Convert to Transaction struct for IsSponsorable check
	gasLimit := tx.Gas()
	sponsorableTx := paymasterclient.TransactionArgs{
		To:    &TokenContractAddress,
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
		_, err := paymasterClient.SendRawTransaction(context.Background(), txInput, &paymasterclient.TransactionOptions{UserAgent: "MegaFuel/v1.2.2"})
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

func cexDoPrivatePolicyGaslessWithdrawl() {
	withdrawAmount := big.NewInt(1e17)

	// Create a PaymasterClient (for transaction sending)
	url := fmt.Sprintf("%s/%s", SponsorURL, ChainID)
	privatePaymasterClient, err := paymasterclient.NewPrivatePaymaster(context.Background(), url, PrivatePolicyUUID.String())
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(HotwalletPrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	fromAddress := getAddressFromPrivateKey(HotwalletPrivateKey)

	// Create ERC20 transfer data
	data, err := createERC20TransferData(WithdrawRecipientAddress, withdrawAmount)
	if err != nil {
		log.Fatalf("Failed to create ERC20 transfer data: %v", err)
	}

	// Get the latest nonce for the from address
	nonce, err := privatePaymasterClient.GetTransactionCount(context.Background(), fromAddress, rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber))
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create the transaction
	gasPrice := big.NewInt(0)
	tx := types.NewTransaction(nonce, TokenContractAddress, big.NewInt(0), 300000, gasPrice, data)

	// Get the chain ID
	chainID, err := privatePaymasterClient.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	txInput, err := signedTx.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to marshal transaction: %v", err)
	}

	// Convert to Transaction struct for IsSponsorable check
	gasLimit := tx.Gas()
	sponsorableTx := paymasterclient.TransactionArgs{
		To:    &TokenContractAddress,
		From:  fromAddress,
		Value: (*hexutil.Big)(big.NewInt(0)),
		Gas:   (*hexutil.Uint64)(&gasLimit),
		Data:  (*hexutil.Bytes)(&data),
	}

	// Check if the transaction is sponsorable
	sponsorableInfo, err := privatePaymasterClient.IsSponsorable(context.Background(), sponsorableTx)
	if err != nil {
		log.Fatalf("Error checking sponsorable status: %v", err)
	}

	jsonInfo, _ := json.MarshalIndent(sponsorableInfo, "", "  ")
	fmt.Printf("Sponsorable Information:\n%s\n", string(jsonInfo))

	if sponsorableInfo.Sponsorable {
		// Send the transaction using PaymasterClient
		_, err = privatePaymasterClient.SendRawTransaction(context.Background(), txInput, &paymasterclient.TransactionOptions{UserAgent: "MegaFuel/v1.2.2"})
		if err != nil {
			log.Fatalf("Failed to send sponsorable transaction: %v", err)
		}
		fmt.Printf("Sponsorable transaction sent: %s\n", signedTx.Hash())
	} else {
		fmt.Println("Transaction is not sponsorable. You may need to send it as a regular transaction.")
	}
}

func createERC20TransferData(to common.Address, amount *big.Int) ([]byte, error) {
	transferFnSignature := []byte("transfer(address,uint256)")
	methodID := crypto.Keccak256(transferFnSignature)[:4]
	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return data, nil
}
