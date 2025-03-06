const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const BlockchainService = require('./interact');

const PROTO_PATH = './banking.proto';

// Load proto file
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const bankingProto = grpc.loadPackageDefinition(packageDefinition).BankingService;

// Blockchain service setup
const blockchainService = new BlockchainService(
  './go/src/github.com/AadithyanRaju/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/connection-org1.json', 
  './wallet',
  'admin'
);

// Implement gRPC methods
const implementations = {
  CreateAccount: async (call, callback) => {
    try {
      const { userID, name, aadhaarHash, email, passwordHash, phoneNumber, role, balance } = call.request;
      await blockchainService.createAccount(userID, name, aadhaarHash, email, passwordHash, phoneNumber, role, balance);
      callback(null, { message: 'Account created successfully' });
    } catch (error) {
      callback(error, null);
    }
  },
  GetAccount: async (call, callback) => {
    try {
      const { userID } = call.request;
      const account = await blockchainService.getAccount(userID);
      callback(null, account);
    } catch (error) {
      callback(error, null);
    }
  },
  Deposit: async (call, callback) => {
    try {
      const { userID, amount, referenceNumber, timestamp } = call.request;
      await blockchainService.deposit(userID, amount, referenceNumber, timestamp);
      callback(null, { message: 'Deposit successful' });
    } catch (error) {
      callback(error, null);
    }
  },
  Withdraw: async (call, callback) => {
    try {
      const { userID, amount, referenceNumber, timestamp } = call.request;
      await blockchainService.withdraw(userID, amount, referenceNumber, timestamp);
      callback(null, { message: 'Withdrawal successful' });
    } catch (error) {
      callback(error, null);
    }
  },
  CreateTransfer: async (call, callback) => {
    try {
      const { senderID, receiverID, amount, referenceNumber, timestamp } = call.request;
      await blockchainService.createTransfer(senderID, receiverID, amount, referenceNumber, timestamp);
      callback(null, { message: 'Transfer created successfully' });
    } catch (error) {
      callback(error, null);
    }
  },
  GetTransfer: async (call, callback) => {
    try {
      const { senderID, receiverID, referenceNumber } = call.request;
      const transfer = await blockchainService.getTransfer(senderID, receiverID, referenceNumber);
      callback(null, transfer);
    } catch (error) {
      callback(error, null);
    }
  },
  GetAllAccounts: async (call, callback) => {
    try {
      const accounts = await blockchainService.getAllAccounts();
      callback(null, { accounts });
    } catch (error) {
      handleError(error, callback);
    }
  },
  GetAllTransfers: async (call, callback) => {
    try {
      const transfers = await blockchainService.getAllTransfers();
      callback(null, { transfers });
    } catch (error) {
      handleError(error, callback);
    }
  },
  GetAllKeys: async (call, callback) => {
    try {
        const { gateway, contract } = await blockchainService._connect();
        const result = await contract.evaluateTransaction('GetAllKeys');
        console.log('Ledger Keys:', result.toString());
        callback(null, { keys: JSON.parse(result.toString()) });
    } catch (error) {
        callback(error, null);
    }
  },
  GetTransferByStateKey: async (call, callback) => {
    try {
      const { stateKey } = call.request;
      const transfer = await blockchainService.getTransferByStateKey(stateKey);
      callback(null, transfer);
    } catch (error) {
      callback(error, null);
    }
  }
   
};

const server = new grpc.Server();
server.addService(bankingProto.service, implementations); 
const PORT = '50051';
server.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), () => {
  console.log(`gRPC server running on port ${PORT}`);
  server.start();
});