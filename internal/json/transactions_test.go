package json_test

import (
	"encoding/json"
	"testing"

	jsonpkg "github.com/anothertobi/viseca-exporter/internal/json"
	"github.com/anothertobi/viseca-exporter/pkg/viseca"
	"github.com/stretchr/testify/assert"
)

var inputTransaction = viseca.Transaction{
	TransactionID:    "AUTH8c919db2-1c23-43f1-8862-61c31336d9b6",
	CardID:           "0000000AAAAA0000",
	MaskedCardNumber: "XXXXXXXXXXXX0000",
	CardName:         "Mastercard",
	Date:             "2021-10-20T17:05:44",
	ShowTimestamp:    true,
	Amount:           50.55,
	Currency:         "CHF",
	OriginalAmount:   50.55,
	OriginalCurrency: "CHF",
	MerchantName:     "Aldi Suisse 00",
	PrettyName:       "ALDI",
	MerchantPlace:    "",
	IsOnline:         false,
	PFMCategory: viseca.PFMCategory{
		ID:                  "cv_groceries",
		Name:                "Lebensmittel",
		LightColor:          "#E2FDD3",
		MediumColor:         "#A5D58B",
		Color:               "#51A127",
		ImageURL:            "https://api.one.viseca.ch/v1/media/categories/icon_with_background/ic_cat_tile_groceries_v2.png",
		TransparentImageURL: "https://api.one.viseca.ch/v1/media/categories/icon_without_background/ic_cat_tile_groceries_v2.png",
	},
	StateType: "authorized",
	Details:   "Aldi Suisse 00",
	Type:      "merchant",
	IsBilled:  false,
	Links: viseca.TransactionLinks{
		Transactiondetails: "/v1/card/0000000AAAAA0000/transaction/AUTH8c919db2-1c23-43f1-8862-61c31336d9b6",
	},
}

func TestTransactionString(t *testing.T) {
	result := jsonpkg.TransactionString(inputTransaction)
	
	// Verify the result can be parsed back to JSON
	var parsedTransaction viseca.Transaction
	err := json.Unmarshal([]byte(result), &parsedTransaction)
	assert.NoError(t, err, "Should be valid JSON")
	
	// Verify fields are preserved correctly
	assert.Equal(t, inputTransaction.TransactionID, parsedTransaction.TransactionID)
	assert.Equal(t, inputTransaction.MerchantName, parsedTransaction.MerchantName)
	assert.Equal(t, inputTransaction.PrettyName, parsedTransaction.PrettyName)
	assert.Equal(t, inputTransaction.Amount, parsedTransaction.Amount)
	assert.Equal(t, inputTransaction.Currency, parsedTransaction.Currency)
	assert.Equal(t, inputTransaction.OriginalAmount, parsedTransaction.OriginalAmount)
	assert.Equal(t, inputTransaction.OriginalCurrency, parsedTransaction.OriginalCurrency)
}

func TestTransactionsString(t *testing.T) {
	inputTransactions := []viseca.Transaction{inputTransaction}
	result := jsonpkg.TransactionsString(inputTransactions)
	
	// Verify the result can be parsed back to JSON
	var parsedTransactions []viseca.Transaction
	err := json.Unmarshal([]byte(result), &parsedTransactions)
	assert.NoError(t, err, "Should be valid JSON")
	
	// Verify we have the right number of transactions
	assert.Equal(t, 1, len(parsedTransactions))
	
	// Verify fields are preserved correctly for the first transaction
	assert.Equal(t, inputTransaction.TransactionID, parsedTransactions[0].TransactionID)
	assert.Equal(t, inputTransaction.MerchantName, parsedTransactions[0].MerchantName)
	assert.Equal(t, inputTransaction.Amount, parsedTransactions[0].Amount)
}

func TestForeignCurrencyTransaction(t *testing.T) {
	foreignTransaction := viseca.Transaction{
		TransactionID:     "TRX2025051200004466612",
		Date:              "2025-05-12T09:01:20+02:00",
		MerchantName:      "CLAUDE.AI SUBSCRIPTION",
		PrettyName:        "Claude.ai",
		IsOnline:          true,
		Amount:            17.15,
		Currency:          "CHF",
		OriginalAmount:    20.00,
		OriginalCurrency:  "USD",
		PFMCategory: viseca.PFMCategory{
			ID:   "entertainment_and_leisure",
			Name: "Entertainment & Leisure",
		},
		StateType: "booked",
		Details:   "CLAUDE.AI SUBSCRIPTION",
		Type:      "merchant",
		IsBilled:  true,
	}

	result := jsonpkg.TransactionString(foreignTransaction)
	
	// Verify the result can be parsed back to JSON
	var parsedTransaction viseca.Transaction
	err := json.Unmarshal([]byte(result), &parsedTransaction)
	assert.NoError(t, err, "Should be valid JSON")
	
	// Verify the currency fields are preserved correctly
	assert.Equal(t, "CHF", parsedTransaction.Currency)
	assert.Equal(t, "USD", parsedTransaction.OriginalCurrency)
	assert.Equal(t, 17.15, parsedTransaction.Amount)
	assert.Equal(t, 20.00, parsedTransaction.OriginalAmount)
}