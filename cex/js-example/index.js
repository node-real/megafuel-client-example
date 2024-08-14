const ethers = require('ethers');

// Replace with the cex hot wallet private key (be cautious with private keys!)
const hotwalletPrivateKey = 'HOT_WALLET_PRIVATE_KEY';
// replace with user ERC20 withdraw address
const userWithdrawAddress = 'USER_WITHDRAW_ADDRESS';
// ERC20 token contract address (replace with the address of the token you want to send)
const erc20TokenAddress = 'TOKEN_CONTRACT_ADDRESS';
const sponsorEndpoint = 'https://open-platform.nodereal.io/{SPONSOR_API_KEY}/megafuel-testnet';
const policyID = 'SPONSOR_POLICY_ID'


class SponsorProvider extends ethers.providers.JsonRpcProvider {
  constructor(url) {
    super(url);
  }

  async addToWhitelist(params) {
    return this.send('pm_addToWhitelist', [params]);
  }

  async removeFromWhitelist(params) {
    return this.send('pm_rmFromWhitelist', [params]);
  }

  async emptyWhitelist(params) {
    return this.send('pm_emptyWhitelist', [params]);
  }

  async getWhitelist(params) {
    return this.send('pm_getWhitelist', [params]);
  }
}

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

async function cexDoGaslessWithdrawTx() {

  // Provider for assembling the transaction (e.g., mainnet)
  const assemblyProvider = new ethers.providers.JsonRpcProvider('https://bsc-testnet-dataseed.bnbchain.org');

  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterProvider = new PaymasterProvider('https://bsc-megafuel-testnet.nodereal.io');

  const wallet = new ethers.Wallet(hotwalletPrivateKey, assemblyProvider);
  // ERC20 token ABI (only including the transfer function)
  const tokenAbi = [
    "function transfer(address to, uint256 amount) returns (bool)"
  ];

  // Create contract instance
  const tokenContract = new ethers.Contract(erc20TokenAddress, tokenAbi, wallet);

  // Transaction details
  const tokenAmount = ethers.utils.parseUnits('1.0', 18); // Amount of tokens to send (adjust decimals as needed)

  try {
    // Get the current nonce for the sender's address
    const nonce = await assemblyProvider.getTransactionCount(wallet.address);

    // Create the transaction object
    const transaction = await tokenContract.populateTransaction.transfer(userWithdrawAddress, tokenAmount);

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

async function sponsorSetUpPolicyRules() {
  const client = new SponsorProvider(sponsorEndpoint);

  const wallet = new ethers.Wallet(hotwalletPrivateKey)
  // sponsor the tx that interact with the stable coin ERC20 contract
  try {
    // You can empty the policy rules before re-try.
    // await client.emptyWhitelist({
    // policyUuid: policyID,
    //  whitelistType: "FromAccountWhitelist",
    // });
    //await client.emptyWhitelist({
    //  policyUuid: policyID,
    //  whitelistType: "ToAccountWhitelist",
    // });
    // sponsor the tx that interact with the stable coin ERC20 contract
    const res1 = await client.addToWhitelist({
      policyUuid: policyID,
      whitelistType: "ToAccountWhitelist",
      values: [erc20TokenAddress]
    });
    console.log("Added ERC20 contract address to whitelist ", res1);

    // sponsor the tx that sent by hotwallet
    const res2 = await client.addToWhitelist({
      policyUuid: policyID,
      whitelistType: "FromAccountWhitelist",
      values: [wallet.address]
    });
    console.log("Added hotwallet to whitelist ", res2);
  } catch (error){
    console.error("Error:", error)
  }
}

async function main() {
  try {
    await sponsorSetUpPolicyRules();
    await cexDoGaslessWithdrawTx();
  } catch (error) {
    console.error("Error:", error);
  }
}

main();