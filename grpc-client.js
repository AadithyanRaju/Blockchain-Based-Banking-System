import { loadPackageDefinition, credentials } from '@grpc/grpc-js';
import { loadSync } from '@grpc/proto-loader';
import { createHash } from 'crypto';
import { create } from 'domain';

function generateTimestamp() {
  return new Date().toISOString();
}

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

const client = new bankingProto('192.168.168.128:50051', credentials.createInsecure());


// Create Accounts

async function createAccount(userID, name, aadhaar, email, password, phoneNumber, role, balance) {
  client.CreateAccount({
    userID: userID,
    name: name,
    aadhaarHash: hashWithSHA256(aadhaar),
    email: email,
    passwordHash: hashWithSHA256(password),
    phoneNumber: phoneNumber,
    role: role,
    balance: balance
  }, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response.message);
    }
  });
}

async function getAccount(userID) {
  client.GetAccount({ userID: userID }, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}

async function transfer(userID, receiverID, amount, reference) {
  client.CreateTransfer({ 
    senderID: userID, 
    receiverID: receiverID, 
    amount: amount, 
    referenceNumber: reference,
    timestamp: generateTimestamp()
  }, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response.message);
    }
  });
}

async function getTransfer(senderID, receiverID, reference) {
  client.GetTransfer({ senderID: senderID, receiverID: receiverID, referenceNumber: reference }, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}

async function getAllAccounts() {
  client.GetAllAccounts({}, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}

async function getAllTransfers() {
  client.GetAllTransfers({}, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}

async function getAllKeys() {
  client.GetAllKeys({}, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}
async function getTransferByStateKey(key) {
  client.GetTransferByStateKey({ key: key }, (err, response) => {
    if (err) {
      console.error('Error:', err);
    } else {
      console.log(response);
    }
  });
}

// Create accounts
// await createAccount('donator1', 'Donator 1', '111122223333', 'don1@p.com', 'don1', '1234567890', 'user', 10000000.0);
// await createAccount('donator2', 'Donator 2', '111122223334', 'don2@p.com', 'don2', '1234567890', 'user', 10000000.0);
// await createAccount('collector', 'Collector', '111122223335', 'col@p.com', 'col', '1234567890', 'user', 500.0);

// Transfer funds
// await transfer('donator1', 'collector', 1000, '000000000001');
// await getTransfer('don1', 'col', '000000000001');

// Get account details
// await getAccount('donator1');
// await getAccount('donator2');
// await getAccount('collector');

//Degub
// await getAllTransfers();
// await getAllKeys();
await getTransferByStateKey('TRANSACTION_TRANSFER_donator1_collector_000000000002');