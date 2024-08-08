# Paymaster Client Example

This repository contains a Go and Js application demonstrating how to interact with an BSC testnet using a paymaster 
service. The application showcases how to create, sign, and send ERC20 token transfers while leveraging a paymaster 
for potential gas sponsorship.

## Workflow

- Connect to an Ethereum network using a paymaster client
- Create and sign ERC20 token transfer transactions
- Check if a transaction is sponsorable
- Send transactions through a paymaster client to paymaster endpoint

## Network Endpoint

BSC testnet: https://bsc-paymaster-testnet.nodereal.io

## Example

Please get ERC20 token for test before you start:
1. Visit Faucet: https://www.bnbchain.org/en/testnet-faucet
2. Claim any kind of ERC20 token except BNB.
![image](./assets/img.png)

- [Js Example](./js-example/readme.md)
- [Go Example](./go-example)

## More docs
- [Paymaster Overview](https://docs.nodereal.io/docs/maganode-paymaster-overview)
- [Sponsor Policy Management](https://docs.nodereal.io/docs/meganode-paymaster-policy-management)
- [Wallet Integration Guide](https://docs.nodereal.io/docs/wallet-integration)
- [Paymaster API Spec](https://docs.nodereal.io/docs/meganode-paymaster-api)