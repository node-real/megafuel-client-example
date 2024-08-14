# GO Example
This repository contains a Go application demonstrating:
1. Payment Gateway manage the sponsor policy to sponsor any transaction that send BEP20 to them.
2. User send ERC20 token transfers to payment gateway without pay gas fee through paymaster.

## Quick Start

The example is performed on BSC testnet, please ensure you have some test ERC20 on BSC testnet. (You can get some
from the official faucet)

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the file
   Before running the application, you need to edit the `main.go` to set up the following:

   - Set "PaymentTokenContractAddress" to the ERC20 token contract address users will use for payment.
   - Set "PaymentRecipientAddress" to the receiver address for the Payment Gateway's generated payment link.
   - Set "PaymentSponsorPolicyId" to the policy ID created by the Payment Gateway on Meganode Paymaster, create one 
   from [here](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines) if you don't have it.
   - Set "SponsorAPIEndpoint" to the API key created by the Payment Gateway in the Nodereal MegaNode Console.
     create one from [here](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines) if you don't have it.
   - Set "UserPrivateKey" to the user's account private key, ensuring this wallet contains the required ERC20 tokens.

3. Run the example
   ```
   $ go run ./
   ```


