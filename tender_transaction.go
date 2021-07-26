package goshopify

import (
	"fmt"
	"net/http"
	"time"
)

const tenderTransactionBasePath = "tender_transactions"
const tenderTransactionResourceName = "tender_transactions"

// linkRegex is used to extract pagination links from product search results.
//var linkRegex = regexp.MustCompile(`^ *<([^>]+)>; rel="(previous|next)" *$`)

type TenderTransactionService interface {
	List(interface{}) ([]TenderTransaction, error)
	ListWithPagination(interface{}) ([]TenderTransaction, *Pagination, error)
	Get(int64, interface{}) (TenderTransaction, error)
}

type TenderTransactionServiceOp struct {
	client *Client
}

type TenderTransaction struct {
	ID							int64		`json:"id,omitempty"`
	OrderID						int64		`json:"order_id,omitempty"`
	Amount						int64		`json:"amount"`
	Currency					string		`json:"currency,omitempty"`
	UserID						*int64		`json:"user_id,omitempty"`
	Test						bool		`json:"test,omitempty"`
	ProcessedAt					*time.Time	`json:"processed_at,omitempty"`
	RemoteReference				string		`json:"remote_efernece,omitempty"`
	PaymentDetails				*TenderTransactionPaymentDetails	`json:"payment_details,omitempty"`
	PaymentMethod				string		`json:"payment_method"`
}

type TenderTransactionPaymentDetails struct {
	CreditCardNumber	string	`json:"credit_card_number,omitempty"`
	CreditCardCompany	string	`json:"credit_card_company,omitempty"`
}

type TenderTransactionListOptions struct {
	ListOptions
	SinceID			int64		`url:"since_id,omitempty"`
	ProcessedAtMin	time.Time	`url:"processed_at_min,omitempty"`
	ProcessedAtMax	time.Time	`url:"processed_at_max,omitempty"`
	Processed		time.Time	`url:"processed_at,omitempty"`
	Order			string		`url:"order,omitempty"`
}

// Represents the result from the tender_transactions.json endpoint
type TenderTransactionsResource struct {
	TenderTransactions []TenderTransaction `json:"tender_transactions"`
}

type TenderTransactionResource struct {
	TenderTransaction TenderTransaction `json:"tender_transaction"`
}

// List products
func (s *TenderTransactionServiceOp) List(options interface{}) ([]TenderTransaction, error) {
	transactions, _, err := s.ListWithPagination(options)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

// ListWithPagination lists products and return pagination to retrieve next/previous results.
func (s *TenderTransactionServiceOp) ListWithPagination(options interface{}) ([]TenderTransaction, *Pagination, error) {
	path := fmt.Sprintf("%s.json", tenderTransactionBasePath)
	resource := new(TenderTransactionsResource)
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

	return resource.TenderTransactions, pagination, nil
}

// Get individual product
func (s *TenderTransactionServiceOp) Get(transactionID int64, options interface{}) (TenderTransaction, error) {
	path := fmt.Sprintf("%s/%d.json", tenderTransactionBasePath, transactionID)
	resource := new(TenderTransactionResource)
	err := s.client.Get(path, resource, options)
	return resource.TenderTransaction, err
}
