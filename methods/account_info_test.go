package methods

import (
	"testing"

	"github.com/xrpscan/xrpl-go/models"
)

func TestAccountInfoRequest_Fields(t *testing.T) {
	request := AccountInfoRequest{
		BaseRequest: models.BaseRequest{
			Command: "account_info",
		},
		Account:     "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
		LedgerIndex: "validated",
		LedgerHash:  "abcd1234",
		Queue:       true,
		SignerLists: false,
		Strict:      true,
	}

	if request.Account != "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn" {
		t.Errorf("Expected Account to be 'rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn', got %s", request.Account)
	}

	if request.LedgerIndex != "validated" {
		t.Errorf("Expected LedgerIndex to be 'validated', got %s", request.LedgerIndex)
	}

	if request.LedgerHash != "abcd1234" {
		t.Errorf("Expected LedgerHash to be 'abcd1234', got %s", request.LedgerHash)
	}

	if !request.Queue {
		t.Error("Expected Queue to be true")
	}

	if request.SignerLists {
		t.Error("Expected SignerLists to be false")
	}

	if !request.Strict {
		t.Error("Expected Strict to be true")
	}

	if request.Command != "account_info" {
		t.Errorf("Expected Command to be 'account_info', got %s", request.Command)
	}
}

func TestQueueTransaction_Fields(t *testing.T) {
	queueTx := QueueTransaction{
		AuthChange:    true,
		Fee:           "12",
		FeeLevel:      "256",
		MaxSpendDrops: "10000000",
		Seq:           5,
	}

	if !queueTx.AuthChange {
		t.Error("Expected AuthChange to be true")
	}

	if queueTx.Fee != "12" {
		t.Errorf("Expected Fee to be '12', got %s", queueTx.Fee)
	}

	if queueTx.FeeLevel != "256" {
		t.Errorf("Expected FeeLevel to be '256', got %s", queueTx.FeeLevel)
	}

	if queueTx.MaxSpendDrops != "10000000" {
		t.Errorf("Expected MaxSpendDrops to be '10000000', got %s", queueTx.MaxSpendDrops)
	}

	if queueTx.Seq != 5 {
		t.Errorf("Expected Seq to be 5, got %d", queueTx.Seq)
	}
}

func TestQueueData_Fields(t *testing.T) {
	queueData := QueueData{
		TxnCount:           3,
		AuthChangeQueued:   true,
		LowestSequence:     1,
		HighestSequence:    5,
		MaxSpendDropsTotal: "30000000",
		Transactions: []QueueTransaction{
			{
				AuthChange:    false,
				Fee:           "10",
				FeeLevel:      "128",
				MaxSpendDrops: "5000000",
				Seq:           1,
			},
			{
				AuthChange:    true,
				Fee:           "15",
				FeeLevel:      "512",
				MaxSpendDrops: "10000000",
				Seq:           3,
			},
		},
	}

	if queueData.TxnCount != 3 {
		t.Errorf("Expected TxnCount to be 3, got %d", queueData.TxnCount)
	}

	if !queueData.AuthChangeQueued {
		t.Error("Expected AuthChangeQueued to be true")
	}

	if queueData.LowestSequence != 1 {
		t.Errorf("Expected LowestSequence to be 1, got %d", queueData.LowestSequence)
	}

	if queueData.HighestSequence != 5 {
		t.Errorf("Expected HighestSequence to be 5, got %d", queueData.HighestSequence)
	}

	if queueData.MaxSpendDropsTotal != "30000000" {
		t.Errorf("Expected MaxSpendDropsTotal to be '30000000', got %s", queueData.MaxSpendDropsTotal)
	}

	if len(queueData.Transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(queueData.Transactions))
	}

	if queueData.Transactions[0].Fee != "10" {
		t.Errorf("Expected first transaction fee to be '10', got %s", queueData.Transactions[0].Fee)
	}

	if queueData.Transactions[1].AuthChange != true {
		t.Error("Expected second transaction AuthChange to be true")
	}
}

func TestAccountInfoResult_Fields(t *testing.T) {
	result := AccountInfoResult{
		AccountData:        "test_account_data",
		SignerLists:        "test_signer_lists",
		LedgerCurrentIndex: 12345,
		LedgerIndex:        12340,
		QueueData: QueueData{
			TxnCount:        1,
			LowestSequence:  1,
			HighestSequence: 1,
		},
		Validated: true,
	}

	if result.AccountData != "test_account_data" {
		t.Errorf("Expected AccountData to be 'test_account_data', got %s", result.AccountData)
	}

	if result.SignerLists != "test_signer_lists" {
		t.Errorf("Expected SignerLists to be 'test_signer_lists', got %s", result.SignerLists)
	}

	if result.LedgerCurrentIndex != 12345 {
		t.Errorf("Expected LedgerCurrentIndex to be 12345, got %d", result.LedgerCurrentIndex)
	}

	if result.LedgerIndex != 12340 {
		t.Errorf("Expected LedgerIndex to be 12340, got %d", result.LedgerIndex)
	}

	if !result.Validated {
		t.Error("Expected Validated to be true")
	}

	if result.QueueData.TxnCount != 1 {
		t.Errorf("Expected QueueData.TxnCount to be 1, got %d", result.QueueData.TxnCount)
	}
}

func TestAccountInfoResponse_Fields(t *testing.T) {
	response := AccountInfoResponse{
		BaseResponse: models.BaseResponse{
			Status: "success",
			Type:   "response",
		},
		Result: AccountInfoResult{
			LedgerCurrentIndex: 12345,
			Validated:          true,
		},
	}

	if response.Status != "success" {
		t.Errorf("Expected Status to be 'success', got %s", response.Status)
	}

	if response.Type != "response" {
		t.Errorf("Expected Type to be 'response', got %s", response.Type)
	}

	if response.Result.LedgerCurrentIndex != 12345 {
		t.Errorf("Expected Result.LedgerCurrentIndex to be 12345, got %d", response.Result.LedgerCurrentIndex)
	}

	if !response.Result.Validated {
		t.Error("Expected Result.Validated to be true")
	}
}