package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"github.com/node-real/megafuel-go-sdk/pkg/paymasterclient"
)

var (
	PaymasterURL string

	TokenContractAddress common.Address
	RecipientAddress     common.Address
	PrivateKey           string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PaymasterURL = os.Getenv("PAYMASTER_URL")
	TokenContractAddress = common.HexToAddress(os.Getenv("TOKEN_CONTRACT_ADDRESS"))
	RecipientAddress = common.HexToAddress(os.Getenv("RECIPIENT_ADDRESS"))
	PrivateKey = os.Getenv("YOUR_PRIVATE_KEY")
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

func main() {
	walletUserDoGaslessTx()
}

func walletUserDoGaslessTx() {
	// Create a PaymasterClient (for transaction sending)
	paymasterClient, err := paymasterclient.New(context.Background(), PaymasterURL)
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// Amount of tokens to transfer (adjust based on token decimals)
	amount := big.NewInt(1e18) // 1 token for a token with 18 decimals

	// Create ERC20 transfer data
	data, err := createERC20TransferData(RecipientAddress, amount)
	if err != nil {
		log.Fatalf("Failed to create ERC20 transfer data: %v", err)
	}

	blockNumber := rpc.PendingBlockNumber
	nonce, err := paymasterClient.GetTransactionCount(context.Background(), fromAddress, rpc.BlockNumberOrHash{BlockNumber: &blockNumber})
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

	fmt.Printf("Sponsorable Information:\n%+v\n", sponsorableInfo)

	if sponsorableInfo.Sponsorable {
		// We strongly encourage you to set the UserAgent value. It should represent
		// your wallet name or brand name. This information is for further statistical
		// analysis and insight. Setting a unique UserAgent will help MegaFuel to
		// better understand wallet usage patterns and improve service.
		_, err = paymasterClient.SendRawTransaction(context.Background(), txInput, &paymasterclient.TransactionOptions{UserAgent: "myWalletName/v1.0.0"})
		if err != nil {
			log.Fatalf("Failed to send sponsorable transaction: %v", err)
		}
		fmt.Printf("Sponsorable transaction sent: %s\n", signedTx.Hash())
	} else {
		fmt.Println("Transaction is not sponsorable. You may need to send it as a regular transaction.")
	}
}
