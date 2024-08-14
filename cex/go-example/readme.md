# GO Example
This repository contains a Go application demonstrating:
1. Sponsor manage the policy to sponsor any transaction sent by Cex hotwallets.
2. Cex do token withdrawal without pay gas fee through paymaster.

## Quick Start

The example is performed on BSC testnet, please ensure you have some test ERC20 on BSC testnet. (You can get some
from the official faucet)

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the file
   Before running the application, you need to edit the `main.go` to set up the following:

   - Set "TokenContractAddress" to the ERC20 token contract address that users want to withdraw.
   - Set "WithdrawRecipientAddress" to the receiver address of user's withdrawal request.
   - Set "SponsorPolicyId" to the policy ID created by the sponsor on Megafuel Paymaster, create one 
   from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
   - Set "SponsorAPIEndpoint" to the API key created by the sponsor in the Nodereal dashboard.
     create one from [here](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines) if you don't have it.
   - Set "HotwalletPrivateKey" to the Cex's hotwallet private key, ensuring this wallet contains the required ERC20 tokens.

3. Run the example
   ```
   $ go run ./
   ```


