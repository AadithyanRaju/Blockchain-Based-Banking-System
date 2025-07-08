# Blockchain-Based Banking System

## 📘 Introduction

### Overview
This project is a **Blockchain-Based Banking System** built on **Hyperledger Fabric**. It provides secure and transparent account handling, transactions, and auditability through decentralized ledger technology. Core features include:
- Immutable transaction logs
- Role-based account operations
- Cryptographic verification

---

## ⚙️ Prerequisites

Make sure you have the following installed:
- [Node.js](https://nodejs.org/)
- [Go](https://go.dev/)
- [Docker & Docker Compose](https://docs.docker.com/)
- Hyperledger Fabric binaries 

---

## NOTE
Read all the code used in the project, for safety reasons.

---

## 🛠️ Setup Instructions

1. **Clone the Repository**
   ```bash
   $ git clone https://github.com/AadithyanRaju/Blockchain-Based-Banking-System.git
   $ cd Blockchain-Based-Banking-System
   ```

2. **Start the Blockchain Network**
    Run the bash script
    ```bash
    $ bash main.sh
    ```
    Menu:
    1. Install Required Packages ( Debian and Arch based Destros)
    2. Setup the workspace (Creates and updates edited documents)
    3. Start the Blockchain
    4. Update Chaincode to the Blockchain
    5. Stop the Blockchain
    6. List Docker Containers
    7. Enter Bash/Zsh Shell (Default: Bash, enter zsh for zsh)
    8. Reset the workspace (Uses Sudo, use only if it is nessessary)
    9. Exit

    - Use option 2 to setup workspace.
    - Use option 3 to Start Blockchain network & Update Chaincode to the network.
    - Use option 7 to enter cli-mode on scripts folder.

4. **Enroll Admin and Register Users**
   ```bash
   $ node enrollAdmin.js
   $ node registerUser.js <Authorized Username>
   ```

5. **Interact with Blockchian**

    Usage: node interact.js <Autherized User> <Action>

    Actions: CreateAccount, Balance, Deposit, Withdraw, Transfer, GetAllAccounts, DeleteAccount, QueryTransaction

    CreateAccount: ```node interact.js <Autherized User> CreateAccount <AccountName> <InitialBalance>```

    Balance: ```node interact.js <Autherized User> Balance <AccountName>```

    Deposit: ```node interact.js <Autherized User> Deposit <AccountName> <Amount>```

    Withdraw: ```node interact.js <Autherized User> Withdraw <AccountName> <Amount>```

    Transfer: ```node interact.js <Autherized User> Transfer <FromAccount> <ToAccount> <Amount>```

    GetAllAccounts: ```node interact.js <Autherized User> GetAllAccounts```

    DeleteAccount: ```node interact.js <Autherized User> DeleteAccount <AccountName>```

    QueryTransaction: ```node interact.js <Autherized User> QueryTransaction <AccountName>```

---

## 📝 License
This project is licensed under the MIT License.
