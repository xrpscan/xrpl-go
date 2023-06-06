package models

// Set of common fields for every transaction
type BaseTransaction struct {
	Account            string
	TransactionType    string
	Fee                string
	Sequence           int64
	AccountTxnID       string
	Flags              int64
	LastLedgerSequence int64
	Memos              []Memo
	Signers            []Signer
	SourceTag          int64
	SigningPubKey      string
	TicketSequence     int64
	TxnSignature       string
}

// A Payment transaction represents a transfer of value from one account to
// another.
//
// TransactionType: 'Payment'
type TransactionPayment struct {
	BaseTransaction
	Amount         Amount
	Destination    string
	DestinationTag int64
	InvoiceID      string
	Paths          []Path
	SendMax        Amount
	DeliverMin     Amount
	Flags          int64
}

type PaymentFlags struct {
	GlobalFlags
	TfNoDirectRipple bool `json:"tfNoDirectRipple,omitempty"`
	TfPartialPayment bool `json:"tfPartialPayment,omitempty"`
	TfLimitQuality   bool `json:"tfLimitQuality,omitempty"`
}

// The NFTokenOfferAccept transaction is used to accept offers to buy or sell
// an NFToken
//
// TransactionType: 'NFTokenAcceptOffer'
type TransactionNFTokenAcceptOffer struct {
	BaseTransaction
	NFTokenSellOffer string
	NFTokenBuyOffer  string
	NFTokenBrokerFee Amount
}

// The NFTokenBurn transaction is used to remove an NFToken object from the
// NFTokenPage in which it is being held, effectively removing the token from
// the ledger ("burning" it).
//
// TransactionType: 'NFTokenBurn'
type TransactionNFTokenBurn struct {
	BaseTransaction
	Account   string
	NFTokenID string
	Owner     string
}

// The NFTokenCancelOffer transaction deletes existing NFTokenOffer objects.
// It is useful if you want to free up space on your account to lower your
// reserve requirement.
//
// TransactionType: 'NFTokenCancelOffer'
type TransactionNFTokenCancelOffer struct {
	BaseTransaction
	NFTokenOffers []string
}

// The NFTokenCreateOffer transaction creates either an offer to buy an
// NFT the submitting account does not own, or an offer to sell an NFT
// the submitting account does own.
//
// TransactionType: 'NFTokenCreateOffer'
type TransactionNFTokenCreateOffer struct {
	BaseTransaction
	NFTokenID   string
	Amount      Amount
	Owner       string
	Expiration  int64
	Destination string
	Flags       int64
}

type NFTokenCreateOfferFlags struct {
	GlobalFlags
	TfSellNFToken bool `json:"tfSellNFToken,omitempty"`
}

// The NFTokenMint transaction creates an NFToken object and adds it to the
// relevant NFTokenPage object of the minter. If the transaction is
// successful, the newly minted token will be owned by the minter account
// specified by the transaction.
//
// TransactionType: 'NFTokenMint'
type TransactionNFTokenMint struct {
	BaseTransaction
	NFTokenTaxon int64
	Issuer       string
	TransferFee  int64
	URI          string
	Flags        int64
}

type NFTokenMintFlags struct {
	GlobalFlags
	TfBurnable     bool `json:"tfBurnable,omitempty"`
	TfOnlyXRP      bool `json:"tfOnlyXRP,omitempty"`
	TfTrustLine    bool `json:"tfTrustLine,omitempty"`
	TfTransferable bool `json:"tfTransferable,omitempty"`
}

// An AccountDelete transaction deletes an account and any objects it owns in
// the XRP Ledger, if possible, sending the account's remaining XRP to a
// specified destination account.
//
// TransactionType: 'AccountDelete'
type TransactionAccountDelete struct {
	BaseTransaction
	Destination    string
	DestinationTag int64
}

// Map of flags to boolean values representing {@link AccountSet} transaction
// flags.
type AccountSetFlags struct {
	GlobalFlags
	TfRequireDestTag  bool `json:"tfRequireDestTag,omitempty"`
	TfOptionalDestTag bool `json:"tfOptionalDestTag,omitempty"`
	TfRequireAuth     bool `json:"tfRequireAuth,omitempty"`
	TfOptionalAuth    bool `json:"tfOptionalAuth,omitempty"`
	TfDisallowXRP     bool `json:"tfDisallowXRP,omitempty"`
	TfAllowXRP        bool `json:"tfAllowXRP,omitempty"`
}

// An AccountSet transaction modifies the properties of an account in the XRP
// Ledger
//
// TransactionType: 'AccountSet'
type TransactionAccountSet struct {
	BaseTransaction
	Flags         int64
	ClearFlag     int64
	Domain        string
	EmailHash     string
	MessageKey    string
	SetFlag       int64
	TransferRate  int64
	TickSize      int64
	NFTokenMinter string
}

const (
	MIN_TICK_SIZE = 3
	MAX_TICK_SIZE = 15
)

// Cancels an unredeemed Check, removing it from the ledger without sending any
// money. The source or the destination of the check can cancel a Check at any
// time using this transaction type. If the Check has expired, any address can
// cancel it.
//
// TransactionType: 'CheckCancel'
type TransactionCheckCancel struct {
	BaseTransaction
	CheckID string
}

// Attempts to redeem a Check object in the ledger to receive up to the amount
// authorized by the corresponding CheckCreate transaction. Only the Destination
// address of a Check can cash it with a CheckCash transaction.
//
// TransactionType: 'CheckCash'
type TransactionCheckCash struct {
	BaseTransaction
	CheckID    string
	Amount     Amount
	DeliverMin Amount
}

// Create a Check object in the ledger, which is a deferred payment that can be
// cashed by its intended destination. The sender of this transaction is the
// sender of the Check.
//
// TransactionType: 'CheckCreate'
type TransactionCheckCreate struct {
	BaseTransaction
	Destination    string
	SendMax        Amount
	DestinationTag int64
	Expiration     int64
	InvoiceID      string
}

// A DepositPreauth transaction gives another account pre-approval to deliver
// payments to the sender of this transaction. This is only useful if the sender
// of this transaction is using (or plans to use) Deposit Authorization.
//
// TransactionType: 'DepositPreauth'
type TransactionDepositPreauth struct {
	BaseTransaction
	Authorize   string
	Unauthorize string
}

// Return escrowed XRP to the sender.
//
// TransactionType: 'EscrowCancel'
type TransactionEscrowCancel struct {
	BaseTransaction
	Owner         string
	OfferSequence int64
}

// Sequester XRP until the escrow process either finishes or is canceled.
//
// TransactionType: 'EscrowCreate'
type TransactionEscrowCreate struct {
	BaseTransaction
	Amount         Amount
	Destination    string
	CancelAfter    int64
	FinishAfter    int64
	Condition      string
	DestinationTag int64
}

// Deliver XRP from a held payment to the recipient.
//
// TransactionType: 'EscrowFinish'
type TransactionEscrowFinish struct {
	BaseTransaction
	Owner         string
	OfferSequence int64
	Condition     string
	Fulfillment   string
}

// An OfferCancel transaction removes an Offer object from the XRP Ledger.
//
// TransactionType: 'OfferCancel'
type TransactionOfferCancel struct {
	BaseTransaction
	OfferSequence int64
}

// An OfferCreate transaction is effectively a limit order . It defines an
// intent to exchange currencies, and creates an Offer object if not completely.
// Fulfilled when placed. Offers can be partially fulfilled.
//
// TransactionType: 'OfferCreate'
type TransactionOfferCreate struct {
	BaseTransaction
	Flags         int64
	Expiration    int64
	OfferSequence int64
	TakerGets     Amount
	TakerPays     Amount
}

type OfferCreateFlags struct {
	GlobalFlags
	TfPassive           bool `json:"tfPassive,omitempty"`
	TfImmediateOrCancel bool `json:"tfImmediateOrCancel,omitempty"`
	TfFillOrKill        bool `json:"tfFillOrKill,omitempty"`
	TfSell              bool `json:"tfSell,omitempty"`
}

// Claim XRP from a payment channel, adjust the payment channel's expiration,
// or both.
//
// TransactionType: 'PaymentChannelClaim'
type TransactionPaymentChannelClaim struct {
	BaseTransaction
	Flags     int64
	Channel   string
	Balance   string
	Amount    string
	Signature string
	PublicKey string
}

type PaymentChannelClaimFlags struct {
	GlobalFlags
	TfRenew bool `json:"tfRenew,omitempty"`
	TfClose bool `json:"tfClose,omitempty"`
}

// Create a unidirectional channel and fund it with XRP. The address sending
// this transaction becomes the "source address" of the payment channel.
//
// TransactionType: 'PaymentChannelCreate'
type TransactionPaymentChannelCreate struct {
	BaseTransaction
	Amount         string
	Destination    string
	SettleDelay    int64
	PublicKey      string
	CancelAfter    int64
	DestinationTag int64
}

// Add additional XRP to an open payment channel, and optionally update the
// expiration time of the channel. Only the source address of the channel can
// use this transaction.
//
// TransactionType: 'PaymentChannelFund'
type TransactionPaymentChannelFund struct {
	BaseTransaction
	Channel    string
	Amount     string
	Expiration int64
}

// A SetRegularKey transaction assigns, changes, or removes the regular key
// pair associated with an account.
//
// TransactionType: 'SetRegularKey'
type TransactionSetRegularKey struct {
	BaseTransaction
	RegularKey string
}

// The SignerListSet transaction creates, replaces, or removes a list of
// signers that can be used to multi-sign a transaction.
//
// TransactionType: 'SignerListSet'
type TransactionSignerListSet struct {
	BaseTransaction
	SignerQuorum  int64
	SignerEntries []SignerEntry
}

// A TicketCreate transaction sets aside one or more sequence numbers as
// Tickets.
//
// TransactionType: 'TicketCreate'
type TransactionTicketCreate struct {
	BaseTransaction
	TicketCount int64
}

const MAX_TICKETS = 250

// Create or modify a trust line linking two accounts.
//
// TransactionType: 'TrustSet'
type TransactionTrustSet struct {
	BaseTransaction
	LimitAmount IssuedCurrencyAmount
	QualityIn   int64
	QualityOut  int64
	Flags       int64
}

type TrustSetFlags struct {
	GlobalFlags
	TfSetfAuth      bool `json:"tfSetfAuth,omitempty"`
	TfSetNoRipple   bool `json:"tfSetNoRipple,omitempty"`
	TfClearNoRipple bool `json:"tfClearNoRipple,omitempty"`
	TfSetFreeze     bool `json:"tfSetFreeze,omitempty"`
	TfClearFreeze   bool `json:"tfClearFreeze,omitempty"`
}

type Transaction struct {
	TransactionAccountDelete
	TransactionAccountSet
	TransactionCheckCancel
	TransactionCheckCash
	TransactionCheckCreate
	TransactionDepositPreauth
	TransactionEscrowCancel
	TransactionEscrowCreate
	TransactionEscrowFinish
	TransactionNFTokenAcceptOffer
	TransactionNFTokenBurn
	TransactionNFTokenCancelOffer
	TransactionNFTokenCreateOffer
	TransactionNFTokenMint
	TransactionOfferCancel
	TransactionOfferCreate
	TransactionPayment
	TransactionPaymentChannelClaim
	TransactionPaymentChannelCreate
	TransactionPaymentChannelFund
	TransactionSetRegularKey
	TransactionSignerListSet
	TransactionTicketCreate
	TransactionTrustSet
}
