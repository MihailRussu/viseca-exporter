package json

import (
	"encoding/json"

	"github.com/anothertobi/viseca-exporter/pkg/viseca"
)

// TransactionsString returns a JSON representation of the transactions.
func TransactionsString(transactions []viseca.Transaction) string {
	bytes, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return "{\"error\": \"Failed to marshal transactions\"}"
	}
	return string(bytes)
}

// TransactionString returns a JSON representation of a single transaction.
func TransactionString(transaction viseca.Transaction) string {
	bytes, err := json.MarshalIndent(transaction, "", "  ")
	if err != nil {
		return "{\"error\": \"Failed to marshal transaction\"}"
	}
	return string(bytes)
}