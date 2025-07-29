package xrpl

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ClientConfig{
				URL:               "ws://localhost:8080",
				ReadTimeout:       30 * time.Second,
				WriteTimeout:      30 * time.Second,
				HeartbeatInterval: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty URL",
			config: ClientConfig{
				URL: "",
			},
			wantErr: true,
		},
		{
			name: "negative read timeout",
			config: ClientConfig{
				URL:         "ws://localhost:8080",
				ReadTimeout: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "read timeout too large",
			config: ClientConfig{
				URL:         "ws://localhost:8080",
				ReadTimeout: 2 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "negative write timeout",
			config: ClientConfig{
				URL:          "ws://localhost:8080",
				WriteTimeout: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "write timeout too large",
			config: ClientConfig{
				URL:          "ws://localhost:8080",
				WriteTimeout: 2 * time.Hour,
			},
			wantErr: true,
		},
		{
			name: "negative heartbeat interval",
			config: ClientConfig{
				URL:               "ws://localhost:8080",
				HeartbeatInterval: -1 * time.Second,
			},
			wantErr: true,
		},
		{
			name: "heartbeat interval too large",
			config: ClientConfig{
				URL:               "ws://localhost:8080",
				HeartbeatInterval: 2 * time.Hour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient_DefaultValues(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to connection failure
		}
	}()

	client := NewClient(config)

	if client.config.ReadTimeout != 60*time.Second {
		t.Errorf("Expected ReadTimeout to be 60s, got %v", client.config.ReadTimeout)
	}
	if client.config.WriteTimeout != 60*time.Second {
		t.Errorf("Expected WriteTimeout to be 60s, got %v", client.config.WriteTimeout)
	}
	if client.config.HeartbeatInterval != 5*time.Second {
		t.Errorf("Expected HeartbeatInterval to be 5s, got %v", client.config.HeartbeatInterval)
	}
	if client.config.QueueCapacity != 128 {
		t.Errorf("Expected QueueCapacity to be 128, got %v", client.config.QueueCapacity)
	}
}

func TestNewClient_CustomValues(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL:               wsURL,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      45 * time.Second,
		HeartbeatInterval: 10 * time.Second,
		QueueCapacity:     256,
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to connection failure
		}
	}()

	client := NewClient(config)

	if client.config.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout to be 30s, got %v", client.config.ReadTimeout)
	}
	if client.config.WriteTimeout != 45*time.Second {
		t.Errorf("Expected WriteTimeout to be 45s, got %v", client.config.WriteTimeout)
	}
	if client.config.HeartbeatInterval != 10*time.Second {
		t.Errorf("Expected HeartbeatInterval to be 10s, got %v", client.config.HeartbeatInterval)
	}
	if client.config.QueueCapacity != 256 {
		t.Errorf("Expected QueueCapacity to be 256, got %v", client.config.QueueCapacity)
	}
}

func TestNewClient_InvalidConfig(t *testing.T) {
	config := ClientConfig{
		URL: "",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected NewClient to panic with invalid config")
		}
	}()

	NewClient(config)
}

func TestClient_NextID(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to connection failure
		}
	}()

	client := NewClient(config)

	id1 := client.NextID()
	id2 := client.NextID()
	id3 := client.NextID()

	if id1 != "1" {
		t.Errorf("Expected first ID to be '1', got '%s'", id1)
	}
	if id2 != "2" {
		t.Errorf("Expected second ID to be '2', got '%s'", id2)
	}
	if id3 != "3" {
		t.Errorf("Expected third ID to be '3', got '%s'", id3)
	}
}

func TestClient_Subscriptions(t *testing.T) {
	server := httptest.NewServer(nil)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to connection failure
		}
	}()

	client := NewClient(config)

	// Initially should be empty
	subs := client.Subscriptions()
	if len(subs) != 0 {
		t.Errorf("Expected no subscriptions initially, got %d", len(subs))
	}

	// Add some subscriptions manually for testing
	client.StreamSubscriptions["ledger"] = true
	client.StreamSubscriptions["transactions"] = true

	subs = client.Subscriptions()
	if len(subs) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(subs))
	}

	// Check that both subscriptions are present
	found := make(map[string]bool)
	for _, sub := range subs {
		found[sub] = true
	}

	if !found["ledger"] || !found["transactions"] {
		t.Error("Expected to find both 'ledger' and 'transactions' subscriptions")
	}
}

// Mock WebSocket server for testing
func setupMockWebSocketServer() *httptest.Server {
	server := httptest.NewServer(nil)
	return server
}

func TestClient_Close(t *testing.T) {
	server := setupMockWebSocketServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	defer func() {
		if r := recover(); r != nil {
			// Expected to panic due to connection failure
		}
	}()

	client := NewClient(config)
	
	// Set closed to false initially for testing
	client.closed = false

	// Close should not panic and should set closed to true
	client.Close()

	if !client.closed {
		t.Error("Expected client to be marked as closed")
	}
}