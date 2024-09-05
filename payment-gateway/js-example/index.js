const ethers = require('ethers');

// Replace with your private key (be cautious with private keys!)
const userPrivateKey = 'USER_PRIVATE_KEY';
// replace with your ERC20 receiver
const paymentReceiverAddress = 'PAYMENT_RECIPIENT_ADDRESS';
// ERC20 token contract address (replace with the address of the token you want to send)
const erc20TokenAddress = 'TOKEN_CONTRACT_ADDRESS';
const policyID = 'SPONSOR_POLICY_ID'

const paymasterEndpoint = 'https://bsc-megafuel.nodereal.io';
const sponsorEndpoint = 'https://open-platform.nodereal.io/{SPONSOR_API_KEY}/megafuel';

// testnet endpoint
// const paymasterEndpoint = 'https://bsc-megafuel-testnet.nodereal.io';
// const sponsorEndpoint = 'https://open-platform.nodereal.io/{SPONSOR_API_KEY}/megafuel-testnet';

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

async function userDoGaslessPayment() {
  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterProvider = new PaymasterProvider(paymasterEndpoint);

  const wallet = new ethers.Wallet(userPrivateKey);
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
    const nonce = await paymasterProvider.getTransactionCount(wallet.address);
    const network = await paymasterProvider.getNetwork();

    // Create the transaction object
    const transaction = await tokenContract.populateTransaction.transfer(paymentReceiverAddress, tokenAmount);

    // Add nonce and gas settings
    transaction.chainId = network.chainId;
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

async function paymentGatewaySetUpPolicyRules() {
  const client = new SponsorProvider(sponsorEndpoint);

  // sponsor the tx that interact with the stable coin ERC20 contract
  try {
    // sponsor the tx that interact with the stable coin ERC20 contract
    const res1 = await client.addToWhitelist({
      policyUuid: policyID,
      whitelistType: "ToAccountWhitelist",
      values: [erc20TokenAddress]
    });
    console.log("Added ERC20 contract address  to whitelist ", res1);

    // sponsor the tx that call the "transfer" interface of ERC20 contract
    const res2 = await client.addToWhitelist({
      policyUuid: policyID,
      whitelistType: "ContractMethodSigWhitelist",
      values: ["0xa9059cbb"]
    });
    console.log("Added 'transfer' contract method  to whitelist ", res2);

    // sponsor the tx that "transfer" stable coin to particular receiver account
    const res3 = await client.addToWhitelist({
      policyUuid: policyID,
      whitelistType: "BEP20ReceiverWhiteList",
      values: [paymentReceiverAddress]
    });
    console.log("Added BEP20 transfer receiver to whitelist ", res3);
  } catch (error){
    console.error("Error:", error)
  }

  try {
    const params = {
      policyUuid: policyID,
      whitelistType: "BEP20ReceiverWhiteList",
      offset: 0,
      limit: 1000
    };

    const result = await client.getWhitelist(params);
    console.log("Whitelist addresses:", result);
  } catch (error) {
    console.error("Error:", error);
  }
}

async function main() {
  try {
    await paymentGatewaySetUpPolicyRules();
    await userDoGaslessPayment();
  } catch (error) {
    console.error("Error:", error);
  }
}

main();