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
)

const PaymentTokenContractAddress = "0x.."
const PaymentRecipientAddress = "0x.."
const PaymentSponsorPolicyId = ".."
const SponsorAPIEndpoint = "https://open-platform.nodereal.io/{Your_API_key}/eoa-paymaster-testnet"
const UserPrivateKey = "..."

func main() {
	receiver := common.HexToAddress(PaymentRecipientAddress)
	payAmount := big.NewInt(1e17)

	paymentGatewaySetUpPolicyRules(receiver)

	userDoGaslessPayment(receiver, payAmount)
}

func paymentGatewaySetUpPolicyRules(receiver common.Address) {
	sponsorClient, err := NewSponsorClient(SponsorAPIEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	// sponsor the tx that interact with the stable coin ERC20 contract
	success, err := sponsorClient.AddToWhitelist(context.Background(), WhitelistParams{
		PolicyUUID:    PaymentSponsorPolicyId,
		WhitelistType: ToAccountWhitelist,
		Values:        []string{PaymentTokenContractAddress},
	})
	if err != nil || !success {
		log.Fatal("failed to add token contract whitelist")
	}

	// sponsor the tx that call the "transfer" interface of ERC20 contract
	success, err = sponsorClient.AddToWhitelist(context.Background(), WhitelistParams{
		PolicyUUID:    PaymentSponsorPolicyId,
		WhitelistType: ContractMethodSigWhitelist,
		Values:        []string{"0xa9059cbb"},
	})
	if err != nil || !success {
		log.Fatal("failed to add contract method whitelist")
	}

	// sponsor the tx that "transfer" stable coin to particular receiver account
	success, err = sponsorClient.AddToWhitelist(context.Background(), WhitelistParams{
		PolicyUUID:    PaymentSponsorPolicyId,
		WhitelistType: BEP20ReceiverWhiteList,
		Values:        []string{receiver.String()},
	})

	if err != nil || !success {
		log.Fatal("failed to add payment receiver whitelist")
	}

	receiverWhitelist := GetWhitelistParams{
		PolicyUUID:    PaymentSponsorPolicyId,
		WhitelistType: BEP20ReceiverWhiteList,
		Offset:        0,
		Limit:         1000,
	}

	// get the receiver whitelist, the Payment Gateway may need to update the whitelist according to its need.
	result, err := sponsorClient.GetWhitelist(context.Background(), receiverWhitelist)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Whitelist addresses:")
	for _, address := range result {
		fmt.Println(address)
	}
}

func userDoGaslessPayment(receiver common.Address, amount *big.Int) {
	// Connect to an Ethereum node (for transaction assembly)
	client, err := ethclient.Dial("https://bsc-testnet-dataseed.bnbchain.org")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}
	// Create a PaymasterClient (for transaction sending)
	paymasterClient, err := NewPaymasterClient("https://bsc-paymaster-testnet.nodereal.io")
	if err != nil {
		log.Fatalf("Failed to create PaymasterClient: %v", err)
	}

	// Load your private key
	privateKey, err := crypto.HexToECDSA(UserPrivateKey)
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
	tokenAddress := common.HexToAddress(PaymentTokenContractAddress)

	// Create ERC20 transfer data
	data, err := createERC20TransferData(receiver, amount)
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
