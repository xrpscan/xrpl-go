package xrpl

import (
	"testing"
)

func TestStreamResponseType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ledger stream type",
			input:    StreamTypeLedger,
			expected: "ledgerClosed",
		},
		{
			name:     "transaction stream type",
			input:    StreamTypeTransaction,
			expected: "transaction",
		},
		{
			name:     "transactions proposed stream type",
			input:    StreamTypeTransactionsProposed,
			expected: "transaction",
		},
		{
			name:     "validations stream type",
			input:    StreamTypeValidations,
			expected: "validationReceived",
		},
		{
			name:     "manifests stream type",
			input:    StreamTypeManifests,
			expected: "manifestReceived",
		},
		{
			name:     "peer status stream type",
			input:    StreamTypePeerStatus,
			expected: "peerStatusChange",
		},
		{
			name:     "consensus stream type",
			input:    StreamTypeConsensus,
			expected: "consensusPhase",
		},
		{
			name:     "path find stream type",
			input:    StreamTypePathFind,
			expected: "path_find",
		},
		{
			name:     "server stream type",
			input:    StreamTypeServer,
			expected: "serverStatus",
		},
		{
			name:     "response stream type",
			input:    StreamTypeResponse,
			expected: "response",
		},
		{
			name:     "unknown stream type",
			input:    "unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StreamResponseType(tt.input)
			if result != tt.expected {
				t.Errorf("StreamResponseType(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStreamTypeConstants(t *testing.T) {
	// Test that all stream type constants are defined and not empty
	streamTypes := []string{
		StreamTypeLedger,
		StreamTypeTransaction,
		StreamTypeTransactionsProposed,
		StreamTypeValidations,
		StreamTypeManifests,
		StreamTypePeerStatus,
		StreamTypeConsensus,
		StreamTypePathFind,
		StreamTypeServer,
		StreamTypeResponse,
	}

	expectedValues := []string{
		"ledger",
		"transactions",
		"transactions_proposed",
		"validations",
		"manifests",
		"peer_status",
		"consensus",
		"path_find",
		"server",
		"response",
	}

	if len(streamTypes) != len(expectedValues) {
		t.Errorf("Number of stream types (%d) doesn't match expected (%d)", len(streamTypes), len(expectedValues))
	}

	for i, streamType := range streamTypes {
		if streamType == "" {
			t.Errorf("Stream type at index %d should not be empty", i)
		}
		if i < len(expectedValues) && streamType != expectedValues[i] {
			t.Errorf("Stream type at index %d: got %s, want %s", i, streamType, expectedValues[i])
		}
	}
}