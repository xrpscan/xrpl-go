package xrpl

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		
		// Parse the incoming request
		var req BaseRequest
		if err := json.Unmarshal(message, &req); err != nil {
			continue
		}
		
		// Create a mock response
		response := BaseResponse{
			"id":     req["id"],
			"status": "success",
			"type":   "response",
			"result": map[string]interface{}{
				"command": req["command"],
			},
		}
		
		responseBytes, _ := json.Marshal(response)
		err = c.WriteMessage(mt, responseBytes)
		if err != nil {
			break
		}
	}
}

func setupTestServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(echoHandler))
	return server
}

func TestClient_Subscribe(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	client := NewClient(config)
	defer client.Close()

	// Give some time for connection to establish  
	time.Sleep(100 * time.Millisecond)

	streams := []string{"ledger", "transactions"}
	response, err := client.Subscribe(streams)

	if err != nil {
		t.Errorf("Subscribe failed: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}

	// Check that subscriptions were added
	subs := client.Subscriptions()
	if len(subs) != 2 {
		t.Errorf("Expected 2 subscriptions, got %d", len(subs))
	}

	found := make(map[string]bool)
	for _, sub := range subs {
		found[sub] = true
	}

	if !found["ledger"] || !found["transactions"] {
		t.Error("Expected to find both 'ledger' and 'transactions' subscriptions")
	}
}

func TestClient_Unsubscribe(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	client := NewClient(config)
	defer client.Close()

	// Give some time for connection to establish
	time.Sleep(100 * time.Millisecond)

	// First subscribe to some streams
	streams := []string{"ledger", "transactions", "validations"}
	_, err := client.Subscribe(streams)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Verify subscriptions were added
	subs := client.Subscriptions()
	if len(subs) != 3 {
		t.Errorf("Expected 3 subscriptions after subscribe, got %d", len(subs))
	}

	// Now unsubscribe from some streams
	toUnsubscribe := []string{"ledger", "validations"}
	response, err := client.Unsubscribe(toUnsubscribe)

	if err != nil {
		t.Errorf("Unsubscribe failed: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}

	// Check that subscriptions were removed
	subs = client.Subscriptions()
	if len(subs) != 1 {
		t.Errorf("Expected 1 subscription after unsubscribe, got %d", len(subs))
	}

	if subs[0] != "transactions" {
		t.Errorf("Expected remaining subscription to be 'transactions', got '%s'", subs[0])
	}
}

func TestClient_Request(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	client := NewClient(config)
	defer client.Close()

	// Give some time for connection to establish
	time.Sleep(100 * time.Millisecond)

	req := BaseRequest{
		"command":      "account_info",
		"account":      "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
		"ledger_index": "current",
	}

	response, err := client.Request(req)

	if err != nil {
		t.Errorf("Request failed: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}

	// Check that the response contains expected fields
	if status, ok := response["status"]; !ok || status != "success" {
		t.Error("Expected response status to be 'success'")
	}

	if responseType, ok := response["type"]; !ok || responseType != "response" {
		t.Error("Expected response type to be 'response'")
	}

	// Check that request ID was added
	if _, ok := response["id"]; !ok {
		t.Error("Expected response to contain an id field")
	}
}

func TestClient_Request_WithID(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	config := ClientConfig{
		URL: wsURL,
	}

	client := NewClient(config)
	defer client.Close()

	// Give some time for connection to establish
	time.Sleep(100 * time.Millisecond)

	// Test that auto-generated ID is used when none provided
	req := BaseRequest{
		"command": "ping",
	}

	response, err := client.Request(req)

	if err != nil {
		t.Errorf("Request failed: %v", err)
	}

	if response == nil {
		t.Error("Expected non-nil response")
	}

	// The ID should be auto-generated (should be "1" as it's the first request)
	if id, ok := response["id"]; !ok || id != "1" {
		t.Errorf("Expected response id to be '1', got %v", id)
	}
}