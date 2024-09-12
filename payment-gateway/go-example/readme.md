# GO Example

This repository contains a Go application demonstrating:

1. Payment Gateway manage the sponsor policy to sponsor any transaction that send BEP20 to them.
2. User send ERC20 token transfers to payment gateway without pay gas fee through paymaster.

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
    - 'SPONSOR_URL' to the API key created by the sponsor in the Nodereal dashboard. create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'POLICY_UUID' to the policy ID created by the sponsor on Megafuel Paymaster, create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'TOKEN_CONTRACT_ADDRESS' to the ERC20 token contract address that users want to withdraw.
    - 'RECIPIENT_ADDRESS' to the receiver address for the Payment Gateway's generated payment link.
    - 'USER_PRIVATE_KEY' to the user's account private key, ensuring this wallet contains the required ERC20 tokens.

3. Run the example
   ```
   $ go run ./
   ```


