import { loadPackageDefinition, credentials } from '@grpc/grpc-js';
import { loadSync } from '@grpc/proto-loader';
import { createHash } from 'crypto';

const PROTO_PATH = './banking.proto';
const packageDefinition = loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const bankingProto = loadPackageDefinition(packageDefinition).BankingService;

function hashWithSHA256(data) {
  const hash = createHash('sha256');
  hash.update(data);
  return hash.digest('hex');
}

const client = new bankingProto('172.16.241.128:50051', credentials.createInsecure());

// Example: Create Account
/*
client.CreateAccount({
  userID: 'johndoe',
  name: 'John Doe',
  aadhaarHash: hashWithSHA256('111122223333'),
  email: 'john@example.com',
  passwordHash: hashWithSHA256('johndoe'),
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
*/
client.GetAccount({ userID: 'johndoe' }, (err, response) => {
  if (err) {
    console.error('Error:', err);
  } else {
    console.log(response);
  }
});