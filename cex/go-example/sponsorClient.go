package main

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

type WhitelistType string

const (
	FromAccountWhitelist       WhitelistType = "FromAccountWhitelist"
	ToAccountWhitelist         WhitelistType = "ToAccountWhitelist"
	ContractMethodSigWhitelist WhitelistType = "ContractMethodSigWhitelist"
	BEP20ReceiverWhiteList     WhitelistType = "BEP20ReceiverWhiteList"
)

// WhitelistParams represents the parameters for the pm_addToWhitelist and pm_rmFromWhitelist methods
type WhitelistParams struct {
	PolicyUUID    string        `json:"policyUuid"`
	WhitelistType WhitelistType `json:"whitelistType"`
	Values        []string      `json:"values"`
}

// EmptyWhitelistParams represents the parameters for the pm_emptyWhitelist method
type EmptyWhitelistParams struct {
	PolicyUUID    string        `json:"policyUuid"`
	WhitelistType WhitelistType `json:"whitelistType"`
}

// GetWhitelistParams represents the parameters for the pm_getWhitelist method
type GetWhitelistParams struct {
	PolicyUUID    string        `json:"policyUuid"`
	WhitelistType WhitelistType `json:"whitelistType"`
	Offset        int           `json:"offset"`
	Limit         int           `json:"limit"`
}

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c *rpc.Client
}

// AddToWhitelist calls the pm_addToWhitelist method
func (ec *Client) AddToWhitelist(ctx context.Context, params WhitelistParams) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "pm_addToWhitelist", params)
	return result, err
}

// RemoveFromWhitelist calls the pm_rmFromWhitelist method
func (ec *Client) RemoveFromWhitelist(ctx context.Context, params WhitelistParams) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "pm_rmFromWhitelist", params)
	return result, err
}

// EmptyWhitelist calls the pm_emptyWhitelist method
func (ec *Client) EmptyWhitelist(ctx context.Context, params EmptyWhitelistParams) (bool, error) {
	var result bool
	err := ec.c.CallContext(ctx, &result, "pm_emptyWhitelist", params)
	return result, err
}

// GetWhitelist calls the pm_getWhitelist method
func (ec *Client) GetWhitelist(ctx context.Context, params GetWhitelistParams) ([]string, error) {
	var result []string
	err := ec.c.CallContext(ctx, &result, "pm_getWhitelist", params)
	return result, err
}

// NewSponsorClient creates a client that uses the given URL.
func NewSponsorClient(url string) (*Client, error) {
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{c: rpcClient}, nil
}
