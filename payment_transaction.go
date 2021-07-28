package goshopify

import (
	"fmt"
	"net/http"
	"time"
)

const paymentTransactionBasePath = "shopify_payments/balance/transactions"
const paymentTransactionResourceName = "transactions"

// linkRegex is used to extract pagination links from product search results.
//var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

type PaymentTransactionService interface {
	List(interface{}) ([]PaymentTransaction, error)
	ListWithPagination(interface{}) ([]PaymentTransaction, *Pagination, error)
	Get(int64, interface{}) (PaymentTransaction, error)
}

type PaymentTransactionServiceOp struct {
	client *Client
}

type PaymentTransaction struct {
	ID							int64		`json:"id,omitempty"`
	Type						string		`json:"type,omitempty"`
	Test						bool		`json:"test,omitempty"`
	PayoutID					int64		`json:"payout_id,omitempty"`
	PayoutStatus				string		`json:"payout_status,omitempty"`
	Currency					string		`json:"currency,omitempty"`
	Amount 						string		`json:"amount,omitempty"`
	Fee							int64		`json:"fee,omitempty"`
	Net							int64		`json:"net,omitempty"`
	SourceID					int64		`json:"source_id,omitempty"`
	SourceType					string		`json:"source_type,omitempty"`
	SourceOrderID				int64		`json:"source_order_id.omitempty"`
	SourceOrderTransactionID	int64		`json:"source_order_transaction_id.omitempty"`
	ProcessedAt					*time.Time	`json:"processed_at,omitempty"`
}

type PaymentTransactionListOptions struct {
	ListOptions
	LastID			int64	`url:"last_id,omitempty"`
	Test			bool	`url:"test,omitempty"`
	PayoutID		int64	`url:"payout_id,omitempty"`
	PayoutStatus	int64	`url:"payout_status,omitempty"`
}

// Represents the result from the paymenttransactions.json endpoint
type PaymentTransactionsResource struct {
	PaymentTransactions []PaymentTransaction `json:"transactions"`
}

type PaymentTransactionResource struct {
	PaymentTransaction PaymentTransaction `json:"transaction"`
}

// List products
func (s *PaymentTransactionServiceOp) List(options interface{}) ([]PaymentTransaction, error) {
	transactions, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (s *PaymentTransactionServiceOp) ListWithPagination(options interface{}) ([]PaymentTransaction, *Pagination, error) {
	path := fmt.Sprintf("%s.json", paymentTransactionBasePath)
	resource := new(PaymentTransactionsResource)
	headers := http.Header{}

	headers, err := s.client.createAndDoGetHeaders("GET", path, nil, options, resource)
	if err != nil {
		return nil, nil, err
	}

	// Extract pagination info from header
	linkHeader := headers.Get("Link")

	pagination, err := extractPagination(linkHeader)
	if err != nil {
		return nil, nil, err
	}

	return resource.PaymentTransactions, pagination, nil
}

// Get individual product
func (s *PaymentTransactionServiceOp) Get(transactionID int64, options interface{}) (PaymentTransaction, error) {
	path := fmt.Sprintf("%s/%d.json", paymentTransactionBasePath, transactionID)
	resource := new(PaymentTransactionResource)
	err := s.client.Get(path, resource, options)
	return resource.PaymentTransaction, err
}
