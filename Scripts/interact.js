const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

async function main() {
    try {
        // Load connection profile
        const ccpPath = path.resolve(__dirname, '..','..', 'fabric-samples', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create wallet
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const username = process.argv[2];
        const identity = await wallet.get(username);
        if (!identity) {
            console.log('An identity for the user "'+username+'" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        const timestamp = new Date().toISOString();
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: username, discovery: { enabled: true, asLocalhost: true } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('banking');
        switch (process.argv[3]) {
            case 'Help':
                console.log('Usage: node interact.js <Autherized User> <Action>');
                console.log('Actions: CreateAccount, Balance, Deposit, Withdraw, Transfer, GetAllAccounts, DeleteAccount, QueryTransaction\n');
                console.log('CreateAccount: node interact.js <Autherized User> CreateAccount <AccountName> <InitialBalance>');
                console.log('Balance: node interact.js <Autherized User> Balance <AccountName>');
                console.log('Deposit: node interact.js <Autherized User> Deposit <AccountName> <Amount>');
                console.log('Withdraw: node interact.js <Autherized User> Withdraw <AccountName> <Amount>');
                console.log('Transfer: node interact.js <Autherized User> Transfer <FromAccount> <ToAccount> <Amount>');
                console.log('GetAllAccounts: node interact.js <Autherized User> GetAllAccounts');
                console.log('DeleteAccount: node interact.js <Autherized User> DeleteAccount <AccountName>');
                console.log('QueryTransaction: node interact.js <Autherized User> QueryTransaction <AccountName>');
                break;
            case 'Balance':
            case 'bal':
                const result1 = await contract.evaluateTransaction('QueryAccount', process.argv[4]);
                console.table(JSON.parse(result1.toString()));
                break;
            case 'CreateAccount':
            case 'ca':
                const userExists1 = await contract.evaluateTransaction('UserExists', process.argv[4]);
                if (userExists1.toString() === 'true') {
                    console.log(`Account "${process.argv[4]}" already exists.`);
                    return;
                }
                await contract.submitTransaction('CreateAccount', process.argv[4], process.argv[5]);
                console.log(`Transaction has been submitted: Account created ${process.argv[4]} with balance ${process.argv[5]}`);
                break;
            case 'Deposit':
            case 'dep':
                await contract.submitTransaction('Deposit', process.argv[4], process.argv[5], timestamp);
                console.log('Transaction has been submitted: Deposit successful');
                break;
            case 'Withdraw':
            case 'with':
                try {
                    await contract.submitTransaction('Withdraw', process.argv[4], process.argv[5], timestamp);
                    console.log('Transaction has been submitted: Withdraw successful');
                } catch (error) {
                    console.error(`Failed to submit transaction: ${error}`);
                }
                break;
            case 'Transfer':
            case 'transfer':
                try {
                    await contract.submitTransaction('Transfer', process.argv[4], process.argv[5], process.argv[6], timestamp);
                    console.log(`Transaction has been submitted: Transfer from ${process.argv[4]} to ${process.argv[5]} of amount ${process.argv[6]} successful`);
                } catch (error) {
                    console.error(`Failed to submit transaction: ${error}`);
                }
                break;
            case 'GetAllAccounts':
            case 'gaa':
                const allAccountsResult = await contract.evaluateTransaction('GetAllAccounts');
                console.table(JSON.parse(allAccountsResult.toString()));
                break;
            case 'DeleteAccount':
            case 'da':
                const userExists = await contract.evaluateTransaction('UserExists', process.argv[4]);
                if (userExists.toString() === 'false') {
                    console.log(`Account "${process.argv[4]}" does not exist.`);
                    return;
                }
                await contract.submitTransaction('DeleteAccount', process.argv[4]);
                console.log('Transaction has been submitted: Account deleted');
                break;
            case 'QueryTransaction':
            case 'qt':
                const result = await contract.evaluateTransaction('QueryTransactions', process.argv[4]);
                console.table(JSON.parse(result.toString()));
                break;
            default:
                console.log('Invalid action');
        }

        // Disconnect from the gateway.
        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
