# Paymaster Example

This repository hosts a collection of examples implemented in both Golang and JavaScript for the [Meganode Paymaster](https://docs.nodereal.io/docs/meganode-paymaster-overview).
The client implementation follows the API standards of [BEP-414](https://github.com/bnb-chain/BEPs/blob/master/BEPs/BEP-414.md).
The examples include:

- Wallet integration
- Centralized Exchange (CEX) integration
- Payment gateway integration


## Network Endpoint
BSC testnet: https://bsc-paymaster-testnet.nodereal.io

## Quick Start

Please get ERC20 token for test before you start:
1. Visit Faucet: https://www.bnbchain.org/en/testnet-faucet
2. Claim any kind of ERC20 token except BNB.
![image](./assets/img.png)
3. Follow the detailed instructions to run examples:

- [For wallet integration](./wallet-user/readme.md)
- [For payment gateway integration](./payment-gateway/readme.md)
- [For Cex integration](./cex/readme.md)

## More docs
- [Paymaster Overview](https://docs.nodereal.io/docs/maganode-paymaster-overview)
- [Sponsor Policy Management](https://docs.nodereal.io/docs/meganode-paymaster-policy-management)
- [Wallet Integration Guide](https://docs.nodereal.io/docs/wallet-integration)
- [Paymaster API Spec](https://docs.nodereal.io/docs/meganode-paymaster-api)