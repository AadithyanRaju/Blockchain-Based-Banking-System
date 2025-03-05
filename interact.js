const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

class BlockchainService {
    constructor(connectionProfilePath, walletPath, identity) {
        this.connectionProfilePath = connectionProfilePath;
        this.walletPath = walletPath;
        this.identity = identity;
    }

    async _connect() {
        const ccp = JSON.parse(fs.readFileSync(this.connectionProfilePath, 'utf8'));
        const wallet = await Wallets.newFileSystemWallet(this.walletPath);

        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: this.identity,
            discovery: { enabled: true, asLocalhost: true }
        });

        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('banking');

        return { gateway, contract };
    }

    async createAccount(userID, name, aadhaarHash, email, passwordHash, phoneNumber, role, balance) {
        const { gateway, contract } = await this._connect();
        try {
            await contract.submitTransaction(
                'CreateAccount',
                userID,
                name,
                aadhaarHash,
                email,
                passwordHash,
                phoneNumber,
                role,
                balance.toString()
            );
            console.log('Account created successfully');
        } catch (error) {
            console.error('Error creating account:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async getAccount(userID) {
        const { gateway, contract } = await this._connect();
        try {
            const result = await contract.evaluateTransaction('GetAccount', userID);
            console.log('Account fetched successfully');
            return JSON.parse(result.toString());
        } catch (error) {
            console.error('Error fetching account:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async deposit(userID, amount, referenceNumber, timestamp) {
        const { gateway, contract } = await this._connect();
        try {
            await contract.submitTransaction(
                'Deposit',
                userID,
                amount.toString(),
                referenceNumber,
                timestamp
            );
            console.log('Deposit successful');
        } catch (error) {
            console.error('Error during deposit:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async withdraw(userID, amount, referenceNumber, timestamp) {
        const { gateway, contract } = await this._connect();
        try {
            await contract.submitTransaction(
                'Withdraw',
                userID,
                amount.toString(),
                referenceNumber,
                timestamp
            );
            console.log('Withdrawal successful');
        } catch (error) {
            console.error('Error during withdrawal:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async createTransfer(senderID, receiverID, amount, referenceNumber, timestamp) {
        const { gateway, contract } = await this._connect();
        try {
            await contract.submitTransaction(
                'TransferFunds',
                senderID,
                receiverID,
                amount.toString(),
                referenceNumber,
                timestamp
            );
            console.log('Transfer created successfully');
        } catch (error) {
            console.error('Error creating transfer:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async getTransfer(senderID, receiverID, referenceNumber) {
        const { gateway, contract } = await this._connect();
        try {
            const result = await contract.evaluateTransaction(
                'GetTransfer',
                senderID,
                receiverID,
                referenceNumber
            );
            console.log('Transfer fetched successfully');
            return JSON.parse(result.toString());
        } catch (error) {
            console.error('Error fetching transfer:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async getAllAccounts() {
        const { gateway, contract } = await this._connect();
        try {
            const result = await contract.evaluateTransaction('GetAllAccounts');
            console.log('All accounts fetched successfully');
            return JSON.parse(result.toString());
        } catch (error) {
            console.error('Error fetching all accounts:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }

    async getAllTransactions() {
        const { gateway, contract } = await this._connect();
        try {
            const result = await contract.evaluateTransaction('GetAllTransactions');
            console.log('All transactions fetched successfully');
            return JSON.parse(result.toString());
        } catch (error) {
            console.error('Error fetching all transactions:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }
    
    async getAllTransfers() {
        const { gateway, contract } = await this._connect();
        try {
            const result = await contract.evaluateTransaction('GetAllTransfers');
            console.log('All transfers fetched successfully');
            return JSON.parse(result.toString());
        } catch (error) {
            console.error('Error fetching all transfers:', error);
            throw error;
        } finally {
            gateway.disconnect();
        }
    }
}

module.exports = BlockchainService;
