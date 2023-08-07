package xrpl

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func (c *Client) Subscribe(streams []string) (BaseResponse, error) {
	req := BaseRequest{
		"command": "subscribe",
		"streams": streams,
	}
	res, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	for _, stream := range streams {
		c.StreamSubscriptions[stream] = true
	}
	c.mutex.Unlock()

	return res, nil
}

func (c *Client) Unsubscribe(streams []string) (BaseResponse, error) {
	req := BaseRequest{
		"command": "unsubscribe",
		"streams": streams,
	}
	res, err := c.Request(req)
	if err != nil {
		return nil, err
	}

	c.mutex.Lock()
	for _, stream := range streams {
		delete(c.StreamSubscriptions, stream)
	}
	c.mutex.Unlock()

	return res, nil
}

// Send a websocket request. This method takes a BaseRequest object and automatically adds
// incremental request ID to it.
//
// Example usage:
//
//	req := BaseRequest{
//		"command": "account_info",
//		"account": "rG1QQv2nh2gr7RCZ1P8YYcBUKCCN633jCn",
//		"ledger_index": "current",
//	}
//
//	err := client.Request(req, func(){})
func (c *Client) Request(req BaseRequest) (BaseResponse, error) {
	requestId := c.NextID()
	req["id"] = requestId
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan BaseResponse, 1)

	c.mutex.Lock()
	c.requestQueue[requestId] = ch
	err = c.connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return nil, err
	}
	c.mutex.Unlock()

	res := <-ch
	return res, nil
}
