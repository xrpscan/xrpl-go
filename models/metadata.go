package models

type CreatedNode struct {
	CreatedNode CreatedNodeMap
}

type CreatedNodeMap struct {
	LedgerEntryType string
	LedgerIndex     string
	NewFields       map[string]interface{}
}

type ModifiedNode struct {
	ModifiedNode ModifiedNodeMap
}

type ModifiedNodeMap struct {
	LedgerEntryType   string
	LedgerIndex       string
	FinalFields       map[string]interface{}
	PreviousFields    map[string]interface{}
	PreviousTxnID     string
	PreviousTxnLgrSeq int64
}

type DeletedNode struct {
	DeletedNode DeletedNodeMap
}

type DeletedNodeMap struct {
	LedgerEntryType string
	LedgerIndex     string
	FinalFields     map[string]interface{}
}

type TransactionMetadata struct {
	AffectedNodes     []map[string]interface{}
	DeliveredAmount   Amount
	Delivered_Amount  Amount `json:"delivered_amount,omitempty"`
	TransactionIndex  int64
	TransactionResult string
}
