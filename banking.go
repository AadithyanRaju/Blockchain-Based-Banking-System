package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Account represents a user account in the banking system
type Account struct {
	UserID       string  `json:"userId"`
	Name         string  `json:"name"`
	AadhaarHash  string  `json:"aadhaarHash"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"passwordHash"`
	PhoneNumber  string  `json:"phoneNumber"`
	Status       string  `json:"status"` // Active, Suspended, Closed
	Role         string  `json:"role"`   // Customer, Admin
	Balance      float64 `json:"balance"`
}

// Transaction represents a transaction record
type Transaction struct {
	UserID         string    `json:"userId"`
	ReferenceNumber string    `json:"referenceNumber"`
	Type           string    `json:"type"` // Credit or Debit
	Amount         float64   `json:"amount"`
	Timestamp      time.Time `json:"timestamp"`
}

// SmartContract provides functions for managing accounts and transactions
type SmartContract struct {
	contractapi.Contract
}

// Prefixes for state keys
const (
	AccountPrefix     = "USER_"
	TransactionPrefix = "TRANSACTION_"
)

// CreateAccount adds a new account to the ledger
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, userID, name, aadhaarHash, email, passwordHash, phoneNumber, role string) error {
	accountKey := AccountPrefix + userID

	existingAccount, err := ctx.GetStub().GetState(accountKey)
	if err != nil {
		return fmt.Errorf("failed to check if account exists: %v", err)
	}
	if existingAccount != nil {
		return fmt.Errorf("account already exists")
	}

	account := Account{
		UserID:       userID,
		Name:         name,
		AadhaarHash:  aadhaarHash,
		Email:        email,
		PasswordHash: passwordHash,
		PhoneNumber:  phoneNumber,
		Status:       "Active",
		Role:         role,
		Balance:      0.0,
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %v", err)
	}

	return ctx.GetStub().PutState(accountKey, accountJSON)
}

// GetAccount retrieves an account by its userID
func (s *SmartContract) GetAccount(ctx contractapi.TransactionContextInterface, userID string) (*Account, error) {
	accountKey := AccountPrefix + userID

	accountJSON, err := ctx.GetStub().GetState(accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read account: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("account does not exist")
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %v", err)
	}

	return &account, nil
}

// CreateTransaction records a new transaction and updates the account balance
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, userID, referenceNumber, transactionType string, amount float64) error {
	account, err := s.GetAccount(ctx, userID)
	if err != nil {
		return err
	}

	if account.Status != "Active" {
		return fmt.Errorf("account is not active")
	}

	if transactionType == "Debit" && account.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	transaction := Transaction{
		UserID:         userID,
		ReferenceNumber: referenceNumber,
		Type:           transactionType,
		Amount:         amount,
		Timestamp:      time.Now(),
	}

	if transactionType == "Credit" {
		account.Balance += amount
	} else if transactionType == "Debit" {
		account.Balance -= amount
	} else {
		return fmt.Errorf("invalid transaction type")
	}

	accountKey := AccountPrefix + userID
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal updated account: %v", err)
	}

	transactionKey := fmt.Sprintf("%s%s_%s", TransactionPrefix, userID, referenceNumber)
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	err = ctx.GetStub().PutState(accountKey, accountJSON)
	if err != nil {
		return fmt.Errorf("failed to update account: %v", err)
	}

	return ctx.GetStub().PutState(transactionKey, transactionJSON)
}

// GetTransaction retrieves a transaction by its key
func (s *SmartContract) GetTransaction(ctx contractapi.TransactionContextInterface, userID, referenceNumber string) (*Transaction, error) {
	transactionKey := fmt.Sprintf("%s%s_%s", TransactionPrefix, userID, referenceNumber)

	transactionJSON, err := ctx.GetStub().GetState(transactionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction: %v", err)
	}
	if transactionJSON == nil {
		return nil, fmt.Errorf("transaction does not exist")
	}

	var transaction Transaction
	err = json.Unmarshal(transactionJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %v", err)
	}

	return &transaction, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating banking chaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting banking chaincode: %v\n", err)
	}
}
