# Js Example 
This is a Js script to send Ethereum transactions with 0 gas price through Meganode Paymaster.

## Quick Start

1. Install the dependency.
    ```shell
    $npm install
    ```

2. Edit the script.
   
    Open index.js and replace the following placeholders:
   - 'YOUR_PRIVATE_KEY' with your Ethereum private key
   - 'TOKEN_CONTRACT_ADDRESS' with the address of the ERC20 token you want to send
   - 'RECIPIENT_ADDRESS' with the Ethereum address you want to send tokens to

    Please ensure that the unlocked account is allowed to send gasless tx.

3. Run script
    ```shell
    $npm start
    ```