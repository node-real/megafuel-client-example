# Payment Gateway Example

In payment scenarios, it's most common for users to pay with stablecoins. However, users might not possess the 
gas token of the specific blockchain network. Therefore, Meganode allows payment gateways or other third parties 
to pay the gas fee on behalf of the users.

## Prepare Work

Before getting started, the payment gateway, acting as a gas fee sponsor, needs to first register as a user on 
Nodereal and then apply to create a policy. The specific process can be referred to in [this document](https://docs.nodereal.io/docs/meganode-paymaster-sponsor-guidelines).

After the application is approved, Nodereal will send an email to the Payment Gateway with the ID of 
the policy created for them.

## Configure the Policy

The scripts in the example demonstrate how the payment gateway should configure the policy, as well as how users can 
send transactions with 0 gas price:

- The payment gateway sets the policy rule through API: sponsor any transaction that sends a specific ERC20 token 
to a fixed set of payment gateway receiving addresses.
- User send the 0 gas price transaction through Paymaster endpoint.



