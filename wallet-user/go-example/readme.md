# GO Example

This repository contains a Go application demonstrating how to send gasless tx on BSC through Megafuel
paymaster.

## Quick Start

The example is performed on BSC testnet or BSC mainnet, please ensure you have some test ERC20 on them. (You can get
some from the official faucet when using testnet)

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the file
   Before running the application, you need to edit the `.env` to set up the following:

    - 'PAYMASTER_URL' with the Paymaster URL.
    - 'CHAIN_URL' with the BSC testnet or BSC mainnet chain URL.
    - 'TOKEN_CONTRACT_ADDRESS' to the ERC20 token contract address that users want to withdraw.
    - 'RECIPIENT_ADDRESS' to the address of the ERC20 token contract you want to interact with.
    - 'YOUR_PRIVATE_KEY' to the private key of your Ethereum account.

3. Run the example
   ```
   $ go run main.go
   ```


