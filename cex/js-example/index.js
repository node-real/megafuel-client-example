import 'dotenv/config'
import {ethers} from 'ethers'
import {PaymasterClient, SponsorClient, WhitelistType} from 'megafuel-js-sdk'

async function cexDoGaslessWithdrawTx() {
  const policyUUID = process.env.POLICY_UUID
  const sponsorUrl = process.env.SPONSOR_URL
  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterClient = PaymasterClient.newPrivatePaymaster(sponsorUrl, policyUUID)
  const network = await paymasterClient.getNetwork()
  const wallet = new ethers.Wallet(process.env.HOTWALLET_PRIVATE_KEY)
  // ERC20 token ABI (only including the transfer function)
  const tokenAbi = ['function transfer(address to, uint256 amount) returns (bool)']
  // Create contract instance
  const tokenContract = new ethers.Contract(process.env.TOKEN_CONTRACT_ADDRESS, tokenAbi, wallet)
  // Transaction details
  const tokenAmount = ethers.parseUnits('1.0', 18) // Amount of tokens to send (adjust decimals as needed)
  // Create the transaction object
  const transaction = await tokenContract.transfer.populateTransaction(process.env.WITHDRAW_RECIPIENT_ADDRESS, tokenAmount)
  const nonce = await paymasterClient.getTransactionCount(wallet.address, 'pending')

  // Add nonce and gas settings
  transaction.from = wallet.address
  transaction.nonce = nonce
  transaction.gasLimit = 100000 // Adjust gas limit as needed for token transfers
  transaction.chainId = network.chainId
  transaction.gasPrice = 0 // Set gas price to 0

  const safeTransaction = {
    ...transaction,
    gasLimit: transaction.gasLimit.toString(),
    chainId: transaction.chainId.toString(),
    gasPrice: transaction.gasPrice.toString(),
  }

  try {
    const sponsorableInfo = await paymasterClient.isSponsorable(safeTransaction)
    console.log('Sponsorable Information:', sponsorableInfo)
  } catch (error) {
    console.error('Error checking sponsorable status:', error)
  }

  try {
    // Sign the transaction
    const signedTx = await wallet.signTransaction(transaction)
    // Send the raw transaction using the sending provider
    const tx = await paymasterClient.sendRawTransaction(signedTx)
    console.log('Transaction sent:', tx)
  } catch (error) {
    console.error('Error sending transaction:', error)
  }
}

async function sponsorSetUpPolicyRules() {
  const client = new SponsorClient(process.env.SPONSOR_URL)

  const wallet = new ethers.Wallet(process.env.HOTWALLET_PRIVATE_KEY)
  // sponsor the tx that interact with the stable coin ERC20 contract
  try {
    // You can empty the policy rules before re-try.
    await client.emptyWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.FromAccountWhitelist,
    });
    await client.emptyWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.ToAccountWhitelist,
    });

    const res1 = await client.addToWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.ToAccountWhitelist,
      Values: [process.env.TOKEN_CONTRACT_ADDRESS]
    });
    console.log("Added ERC20 contract address to whitelist ", res1);

    // sponsor the tx that sent by hotwallet
    const res2 = await client.addToWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.FromAccountWhitelist,
      Values: [wallet.address]
    });
    console.log("Added hotwallet to whitelist ", res2);
  } catch (error){
    console.error("Error:", error)
  }
}

async function main() {
  try {
    await sponsorSetUpPolicyRules()
    await cexDoGaslessWithdrawTx()
  } catch (error) {
    console.error('Error:', error)
  }
}

main()
