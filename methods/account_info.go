package methods

import "github.com/xrpscan/xrpl-go/models"

type AccountInfoRequest struct {
	models.BaseRequest
	Account     string `json:"account,omitempty"`
	LedgerHash  string `json:"ledger_hash,omitempty"`
	LedgerIndex string `json:"ledger_index,omitempty"`
	Queue       bool   `json:"queue,omitempty"`
	SignerLists bool   `json:"signer_lists,omitempty"`
	Strict      bool   `json:"strict,omitempty"`
}

type QueueTransaction struct {
	AuthChange    bool   `json:"auth_change,omitempty"`
	Fee           string `json:"fee,omitempty"`
	FeeLevel      string `json:"fee_level,omitempty"`
	MaxSpendDrops string `json:"max_spend_drops,omitempty"`
	Seq           int    `json:"seq,omitempty"`
}

type QueueData struct {
	TxnCount           int                `json:"txn_count,omitempty"`
	AuthChangeQueued   bool               `json:"auth_change_queued,omitempty"`
	LowestSequence     int                `json:"lowest_sequence,omitempty"`
	HighestSequence    int                `json:"highest_sequence,omitempty"`
	MaxSpendDropsTotal string             `json:"max_spend_drops_total,omitempty"`
	Transactions       []QueueTransaction `json:"transactions,omitempty"`
}

type AccountInfoResponse struct {
	models.BaseResponse
	Result AccountInfoResult `json:"result,omitempty"`
}

type AccountInfoResult struct {
	AccountData        string    `json:"account_data,omitempty"` //todo
	SignerLists        string    `json:"signer_lists,omitempty"` //todo
	LedgerCurrentIndex int       `json:"ledger_current_index,omitempty"`
	LedgerIndex        int       `json:"ledger_index,omitempty"`
	QueueData          QueueData `json:"queue_data,omitempty"`
	Validated          bool      `json:"validated,omitempty"`
}
