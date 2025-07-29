package xrpl

import (
	"encoding/json"
	"testing"
	"time"
)

func TestClient_resolveStream(t *testing.T) {
	config := ClientConfig{
		URL:           "ws://localhost:8080",
		QueueCapacity: 10,
	}

	// Create client without connecting
	client := &Client{
		config:              config,
		heartbeatDone:       make(chan bool),
		handlerDone:         make(chan bool),
		StreamLedger:        make(chan []byte, config.QueueCapacity),
		StreamTransaction:   make(chan []byte, config.QueueCapacity),
		StreamValidation:    make(chan []byte, config.QueueCapacity),
		StreamManifest:      make(chan []byte, config.QueueCapacity),
		StreamPeerStatus:    make(chan []byte, config.QueueCapacity),
		StreamConsensus:     make(chan []byte, config.QueueCapacity),
		StreamPathFind:      make(chan []byte, config.QueueCapacity),
		StreamServer:        make(chan []byte, config.QueueCapacity),
		StreamDefault:       make(chan []byte, config.QueueCapacity),
		StreamSubscriptions: make(map[string]bool),
		requestQueue:        make(map[string](chan<- BaseResponse)),
		nextId:              0,
	}

	tests := []struct {
		name         string
		message      []byte
		expectedChan string
	}{
		{
			name:         "ledger stream",
			message:      []byte(`{"type": "ledgerClosed"}`),
			expectedChan: "ledger",
		},
		{
			name:         "transaction stream",
			message:      []byte(`{"type": "transaction"}`),
			expectedChan: "transaction",
		},
		{
			name:         "validation stream",
			message:      []byte(`{"type": "validationReceived"}`),
			expectedChan: "validation",
		},
		{
			name:         "manifest stream",
			message:      []byte(`{"type": "manifestReceived"}`),
			expectedChan: "manifest",
		},
		{
			name:         "peer status stream",
			message:      []byte(`{"type": "peerStatusChange"}`),
			expectedChan: "peerStatus",
		},
		{
			name:         "consensus stream",
			message:      []byte(`{"type": "consensusPhase"}`),
			expectedChan: "consensus",
		},
		{
			name:         "path find stream",
			message:      []byte(`{"type": "path_find"}`),
			expectedChan: "pathFind",
		},
		{
			name:         "server stream",
			message:      []byte(`{"type": "serverStatus"}`),
			expectedChan: "server",
		},
		{
			name:         "default stream",
			message:      []byte(`{"type": "unknown"}`),
			expectedChan: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear channels
			select {
			case <-client.StreamLedger:
			default:
			}
			select {
			case <-client.StreamTransaction:
			default:
			}
			select {
			case <-client.StreamValidation:
			default:
			}
			select {
			case <-client.StreamManifest:
			default:
			}
			select {
			case <-client.StreamPeerStatus:
			default:
			}
			select {
			case <-client.StreamConsensus:
			default:
			}
			select {
			case <-client.StreamPathFind:
			default:
			}
			select {
			case <-client.StreamServer:
			default:
			}
			select {
			case <-client.StreamDefault:
			default:
			}

			// Call resolveStream
			client.resolveStream(tt.message)

			// Check that message was sent to correct channel
			timeout := time.After(100 * time.Millisecond)
			
			switch tt.expectedChan {
			case "ledger":
				select {
				case msg := <-client.StreamLedger:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamLedger")
				}
			case "transaction":
				select {
				case msg := <-client.StreamTransaction:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamTransaction")
				}
			case "validation":
				select {
				case msg := <-client.StreamValidation:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamValidation")
				}
			case "manifest":
				select {
				case msg := <-client.StreamManifest:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamManifest")
				}
			case "peerStatus":
				select {
				case msg := <-client.StreamPeerStatus:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamPeerStatus")
				}
			case "consensus":
				select {
				case msg := <-client.StreamConsensus:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamConsensus")
				}
			case "pathFind":
				select {
				case msg := <-client.StreamPathFind:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamPathFind")
				}
			case "server":
				select {
				case msg := <-client.StreamServer:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamServer")
				}
			case "default":
				select {
				case msg := <-client.StreamDefault:
					if string(msg) != string(tt.message) {
						t.Errorf("Expected message %s, got %s", string(tt.message), string(msg))
					}
				case <-timeout:
					t.Error("Timeout waiting for message in StreamDefault")
				}
			}
		})
	}
}

func TestClient_resolveStream_Response(t *testing.T) {
	config := ClientConfig{
		URL:           "ws://localhost:8080",
		QueueCapacity: 10,
	}

	client := &Client{
		config:              config,
		heartbeatDone:       make(chan bool),
		handlerDone:         make(chan bool),
		StreamLedger:        make(chan []byte, config.QueueCapacity),
		StreamTransaction:   make(chan []byte, config.QueueCapacity),
		StreamValidation:    make(chan []byte, config.QueueCapacity),
		StreamManifest:      make(chan []byte, config.QueueCapacity),
		StreamPeerStatus:    make(chan []byte, config.QueueCapacity),
		StreamConsensus:     make(chan []byte, config.QueueCapacity),
		StreamPathFind:      make(chan []byte, config.QueueCapacity),
		StreamServer:        make(chan []byte, config.QueueCapacity),
		StreamDefault:       make(chan []byte, config.QueueCapacity),
		StreamSubscriptions: make(map[string]bool),
		requestQueue:        make(map[string](chan<- BaseResponse)),
		nextId:              0,
	}

	// Create a channel to receive the response
	responseChan := make(chan BaseResponse, 1)
	requestId := "123"
	
	client.requestQueue[requestId] = responseChan

	// Create a response message
	responseMessage := map[string]interface{}{
		"type":   "response",
		"id":     requestId,
		"status": "success",
	}
	message, _ := json.Marshal(responseMessage)

	// Call resolveStream
	client.resolveStream(message)

	// Check that response was sent to the correct channel
	timeout := time.After(100 * time.Millisecond)
	select {
	case response := <-responseChan:
		if response["id"] != requestId {
			t.Errorf("Expected response id %s, got %v", requestId, response["id"])
		}
		if response["status"] != "success" {
			t.Errorf("Expected response status 'success', got %v", response["status"])
		}
	case <-timeout:
		t.Error("Timeout waiting for response")
	}

	// Check that request was removed from queue
	if _, exists := client.requestQueue[requestId]; exists {
		t.Error("Expected request to be removed from queue")
	}
}

func TestClient_handlePong(t *testing.T) {
	// This test is skipped because handlePong requires a valid websocket connection
	// and testing it properly would require setting up a full websocket server
	// The function primarily updates connection deadlines which is an internal implementation detail
	t.Skip("handlePong requires valid websocket connection - integration test needed")
}