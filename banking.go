package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Account struct for storing user account details
type Account struct {
	UserID       string  `json:"userID"`
	Name         string  `json:"name"`
	AadhaarHash  string  `json:"aadhaarHash"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"passwordHash"`
	PhoneNumber  string  `json:"phoneNumber"`
	Role         string  `json:"role"`
	Balance      float64 `json:"balance"`
}

// Transaction struct for storing individual transactions
type Transaction struct {
	UserID         string  `json:"userID"`
	ReferenceNumber string  `json:"referenceNumber"`
	Type           string  `json:"type"`
	Amount         float64 `json:"amount"`
	Timestamp      string  `json:"timestamp"`
}

// Transfer struct for storing transfer transactions
type Transfer struct {
	SenderID        string  `json:"senderID"`
	ReceiverID      string  `json:"receiverID"`
	ReferenceNumber string  `json:"referenceNumber"`
	Amount          float64 `json:"amount"`
	Timestamp       string  `json:"timestamp"`
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

// CreateAccount creates a new account in the ledger
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, userID, name, aadhaarHash, email, passwordHash, phoneNumber, role string, balance float64) error {
	accountKey := AccountPrefix + userID

	existingAccount, err := ctx.GetStub().GetState(accountKey)
	if err != nil {
		return fmt.Errorf("failed to check for existing account: %v", err)
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
		Role:         role,
		Balance:      balance,
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %v", err)
	}

	return ctx.GetStub().PutState(accountKey, accountJSON)
}

// GetAccount retrieves an account by user ID
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

// TransferFunds transfers an amount from one account to another
func (s *SmartContract) TransferFunds(ctx contractapi.TransactionContextInterface, senderID, receiverID, amountStr, referenceNumber string) error {
	// Parse amount
	amount, err := parseAmount(amountStr)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	// Fetch sender account
	senderAccount, err := s.GetAccount(ctx, senderID)
	if err != nil {
		return fmt.Errorf("failed to get sender account: %v", err)
	}
	if senderAccount.Balance < amount {
		return fmt.Errorf("insufficient funds in sender account")
	}

	// Fetch receiver account
	receiverAccount, err := s.GetAccount(ctx, receiverID)
	if err != nil {
		return fmt.Errorf("failed to get receiver account: %v", err)
	}

	// Perform transfer
	senderAccount.Balance -= amount
	receiverAccount.Balance += amount

	// Update accounts in the ledger
	if err := s.updateAccount(ctx, senderAccount); err != nil {
		return fmt.Errorf("failed to update sender account: %v", err)
	}
	if err := s.updateAccount(ctx, receiverAccount); err != nil {
		return fmt.Errorf("failed to update receiver account: %v", err)
	}

	// Record transactions
	timestamp := time.Now().String()
	senderTransaction := Transaction{
		UserID:          senderID,
		ReferenceNumber: referenceNumber,
		Type:            "Debit",
		Amount:          amount,
		Timestamp:       timestamp,
	}
	receiverTransaction := Transaction{
		UserID:          receiverID,
		ReferenceNumber: referenceNumber,
		Type:            "Credit",
		Amount:          amount,
		Timestamp:       timestamp,
	}

	if err := s.recordTransaction(ctx, senderTransaction, senderID, receiverID); err != nil {
		return fmt.Errorf("failed to record sender transaction: %v", err)
	}
	if err := s.recordTransaction(ctx, receiverTransaction, senderID, receiverID); err != nil {
		return fmt.Errorf("failed to record receiver transaction: %v", err)
	}

	return nil
}

// Helper methods

func (s *SmartContract) updateAccount(ctx contractapi.TransactionContextInterface, account *Account) error {
	accountKey := AccountPrefix + account.UserID
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %v", err)
	}

	return ctx.GetStub().PutState(accountKey, accountJSON)
}

func (s *SmartContract) recordTransaction(ctx contractapi.TransactionContextInterface, transaction Transaction, senderID, receiverID string) error {
	transactionKey := fmt.Sprintf("%s%s_%s_%s", TransactionPrefix, senderID, receiverID, transaction.ReferenceNumber)
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %v", err)
	}

	return ctx.GetStub().PutState(transactionKey, transactionJSON)
}

func parseAmount(amountStr string) (float64, error) {
	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)
	return amount, err
}

// Deposit funds into the account
func (s *SmartContract) Deposit(ctx contractapi.TransactionContextInterface, userID string, amount string, referenceNumber string) error {
    accountKey := "USER_" + userID
    accountBytes, err := ctx.GetStub().GetState(accountKey)
    if err != nil || accountBytes == nil {
        return fmt.Errorf("account not found")
    }

    var account Account
    json.Unmarshal(accountBytes, &account)
    depositAmount, _ := strconv.ParseFloat(amount, 64)
    account.Balance += depositAmount

    accountBytes, _ = json.Marshal(account)
    ctx.GetStub().PutState(accountKey, accountBytes)

    transaction := Transaction{
        UserID:         userID,
        ReferenceNumber: referenceNumber,
        Type:           "Deposit",
        Amount:         depositAmount,
		Timestamp:      time.Now().String(),
    }
    transactionKey := fmt.Sprintf("TRANSACTION_%s_%s", userID, referenceNumber)
    transactionBytes, _ := json.Marshal(transaction)
    ctx.GetStub().PutState(transactionKey, transactionBytes)

    return nil
}

// Withdraw funds from the account
func (s *SmartContract) Withdraw(ctx contractapi.TransactionContextInterface, userID string, amount string, referenceNumber string) error {
    accountKey := "USER_" + userID
    accountBytes, err := ctx.GetStub().GetState(accountKey)
    if err != nil || accountBytes == nil {
        return fmt.Errorf("account not found")
    }

    var account Account
    json.Unmarshal(accountBytes, &account)
    withdrawAmount, _ := strconv.ParseFloat(amount, 64)
    if account.Balance < withdrawAmount {
        return fmt.Errorf("insufficient balance")
    }
    account.Balance -= withdrawAmount

    accountBytes, _ = json.Marshal(account)
    ctx.GetStub().PutState(accountKey, accountBytes)

    transaction := Transaction{
        UserID:         userID,
        ReferenceNumber: referenceNumber,
        Type:           "Withdraw",
        Amount:         withdrawAmount,
        Timestamp:      time.Now().String(),
    }
    transactionKey := fmt.Sprintf("TRANSACTION_%s_%s", userID, referenceNumber)
    transactionBytes, _ := json.Marshal(transaction)
    ctx.GetStub().PutState(transactionKey, transactionBytes)

    return nil
}

// CreateTransfer creates a new transfer with states for sender and receiver
func (s *SmartContract) CreateTransfer(ctx contractapi.TransactionContextInterface, senderID string, receiverID string, amount string, referenceNumber string) error {
	// Parse amount
	transferAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return fmt.Errorf("invalid amount: %v", err)
	}

	// Fetch sender account
	senderKey := "USER_" + senderID
	senderBytes, err := ctx.GetStub().GetState(senderKey)
	if err != nil || senderBytes == nil {
		return fmt.Errorf("sender account not found")
	}

	// Fetch receiver account
	receiverKey := "USER_" + receiverID
	receiverBytes, err := ctx.GetStub().GetState(receiverKey)
	if err != nil || receiverBytes == nil {
		return fmt.Errorf("receiver account not found")
	}

	// Unmarshal accounts
	var sender Account
	var receiver Account
	json.Unmarshal(senderBytes, &sender)
	json.Unmarshal(receiverBytes, &receiver)

	// Check sender balance
	if sender.Balance < transferAmount {
		return fmt.Errorf("insufficient balance for transfer")
	}

	// Update balances
	sender.Balance -= transferAmount
	receiver.Balance += transferAmount

	// Save updated accounts
	senderBytes, _ = json.Marshal(sender)
	ctx.GetStub().PutState(senderKey, senderBytes)

	receiverBytes, _ = json.Marshal(receiver)
	ctx.GetStub().PutState(receiverKey, receiverBytes)

	// Create transfer record
	timestamp := time.Now().String()
	transfer := Transfer{
		SenderID:        senderID,
		ReceiverID:      receiverID,
		ReferenceNumber: referenceNumber,
		Amount:          transferAmount,
		Timestamp:       timestamp,
	}

	// Save transfer records for sender and receiver
	senderTransferKey := fmt.Sprintf("TRANSACTION_TRANSFER_%s_%s_%s", senderID, receiverID, referenceNumber)
	receiverTransferKey := fmt.Sprintf("TRANSACTION_TRANSFER_%s_%s_%s", receiverID, senderID, referenceNumber)

	transferBytes, _ := json.Marshal(transfer)
	ctx.GetStub().PutState(senderTransferKey, transferBytes)
	ctx.GetStub().PutState(receiverTransferKey, transferBytes)

	return nil
}

// GetTransfer fetches a transfer transaction by key
func (s *SmartContract) GetTransfer(ctx contractapi.TransactionContextInterface, senderID string, receiverID string, referenceNumber string) (*Transfer, error) {
	transferKey := fmt.Sprintf("TRANSACTION_TRANSFER_%s_%s_%s", senderID, receiverID, referenceNumber)
	transferBytes, err := ctx.GetStub().GetState(transferKey)
	if err != nil || transferBytes == nil {
		return nil, fmt.Errorf("transfer not found")
	}

	var transfer Transfer
	err = json.Unmarshal(transferBytes, &transfer)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling transfer: %v", err)
	}

	return &transfer, nil
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
