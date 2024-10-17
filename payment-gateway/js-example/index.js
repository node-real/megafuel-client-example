import 'dotenv/config';
import {ethers} from "ethers";
import {PaymasterClient, SponsorClient, WhitelistType} from 'megafuel-js-sdk';

async function userDoGaslessPayment() {
  // Provider for sending the transaction (e.g., could be a different network or provider)
  const paymasterClient = PaymasterClient.new(process.env.PAYMASTER_URL);
  const network = await paymasterClient.getNetwork()
  const wallet = new ethers.Wallet(process.env.USER_PRIVATE_KEY);
  // ERC20 token ABI (only including the transfer function)
  const tokenAbi = ["function transfer(address to, uint256 amount) returns (bool)"];
  // Create contract instance
  const tokenContract = new ethers.Contract(process.env.TOKEN_CONTRACT_ADDRESS, tokenAbi, wallet);
  // Transaction details
  const tokenAmount = ethers.parseUnits('1.0', 18); // Amount of tokens to send (adjust decimals as needed)
  // Create the transaction object
  const transaction = await tokenContract.transfer.populateTransaction(process.env.RECIPIENT_ADDRESS, tokenAmount);
  const nonce = await paymasterClient.getTransactionCount(wallet.address, 'pending');

  // Add nonce and gas settings
  transaction.from = wallet.address;
  transaction.nonce = nonce;
  transaction.gasLimit = 100000; // Adjust gas limit as needed for token transfers
  transaction.chainId = network.chainId;
  transaction.gasPrice = 0; // Set gas price to 0

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
    // We strongly encourage you to set the UserAgent value. It should represent
    // your wallet name or brand name. This information is for further statistical
    // analysis and insight. Setting a unique UserAgent will help MegaFuel to
    // better understand wallet usage patterns and improve service.
    const txOpt = {
      UserAgent: "myWalletName/v1.0.0"
    }
    // Send the raw transaction using the sending provider
    const tx = await paymasterClient.sendRawTransaction(signedTx, txOpt);
    console.log('Transaction sent:', tx);
  } catch (error) {
    console.error('Error sending transaction:', error);
  }
}

async function paymentGatewaySetUpPolicyRules() {
  const paymasterClient = new PaymasterClient(process.env.PAYMASTER_URL)
  const network = await paymasterClient.getNetwork();
  const client = new SponsorClient(process.env.SPONSOR_URL, null,
      { staticNetwork: ethers.Network.from(network.chainId) });

  try {
    // sponsor the tx that interact with the stable coin ERC20 contract
    const res1 = await client.addToWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.ToAccountWhitelist,
      Values: [process.env.TOKEN_CONTRACT_ADDRESS]
    });
    console.log("Added ERC20 contract address  to whitelist ", res1);

    // sponsor the tx that call the "transfer" interface of ERC20 contract
    const res2 = await client.addToWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.ContractMethodSigWhitelist,
      Values: ["0xa9059cbb"]
    });
    console.log("Added 'transfer' contract method  to whitelist ", res2);

    // // sponsor the tx that "transfer" stable coin to particular receiver account
    const res3 = await client.addToWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.BEP20ReceiverWhiteList,
      Values: [process.env.RECIPIENT_ADDRESS]
    });
    console.log("Added BEP20 transfer receiver to whitelist ", res3);

    const result = await client.getWhitelist({
      PolicyUUID: process.env.POLICY_UUID,
      WhitelistType: WhitelistType.BEP20ReceiverWhiteList,
      Offset: 0,
      Limit: 1000
    });
    console.log("Whitelist addresses:", result);
  } catch (error){
    console.error("Error:", error)
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
