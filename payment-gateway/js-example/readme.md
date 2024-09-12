# Js Example

This repository contains a Js application demonstrating:

1. Payment Gateway manage the sponsor policy to sponsor any transaction that send BEP20 to them.
2. User send ERC20 token transfers to payment gateway without pay gas fee through Megafuel paymaster.

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
    - 'RECIPIENT_ADDRESS' with the receiver account of the withdrawal request.
    - 'USER_PRIVATE_KEY' with the user's account private key that will do token withdrawal.

3. Run script
    ```shell
    $ npm start
    ```
