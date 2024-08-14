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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const YourPrivateKey = ""
const TokenContractAddress = "0x.."
const RecipientAddress = "0x.."

type PaymasterClient struct {
	*ethclient.Client
	rpcClient *rpc.Client
}

type SponsorableInfo struct {
	Sponsorable    bool   `json:"sponsorable"`
	SponsorName    string `json:"sponsorName"`
	SponsorIcon    string `json:"sponsorIcon"`
	SponsorWebsite string `json:"sponsorWebsite"`
}

type Transaction struct {
	To    *common.Address `json:"to"`
	From  common.Address  `json:"from"`
	Value *hexutil.Big    `json:"value"`
	Gas   *hexutil.Uint64 `json:"gas"`
	Data  *hexutil.Bytes  `json:"data"`
}

func NewPaymasterClient(url string) (*PaymasterClient, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}

	ethClient := ethclient.NewClient(rpcClient)

	return &PaymasterClient{
		Client:    ethClient,
		rpcClient: rpcClient,
	}, nil
}

func (c *PaymasterClient) IsSponsorable(ctx context.Context, tx Transaction) (*SponsorableInfo, error) {
	var result SponsorableInfo
	err := c.rpcClient.CallContext(ctx, &result, "pm_isSponsorable", tx)
	if err != nil {
		return nil, err
	}
	return &result, nil
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

	// Connect to an Ethereum node (for transaction assembly)
	client, err := ethclient.Dial("https://bsc-testnet-dataseed.bnbchain.org")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}
	// Create a PaymasterClient (for transaction sending)
	paymasterClient, err := NewPaymasterClient("https://bsc-megafuel-testnet.nodereal.io")
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(YourPrivateKey)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Token contract address
	tokenAddress := common.HexToAddress(TokenContractAddress)

	// Recipient address
	toAddress := common.HexToAddress(RecipientAddress)

	// Amount of tokens to transfer (adjust based on token decimals)
	amount := big.NewInt(1000000000000000000) // 1 token for a token with 18 decimals

	// Create ERC20 transfer data
	data, err := createERC20TransferData(toAddress, amount)
	if err != nil {
		log.Fatalf("Failed to create ERC20 transfer data: %v", err)
	}

	// Get the latest nonce for the from address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Create the transaction
	gasPrice := big.NewInt(0)
	tx := types.NewTransaction(nonce, tokenAddress, big.NewInt(0), 300000, gasPrice, data)

	// Get the chain ID
	chainID, err := client.ChainID(context.Background())
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
