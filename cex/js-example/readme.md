# Js Example

This repository contains a Js application demonstrating:

1. Sponsor manage the policy to sponsor any transaction sent by Cex hotwallets.
2. Cex do token withdrawal without pay gas fee through paymaster.

## Quick Start

The example is performed on BSC testnet, please ensure you have some test ERC20 on BSC testnet. (You can get some
from the official faucet)

1. Install the dependency.
    ```shell
    $ npm install
    ```

2. Configure the `.env`.
   Open `.env` and replace the following placeholders:

    - 'PAYMASTER_URL' with the Paymaster URL.
    - 'SPONSOR_URL' to the API key created by the sponsor in the Nodereal dashboard. create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'POLICY_UUID' to the policy UUID created by the sponsor on Megafuel Paymaster, create one
      from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
    - 'TOKEN_CONTRACT_ADDRESS' with the address of the ERC20 token user want to withdraw.
    - 'WITHDRAW_RECIPIENT_ADDRESS' with the receiver account of the withdrawal request.
    - 'HOTWALLET_PRIVATE_KEY' with the Cex's hotwallet private key that will do token withdrawal.

3. Run script
    ```shell
    $ npm start
    ```
