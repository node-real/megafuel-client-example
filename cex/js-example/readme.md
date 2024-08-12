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

2. Edit the script.
   
    Open index.js and replace the following placeholders:
   - 'HOT_WALLET_PRIVATE_KEY' with the Cex's hotwallet private key that will do token withdrawal.
   - 'TOKEN_CONTRACT_ADDRESS' with the address of the ERC20 token user want to withdraw.
   - 'USER_WITHDRAW_ADDRESS' with the receiver account of the withdrawal request.
   - 'TOKEN_CONTRACT_ADDRESS' with the ERC20 token that the Payment Gateway support.
   - 'SPONSOR_API_KEY' to the API key created by the Sponsor in the Nodereal MegaNode Console. create one 
     from [here](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines) if you don't have it.
   - 'SPONSOR_POLICY_ID' to the policy ID created by the sponsor on Meganode Paymaster, create one
     from [here](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines) if you don't have it.

3. Run script
    ```shell
    $ npm start
    ```