# GO Example
This repository contains a Go application demonstrating how to interact with an BSC network through Mega 
Paymaster. The application showcases how to create, sign, and send ERC20 token transfers while 
leveraging a paymaster for potential gas sponsorship.


## Quick Start

1. Install dependencies
    ```shell
    $ go mod tidy
    ```
2. Configure the file
   Before running the application, you need to set up the following:

   - Set `"YOUR_PRIVATE_KEY"` with the private key of your Ethereum account.
   - Set `"TOKEN_CONTRACT_ADDRESS"` to the address of the ERC20 token contract you want to interact with.
   - Set `"RECIPIENT_ADDRESS"` to the Ethereum address you want to send tokens to.

3. Run the example
   ```
   $ go run main.go
   ```


