package models

type LedgerRequest struct {
	BaseRequest
	LedgerHash   string      `json:"ledger_hash,omitempty"`
	LedgerIndex  LedgerIndex `json:"ledger_index,omitempty"`
	Full         bool        `json:"full,omitempty"`
	Accounts     bool        `json:"accounts,omitempty"`
	Transactions bool        `json:"transactions,omitempty"`
	Expand       bool        `json:"expand,omitempty"`
	OwnerFunds   bool        `json:"owner_funds,omitempty"`
	Binary       bool        `json:"binary,omitempty"`
	Queue        bool        `json:"queue,omitempty"`
}

// type ModifiedMetadata struct {
// 	TransactionMetadata
// 	OwnerFunds string
// }

// type ModifiedOfferCreateTransaction struct {
// 	Transaction Transaction
// 	Metadata    ModifiedMetadata
// }

type LedgerQueueData struct {
	Account          string      `json:"account,omitempty"`
	Tx               interface{} `json:"tx,omitempty"`
	RetriesRemaining int         `json:"retries_remaining,omitempty"`
	PreflightResult  string      `json:"preflight_result,omitempty"`
	LastResult       string      `json:"last_result,omitempty"`
	AuthChange       bool        `json:"auth_change,omitempty"`
	Fee              string      `json:"fee,omitempty"`
	FeeLevel         string      `json:"fee_level,omitempty"`
	MaxSpendDrops    string      `json:"max_spend_drops,omitempty"`
}

type BinaryLedger struct {
	AccountState []string `json:"accountState,omitempty"`
	Transactions []string `json:"transactions,omitempty"`
}

type LedgerResponse struct {
	BaseResponse
	Result LedgerResult
}

type LedgerResult struct {
	Ledger      interface{}       `json:"ledger,omitempty"`
	LedgerHash  string            `json:"ledger_hash,omitempty"`
	LedgerIndex int               `json:"ledger_index,omitempty"`
	QueueData   []LedgerQueueData `json:"queue_data,omitempty"`
	Validated   bool              `json:"validated,omitempty"`
}
