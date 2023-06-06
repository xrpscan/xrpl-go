package models

type LedgerStream struct {
	Type             string `json:"type,omitempty"` // default: ledgerClosed
	FeeBase          uint64 `json:"fee_base,omitempty"`
	FeeRef           uint64 `json:"fee_ref,omitempty"`
	LedgerHash       string `json:"ledger_hash,omitempty"`
	LedgerIndex      uint64 `json:"ledger_index,omitempty"`
	LedgerTime       uint64 `json:"ledger_time,omitempty"`
	ReserveBase      uint64 `json:"reserve_base,omitempty"`
	ReserveInc       uint64 `json:"reserve_inc,omitempty"`
	TxnCount         uint64 `json:"txn_count,omitempty"`
	ValidatedLedgers string `json:"validated_ledgers,omitempty"`
}

type ValidationStream struct {
	Type                string   `json:"type,omitempty"` // default: validationReceived
	Amendments          []string `json:"amendments,omitempty"`
	BaseFee             uint64   `json:"base_fee,omitempty"`
	Cookie              string   `json:"cookie,omitempty"`
	Data                string   `json:"data,omitempty"`
	Flags               uint64   `json:"flags,omitempty"`
	Full                bool     `json:"full,omitempty"`
	LedgerHash          string   `json:"ledger_hash,omitempty"`
	LedgerIndex         uint64   `json:"ledger_index,omitempty"`
	LoadFee             uint64   `json:"load_fee,omitempty"`
	MasterKey           string   `json:"master_key,omitempty"`
	ReserveBase         uint64   `json:"reserve_base,omitempty"`
	ReserveInc          uint64   `json:"reserve_inc,omitempty"`
	Signature           string   `json:"signature,omitempty"`
	SigningTime         uint64   `json:"signing_time,omitempty"`
	ValidationPublicKey string   `json:"validation_public_key,omitempty"`
}

type TransactionStream struct {
	Type                string `json:"type,omitempty"` // default: transaction
	Status              string `json:"status,omitempty"`
	EngineResult        string `json:"engine_result,omitempty"`
	EngineResultCode    int64  `json:"engine_result_code,omitempty"`
	EngineResultMessage string `json:"engine_result_message,omitempty"`
	LedgerCurrentIndex  uint64 `json:"ledger_current_index,omitempty"`
	LedgerHash          string `json:"ledger_hash,omitempty"`
	LedgerIndex         uint64 `json:"ledger_index,omitempty"`
	Meta                string `json:"meta,omitempty"`
	Transaction         string `json:"transaction,omitempty"`
	Validated           bool   `json:"validated,omitempty"`
	// Meta                TransactionMeta `json:"meta,omitempty"`
	// Transaction         Transaction     `json:"transaction,omitempty"`
}

type PeerStatusStream struct {
	Type           string `json:"type,omitempty"` // default: peerStatusChange
	Action         string `json:"action,omitempty"`
	Date           uint64 `json:"date,omitempty"`
	LedgerHash     string `json:"ledger_hash,omitempty"`
	LedgerIndex    uint64 `json:"ledger_index,omitempty"`
	LedgerIndexMax uint64 `json:"ledger_index_max,omitempty"`
	LedgerIndexMin uint64 `json:"ledger_index_min,omitempty"`
}

type OrderBookStream struct {
	Type                string `json:"type,omitempty"` // default: transaction
	Status              string `json:"status,omitempty"`
	EngineResult        string `json:"engine_result,omitempty"`
	EngineResultCode    int64  `json:"engine_result_code,omitempty"`
	EngineResultMessage string `json:"engine_result_message,omitempty"`
	LedgerCurrentIndex  uint64 `json:"ledger_current_index,omitempty"`
	LedgerHash          string `json:"ledger_hash,omitempty"`
	LedgerIndex         uint64 `json:"ledger_index,omitempty"`
	Meta                string `json:"meta,omitempty"`
	Transaction         string `json:"transaction,omitempty"`
	Validated           bool   `json:"validated,omitempty"`
}

type ConsensusStream struct {
	Type      string `json:"type,omitempty"` // default: consensusPhase
	Consensus string `json:"consensus,omitempty"`
}

type PathFindStream struct {
	Type               string `json:"type,omitempty"` // default: path_find
	SourceAccount      string `json:"source_account,omitempty"`
	DestinationAccount string `json:"destination_account,omitempty"`
	DestinationAmount  string `json:"destination_amount,omitempty"`
	FullReply          bool   `json:"full_reply,omitempty"`
	Id                 string `json:"id,omitempty"`
	SendMax            string `json:"send_max,omitempty"`
}
