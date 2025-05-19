package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/anothertobi/viseca-exporter/internal/csv"
	"github.com/anothertobi/viseca-exporter/internal/json"
	"github.com/anothertobi/viseca-exporter/pkg/viseca"
)

const sessionCookieName = "AL_SESS-S"

// arg0: cardID
// arg1: sessionCookie (e.g. `AL_SESS-S=...`)
// arg2: [optional] output format (json or csv, defaults to csv)
func main() {
	if len(os.Args) < 3 {
		log.Fatal("card ID and session cookie args required")
	}
	visecaClient, err := initClient(os.Args[2])
	if err != nil {
		log.Fatalf("error initializing Viseca API client: %v", err)
	}

	ctx := context.Background()

	transactions, err := visecaClient.ListAllTransactions(ctx, os.Args[1])
	if err != nil {
		log.Fatalf("error listing all transactions: %v", err)
	}
	
	// Determine output format
	outputFormat := "csv" // Default format
	if len(os.Args) >= 4 {
		if os.Args[3] == "json" {
			outputFormat = "json"
		}
	}
	
	// Output in the selected format
	if outputFormat == "json" {
		fmt.Println(json.TransactionsString(transactions))
	} else {
		fmt.Println(csv.TransactionsString(transactions))
	}
}

func initClient(sessionCookie string) (*viseca.Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := http.Client{
		Jar: cookieJar,
	}
	visecaClient := viseca.NewClient(&httpClient)
	cookie := &http.Cookie{
		Name:  sessionCookieName,
		Value: extractSessionCookieValue(sessionCookie),
	}
	httpClient.Jar.SetCookies(visecaClient.BaseURL, []*http.Cookie{cookie})

	return visecaClient, nil
}

func extractSessionCookieValue(sessionCookie string) string {
	return strings.TrimPrefix(sessionCookie, fmt.Sprintf("%s=", sessionCookieName))
}
