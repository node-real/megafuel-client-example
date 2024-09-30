# Cex Deposit Wallet Collection Example

MegaFuel provides a gasless solution for EOA (Externally Owned Account) users. 
By integrating with MegaFuel, CEXs can achieve deposit wallet collection in a single transaction with minimal 
modifications to their existing systems. This integration simultaneously reduces complexity, lowers costs, and enhances g
as token utilization efficiency.

This example demonstrates how a CEX should manage the sponsor policy in such a scenario.

## Prepare Work

Before getting started, the Cex, needs to first register as a user on 
Nodereal and then apply to create a policy. The specific process can be referred to in [this document](https://docs.nodereal.io/docs/megafuel-sponsor-guidelines).

After the application is approved, Nodereal will email the Sponsor with the ID of 
the policy created for them.

## Configure the Policy

The scripts in the example demonstrate how the sponsor should configure the policy, as well as how Cex can 
send transactions with 0 gas price:

- The sponsor sets the policy rule through API:  Any transactions that send particular tokens to the consolidation/hot 
wallets of the Cex can be sponsored.
- Cex send the 0 gas price transaction through Paymaster endpoint.



