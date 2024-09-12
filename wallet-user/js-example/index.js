import 'dotenv/config';
import {ethers} from "ethers";
import { PaymasterClient } from 'megafuel-js-sdk';

async function sendERC20Transaction() {
  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterClient = new PaymasterClient(process.env.PAYMASTER_URL);
  const network = await paymasterClient.getNetwork()
  const wallet = new ethers.Wallet(process.env.YOUR_PRIVATE_KEY);
  // ERC20 token ABI (only including the transfer function)
  const tokenAbi = ["function transfer(address,uint256) returns (bool)"];
  // Create contract instance
  const tokenContract = new ethers.Contract(process.env.TOKEN_CONTRACT_ADDRESS, tokenAbi, wallet);
  // Transaction details
  const tokenAmount = ethers.parseUnits('1.0', 18); // Amount of tokens to send (adjust decimals as needed)
  // Create the transaction object
  const transaction = await tokenContract.transfer.populateTransaction(process.env.RECIPIENT_ADDRESS, tokenAmount)
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
    const sponsorableInfo = await paymasterClient.isSponsorable(safeTransaction);
    console.log('Sponsorable Information:', sponsorableInfo);
  } catch (error) {
    console.error('Error checking sponsorable status:', error);
  }

    try {
    // Sign the transaction
    const signedTx = await wallet.signTransaction(transaction);
    // Send the raw transaction using the sending provider
    const tx = await paymasterClient.sendRawTransaction(signedTx);
    console.log('Transaction sent:', tx);

  } catch (error) {
    console.error('Error sending transaction:', error);
  }
}

sendERC20Transaction();
