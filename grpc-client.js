const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = './banking.proto';
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const bankingProto = grpc.loadPackageDefinition(packageDefinition).BankingService;

const client = new bankingProto('localhost:50051', grpc.credentials.createInsecure());

// Example: Create Account
client.CreateAccount({
  userID: 'USER_001',
  name: 'John Doe',
  aadhaarHash: '####',
  email: 'john@example.com',
  passwordHash: '####',
  phoneNumber: '1234567890',
  role: 'user',
  balance: 1000.0
}, (err, response) => {
  if (err) {
    console.error('Error:', err);
  } else {
    console.log(response.message);
  }
});
