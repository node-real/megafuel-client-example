# GO Example

This repository contains a Go application demonstrating:

1. Sponsor manage the policy to sponsor any transaction sent particular tokens to the consolidation/hot
   wallets of the Cex
2. Cex do token transfer without pay gas fee through MegaFuel.

## Quick Start

The example is performed on BSC testnet or BSC mainnet, please ensure you have some test ERC20 on them. (You can get
some from the official faucet when using testnet)

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the .env file
   Before running the application, you need to edit the `.env` to set up the following:

    - 'PAYMASTER_URL' with the Paymaster URL.
    - 'CHAIN_URL' with the BSC testnet or BSC mainnet chain URL.
    - 'SPONSOR_URL' to the API key created by the sponsor in the Nodereal dashboard. create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'POLICY_UUID' to the policy ID created by the sponsor on Megafuel Paymaster, create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'TOKEN_CONTRACT_ADDRESS' to the ERC20 token contract address that Cex want to collect.
    - 'CONSOLIDATION_WALLET_ADDRESS' to the consolidation/hot wallet of Cex.
    - 'DEPOSIT_WALLET_PRIVATE_KEY' to the Cex's deposit wallet private key, ensuring this wallet contains the required ERC20
      tokens.

3. Run the example
   ```
   $ go run ./
   ```


