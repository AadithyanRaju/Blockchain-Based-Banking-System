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
const createAccount = async (call, callback) => {
  try {
    const { userID, name, aadhaarHash, email, passwordHash, phoneNumber, role, balance } = call.request;
    await blockchainService.createAccount(userID, name, aadhaarHash, email, passwordHash, phoneNumber, role, balance);
    callback(null, { message: 'Account created successfully' });
  } catch (error) {
    callback(error, null);
  }
};

const getAccount = async (call, callback) => {
  try {
    const { userID } = call.request;
    const account = await blockchainService.getAccount(userID);
    callback(null, account);
  } catch (error) {
    callback(error, null);
  }
};

const deposit = async (call, callback) => {
  try {
    const { userID, amount, referenceNumber } = call.request;
    await blockchainService.deposit(userID, amount, referenceNumber);
    callback(null, { message: 'Deposit successful' });
  } catch (error) {
    callback(error, null);
  }
};

const withdraw = async (call, callback) => {
  try {
    const { userID, amount, referenceNumber } = call.request;
    await blockchainService.withdraw(userID, amount, referenceNumber);
    callback(null, { message: 'Withdrawal successful' });
  } catch (error) {
    callback(error, null);
  }
};

const createTransfer = async (call, callback) => {
  try {
    const { senderID, receiverID, amount, referenceNumber } = call.request;
    await blockchainService.createTransfer(senderID, receiverID, amount, referenceNumber);
    callback(null, { message: 'Transfer created successfully' });
  } catch (error) {
    callback(error, null);
  }
};

const getTransfer = async (call, callback) => {
  try {
    const { senderID, receiverID, referenceNumber } = call.request;
    const transfer = await blockchainService.getTransfer(senderID, receiverID, referenceNumber);
    callback(null, transfer);
  } catch (error) {
    callback(error, null);
  }
};

// Start gRPC server
const server = new grpc.Server();
server.addService(bankingProto, {
  CreateAccount: createAccount,
  GetAccount: getAccount,
  Deposit: deposit,
  Withdraw: withdraw,
  CreateTransfer: createTransfer,
  GetTransfer: getTransfer
});

const PORT = '50051';
server.bindAsync(`0.0.0.0:${PORT}`, grpc.ServerCredentials.createInsecure(), () => {
  console.log(`gRPC server running on port ${PORT}`);
  server.start();
});
