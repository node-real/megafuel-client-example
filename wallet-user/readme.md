# Wallet Client Example

This repository contains a Go and Js application demonstrating how to interact with an BSC testnet using a paymaster
service. The application showcases how to create, sign, and send ERC20 token transfers while leveraging a paymaster
for potential gas sponsorship.

## Workflow

- Connect to an BSC/opBNB network using a paymaster client
- Create an ERC20 token transfer transactions
- Check if a transaction is sponsorable
- Set the gas price to 0 and sign transaction
- Send transactions through a paymaster client