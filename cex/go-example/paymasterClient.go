package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

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
