const ethers = require('ethers');

// Replace with your private key (be cautious with private keys!)
const privateKey = 'YOUR_PRIVATE_KEY';
// replace with your ERC20 receiver
const toAddress = 'RECIPIENT_ADDRESS';
// ERC20 token contract address (replace with the address of the token you want to send)
const tokenAddress = 'TOKEN_CONTRACT_ADDRESS';

const web3ProviderEndpoint = 'https://bsc-dataseed.bnbchain.org';
const paymasterEndpoint = 'https://bsc-megafuel.nodereal.io';

// testnet endpoint
// const web3ProviderEndpoint = 'https://bsc-testnet-dataseed.bnbchain.org';
// const paymasterEndpoint = 'https://bsc-megafuel-testnet.nodereal.io';


class PaymasterProvider extends ethers.providers.JsonRpcProvider {
  constructor(url) {
    super(url);
  }

  async isSponsorable(transaction) {
    const params = [{
      to: transaction.to,
      from: transaction.from,
      value: transaction.value != null ? ethers.utils.hexlify(transaction.value) : '0x0',
      gas: ethers.utils.hexlify(transaction.gasLimit || 0),
      data: transaction.data || '0x'
    }];

    const result = await this.send('pm_isSponsorable', params);
    return result;
  }
}

async function sendERC20Transaction() {

  // Provider for assembling the transaction (e.g., mainnet)
  const assemblyProvider = new ethers.providers.JsonRpcProvider(web3ProviderEndpoint);

  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterProvider = new PaymasterProvider(paymasterEndpoint);

  const wallet = new ethers.Wallet(privateKey, assemblyProvider);
  // ERC20 token ABI (only including the transfer function)
  const tokenAbi = [
    "function transfer(address to, uint256 amount) returns (bool)"
  ];

  // Create contract instance
  const tokenContract = new ethers.Contract(tokenAddress, tokenAbi, wallet);

  // Transaction details
  const tokenAmount = ethers.utils.parseUnits('1.0', 18); // Amount of tokens to send (adjust decimals as needed)

  try {
    // Get the current nonce for the sender's address
    const nonce = await assemblyProvider.getTransactionCount(wallet.address);

    // Create the transaction object
    const transaction = await tokenContract.populateTransaction.transfer(toAddress, tokenAmount);

    // Add nonce and gas settings
    transaction.nonce = nonce;
    transaction.gasPrice = 0; // Set gas price to 0
    transaction.gasLimit = 100000; // Adjust gas limit as needed for token transfers

    try {
      const sponsorableInfo = await paymasterProvider.isSponsorable(transaction);
      console.log('Sponsorable Information:', sponsorableInfo);
    } catch (error) {
      console.error('Error checking sponsorable status:', error);
    }

    // Sign the transaction
    const signedTx = await wallet.signTransaction(transaction);

    // Send the raw transaction using the sending provider
    const tx = await paymasterProvider.send('eth_sendRawTransaction', [signedTx]);
    console.log('Transaction sent:', tx);

  } catch (error) {
    console.error('Error sending transaction:', error);
  }
}

sendERC20Transaction();
