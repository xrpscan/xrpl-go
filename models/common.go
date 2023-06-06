package models

type LedgerIndex int

type Currency struct {
	Currency string `json:"currency,omitempty"`
}

type IssuedCurrency struct {
	Currency `json:"currency,omitempty"`
	Issuer   string `json:"issuer,omitempty"`
}

type IssuedCurrencyAmount struct {
	IssuedCurrency `json:"issued_currency,omitempty"`
	Value          string `json:"value,omitempty"`
}

type Amount IssuedCurrencyAmount

type Signer struct {
	Signer SignerMap `json:"signer,omitempty"`
}

type SignerMap struct {
	Account       string `json:"account,omitempty"`
	TxnSignature  string `json:"txn_signature,omitempty"`
	SigningPubKey string `json:"signing_pub_key,omitempty"`
}

type Memo struct {
	Memo MemoMap `json:"memo,omitempty"`
}

type MemoMap struct {
	MemoData   string `json:"memo_data,omitempty"`
	MemoType   string `json:"memo_type,omitempty"`
	MemoFormat string `json:"memo_format,omitempty"`
}

type StreamType string

type PathStep struct {
	Account  string `json:"account,omitempty"`
	Currency string `json:"currency,omitempty"`
	Issuer   string `json:"issuer,omitempty"`
}

type Path []PathStep

type SignerEntry struct {
	SignerEntry SignerEntryMap `json:"signer_entry,omitempty"`
}

type SignerEntryMap struct {
	Account       string `json:"account,omitempty"`
	SignerWeight  int16  `json:"signer_weight,omitempty"`
	WalletLocator string `json:"wallet_locator,omitempty"`
}

type ResponseOnlyTxInfo struct {
	Date        int    `json:"date,omitempty"`
	Hash        string `json:"hash,omitempty"`
	LedgerIndex int    `json:"ledger_index,omitempty"`
	InLedger    int    `json:"in_ledger,omitempty"`
}

type NFTOffer struct {
	Amount        Amount `json:"amount,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	NftOfferIndex string `json:"nft_offer_index,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Destination   string `json:"destination,omitempty"`
	Expiration    int    `json:"expiration,omitempty"`
}

type GlobalFlags struct {
}
