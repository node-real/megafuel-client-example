# Cex Withdraw Example

When CEX users withdraw cryptocurrencies, they often need to pay corresponding network gas fees. However, 
some tokens, in order to promote on-chain adoption, can reduce the withdrawal fees to zero. This way, the fees 
can be sponsored by a third party, such as the token issuer.

This example demonstrates how a CEX should manage the sponsor policy in such a scenario.


## Prepare Work

Before getting started, the gas fee sponsor, needs to first register as a user on 
Nodereal and then apply to create a policy. The specific process can be referred to in [this document](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines).

After the application is approved, Nodereal will email the Sponsor with the ID of 
the policy created for them.

And the sponsor will add the hot-wallets of Cex into the whitelist of the sponsor policy.

## Configure the Policy

The scripts in the example demonstrate how the sponsor should configure the policy, as well as how users can 
send transactions with 0 gas price:

- The sponsor sets the policy rule through API: sponsor any transaction that sends a specific ERC20 token from a list of
Cex's withdraw hot wallets.
- Cex send the 0 gas price transaction through Paymaster endpoint.



