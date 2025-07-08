package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Account struct {
    ID           string        `json:"id"`
    Balance      int           `json:"balance"`
    Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
    ID        string `json:"id"`
    Type      string `json:"type"` // "deposit", "withdraw", or "transfer"
    Amount    int    `json:"amount"`
    Timestamp string `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    accounts := []Account{
        {ID: "acc1", Balance: 1000, Transactions: []Transaction{}},
        {ID: "acc2", Balance: 2000, Transactions: []Transaction{}},
    }

    for _, account := range accounts {
        accountJSON, err := json.Marshal(account)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(account.ID, accountJSON)
        if err != nil {
            return fmt.Errorf("failed to put account to world state. %v", err)
        }
    }

    return nil
}

func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, balance int) error {
    account := Account{
        ID:           id,
        Balance:      balance,
        Transactions: []Transaction{},
    }

    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, accountJSON)
}

func (s *SmartContract) QueryAccount(ctx contractapi.TransactionContextInterface, id string) (*Account, error) {
    accountJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return nil, fmt.Errorf("failed to read from world state. %v", err)
    }
    if accountJSON == nil {
        return nil, fmt.Errorf("the account %s does not exist", id)
    }

    var account Account
    err = json.Unmarshal(accountJSON, &account)
    if err != nil {
        return nil, err
    }

    return &account, nil
}

func (s *SmartContract) Deposit(ctx contractapi.TransactionContextInterface, id string, amount int, timestamp string) error {
    account, err := s.QueryAccount(ctx, id)
    if err != nil {
        return err
    }

    account.Balance += amount

    transaction := Transaction{
        ID:        ctx.GetStub().GetTxID(),
        Type:      "deposit",
        Amount:    amount,
        Timestamp: timestamp,
    }
    account.Transactions = append(account.Transactions, transaction)

    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, accountJSON)
}

func (s *SmartContract) Withdraw(ctx contractapi.TransactionContextInterface, id string, amount int, timestamp string) error {
    account, err := s.QueryAccount(ctx, id)
    if err != nil {
        return err
    }

    if account.Balance < amount {
        return fmt.Errorf("insufficient funds")
    }

    account.Balance -= amount

    transaction := Transaction{
        ID:        ctx.GetStub().GetTxID(),
        Type:      "withdraw",
        Amount:    amount,
        Timestamp: timestamp,
    }
    account.Transactions = append(account.Transactions, transaction)

    accountJSON, err := json.Marshal(account)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, accountJSON)
}

func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, fromID string, toID string, amount int, timestamp string) error {
    fromAccount, err := s.QueryAccount(ctx, fromID)
    if err != nil {
        return err
    }

    if fromAccount.Balance < amount {
        return fmt.Errorf("insufficient funds")
    }

    toAccount, err := s.QueryAccount(ctx, toID)
    if err != nil {
        return err
    }

    fromAccount.Balance -= amount
    toAccount.Balance += amount

    transferID := ctx.GetStub().GetTxID()

    fromTransaction := Transaction{
        ID:        transferID,
        Type:      "transfer_out",
        Amount:    amount,
        Timestamp: timestamp,
    }
    toTransaction := Transaction{
        ID:        transferID,
        Type:      "transfer_in",
        Amount:    amount,
        Timestamp: timestamp,
    }

    fromAccount.Transactions = append(fromAccount.Transactions, fromTransaction)
    toAccount.Transactions = append(toAccount.Transactions, toTransaction)

    fromAccountJSON, err := json.Marshal(fromAccount)
    if err != nil {
        return err
    }

    toAccountJSON, err := json.Marshal(toAccount)
    if err != nil {
        return err
    }

    err = ctx.GetStub().PutState(fromID, fromAccountJSON)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(toID, toAccountJSON)
}

func (s *SmartContract) UserExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    accountJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state. %v", err)
    }
    return accountJSON != nil, nil
}

func (s *SmartContract) GetAllAccounts(ctx contractapi.TransactionContextInterface) ([]*Account, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var accounts []*Account
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var account Account
        err = json.Unmarshal(queryResponse.Value, &account)
        if err != nil {
            return nil, err
        }

        accounts = append(accounts, &account)
    }

    return accounts, nil
}

func (s *SmartContract) DeleteAccount(ctx contractapi.TransactionContextInterface, id string) error {
    return ctx.GetStub().DelState(id)
}

func (s *SmartContract) QueryTransactions(ctx contractapi.TransactionContextInterface, id string) ([]Transaction, error) {
    account, err := s.QueryAccount(ctx, id)
    if err != nil {
        return nil, err
    }
    return account.Transactions, nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error creating banking chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting banking chaincode: %s", err.Error())
        return
    }
}