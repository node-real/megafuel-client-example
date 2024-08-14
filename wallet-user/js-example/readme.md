# Js Example 
This is a Js script to send BSC transactions with 0 gas price through Meganode Paymaster.

## Quick Start

The example is performed on BSC testnet, please ensure you have some test ERC20 on BSC testnet. (You can get some 
from the official faucet)

1. Install the dependency.
    ```shell
    $ npm install
    ```

2. Edit the script.
   
    Open index.js and replace the following placeholders:
   - 'YOUR_PRIVATE_KEY' with your BSC private key
   - 'TOKEN_CONTRACT_ADDRESS' with the address of the ERC20 token you want to send
   - 'RECIPIENT_ADDRESS' with the BSC address you want to send tokens to

    Please ensure that the unlocked account is allowed to send gasless tx.

3. Run script
    ```shell
    $ npm start
    ```