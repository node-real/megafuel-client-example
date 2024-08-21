# GO Example
This repository contains a Go application demonstrating how to send gasless tx on BSC through Megafuel 
paymaster.

## Quick Start
If the example is performed on BSC testnet, please ensure you have some test ERC20 on BSC testnet. (You can get some
from the official faucet)

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the file
   Before running the application, you need to edit the main.go to set up the following:

   - Set `"YourPrivateKey"` with the private key of your Ethereum account.
   - Set `"TokenContractAddress"` to the address of the ERC20 token contract you want to interact with.
   - Set `"RecipientAddress"` to the Ethereum address you want to send tokens to.

3. Run the example
   ```
   $ go run main.go
   ```


